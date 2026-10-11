import { SignUp } from '@clerk/clerk-react';
import { createFileRoute } from '@tanstack/react-router';

export const Route = createFileRoute('/sign-up/$')({
  component: SignUpPage,
});

// The splat route lets Clerk use sub-paths such as /sign-up/verify-email-address.
function SignUpPage() {
  return (
    <main className="flex min-h-dvh items-center justify-center p-4">
      <SignUp routing="path" path="/sign-up" signInUrl="/sign-in" />
    </main>
  );
}
