'use client';
import Link from 'next/link';
import { usePathname } from 'next/navigation';

import { cn } from 'cn';

import { Button } from '../ui/button';

const navItems = [
  {
    title: 'Trainings',
    href: '/trainings',
  },
  {
    title: 'Set Schedule',
    href: '/set-schedule',
  },
];

export default function NavLinks() {
  const pathname = usePathname();

  return (
    <nav className="flex gap-1">
      {navItems.map((item) => {
        const isActive = pathname.startsWith(item.href);
        return (
          <Button
            key={item.href}
            variant="ghost"
            nativeButton={false}
            className={cn(isActive && 'bg-accent text-accent-foreground')}
            render={
              <Link
                href={item.href}
                aria-current={isActive ? 'page' : undefined}
              >
                {item.title}
              </Link>
            }
          />
        );
      })}
    </nav>
  );
}
