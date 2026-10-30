package main

import "flag"

type config struct {
	addr string
	dsn  string
}

func loadCfg() *config {
	var cfg config

	flag.StringVar(&cfg.addr, "addr", ":5000", "HTTP network address")
	flag.StringVar(&cfg.dsn, "dsn", "postgres://user:password@localhost:5432/rssagg?sslmode=disable", "Database Source Name")
	flag.Parse()

	return &cfg
}
