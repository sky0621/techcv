import { HTTPError } from 'ky';

export interface ApiErrorDetail {
  field?: string;
  code?: string;
  message?: string;
}

export interface ApiErrorBody {
  requestId: string;
  code?: string;
  message?: string;
  details?: ApiErrorDetail[];
}

export interface ApiErrorResponse {
  status: 'error';
  error?: ApiErrorBody;
}

export class APIError extends Error {
  public readonly requestId?: string;
  public readonly code?: string;
  public readonly status?: number;
  public readonly details: ApiErrorDetail[];

  constructor(options: { message: string; code?: string; requestId?: string; status?: number; details?: ApiErrorDetail[] }) {
    super(options.message);
    this.name = 'APIError';
    this.code = options.code;
    this.requestId = options.requestId;
    this.status = options.status;
    this.details = options.details ?? [];
  }
}

export const parseKyError = async (error: unknown): Promise<never> => {
  if (error instanceof HTTPError) {
    try {
      const payload = (await error.response.json()) as ApiErrorResponse;
      const body = payload.error;
      throw new APIError({
        message: body?.message ?? error.message,
        code: body?.code,
        requestId: body?.requestId,
        status: error.response.status,
        details: body?.details ?? []
      });
    } catch (parseError) {
      throw new APIError({
        message: error.message,
        status: error.response.status
      });
    }
  }

  throw error;
};
