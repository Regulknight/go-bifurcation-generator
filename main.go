// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bifurcation-generator/converter"
	"bifurcation-generator/websocketserver"
	"flag"
	"log"
	"net/http"
)

var addr = flag.String("addr", ":8083", "http service address")

func serveHome(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL)
	if r.URL.Path != "/ws_check" {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	http.ServeFile(w, r, "home.html")
}

func broadcastMessage(subsequenceChan <-chan []float64, hub *websocketserver.Hub) {
	for {
		for client := range hub.Clients {
			byteChan := converter.GetFloatSliceConverter(subsequenceChan)
			msg := <-byteChan
			select {
			case client.Send <- msg:
			}
		}

	}
}

func main() {
	flag.Parse()

	hub := websocketserver.NewHub()
	go hub.Run()

	go broadcastMessage(getBifurcationCyclesChannel(), hub)

	http.HandleFunc("/", serveHome)
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		websocketserver.ServeWs(hub, w, r)
	})
	err := http.ListenAndServe(*addr, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
