package repository

import "github.com/yamakenji24/golang-auth/domain/entity"

type AuthRepository interface {
	StoreAuthData(state string, data entity.AuthData) error
	GetAuthData(state string) (entity.AuthData, bool)
}

type AuthleteClient interface {
	RequestAuthorization(params map[string]string) (*entity.AuthResponse, error)
	IssueAuthorization(ticket string) (*entity.AuthResponse, error)
	ExchangeToken(params map[string]string) (*entity.TokenResponse, error)
}
