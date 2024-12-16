package domain

import "time"

type TaskStatus string

const (
	StatusPending    TaskStatus = "Pending"
	StatusInProgress TaskStatus = "In Progress"
	StatusCompleted  TaskStatus = "Completed"
)

type Task struct {
	ID          string     `json:"id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	AssigneeID  *string    `json:"assignee_id"`
	Status      TaskStatus `json:"status"`
	CreatedAt   time.Time  `json:"created_at"`
	CreatedBy   string     `json:"created_by"`
	UpdatedAt   time.Time  `json:"updated_at"`
	UpdatedBy   string     `json:"updated_by"`
	DueDate     time.Time  `json:"due_date"`
}

type CreateTaskRequest struct {
	Title       string    `json:"title"`
	Description string    `json:"description"`
	DueDate     time.Time `json:"due_date"`
}

type TaskSummary struct {
	EmployeeID     string `json:"employee_id"`
	TotalTasks     int    `json:"total_tasks"`
	CompletedTasks int    `json:"completed_tasks"`
}
