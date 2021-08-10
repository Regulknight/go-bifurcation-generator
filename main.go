// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"github.com/regulknight/go-bifurcation-generator/bifurcation"
	"github.com/regulknight/go-bifurcation-generator/converter"
	"github.com/regulknight/go-bifurcation-generator/iterator"
	"github.com/regulknight/go-bifurcation-generator/searcher"
	"github.com/regulknight/go-bifurcation-generator/websocketserver"
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

	cycleSearcher := searcher.NewCycleSearcher(bifurcation.NewBifurcation(bifurcation.DefaultBifurcationFunction(), iterator.NewSegmentIterator(0.0, 3.9, 0.1), 0.4).GetBifurcationChannel())

	hub := websocketserver.NewHub()
	go hub.Run()

	go broadcastMessage(cycleSearcher.GetCyclesChannel(), hub)

	http.HandleFunc("/", serveHome)
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		websocketserver.ServeWs(hub, w, r)
	})
	err := http.ListenAndServe(*addr, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
