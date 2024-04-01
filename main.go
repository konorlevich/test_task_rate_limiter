package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os/signal"
	"syscall"

	log "github.com/sirupsen/logrus"

	"github.com/konorlevich/test_task_rate_limiter/internal/server/handler"
)

var port = 8080

func main() {

	l := log.New().WithFields(log.Fields{
		"port": port,
	})
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	defer l.Println("got interruption signal")

	server := &http.Server{Addr: fmt.Sprintf(":%d", port), Handler: handler.NewHandler()}

	go func() {
		l.Printf("listening to port %d", port)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			l.WithError(err).Fatal("listen and serve returned err")
		}
	}()
	defer func() {
		if err := server.Shutdown(context.Background()); err != nil {
			l.WithError(err).Error("handler shutdown returned an err")
		}
	}()

	<-ctx.Done()

}
