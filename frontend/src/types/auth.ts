export interface LoginRequest {
  state: string;
  username: string;
  password: string;
}

export interface LoginResponse {
  redirect_url: string;
}

export interface SessionResponse {
  access_token: string;
}

export interface UserInfo {
  sub: string;
  name: string;
  email: string;
}

export interface AuthState {
  isAuthenticated: boolean | null;
  user: UserInfo | null;
  accessToken: string | null;
  isLoading: boolean;
}
