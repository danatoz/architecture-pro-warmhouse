package main

import (
	"auth/db"
	_ "auth/docs"
	"auth/handlers"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func main() {
	// Set up database connection
	dbURL := getEnv("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/smarthome_auth")
	database, err := db.New(dbURL)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	defer database.Close()

	log.Println("Connected to database successfully")

	router := gin.Default()

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	})
	// API routes
	apiRoutes := router.Group("/api/v1")

	// Register user routes
	userHandler := handlers.NewUserHandler(database)
	userHandler.RegisterRoutes(apiRoutes)

	// Register home routes
	homeHandler := handlers.NewHomeHandler(database)
	homeHandler.RegisterRoutes(apiRoutes)

	router.Run(":5001") // Запуск на порту 8080
}

// getEnv gets an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
