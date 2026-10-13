import type React from 'react';

import { Navigate } from '@tanstack/react-router';

import { type Role, useRole } from '@/lib/role';

// Renders the children only for the allowed role; other roles go back to the trainings page.
export default function RoleGate(props: {
  allow: Role;
  children: React.ReactNode;
}) {
  const { role } = useRole();

  if (!role) return null;
  if (role !== props.allow) return <Navigate to="/trainings" replace />;

  return props.children;
}
