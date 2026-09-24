package decision

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

func (r *PostgresRepository) Create(ctx context.Context, title, status, ownerID, workspaceID string) (*Decision, error) {
	var d Decision
	err := r.pool.QueryRow(ctx,
		`INSERT INTO decisions (title, status, owner_id, workspace_id)
		 VALUES ($1, $2, $3, $4)
		 RETURNING id::text, title, status, owner_id::text, workspace_id::text, created_at`,
		title, status, ownerID, workspaceID,
	).Scan(&d.ID, &d.Title, &d.Status, &d.OwnerID, &d.WorkspaceID, &d.CreatedAt)

	if err != nil {
		return nil, err
	}
	return &d, nil
}

func (r *PostgresRepository) ListByWorkspace(ctx context.Context, workspaceID string) ([]*Decision, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id::text, title, status, owner_id::text, workspace_id::text, created_at
		 FROM decisions WHERE workspace_id = $1`,
		workspaceID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []*Decision
	for rows.Next() {
		var d Decision
		if err := rows.Scan(&d.ID, &d.Title, &d.Status, &d.OwnerID, &d.WorkspaceID, &d.CreatedAt); err != nil {
			return nil, err
		}
		result = append(result, &d)
	}
	return result, rows.Err()
}

func (r *PostgresRepository) GetByID(ctx context.Context, id string) (*Decision, error) {
	var d Decision
	err := r.pool.QueryRow(ctx,
		`SELECT id::text, title, status, owner_id::text, workspace_id::text, created_at
		 FROM decisions WHERE id = $1`,
		id,
	).Scan(&d.ID, &d.Title, &d.Status, &d.OwnerID, &d.WorkspaceID, &d.CreatedAt)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &d, nil
}

func (r *PostgresRepository) Update(ctx context.Context, d *Decision) (*Decision, error) {
	cmdTag, err := r.pool.Exec(ctx,
		`UPDATE decisions SET title = $1, status = $2 WHERE id = $3`,
		d.Title, d.Status, d.ID,
	)
	if err != nil {
		return nil, err
	}
	if cmdTag.RowsAffected() == 0 {
		return nil, ErrNotFound
	}
	return d, nil
}

func (r *PostgresRepository) Delete(ctx context.Context, id string) error {
	cmdTag, err := r.pool.Exec(ctx, `DELETE FROM decisions WHERE id = $1`, id)
	if err != nil {
		return err
	}
	if cmdTag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *PostgresRepository) SlowOperation(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}
	_, err := r.pool.Exec(ctx, `SELECT pg_sleep(5)`)
	return err
}