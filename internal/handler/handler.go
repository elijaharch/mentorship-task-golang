package handler

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	calculation "github.com/elijaharch/mentorship-task-golang/internal/domain"
)

type Service interface {
	Create(ctx context.Context, input calculation.Input) (calculation.Calculation, error)
	Get(ctx context.Context, id int64) (calculation.Calculation, error)
	Update(ctx context.Context, id int64, input calculation.Input) (calculation.Calculation, error)
	// List(ctx context.Context, options calculation.ListOptions) ([]calculation.Calculation, error)
	Delete(ctx context.Context, id int64) error
}

type Handler struct {
	svc    Service
	logger *slog.Logger
}

func New(svc Service, logger *slog.Logger) *Handler {
	return &Handler{
		svc:    svc,
		logger: logger,
	}
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var req calculationRequest

	if err := decodeJSON(w, r, &req); err != nil {
		h.writeError(
			w,
			r,
			http.StatusBadRequest,
			"invalid_body",
			"invalid request body",
		)
		return
	}

	calc, err := h.svc.Create(r.Context(), req.toInput())
	if err != nil {
		h.writeServiceError(w, r, "create calculation", err)
		return
	}

	if err := writeJSON(w, http.StatusCreated, newCalculationResponse(calc)); err != nil {
		h.logger.ErrorContext(
			r.Context(),
			"write create calculation response",
			"err",
			err,
		)
	}
}

func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	reqID, err := parseID(r, "id")
	if err != nil || reqID <= 0 {
		h.writeError(w, r, http.StatusBadRequest, "invalid_id", "invalid calculation id")
		return
	}

	calc, err := h.svc.Get(r.Context(), reqID)
	if err != nil {
		h.writeServiceError(w, r, "get calculation", err)
		return
	}

	if err := writeJSON(w, http.StatusOK, newCalculationResponse(calc)); err != nil {
		h.logger.ErrorContext(
			r.Context(),
			"write get calculation response",
			"err",
			err,
		)
	}
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	reqID, err := parseID(r, "id")
	if err != nil || reqID <= 0 {
		h.writeError(w, r, http.StatusBadRequest, "invalid_id", "invalid calculation id")
		return
	}

	var req calculationRequest

	if err = decodeJSON(w, r, &req); err != nil {
		h.writeError(
			w,
			r,
			http.StatusBadRequest,
			"invalid_body",
			"invalid request body",
		)
		return
	}

	calc, err := h.svc.Update(r.Context(), reqID, req.toInput())
	if err != nil {
		h.writeServiceError(w, r, "update calculation", err)
		return
	}

	if err := writeJSON(w, http.StatusOK, newCalculationResponse(calc)); err != nil {
		h.logger.ErrorContext(
			r.Context(),
			"write update calculation response",
			"err",
			err,
		)
	}
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	reqID, err := parseID(r, "id")
	if err != nil || reqID <= 0 {
		h.writeError(w, r, http.StatusBadRequest, "invalid_id", "invalid calculation id")
		return
	}

	err = h.svc.Delete(r.Context(), reqID)
	if err != nil {
		h.writeServiceError(w, r, "delete calculation", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func parseID(r *http.Request, key string) (int64, error) {
	raw := r.PathValue(key)
	n, err := strconv.ParseInt(raw, 10, 64)
	if err != nil || n <= 0 {
		return 0, errors.New("invalid id")
	}
	return n, nil
}
