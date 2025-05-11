package memory

import (
	"sync"

	"github.com/yamakenji24/golang-auth/domain/entity"
	"github.com/yamakenji24/golang-auth/interface/repository"
)

type authRepository struct {
	data sync.Map
}

func NewAuthRepository() repository.AuthRepository {
	return &authRepository{}
}

func (r *authRepository) StoreAuthData(state string, data entity.AuthData) error {
	r.data.Store(state, data)
	return nil
}

func (r *authRepository) GetAuthData(state string) (entity.AuthData, bool) {
	if data, ok := r.data.Load(state); ok {
		if authData, ok := data.(entity.AuthData); ok {
			return authData, true
		}
	}
	return entity.AuthData{}, false
}
