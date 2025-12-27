const DEFAULT_API_BASE_URL = 'http://localhost:8080';

const getBaseUrl = () => {
  const envValue = process.env.VITE_API_BASE_URL;
  return (envValue && envValue.trim().length > 0 ? envValue : DEFAULT_API_BASE_URL).replace(/\/$/, '');
};

const getEnv = (key: string, fallback?: string) => {
  const value = process.env[key];
  if (value && value.trim().length > 0) {
    return value;
  }
  if (fallback) {
    return fallback;
  }
  throw new Error(`Environment variable ${key} is required for tests`);
};

export const appConfig = {
  apiBaseUrl: getBaseUrl(),
  googleClientId: process.env.VITE_GOOGLE_CLIENT_ID,
  firebase: {
    apiKey: getEnv('VITE_FIREBASE_API_KEY', 'test-api-key'),
    authDomain: getEnv('VITE_FIREBASE_AUTH_DOMAIN', 'localhost'),
    projectId: getEnv('VITE_FIREBASE_PROJECT_ID', 'techcv-dev'),
    appId: getEnv('VITE_FIREBASE_APP_ID', 'test-app-id'),
    messagingSenderId: process.env.VITE_FIREBASE_MESSAGING_SENDER_ID,
    measurementId: process.env.VITE_FIREBASE_MEASUREMENT_ID
  }
};
