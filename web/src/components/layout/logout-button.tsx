import { useClerk } from '@clerk/clerk-react';

import { Button } from '../ui/button';

export default function LogoutButton(props: { className?: string }) {
  const { signOut } = useClerk();

  return (
    <Button
      size="lg"
      variant="outline"
      className={props.className}
      onClick={() => void signOut()}
    >
      Logout
    </Button>
  );
}
