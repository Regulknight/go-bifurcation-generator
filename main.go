// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"net/http"
	"bufio"
	"log"
	"os"
)

var addr = flag.String("addr", ":8080", "http service address")

func serveHome(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL)
	if r.URL.Path != "/" {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	http.ServeFile(w, r, "home.html")
}

func broadcast_message(reader *bufio.Reader, hub *Hub) {
	for {
		message, err := reader.ReadBytes('\n')

		if err != nil {
			log.Println(err)
		}

		for client := range hub.clients {
			select {
				case client.send <- message:
			}
		}
	}
}

func main() {
	flag.Parse()
	
	hub := newHub()
	go hub.run()
	
	reader := bufio.NewReader(os.Stdin)
	go broadcast_message(reader, hub)


	http.HandleFunc("/", serveHome)
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWs(hub, w, r)
	})
	err := http.ListenAndServe(*addr, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
