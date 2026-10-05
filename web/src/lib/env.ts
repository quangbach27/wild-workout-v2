import * as z from 'zod';

const schema = z.object({
  VITE_API_URL: z.url().default('http://localhost:4000'),
  VITE_MOCK_USER: z.enum(['trainer', 'attendee']).default('trainer'),
});

export const env = schema.parse(import.meta.env);
