package comment

import (
	"context"
	"errors"
)

var ErrInvalidTarget = errors.New("comment target must be 'decision' or 'task'")

// DecisionWorkspaceFinder resolves a decision's workspace directly.
type DecisionWorkspaceFinder interface {
	GetWorkspaceID(ctx context.Context, decisionID string) (string, error)
}

// TaskFinder resolves a task's own DecisionID, so the resolver can chain
// through to the decision, and from there to the workspace.
type TaskFinder interface {
	GetDecisionID(ctx context.Context, taskID string) (string, error)
}

// Resolver implements WorkspaceResolver by branching on target type.
type Resolver struct {
	decisions DecisionWorkspaceFinder
	tasks     TaskFinder
}

func NewResolver(decisions DecisionWorkspaceFinder, tasks TaskFinder) *Resolver {
	return &Resolver{decisions: decisions, tasks: tasks}
}

func (r *Resolver) GetWorkspaceID(ctx context.Context, targetType, targetID string) (string, error) {
	switch targetType {
	case "decision":
		return r.decisions.GetWorkspaceID(ctx, targetID)
	case "task":
		decisionID, err := r.tasks.GetDecisionID(ctx, targetID)
		if err != nil {
			return "", err
		}
		return r.decisions.GetWorkspaceID(ctx, decisionID)
	default:
		return "", ErrInvalidTarget
	}
}