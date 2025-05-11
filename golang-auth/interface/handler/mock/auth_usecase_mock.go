package mock

import (
	"github.com/yamakenji24/golang-auth/domain/entity"
)

type MockAuthUseCase struct {
	GetAuthorizationURLFunc   func() (string, error)
	LoginFunc                 func(req entity.AuthRequest) (string, error)
	GetAuthDataFunc           func(state string) (entity.AuthData, bool)
	ExchangeCodeForTokensFunc func(code, codeVerifier string) (entity.Tokens, error)
	StoreSessionFunc          func(sessionID, accessToken string) error
	GetAccessTokenFunc        func(sessionID string) (string, error)
	GetUserInfoFunc           func(accessToken string) (entity.UserInfo, error)
	DeleteSessionFunc         func(sessionID string) error
}

func NewMockAuthUseCase() *MockAuthUseCase {
	return &MockAuthUseCase{}
}

func (m *MockAuthUseCase) GetAuthorizationURL() (string, error) {
	if m.GetAuthorizationURLFunc != nil {
		return m.GetAuthorizationURLFunc()
	}
	return "", nil
}

func (m *MockAuthUseCase) Login(req entity.AuthRequest) (string, error) {
	if m.LoginFunc != nil {
		return m.LoginFunc(req)
	}
	return "", nil
}

func (m *MockAuthUseCase) GetAuthData(state string) (entity.AuthData, bool) {
	if m.GetAuthDataFunc != nil {
		return m.GetAuthDataFunc(state)
	}
	return entity.AuthData{}, false
}

func (m *MockAuthUseCase) ExchangeCodeForTokens(code, codeVerifier string) (entity.Tokens, error) {
	if m.ExchangeCodeForTokensFunc != nil {
		return m.ExchangeCodeForTokensFunc(code, codeVerifier)
	}
	return entity.Tokens{}, nil
}

func (m *MockAuthUseCase) StoreSession(sessionID, accessToken string) error {
	if m.StoreSessionFunc != nil {
		return m.StoreSessionFunc(sessionID, accessToken)
	}
	return nil
}

func (m *MockAuthUseCase) GetAccessToken(sessionID string) (string, error) {
	if m.GetAccessTokenFunc != nil {
		return m.GetAccessTokenFunc(sessionID)
	}
	return "", nil
}

func (m *MockAuthUseCase) GetUserInfo(accessToken string) (entity.UserInfo, error) {
	if m.GetUserInfoFunc != nil {
		return m.GetUserInfoFunc(accessToken)
	}
	return entity.UserInfo{}, nil
}

func (m *MockAuthUseCase) DeleteSession(sessionID string) error {
	if m.DeleteSessionFunc != nil {
		return m.DeleteSessionFunc(sessionID)
	}
	return nil
}
