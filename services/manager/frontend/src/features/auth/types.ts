export interface RegisterPayload {
  email: string;
  password: string;
  passwordConfirmation: string;
}

export interface RegisterResponse {
  message: string;
  expiresAt: string;
}

export interface VerifyPayload {
  token: string;
}

export interface VerifiedUser {
  id: string;
  email: string;
  name?: string | null;
  bio?: string | null;
  isActive: boolean;
  emailVerifiedAt: string;
  lastLoginAt?: string | null;
  createdAt: string;
  updatedAt: string;
}

export interface VerifyResponse {
  message: string;
  authToken: string;
  user: VerifiedUser;
}
