package mock

import (
	"github.com/yamakenji24/golang-auth/domain/entity"
)

type MockAuthleteClient struct {
	AuthResponse  *entity.AuthResponse
	TokenResponse *entity.TokenResponse
	UserInfo      *entity.UserInfo
	Error         error
}

func NewMockAuthleteClient() *MockAuthleteClient {
	return &MockAuthleteClient{}
}

func (m *MockAuthleteClient) RequestAuthorization(params map[string]string) (*entity.AuthResponse, error) {
	if m.Error != nil {
		return nil, m.Error
	}
	return m.AuthResponse, nil
}

func (m *MockAuthleteClient) IssueAuthorization(ticket string) (*entity.AuthResponse, error) {
	if m.Error != nil {
		return nil, m.Error
	}
	return m.AuthResponse, nil
}

func (m *MockAuthleteClient) ExchangeToken(params map[string]string) (*entity.TokenResponse, error) {
	if m.Error != nil {
		return nil, m.Error
	}
	return m.TokenResponse, nil
}

func (m *MockAuthleteClient) GetUserInfo(accessToken string) (entity.UserInfo, error) {
	if m.Error != nil {
		return entity.UserInfo{}, m.Error
	}
	return *m.UserInfo, nil
}
