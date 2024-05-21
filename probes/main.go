package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

func newServer() *http.Server {
	port := os.Getenv("PORT")
	mux := http.NewServeMux()
	s := &http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: mux,
	}
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.URL.Path)
		w.Write([]byte("Hello, World!"))
	})
	mux.HandleFunc("/probes/liveness", func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.URL.Path)
		w.Write([]byte("OK"))
	})
	mux.HandleFunc("/probes/liveness/fail", func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.URL.Path)
		w.WriteHeader(http.StatusInternalServerError)
	})
	mux.HandleFunc("/probes/readiness", func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.URL.Path)
		w.Write([]byte("OK"))
	})
	mux.HandleFunc("/probes/readiness/fail", func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.URL.Path)
		w.WriteHeader(http.StatusInternalServerError)
	})
	var probeCount int32 = 0
	var lock sync.Mutex
	incrementProbe := func() int32 {
		lock.Lock()
		defer lock.Unlock()
		if probeCount == 8 {
			probeCount = 0
		} else {
			probeCount++
		}
		return probeCount
	}
	mux.HandleFunc("/probes/readiness/count", func(w http.ResponseWriter, r *http.Request) {
		count := incrementProbe()
		log.Println(r.URL.Path, count)
		switch count / 4 {
		case 1:
			log.Println("DOWN")
			w.WriteHeader(http.StatusInternalServerError)
		default:
			log.Println("OK")
			w.Write([]byte("OK"))
		}
	})
	mux.HandleFunc("/probes/readiness/slow", func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.URL.Path)
		time.Sleep(3 * time.Second)
		w.Write([]byte("OK"))
	})
	mux.HandleFunc("/probes/startup", func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.URL.Path)
		w.Write([]byte("OK"))
	})
	mux.HandleFunc("/probes/startup/fail", func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.URL.Path)
		w.WriteHeader(http.StatusInternalServerError)
	})
	mux.HandleFunc("/probes/startup/slow", func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.URL.Path)
		time.Sleep(5 * time.Second)
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
