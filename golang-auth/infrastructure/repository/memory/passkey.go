package memory

import (
	"errors"
	"sync"
	"time"

	"github.com/yamakenji24/golang-auth/domain/entity"
)

// passkeyRepository はメモリ上でパスキー情報を管理するリポジトリの実装です
type passkeyRepository struct {
	passkeys map[string]*entity.Passkey // パスキー情報を保持するマップ
	mu       sync.RWMutex               // 並行アクセス制御用のミューテックス
}

// NewPasskeyRepository は新しいパスキーリポジトリを作成します
func NewPasskeyRepository() *passkeyRepository {
	return &passkeyRepository{
		passkeys: make(map[string]*entity.Passkey),
	}
}

// FindByID は指定されたIDのパスキーを取得します
func (r *passkeyRepository) FindByID(id string) (*entity.Passkey, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	passkey, exists := r.passkeys[id]
	if !exists {
		return nil, errors.New("passkey not found")
	}
	return passkey, nil
}

// FindByUserID は指定されたユーザーIDに紐づくパスキーを取得します
func (r *passkeyRepository) FindByUserID(userID string) ([]*entity.Passkey, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var result []*entity.Passkey
	for _, passkey := range r.passkeys {
		if passkey.UserID == userID {
			result = append(result, passkey)
		}
	}
	return result, nil
}

// FindByCredentialID は指定されたクレデンシャルIDのパスキーを取得します
func (r *passkeyRepository) FindByCredentialID(credentialID string) (*entity.Passkey, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, passkey := range r.passkeys {
		if passkey.CredentialID == credentialID {
			return passkey, nil
		}
	}
	return nil, errors.New("passkey not found")
}

// Save はパスキー情報を保存します
func (r *passkeyRepository) Save(passkey *entity.Passkey) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	now := time.Now()
	if passkey.CreatedAt.IsZero() {
		passkey.CreatedAt = now
	}
	passkey.UpdatedAt = now

	r.passkeys[passkey.ID] = passkey
	return nil
}

// Delete は指定されたIDのパスキーを削除します
func (r *passkeyRepository) Delete(id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.passkeys[id]; !exists {
		return errors.New("passkey not found")
	}
	delete(r.passkeys, id)
	return nil
}
