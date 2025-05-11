package repository

import "github.com/yamakenji24/golang-auth/domain/entity"

type PasskeyRepository interface {
	SaveCredential(credential *entity.Credential) error
	GetCredential(id string) (*entity.Credential, error)
	GetCredentialsByUsername(username string) ([]*entity.Credential, error)
}
