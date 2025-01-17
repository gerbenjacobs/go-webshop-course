package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gerbenjacobs/go-webshop-course/handler"
	"github.com/gerbenjacobs/go-webshop-course/services"
	"github.com/gerbenjacobs/go-webshop-course/storage"
	"github.com/lmittmann/tint"
)

var address = "localhost:8000"

func main() {
	// handle shutdown signals
	shutdown := make(chan os.Signal, 3)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	// create logger
	output := os.Stdout
	tintOpt := &tint.Options{
		Level: slog.LevelDebug,
	}
	logger := slog.New(tint.NewHandler(output, tintOpt))

	// create our dependencies
	productRepo := storage.NewProductRepo()
	productSvc := services.NewProductService(productRepo)
	basketSvc := services.NewBasketService(storage.NewBasketRepo())
	deps := handler.Dependencies{
		Product: productSvc,
		Basket:  basketSvc,
	}

	// create a handler and server
	app := handler.New(logger, deps)
	srv := &http.Server{
		Addr:         address,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		Handler:      app,
	}

	// start running the server
	go func() {
		logger.Info("Server started", "address", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error("failed to listen", "error", err)
			os.Exit(1)
		}
	}()

	// wait for shutdown signals
	<-shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		logger.Error("Server shutdown failed", "error", err)
	}
	logger.Info("Server stopped successfully")
}
