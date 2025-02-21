package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/edgejay/pify-player/api/internal/handlers"
)

func main() {
	port := os.Getenv("PORT")
	corsOrigins := strings.Split(os.Getenv("CORS_ORIGINS"), ",")
	sslDomain := os.Getenv("SSL_DOMAIN")

	if _, err := strconv.Atoi(port); err != nil {
		port = "8080"
	}

	e := echo.New()

	// Enable logging middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	if len(corsOrigins) > 0 {
		log.Printf("Setting up CORS for origins: %v\n", corsOrigins)

		// Enable CORS middleware
		e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
			AllowOrigins: corsOrigins, // Change this to specific domains for security
			AllowMethods: []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete},
		}))
	}

	// Handlers
	apiGroup := e.Group("/api")
	authGroup := apiGroup.Group("/auth")
	handlers.SetAuthRoutes(authGroup)

	// Channel to listen for OS signals for graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	// Start server in a goroutine
	go func() {
		log.Printf("Server is running on port %s...\n", port)

		// Start the server
		if err := e.StartTLS(fmt.Sprintf(":%s", port), fmt.Sprintf("./certs/%s.pem", sslDomain), fmt.Sprintf("./certs/%s.key.pem", sslDomain)); err != nil && err != http.ErrServerClosed {
			e.Logger.Fatal("Shutting down the server")
		}
	}()

	// Wait for termination signal
	<-stop
	log.Println("\nShutting down server gracefully...")

	// Shutdown with timeout
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := e.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server shutdown complete.")
}
