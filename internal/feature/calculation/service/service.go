package service

import (
	"context"
	"math"

	"github.com/elijaharch/mentorship-task-golang/internal/feature/calculation/domain"
)

type CalculationRepository interface {
	Create(ctx context.Context, calc domain.Calculation) (domain.Calculation, error)
	Get(ctx context.Context, id int64) (domain.Calculation, error)
	Update(ctx context.Context, id int64, calc domain.Calculation) (domain.Calculation, error)
	Delete(ctx context.Context, id int64) error
}

type CalculationService struct {
	repo CalculationRepository
}

func New(repo CalculationRepository) *CalculationService {
	return &CalculationService{repo: repo}
}

func (s *CalculationService) Create(ctx context.Context, input domain.Input) (domain.Calculation, error) {
	calc, err := newCalculation(input)
	if err != nil {
		return domain.Calculation{}, err
	}

	return s.repo.Create(ctx, calc)
}

func (s *CalculationService) Get(ctx context.Context, id int64) (domain.Calculation, error) {
	return s.repo.Get(ctx, id)
}

func (s *CalculationService) Update(ctx context.Context, id int64, input domain.Input) (domain.Calculation, error) {
	calc, err := newCalculation(input)
	if err != nil {
		return domain.Calculation{}, err
	}

	return s.repo.Update(ctx, id, calc)
}

func (s *CalculationService) Delete(ctx context.Context, id int64) error {
	return s.repo.Delete(ctx, id)
}

func validate(input domain.Input) error {
	if !input.Operation.Valid() {
		return domain.ErrInvalidOperation
	}
	if math.IsNaN(input.A) || math.IsInf(input.A, 0) || math.IsNaN(input.B) || math.IsInf(input.B, 0) {
		return domain.ErrInvalidNumber
	}
	if input.Operation == domain.OperationDivide && input.B == 0 {
		return domain.ErrDivisionByZero
	}

	return nil
}

func newCalculation(input domain.Input) (domain.Calculation, error) {
	if err := validate(input); err != nil {
		return domain.Calculation{}, err
	}

	var result float64
	switch input.Operation {
	case domain.OperationAdd:
		result = input.A + input.B
	case domain.OperationSubtract:
		result = input.A - input.B
	case domain.OperationMultiply:
		result = input.A * input.B
	case domain.OperationDivide:
		result = input.A / input.B
	default:
		return domain.Calculation{}, domain.ErrInvalidOperation
	}

	if math.IsNaN(result) || math.IsInf(result, 0) {
		return domain.Calculation{}, domain.ErrInvalidNumber
	}

	calc := domain.Calculation{
		A:         input.A,
		B:         input.B,
		Operation: input.Operation,
		Result:    result,
	}

	return calc, nil
}
