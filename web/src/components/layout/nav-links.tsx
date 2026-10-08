import { Link, type LinkProps } from '@tanstack/react-router';
import { cn } from 'cn';

import { Button } from '../ui/button';

const NavItems: { to: LinkProps['to']; label: string }[] = [
  {
    to: '/trainings',
    label: 'Trainings',
  },
  {
    to: '/set-schedule',
    label: 'Set schedule',
  },
  {
    to: '/trainings/schedule',
    label: 'Schedule training',
  },
];

type NavLinksProps = {
  className?: string;
  // Stacked full-width links (drawer) instead of an inline row
  vertical?: boolean;
  onNavigate?: () => void;
};

export default function NavLinks(props: NavLinksProps) {
  return (
    <nav
      className={cn(
        'text-muted-foreground flex',
        props.vertical && 'flex-col gap-1',
        props.className,
      )}
    >
      {NavItems.map((item) => (
        <Button
          key={item.to}
          variant="ghost"
          size="lg"
          className={cn(
            'data-[status=active]:bg-muted data-[status=active]:text-foreground',
            props.vertical && 'h-11 justify-start text-base',
          )}
          render={
            <Link
              to={item.to}
              activeOptions={{ exact: true }}
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
