package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"time"

	"server_course/api"
	"server_course/api/middleware"
	"server_course/db"
)

func run(ctx context.Context, l *slog.Logger, debugMode bool) error {
	ctx, stop := signal.NotifyContext(ctx, os.Interrupt)
	defer stop()
	

	db, err := db.NewDB("./db")
	if err != nil {
		return err
	}

	middleware := middleware.NewMiddleware(l)

	router := api.NewServer(l, middleware, db)

	srv := &http.Server{
		Addr: ":8080",
		Handler: router,
	}

	go func() {
		l.Info("started listening on server", slog.String("addr", srv.Addr))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Fprintf(os.Stderr, "error listening and serving: %s\n", err)
		}
	}()

	// Listen for the interrupt signal.
	<-ctx.Done()

	// Restore default behavior on the interrupt signal and notify user of shutdown.
	stop()
	l.Info("shutting down gracefully, press Ctrl+C again to force")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		fmt.Fprintln(os.Stderr, "Server forced to shutdown: ", err)
	}

	if debugMode {
		err := db.Reset()
		if err != nil {
			fmt.Fprintln(os.Stderr, "[debug mode] Failed to reset DB: ", err)
			return err
		}
	}

	return nil
}

func main() {
	logger := slog.Default()

	dbg := flag.Bool("debug", false, "Enable debug mode")
	flag.Parse()

	ctx := context.Background()
	if err := run(ctx, logger, *dbg); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}