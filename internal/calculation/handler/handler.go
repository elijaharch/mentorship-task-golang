package handler

import (
	"context"
	"log/slog"

	"github.com/elijaharch/mentorship-task-golang/internal/calculation"
)

type Service interface {
	Create(ctx context.Context, input calculation.Input) (calculation.Calculation, error)
	Get(ctx context.Context, id int64) (calculation.Calculation, error)
	Update(ctx context.Context, id int64, input calculation.Input) (calculation.Calculation, error)
	List(ctx context.Context, options calculation.ListOptions) ([]calculation.Calculation, error)
	Delete(ctx context.Context, id int64) error
}

type Handler struct {
	service Service
	logger  *slog.Logger
}

func New(service Service, logger *slog.Logger) *Handler {
	return &Handler{
		service: service,
		logger:  logger,
	}
}
