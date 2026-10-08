import { Avatar, AvatarFallback } from '../ui/avatar';

// Who is signed in: avatar and role.
export default function UserInfo() {
  return (
    <div className="flex items-center gap-3">
      <Avatar>
        <AvatarFallback>Tr</AvatarFallback>
      </Avatar>
      <span>Trainer</span>
    </div>
  );
}
