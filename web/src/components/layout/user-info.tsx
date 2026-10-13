import { useUser } from '@clerk/clerk-react';

import { useRole } from '@/lib/role';

import { Avatar, AvatarFallback, AvatarImage } from '../ui/avatar';

// Who is signed in: avatar, name and role (the backend sets the role in Clerk's public metadata).
export default function UserInfo() {
  const { user } = useUser();
  const name = user?.fullName ?? '';
  const { role } = useRole();

  return (
    <div className="flex items-center gap-3">
      <Avatar>
        <AvatarImage src={user?.imageUrl} alt={name} />
        <AvatarFallback>{name.slice(0, 2).toUpperCase()}</AvatarFallback>
      </Avatar>
      <div className="flex flex-col leading-tight">
        <span>{name}</span>
        {role && (
          <span className="text-muted-foreground text-xs capitalize">
            {role}
          </span>
        )}
      </div>
    </div>
  );
}
