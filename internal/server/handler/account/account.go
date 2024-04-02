package account

import (
	"net/http"
	"strings"

	log "github.com/sirupsen/logrus"
)

type RateLimiter interface {
	AllowedToCreate(ip string) (bool, error)
	AccountCreated(ip string) error
}

func NewCreateAccountHandler(r RateLimiter, l *log.Entry) http.Handler {
	return RateLimier(r, http.HandlerFunc(CreateAccount), l)
}
func CreateAccount(rw http.ResponseWriter, _ *http.Request) {
	_, _ = rw.Write([]byte(`ok`))
}

func RateLimier(rl RateLimiter, next http.Handler, l *log.Entry) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		ip := strings.Split(r.RemoteAddr, ":")[0]

		if ok, err := rl.AllowedToCreate(ip); err != nil {
			l.WithError(err).Error("can't check if request has been sent before")
		} else if !ok {
			http.Error(rw,
				"You are not allowed to create more than 20 accounts per day",
				http.StatusTooManyRequests)
			return
		}

		if err := rl.AccountCreated(ip); err != nil {
			l.WithError(err).Error("can't save account created flag")
		}
		next.ServeHTTP(rw, r)
	})
}
