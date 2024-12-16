package router

import (
	"kn-assignment/docs"
	"kn-assignment/internal/core/domain"
	authhdl "kn-assignment/internal/handler/auth-hdl"
	taskhdl "kn-assignment/internal/handler/task-hdl"
	"kn-assignment/internal/middleware"

	"kn-assignment/property"
	"net/http"

	"github.com/gin-gonic/gin"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type HandlerList struct {
	TaskHandler taskhdl.Handler
	AuthHandler authhdl.Handler
}

const serviceBaseURL = "/api/v1"

func InitRouter(app *gin.Engine, h HandlerList) {
	docs.SwaggerInfo.Title = property.Get().Server.ServiceName
	docs.SwaggerInfo.Description = property.Get().Server.ServiceDescription
	docs.SwaggerInfo.Version = property.Get().Server.ApiDocsVersion
	docs.SwaggerInfo.Host = property.Get().Server.Host + ":" + property.Get().Server.Port
	docs.SwaggerInfo.Schemes = []string{property.Get().Server.ApiDocsSchema}
	docs.SwaggerInfo.BasePath = serviceBaseURL

	// Add API doc endpoint to router
	if property.Get().Server.ApiDocs {
		docPath := "/docs"
		app.GET(docPath, func(ctx *gin.Context) { ctx.Redirect(http.StatusTemporaryRedirect, docPath+"/swagger/index.html") })
		app.GET(docPath+"/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}

	v1 := app.Group("/api/v1")

	// common routes
	v1.GET("/tasks", middleware.AuthMiddleware(), h.TaskHandler.GetAllTasks)

	// auth routes
	auth := v1.Group("/auth")
	auth.POST("/register", h.AuthHandler.Register)
	auth.POST("/login", h.AuthHandler.Login)
	auth.POST("/refresh-token", h.AuthHandler.RefreshToken)

	// employee routes
	employee := v1.Group("/")
	employee.Use(middleware.AuthMiddleware())
	employee.GET("/tasks/assignee/:assigneeID", h.TaskHandler.GetTasksByAssignee)
	employee.PATCH("/tasks/:taskID/status", h.TaskHandler.UpdateTaskStatus)

	// employer routes
	employer := v1.Group("/")
	employer.Use(middleware.AuthMiddleware())
	employer.Use(middleware.RoleMiddleware(domain.RoleEmployer))
	employer.POST("/tasks", h.TaskHandler.CreateTask)
	employer.PATCH("/tasks/:taskID/assign", h.TaskHandler.AssignTask)
	employer.GET("/tasks/summary", h.TaskHandler.GetTaskSummary)
	employer.PATCH("/tasks/:taskID", middleware.AuthMiddleware(), h.TaskHandler.UpdateTask)
	employer.DELETE("/tasks/:taskID", middleware.AuthMiddleware(), h.TaskHandler.DeleteTask)
}
