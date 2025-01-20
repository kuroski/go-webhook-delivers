package server

import (
	"net/http"
)

func (srv *Server) routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /up", srv.upHandler)
	mux.HandleFunc("POST /webhook", srv.webhookHandler)
	mux.HandleFunc("GET /channel/{channel}", srv.payloadUrlHandler)
	mux.HandleFunc("POST /channel/{channel}", srv.forwardRequestHandler)
	return mux
}
