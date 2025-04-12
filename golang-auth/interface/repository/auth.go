package repository

import "github.com/yamakenji24/golang-auth/domain/entity"

type AuthRepository interface {
	StoreAuthData(state string, data entity.AuthData) error
	GetAuthData(state string) (entity.AuthData, bool)
}
