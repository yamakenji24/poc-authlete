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
	ExchangeCodeForTokens(code, codeVerifier string) (*entity.TokenResponse, error)
	GetAuthData(state string) (entity.AuthData, bool)
}

type authUseCase struct {
	authRepo     repository.AuthRepository
	authleteRepo repository.AuthleteClient
	config       *config.Config
}

func NewAuthUseCase(authRepo repository.AuthRepository, authleteRepo repository.AuthleteClient, cfg *config.Config) AuthUseCase {
	return &authUseCase{
		authRepo:     authRepo,
		authleteRepo: authleteRepo,
		config:       cfg,
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

	return fmt.Sprintf("http://localhost:8081/auth/login?state=%s", state), nil
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

	return resp.ResponseContent + "&state=" + req.State, nil
}

func (u *authUseCase) ExchangeCodeForTokens(code, codeVerifier string) (*entity.TokenResponse, error) {
	params := map[string]string{
		"grant_type":    "authorization_code",
		"code":          code,
		"redirect_uri":  u.config.AuthleteRedirectURI,
		"code_verifier": codeVerifier,
	}

	return u.authleteRepo.ExchangeToken(params)
}

func (u *authUseCase) GetAuthData(state string) (entity.AuthData, bool) {
	return u.authRepo.GetAuthData(state)
}
