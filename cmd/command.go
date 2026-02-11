package cmd

import (
	"embed"
	"flag"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/naoki-higashi-28/mdprev/internal/dependency"
	"github.com/naoki-higashi-28/mdprev/internal/infrastructure/server"
)

//go:embed all:dist
var distFS embed.FS

var defaultPort = "0"

const (
	defaultAutoCloseGrace  = 5 * time.Second
	defaultShutdownTimeout = 5 * time.Second
)

type options struct {
	host      string
	port      string
	open      bool
	autoClose bool
	root      string
}

func parseOptions() options {
	host := flag.String("host", "127.0.0.1", "bind host")
	port := flag.String("port", defaultPort, "listen port (0 for random)")
	open := flag.Bool("open", true, "open browser automatically")
	autoClose := flag.Bool("auto-close", true, "shutdown when all connections close")
	flag.Parse()

	root := "."
	if flag.NArg() > 0 {
		root = flag.Arg(0)
	}

	return options{
		host:      *host,
		port:      *port,
		open:      *open,
		autoClose: *autoClose,
		root:      root,
	}
}

func resolveRoot(root string) (string, error) {
	absRoot, err := filepath.Abs(root)
	if err != nil {
		return "", fmt.Errorf("resolving root path: %w", err)
	}

	info, err := os.Stat(absRoot)
	if err != nil {
		return "", fmt.Errorf("stating root path: %w", err)
	}
	if !info.IsDir() {
		return "", fmt.Errorf("root path is not a valid directory: %s", absRoot)
	}

	return absRoot, nil
}

func subDistFS() (fs.FS, error) {
	sub, err := fs.Sub(distFS, "dist")
	if err != nil {
		return nil, fmt.Errorf("creating sub filesystem: %w", err)
	}
	return sub, nil
}

func setupConnectionTracker(autoClose bool) (*server.ConnectionTracker, func(), func()) {
	if !autoClose {
		return nil, nil, nil
	}

	tracker := server.NewConnectionTracker(defaultAutoCloseGrace)
	return tracker, tracker.Add, tracker.Remove
}

func execute() error {
	opts := parseOptions()

	absRoot, err := resolveRoot(opts.root)
	if err != nil {
		return err
	}

	staticFS, err := subDistFS()
	if err != nil {
		return err
	}

	tracker, onConnect, onDisconnect := setupConnectionTracker(opts.autoClose)

	mux, subscriber, err := dependency.NewServerMux(absRoot, staticFS, onConnect, onDisconnect)
	if err != nil {
		return fmt.Errorf("initializing server dependencies: %w", err)
	}
	runner := server.NewHTTPServerRunner(server.RunnerConfig{
		Host:            opts.host,
		Port:            opts.port,
		OpenBrowser:     opts.open,
		ShutdownTimeout: defaultShutdownTimeout,
	}, tracker)
	return runner.Run(absRoot, mux, subscriber)
}

func Execute() {
	if err := execute(); err != nil {
		log.Fatalf("Failed to run mdprev: %v", err)
	}
}
