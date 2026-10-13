import { SignedIn, SignedOut, SignUp } from '@clerk/clerk-react';
import { createFileRoute, Navigate } from '@tanstack/react-router';

export const Route = createFileRoute('/sign-up/$')({
  component: SignUpPage,
});

// The splat route lets Clerk use sub-paths such as /sign-up/verify-email-address.
function SignUpPage() {
  return (
    <>
      <SignedIn>
        <Navigate to="/trainings" replace />
      </SignedIn>
      <SignedOut>
        <main className="flex min-h-dvh items-center justify-center p-4">
          <SignUp routing="path" path="/sign-up" signInUrl="/sign-in" />
        </main>
      </SignedOut>
    </>
  );
}
