// Package model defines the core Note domain type and its validation rules.
package model

import (
	"errors"
	"strings"
	"time"
)

// Domain-level sentinel errors.
var (
	ErrEmptyTitle = errors.New("note title must not be empty")
)

// Note is the central data structure managed by the service.
type Note struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	Tags      []string  `json:"tags,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Validate checks that a note satisfies the basic invariants before it is
// persisted. It returns ErrEmptyTitle when the title is blank.
func (n Note) Validate() error {
	if strings.TrimSpace(n.Title) == "" {
		return ErrEmptyTitle
	}
	return nil
}
