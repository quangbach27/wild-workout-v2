declare global {
  interface Window {
    // Set by <ClerkProvider>; only the part used here.
    Clerk?: {
      session?: {
        getToken(options?: { skipCache?: boolean }): Promise<string | null>;
      } | null;
    };
  }
}

// Clerk keeps the long-lived session in an HttpOnly cookie and mints short-lived session JWTs.
// <ClerkProvider> and has a session once the user is signed in, so this also works outside React
export async function getAuthToken(options?: { skipCache?: boolean }) {
  return (await window.Clerk?.session?.getToken(options)) ?? null;
}

export {};
