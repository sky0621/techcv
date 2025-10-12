const DEFAULT_API_BASE_URL = 'http://localhost:8080';

const getBaseUrl = () => {
  const envValue = process.env.VITE_API_BASE_URL;
  return (envValue && envValue.trim().length > 0 ? envValue : DEFAULT_API_BASE_URL).replace(/\/$/, '');
};

export const appConfig = {
  apiBaseUrl: getBaseUrl(),
  googleClientId: process.env.VITE_GOOGLE_CLIENT_ID
};
