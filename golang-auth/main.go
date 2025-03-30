package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/yamakenji24/golang-auth/domain/usecase"
	"github.com/yamakenji24/golang-auth/infrastructure/external/authlete"
	"github.com/yamakenji24/golang-auth/infrastructure/persistence/memory"
	"github.com/yamakenji24/golang-auth/interface/handler"
	"github.com/yamakenji24/golang-auth/pkg/config"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	authRepo := memory.NewAuthRepository()
	authleteClient := authlete.NewClient(cfg)
	authUseCase := usecase.NewAuthUseCase(authRepo, authleteClient, cfg)
	authHandler := handler.NewAuthHandler(authUseCase)

	r := gin.Default()

	// 認可エンドポイント
	r.GET("/auth/authorize", authHandler.Authorize)

	// ログインエンドポイント
	r.POST("/auth/login", authHandler.Login)

	// コールバックエンドポイント
	r.GET("/auth/callback", authHandler.Callback)

	r.Run(":3000")
}
