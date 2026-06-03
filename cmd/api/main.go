package main

import (
	"context"
	"log"
	"os/signal"
	"syscall"
	"time"

	"github.com/mrbananaaa/gosocialize/internal/api"
)

func main() {
	s, err := api.NewServer()
	if err != nil {
		panic(err)
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	go func() {
		if err := s.Run(); err != nil {
			log.Fatal(err)
		}
	}()

	<-ctx.Done()
	log.Println("Shutting down server...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := s.Shutdown(shutdownCtx); err != nil {
		log.Println("Failed to shutdown http server")
		log.Fatal(err)
	}

	log.Println("Server closed gracefully")
}
