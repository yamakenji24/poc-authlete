/* eslint-disable @typescript-eslint/no-explicit-any */
export const base64ToArrayBuffer = (base64: string): ArrayBuffer => {
  const binaryString = atob(base64);
  const bytes = new Uint8Array(binaryString.length);
  for (let i = 0; i < binaryString.length; i++) {
    bytes[i] = binaryString.charCodeAt(i);
  }
  return bytes.buffer;
};

export const arrayBufferToBase64 = (buffer: ArrayBuffer): string => {
  const bytes = new Uint8Array(buffer);
  let binary = "";
  for (let i = 0; i < bytes.byteLength; i++) {
    binary += String.fromCharCode(bytes[i]);
  }
  return btoa(binary);
};

export const convertPublicKeyCredentialCreationOptions = (
  options: any
): PublicKeyCredentialCreationOptions => {
  return {
    ...options,
    challenge: base64ToArrayBuffer(options.challenge as string),
    user: {
      ...options.user,
      id: base64ToArrayBuffer(options.user.id as string),
    },
    excludeCredentials: options.excludeCredentials?.map((credential: any) => ({
      ...credential,
      id: base64ToArrayBuffer(credential.id as string),
    })),
  };
};

export const convertPublicKeyCredentialRequestOptions = (
  options: any
): PublicKeyCredentialRequestOptions => {
  return {
    ...options,
    challenge: base64ToArrayBuffer(options.challenge as string),
    allowCredentials: options.allowCredentials?.map((credential: any) => ({
      ...credential,
      id: base64ToArrayBuffer(credential.id as string),
    })),
  };
};
