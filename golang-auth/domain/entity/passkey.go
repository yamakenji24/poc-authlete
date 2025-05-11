package entity

import "time"

// Passkey はパスキーの情報を保持する構造体です
type Passkey struct {
	ID              string
	UserID          string
	CredentialID    string
	PublicKey       string
	AttestationType string
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

// PasskeyRepository はパスキー情報を管理するリポジトリのインターフェースです
type PasskeyRepository interface {
	FindByID(id string) (*Passkey, error)
	FindByUserID(userID string) ([]*Passkey, error)
	FindByCredentialID(credentialID string) (*Passkey, error)
	Save(passkey *Passkey) error
	Delete(id string) error
}

// WebAuthnRegistrationRequest はパスキー登録開始時のリクエストです
type WebAuthnRegistrationRequest struct {
	Username string `json:"username"`
}

// WebAuthnRegistrationResponse はパスキー登録開始時のレスポンスです
type WebAuthnRegistrationResponse struct {
	PublicKey PublicKeyCredentialCreationOptions `json:"publicKey"`
}

// PublicKeyCredentialCreationOptions はWebAuthnの登録オプションです
type PublicKeyCredentialCreationOptions struct {
	Challenge        string            `json:"challenge"`
	RP               RP                `json:"rp"`
	User             WebAuthnUser      `json:"user"`
	PubKeyCredParams []PubKeyCredParam `json:"pubKeyCredParams"`
	Timeout          int               `json:"timeout"`
	Attestation      string            `json:"attestation"`
}

// RP はRelying Partyの情報です
type RP struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// WebAuthnUser はWebAuthnのユーザー情報です
type WebAuthnUser struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	DisplayName string `json:"displayName"`
}

// PubKeyCredParam は公開鍵クレデンシャルのパラメータです
type PubKeyCredParam struct {
	Type string `json:"type"`
	Alg  int    `json:"alg"`
}

// WebAuthnAuthenticationRequest はパスキー認証開始時のリクエストです
type WebAuthnAuthenticationRequest struct {
	Username string `json:"username"`
}

// WebAuthnAuthenticationResponse はパスキー認証開始時のレスポンスです
type WebAuthnAuthenticationResponse struct {
	PublicKey PublicKeyCredentialRequestOptions `json:"publicKey"`
}

// PublicKeyCredentialRequestOptions はWebAuthnの認証オプションです
type PublicKeyCredentialRequestOptions struct {
	Challenge        string            `json:"challenge"`
	Timeout          int               `json:"timeout"`
	RPID             string            `json:"rpId"`
	AllowCredentials []AllowCredential `json:"allowCredentials"`
}

// AllowCredential は許可されるクレデンシャルの情報です
type AllowCredential struct {
	Type       string   `json:"type"`
	ID         string   `json:"id"`
	Transports []string `json:"transports"`
}

type Credential struct {
	ID              string
	Username        string
	PublicKey       []byte
	UserHandle      []byte
	SignCount       uint32
	Transports      []string
	AttestationType string
	AAGUID          []byte
}
