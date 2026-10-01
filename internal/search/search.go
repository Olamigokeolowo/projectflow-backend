package search

type Results struct {
	Decisions []DecisionResult `json:"decisions"`
	Tasks     []TaskResult     `json:"tasks"`
	Comments  []CommentResult  `json:"comments"`
}

type DecisionResult struct {
	ID    string `json:"id"`
	Title string `json:"title"`
}

type TaskResult struct {
	ID         string `json:"id"`
	Title      string `json:"title"`
	DecisionID string `json:"decision_id"`
}

type CommentResult struct {
	ID         string `json:"id"`
	Body       string `json:"body"`
	TargetType string `json:"target_type"`
	TargetID   string `json:"target_id"`
}