import { createFileRoute, Outlet } from '@tanstack/react-router';

import AppLayoutComponent from '@/components/layout/app-layout';

export const Route = createFileRoute('/_app')({
  component: AppLayout,
});

function AppLayout() {
  return (
    <AppLayoutComponent>
      <Outlet />
    </AppLayoutComponent>
  );
}
