package service

import (
	"context"
	"testing"

	"go-note-api/internal/model"
	"go-note-api/internal/store"
)

func newTestService(t *testing.T) (*Service, *store.Store) {
	t.Helper()
	st := store.New("")
	return New(st), st
}

func TestService_Create(t *testing.T) {
	svc, _ := newTestService(t)
	created, err := svc.Create(context.Background(), model.Note{Title: "hello", Content: "world"})
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	if created.ID == "" {
		t.Fatal("expected non-empty id")
	}
}

func TestService_GetAndUpdate(t *testing.T) {
	svc, _ := newTestService(t)
	created, err := svc.Create(context.Background(), model.Note{Title: "old"})
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	if _, err := svc.Get(context.Background(), created.ID); err != nil {
		t.Fatalf("get: %v", err)
	}
	updated, err := svc.Update(context.Background(), created.ID, model.Note{Title: "new"})
	if err != nil {
		t.Fatalf("update: %v", err)
	}
	if updated.Title != "new" {
		t.Fatalf("got %q, want new", updated.Title)
	}
}

func TestService_Delete(t *testing.T) {
	svc, _ := newTestService(t)
	created, err := svc.Create(context.Background(), model.Note{Title: "t"})
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	if err := svc.Delete(context.Background(), created.ID); err != nil {
		t.Fatalf("delete: %v", err)
	}
}

func TestService_Search(t *testing.T) {
	svc, _ := newTestService(t)
	if _, err := svc.Create(context.Background(), model.Note{Title: "alpha", Content: "beta"}); err != nil {
		t.Fatalf("create: %v", err)
	}
	res, err := svc.Search(context.Background(), "alpha")
	if err != nil {
		t.Fatalf("search: %v", err)
	}
	if len(res) != 1 {
		t.Fatalf("got %d, want 1", len(res))
	}
}

func TestService_Recent(t *testing.T) {
	svc, _ := newTestService(t)
	for i := 0; i < 3; i++ {
		if _, err := svc.Create(context.Background(), model.Note{Title: "t"}); err != nil {
			t.Fatalf("create: %v", err)
		}
	}
	if got := svc.Recent(2); len(got) != 2 {
		t.Fatalf("got %d, want 2", len(got))
	}
}
