const DEFAULT_API_BASE_URL = 'http://localhost:8080';

const requireEnv = (key: keyof ImportMetaEnv) => {
  const value = import.meta.env[key];
  if (!value || value.trim().length === 0) {
    throw new Error(`Environment variable ${key} is required`);
  }
  return value;
};

export const appConfig = {
  apiBaseUrl: (import.meta.env.VITE_API_BASE_URL ?? DEFAULT_API_BASE_URL).replace(/\/$/, ''),
  googleClientId: import.meta.env.VITE_GOOGLE_CLIENT_ID,
  firebase: {
    apiKey: requireEnv('VITE_FIREBASE_API_KEY'),
    authDomain: requireEnv('VITE_FIREBASE_AUTH_DOMAIN'),
    projectId: requireEnv('VITE_FIREBASE_PROJECT_ID'),
    appId: requireEnv('VITE_FIREBASE_APP_ID'),
    messagingSenderId: import.meta.env.VITE_FIREBASE_MESSAGING_SENDER_ID,
    measurementId: import.meta.env.VITE_FIREBASE_MEASUREMENT_ID
  }
};
