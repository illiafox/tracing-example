package server

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel/trace"
	"time"
	"trace-example/models"
	"trace-example/storage"
)

type FiberHandler struct {
	notesStorage storage.NotesStorage
	tracer       trace.Tracer
}

func NewFiberHandler(notesStorage storage.NotesStorage, tracer trace.Tracer) FiberHandler {
	return FiberHandler{notesStorage: notesStorage, tracer: tracer}
}

func (h FiberHandler) CreateNote(fiberctx *fiber.Ctx) error {
	ctx, span := h.tracer.Start(fiberctx.UserContext(), "GetNote")
	defer span.End()

	input := struct {
		Title   string `json:"title"`
		Content string `json:"content"`
	}{}

	span.AddEvent("parse body")
	if err := fiberctx.BodyParser(&input); err != nil {
		return err
	}

	noteID := uuid.New()
	span.AddEvent("call redis")
	err := h.notesStorage.Store(ctx, models.Note{
		NoteID:  noteID,
		Title:   input.Title,
		Content: input.Content,
		Created: time.Now(),
	})
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	span.AddEvent("write note_id")
	return fiberctx.JSON(map[string]any{
		"note_id": noteID,
	})
}

func (h FiberHandler) GetNote(fiberctx *fiber.Ctx) error {
	ctx, span := h.tracer.Start(fiberctx.UserContext(), "GetNote")
	defer span.End()

	span.AddEvent("parse note_id")
	noteID, err := uuid.Parse(fiberctx.Query("note_id"))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	span.AddEvent("call redis storage")
	note, err := h.notesStorage.Get(ctx, noteID)
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			return fiber.NewError(fiber.StatusNotFound, err.Error())
		}
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	span.AddEvent("write note in json")
	return fiberctx.JSON(note)
}
