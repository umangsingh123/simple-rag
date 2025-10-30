package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// Features:
// - Starts HTTP server
// - Handles graceful shutdown (Ctrl+C)
// - Configures timeouts
// - Proper error handling
type Server struct {
	httpServer *http.Server
}

// Dependency Injection: Takes addr and handler as parameters (separation of concerns)
func NewServer(addr string, handler http.Handler) *Server {
	return &Server{
		httpServer: &http.Server{
			Addr:         addr,
			Handler:      handler,
			ReadTimeout:  15 * time.Second,
			WriteTimeout: 30 * time.Second,
			IdleTimeout:  60 * time.Second,
		},
	}
}

// Non-blocking start: Uses goroutine to start server in background
// Error handling: Checks for errors other than http.ErrServerClosed (normal shutdown)
// Fatal errors: Exits process on serious server failures
// Immediate return: Allows caller to proceed with other initialization

func (s *Server) Start() error {
	fmt.Println(":: Starting server on::", s.httpServer.Addr)

	// Start server in background
	go func() {
		if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Printf(":::::::::: Server failed::::::::: %v\n", err)
			os.Exit(1)
		}
	}()

	return nil
}

/*
- Graceful shutdown: Uses http.Server.Shutdown() which stops accepting new requests and waits for active ones to complete
- Timeout context: Prevents indefinite waiting during shutdown
- Proper cleanup: defer cancel() ensures context cancellation
- Error wrapping: Returns formatted error if shutdown fails
*/
func (s *Server) Shutdown(timeout time.Duration) error {
	fmt.Println(":::::::::: Shutting down server...::::::::::")

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	if err := s.httpServer.Shutdown(ctx); err != nil {
		return fmt.Errorf("shutdown failed: %v", err)
	}

	fmt.Println(":::::::::: Server shutdown gracefully ::::::::::")
	return nil
}

/*
- Signal handling: Listens for SIGINT (Ctrl+C) and SIGTERM (container shutdown)

- Blocking wait: Stops main goroutine until shutdown signal received

- Graceful termination: Gives active requests 10 seconds to complete

- Emergency exit: Force exits if graceful shutdown fails
*/
func (s *Server) WaitForShutdown() {
	// Wait for interrupt signal
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)

	<-signalChan
	fmt.Println("\n ::::::::::Received shutdown signal::::::::::")

	// Graceful shutdown with 10 second timeout
	if err := s.Shutdown(10 * time.Second); err != nil {
		fmt.Printf(":::::::::: Shutdown error: %v\n", err)
		os.Exit(1)
	}
}
