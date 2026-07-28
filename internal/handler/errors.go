package handler

import (
	"errors"
	"net/http"

	calculation "github.com/elijaharch/mentorship-task-golang/internal/domain"
)

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
