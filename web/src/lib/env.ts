import { z } from 'zod';

const schema = z.object({
  VITE_API_URL: z.url().default('http://localhost:4000'),
  VITE_FUTURE_WEEKS: z.coerce.number().int().positive().default(6),
  VITE_CLERK_PUBLISHABLE_KEY: z
    .string()
    .trim()
    .min(1, 'Missing VITE_CLERK_PUBLISHABLE_KEY')
    .startsWith('pk_', 'Clerk publishable key must start with "pk_"'),
});

const parsed = schema.safeParse(import.meta.env);
if (!parsed.success) {
  console.error('Invalid environment variables:', z.treeifyError(parsed.error));
  throw new Error('Invalid environment variables');
}

export const env = parsed.data;
