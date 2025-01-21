package router

import (
	"ToDoList_App/handler"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine, h *handler.Handler) {
	// User routes
	userGroup := router.Group("/users")
	{
		userGroup.POST("/", h.CreateUser)
		userGroup.GET("/", h.GetUsers)
		userGroup.PUT("/", h.UpdateUser)
	}

	// Task routes
	taskGroup := router.Group("/tasks")
	{
		taskGroup.POST("/", h.CreateTask)
		taskGroup.GET("/", h.GetTasks)
		taskGroup.PUT("/", h.UpdateTask)
		taskGroup.DELETE("/:id", h.DeleteTask)
	}

	// Tag routes
	tagGroup := router.Group("/tags")
	{
		tagGroup.POST("/", h.CreateTag)
	}
}
