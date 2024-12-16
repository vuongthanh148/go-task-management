package port

import (
	"context"
	"kn-assignment/internal/core/domain"
)

type TaskRepository interface {
	CreateTask(ctx context.Context, task domain.CreateTaskRequest, userId string) error
	GetTasksByAssignee(ctx context.Context, assigneeID string) ([]domain.Task, error)
	UpdateTaskStatus(ctx context.Context, taskID string, status domain.TaskStatus, userId string) error
	GetAllTasks(ctx context.Context, filter map[string]string, sort, order string) ([]domain.Task, error)
	GetTaskSummary(ctx context.Context) ([]domain.TaskSummary, error)
	AssignTask(ctx context.Context, taskID, assigneeID string) error // New method for assigning tasks
	GetTaskByID(ctx context.Context, taskID string) (domain.Task, error)
	UpdateTask(ctx context.Context, taskID string, name, description *string) error
	DeleteTask(ctx context.Context, taskID string) error
}

type AuthRepository interface {
	CreateUser(ctx context.Context, user domain.CreateUserRequest) error
	GetUserByUsername(ctx context.Context, username string) (*domain.User, error)
	UpdateUser(ctx context.Context, user domain.User) error               // Added method
	GetUserByID(ctx context.Context, userID string) (*domain.User, error) // New method for getting user by ID
}

type UserRepository interface {
	CreateUser(ctx context.Context, user domain.User) error
	GetUserByUsername(ctx context.Context, username string) (domain.User, error)
	GetUserByID(ctx context.Context, userID string) (domain.User, error)
	UpdateUser(ctx context.Context, user domain.User) error
}
