package entity

type AuthData struct {
	CodeVerifier string
	Ticket       string
}

type AuthRequest struct {
	State    string `json:"state"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthResponse struct {
	Ticket          string `json:"ticket"`
	ResponseContent string `json:"responseContent"`
}

type TokenResponse struct {
	AccessToken     string `json:"accessToken"`
	RefreshToken    string `json:"refreshToken"`
	IdToken         string `json:"idToken"`
	ResponseContent string `json:"responseContent"`
}

type Tokens struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	IDToken      string `json:"id_token"`
}

type UserInfo struct {
	Sub       string `json:"sub"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	Picture   string `json:"picture"`
	UpdatedAt int64  `json:"updated_at"`
}
