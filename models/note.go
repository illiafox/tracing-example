package models

import (
	"errors"
	"time"
)
import "github.com/google/uuid"

var ErrNotFound = errors.New("note not found")

type Note struct {
	NoteID  uuid.UUID `json:"note_id"`
	Title   string    `json:"title"`
	Content string    `json:"content"`
	Created time.Time `json:"created"`
}
