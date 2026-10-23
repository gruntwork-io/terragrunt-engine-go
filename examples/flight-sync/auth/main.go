// Command flight-sync-auth stands in for the auth server the example engine
// exchanges credentials with. Every POST to /exchange returns a new random
// credential and prints a running count, so the count shows how many engines
// exchanged during a run.
package main

import (
	"crypto/rand"
	"flag"
	"log"
	"net/http"
	"sync/atomic"
	"time"
)

const readHeaderTimeout = 5 * time.Second

func main() {
	addr := flag.String("addr", "127.0.0.1:8080", "address to listen on")

	flag.Parse()

	var exchanges atomic.Int64

	mux := http.NewServeMux()
	mux.HandleFunc("POST /exchange", func(w http.ResponseWriter, _ *http.Request) {
		log.Printf("exchange %d", exchanges.Add(1))

		if _, err := w.Write([]byte(rand.Text())); err != nil {
			log.Printf("writing response: %v", err)
		}
	})

	server := &http.Server{
		Addr:              *addr,
		Handler:           mux,
		ReadHeaderTimeout: readHeaderTimeout,
	}

	log.Printf("listening on %s", *addr)
	log.Fatal(server.ListenAndServe())
}
