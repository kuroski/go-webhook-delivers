package server

import (
	"context"
	"github.com/go-telegram/bot"
	"github.com/kuroski/go-webhook-deliveries/internal/logger"
	"github.com/kuroski/go-webhook-deliveries/internal/workflowmanager"
	"github.com/tmaxmax/go-sse"
	"log/slog"
	"net/http"
	"os"
)

type Server struct {
	sseServer       *sse.Server
	server          http.Server
	log             *slog.Logger
	ctx             context.Context
	bot             *bot.Bot
	workflowManager *workflowmanager.WorkflowManager
}

func NewServer(addr string, ctx context.Context, bot *bot.Bot) *Server {
	server := &Server{
		ctx: ctx,
		bot: bot,
		log: logger.NewLogger(),
		sseServer: &sse.Server{
			OnSession: func(s *sse.Session) (sse.Subscription, bool) {
				channel := s.Req.PathValue("channel")
				return sse.Subscription{
					Client: s,
					Topics: append([]string{channel}, sse.DefaultTopic),
				}, true
			},
		},
		workflowManager: workflowmanager.NewWorkflowManager(),
	}
	server.server = http.Server{
		Addr:     addr,
		Handler:  server.routes(),
		ErrorLog: slog.NewLogLogger(logger.NewLogger().Handler(), slog.LevelError),
	}
	return server
}

func (srv *Server) Start() {
	srv.log.Info("starting server", slog.String("addr", srv.server.Addr))

	err := srv.server.ListenAndServe()
	if err != nil {
		srv.log.Error(err.Error())
		os.Exit(1)
	}
}
