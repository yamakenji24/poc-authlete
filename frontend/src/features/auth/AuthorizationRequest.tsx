import React, { useEffect, useState } from "react";
import { authApi } from "../../api/auth";

export const AuthorizationRequest: React.FC = () => {
  const [error, setError] = useState<string>("");

  useEffect(() => {
    const requestAuthorization = async () => {
      try {
        await authApi.requestAuthorization();
      } catch (err) {
        console.error(err);
        setError("認可リクエストに失敗しました。");
      }
    };
    requestAuthorization();
  }, []);

  if (error) {
    return (
      <div className="min-h-screen flex items-center justify-center bg-gray-50">
        <div className="max-w-md w-full space-y-8 p-8 bg-white rounded-lg shadow">
          <div className="text-center">
            <h2 className="text-3xl font-bold text-gray-900">エラー</h2>
            <p className="mt-2 text-red-600">{error}</p>
          </div>
        </div>
      </div>
    );
  }

  return (
    <div className="min-h-screen flex items-center justify-center bg-gray-50">
      <div className="max-w-md w-full space-y-8 p-8 bg-white rounded-lg shadow">
        <div className="text-center">
          <h2 className="text-3xl font-bold text-gray-900">認可リクエスト中</h2>
          <p className="mt-2 text-gray-600">
            認可サーバーにリダイレクトしています...
          </p>
        </div>
      </div>
    </div>
  );
};
