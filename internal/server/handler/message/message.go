package message

import (
	"net/http"

	log "github.com/sirupsen/logrus"

	"github.com/konorlevich/test_task_rate_limiter/internal/server/middleware"
)

func NewSendMessageHandler(rl RateLimiter, l *log.Entry) http.Handler {

	return middleware.CheckAuth(RateLimier(rl, http.HandlerFunc(SendMessage), l))
}

func SendMessage(rw http.ResponseWriter, _ *http.Request) {
	_, _ = rw.Write([]byte("ok"))
}

type RateLimiter interface {
	AllowedToSend(username string) (bool, error)
	MessageSent(username string) error
}

func RateLimier(rl RateLimiter, next http.Handler, l *log.Entry) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		username := r.PathValue(middleware.FieldNameUserName)

		if ok, err := rl.AllowedToSend(username); err != nil {
			l.WithError(err).Error("can't check if request has been sent before")
		} else if !ok {
			http.Error(rw, "You send a new message too soo", http.StatusTooManyRequests)
			return
		}

		if err := rl.MessageSent(username); err != nil {
			l.WithError(err).Error("can't save message sent flag")
		}
		next.ServeHTTP(rw, r)
	})
}
