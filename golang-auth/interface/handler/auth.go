package handler

import (
	"fmt"
	"math/rand"
	"net/http"

	"encoding/base64"

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
	fmt.Println("Authorize")
	url, err := h.authUseCase.GetAuthorizationURL()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Redirect(301, url)
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

	c.JSON(http.StatusOK, gin.H{
		"redirect_url": redirectURI,
	})
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
	fmt.Println("tokens: ", tokens)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// ランダムなセッションIDを生成
	sessionID := generateRandomSessionID()

	// セッションIDとアクセストークンを紐付けて保存
	h.authUseCase.StoreSession(sessionID, tokens.AccessToken)

	// セッションIDをCookieに設定
	c.SetCookie(
		"poc-authlete",       // 名前
		sessionID,            // 値
		3600,                 // 有効期限（秒）
		"/",                  // パス
		"poc-authlete.local", // ドメイン
		false,                // Secure
		true,                 // HttpOnly
	)

	c.Redirect(301, "http://poc-authlete.local/dashboard")
}

// ランダムなセッションIDを生成
func generateRandomSessionID() string {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return ""
	}
	return base64.URLEncoding.EncodeToString(b)
}

func (h *AuthHandler) GetSession(c *gin.Context) {
	sessionID, err := c.Cookie("poc-authlete")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Session not found"})
		return
	}

	accessToken, err := h.authUseCase.GetAccessToken(sessionID)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid session"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"access_token": accessToken,
	})
}

func (h *AuthHandler) GetUserInfo(c *gin.Context) {
	sessionID, err := c.Cookie("poc-authlete")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Session not found"})
		return
	}

	accessToken, err := h.authUseCase.GetAccessToken(sessionID)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid session"})
		return
	}

	userInfo, err := h.authUseCase.GetUserInfo(accessToken)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, userInfo)
}

func (h *AuthHandler) Logout(c *gin.Context) {
	sessionID, err := c.Cookie("poc-authlete")
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"message": "Logged out"})
		return
	}

	if err := h.authUseCase.DeleteSession(sessionID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Cookieを削除
	c.SetCookie(
		"poc-authlete",
		"",
		-1,
		"/",
		"poc-authlete.local",
		false,
		true,
	)

	c.JSON(http.StatusOK, gin.H{"message": "Logged out"})
}
