// Inspired by
// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package server

import (
	"flag"
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog/log"

	"github.com/fc92/poker/internal/common/logger"
)

func init() {
	logger.InitLogger()
}

var httpListenAndServe = http.ListenAndServe

func serveHome(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		log.Error().Msgf("URL not supported: %s", r.URL)
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	if r.Method != http.MethodGet {
		log.Error().Msgf("URL %s, method not supported: %s", r.URL, r.Method)
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	http.ServeFile(w, r, "home.html")
}

func StartServer(ws string) error {
	var addr = flag.String("addr", ws, "http service address")
	flag.Parse()
	hub := newHub()
	go hub.run()
	http.HandleFunc("/", serveHome)
	http.Handle("/metrics", promhttp.Handler())
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWs(hub, w, r)
	})
	err := httpListenAndServe(*addr, nil)
	if err != nil {
		return err
	}
	return nil
}