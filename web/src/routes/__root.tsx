import { createRootRoute, Outlet } from '@tanstack/react-router';

import AuthProvider from '@/components/auth-provider';
import { Toaster } from '@/components/ui/sonner';
import { TooltipProvider } from '@/components/ui/tooltip';

export const Route = createRootRoute({
  component: RootComponent,
});

function RootComponent() {
  return (
    <AuthProvider>
      <TooltipProvider delay={500}>
        <Outlet />
        <Toaster position="top-right" richColors />
      </TooltipProvider>
    </AuthProvider>
  );
}
