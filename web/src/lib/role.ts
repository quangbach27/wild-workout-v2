import { useAuth } from '@clerk/clerk-react';

export type Role = 'trainer' | 'attendee';

// The backend maps the same Clerk active organization role ("org:trainer") to its role.
export function useRole() {
  const { orgRole } = useAuth();
  const parsed = orgRole?.split(':')[1];
  const role: Role | undefined =
    parsed === 'trainer' || parsed === 'attendee' ? parsed : undefined;

  return {
    role,
    isTrainer: role === 'trainer',
    isAttendee: role === 'attendee',
  };
}
