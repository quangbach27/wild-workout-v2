import { env } from './env';

// Dev-only tokens. The backend registers the same tokens in cmd/main.go
// when APP_ENV=dev. Replace with real login later.
export const MOCK_TOKENS = {
  trainer: 'mock-trainer-token',
  attendee: 'mock-attendee-token',
} as const;

export function getAuthToken() {
  return MOCK_TOKENS[env.VITE_MOCK_USER];
}
