package router

import (
	"ToDoList_App/handler"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine, h *handler.Handler) {
	// Define routes and associate with handler functions
	router.POST("/create-user", h.CreateUser)
	router.GET("/users", h.GetUsers)
	router.PUT("/user", h.UpdateUser)

	router.POST("/create-task", h.CreateTask)
	router.GET("/tasks", h.GetTasks)
	router.PUT("/task", h.UpdateTask)
	router.DELETE("/delete-task/:id", h.DeleteTask)
}
