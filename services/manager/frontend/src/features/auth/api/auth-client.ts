import ky, { HTTPError } from 'ky';

import { appConfig } from '@/lib/env';

export type AuthResponse = {
  userId: string;
  firebaseUid: string;
  email: string;
  displayName?: string | null;
};

const authClient = ky.create({
  prefixUrl: `${appConfig.apiBaseUrl}/techcv/api/v1`,
  headers: {
    'Content-Type': 'application/json'
  }
});

const postWithToken = (path: string, idToken: string) =>
  authClient
    .post(path, {
      headers: {
        Authorization: `Bearer ${idToken}`
      }
    })
    .json<AuthResponse>();

export const loginWithFirebase = (idToken: string) => postWithToken('auth/firebase/login', idToken);
export const registerWithFirebase = (idToken: string) => postWithToken('auth/firebase/register', idToken);

export const ensureAuthSession = async (idToken: string): Promise<AuthResponse> => {
  try {
    return await loginWithFirebase(idToken);
  } catch (error) {
    if (error instanceof HTTPError && error.response.status === 404) {
      await registerWithFirebase(idToken);
      return loginWithFirebase(idToken);
    }
    throw error;
  }
};
