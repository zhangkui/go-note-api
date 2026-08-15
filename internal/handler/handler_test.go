package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"go-note-api/internal/model"
	"go-note-api/internal/service"
	"go-note-api/internal/store"
)

func newTestHandler(t *testing.T) (*Handler, *store.Store) {
	t.Helper()
	st := store.New("")
	return New(service.New(st)), st
}

func do(h *Handler, method, target string, body interface{}) *httptest.ResponseRecorder {
	var buf []byte
	if body != nil {
		buf, _ = json.Marshal(body)
	}
	req := httptest.NewRequest(method, target, bytes.NewReader(buf))
	rr := httptest.NewRecorder()
	h.Routes().ServeHTTP(rr, req)
	return rr
}

func TestHandler_CreateAndList(t *testing.T) {
	h, _ := newTestHandler(t)
	rr := do(h, http.MethodPost, "/api/notes", model.Note{Title: "hello", Content: "world"})
	if rr.Code != http.StatusCreated {
		t.Fatalf("got %d, want %d: %s", rr.Code, http.StatusCreated, rr.Body.String())
	}
	var created model.Note
	if err := json.Unmarshal(rr.Body.Bytes(), &created); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if created.ID == "" {
		t.Fatal("expected non-empty id")
	}

	rr = do(h, http.MethodGet, "/api/notes", nil)
	if rr.Code != http.StatusOK {
		t.Fatalf("got %d, want %d", rr.Code, http.StatusOK)
	}
	var list []model.Note
	if err := json.Unmarshal(rr.Body.Bytes(), &list); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if len(list) != 1 {
		t.Fatalf("got %d, want 1", len(list))
	}
}

func TestHandler_GetNote(t *testing.T) {
	h, st := newTestHandler(t)
	created, err := st.Create(model.Note{Title: "hi"})
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	rr := do(h, http.MethodGet, "/api/notes/"+created.ID, nil)
	if rr.Code != http.StatusOK {
		t.Fatalf("got %d, want %d", rr.Code, http.StatusOK)
	}
	var got model.Note
	if err := json.Unmarshal(rr.Body.Bytes(), &got); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if got.Title != "hi" {
		t.Fatalf("got %q, want hi", got.Title)
	}
}

func TestHandler_Update(t *testing.T) {
	h, st := newTestHandler(t)
	created, err := st.Create(model.Note{Title: "old"})
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	rr := do(h, http.MethodPut, "/api/notes/"+created.ID, model.Note{Title: "new", Content: "c"})
	if rr.Code != http.StatusOK {
		t.Fatalf("got %d, want %d", rr.Code, http.StatusOK)
	}
}

func TestHandler_Delete(t *testing.T) {
	h, st := newTestHandler(t)
	created, err := st.Create(model.Note{Title: "t"})
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	rr := do(h, http.MethodDelete, "/api/notes/"+created.ID, nil)
	if rr.Code != http.StatusNoContent {
		t.Fatalf("got %d, want %d", rr.Code, http.StatusNoContent)
	}
}

func TestHandler_Search(t *testing.T) {
	h, st := newTestHandler(t)
	if _, err := st.Create(model.Note{Title: "go lang", Content: "x"}); err != nil {
		t.Fatalf("create: %v", err)
	}
	rr := do(h, http.MethodGet, "/api/notes/search?q=go", nil)
	if rr.Code != http.StatusOK {
		t.Fatalf("got %d, want %d", rr.Code, http.StatusOK)
	}
	var res []model.Note
	if err := json.Unmarshal(rr.Body.Bytes(), &res); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if len(res) != 1 {
		t.Fatalf("got %d, want 1", len(res))
	}
}

func TestHandler_Recent(t *testing.T) {
	h, st := newTestHandler(t)
	for i := 0; i < 3; i++ {
		if _, err := st.Create(model.Note{Title: "t"}); err != nil {
			t.Fatalf("create: %v", err)
		}
	}
	rr := do(h, http.MethodGet, "/api/notes/recent?limit=2", nil)
	if rr.Code != http.StatusOK {
		t.Fatalf("got %d, want %d", rr.Code, http.StatusOK)
	}
	var res []model.Note
	if err := json.Unmarshal(rr.Body.Bytes(), &res); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if len(res) != 2 {
		t.Fatalf("got %d, want 2", len(res))
	}
}

func TestHandler_CreateInvalidJSON(t *testing.T) {
	h, _ := newTestHandler(t)
	req := httptest.NewRequest(http.MethodPost, "/api/notes", bytes.NewReader([]byte("not json")))
	rr := httptest.NewRecorder()
	h.Routes().ServeHTTP(rr, req)
	if rr.Code != http.StatusBadRequest {
		t.Fatalf("got %d, want %d", rr.Code, http.StatusBadRequest)
	}
}

func TestHandler_MethodNotAllowed(t *testing.T) {
	h, _ := newTestHandler(t)
	rr := do(h, http.MethodPatch, "/api/notes", nil)
	if rr.Code != http.StatusMethodNotAllowed {
		t.Fatalf("got %d, want %d", rr.Code, http.StatusMethodNotAllowed)
	}
}

func TestHandler_UpdateNotFound(t *testing.T) {
	h, _ := newTestHandler(t)
	rr := do(h, http.MethodPut, "/api/notes/missing", model.Note{Title: "x"})
	if rr.Code != http.StatusNotFound {
		t.Fatalf("got %d, want %d", rr.Code, http.StatusNotFound)
	}
}

func TestHandler_DeleteNotFound(t *testing.T) {
	h, _ := newTestHandler(t)
	rr := do(h, http.MethodDelete, "/api/notes/missing", nil)
	if rr.Code != http.StatusNotFound {
		t.Fatalf("got %d, want %d", rr.Code, http.StatusNotFound)
	}
}
