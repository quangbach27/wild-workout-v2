import Link from 'next/link';

import { cn } from 'cn';

import { Avatar, AvatarFallback } from '../ui/avatar';
import { Button } from '../ui/button';

import Logo from './logo';
import NavLinks from './nav-links';

export default function Header(props: { className?: string }) {
  return (
    <header
      className={cn(
        'bg-background text-foreground dark flex h-16 min-h-3.5 items-center justify-between px-2 md:px-5',
        props.className,
      )}
    >
      <div className="flex">
        <Logo />
        <NavLinks />
      </div>
      <div className="flex items-center gap-4">
        <Avatar>
          <AvatarFallback>TR</AvatarFallback>
        </Avatar>
        <span className="hidden md:block">Trainer</span>
        <Button>Logout</Button>
      </div>
    </header>
  );
}
