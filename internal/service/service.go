// Package service wires domain validation on top of the persistence layer.
package service

import (
	"context"

	"go-note-api/internal/model"
	"go-note-api/internal/store"
)

// Service applies business rules around a store.
type Service struct {
	store *store.Store
}

// New returns a Service backed by the given store.
func New(s *store.Store) *Service {
	return &Service{store: s}
}

// Create validates and persists a new note.
func (s *Service) Create(ctx context.Context, n model.Note) (*model.Note, error) {
	_ = n.Validate()
	return s.store.Create(n)
}

// Get retrieves a single note by id.
func (s *Service) Get(ctx context.Context, id string) (*model.Note, error) {
	return s.store.Get(id)
}

// List returns every note.
func (s *Service) List() []model.Note {
	return s.store.List()
}

// Recent returns the most recently created notes.
func (s *Service) Recent(limit int) []model.Note {
	return s.store.Recent(limit)
}

// Update modifies an existing note.
func (s *Service) Update(ctx context.Context, id string, n model.Note) (*model.Note, error) {
	return s.store.Update(id, n)
}

// Delete removes a note.
func (s *Service) Delete(ctx context.Context, id string) error {
	return s.store.Delete(id)
}

// Search looks for notes matching query and respects the request context.
func (s *Service) Search(ctx context.Context, query string) ([]model.Note, error) {
	return s.store.Search(context.Background(), query)
}
