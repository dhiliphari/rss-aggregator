package main

import (
	"log/slog"
	"net/http"
	"os"
)

type application struct {
	cfg    *config
	logger *slog.Logger
}

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	cfg := loadCfg()

	app := &application{
		cfg:    cfg,
		logger: logger,
	}

	srv := &http.Server{
		Addr:    cfg.addr,
		Handler: app.routes(),
	}

	logger.Info("starting server", "addr", cfg.addr)

	err := srv.ListenAndServe()
	logger.Error(err.Error())
	os.Exit(1)
}
