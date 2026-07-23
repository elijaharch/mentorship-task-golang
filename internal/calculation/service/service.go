package service

import (
	"context"
	"math"

	"github.com/elijaharch/mentorship-task-golang/internal/calculation"
)

type Repository interface {
	Create(ctx context.Context, calc calculation.Calculation) (calculation.Calculation, error)
	Get(ctx context.Context, id int64) (calculation.Calculation, error)
	Update(ctx context.Context, id int64, calc calculation.Calculation) (calculation.Calculation, error)
}

type Service struct {
	repo Repository
}

func New(repo Repository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) Create(ctx context.Context, input calculation.Input) (calculation.Calculation, error) {
	calc, err := newCalculation(input)
	if err != nil {
		return calculation.Calculation{}, err
	}

	return s.repo.Create(ctx, calc)
}

func (s *Service) Get(ctx context.Context, id int64) (calculation.Calculation, error) {
	return s.repo.Get(ctx, id)
}

func (s *Service) Update(ctx context.Context, id int64, input calculation.Input) (calculation.Calculation, error) {
	calc, err := newCalculation(input)
	if err != nil {
		return calculation.Calculation{}, err
	}

	return s.repo.Update(ctx, id, calc)
}

func validate(input calculation.Input) error {
	if !input.Operation.Valid() {
		return calculation.ErrInvalidOperation
	}
	if math.IsNaN(input.A) || math.IsInf(input.A, 0) || math.IsNaN(input.B) || math.IsInf(input.B, 0) {
		return calculation.ErrInvalidNumber
	}
	if input.Operation == calculation.OperationDivide && input.B == 0 {
		return calculation.ErrDivisionByZero
	}

	return nil
}

func newCalculation(input calculation.Input) (calculation.Calculation, error) {
	if err := validate(input); err != nil {
		return calculation.Calculation{}, err
	}

	var result float64
	switch input.Operation {
	case calculation.OperationAdd:
		result = input.A + input.B
	case calculation.OperationSubtract:
		result = input.A - input.B
	case calculation.OperationMultiply:
		result = input.A * input.B
	case calculation.OperationDivide:
		result = input.A / input.B
	default:
		return calculation.Calculation{}, calculation.ErrInvalidOperation
	}

	if math.IsNaN(result) || math.IsInf(result, 0) {
		return calculation.Calculation{}, calculation.ErrInvalidNumber
	}

	calc := calculation.Calculation{
		A:         input.A,
		B:         input.B,
		Operation: input.Operation,
		Result:    result,
	}

	return calc, nil
}
