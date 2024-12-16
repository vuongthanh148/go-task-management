package taskhdl

import (
	"kn-assignment/internal/core/port"

	"github.com/gin-gonic/gin"
)

type Handler interface {
	CreateTask(c *gin.Context)
	GetTasksByAssignee(c *gin.Context)
	UpdateTaskStatus(c *gin.Context)
	GetAllTasks(c *gin.Context)
	GetTaskSummary(c *gin.Context)
	AssignTask(c *gin.Context)
	UpdateTask(c *gin.Context)
	DeleteTask(c *gin.Context)
}

type handler struct {
	svc port.TaskService
}

func New(svc port.TaskService) Handler {
	return &handler{
		svc: svc,
	}
}
