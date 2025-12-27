import ky from 'ky';

import { appConfig } from '@/lib/env';

let authToken: string | null = null;

export const setApiAuthToken = (token: string | null) => {
  authToken = token;
};

export const apiClient = ky.create({
  prefixUrl: appConfig.apiBaseUrl,
  credentials: 'include',
  headers: {
    'Content-Type': 'application/json'
  },
  hooks: {
    beforeRequest: [
      (request) => {
        if (authToken) {
          request.headers.set('Authorization', `Bearer ${authToken}`);
        }
      }
    ]
  }
});
