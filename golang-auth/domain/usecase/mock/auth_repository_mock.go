package mock

import (
	"github.com/yamakenji24/golang-auth/domain/entity"
)

type MockAuthRepository struct {
	AuthDataMap map[string]entity.AuthData
}

func NewMockAuthRepository() *MockAuthRepository {
	return &MockAuthRepository{
		AuthDataMap: make(map[string]entity.AuthData),
	}
}

func (m *MockAuthRepository) StoreAuthData(state string, authData entity.AuthData) error {
	m.AuthDataMap[state] = authData
	return nil
}

func (m *MockAuthRepository) GetAuthData(state string) (entity.AuthData, bool) {
	authData, ok := m.AuthDataMap[state]
	return authData, ok
}
