package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/kodeyeen/shortify/docs"
	"github.com/kodeyeen/shortify/internal/config"
	"github.com/kodeyeen/shortify/internal/delivery/http/httpmw"
	httpdel "github.com/kodeyeen/shortify/internal/delivery/http/v1"
	"github.com/kodeyeen/shortify/internal/generation/rand"
	"github.com/kodeyeen/shortify/internal/persistence"
	"github.com/kodeyeen/shortify/internal/persistence/inmemory"
	"github.com/kodeyeen/shortify/internal/persistence/postgres"
	"github.com/kodeyeen/shortify/internal/url"
	httpswagger "github.com/swaggo/http-swagger/v2"
)

//	@title			Shortify API
//	@version		1.0
//	@description	This is the Ozon internship assignment.
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	Ruslan Iskandarov
//	@contact.url	https://www.t.me/ixderious
//	@contact.email	scanderoff@gmail.com

// @license.name	BSD 3-Clause "New" or "Revised" License
func main() {
	cfg := config.MustLoad()

	ctx := context.Background()

	log := newLogger(cfg.Env)

	log.Info("starting shortify", slog.String("env", cfg.Env))
	log.Debug("debug log level enabled")

	var urlRepo url.Repository

	switch cfg.PersistenceType {
	case config.PersistenceTypeInmemory:
		urlRepo = inmemory.NewURLRepository()
	case config.PersistenceTypePostgres:
		connString := persistence.NewConnString(
			"postgres",
			cfg.Postgres.Username,
			cfg.Postgres.Password,
			cfg.Postgres.Host,
			cfg.Postgres.Database,
		)

		dbpool, err := pgxpool.New(ctx, connString)
		if err != nil {
			log.Error("failed to create dbpool", slog.String("error", err.Error()))
			os.Exit(1)
		}
		defer dbpool.Close()

		err = dbpool.Ping(ctx)
		if err != nil {
			log.Error("failed to ping db", slog.String("error", err.Error()))
			os.Exit(1)
		}

		urlRepo = postgres.NewURLRepository(dbpool)
	}

	aliasPrvr := rand.NewAliasProvider(cfg.Alias.Charset, cfg.Alias.Length)
	urlSvc := url.NewService(urlRepo, aliasPrvr, log)
	urlClr := httpdel.NewURLController(urlSvc, log)

	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)
	router.Use(httpmw.NewLogger(log))

	router.Route("/api/v1", func(r chi.Router) {
		r.Post("/urls", urlClr.Create)
		r.Get("/urls/{alias}", urlClr.GetByAlias)
	})

	router.Get("/swagger/*", httpswagger.Handler(
		httpswagger.URL(fmt.Sprintf("http://localhost:%d/swagger/doc.json", cfg.HTTPServer.Port)),
	))

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	addr := fmt.Sprintf(":%d", cfg.HTTPServer.Port)

	log.Info("starting api server", slog.String("address", addr))

	srv := &http.Server{
		Addr:         addr,
		Handler:      router,
		ReadTimeout:  cfg.HTTPServer.ReadTimeout,
		WriteTimeout: cfg.HTTPServer.WriteTimeout,
		IdleTimeout:  cfg.HTTPServer.IdleTimeout,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Error("failed to start server")
		}
	}()

	log.Info("server started")

	<-done
	log.Info("stopping server")

	ctx, cancel := context.WithTimeout(context.Background(), cfg.HTTPServer.ShutdownTimeout)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Error("failed to stop server", slog.String("error", err.Error()))

		os.Exit(1)
	}

	log.Info("server stopped")
}

func newLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case config.EnvLocal:
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
				Level: slog.LevelDebug,
			}),
		)
	case config.EnvDev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
				Level: slog.LevelDebug,
			}),
		)
	case config.EnvProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
				Level: slog.LevelInfo,
			}),
		)
	}

	return log
}
