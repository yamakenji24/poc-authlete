package memory

import (
	"errors"
	"sync"
	"time"

	"github.com/yamakenji24/golang-auth/domain/entity"
	"github.com/yamakenji24/golang-auth/interface/repository"
)

type userRepository struct {
	users map[string]*entity.User
	mu    sync.RWMutex
}

func NewUserRepository() repository.UserRepository {
	return &userRepository{
		users: make(map[string]*entity.User),
	}
}

func (r *userRepository) FindByID(id string) (*entity.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	user, exists := r.users[id]
	if !exists {
		return nil, errors.New("user not found")
	}
	return user, nil
}

func (r *userRepository) FindByUsername(username string) (*entity.User, error) {
	// r.mu.RLock()
	// defer r.mu.RUnlock()

	// for _, user := range r.users {
	// 	if user.Username == username {
	// 		return user, nil
	// 	}
	// }
	// return nil, errors.New("user not found")
	return &entity.User{
		ID:        "1",
		Username:  "John Doe",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}, nil
}

func (r *userRepository) Save(user *entity.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	now := time.Now()
	if user.CreatedAt.IsZero() {
		user.CreatedAt = now
	}
	user.UpdatedAt = now

	r.users[user.ID] = user
	return nil
}

func (r *userRepository) Delete(id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.users[id]; !exists {
		return errors.New("user not found")
	}
	delete(r.users, id)
	return nil
}
