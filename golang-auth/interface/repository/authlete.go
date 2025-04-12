package repository

import "github.com/yamakenji24/golang-auth/domain/entity"

type AuthleteClient interface {
	RequestAuthorization(params map[string]string) (*entity.AuthResponse, error)
	IssueAuthorization(ticket string) (*entity.AuthResponse, error)
	ExchangeToken(params map[string]string) (*entity.TokenResponse, error)
	GetUserInfo(accessToken string) (entity.UserInfo, error)
}
