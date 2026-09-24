package workspace

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

func (r *PostgresRepository) Create(ctx context.Context, name, createdBy string) (*Workspace, error) {
	var w Workspace

	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	err = tx.QueryRow(ctx,
		`INSERT INTO workspaces (name, created_by) VALUES ($1, $2)
		 RETURNING id::text, name, created_by::text, created_at`,
		name, createdBy,
	).Scan(&w.ID, &w.Name, &w.CreatedBy, &w.CreatedAt)
	if err != nil {
		return nil, err
	}

	_, err = tx.Exec(ctx,
		`INSERT INTO workspace_members (workspace_id, user_id) VALUES ($1, $2)`,
		w.ID, createdBy,
	)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	return &w, nil
}

func (r *PostgresRepository) GetByID(ctx context.Context, id string) (*Workspace, error) {
	var w Workspace
	err := r.pool.QueryRow(ctx,
		`SELECT id::text, name, created_by::text, created_at FROM workspaces WHERE id = $1`,
		id,
	).Scan(&w.ID, &w.Name, &w.CreatedBy, &w.CreatedAt)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &w, nil
}

func (r *PostgresRepository) ListForUser(ctx context.Context, userID string) ([]*Workspace, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT w.id::text, w.name, w.created_by::text, w.created_at
		 FROM workspaces w
		 JOIN workspace_members wm ON wm.workspace_id = w.id
		 WHERE wm.user_id = $1`,
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []*Workspace
	for rows.Next() {
		var w Workspace
		if err := rows.Scan(&w.ID, &w.Name, &w.CreatedBy, &w.CreatedAt); err != nil {
			return nil, err
		}
		result = append(result, &w)
	}
	return result, rows.Err()
}

func (r *PostgresRepository) AddMember(ctx context.Context, workspaceID, userID string) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO workspace_members (workspace_id, user_id) VALUES ($1, $2)
		 ON CONFLICT (workspace_id, user_id) DO NOTHING`,
		workspaceID, userID,
	)
	return err
}

func (r *PostgresRepository) IsMember(ctx context.Context, workspaceID, userID string) (bool, error) {
	var exists bool
	err := r.pool.QueryRow(ctx,
		`SELECT EXISTS(
			SELECT 1 FROM workspace_members WHERE workspace_id = $1 AND user_id = $2
		 )`,
		workspaceID, userID,
	).Scan(&exists)
	return exists, err
}

func (r *PostgresRepository) ListMembers(ctx context.Context, workspaceID string) ([]string, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT user_id::text FROM workspace_members WHERE workspace_id = $1`,
		workspaceID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var members []string
	for rows.Next() {
		var userID string
		if err := rows.Scan(&userID); err != nil {
			return nil, err
		}
		members = append(members, userID)
	}
	return members, rows.Err()
}