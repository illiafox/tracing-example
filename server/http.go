package server

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"time"
	"trace-example/models"
	"trace-example/storage"
)

type FiberHandler struct {
	notesStorage storage.NotesStorage
}

func NewFiberHandler(notesStorage storage.NotesStorage) FiberHandler {
	return FiberHandler{notesStorage: notesStorage}
}

func (h FiberHandler) CreateNote(fiberctx *fiber.Ctx) error {
	input := struct {
		Title   string `json:"title"`
		Content string `json:"content"`
	}{}

	if err := fiberctx.BodyParser(&input); err != nil {
		return err
	}

	noteID := uuid.New()
	err := h.notesStorage.Store(fiberctx.UserContext(), models.Note{
		NoteID:  noteID,
		Title:   input.Title,
		Content: input.Content,
		Created: time.Now(),
	})
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return fiberctx.JSON(map[string]any{
		"note_id": noteID,
	})
}

func (h FiberHandler) GetNote(fiberctx *fiber.Ctx) error {
	noteID, err := uuid.Parse(fiberctx.Query("note_id"))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	note, err := h.notesStorage.Get(fiberctx.UserContext(), noteID)
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			return fiber.NewError(fiber.StatusNotFound, err.Error())
		}
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return fiberctx.JSON(note)
}
