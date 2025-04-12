import axios from "axios";
import {
  LoginRequest,
  LoginResponse,
  SessionResponse,
  UserInfo,
} from "../types/auth";

const API_BASE_URL = "http://poc-authlete.local/api";

export const authApi = {
  // 認可リクエストを送信
  requestAuthorization: async (): Promise<void> => {
    window.location.replace(`${API_BASE_URL}/auth/authorize`);
  },

  // ログイン処理
  login: async (data: LoginRequest): Promise<LoginResponse> => {
    const response = await axios.post(`${API_BASE_URL}/auth/login`, data, {
      withCredentials: true,
    });
    console.log("login: ", response.data);
    return response.data;
  },

  getSession: async (): Promise<SessionResponse> => {
    const response = await axios.get(`${API_BASE_URL}/auth/session`, {
      withCredentials: true,
    });
    return response.data;
  },

  getUserInfo: async (): Promise<UserInfo> => {
    const response = await axios.get(`${API_BASE_URL}/auth/userinfo`, {
      withCredentials: true,
    });
    return response.data;
  },

  logout: async (): Promise<void> => {
    await axios.post(
      `${API_BASE_URL}/auth/logout`,
      {},
      {
        withCredentials: true,
      }
    );
  },
};
