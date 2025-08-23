package main

import (
	"context"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/f1monkey/spellchecker"
	server "github.com/f1monkey/spellchecker-web"
	"github.com/f1monkey/spellchecker-web/internal/logger"
)

var GitCommit string = "dev"

const (
	defaultServerAddr = "localhost:8011"

	defaultDBPath   = "spellchecker.bin"
	defaultAlphabet = spellchecker.DefaultAlphabet
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
		logger.New(GitCommit, os.Getenv("LOG_LEVEL")),
	)

	sc, err := initSpellchecker()
	if err != nil {
		logger.FromContext(ctx).Error("init spellchecker error", "error", err)
		os.Exit(1)
	}

	server := server.NewServer(ctx, sc)

	addr := defaultServerAddr
	if a := os.Getenv("HTTP_ADDR"); a != "" {
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

func initSpellchecker() (*spellchecker.Spellchecker, error) {
	dbPath := defaultDBPath
	if dbp := os.Getenv("SPELLCHECKER_DB_PATH"); dbp != "" {
		dbPath = dbp
	}

	f, err := os.Open(dbPath)
	if os.IsNotExist(err) {
		alphabet := defaultAlphabet
		if ab := os.Getenv("SPELLCHECKER_ALPHABET"); ab != "" {
			alphabet = ab
		}

		return spellchecker.New(alphabet)
	} else if err != nil {
		return nil, err
	}

	defer f.Close()

	return spellchecker.Load(f)
}
