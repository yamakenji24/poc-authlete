package usecase

import (
	"crypto/rand"
	"encoding/base64"
	"errors"

	"github.com/yamakenji24/golang-auth/domain/entity"
	"github.com/yamakenji24/golang-auth/interface/repository"
)

// PasskeyUseCase はパスキー認証のユースケースを実装します
type PasskeyUseCase struct {
	passkeyRepo repository.PasskeyRepository
	userRepo    repository.UserRepository
}

// NewPasskeyUseCase は新しいパスキーユースケースを作成します
func NewPasskeyUseCase(passkeyRepo repository.PasskeyRepository, userRepo repository.UserRepository) *PasskeyUseCase {
	return &PasskeyUseCase{
		passkeyRepo: passkeyRepo,
		userRepo:    userRepo,
	}
}

// StartRegistration はパスキー登録を開始します
func (u *PasskeyUseCase) StartRegistration(username string) (*entity.WebAuthnRegistrationResponse, error) {
	// ユーザーの存在確認
	user, err := u.userRepo.FindByUsername(username)
	if err != nil {
		return nil, errors.New("user not found")
	}

	// チャレンジの生成
	challenge := make([]byte, 32)
	if _, err := rand.Read(challenge); err != nil {
		return nil, err
	}

	// 登録オプションの生成
	options := &entity.WebAuthnRegistrationResponse{
		PublicKey: entity.PublicKeyCredentialCreationOptions{
			Challenge: base64.RawURLEncoding.EncodeToString(challenge),
			RP: entity.RP{
				ID:   "localhost",
				Name: "Passkey Demo",
			},
			User: entity.WebAuthnUser{
				ID:          user.ID,
				Name:        user.Username,
				DisplayName: user.Username,
			},
			PubKeyCredParams: []entity.PubKeyCredParam{
				{
					Type: "public-key",
					Alg:  -7, // ES256
				},
			},
			Timeout:     60000,
			Attestation: "none",
		},
	}

	return options, nil
}

// CompleteRegistration はパスキー登録を完了します
func (u *PasskeyUseCase) CompleteRegistration(userID string, credentialID string, publicKey string, attestationType string) error {
	credential := &entity.Credential{
		ID:              generateRandomString(32),
		Username:        userID, // ユーザーIDをユーザー名として使用
		PublicKey:       []byte(publicKey),
		UserHandle:      []byte(userID),
		SignCount:       0,
		Transports:      []string{"internal"},
		AttestationType: attestationType,
		AAGUID:          []byte{}, // 必要に応じて設定
	}

	return u.passkeyRepo.SaveCredential(credential)
}

// StartAuthentication はパスキー認証を開始します
func (u *PasskeyUseCase) StartAuthentication(username string) (*entity.WebAuthnAuthenticationResponse, error) {
	// ユーザーの存在確認
	if _, err := u.userRepo.FindByUsername(username); err != nil {
		return nil, errors.New("user not found")
	}

	// ユーザーのパスキーを取得
	credentials, err := u.passkeyRepo.GetCredentialsByUsername(username)
	if err != nil {
		return nil, errors.New("no passkeys found")
	}

	// チャレンジの生成
	challenge := make([]byte, 32)
	if _, err := rand.Read(challenge); err != nil {
		return nil, err
	}

	// 認証オプションの生成
	allowCredentials := make([]entity.AllowCredential, len(credentials))
	for i, credential := range credentials {
		allowCredentials[i] = entity.AllowCredential{
			Type:       "public-key",
			ID:         credential.ID,
			Transports: credential.Transports,
		}
	}

	options := &entity.WebAuthnAuthenticationResponse{
		PublicKey: entity.PublicKeyCredentialRequestOptions{
			Challenge:        base64.RawURLEncoding.EncodeToString(challenge),
			Timeout:          60000,
			RPID:             "localhost",
			AllowCredentials: allowCredentials,
		},
	}

	return options, nil
}

// CompleteAuthentication はパスキー認証を完了します
func (u *PasskeyUseCase) CompleteAuthentication(credentialID string) error {
	// クレデンシャルの取得
	credential, err := u.passkeyRepo.GetCredential(credentialID)
	if err != nil {
		return err
	}
	if credential == nil {
		return errors.New("credential not found")
	}

	// クレデンシャルの検証
	// TODO: 実際の検証ロジックを実装

	return nil
}

// generateRandomString はランダムな文字列を生成します
func generateRandomString(length int) string {
	b := make([]byte, length)
	rand.Read(b)
	return base64.RawURLEncoding.EncodeToString(b)
}
