package handler

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yamakenji24/golang-auth/domain/usecase"
)

// PasskeyHandler はパスキー認証のHTTPハンドラーを実装します
type PasskeyHandler struct {
	passkeyUseCase *usecase.PasskeyUseCase
}

// NewPasskeyHandler は新しいパスキーハンドラーを作成します
func NewPasskeyHandler(passkeyUseCase *usecase.PasskeyUseCase) *PasskeyHandler {
	return &PasskeyHandler{
		passkeyUseCase: passkeyUseCase,
	}
}

// StartRegistration はパスキー登録を開始するハンドラーです
func (h *PasskeyHandler) StartRegistration(c *gin.Context) {
	var req struct {
		Username string `json:"username"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	fmt.Println("StartRegistration: ", req.Username)

	options, err := h.passkeyUseCase.StartRegistration(req.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, options)
}

// CompleteRegistration はパスキー登録を完了するハンドラーです
func (h *PasskeyHandler) CompleteRegistration(c *gin.Context) {
	var req struct {
		UserID          string `json:"userId"`
		CredentialID    string `json:"credentialId"`
		PublicKey       string `json:"publicKey"`
		AttestationType string `json:"attestationType"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.passkeyUseCase.CompleteRegistration(req.UserID, req.CredentialID, req.PublicKey, req.AttestationType); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success"})
}

// StartAuthentication はパスキー認証を開始するハンドラーです
func (h *PasskeyHandler) StartAuthentication(c *gin.Context) {
	var req struct {
		Username string `json:"username"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	options, err := h.passkeyUseCase.StartAuthentication(req.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, options)
}

// CompleteAuthentication はパスキー認証を完了するハンドラーです
func (h *PasskeyHandler) CompleteAuthentication(c *gin.Context) {
	var req struct {
		CredentialID string `json:"credentialId"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.passkeyUseCase.CompleteAuthentication(req.CredentialID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success"})
}
