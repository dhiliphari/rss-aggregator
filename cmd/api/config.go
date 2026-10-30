package main

import "flag"

type config struct {
	addr string
}

func loadCfg() *config {
	var cfg config

	flag.StringVar(&cfg.addr, "addr", ":5000", "HTTP network address")
	flag.Parse()

	return &cfg
}
