import { SignedIn, SignedOut } from '@clerk/clerk-react';
import { Link } from '@tanstack/react-router';
import { ArrowRight, CalendarCheck, Users, Wallet } from 'lucide-react';

import LogoMark from '@/components/layout/logo-mark';
import { Button } from '@/components/ui/button';
import {
  Card,
  CardDescription,
  CardHeader,
  CardTitle,
} from '@/components/ui/card';

const features = [
  {
    icon: Users,
    title: 'Choose your trainer',
    description: 'Browse trainers and pick the one that fits your goals.',
  },
  {
    icon: CalendarCheck,
    title: 'Book in seconds',
    description: 'See free hours across the coming weeks and grab a slot.',
  },
  {
    icon: Wallet,
    title: 'Simple credits',
    description: 'Top up your balance and cancel up to 24h before a session.',
  },
];

function AuthButtons({ size }: { size: 'sm' | 'lg' }) {
  return (
    <>
      <SignedOut>
        <Button
          size={size}
          variant="outline"
          render={<Link to="/sign-up/$" params={{ _splat: '' }} />}
        >
          Sign up
        </Button>
        <Button
          size={size}
          render={<Link to="/sign-in/$" params={{ _splat: '' }} />}
        >
          Log in
          {size === 'lg' && <ArrowRight />}
        </Button>
      </SignedOut>
      <SignedIn>
        <Button size={size} render={<Link to="/trainings" />}>
          Trainings
          {size === 'lg' && <ArrowRight />}
        </Button>
      </SignedIn>
    </>
  );
}

function Header() {
  return (
    <header className="bg-background text-foreground dark flex h-16 items-center justify-between px-4 md:px-6">
      <div className="flex items-center gap-2">
        <LogoMark />
        <span className="hidden text-xl font-bold tracking-wider uppercase sm:block">
          Wild workout
        </span>
      </div>
      <div className="flex items-center gap-2">
        <AuthButtons size="sm" />
      </div>
    </header>
  );
}

function Hero() {
  return (
    <section className="flex flex-col items-center gap-6 px-4 pt-16 pb-12 text-center md:pt-28 md:pb-20">
      <h1 className="max-w-3xl text-4xl font-bold tracking-tight text-balance md:text-6xl">
        Train wild.
        <br />
        <span className="text-hour-action">Book easy.</span>
      </h1>
      <p className="text-muted-foreground max-w-xl text-base text-balance md:text-lg">
        Find a trainer, pick a free hour and schedule your next workout in a few
        clicks.
      </p>
      <div className="flex flex-wrap items-center justify-center gap-3">
        <AuthButtons size="lg" />
      </div>
    </section>
  );
}

function Features() {
  return (
    <section className="mx-auto grid w-full max-w-5xl gap-4 px-4 pb-16 md:grid-cols-3">
      {features.map(({ icon: Icon, title, description }) => (
        <Card key={title}>
          <CardHeader>
            <div className="bg-hour-open text-hour-open-foreground mb-2 flex size-10 items-center justify-center rounded-lg">
              <Icon className="size-5" />
            </div>
            <CardTitle className="text-base">{title}</CardTitle>
            <CardDescription>{description}</CardDescription>
          </CardHeader>
        </Card>
      ))}
    </section>
  );
}

export default function LandingPage() {
  return (
    <div className="bg-background flex min-h-dvh flex-col">
      <Header />
      <main className="flex-1">
        <Hero />
        <Features />
      </main>
      <footer className="text-muted-foreground p-4 text-center text-sm">
        © Wild Workouts 2026
      </footer>
    </div>
  );
}
