package main

import (
	"log"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/yamakenji24/golang-auth/domain/usecase"
	"github.com/yamakenji24/golang-auth/infrastructure/external/authlete"
	"github.com/yamakenji24/golang-auth/infrastructure/persistence/memory"
	"github.com/yamakenji24/golang-auth/interface/handler"
	"github.com/yamakenji24/golang-auth/pkg/config"
)

func main() {
	r := gin.Default()

	// CORSの設定
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://poc-authlete.local"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * 60 * 60,
	}))

	// 依存性の注入
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	authleteClient := authlete.NewClient(cfg)
	authRepo := memory.NewAuthRepository()
	authUseCase := usecase.NewAuthUseCase(authRepo, authleteClient, cfg, authleteClient)
	authHandler := handler.NewAuthHandler(authUseCase)

	// ルーティング
	api := r.Group("/api")
	{
		auth := api.Group("/auth")
		{
			auth.GET("/authorize", authHandler.Authorize)
			auth.POST("/login", authHandler.Login)
			auth.GET("/callback", authHandler.Callback)
			auth.GET("/session", authHandler.GetSession)
			auth.GET("/userinfo", authHandler.GetUserInfo)
			auth.POST("/logout", authHandler.Logout)
		}
	}

	r.Run(":3000")
}
