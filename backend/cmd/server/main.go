package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"calculator/internal/api"
)

func main() {
	address := os.Getenv("ADDR")
	if address == "" {
		address = "127.0.0.1:8080"
	}
	server := &http.Server{
		Addr: address, Handler: api.NewHandler(),
		ReadHeaderTimeout: 5 * time.Second, ReadTimeout: 10 * time.Second,
		WriteTimeout: 10 * time.Second, IdleTimeout: 60 * time.Second, MaxHeaderBytes: 16 * 1024,
	}
	log.Printf("Calculator API listening on %s", address)
	log.Fatal(server.ListenAndServe())
}
