package port

import (
	"context"
	"kn-assignment/internal/core/domain"
)

type TaskService interface {
	CreateTask(ctx context.Context, task domain.CreateTaskRequest, userId string) error
	AssignTask(ctx context.Context, taskID, assigneeID string) error
	GetTasksByAssignee(ctx context.Context, assigneeID string) ([]domain.Task, error)
	UpdateTaskStatus(ctx context.Context, taskID string, status domain.TaskStatus, assignee string) error
	GetAllTasks(ctx context.Context, userRole, userID string, filter map[string]string, sort, order string) ([]domain.Task, error)
	GetTaskSummary(ctx context.Context) ([]domain.TaskSummary, error)
	VerifyTaskAssignment(ctx context.Context, taskID, userID string) (bool, error)
	UpdateTask(ctx context.Context, taskID string, name, description *string) error
	DeleteTask(ctx context.Context, taskID string) error
}

type AuthService interface {
	RegisterUser(ctx context.Context, user domain.CreateUserRequest) error
	AuthenticateUser(ctx context.Context, username, password string) (domain.LoginResponse, error)
}
