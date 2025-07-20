package main

import (
	"context"
	"flag"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/Sidney-Bernardin/Chat-Alloy/internal"
	"github.com/Sidney-Bernardin/Chat-Alloy/internal/repos/postgres"
	"github.com/Sidney-Bernardin/Chat-Alloy/internal/repos/redis"
	"github.com/Sidney-Bernardin/Chat-Alloy/internal/service"
	"github.com/Sidney-Bernardin/Chat-Alloy/internal/web"
	"github.com/Sidney-Bernardin/Chat-Alloy/internal/web/pages/home"
	"github.com/Sidney-Bernardin/Chat-Alloy/internal/web/users"
	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
)

var (
	logJSON = flag.Bool("log-json", false, "JSON logger")
)

func main() {
	flag.Parse()
	errs, ctx := errgroup.WithContext(context.Background())

	// Create logger.
	log := slog.New(createLogHandler())

	// Create configuration.
	cfg, err := internal.NewConfig()
	if err != nil {
		log.Error("Cannot create configuration", "err", err.Error())
		return
	}

	// Create postgres repository.
	postgresRepo, err := postgres.New(ctx, cfg)
	if err != nil {
		log.Error("Cannot create postgres repository", "err", err.Error())
		return
	}

	// Create redis repository.
	redisRepo, err := redis.New(ctx, cfg)
	if err != nil {
		log.Error("Cannot create redis repository", "err", err.Error())
		return
	}

	/////

	svc := &service.Service{
		Config:   cfg,
		Logger:   log,
		Postgres: postgresRepo,
		Redis:    redisRepo,
	}

	// Start new web server.
	errs.Go(func() error {
		return errors.Wrap(startWebServer(ctx, &web.Server{
			Server: &http.Server{
				Addr: cfg.ADDR,
			},
			Config:  cfg,
			Logger:  log,
			Service: svc,
		}), "web server failed")
	})

	// Wait for signal interrupts.
	errs.Go(func() error {
		sigCtx, sigCancel := signal.NotifyContext(ctx, os.Interrupt, os.Kill)
		<-sigCtx.Done()
		sigCancel()
		return errors.Wrap(sigCtx.Err(), "signal interrupt")
	})

	log.Error("Done", "err", errs.Wait().Error())
}

func createLogHandler() (h slog.Handler) {
	if *logJSON {
		h = slog.NewJSONHandler(os.Stderr, nil)
	} else {
		h = slog.NewTextHandler(os.Stderr, nil)
	}
	return h
}

func startWebServer(ctx context.Context, svr *web.Server) error {
	go func() {
		<-ctx.Done()
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := svr.Shutdown(ctx); err != nil {
			svr.Logger.Error("Cannot shutdown web server", "err", err.Error())
			return
		}
	}()

	r := http.NewServeMux()
	r.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("tmp/dist"))))
	r.Handle("/", home.Handler(svr))
	r.Handle("/signup", users.HandleSignup(svr))
	svr.Handler = svr.MWLog(r)

	return svr.ListenAndServe()
}
