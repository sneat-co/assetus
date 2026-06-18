// Command assetusd is the assetus backend service.
//
// This entrypoint exposes a health endpoint; the assetus domain capabilities are
// served as a Sneat space-module extension (see package assetusext) mounted by
// the host Sneat backend.
package main

import (
	"log"
	"net/http"
	"os"

	"github.com/sneat-co/assetus/backend/internal/health"
)

func main() {
	addr := os.Getenv("ASSETUS_ADDR")
	if addr == "" {
		addr = ":8080"
	}

	mux := http.NewServeMux()
	mux.Handle("GET /health", health.Handler())

	log.Printf("assetusd listening on %s", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("assetusd failed: %v", err)
	}
}
