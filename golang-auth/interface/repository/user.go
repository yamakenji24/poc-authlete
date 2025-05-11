package repository

import "github.com/yamakenji24/golang-auth/domain/entity"

type UserRepository interface {
	FindByID(id string) (*entity.User, error)
	FindByUsername(username string) (*entity.User, error)
	Save(user *entity.User) error
	Delete(id string) error
}
