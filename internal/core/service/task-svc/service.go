package tasksvc

import "kn-assignment/internal/core/port"

type service struct {
	taskRepo port.TaskRepository
	userRepo port.UserRepository
}

func New(taskRepository port.TaskRepository, userRepo port.UserRepository) port.TaskService {
	return &service{taskRepo: taskRepository, userRepo: userRepo}
}
