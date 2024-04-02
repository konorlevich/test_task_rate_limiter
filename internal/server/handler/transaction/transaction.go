package transaction

import (
	"net/http"

	log "github.com/sirupsen/logrus"

	"github.com/konorlevich/test_task_rate_limiter/internal/server/middleware"
)

const (
	httpCodeFailedTransaction = http.StatusConflict
)

type RateLimiter interface {
	AllowedToSend(username string) (bool, error)
	TransactionFailed(username string) error
}

func NewMakeTransactionHandler(r RateLimiter, l *log.Entry) http.Handler {
	return middleware.CheckAuth(
		RateLimier(r, http.HandlerFunc(MakeTransaction), l),
	)
}
func MakeTransaction(rw http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
	}
	f := r.PostForm
	success := f.Get("success")
	if success == "1" {
		_, _ = rw.Write([]byte(`"{"status":"ok"}"`))
		return
	}
	http.Error(rw, "transaction failed", httpCodeFailedTransaction)
}

func RateLimier(rl RateLimiter, next http.Handler, l *log.Entry) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		username := r.PathValue(middleware.FieldNameUserName)

		if ok, err := rl.AllowedToSend(username); err != nil {
			l.WithError(err).Error("can't check if request has been sent before")
		} else if !ok {
			http.Error(rw, "You are not allowed to make a transaction today", http.StatusForbidden)
			return
		}

		next.ServeHTTP(
			NewCustomResponseWriter(rl, rw, username, l),
			r,
		)
	})
}

func NewCustomResponseWriter(
	r RateLimiter, //TODO: interface segregation
	w http.ResponseWriter,
	u string,
	l *log.Entry,
) http.ResponseWriter {
	return &customResponseWriter{r: r, w: w, u: u, l: l}
}

type customResponseWriter struct {
	r RateLimiter
	w http.ResponseWriter
	u string
	l *log.Entry
}

func (c customResponseWriter) Header() http.Header {
	return c.w.Header()
}

func (c customResponseWriter) Write(bytes []byte) (int, error) {
	return c.w.Write(bytes)
}

func (c customResponseWriter) WriteHeader(statusCode int) {
	c.l.WithField("status", statusCode).Info("transaction finished")
	if statusCode == httpCodeFailedTransaction {
		if err := c.r.TransactionFailed(c.u); err != nil {
			c.l.WithField("username", c.u).WithError(err).Error("can't save failed transaction")
		}
	}
	c.w.WriteHeader(statusCode)
}
