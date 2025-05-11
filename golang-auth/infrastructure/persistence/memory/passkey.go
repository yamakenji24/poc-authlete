package memory

import (
	"github.com/yamakenji24/golang-auth/domain/entity"
	"github.com/yamakenji24/golang-auth/interface/repository"
)

type PasskeyRepository struct {
	credentials map[string]*entity.Credential
}

func NewPasskeyRepository() repository.PasskeyRepository {
	return &PasskeyRepository{
		credentials: make(map[string]*entity.Credential),
	}
}

func (r *PasskeyRepository) SaveCredential(credential *entity.Credential) error {
	r.credentials[credential.ID] = credential
	return nil
}

func (r *PasskeyRepository) GetCredential(id string) (*entity.Credential, error) {
	if credential, exists := r.credentials[id]; exists {
		return credential, nil
	}
	return nil, nil
}

func (r *PasskeyRepository) GetCredentialsByUsername(username string) ([]*entity.Credential, error) {
	var credentials []*entity.Credential
	for _, credential := range r.credentials {
		if credential.Username == username {
			credentials = append(credentials, credential)
		}
	}
	return credentials, nil
}
