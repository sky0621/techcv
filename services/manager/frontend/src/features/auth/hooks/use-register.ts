import { useMutation } from '@tanstack/react-query';

import { apiClient } from '@/lib/api-client';
import { APIError, parseKyError } from '@/lib/api-error';
import { RegisterPayload, RegisterResponse } from '@/features/auth/types';

interface RegisterApiResponse {
  status: 'success';
  data: {
    message: string;
    expires_at: string;
  };
}

const registerUser = async (payload: RegisterPayload): Promise<RegisterResponse> => {
  try {
    const response = await apiClient
      .post('techcv/api/v1/auth/register', {
        json: {
          email: payload.email,
          password: payload.password,
          password_confirmation: payload.passwordConfirmation
        }
      })
      .json<RegisterApiResponse>();

    return {
      message: response.data.message,
      expiresAt: response.data.expires_at
    };
  } catch (error) {
    return parseKyError(error);
  }
};

export const useRegisterUser = () =>
  useMutation<RegisterResponse, APIError, RegisterPayload>({
    mutationFn: registerUser
  });
