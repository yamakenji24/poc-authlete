package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/yamakenji24/golang-auth/domain/entity"
	"github.com/yamakenji24/golang-auth/interface/handler/mock"
)

func setupTestRouter() (*gin.Engine, *mock.MockAuthUseCase) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	mockUseCase := mock.NewMockAuthUseCase()
	authHandler := NewAuthHandler(mockUseCase)

	api := router.Group("/api")
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

	return router, mockUseCase
}

func TestAuthorize(t *testing.T) {
	router, mockUseCase := setupTestRouter()

	// モックの設定
	expectedURL := "https://poc-authlete.local/auth/login?state=test-state"
	mockUseCase.GetAuthorizationURLFunc = func() (string, error) {
		return expectedURL, nil
	}

	// テストリクエストの作成
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/auth/authorize", nil)
	router.ServeHTTP(w, req)

	// アサーション
	assert.Equal(t, http.StatusFound, w.Code)
	assert.Equal(t, expectedURL, w.Header().Get("Location"))
}

func TestLogin(t *testing.T) {
	router, mockUseCase := setupTestRouter()

	// モックの設定
	expectedRedirectURL := "https://poc-authlete.local/callback?code=test-code&state=test-state"
	mockUseCase.LoginFunc = func(req entity.AuthRequest) (string, error) {
		return expectedRedirectURL, nil
	}

	// テストリクエストの作成
	loginReq := entity.AuthRequest{
		State: "test-state",
	}
	reqBody, _ := json.Marshal(loginReq)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/auth/login", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	// アサーション
	assert.Equal(t, http.StatusOK, w.Code)
	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, expectedRedirectURL, response["redirect_url"])
}

func TestCallback(t *testing.T) {
	router, mockUseCase := setupTestRouter()

	// モックの設定
	mockUseCase.GetAuthDataFunc = func(state string) (entity.AuthData, bool) {
		return entity.AuthData{
			CodeVerifier: "test-code-verifier",
			Ticket:       "test-ticket",
		}, true
	}

	mockUseCase.ExchangeCodeForTokensFunc = func(code, codeVerifier string) (entity.Tokens, error) {
		return entity.Tokens{
			AccessToken:  "test-access-token",
			RefreshToken: "test-refresh-token",
			IDToken:      "test-id-token",
		}, nil
	}

	mockUseCase.StoreSessionFunc = func(sessionID, accessToken string) error {
		return nil
	}

	// テストリクエストの作成
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/auth/callback?state=test-state&code=test-code", nil)
	router.ServeHTTP(w, req)

	// アサーション
	assert.Equal(t, http.StatusFound, w.Code)
	assert.Equal(t, "https://poc-authlete.local/dashboard", w.Header().Get("Location"))
}

func TestGetUserInfo(t *testing.T) {
	router, mockUseCase := setupTestRouter()

	// モックの設定
	mockUseCase.GetAccessTokenFunc = func(sessionID string) (string, error) {
		return "test-access-token", nil
	}

	mockUseCase.GetUserInfoFunc = func(accessToken string) (entity.UserInfo, error) {
		return entity.UserInfo{
			Sub:       "test-sub",
			Name:      "Test User",
			Email:     "test@example.com",
			Picture:   "https://example.com/picture.jpg",
			UpdatedAt: 1234567890,
		}, nil
	}

	// テストリクエストの作成
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/auth/userinfo", nil)
	req.AddCookie(&http.Cookie{
		Name:  "poc-authlete",
		Value: "test-session-id",
	})
	router.ServeHTTP(w, req)

	// アサーション
	assert.Equal(t, http.StatusOK, w.Code)
	var userInfo entity.UserInfo
	err := json.Unmarshal(w.Body.Bytes(), &userInfo)
	assert.NoError(t, err)
	assert.Equal(t, "test-sub", userInfo.Sub)
	assert.Equal(t, "Test User", userInfo.Name)
	assert.Equal(t, "test@example.com", userInfo.Email)
}

func TestLogout(t *testing.T) {
	router, mockUseCase := setupTestRouter()

	// モックの設定
	mockUseCase.DeleteSessionFunc = func(sessionID string) error {
		return nil
	}

	// テストリクエストの作成
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/auth/logout", nil)
	req.AddCookie(&http.Cookie{
		Name:  "poc-authlete",
		Value: "test-session-id",
	})
	router.ServeHTTP(w, req)

	// アサーション
	assert.Equal(t, http.StatusOK, w.Code)
	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "Logged out", response["message"])

	// Cookieが削除されていることを確認
	cookies := w.Result().Cookies()
	var found bool
	for _, cookie := range cookies {
		if cookie.Name == "poc-authlete" && cookie.Value == "" {
			found = true
			break
		}
	}
	assert.True(t, found)
}
