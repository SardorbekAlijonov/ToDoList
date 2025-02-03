package handler_test

import (
	"ToDoList_App/handler"
	"ToDoList_App/models"
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func SetupTestDB() *gorm.DB {
	dsn := "host=localhost user=postgres password=1506 dbname=test_db port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		log.Fatalf("Failed to connect to test database: %v", err)
	}
	db.AutoMigrate(&models.User{}, &models.Task{}, &models.Tag{})
	return db
}

func SetupRouter(db *gorm.DB) *gin.Engine {
	router := gin.Default()
	h := &handler.Handler{DB: db}
	router.POST("/users/", h.CreateUser)
	router.GET("/users/", h.GetUsers)
	router.PUT("/users/update", h.UpdateUser)
	router.POST("/tasks/", h.CreateTask)
	router.GET("/tasks/", h.GetTasks)
	router.PUT("/tasks/", h.UpdateTask)
	router.DELETE("/tasks/:id", h.DeleteTask)
	router.POST("/tags/", h.CreateTag)
	return router
}

func TestInvalidUserCreation(t *testing.T) {
	db := SetupTestDB()
	router := SetupRouter(db)

	invalidUser := `{"name":""}`
	req, _ := http.NewRequest("POST", "/users/", bytes.NewBuffer([]byte(invalidUser)))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
}

func TestNonExistentUserUpdate(t *testing.T) {
	db := SetupTestDB()
	router := SetupRouter(db)

	updateData := `{"id":999,"name":"Ghost User","email":"ghost@example.com"}`
	req, _ := http.NewRequest("PUT", "/users/update", bytes.NewBuffer([]byte(updateData)))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status 500, got %d", w.Code)
	}
}

func TestFetchEmptyTasks(t *testing.T) {
	db := SetupTestDB()
	router := SetupRouter(db)

	req, _ := http.NewRequest("GET", "/tasks/", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var tasks []models.Task
	json.Unmarshal(w.Body.Bytes(), &tasks)

	if len(tasks) != 0 {
		t.Errorf("Expected empty task list, got %d tasks", len(tasks))
	}
}

func TestInvalidTaskDeletion(t *testing.T) {
	db := SetupTestDB()
	router := SetupRouter(db)

	req, _ := http.NewRequest("DELETE", "/tasks/999", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status 500, got %d", w.Code)
	}
}

func TestCreateAndRetrieveTag(t *testing.T) {
	db := SetupTestDB()
	router := SetupRouter(db)

	tagData := `{"name":"Urgent"}`
	req, _ := http.NewRequest("POST", "/tags/", bytes.NewBuffer([]byte(tagData)))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("Expected status 201, got %d", w.Code)
	}
}
