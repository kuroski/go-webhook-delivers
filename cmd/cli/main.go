package main

import (
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"github.com/cenkalti/backoff/v5"
	"github.com/kuroski/go-webhook-deliveries/internal/logger"
	"github.com/kuroski/go-webhook-deliveries/internal/model"
	"github.com/tmaxmax/go-sse"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	source := flag.String("source", "", "The source server that we will listen for events")
	target := flag.String("target", "", "The target server that the request will be forwarded to")
	flag.Parse()

	log := logger.NewLogger()

	if *source == "" || *target == "" {
		log.Error("Required flags --source and --target must be set")
		flag.Usage()
		os.Exit(1)
		return
	}

	sourceURL, err := url.Parse(*source)
	if err != nil {
		log.Error("--source argument must be a valid url", "source", *source)
	}

	targetURL, err := url.Parse(*target)
	if err != nil {
		log.Error("--target argument must be a valid url", "target", *target)
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	client := sse.DefaultClient

	r, _ := http.NewRequestWithContext(ctx, http.MethodGet, sourceURL.String(), http.NoBody)
	conn := client.NewConnection(r)

	conn.SubscribeMessages(func(event sse.Event) {
		var request model.Request
		if err := json.Unmarshal([]byte(event.Data), &request); err != nil {
			log.Error(err.Error())
			return
		}

		log.Info(
			"Message received",
			"headers",
			request.Headers,
			"query",
			request.Query,
			"timestamp",
			request.Timestamp,
			"body",
			request.Body[:100],
		)

		targetURL.RawQuery = request.Query
		req, err := http.NewRequest("POST", targetURL.String(), bytes.NewBuffer(request.Body))
		if err != nil {
			log.Error(err.Error())
			return
		}

		req.Header = request.Headers.Clone()
		res, err := client.HTTPClient.Do(req)
		if err != nil {
			log.Error(err.Error())
			return
		}

		log.Info("Request delivered successfully", "response", res)
	})

	_, err = backoff.Retry(ctx, func() (any, error) {
		log.Info("connecting sse, waiting for events from", "source", sourceURL.String(), "target", targetURL.String())
		if err := conn.Connect(); err != nil {
			log.Info(err.Error())
			return nil, err
		}

		return nil, nil
	}, backoff.WithBackOff(backoff.NewExponentialBackOff()))
	if err != nil {
		log.Error(err.Error())
	}
}
