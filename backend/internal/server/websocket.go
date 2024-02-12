// Inspired by
// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package server

import (
	"flag"
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/fc92/poker/internal/common/logger"
)

func init() {
	logger.InitLogger()
}

var httpListenAndServe = http.ListenAndServe

func StartServer(ws string) error {
	var addr = flag.String("addr", ws, "http service address")
	flag.Parse()
	hub := newHub()
	go hub.run()
	// Define the directory to serve static files from
	staticDir := "../../frontend/dist"
	// Custom handler to manage routing
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Exclude paths "/metrics" and "/ws" from serving static files
		if r.URL.Path != "/metrics" && r.URL.Path != "/ws" {
			// For any route (except /ws and /metrics), serve the index.html file
			http.ServeFile(w, r, staticDir+"/index.html")
			return
		}

		// Handle other paths ("/metrics", "/ws") as needed
		switch r.URL.Path {
		case "/metrics":
			promhttp.Handler().ServeHTTP(w, r)
		case "/ws":
			serveWs(hub, w, r)
		}
	})
	http.HandleFunc("/assets/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, staticDir+r.URL.Path)
		w.Header().Set("Content-Type", "application/javascript")
	})

	err := httpListenAndServe(*addr, nil)
	if err != nil {
		return err
	}
	return nil
}
