import axios, { type AxiosError, type AxiosRequestConfig } from 'axios';

import { getAuthToken } from '@/lib/auth';
import { env } from '@/lib/env';

import type { ApiErrorBody, ApiErrorDetail } from './types';

export class ApiError extends Error {
  slug: string;
  status: number;
  details?: ApiErrorDetail[];

  constructor(
    message: string,
    opts: { slug: string; status: number; details?: ApiErrorDetail[] },
  ) {
    super(message);
    this.name = 'ApiError';
    this.slug = opts.slug;
    this.status = opts.status;
    this.details = opts.details;
  }
}

export function isApiError(e: unknown): e is ApiError {
  return e instanceof ApiError;
}

export const api = axios.create({
  baseURL: env.VITE_API_URL,
  headers: { 'Content-Type': 'application/json' },
});

api.interceptors.request.use((config) => {
  config.headers.set('Authorization', `Bearer ${getAuthToken()}`);
  return config;
});

// The backend returns the payload as-is (no envelope), so success unwraps
// to the response body and errors are mapped to ApiError.
api.interceptors.response.use(
  (response) => response.data,
  (error: AxiosError<ApiErrorBody>) => {
    // 0 means no HTTP response (network error).
    const status = error.response?.status ?? 0;
    const body = error.response?.data;

    if (body && typeof body.message === 'string') {
      return Promise.reject(
        new ApiError(body.message, {
          slug: body.slug,
          details: body.details,
          status,
        }),
      );
    }

    return Promise.reject(
      new ApiError(error.message || 'Network error', {
        slug: 'network-error',
        status,
      }),
    );
  },
);

// The response interceptor changes the runtime value, so axios' default
// types no longer apply. Feature code uses `http`, which types the body.
export const http = {
  get: <T>(url: string, config?: AxiosRequestConfig) =>
    api.get<T, T>(url, config),
  post: <T, B = unknown>(url: string, body?: B, config?: AxiosRequestConfig) =>
    api.post<T, T, B>(url, body, config),
  put: <T, B = unknown>(url: string, body?: B, config?: AxiosRequestConfig) =>
    api.put<T, T, B>(url, body, config),
  delete: <T = void>(url: string, config?: AxiosRequestConfig) =>
    api.delete<T, T>(url, config),
};
