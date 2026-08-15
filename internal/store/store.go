// Package store provides a file-backed implementation of note persistence.
package store

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"go-note-api/internal/model"
)

// ErrNotFound is returned when a note cannot be located by its identifier.
var ErrNotFound = errors.New("note not found")

// Store keeps notes in memory and optionally mirrors them to a JSON file so
// that data survives process restarts.
type Store struct {
	mu     sync.RWMutex
	notes  map[string]*model.Note
	nextID int64
	file   string
}

// New creates a Store. When file is non-empty the notes are loaded from that
// path and every mutation is flushed back to it.
func New(file string) *Store {
	s := &Store{
		notes: make(map[string]*model.Note),
		file:  file,
	}
	if file != "" {
		s.load()
	}
	return s
}

// load reads previously persisted notes from disk. Missing or unreadable files
// are treated as an empty store.
func (s *Store) load() {
	data, err := os.ReadFile(s.file)
	if err != nil {
		return
	}
	var notes map[string]*model.Note
	if err := json.Unmarshal(data, &notes); err != nil {
		return
	}
	var maxID int64
	for id, n := range notes {
		if n == nil {
			delete(notes, id)
			continue
		}
		if v, err := strconv.ParseInt(id, 10, 64); err == nil && v > maxID {
			maxID = v
		}
	}
	s.notes = notes
	s.nextID = maxID
}

// Create inserts a new note, assigning it the next available identifier.
func (s *Store) Create(n model.Note) (*model.Note, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.nextID++
	id := strconv.FormatInt(s.nextID, 10)
	now := time.Now()
	n.ID = id
	n.CreatedAt = now
	n.UpdatedAt = now
	cp := n
	s.notes[id] = &cp
	s.persistLocked()

	result := cp
	return &result, nil
}

// Get returns a copy of the note with the given identifier.
func (s *Store) Get(id string) (*model.Note, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	n, ok := s.notes[id]
	if !ok {
		return nil, ErrNotFound
	}
	cp := *n
	return &cp, nil
}

// List returns every note currently stored.
func (s *Store) List() []model.Note {
	s.mu.RLock()
	defer s.mu.RUnlock()

	notes := make([]model.Note, 0, len(s.notes))
	for _, n := range s.notes {
		notes = append(notes, *n)
	}
	return notes
}

// Recent returns up to limit notes ordered from newest to oldest.
func (s *Store) Recent(limit int) []model.Note {
	s.mu.RLock()
	defer s.mu.RUnlock()

	all := make([]model.Note, 0, len(s.notes))
	for _, n := range s.notes {
		all = append(all, *n)
	}
	sort.Slice(all, func(i, j int) bool {
		return all[i].CreatedAt.After(all[j].CreatedAt)
	})
	if limit < 0 {
		limit = 0
	}
	return all[:limit]
}

// Update overwrites the mutable fields of an existing note.
func (s *Store) Update(id string, n model.Note) (*model.Note, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	existing, ok := s.notes[id]
	if !ok {
		return nil, ErrNotFound
	}
	existing.Title = n.Title
	existing.Content = n.Content
	existing.Tags = n.Tags
	existing.UpdatedAt = time.Now()
	s.persistLocked()

	cp := *existing
	return &cp, nil
}

// Delete removes the note with the given identifier.
func (s *Store) Delete(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.notes[id]; !ok {
		return ErrNotFound
	}
	delete(s.notes, id)
	s.persistLocked()
	return nil
}

// Search returns notes whose title or content contain query (case-insensitive).
// It honours the provided context and aborts early when it has been cancelled.
func (s *Store) Search(ctx context.Context, query string) ([]model.Note, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	s.mu.RLock()
	defer s.mu.RUnlock()

	q := strings.ToLower(query)
	var results []model.Note
	for _, n := range s.notes {
		if strings.Contains(strings.ToLower(n.Title), q) ||
			strings.Contains(strings.ToLower(n.Content), q) {
			results = append(results, *n)
		}
	}
	return results, nil
}

// persistLocked writes the current state to disk. It must be called while
// holding the write lock.
func (s *Store) persistLocked() {
	if s.file == "" {
		return
	}
	data, err := json.Marshal(s.notes)
	if err != nil {
		return
	}
	_ = os.WriteFile(s.file, data, 0o644)
}
