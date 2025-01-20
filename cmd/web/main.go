package main

import (
	"context"
	"encoding/base64"
	"flag"
	"github.com/bradleyfalzon/ghinstallation/v2"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/google/go-github/v67/github"
	_ "github.com/joho/godotenv/autoload"
	"github.com/kuroski/go-webhook-deliveries/internal/logger"
	"github.com/kuroski/go-webhook-deliveries/internal/server"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
)

func main() {
	addr := flag.String("addr", ":3000", "HTTP network address")
	flag.Parse()

	log := logger.NewLogger()

	privateKey, _ := base64.StdEncoding.DecodeString(os.Getenv("GITHUB_APP_PRIVATE_KEY"))
	appID, _ := strconv.ParseInt(os.Getenv("GITHUB_APP_ID"), 10, 64)
	installationID, _ := strconv.ParseInt(os.Getenv("GITHUB_APP_INSTALLATION_ID"), 10, 64)

	tr := http.DefaultTransport
	itr, err := ghinstallation.New(tr, appID, installationID, privateKey)
	if err != nil {
		panic(err)
	}

	githubClient := github.NewClient(&http.Client{Transport: itr})

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	opts := []bot.Option{
		bot.WithCallbackQueryDataHandler("cancel_", bot.MatchTypePrefix, func(ctx context.Context, b *bot.Bot, update *models.Update) {
			data := strings.Split(update.CallbackQuery.Data, "_")
			if len(data) == 3 {
				r := strings.Split(data[2], "/")
				id, err := strconv.ParseInt(data[1], 10, 64)
				if err != nil {
					log.Error(err.Error())
					return
				}

				log.Info("----- CANCEL CLICKED", data, r, id)
				if _, err := githubClient.Actions.CancelWorkflowRunByID(ctx, r[0], r[1], id); err != nil {
					log.Error(err.Error())
					return
				}

				log.Info("----- WORKFLOW RUN CANCELLED", data)
			}
		}),
	}

	b, err := bot.New(os.Getenv("TELEGRAM_BOT_TOKEN"), opts...)
	if err != nil {
		panic(err)
	}

	go b.Start(ctx)

	go func() {
		srv := server.NewServer(*addr, ctx, b)
		srv.Start()
	}()

	<-ctx.Done()
	log.Info("Shutting down gracefully")
}
