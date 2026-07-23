package repository

import (
	"context"
	"fmt"

	"github.com/elijaharch/mentorship-task-golang/internal/calculation"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	pool *pgxpool.Pool
}

func New(pool *pgxpool.Pool) *Repository {
	return &Repository{
		pool: pool,
	}
}

func (r *Repository) Create(ctx context.Context, calc calculation.Calculation) (calculation.Calculation, error) {
	const query = `
		INSERT INTO numbers (a, b, operation, result)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at`
	err := r.pool.QueryRow(ctx,
		query,
		calc.A,
		calc.B,
		calc.Operation,
		calc.Result,
	).Scan(
		&calc.ID,
		&calc.CreatedAt,
	)
	if err != nil {
		return calculation.Calculation{}, fmt.Errorf("create calculation: %w", err)
	}

	return calc, nil
}
