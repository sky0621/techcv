import ky from 'ky';

import { appConfig } from '@/lib/env';

export const apiClient = ky.create({
  prefixUrl: appConfig.apiBaseUrl,
  credentials: 'include',
  headers: {
    'Content-Type': 'application/json'
  }
});
