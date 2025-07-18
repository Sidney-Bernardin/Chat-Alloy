package main

import (
	"context"
	"flag"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/Sidney-Bernardin/Chat-Alloy/server"
	"github.com/Sidney-Bernardin/Chat-Alloy/server/handlers"
	"github.com/Sidney-Bernardin/Chat-Alloy/server/service"

	"github.com/pkg/errors"
)

var (
	logJSON = flag.Bool("log-json", false, "JSON logger")
)

func main() {
	flag.Parse()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Create logger.
	log := slog.New(createLogHandler())

	// Create configuration.
	cfg, err := server.NewConfig()
	if err != nil {
		log.Error("Cannot create configuration", "err", err.Error())
		return
	}

	svc := service.New(cfg, log, nil)
	handler := handlers.New(cfg, log, svc)

	go func() {
		defer cancel()
		if err = handler.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Error("Cannot listen and serve", "err", err.Error())
		}
	}()

	log.Info("Ready", "addr", cfg.ADDR)

	sigCtx, sigCancel := signal.NotifyContext(ctx, os.Interrupt, os.Kill)
	<-sigCtx.Done()
	sigCancel()

	// Shutdown the server.
	ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
	if err := handler.Shutdown(ctx); err != nil {
		log.Error("Cannot shutdown API", "err", err.Error())
		return
	}

	log.Info("Done")
}

func createLogHandler() (h slog.Handler) {
	if *logJSON {
		h = slog.NewJSONHandler(os.Stderr, nil)
	} else {
		h = slog.NewTextHandler(os.Stderr, nil)
	}
	return h
}
