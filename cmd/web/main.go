package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"regexp"
	"syscall"
	"time"

	server "github.com/f1monkey/spellchecker-web"
	"github.com/f1monkey/spellchecker-web/internal/logger"
	"github.com/f1monkey/spellchecker-web/internal/spellchecker"
)

var GitCommit string = "dev"

const (
	defaultServerAddr = "localhost:8011"
)

func main() {
	local, err := time.LoadLocation("UTC")
	if err != nil {
		panic(err)
	}
	time.Local = local

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer cancel()

	ctx = logger.WithContext(
		ctx,
		logger.New(GitCommit, os.Getenv("SPELLCHECKER_LOG_LEVEL")),
	)

	registry, err := initRegistry(ctx)
	if err != nil {
		logger.FromContext(ctx).Error("init spellchecker error", "error", err)
		os.Exit(1)
	}

	splitter, err := initWordSpliter()
	if err != nil {
		logger.FromContext(ctx).Error("init spellchecker error", "error", err)
		os.Exit(1)
	}

	defer registry.SaveAll(ctx)

	server := server.NewServer(ctx, registry, splitter)

	addr := defaultServerAddr
	if a := os.Getenv("SPELLCHECKER_HTTP_ADDR"); a != "" {
		addr = a
	}

	srv := &http.Server{
		Addr:    addr,
		Handler: server,
		BaseContext: func(net.Listener) context.Context {
			return ctx
		},
	}

	errors := make(chan error)
	go func() {
		logger.FromContext(ctx).Info("http server started", "address", addr)
		errors <- srv.ListenAndServe()
	}()

	select {
	case err := <-errors:
		if err != nil {
			logger.FromContext(ctx).Error("http server stopped", "error", err)
			os.Exit(1)
		}
	case <-ctx.Done():
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()

		if err := srv.Shutdown(ctx); err != nil {
			logger.FromContext(ctx).Error("http server graceful shutdown failed", "error", err)
			if err := srv.Close(); err != nil {
				logger.FromContext(ctx).Error("http server forced close failed", "error", err)
				os.Exit(1)
			}
		}
	}
}

func initRegistry(ctx context.Context) (*spellchecker.Registry, error) {
	dir := os.Getenv("SPELLCHECKER_DIR")
	if dir == "" {
		return nil, fmt.Errorf("env SPELLCHECKER_DIR must be provided")
	}

	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("unable to create dir %s: %w", dir, err)
	}

	var saveInterval time.Duration

	saveIntervalStr := os.Getenv("SPELLCHECKER_AUTOSAVE_INTERVAL")
	if saveIntervalStr != "" {
		i, err := time.ParseDuration(saveIntervalStr)
		if err != nil {
			return nil, fmt.Errorf("invalid SPELLCHECKER_AUTOSAVE_INTERVAL: %w", err)
		}

		saveInterval = i
	}

	result, err := spellchecker.NewRegistry(ctx, dir)
	if err != nil {
		return nil, err
	}

	result.AutoSave(ctx, saveInterval)

	return result, nil
}

var defaultRegexp = regexp.MustCompile(`['\pL]+`)

func initWordSpliter() (*regexp.Regexp, error) {
	value := os.Getenv("SPELLCHECKER_WORD_SPLIT_REGEXP")
	if value == "" {
		return defaultRegexp, nil
	}

	result, err := regexp.Compile(value)
	if err != nil {
		return nil, fmt.Errorf("invalid SPELLCHECKER_WORD_SPLIT_REGEXP: %w", err)
	}

	return result, nil
}
