package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/edgejay/pify-player/api/internal/server"
	"github.com/edgejay/pify-player/api/internal/utils"
)

func main() {
	serverSettings := utils.GetServerSettings()
	port := serverSettings.Port
	corsOrigins := serverSettings.CorsOrigins
	sslDomain := serverSettings.SslDomain

	server := server.NewServer(port, corsOrigins, sslDomain)

	// Channel to listen for OS signals for graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	// Start server in a goroutine
	go func() {
		server.Start()
	}()

	// Wait for termination signal
	<-stop
	log.Println("\nShutting down server gracefully...")

	// Shutdown with timeout
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	server.Shutdown(shutdownCtx)

	log.Println("Server shutdown complete.")
}
