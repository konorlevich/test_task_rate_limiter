package handler

import (
	"net/http"
)

func NewHandler() http.Handler {
	handler := http.NewServeMux()

	handler.Handle("POST /message", nil)
	handler.Handle("POST /account", nil)
	handler.Handle("POST /transaction", nil)
	return handler
}
