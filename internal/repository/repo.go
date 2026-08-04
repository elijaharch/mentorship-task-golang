package repository

import (
	"context"
	"errors"
	"fmt"

	calculation "github.com/elijaharch/mentorship-task-golang/internal/domain"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	pool *pgxpool.Pool
}

func New(pool *pgxpool.Pool) *Repository {
	return &Repository{pool: pool}
}

func (r *Repository) Create(ctx context.Context, calc calculation.Calculation) (calculation.Calculation, error) {
	const query = `
		INSERT INTO numbers (a, b, operation, result, command_id)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at`
	err := r.pool.QueryRow(ctx,
		query,
		calc.A,
		calc.B,
		calc.Operation,
		calc.Result,
		calc.CommandID,
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
		SELECT id, a, b, operation, result, command_id, created_at
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
		&calc.CommandID,
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
		SET a=$1, b=$2, operation=$3, result=$4, command_id=$5
		WHERE id=$6
		RETURNING id, a, b, operation, result, command_id, created_at`

	err := r.pool.QueryRow(ctx,
		query,
		calc.A,
		calc.B,
		calc.Operation,
		calc.Result,
		calc.CommandID,
		id,
	).Scan(
		&calc.ID,
		&calc.A,
		&calc.B,
		&calc.Operation,
		&calc.Result,
		&calc.CommandID,
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

func (r *Repository) Delete(ctx context.Context, id int64) error {
	const query = `
		DELETE FROM numbers
		WHERE id = $1`

	commandTag, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("delete calculation: %w", err)
	}
	if commandTag.RowsAffected() == 0 {
		return calculation.ErrNotFound
	}

	return nil
}
