package handler

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"strconv"

	calculation "github.com/elijaharch/mentorship-task-golang/internal/features/calculation/domain"
)

const maxRequestBodyBytes = 1 << 20

type Service interface {
	Create(ctx context.Context, input calculation.Input) (calculation.Calculation, error)
	Get(ctx context.Context, id int64) (calculation.Calculation, error)
	Update(ctx context.Context, id int64, input calculation.Input) (calculation.Calculation, error)
	// List(ctx context.Context, options calculation.ListOptions) ([]calculation.Calculation, error)
	// Delete(ctx context.Context, id int64) error
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

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, maxRequestBodyBytes)

	var req calculationRequest

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&req); err != nil {
		h.respondError(
			w,
			r,
			http.StatusBadRequest,
			"invalid_body",
			"invalid request body",
		)
		return
	}
	// after 1st obj only whitespaces allowed
	if err := decoder.Decode(&struct{}{}); !errors.Is(err, io.EOF) {
		h.respondError(
			w,
			r,
			http.StatusBadRequest,
			"invalid_body",
			"request body must contain a single JSON object",
		)
		return
	}

	calc, err := h.service.Create(r.Context(), req.toInput())
	if err != nil {
		h.respondServiceError(w, r, "create calculation", err)
		return
	}

	if err := writeJSON(
		w,
		http.StatusCreated,
		newCalculationResponse(calc),
	); err != nil {
		h.logger.ErrorContext(
			r.Context(),
			"write create calculation response",
			"err",
			err,
		)
	}
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	reqID, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil || reqID <= 0 {
		h.respondError(
			w,
			r,
			http.StatusBadRequest,
			"invalid_id",
			"invalid calculation id",
		)
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, maxRequestBodyBytes)

	var req calculationRequest
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&req); err != nil {
		h.respondError(
			w,
			r,
			http.StatusBadRequest,
			"invalid_body",
			"invalid request body",
		)
		return
	}

	if err := decoder.Decode(&struct{}{}); !errors.Is(err, io.EOF) {
		h.respondError(
			w,
			r,
			http.StatusBadRequest,
			"invalid_body",
			"request body must contain a single JSON object",
		)
		return
	}

	calc, err := h.service.Update(r.Context(), reqID, req.toInput())
	if err != nil {
		h.respondServiceError(w, r, "update calculation", err)
		return
	}

	if err := writeJSON(
		w,
		http.StatusOK,
		newCalculationResponse(calc),
	); err != nil {
		h.logger.ErrorContext(
			r.Context(),
			"write update calculation response",
			"err",
			err,
		)
	}
}

func (h *Handler) respondServiceError(
	w http.ResponseWriter,
	r *http.Request,
	action string,
	err error,
) {
	switch {
	case errors.Is(err, calculation.ErrNotFound):
		h.respondError(
			w,
			r,
			http.StatusNotFound,
			"not_found",
			err.Error(),
		)
	case errors.Is(err, calculation.ErrInvalidOperation):
		h.respondError(
			w,
			r,
			http.StatusUnprocessableEntity,
			"invalid_operation",
			err.Error(),
		)
	case errors.Is(err, calculation.ErrDivisionByZero):
		h.respondError(
			w,
			r,
			http.StatusUnprocessableEntity,
			"division_by_zero",
			err.Error(),
		)
	case errors.Is(err, calculation.ErrInvalidNumber):
		h.respondError(
			w,
			r,
			http.StatusUnprocessableEntity,
			"invalid_number",
			err.Error(),
		)
	default:
		h.logger.ErrorContext(
			r.Context(),
			action+" failed",
			"err",
			err,
		)

		h.respondError(
			w,
			r,
			http.StatusInternalServerError,
			"internal_error",
			"internal server error",
		)
	}
}

func writeJSON(w http.ResponseWriter, status int, value any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	return json.NewEncoder(w).Encode(value)
}

func (h *Handler) respondError(
	w http.ResponseWriter,
	r *http.Request,
	status int,
	code string,
	message string,
) {
	response := errorResponse{
		Error: errorDetail{
			Code:    code,
			Message: message,
		},
	}

	if err := writeJSON(w, status, response); err != nil {
		h.logger.ErrorContext(
			r.Context(),
			"write error response",
			"err",
			err,
		)
	}
}
