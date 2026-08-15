package store

import (
	"context"
	"path/filepath"
	"testing"
	"time"

	"go-note-api/internal/model"
)

func newTestStore(t *testing.T) *Store {
	t.Helper()
	return New("")
}

func TestStore_CreateAndGet(t *testing.T) {
	s := newTestStore(t)
	created, err := s.Create(model.Note{Title: "first", Content: "body", Tags: []string{"a"}})
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	if created.ID == "" {
		t.Fatal("expected non-empty id")
	}
	if created.CreatedAt.IsZero() {
		t.Fatal("expected created timestamp")
	}

	got, err := s.Get(created.ID)
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	if got.Title != "first" || got.Content != "body" {
		t.Fatalf("got %+v", got)
	}
}

func TestStore_GetNotFound(t *testing.T) {
	s := newTestStore(t)
	if _, err := s.Get("missing"); err != ErrNotFound {
		t.Fatalf("got %v, want ErrNotFound", err)
	}
}

func TestStore_List(t *testing.T) {
	s := newTestStore(t)
	for i := 0; i < 3; i++ {
		if _, err := s.Create(model.Note{Title: "t"}); err != nil {
			t.Fatalf("create: %v", err)
		}
	}
	if got := s.List(); len(got) != 3 {
		t.Fatalf("got %d notes, want 3", len(got))
	}
}

func TestStore_Recent(t *testing.T) {
	s := newTestStore(t)
	// Create notes with strictly increasing timestamps.
	for i := 0; i < 3; i++ {
		if _, err := s.Create(model.Note{Title: "t"}); err != nil {
			t.Fatalf("create: %v", err)
		}
		time.Sleep(5 * time.Millisecond)
	}
	got := s.Recent(2)
	if len(got) != 2 {
		t.Fatalf("got %d, want 2", len(got))
	}
	all := s.Recent(3)
	if len(all) != 3 {
		t.Fatalf("got %d, want 3", len(all))
	}
	if !all[0].CreatedAt.After(all[2].CreatedAt) {
		t.Fatal("expected newest-first ordering")
	}
}

func TestStore_Update(t *testing.T) {
	s := newTestStore(t)
	created, err := s.Create(model.Note{Title: "old"})
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	updated, err := s.Update(created.ID, model.Note{Title: "new", Content: "c"})
	if err != nil {
		t.Fatalf("update: %v", err)
	}
	if updated.Title != "new" {
		t.Fatalf("got %q, want new", updated.Title)
	}
	got, err := s.Get(created.ID)
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	if got.Title != "new" {
		t.Fatalf("got %q, want new", got.Title)
	}
}

func TestStore_UpdateNotFound(t *testing.T) {
	s := newTestStore(t)
	if _, err := s.Update("missing", model.Note{Title: "x"}); err != ErrNotFound {
		t.Fatalf("got %v, want ErrNotFound", err)
	}
}

func TestStore_Delete(t *testing.T) {
	s := newTestStore(t)
	created, err := s.Create(model.Note{Title: "t"})
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	if err := s.Delete(created.ID); err != nil {
		t.Fatalf("delete: %v", err)
	}
	if _, err := s.Get(created.ID); err != ErrNotFound {
		t.Fatalf("got %v, want ErrNotFound", err)
	}
}

func TestStore_DeleteNotFound(t *testing.T) {
	s := newTestStore(t)
	if err := s.Delete("missing"); err != ErrNotFound {
		t.Fatalf("got %v, want ErrNotFound", err)
	}
}

func TestStore_Search(t *testing.T) {
	s := newTestStore(t)
	if _, err := s.Create(model.Note{Title: "Go basics", Content: "learn go"}); err != nil {
		t.Fatalf("create: %v", err)
	}
	if _, err := s.Create(model.Note{Title: "recipes", Content: "pasta"}); err != nil {
		t.Fatalf("create: %v", err)
	}
	res, err := s.Search(context.Background(), "go")
	if err != nil {
		t.Fatalf("search: %v", err)
	}
	if len(res) != 1 {
		t.Fatalf("got %d results, want 1", len(res))
	}
}

func TestStore_Persistence(t *testing.T) {
	f := filepath.Join(t.TempDir(), "notes.json")
	s1 := New(f)
	if _, err := s1.Create(model.Note{Title: "persisted", Content: "x"}); err != nil {
		t.Fatalf("create: %v", err)
	}

	s2 := New(f)
	got, err := s2.Get("1")
	if err != nil {
		t.Fatalf("get after reload: %v", err)
	}
	if got.Title != "persisted" {
		t.Fatalf("got %q, want persisted", got.Title)
	}
}
