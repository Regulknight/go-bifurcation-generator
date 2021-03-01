// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bifurcation-generator/generator"
	"bifurcation-generator/subsequencesearcher"
	"bifurcation-generator/websocketserver"
	"flag"
	"log"
	"net/http"
	"strconv"
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

func getFloatSliceConverter(float64Chan <-chan []float64) <-chan []byte {
	out := make(chan []byte)

	go func() {
		dst := []byte{}
		floatSlice := <-float64Chan
		for i := 0; i < len(floatSlice); i++ {
			dst = append(dst[:], strconv.AppendFloat(dst, floatSlice[i], 'E', -1, 32)[:]...)
		}

		out <- dst
	}()

	return out
}

func broadcastMessage(subsequenceChan <-chan []float64, hub *websocketserver.Hub) {
	for {
		for client := range hub.Clients {
			byteChan := getFloatSliceConverter(subsequenceChan)
			msg := <-byteChan
			select {
			case client.Send <- msg:
			}
		}

	}
}

func getBifurcationCyclesChannel() <-chan []float64 {
	out := make(chan []float64)

	go func() {
		bifurcationSequenceGenerator := generator.GetBifurcationSequenceGenerator(0.4, 0.0)
		for r := 0.0; r < 3.9; r += 0.1 {
			bifurcationSequenceGenerator = generator.GetBifurcationSequenceGenerator(0.4, r)
			calculationSlice, ok := <-bifurcationSequenceGenerator

			for ok {
				subsequence := subsequencesearcher.IsContainsSubsequences(calculationSlice)

				if subsequence != nil || len(calculationSlice) > 40000 {
					out <- subsequence
					ok = false
				} else {
					calculationSlice, ok = <-bifurcationSequenceGenerator
				}
			}

		}
	}()

	return out
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
