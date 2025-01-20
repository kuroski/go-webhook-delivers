package server

import (
	"encoding/json"
	"fmt"
	"github.com/google/go-github/v67/github"
	"github.com/kuroski/go-webhook-deliveries/internal/model"
	"github.com/tmaxmax/go-sse"
	"io"
	"net/http"
	"os"
	"time"
)

func (srv *Server) upHandler(w http.ResponseWriter, r *http.Request) {
	if _, err := fmt.Fprintf(w, "Hello World!!"); err != nil {
		srv.serverError(w, r, err)
	}
}

func (srv *Server) payloadUrlHandler(w http.ResponseWriter, r *http.Request) {
	srv.sseServer.ServeHTTP(w, r)
}

func (srv *Server) forwardRequestHandler(w http.ResponseWriter, r *http.Request) {
	channel := r.PathValue("channel")
	body, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		srv.serverError(w, r, err)
		return
	}

	req := model.Request{
		Headers:   r.Header,
		Body:      body,
		Query:     r.URL.RawQuery,
		Timestamp: time.Now(),
	}

	payload, err := json.Marshal(req)
	if err != nil {
		srv.serverError(w, r, err)
		return
	}

	message := &sse.Message{}
	message.AppendData(string(payload))
	srv.log.Info("Hook received, forwarding the request", "req", req, "body", body)

	if err = srv.sseServer.Publish(message, channel); err != nil {
		srv.serverError(w, r, err)
	}
	if _, err := fmt.Fprintf(w, "Message forwarded"); err != nil {
		srv.serverError(w, r, err)
	}
}

func (srv *Server) webhookHandler(w http.ResponseWriter, r *http.Request) {
	srv.log.Info("-- Received request", "method", r.Method, "url", r.URL.String(), "headers", r.Header)

	payload, err := github.ValidatePayload(r, nil)
	if err != nil {
		srv.serverError(w, r, err)
		return
	}
	event, err := github.ParseWebHook(github.WebHookType(r), payload)
	if err != nil {
		srv.serverError(w, r, err)
		return
	}

	if event, ok := event.(*github.WorkflowRunEvent); ok {
		err := srv.workflowManager.HandleProgress(srv.ctx, srv.bot, *event.WorkflowRun, os.Getenv("TELEGRAM_CHAT_ID"))
		if err != nil {
			srv.log.Error(err.Error())
		}
	}
}
