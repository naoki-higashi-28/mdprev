package server

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"time"
)

// RunnerConfig defines runtime options for the HTTP server.
type RunnerConfig struct {
	Host            string
	Port            string
	OpenBrowser     bool
	ShutdownTimeout time.Duration
}

// HTTPServerRunner runs and manages the HTTP server lifecycle.
type HTTPServerRunner struct {
	config  RunnerConfig
	tracker *ConnectionTracker
}

// NewHTTPServerRunner creates a new HTTPServerRunner.
func NewHTTPServerRunner(config RunnerConfig, tracker *ConnectionTracker) *HTTPServerRunner {
	return &HTTPServerRunner{config: config, tracker: tracker}
}

// Run starts the server and blocks until it stops.
func (r *HTTPServerRunner) Run(root string, handler http.Handler, closer io.Closer) error {
	if closer != nil {
		defer func() {
			if err := closer.Close(); err != nil {
				log.Printf("Failed to close resource: %v", err)
			}
		}()
	}

	addr := net.JoinHostPort(r.config.Host, r.config.Port)
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("listening on %s: %w", addr, err)
	}

	url := fmt.Sprintf("http://%s", ln.Addr().String())
	fmt.Printf("mdprev serving %s on %s\n", root, url)

	if r.config.OpenBrowser {
		OpenBrowser(url)
	}

	srv := &http.Server{Handler: handler}
	r.startAutoShutdown(srv)

	if err := srv.Serve(ln); err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("serving http: %w", err)
	}
	return nil
}

func (r *HTTPServerRunner) startAutoShutdown(srv *http.Server) {
	if r.tracker == nil {
		return
	}

	go func() {
		<-r.tracker.Done()
		fmt.Println("\nNo active connections. Shutting down...")
		ctx, cancel := context.WithTimeout(context.Background(), r.config.ShutdownTimeout)
		defer cancel()
		_ = srv.Shutdown(ctx)
	}()
}
