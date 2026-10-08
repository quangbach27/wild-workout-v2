import { useState } from 'react';

import { MenuIcon } from 'lucide-react';

import { Button } from '../ui/button';
import { Sheet, SheetContent, SheetTitle, SheetTrigger } from '../ui/sheet';

import NavLinks from './nav-links';
import UserInfo from './user-info';

// Below md the nav, user and logout live in a drawer instead of the header.
export default function MobileNav() {
  const [open, setOpen] = useState(false);

  return (
    <Sheet open={open} onOpenChange={setOpen}>
      <SheetTrigger
        render={
          <Button
            variant="ghost"
            size="icon-lg"
            aria-label="Open menu"
            className="md:hidden"
          />
        }
      >
        <MenuIcon />
      </SheetTrigger>
      <SheetContent side="left" className="dark w-72 p-4">
        <SheetTitle className="sr-only">Menu</SheetTitle>
        <UserInfo />
        <NavLinks vertical onNavigate={() => setOpen(false)} />
        <Button size="lg" variant="outline" className="mt-auto w-full">
          Logout
        </Button>
      </SheetContent>
    </Sheet>
  );
}
