const DEFAULT_API_BASE_URL = 'http://localhost:8080';

export const appConfig = {
  apiBaseUrl: (import.meta.env.VITE_API_BASE_URL ?? DEFAULT_API_BASE_URL).replace(/\/$/, ''),
  googleClientId: import.meta.env.VITE_GOOGLE_CLIENT_ID
};
