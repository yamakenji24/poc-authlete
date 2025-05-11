package usecase

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/yamakenji24/golang-auth/domain/entity"
	"github.com/yamakenji24/golang-auth/domain/usecase/mock"
	"github.com/yamakenji24/golang-auth/pkg/config"
)

func TestGetAuthorizationURL(t *testing.T) {
	// テストケースの準備
	mockAuthRepo := mock.NewMockAuthRepository()
	mockAuthleteClient := mock.NewMockAuthleteClient()
	cfg := &config.Config{
		AuthleteClientID:    "test-client-id",
		AuthleteRedirectURI: "http://localhost:3000/callback",
	}

	// モックの設定
	mockAuthleteClient.AuthResponse = &entity.AuthResponse{
		Ticket:          "test-ticket",
		ResponseContent: "test-response",
	}

	// ユースケースの作成
	authUseCase := NewAuthUseCase(mockAuthRepo, mockAuthleteClient, cfg, mockAuthleteClient)

	// テスト実行
	url, err := authUseCase.GetAuthorizationURL()

	// アサーション
	assert.NoError(t, err)
	assert.Contains(t, url, "https://poc-authlete.local/auth/login")
	assert.Contains(t, url, "state=")

	// 認証データが正しく保存されているか確認
	state := url[len(url)-32:] // stateの長さは32文字
	authData, ok := mockAuthRepo.GetAuthData(state)
	assert.True(t, ok)
	assert.Equal(t, "test-ticket", authData.Ticket)
	assert.NotEmpty(t, authData.CodeVerifier)
}

func TestLogin(t *testing.T) {
	// テストケースの準備
	mockAuthRepo := mock.NewMockAuthRepository()
	mockAuthleteClient := mock.NewMockAuthleteClient()
	cfg := &config.Config{}

	// モックの設定
	state := "test-state"
	authData := entity.AuthData{
		CodeVerifier: "test-code-verifier",
		Ticket:       "test-ticket",
	}
	mockAuthRepo.StoreAuthData(state, authData)
	mockAuthleteClient.AuthResponse = &entity.AuthResponse{
		ResponseContent: "test-response",
	}

	// ユースケースの作成
	authUseCase := NewAuthUseCase(mockAuthRepo, mockAuthleteClient, cfg, mockAuthleteClient)

	// テスト実行
	req := entity.AuthRequest{State: state}
	response, err := authUseCase.Login(req)

	// アサーション
	assert.NoError(t, err)
	assert.Contains(t, response, "test-response")
	assert.Contains(t, response, state)
}

func TestExchangeCodeForTokens(t *testing.T) {
	// テストケースの準備
	mockAuthRepo := mock.NewMockAuthRepository()
	mockAuthleteClient := mock.NewMockAuthleteClient()
	cfg := &config.Config{
		AuthleteRedirectURI: "http://localhost:3000/callback",
	}

	// モックの設定
	mockAuthleteClient.TokenResponse = &entity.TokenResponse{
		AccessToken:     "test-access-token",
		RefreshToken:    "test-refresh-token",
		IdToken:         "test-id-token",
		ResponseContent: "test-response-content",
	}

	// ユースケースの作成
	authUseCase := NewAuthUseCase(mockAuthRepo, mockAuthleteClient, cfg, mockAuthleteClient)

	// テスト実行
	tokens, err := authUseCase.ExchangeCodeForTokens("test-code", "test-code-verifier")

	// アサーション
	assert.NoError(t, err)
	assert.Equal(t, "test-access-token", tokens.AccessToken)
	assert.Equal(t, "test-refresh-token", tokens.RefreshToken)
	assert.Equal(t, "test-id-token", tokens.IDToken)
}

func TestGetUserInfo(t *testing.T) {
	// テストケースの準備
	mockAuthRepo := mock.NewMockAuthRepository()
	mockAuthleteClient := mock.NewMockAuthleteClient()
	cfg := &config.Config{}

	// モックの設定
	mockAuthleteClient.UserInfo = &entity.UserInfo{
		Sub:       "test-sub",
		Name:      "Test User",
		Email:     "test@example.com",
		Picture:   "https://example.com/picture.jpg",
		UpdatedAt: 1234567890,
	}

	// ユースケースの作成
	authUseCase := NewAuthUseCase(mockAuthRepo, mockAuthleteClient, cfg, mockAuthleteClient)

	// テスト実行
	userInfo, err := authUseCase.GetUserInfo("test-access-token")

	// アサーション
	assert.NoError(t, err)
	assert.Equal(t, "test-sub", userInfo.Sub)
	assert.Equal(t, "Test User", userInfo.Name)
	assert.Equal(t, "test@example.com", userInfo.Email)
	assert.Equal(t, "https://example.com/picture.jpg", userInfo.Picture)
	assert.Equal(t, int64(1234567890), userInfo.UpdatedAt)
}
