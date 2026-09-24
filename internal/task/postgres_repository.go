package task

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

func (r *PostgresRepository) Create(ctx context.Context, title, decisionID, assigneeID, createdBy string) (*Task, error) {
	var t Task
	err := r.pool.QueryRow(ctx,
		`INSERT INTO tasks (title, decision_id, assignee_id, created_by)
		 VALUES ($1, $2, $3, $4)
		 RETURNING id::text, title, status, decision_id::text, assignee_id::text, created_by::text, created_at`,
		title, decisionID, assigneeID, createdBy,
	).Scan(&t.ID, &t.Title, &t.Status, &t.DecisionID, &t.AssigneeID, &t.CreatedBy, &t.CreatedAt)

	if err != nil {
		return nil, err
	}
	return &t, nil
}

func (r *PostgresRepository) ListByDecision(ctx context.Context, decisionID string) ([]*Task, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id::text, title, status, decision_id::text, assignee_id::text, created_by::text, created_at
		 FROM tasks WHERE decision_id = $1`,
		decisionID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []*Task
	for rows.Next() {
		var t Task
		if err := rows.Scan(&t.ID, &t.Title, &t.Status, &t.DecisionID, &t.AssigneeID, &t.CreatedBy, &t.CreatedAt); err != nil {
			return nil, err
		}
		result = append(result, &t)
	}
	return result, rows.Err()
}

func (r *PostgresRepository) GetByID(ctx context.Context, id string) (*Task, error) {
	var t Task
	err := r.pool.QueryRow(ctx,
		`SELECT id::text, title, status, decision_id::text, assignee_id::text, created_by::text, created_at
		 FROM tasks WHERE id = $1`,
		id,
	).Scan(&t.ID, &t.Title, &t.Status, &t.DecisionID, &t.AssigneeID, &t.CreatedBy, &t.CreatedAt)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &t, nil
}

func (r *PostgresRepository) Update(ctx context.Context, t *Task) (*Task, error) {
	cmdTag, err := r.pool.Exec(ctx,
		`UPDATE tasks SET status = $1, assignee_id = $2 WHERE id = $3`,
		t.Status, t.AssigneeID, t.ID,
	)
	if err != nil {
		return nil, err
	}
	if cmdTag.RowsAffected() == 0 {
		return nil, ErrNotFound
	}
	return t, nil
}

func (r *PostgresRepository) Delete(ctx context.Context, id string) error {
	cmdTag, err := r.pool.Exec(ctx, `DELETE FROM tasks WHERE id = $1`, id)
	if err != nil {
		return err
	}
	if cmdTag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}