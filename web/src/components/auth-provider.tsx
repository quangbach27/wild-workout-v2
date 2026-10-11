import type React from 'react';

import { ClerkProvider } from '@clerk/clerk-react';
import { useRouter } from '@tanstack/react-router';

import { env } from '@/lib/env';

// Clerk's provider, wired to the router so Clerk navigates client-side. It must render inside the
// router (it is mounted by the root route).
export default function AuthProvider(props: { children: React.ReactNode }) {
  const router = useRouter();

  return (
    <ClerkProvider
      publishableKey={env.VITE_CLERK_PUBLISHABLE_KEY}
      signInUrl="/sign-in"
      signUpUrl="/sign-up"
      signInFallbackRedirectUrl="/trainings"
      signUpFallbackRedirectUrl="/trainings"
      afterSignOutUrl="/sign-in"
      routerPush={(to) => void router.navigate({ href: to })}
      routerReplace={(to) => void router.navigate({ href: to, replace: true })}
    >
      {props.children}
    </ClerkProvider>
  );
}
