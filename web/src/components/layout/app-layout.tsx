import type React from 'react';

import { Avatar, AvatarFallback } from '../ui/avatar';
import { Button } from '../ui/button';

import NavLinks from './nav-links';

export default function AppLayout(props: { children: React.ReactNode }) {
  return (
    <div className="flex h-dvh flex-col">
      <header className="bg-background text-foreground dark flex h-16 items-center justify-between gap-5 px-2 md:px-6">
        <div className="flex gap-2">
          <svg width="44" height="24" viewBox="0 0 48 26" fill="none">
            <path
              d="M3 5.5 L8.74 23 L13.5 1.8 L18.26 23 L24 5.5"
              stroke="#E0491F"
              strokeWidth="3.6"
              strokeLinecap="round"
              strokeLinejoin="miter"
              strokeMiterlimit="8"
            ></path>
            <path
              d="M3 5.5 L8.74 23 L13.5 1.8 L18.26 23 L24 5.5"
              transform="translate(21 0)"
              stroke="#F7F5F1"
              strokeWidth="3.6"
              strokeLinecap="round"
              strokeLinejoin="miter"
              strokeMiterlimit="8"
            ></path>
          </svg>
          <h1 className="hidden text-xl font-bold tracking-wider uppercase md:block">
            Wild workout
          </h1>
        </div>
        <NavLinks />
        <div className="flex items-center gap-3">
          <Avatar className="hidden md:block">
            <AvatarFallback>Tr</AvatarFallback>
          </Avatar>
          <span className="hidden md:block">Trainer</span>
          <Button size="lg" variant="outline">
            Logout
          </Button>
        </div>
      </header>

      <div className="bg-background container-app flex min-h-0 flex-1 flex-col px-3 md:px-8">
        {props.children}
      </div>

      <footer className="text-muted-foreground space-y-2 p-2 md:p-8">
        <div className="hidden gap-3 md:flex">
          <span>© Wild Workouts 2026</span>
          <a
            href="https://github.com/quangbach27/wild-workout-v2"
            target="_blank"
            rel="noopener"
            className="flex items-center gap-2"
          >
            <svg
              width="14"
              height="14"
              viewBox="0 0 16 16"
              fill="currentColor"
              aria-hidden="true"
              className="inline-block"
            >
              <path
                fill-rule="evenodd"
                clip-rule="evenodd"
                d="M8 0C3.58 0 0 3.58 0 8c0 3.54 2.29 6.53 5.47 7.59.4.07.55-.17.55-.38 0-.19-.01-.82-.01-1.49-2.01.37-2.53-.49-2.69-.94-.09-.23-.48-.94-.82-1.13-.28-.15-.68-.52-.01-.53.63-.01 1.08.58 1.23.82.72 1.21 1.87.87 2.33.66.07-.52.28-.87.51-1.07-1.78-.2-3.64-.89-3.64-3.95 0-.87.31-1.59.82-2.15-.08-.2-.36-1.02.08-2.12 0 0 .67-.21 2.2.82.64-.18 1.32-.27 2-.27.68 0 1.36.09 2 .27 1.53-1.04 2.2-.82 2.2-.82.44 1.1.16 1.92.08 2.12.51.56.82 1.27.82 2.15 0 3.07-1.87 3.75-3.65 3.95.29.25.54.73.54 1.48 0 1.07-.01 1.93-.01 2.2 0 .21.15.46.55.38A8.01 8.01 0 0 0 16 8c0-4.42-3.58-8-8-8Z"
              ></path>
            </svg>
            <span className="font-bold">Source on GitHub</span>
          </a>
        </div>
        <p>
          * Cancelling less than 24h before a session doesn't return the credit
        </p>
      </footer>
    </div>
  );
}
