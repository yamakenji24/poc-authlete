package handler

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yamakenji24/golang-auth/domain/entity"
	"github.com/yamakenji24/golang-auth/domain/usecase"
)

type AuthHandler struct {
	authUseCase usecase.AuthUseCase
}

func NewAuthHandler(authUseCase usecase.AuthUseCase) *AuthHandler {
	return &AuthHandler{
		authUseCase: authUseCase,
	}
}

func (h *AuthHandler) Authorize(c *gin.Context) {
	url, err := h.authUseCase.GetAuthorizationURL()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	fmt.Println("url:", url)

	c.Redirect(http.StatusTemporaryRedirect, url)
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req entity.AuthRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	redirectURI, err := h.authUseCase.Login(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	fmt.Println("redirectURI:", redirectURI)

	c.Redirect(http.StatusTemporaryRedirect, redirectURI)
}

func (h *AuthHandler) Callback(c *gin.Context) {
	state := c.Query("state")
	code := c.Query("code")

	if state == "" || code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	authData, ok := h.authUseCase.GetAuthData(state)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "AuthData not found"})
		return
	}

	tokens, err := h.authUseCase.ExchangeCodeForTokens(code, authData.CodeVerifier)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, tokens)
}
