import { SignedIn, SignedOut, SignIn } from '@clerk/clerk-react';
import { createFileRoute, Navigate } from '@tanstack/react-router';

export const Route = createFileRoute('/sign-in/$')({
  component: SignInPage,
});

// The splat route lets Clerk use sub-paths such as /sign-in/factor-one.
function SignInPage() {
  return (
    <>
      <SignedIn>
        <Navigate to="/trainings" replace />
      </SignedIn>
      <SignedOut>
        <main className="flex min-h-dvh items-center justify-center p-4">
          <SignIn routing="path" path="/sign-in" signUpUrl="/sign-up" />
        </main>
      </SignedOut>
    </>
  );
}
