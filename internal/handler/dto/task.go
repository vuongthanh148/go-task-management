package dto

import (
	"kn-assignment/internal/core/domain"
	"time"
)

type AssignTaskRequest struct {
	AssigneeID string `json:"assignee_id"`
}

type UpdateTaskStatusRequest struct {
	Status domain.TaskStatus `json:"status"`
}
type CreateTaskRequest struct {
	Title       string    `json:"title" example:"New Task"`
	Description string    `json:"description" example:"This is a new task"`
	DueDate     time.Time `json:"due_date" example:"2024-12-31T23:59:59Z"`
}

func (s *CreateTaskRequest) ToDomain() domain.CreateTaskRequest {
	return domain.CreateTaskRequest{
		Title:       s.Title,
		Description: s.Description,
		DueDate:     s.DueDate,
	}
}

type UpdateTaskRequest struct {
	Name        *string `json:"name,omitempty"`
	Description *string `json:"description,omitempty"`
}
