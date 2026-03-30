package dto

type CreateTaskRequest struct {
	Title       string `json:"title"`
	Schedule    string `json:"schedule"`
	Priority    string `json:"priority"`
	Description string `json:"description"`
}
