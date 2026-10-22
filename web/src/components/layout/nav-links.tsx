import { useMemo } from 'react';

import { Link } from '@tanstack/react-router';
import { cn } from 'cn';

import { Button } from '../ui/button';
import { useRole, type Role } from '@/hooks/use-role';

const navItems: { to: string; label: string; roles: Role[] }[] = [
  {
    to: '/trainings',
    label: 'Trainings',
    roles: ['trainer', 'attendee'],
  },
  {
    to: '/set-schedule',
    label: 'Set schedule',
    roles: ['trainer'],
  },
  {
    to: '/trainings/schedule',
    label: 'Schedule training',
    roles: ['attendee'],
  },
];

type NavLinksProps = {
  className?: string;
  // Stacked full-width links (drawer) instead of an inline row
  vertical?: boolean;
  onNavigate?: () => void;
};

export default function NavLinks(props: NavLinksProps) {
  const { role } = useRole();

  // Until the role is known only the page every role has is shown
  const items = useMemo(
    () => (role ? navItems.filter((item) => item.roles.includes(role)) : []),
    [role],
  );

  return (
    <nav
      className={cn(
        'text-muted-foreground flex',
        props.vertical && 'flex-col gap-1',
        props.className,
      )}
    >
      {items.map((item) => (
        <Button
          key={item.to}
          variant="ghost"
          size="lg"
          className={cn(
            'data-[status=active]:bg-muted data-[status=active]:text-foreground',
            props.vertical && 'h-11 justify-start text-base',
          )}
          nativeButton={false}
          render={
            <Link
              to={item.to}
              activeOptions={{ exact: true, includeSearch: false }}
              onClick={props.onNavigate}
            >
              {item.label}
            </Link>
          }
        />
      ))}
    </nav>
  );
}
