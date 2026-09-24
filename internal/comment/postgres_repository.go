package comment

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresRepository(pool *pgxpool.Pool) *PostgresRepository {
	return &PostgresRepository{pool: pool}
}

func (r *PostgresRepository) Create(ctx context.Context, body, authorID, targetType, targetID string) (*Comment, error) {
	var c Comment
	err := r.pool.QueryRow(ctx,
		`INSERT INTO comments (body, author_id, target_type, target_id)
		 VALUES ($1, $2, $3, $4)
		 RETURNING id::text, body, author_id::text, target_type, target_id::text, created_at`,
		body, authorID, targetType, targetID,
	).Scan(&c.ID, &c.Body, &c.AuthorID, &c.TargetType, &c.TargetID, &c.CreatedAt)

	if err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *PostgresRepository) ListByTarget(ctx context.Context, targetType, targetID string) ([]*Comment, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id::text, body, author_id::text, target_type, target_id::text, created_at
		 FROM comments WHERE target_type = $1 AND target_id = $2`,
		targetType, targetID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []*Comment
	for rows.Next() {
		var c Comment
		if err := rows.Scan(&c.ID, &c.Body, &c.AuthorID, &c.TargetType, &c.TargetID, &c.CreatedAt); err != nil {
			return nil, err
		}
		result = append(result, &c)
	}
	return result, rows.Err()
}

func (r *PostgresRepository) GetByID(ctx context.Context, id string) (*Comment, error) {
	var c Comment
	err := r.pool.QueryRow(ctx,
		`SELECT id::text, body, author_id::text, target_type, target_id::text, created_at
		 FROM comments WHERE id = $1`,
		id,
	).Scan(&c.ID, &c.Body, &c.AuthorID, &c.TargetType, &c.TargetID, &c.CreatedAt)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &c, nil
}

func (r *PostgresRepository) Delete(ctx context.Context, id string) error {
	cmdTag, err := r.pool.Exec(ctx, `DELETE FROM comments WHERE id = $1`, id)
	if err != nil {
		return err
	}
	if cmdTag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}