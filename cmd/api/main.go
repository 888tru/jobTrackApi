package main

import (
	"log"
	"net/http"

	"github.com/888tru/jobtrack-api/internal/config"
	"github.com/888tru/jobtrack-api/internal/handler"
	"github.com/888tru/jobtrack-api/internal/middleware"
	"github.com/888tru/jobtrack-api/internal/model"
	"github.com/888tru/jobtrack-api/internal/repository"
	"github.com/888tru/jobtrack-api/internal/service"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	cfg := config.Load()

	db, err := gorm.Open(postgres.Open(cfg.DatabaseURL), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect to database:", err)
	}

	if err := db.AutoMigrate(&model.User{}, &model.JobApplication{}, &model.Interview{}); err != nil {
		log.Fatal("failed to migrate:", err)
	}

	userRepo := repository.NewUserRepository(db)
	appRepo := repository.NewJobApplicationRepository(db)
	interviewRepo := repository.NewInterviewRepository(db)

	authSvc := service.NewAuthService(userRepo, cfg.JWTSecret)
	statsSvc := service.NewStatsService(appRepo, interviewRepo)

	authH := handler.NewAuthHandler(authSvc)
	appH := handler.NewJobApplicationHandler(appRepo)
	interviewH := handler.NewInterviewHandler(interviewRepo, appRepo)
	statsH := handler.NewStatsHandler(statsSvc)

	router := gin.Default()

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	v1 := router.Group("/api/v1")

	auth := v1.Group("/auth")
	{
		auth.POST("/register", authH.Register)
		auth.POST("/login", authH.Login)
	}

	protected := v1.Group("")
	protected.Use(middleware.Auth(cfg.JWTSecret))
	{
		protected.GET("/applications", appH.List)
		protected.POST("/applications", appH.Create)
		protected.GET("/applications/:id", appH.Get)
		protected.PUT("/applications/:id", appH.Update)
		protected.DELETE("/applications/:id", appH.Delete)

		protected.GET("/applications/:id/interviews", interviewH.List)
		protected.POST("/applications/:id/interviews", interviewH.Create)
		protected.PUT("/applications/:id/interviews/:interview_id", interviewH.Update)
		protected.DELETE("/applications/:id/interviews/:interview_id", interviewH.Delete)

		protected.GET("/stats", statsH.GetStats)
	}

	log.Printf("server running on :%s", cfg.Port)
	if err := router.Run(":" + cfg.Port); err != nil {
		log.Fatal(err)
	}
}
