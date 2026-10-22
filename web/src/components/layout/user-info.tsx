import { useAuth } from '@clerk/clerk-react';
import { useQuery } from '@tanstack/react-query';

import { getUserOptions } from '@/features/users/api/get-user-info';

import { Avatar, AvatarFallback } from '../ui/avatar';

// Who is signed in: avatar, name and role (the backend sets the role in Clerk's public metadata).
export default function UserInfo() {
  const { userId } = useAuth();
  const { data: user } = useQuery(getUserOptions({ userId: userId! }));

  return (
    <div className="flex items-center gap-3">
      <Avatar>
        <AvatarFallback className="bg-green-700 text-white">
          {user?.displayName.slice(0, 2).toUpperCase()}
        </AvatarFallback>
      </Avatar>
      <div className="flex flex-col leading-tight">
        <span>{user?.displayName}</span>
        <span className="text-muted-foreground text-xs capitalize">
          {user?.role}
        </span>
      </div>
    </div>
  );
}
