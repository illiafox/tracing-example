package storage

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v9"
	"github.com/google/uuid"
	"trace-example/models"
)

type NotesStorage struct {
	client redis.UniversalClient
}

func NewNotesStorage(client redis.UniversalClient) NotesStorage {
	return NotesStorage{client: client}
}

func (s NotesStorage) Store(ctx context.Context, note models.Note) error {
	data, err := json.Marshal(note)
	if err != nil {
		return fmt.Errorf("marshal note: %w", err)
	}

	if err = s.client.Set(ctx, note.NoteID.String(), data, -1).Err(); err != nil {
		return fmt.Errorf("redis set: %w", err)
	}

	return nil
}

func (s NotesStorage) Get(ctx context.Context, noteID uuid.UUID) (*models.Note, error) {
	data, err := s.client.Get(ctx, noteID.String()).Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil, models.ErrNotFound
		}
		return nil, fmt.Errorf("redis get: %w", err)
	}

	var note models.Note
	if err = json.Unmarshal(data, &note); err != nil {
		return nil, fmt.Errorf("unmarshal note: %w", err)
	}

	return &note, nil
}
