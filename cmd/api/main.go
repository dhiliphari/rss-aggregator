package main

import (
	"database/sql"
	"log/slog"
	"net/http"
	"os"

	_ "github.com/lib/pq"
)

type application struct {
	cfg    *config
	logger *slog.Logger
}

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	cfg := loadCfg()

	db, err := openDB(cfg.dsn)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
	defer db.Close()

	app := &application{
		cfg:    cfg,
		logger: logger,
	}

	srv := &http.Server{
		Addr:    cfg.addr,
		Handler: app.routes(),
	}

	logger.Info("starting server", "addr", cfg.addr)

	err = srv.ListenAndServe()
	logger.Error(err.Error())
	os.Exit(1)
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}
