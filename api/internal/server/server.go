package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/edgejay/pify-player/api/internal/handlers"
)

type Server struct {
	port        string
	corsOrigins []string
	sslDomain   string
	e           *echo.Echo
}

func NewServer(port string, corsOrigins []string, sslDomain string) *Server {
	return &Server{port, corsOrigins, sslDomain, nil}
}

func (svr *Server) Start() {
	if _, err := strconv.Atoi(svr.port); err != nil {
		svr.port = "8080"
	}

	svr.e = echo.New()

	// Enable logging middleware
	svr.e.Use(middleware.Logger())
	svr.e.Use(middleware.Recover())

	if len(svr.corsOrigins) > 0 {
		log.Printf("Setting up CORS for origins: %v\n", svr.corsOrigins)

		// Enable CORS middleware
		svr.e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
			AllowOrigins: svr.corsOrigins, // Change this to specific domains for security
			AllowMethods: []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete},
		}))
	}

	// Handlers
	// Handlers
	apiGroup := svr.e.Group("/api")
	authGroup := apiGroup.Group("/auth")
	handlers.SetAuthRoutes(authGroup)

	log.Printf("Server is running on port %s...\n", svr.port)

	// Start the server
	if err := svr.e.StartTLS(
		fmt.Sprintf(":%s", svr.port),
		fmt.Sprintf("./certs/%s.pem", svr.sslDomain),
		fmt.Sprintf("./certs/%s.key.pem", svr.sslDomain),
	); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Server failed to start: %v\n", err)
	}
}

func (svr *Server) Shutdown(ctx context.Context) {
	if err := svr.e.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}
}
