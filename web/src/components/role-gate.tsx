import type React from 'react';

import { useClerk } from '@clerk/clerk-react';
import { Navigate } from '@tanstack/react-router';
import { cn } from 'cn';
import { Loader2 } from 'lucide-react';

import { useRole, type Role } from '@/hooks/use-role';

import { Button } from './ui/button';

const NoRoleMessage = () => {
  const { signOut } = useClerk();

  return (
    <div className="flex min-h-[60vh] flex-col items-center justify-center gap-4 px-4 text-center">
      <h1 className="text-xl font-semibold">No role assigned</h1>
      <p className="max-w-md text-sm text-gray-600">
        Your account doesn't have a role in this organization yet. Please
        contact your organization admin to get access.
      </p>

      <div className="flex gap-2">
        <Button
          onClick={() => window.location.reload()}
          // className="rounded-md border px-4 py-2 text-sm"
        >
          Try again
        </Button>
        <Button
          onClick={() => signOut({ redirectUrl: '/' })}
          // className="rounded-md bg-black px-4 py-2 text-sm text-white"
        >
          Sign out
        </Button>
      </div>
    </div>
  );
};

const Spinner = ({ className }: { className?: string }) => {
  return (
    <Loader2
      role="status"
      aria-label="Loading"
      className={cn('text-muted-foreground size-6 animate-spin', className)}
    />
  );
};

const PageSpinner = () => {
  return (
    <div className="flex min-h-[60vh] items-center justify-center">
      <Spinner />
    </div>
  );
};

// Renders the children only for the allowed role; other roles go back to the trainings page.
export default function RoleGate(props: {
  allow: Role;
  children: React.ReactNode;
}) {
  const { role, isLoading } = useRole();

  if (isLoading) return <PageSpinner />;

  if (!role) return <NoRoleMessage />;

  if (role !== props.allow) return <Navigate to="/trainings" replace />;

  return props.children;
}
