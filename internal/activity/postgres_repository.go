package activity

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresRepository(pool *pgxpool.Pool) *PostgresRepository {
	return &PostgresRepository{pool: pool}
}

func (r *PostgresRepository) Create(ctx context.Context, workspaceID, actorID, verb, targetType, targetID, summary string) (*Activity, error) {
	var a Activity
	err := r.pool.QueryRow(ctx,
		`INSERT INTO activity (workspace_id, actor_id, verb, target_type, target_id, summary)
		 VALUES ($1, $2, $3, $4, $5, $6)
		 RETURNING id::text, workspace_id::text, actor_id::text, verb, target_type, target_id::text, summary, created_at`,
		workspaceID, actorID, verb, targetType, targetID, summary,
	).Scan(&a.ID, &a.WorkspaceID, &a.ActorID, &a.Verb, &a.TargetType, &a.TargetID, &a.Summary, &a.CreatedAt)

	if err != nil {
		return nil, err
	}
	return &a, nil
}

func (r *PostgresRepository) ListByWorkspace(ctx context.Context, workspaceID string) ([]*Activity, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id::text, workspace_id::text, actor_id::text, verb, target_type, target_id::text, summary, created_at
		 FROM activity WHERE workspace_id = $1
		 ORDER BY created_at DESC`,
		workspaceID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []*Activity
	for rows.Next() {
		var a Activity
		if err := rows.Scan(&a.ID, &a.WorkspaceID, &a.ActorID, &a.Verb, &a.TargetType, &a.TargetID, &a.Summary, &a.CreatedAt); err != nil {
			return nil, err
		}
		result = append(result, &a)
	}
	return result, rows.Err()
}