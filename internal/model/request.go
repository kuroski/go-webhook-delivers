package model

import (
	"net/http"
	"time"
)

type Request struct {
	Headers   http.Header
	Body      []byte
	Query     string
	Timestamp time.Time
}
