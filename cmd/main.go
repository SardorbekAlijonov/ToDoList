package main

import (
	_ "ToDoList_App/docs"
	handler "ToDoList_App/handler"
	"ToDoList_App/models"
	"ToDoList_App/router"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
)

var DB *gorm.DB

// Helper for default fallback
func defaultIfEmpty(value, fallback string) string {
	if value == "" {
		return fallback
	}
	return value
}

func initDatabase() *gorm.DB {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found")
	}

	// Build DSN
	dsn := "host=" + os.Getenv("DB_HOST") +
		" user=" + os.Getenv("DB_USER") +
		" password=" + os.Getenv("DB_PASSWORD") +
		" dbname=" + os.Getenv("DB_NAME") +
		" port=" + os.Getenv("DB_PORT") +
		" TimeZone=" + defaultIfEmpty(os.Getenv("DB_TIMEZONE"), "UTC") +
		" sslmode=" + defaultIfEmpty(os.Getenv("DB_SSLMODE"), "disable")

	log.Println("Connecting to database with DSN:", dsn)

	// Open DB connection
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info), // Enable SQL Logging
	})
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}

	// Auto-migrate models
	log.Println("Starting auto-migration")
	if err := db.AutoMigrate(&models.User{}, &models.Task{}, &models.Tag{}); err != nil {
		log.Fatalf("Migration failed: %v", err)
	}
	log.Println("Auto-migration completed successfully")

	log.Println("Connected to PostgreSQL successfully!")
	return db
}

func main() {
	// Initialize the database
	db := initDatabase()
	h := handler.Handler{DB: db}

	// Ensure database connection cleanup on exit
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("Failed to retrieve database connection: %v", err)
	}
	defer func() {
		if err := sqlDB.Close(); err != nil {
			log.Printf("Error while closing the database connection: %v", err)
		} else {
			log.Println("Database connection closed successfully")
		}
	}()

	log.Println("Starting the application...")
	r := gin.New()

	router.SetupRoutes(r, &h)
	// Swagger documentation route
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	//router.SetupRoutes(r, &h)
	// Start the server
	if err := r.Run("localhost:8080"); err != nil {
		log.Fatalf("Failed to start the server: %v", err)
	}
}
