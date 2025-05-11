import React, { useState } from "react";
import { useAuth } from "../contexts/AuthContext";
import { passkeyApi } from "../api/passkey";
import { convertPublicKeyCredentialCreationOptions } from "../utils/webauthn";

interface Passkey {
  id: string;
  userId: string;
  name?: string;
  createdAt: string;
  lastUsedAt?: string;
}

export const AccountSettings: React.FC = () => {
  const { user } = useAuth();
  const [error, setError] = useState<string | null>(null);
  const [isRegistering, setIsRegistering] = useState(false);
  const [registeredPasskeys, setRegisteredPasskeys] = useState<Passkey[]>([]);

  // パスキー登録処理
  const handleRegisterPasskey = async () => {
    try {
      setIsRegistering(true);
      setError(null);

      const options = await passkeyApi.startRegistration(user?.name || "");
      console.log("options: ", options);
      const convertedOptions = convertPublicKeyCredentialCreationOptions(options.publicKey);
      const credential = (await navigator.credentials.create({
        publicKey: convertedOptions,
      })) as PublicKeyCredential;

      if (!credential) {
        throw new Error("パスキーの登録に失敗しました");
      }

      const response = credential.response as AuthenticatorAttestationResponse;
      const publicKey = response.getPublicKey();
      if (!publicKey) {
        throw new Error("公開鍵の取得に失敗しました");
      }

      const publicKeyStr = btoa(
        String.fromCharCode(...new Uint8Array(publicKey))
      );

      await passkeyApi.completeRegistration(
        user?.sub || "",
        credential.id,
        publicKeyStr,
        "none"
      );

      const passkeys = await passkeyApi.getPasskeys(user?.sub || "");
      setRegisteredPasskeys(passkeys);

      alert("パスキーの登録が完了しました");
    } catch (err) {
      setError("パスキーの登録に失敗しました");
      console.error(err);
    } finally {
      setIsRegistering(false);
    }
  };

  // パスキー削除処理
  const handleDeletePasskey = async (passkeyId: string) => {
    try {
      await passkeyApi.deletePasskey(passkeyId);
      const passkeys = await passkeyApi.getPasskeys(user?.sub || "");
      setRegisteredPasskeys(passkeys);
    } catch (err) {
      setError("パスキーの削除に失敗しました");
      console.error(err);
    }
  };

  return (
    <div className="container mx-auto px-4 py-8">
      <h1 className="text-2xl font-bold mb-6">アカウント設定</h1>

      <div className="bg-white shadow rounded-lg p-6 mb-6">
        <h2 className="text-xl font-semibold mb-4">パスキー設定</h2>
        <p className="mb-4">パスキーを使用して、より安全にログインできます。</p>

        <button
          onClick={handleRegisterPasskey}
          disabled={isRegistering}
          className="bg-blue-500 hover:bg-blue-600 text-white font-bold py-2 px-4 rounded disabled:opacity-50"
        >
          {isRegistering ? "登録中..." : "パスキーを登録"}
        </button>

        {error && <p className="mt-4 text-red-500">{error}</p>}

        {/* 登録済みパスキーの一覧 */}
        {registeredPasskeys.length > 0 && (
          <div className="mt-6">
            <h3 className="text-lg font-semibold mb-2">登録済みのパスキー</h3>
            <ul className="space-y-2">
              {registeredPasskeys.map((passkey) => (
                <li
                  key={passkey.id}
                  className="flex items-center justify-between p-3 bg-gray-50 rounded"
                >
                  <span>{passkey.name || "パスキー"}</span>
                  <button
                    onClick={() => handleDeletePasskey(passkey.id)}
                    className="text-red-500 hover:text-red-600"
                  >
                    削除
                  </button>
                </li>
              ))}
            </ul>
          </div>
        )}
      </div>
    </div>
  );
};
