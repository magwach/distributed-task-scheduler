package dto

type CreateTaskRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Schedule    string `json:"schedule"`
}
