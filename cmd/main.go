// Command go-note-api runs the notes management service.
package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go-note-api/internal/handler"
	"go-note-api/internal/server"
	"go-note-api/internal/service"
	"go-note-api/internal/store"
)

func main() {
	addr := flag.String("addr", ":8080", "HTTP listen address")
	dbFile := flag.String("db", "notes.json", "path to the persistence file")
	flag.Parse()

	st := store.New(*dbFile)
	svc := service.New(st)
	h := handler.New(svc)
	srv := server.New(*addr, h)

	go func() {
		log.Printf("go-note-api listening on %s", *addr)
		if err := srv.Start(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server error: %v", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("shutdown error: %v", err)
	}
}
