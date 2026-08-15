package server

import (
	"context"
	"net/http/httptest"
	"testing"

	"go-note-api/internal/handler"
	"go-note-api/internal/service"
	"go-note-api/internal/store"
)

func TestServer_New(t *testing.T) {
	srv := New("127.0.0.1:0", handler.New(service.New(store.New(""))))
	if srv == nil {
		t.Fatal("expected non-nil server")
	}
}

func TestServer_Integration(t *testing.T) {
	h := handler.New(service.New(store.New("")))
	ts := httptest.NewServer(h.Routes())
	defer ts.Close()

	resp, err := ts.Client().Get(ts.URL + "/api/notes")
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		t.Fatalf("got %d, want 200", resp.StatusCode)
	}
}

func TestServer_Shutdown(t *testing.T) {
	srv := New("127.0.0.1:0", handler.New(service.New(store.New(""))))
	if err := srv.Shutdown(context.Background()); err != nil {
		t.Fatalf("shutdown: %v", err)
	}
}
