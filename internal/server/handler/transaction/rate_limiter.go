package transaction

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	log "github.com/sirupsen/logrus"
)

type CachedRateLimiter struct {
	redis *redis.Client
	l     *log.Entry
}

func NewCachedRateLimiter(r *redis.Client, l *log.Entry) *CachedRateLimiter {
	return &CachedRateLimiter{redis: r, l: l.WithField("module", "account.rate_limiter")}
}

func getCacheKey(username string) string {
	return fmt.Sprintf("transaction_failed_for_%s_%s", username, time.Now().Format("2006-01-02"))
}

func (r *CachedRateLimiter) TransactionFailed(username string) error {
	inc := r.redis.Incr(context.Background(), getCacheKey(username))
	if inc.Err() != nil {
		return inc.Err()
	}
	_, err := inc.Result()
	return err
}

func (r *CachedRateLimiter) AllowedToSend(username string) (bool, error) {
	got := r.redis.Get(context.Background(), getCacheKey(username))
	if got.Err() != nil {
		r.l.WithError(got.Err()).Info("received error")
		if errors.Is(got.Err(), redis.Nil) {
			return true, nil
		}
		return true, got.Err()
	}
	res, err := got.Int()
	if err != nil {
		r.l.WithError(err).Info("received error in res")
		if errors.Is(err, redis.Nil) {
			return true, nil
		}
		return true, err
	}
	r.l.WithField("res", res).Info("received result")
	if res >= 3 {
		return false, nil
	}
	return true, nil
}
