// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
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

func broadcastMessage(calculationResultChannel <-chan *CalculationResult, hub *websocketserver.Hub) {
	for {
		for client := range hub.Clients {
			byteChan := ResultChannelConvertToByteArrayChannel(calculationResultChannel)
			select {
			case client.Send <- <-byteChan:
			}
		}

	}
}

func main() {
	flag.Parse()

	hub := websocketserver.NewHub()
	go hub.Run()

	calculationResultChannel := getCalculationResultChannel()
	go broadcastMessage(calculationResultChannel, hub)

	http.HandleFunc("/", serveHome)
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		websocketserver.ServeWs(hub, w, r)
	})
	err := http.ListenAndServe(*addr, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
