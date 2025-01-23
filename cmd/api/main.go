package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/gorm"
	"log"
	"os"
	"sunrise_project/internal/dao"
	"sunrise_project/internal/handler"
	"sunrise_project/internal/platform"
	"sunrise_project/internal/repository"
	"sunrise_project/internal/service"

	"time"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: Error loading .env file, using environment variables instead")
	}

	var db *gorm.DB
	for i := 0; i < 5; i++ {
		db, err = platform.NewPostgresDB()
		if err == nil {
			break
		}
		log.Printf("Failed to connect to database: %v, retrying in 5 seconds...", err)
		time.Sleep(5 * time.Second)
	}
	if err != nil {
		log.Fatalf("Failed to Migrate Database: %s", err)
	}

	err = db.AutoMigrate(&dao.Location{})
	if err != nil {
		log.Fatalf("Failed to Migrate Database: %s", err)
	}

	locationRepo := repository.NewLocationRepository(db)
	locationService := service.NewLocationService(locationRepo)
	locationHandler := handler.NewLocationHandler(locationService)
	secretHandler := handler.NewSecretHandler()

	r := gin.Default()

	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost"}
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	r.Use(cors.New(config))

	r.GET("/location/:ip", locationHandler.GetLocationByCustomIP)
	r.GET("/locations", locationHandler.GetAllLocations)
	r.GET("/", secretHandler.GetSecretValue)
	r.POST("/location", locationHandler.CreateLocation)
	r.PUT("/location/:ip", locationHandler.UpdateLocation)
	r.DELETE("/location/:ip", locationHandler.DeleteLocation)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Starting server on port - %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %s", err)
	}
}
