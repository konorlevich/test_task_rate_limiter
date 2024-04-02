package message

import (
	"context"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"
	log "github.com/sirupsen/logrus"
)

type CachedRateLimiter struct {
	redis *redis.Client
	l     *log.Entry
}

func NewCachedRateLimiter(r *redis.Client, l *log.Entry) *CachedRateLimiter {
	return &CachedRateLimiter{redis: r, l: l.WithField("module", "message.rate_limiter")}
}

func (r *CachedRateLimiter) MessageSent(username string) error {
	set := r.redis.Set(context.Background(), "message_sent_"+username, 1, 1*time.Second)
	if set.Err() != nil {
		return set.Err()
	}
	_, err := set.Result()
	return err
}

func (r *CachedRateLimiter) AllowedToSend(username string) (bool, error) {
	got := r.redis.Get(context.Background(), "message_sent_"+username)
	if got.Err() != nil {
		r.l.WithError(got.Err()).Info("received error")
		if errors.Is(got.Err(), redis.Nil) {
			return true, nil
		}
		return true, got.Err()
	}
	res, err := got.Result()
	if err != nil {
		r.l.WithError(err).Info("received error in res")
		if errors.Is(err, redis.Nil) {
			return true, nil
		}
		return true, err
	}
	r.l.WithField("res", res).Info("received result")
	return false, nil
}
