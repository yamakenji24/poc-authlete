import React, { createContext, useContext, useState, useEffect } from "react";
import { AuthState } from "../types/auth";
import { authApi } from "../api/auth";

interface AuthContextType extends AuthState {
  login: (username: string, password: string) => Promise<void>;
  logout: () => Promise<void>;
  checkAuth: () => Promise<void>;
}

const AuthContext = createContext<AuthContextType | undefined>(undefined);

export const AuthProvider: React.FC<{ children: React.ReactNode }> = ({
  children,
}) => {
  const [state, setState] = useState<AuthState>({
    isAuthenticated: null, // null: 未確認, true: 認証済み, false: 未認証
    user: null,
    accessToken: null,
    isLoading: false,
  });

  const checkAuth = async () => {
    try {
      const session = await authApi.getSession();
      if (session.access_token) {
        const userInfo = await authApi.getUserInfo();
        setState({
          isAuthenticated: true,
          user: userInfo,
          accessToken: session.access_token,
          isLoading: false,
        });
      } else {
        setState({
          isAuthenticated: false,
          user: null,
          accessToken: null,
          isLoading: false,
        });
      }
    } catch (error) {
      console.error("Session check failed:", error);
      setState({
        isAuthenticated: false,
        user: null,
        accessToken: null,
        isLoading: false,
      });
    }
  };

  useEffect(() => {
    checkAuth();
  }, []);

  const login = async (username: string, password: string) => {
    setState((prev) => ({ ...prev, isLoading: true }));
    try {
      const url = new URL(window.location.href);
      const state = url.searchParams.get("state");
      if (!state) {
        throw new Error("state is not found");
      }
      const response = await authApi.login({ state, username, password });
      window.location.replace(response.redirect_url);
    } catch (error) {
      setState((prev) => ({ ...prev, isLoading: false }));
      throw error;
    }
  };

  const logout = async () => {
    setState((prev) => ({ ...prev, isLoading: true }));
    try {
      await authApi.logout();
      setState({
        isAuthenticated: false,
        user: null,
        accessToken: null,
        isLoading: false,
      });
    } catch (error) {
      setState((prev) => ({ ...prev, isLoading: false }));
      throw error;
    }
  };

  return (
    <AuthContext.Provider value={{ ...state, login, logout, checkAuth }}>
      {children}
    </AuthContext.Provider>
  );
};

export const useAuth = () => {
  const context = useContext(AuthContext);
  if (context === undefined) {
    throw new Error("useAuth must be used within an AuthProvider");
  }
  return context;
};
