import { SignedIn, SignedOut } from '@clerk/clerk-react';
import { createFileRoute, Navigate, Outlet } from '@tanstack/react-router';

import AppLayoutComponent from '@/components/layout/app-layout';

export const Route = createFileRoute('/_app')({
  component: AppLayout,
});

function AppLayout() {
  return (
    <>
      <SignedIn>
        <AppLayoutComponent>
          <Outlet />
        </AppLayoutComponent>
      </SignedIn>
      <SignedOut>
        <Navigate to="/" replace />
      </SignedOut>
    </>
  );
}
