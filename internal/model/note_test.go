package model

import "testing"

func TestNote_Validate(t *testing.T) {
	tests := []struct {
		name    string
		note    Note
		wantErr error
	}{
		{"valid", Note{Title: "hello"}, nil},
		{"empty title", Note{Title: ""}, ErrEmptyTitle},
		{"whitespace title", Note{Title: "   "}, ErrEmptyTitle},
		{"content ignored", Note{Title: "t", Content: "anything"}, nil},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.note.Validate()
			if err != tc.wantErr {
				t.Fatalf("got %v, want %v", err, tc.wantErr)
			}
		})
	}
}
