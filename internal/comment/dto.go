package comment

type CreateCommentRequest struct {
	Body string `json:"body" binding:"required,min=1,max=2000"`
}