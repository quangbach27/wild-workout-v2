import type { components } from './trainers/generated/schema';

// Every module returns the same error shape.
export type ApiErrorBody = components['schemas']['Error'];
export type ApiErrorDetail = components['schemas']['ErrorDetail'];
