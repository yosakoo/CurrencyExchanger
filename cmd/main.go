package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/yosakoo/CurrencyExchanger/internal/config"
	"github.com/yosakoo/CurrencyExchanger/internal/currency"
	"github.com/yosakoo/CurrencyExchanger/pkg/postgres"
)

func main() {
	cfg := config.Load()

	log.Printf("Server starting on port: %s", cfg.Port)

	conn, err := postgres.NewConnection(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Error connect to daabase: %v", err)
	}

	storage := currency.NewStorage(conn)
	mux := currency.NewRouter(storage)

	srv := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: mux,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed: %v", err)
		}
	}()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	<-ctx.Done()

	log.Println("Shutting down server...")

	if err := srv.Shutdown(context.Background()); err != nil {
		log.Printf("Server shutdown failed: %v", err)
		panic(err)
	}

	log.Println("Server stopped.")
}
