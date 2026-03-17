package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/StreamPeek/StreamPeek/internal/config"
	"github.com/StreamPeek/StreamPeek/internal/kafka"
	"github.com/StreamPeek/StreamPeek/internal/server"
)

func main() {
	cfg := config.Load()

	log.Printf("Starting StreamPeek on port %s...", cfg.Port)
	log.Printf("Configured Kafka brokers: %s", cfg.KafkaBrokers)

	// Initialize Kafka Client
	kafkaClient, err := kafka.NewClient(cfg.KafkaBrokers)
	if err != nil {
		log.Fatalf("Fatal error initializing kafka client: %v", err)
	}
	defer kafkaClient.Close()
	log.Println("Successfully initialized Kafka client.")

	// Set up Router
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Register Routes
	produceHandler := server.NewProduceHandler(kafkaClient)
	produceHandler.RegisterRoutes(r)

	// Health Check
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// Configure HTTP Server
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%s", cfg.Port),
		Handler:      r,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	// Graceful Shutdown
	go func() {
		log.Printf("Server listening on %s", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exiting")
}
