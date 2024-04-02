package handler

import (
	"net/http"

	"github.com/redis/go-redis/v9"
	log "github.com/sirupsen/logrus"

	"github.com/konorlevich/test_task_rate_limiter/internal/server/handler/account"
	"github.com/konorlevich/test_task_rate_limiter/internal/server/handler/message"
	"github.com/konorlevich/test_task_rate_limiter/internal/server/handler/transaction"
)

func NewHandler(r *redis.Client, l *log.Entry) http.Handler {
	handler := http.NewServeMux()

	handler.Handle("POST /message", message.NewSendMessageHandler(message.NewCachedRateLimiter(r, l), l))
	handler.Handle("POST /account", account.NewCreateAccountHandler(account.NewCachedRateLimiter(r, l), l))
	handler.Handle("POST /transaction", transaction.NewMakeTransactionHandler(transaction.NewCachedRateLimiter(r, l), l))
	return handler
}
