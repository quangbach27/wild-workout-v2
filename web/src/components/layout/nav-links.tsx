import { Link, type LinkProps } from '@tanstack/react-router';

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
];

export default function NavLinks() {
  return (
    <nav className="text-muted-foreground flex-1">
      {NavItems.map((item) => (
        <Button
          key={item.to}
          variant="ghost"
          size="lg"
          className="data-[status=active]:bg-muted data-[status=active]:text-foreground"
          render={<Link to={item.to}>{item.label}</Link>}
        />
      ))}
    </nav>
  );
}
