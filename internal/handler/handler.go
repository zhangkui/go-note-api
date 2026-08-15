// Package handler exposes the note service over HTTP.
package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"go-note-api/internal/model"
	"go-note-api/internal/service"
	"go-note-api/internal/store"
)

// Handler translates HTTP requests into service calls.
type Handler struct {
	svc *service.Service
}

// New returns a Handler backed by the given service.
func New(svc *service.Service) *Handler {
	return &Handler{svc: svc}
}

// Routes builds the request multiplexer for the notes API.
func (h *Handler) Routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/notes", h.collection)
	mux.HandleFunc("/api/notes/", h.item)
	return mux
}

func (h *Handler) collection(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.List(w, r)
	case http.MethodPost:
		h.Create(w, r)
	default:
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
	}
}

func (h *Handler) item(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/api/notes/")
	switch {
	case path == "search":
		if r.Method == http.MethodGet {
			h.Search(w, r)
		} else {
			writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		}
	case path == "recent":
		if r.Method == http.MethodGet {
			h.Recent(w, r)
		} else {
			writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		}
	case path == "":
		if r.Method == http.MethodGet {
			h.List(w, r)
		} else {
			writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		}
	default:
		switch r.Method {
		case http.MethodGet:
			h.GetNote(w, r, path)
		case http.MethodPut:
			h.Update(w, r, path)
		case http.MethodDelete:
			h.Delete(w, r, path)
		default:
			writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		}
	}
}

// Create handles POST /api/notes.
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var n model.Note
	if err := json.NewDecoder(r.Body).Decode(&n); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	created, err := h.svc.Create(r.Context(), n)
	if err != nil {
		if errors.Is(err, model.ErrEmptyTitle) {
			writeError(w, http.StatusBadRequest, err.Error())
			return
		}
		writeError(w, http.StatusInternalServerError, "internal error")
		return
	}
	writeJSON(w, http.StatusCreated, created)
}

// List handles GET /api/notes.
func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, h.svc.List())
}

// Recent handles GET /api/notes/recent?limit=N.
func (h *Handler) Recent(w http.ResponseWriter, r *http.Request) {
	limit := 10
	if q := r.URL.Query().Get("limit"); q != "" {
		if v, err := strconv.Atoi(q); err == nil {
			limit = v
		}
	}
	writeJSON(w, http.StatusOK, h.svc.Recent(limit))
}

// Search handles GET /api/notes/search?q=...
func (h *Handler) Search(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query().Get("q")
	results, err := h.svc.Search(r.Context(), q)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "search failed")
		return
	}
	writeJSON(w, http.StatusOK, results)
}

// GetNote handles GET /api/notes/{id}.
func (h *Handler) GetNote(w http.ResponseWriter, r *http.Request, id string) {
	note, _ := h.svc.Get(r.Context(), id)
	resp := map[string]interface{}{
		"id":         note.ID,
		"title":      note.Title,
		"content":    note.Content,
		"tags":       note.Tags,
		"created_at": note.CreatedAt,
		"updated_at": note.UpdatedAt,
	}
	writeJSON(w, http.StatusOK, resp)
}

// Update handles PUT /api/notes/{id}.
func (h *Handler) Update(w http.ResponseWriter, r *http.Request, id string) {
	var n model.Note
	if err := json.NewDecoder(r.Body).Decode(&n); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	updated, err := h.svc.Update(r.Context(), id, n)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			writeError(w, http.StatusNotFound, "note not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "internal error")
		return
	}
	writeJSON(w, http.StatusOK, updated)
}

// Delete handles DELETE /api/notes/{id}.
func (h *Handler) Delete(w http.ResponseWriter, r *http.Request, id string) {
	if err := h.svc.Delete(r.Context(), id); err != nil {
		if errors.Is(err, store.ErrNotFound) {
			writeError(w, http.StatusNotFound, "note not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "internal error")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func writeJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}
