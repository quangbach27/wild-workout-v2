import type { TrainerDay } from '@/api/trainers/types';
import { Button } from '@/components/ui/button';
import { cn } from '@/lib/utils';

import { useMakeHoursAvailable } from '../api/make-hours-available';

import HourCell from './hour-cell';

const weekdayFormat = new Intl.DateTimeFormat('en-US', { weekday: 'short' });
const dayFormat = new Intl.DateTimeFormat('en-US', { day: '2-digit' });

export default function DayColumn(props: { day: TrainerDay }) {
  const date = new Date(props.day.date);
  const isToday = date.toDateString() === new Date().toDateString();
  const makeAvailable = useMakeHoursAvailable();
  const closedHours = props.day.hours
    .filter((hour) => hour.status === 'not-availability')
    .map((hour) => hour.hour);

  return (
    <div className="flex flex-col gap-2 md:gap-3">
      <div className="bg-background sticky top-0 z-10 pb-3 text-center">
        <div className="flex h-5 items-center justify-center">
          {isToday && (
            <span className="bg-hour-action rounded-full px-2 text-xs font-semibold text-white">
              Today
            </span>
          )}
        </div>
        <div className="text-muted-foreground text-sm font-semibold tracking-widest uppercase">
          {weekdayFormat.format(date)}
        </div>
        <div
          className={cn('text-3xl font-black', isToday && 'text-hour-action')}
        >
          {dayFormat.format(date)}
        </div>
        <Button
          variant="link"
          className="text-hour-action font-semibold"
          disabled={closedHours.length === 0 || makeAvailable.isPending}
          onClick={() => makeAvailable.mutate(closedHours)}
        >
          Open all
        </Button>
      </div>
      {props.day.hours.map((hour) => (
        <HourCell key={hour.hour} hour={hour} />
      ))}
    </div>
  );
}
