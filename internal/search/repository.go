package search

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	pool *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{pool: pool}
}

func (r *Repository) Search(ctx context.Context, workspaceID, query string) (*Results, error) {
	results := &Results{}

	decRows, err := r.pool.Query(ctx,
		`SELECT id::text, title FROM decisions
		 WHERE workspace_id = $1 AND search_vector @@ plainto_tsquery('english', $2)`,
		workspaceID, query,
	)
	if err != nil {
		return nil, err
	}
	defer decRows.Close()
	for decRows.Next() {
		var d DecisionResult
		if err := decRows.Scan(&d.ID, &d.Title); err != nil {
			return nil, err
		}
		results.Decisions = append(results.Decisions, d)
	}

	taskRows, err := r.pool.Query(ctx,
		`SELECT t.id::text, t.title, t.decision_id::text FROM tasks t
		 JOIN decisions d ON d.id = t.decision_id
		 WHERE d.workspace_id = $1 AND t.search_vector @@ plainto_tsquery('english', $2)`,
		workspaceID, query,
	)
	if err != nil {
		return nil, err
	}
	defer taskRows.Close()
	for taskRows.Next() {
		var t TaskResult
		if err := taskRows.Scan(&t.ID, &t.Title, &t.DecisionID); err != nil {
			return nil, err
		}
		results.Tasks = append(results.Tasks, t)
	}

	commentRows, err := r.pool.Query(ctx,
		`SELECT c.id::text, c.body, c.target_type, c.target_id::text FROM comments c
		 WHERE c.search_vector @@ plainto_tsquery('english', $2)
		 AND (
		   (c.target_type = 'decision' AND c.target_id IN (SELECT id FROM decisions WHERE workspace_id = $1))
		   OR
		   (c.target_type = 'task' AND c.target_id IN (SELECT t.id FROM tasks t JOIN decisions d ON d.id = t.decision_id WHERE d.workspace_id = $1))
		 )`,
		workspaceID, query,
	)
	if err != nil {
		return nil, err
	}
	defer commentRows.Close()
	for commentRows.Next() {
		var c CommentResult
		if err := commentRows.Scan(&c.ID, &c.Body, &c.TargetType, &c.TargetID); err != nil {
			return nil, err
		}
		results.Comments = append(results.Comments, c)
	}

	return results, nil
}