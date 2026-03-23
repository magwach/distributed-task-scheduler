package dto

type CreateTaskRequest struct {
	Title       string `json:"title"`
	Schedule    string `json:"schedule"`
	Description string `json:"description"`
}
