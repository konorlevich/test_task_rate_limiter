package account

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

func getCacheKey(ip string) string {
	return fmt.Sprintf("account_created_from_%s_%s", ip, time.Now().Format("2006-01-02"))
}

func (r *CachedRateLimiter) AccountCreated(ip string) error {
	inc := r.redis.Incr(context.Background(), getCacheKey(ip))
	if inc.Err() != nil {
		return inc.Err()
	}
	_, err := inc.Result()
	return err
}

func (r *CachedRateLimiter) AllowedToCreate(ip string) (bool, error) {
	got := r.redis.Get(context.Background(), getCacheKey(ip))
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
	if res > 20 {
		return false, nil
	}
	return true, nil
}
