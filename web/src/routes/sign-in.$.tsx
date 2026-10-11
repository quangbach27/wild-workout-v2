import { SignIn } from '@clerk/clerk-react';
import { createFileRoute } from '@tanstack/react-router';

export const Route = createFileRoute('/sign-in/$')({
  component: SignInPage,
});

// The splat route lets Clerk use sub-paths such as /sign-in/factor-one.
function SignInPage() {
  return (
    <main className="flex min-h-dvh items-center justify-center p-4">
      <SignIn routing="path" path="/sign-in" signUpUrl="/sign-up" />
    </main>
  );
}
