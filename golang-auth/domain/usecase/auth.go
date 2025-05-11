package usecase

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"

	"github.com/yamakenji24/golang-auth/domain/entity"
	"github.com/yamakenji24/golang-auth/interface/repository"
	"github.com/yamakenji24/golang-auth/pkg/config"
)

type AuthUseCase interface {
	GetAuthorizationURL() (string, error)
	Login(req entity.AuthRequest) (string, error)
	GetAuthData(state string) (entity.AuthData, bool)
	ExchangeCodeForTokens(code, codeVerifier string) (entity.Tokens, error)
	StoreSession(sessionID, accessToken string) error
	GetAccessToken(sessionID string) (string, error)
	GetUserInfo(accessToken string) (entity.UserInfo, error)
	DeleteSession(sessionID string) error
}

type authUseCase struct {
	authRepo       repository.AuthRepository
	authleteRepo   repository.AuthleteClient
	config         *config.Config
	authleteClient repository.AuthleteClient
	authDataMap    map[string]entity.AuthData
	sessionMap     map[string]string
}

func NewAuthUseCase(authRepo repository.AuthRepository, authleteRepo repository.AuthleteClient, cfg *config.Config, authleteClient repository.AuthleteClient) AuthUseCase {
	return &authUseCase{
		authRepo:       authRepo,
		authleteRepo:   authleteRepo,
		config:         cfg,
		authleteClient: authleteClient,
		authDataMap:    make(map[string]entity.AuthData),
		sessionMap:     make(map[string]string),
	}
}

func (u *authUseCase) generateState() string {
	b := make([]byte, 16)
	rand.Read(b)
	return fmt.Sprintf("%x", b)
}

func (u *authUseCase) generateCodeVerifier() string {
	b := make([]byte, 32)
	rand.Read(b)
	return base64.RawURLEncoding.EncodeToString(b)
}

func (u *authUseCase) generateCodeChallenge(codeVerifier string) string {
	h := sha256.New()
	h.Write([]byte(codeVerifier))
	return base64.RawURLEncoding.EncodeToString(h.Sum(nil))
}

func (u *authUseCase) GetAuthorizationURL() (string, error) {
	codeVerifier := u.generateCodeVerifier()
	codeChallenge := u.generateCodeChallenge(codeVerifier)
	state := u.generateState()

	params := map[string]string{
		"response_type":         "code",
		"client_id":             u.config.AuthleteClientID,
		"redirect_uri":          u.config.AuthleteRedirectURI,
		"scope":                 "openid",
		"code_challenge":        codeChallenge,
		"code_challenge_method": "S256",
	}

	resp, err := u.authleteRepo.RequestAuthorization(params)
	if err != nil {
		return "", err
	}

	authData := entity.AuthData{
		CodeVerifier: codeVerifier,
		Ticket:       resp.Ticket,
	}

	if err := u.authRepo.StoreAuthData(state, authData); err != nil {
		return "", err
	}

	return fmt.Sprintf("https://poc-authlete.local/auth/login?state=%s", state), nil
}

func (u *authUseCase) Login(req entity.AuthRequest) (string, error) {
	authData, ok := u.authRepo.GetAuthData(req.State)
	if !ok {
		return "", fmt.Errorf("AuthData not found")
	}

	resp, err := u.authleteRepo.IssueAuthorization(authData.Ticket)
	if err != nil {
		return "", err
	}

	fmt.Println(resp.ResponseContent)

	return resp.ResponseContent + "&state=" + req.State, nil
}

func (u *authUseCase) ExchangeCodeForTokens(code, codeVerifier string) (entity.Tokens, error) {
	params := map[string]string{
		"grant_type":    "authorization_code",
		"code":          code,
		"redirect_uri":  u.config.AuthleteRedirectURI,
		"code_verifier": codeVerifier,
	}

	tokenResponse, err := u.authleteRepo.ExchangeToken(params)
	if err != nil {
		return entity.Tokens{}, err
	}

	return entity.Tokens{
		AccessToken:  tokenResponse.AccessToken,
		RefreshToken: tokenResponse.RefreshToken,
		IDToken:      tokenResponse.IdToken,
	}, nil
}

func (u *authUseCase) GetAuthData(state string) (entity.AuthData, bool) {
	return u.authRepo.GetAuthData(state)
}

// StoreSession セッションIDとアクセストークンを紐付けて保存
func (u *authUseCase) StoreSession(sessionID, accessToken string) error {
	u.sessionMap[sessionID] = accessToken
	return nil
}

// GetAccessToken セッションIDからアクセストークンを取得
func (u *authUseCase) GetAccessToken(sessionID string) (string, error) {
	accessToken, ok := u.sessionMap[sessionID]
	if !ok {
		return "", fmt.Errorf("session not found")
	}
	return accessToken, nil
}

// GetUserInfo アクセストークンからユーザー情報を取得
func (u *authUseCase) GetUserInfo(accessToken string) (entity.UserInfo, error) {
	return u.authleteClient.GetUserInfo(accessToken)
}

// DeleteSession セッションを削除
func (u *authUseCase) DeleteSession(sessionID string) error {
	delete(u.sessionMap, sessionID)
	return nil
}
