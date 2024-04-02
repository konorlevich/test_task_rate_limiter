package handler

import (
	"net/http"

	"github.com/redis/go-redis/v9"
	log "github.com/sirupsen/logrus"

	"github.com/konorlevich/test_task_rate_limiter/internal/server/handler/message"
)

func NewHandler(r *redis.Client, l *log.Entry) http.Handler {
	handler := http.NewServeMux()

	handler.Handle("POST /message", message.NewSendMessageHandler(message.NewCachedRateLimiter(r, l), l))
	handler.Handle("POST /account", nil)
	handler.Handle("POST /transaction", nil)
	return handler
}
