import axios from "axios";

const API_BASE_URL = "https://poc-authlete.local/api";

export interface Passkey {
  id: string;
  userId: string;
  name?: string;
  createdAt: string;
  lastUsedAt?: string;
}

export const passkeyApi = {
  startRegistration: async (username: string) => {
    console.log("startRegistration: ", username);
    const response = await axios.post(
      `${API_BASE_URL}/passkey/register/start`,
      {
        username,
      }
    );
    return response.data;
  },

  completeRegistration: async (
    userId: string,
    credentialId: string,
    publicKey: string,
    attestationType: string
  ) => {
    await axios.post(`${API_BASE_URL}/passkey/register/complete`, {
      userId,
      credentialId,
      publicKey,
      attestationType,
    });
  },

  startAuthentication: async (username: string) => {
    const response = await axios.post(
      `${API_BASE_URL}/passkey/authenticate/start`,
      {
        username,
      }
    );
    return response.data;
  },

  completeAuthentication: async (credentialId: string) => {
    await axios.post(`${API_BASE_URL}/passkey/authenticate/complete`, {
      credentialId,
    });
  },

  // 登録済みパスキーの取得
  getPasskeys: async (userId: string): Promise<Passkey[]> => {
    const response = await axios.get(`${API_BASE_URL}/passkey/list`, {
      params: { userId },
    });
    return response.data;
  },

  // パスキーの削除
  deletePasskey: async (passkeyId: string) => {
    await axios.delete(`${API_BASE_URL}/passkey/${passkeyId}`);
  },
};
