import { useAuth } from '@clerk/clerk-react';
import { useQuery } from '@tanstack/react-query';

import { getUserOptions } from '@/features/users/api/get-user-info';

export type Role = 'trainer' | 'attendee';

const parseRole = (value?: string | null): Role | undefined => {
  const role = value?.split(':').pop();
  return role === 'trainer' || role === 'attendee' ? role : undefined;
};

type RoleReturn = {
  role: Role;
  isLoading: boolean;
  isTrainer: boolean;
  isAttendee: boolean;
};

// The backend maps the same Clerk active organization role ("org:trainer") to its role.
export function useRole(): RoleReturn {
  const { userId, orgRole, isLoaded } = useAuth();
  const clerkRole = parseRole(orgRole);

  const shouldFetch = isLoaded && !!userId && !clerkRole;

  const { data, isLoading } = useQuery({
    ...getUserOptions({ userId: userId! }),
    enabled: shouldFetch,
  });

  const role = (clerkRole ?? data?.role) as Role;

  return {
    role,
    isLoading: !isLoaded || (shouldFetch && isLoading),
    isTrainer: role === 'trainer',
    isAttendee: role === 'attendee',
  };
}
