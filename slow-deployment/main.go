package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func newServer() *http.Server {
	isCrash := os.Getenv("CRASH")
	if isCrash == "true" {
		time.Sleep(20 * time.Second)
		log.Println("Crashing...")
		os.Exit(1)
	}

	port := os.Getenv("PORT")
	mux := http.NewServeMux()
	s := &http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: mux,
	}
	mux.HandleFunc("/probes/startup", func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.URL.Path)
		w.Write([]byte("OK"))
	})
	return s
}

func mainInternal() error {
	s := newServer()
	go func() {
		if err := s.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()

	log.Println("Server is running on", s.Addr)
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()
	<-ctx.Done()

	ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	return s.Shutdown(ctx)
}

func main() {
	if err := mainInternal(); err != nil {
		log.Fatal(err)
	}
}
