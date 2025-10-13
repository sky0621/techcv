import { useMutation } from '@tanstack/react-query';

import { apiClient } from '@/lib/api-client';
import { APIError, parseKyError } from '@/lib/api-error';
import { VerifyPayload, VerifyResponse } from '@/features/auth/types';

interface VerifyApiResponse {
  status: 'success';
  data: {
    message: string;
    auth_token: string;
    user: {
      id: string;
      email: string;
      name?: string | null;
      bio?: string | null;
      is_active: boolean;
      email_verified_at: string;
      last_login_at?: string | null;
      created_at: string;
      updated_at: string;
    };
  };
}

const verifyRegistration = async (payload: VerifyPayload): Promise<VerifyResponse> => {
  try {
    const response = await apiClient
      .post('techcv/api/v1/auth/verify', {
        json: {
          token: payload.token
        }
      })
      .json<VerifyApiResponse>();

    return {
      message: response.data.message,
      authToken: response.data.auth_token,
      user: {
        id: response.data.user.id,
        email: response.data.user.email,
        name: response.data.user.name ?? null,
        bio: response.data.user.bio ?? null,
        isActive: response.data.user.is_active,
        emailVerifiedAt: response.data.user.email_verified_at,
        lastLoginAt: response.data.user.last_login_at ?? null,
        createdAt: response.data.user.created_at,
        updatedAt: response.data.user.updated_at
      }
    };
  } catch (error) {
    return parseKyError(error);
  }
};

export const useVerifyRegistration = () =>
  useMutation<VerifyResponse, APIError, VerifyPayload>({
    mutationFn: verifyRegistration
  });
