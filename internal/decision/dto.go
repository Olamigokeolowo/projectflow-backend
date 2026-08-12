package decision

type CreateDecisionRequest struct {
	Title  string `json:"title" binding:"required,min=3,max=200"`
	Status string `json:"status" binding:"required,oneof=draft proposed accepted rejected"`
}