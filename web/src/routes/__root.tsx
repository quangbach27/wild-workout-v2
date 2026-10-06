import { createRootRoute, Outlet } from '@tanstack/react-router';

import { Toaster } from '@/components/ui/sonner';
import { TooltipProvider } from '@/components/ui/tooltip';

export const Route = createRootRoute({
  component: RootComponent,
});

function RootComponent() {
  return (
    <TooltipProvider delay={500}>
      <Outlet />
      <Toaster position="top-right" richColors />
    </TooltipProvider>
  );
}
