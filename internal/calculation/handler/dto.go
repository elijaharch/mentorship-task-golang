package handler

import (
	"time"

	"github.com/elijaharch/mentorship-task-golang/internal/calculation"
)

type calculationRequest struct {
	A         float64               `json:"a"`
	B         float64               `json:"b"`
	Operation calculation.Operation `json:"operation"`
}

func (r calculationRequest) toInput() calculation.Input {
	return calculation.Input{
		A:         r.A,
		B:         r.B,
		Operation: r.Operation,
	}
}

type calculationResponse struct {
	ID        int64                 `json:"id"`
	A         float64               `json:"a"`
	B         float64               `json:"b"`
	Operation calculation.Operation `json:"operation"`
	Result    float64               `json:"result"`
	CreatedAt time.Time             `json:"created_at"`
}

func newCalculationResponse(c calculation.Calculation) calculationResponse {
	return calculationResponse{
		ID:        c.ID,
		A:         c.A,
		B:         c.B,
		Operation: c.Operation,
		Result:    c.Result,
		CreatedAt: c.CreatedAt,
	}
}

type listResponse struct {
	Items  []calculationResponse `json:"items"`
	Limit  int                   `json:"limit"`
	Offset int                   `json:"offset"`
}

func newListResponse(
	calculations []calculation.Calculation,
	options calculation.ListOptions,
) listResponse {
	items := make([]calculationResponse, 0, len(calculations))

	for _, item := range calculations {
		items = append(items, newCalculationResponse(item))
	}

	return listResponse{
		Items:  items,
		Limit:  options.Limit,
		Offset: options.Offset,
	}
}

type errorDetail struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type errorResponse struct {
	Error errorDetail `json:"error"`
}
