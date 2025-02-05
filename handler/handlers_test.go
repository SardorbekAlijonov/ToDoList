package handler_test

import (
	"ToDoList_App/handler"
	"ToDoList_App/models"
	"ToDoList_App/router"
	"bytes"
	"encoding/json"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func SetupTestDB() *gorm.DB {
	if err := godotenv.Load("../.env"); err != nil {
		log.Println("Warning: .env file not found")
	}

	dsn := "host=" + os.Getenv("DB_HOST") +
		" user=" + os.Getenv("DB_USER") +
		" password=" + os.Getenv("DB_PASSWORD") +
		" dbname=test_db" +
		" port=5432" +
		" sslmode=disable"

	log.Println("Connecting to database with DSN:", dsn)

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
	r := gin.Default()
	h := &handler.Handler{DB: db}
	router.SetupRoutes(r, h)
	return r
}
func TestInvalidUserCreation(t *testing.T) {
	db := SetupTestDB()
	router := SetupRouter(db)

	invalidUser := `{"name":1}`
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

	updateData := `{"id":99,"name":"Ghost User","email":"ghost@example.com"}`
	req, _ := http.NewRequest("PUT", "/users/update", bytes.NewBuffer([]byte(updateData)))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	expectedStatus := http.StatusNotFound // 404 for a non-existent user
	if w.Code != expectedStatus {
		t.Errorf("Expected status %d, got %d", expectedStatus, w.Code)
	}

	expectedMessage := `{"error":"User not found"}`
	if w.Body.String() != expectedMessage {
		t.Errorf("Expected body %s, got %s", expectedMessage, w.Body.String())
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
