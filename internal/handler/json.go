package handler

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	calculation "github.com/elijaharch/mentorship-task-golang/internal/domain"
)

const maxRequestBodyBytes = 1 << 20

func writeJSON(w http.ResponseWriter, status int, value any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if value == nil {
		return nil
	}

	return json.NewEncoder(w).Encode(value)
}

func decodeJSON(w http.ResponseWriter, r *http.Request, dst any) error {
	r.Body = http.MaxBytesReader(w, r.Body, maxRequestBodyBytes)

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(dst); err != nil {
		if errors.Is(err, io.EOF) {
			return errors.New("request body is empty")
		}
		return errors.New("invalid json")
	}

	if err := dec.Decode(&struct{}{}); !errors.Is(err, io.EOF) {
		return errors.New("request body must contain a single JSON object")
	}

	return nil
}

func (h *Handler) writeError(
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

func (h *Handler) writeServiceError(
	w http.ResponseWriter,
	r *http.Request,
	action string,
	err error,
) {
	switch {
	case errors.Is(err, calculation.ErrNotFound):
		h.writeError(
			w,
			r,
			http.StatusNotFound,
			"not_found",
			err.Error(),
		)
	case errors.Is(err, calculation.ErrInvalidOperation):
		h.writeError(
			w,
			r,
			http.StatusUnprocessableEntity,
			"invalid_operation",
			err.Error(),
		)
	case errors.Is(err, calculation.ErrDivisionByZero):
		h.writeError(
			w,
			r,
			http.StatusUnprocessableEntity,
			"division_by_zero",
			err.Error(),
		)
	case errors.Is(err, calculation.ErrInvalidNumber):
		h.writeError(
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

		h.writeError(
			w,
			r,
			http.StatusInternalServerError,
			"internal_error",
			"internal server error",
		)
	}
}
