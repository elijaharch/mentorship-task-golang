package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/elijaharch/mentorship-task-golang/internal/calculation"
	"github.com/jackc/pgx/v5"
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

func (r *Repository) Get(ctx context.Context, id int64) (calculation.Calculation, error) {
	const query = `
		SELECT id, a, b, operation, result, created_at
		FROM numbers
		WHERE id=$1`

	var calc calculation.Calculation
	err := r.pool.QueryRow(ctx,
		query,
		id,
	).Scan(
		&calc.ID,
		&calc.A,
		&calc.B,
		&calc.Operation,
		&calc.Result,
		&calc.CreatedAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return calculation.Calculation{}, calculation.ErrNotFound
	}
	if err != nil {
		return calculation.Calculation{}, fmt.Errorf("get calculation: %w", err)
	}

	return calc, nil
}

func (r *Repository) Update(ctx context.Context, id int64, calc calculation.Calculation) (calculation.Calculation, error) {
	const query = `
		UPDATE numbers
		SET a=$1, b=$2, operation=$3, result=$4
		WHERE id=$5
		RETURNING id, a, b, operation, result, created_at`

	err := r.pool.QueryRow(ctx,
		query,
		calc.A,
		calc.B,
		calc.Operation,
		calc.Result,
		id,
	).Scan(
		&calc.ID,
		&calc.A,
		&calc.B,
		&calc.Operation,
		&calc.Result,
		&calc.CreatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return calculation.Calculation{}, calculation.ErrNotFound
	}
	if err != nil {
		return calculation.Calculation{}, fmt.Errorf("update calculation: %w", err)
	}

	return calc, nil
}
