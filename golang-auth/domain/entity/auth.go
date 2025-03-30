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
