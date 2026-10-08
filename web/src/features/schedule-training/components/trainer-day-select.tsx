import { cn } from 'cn';

import { DateBadge } from '@/components/date-badge';
import { ToggleGroup, ToggleGroupItem } from '@/components/ui/toggle-group';

import { useTrainerAvailability } from '../hook/use-trainer-availability';

import { useBookingActions, useBookingState } from './booking-provider';

const weekday = new Intl.DateTimeFormat('en-US', { weekday: 'short' });
const dayOfMonth = new Intl.DateTimeFormat('en-US', { day: '2-digit' });
const month = new Intl.DateTimeFormat('en-US', { month: 'short' });

export default function TrainerDaySelect(props: { trainerUuid: string }) {
  const { day } = useBookingState();
  const { selectDay } = useBookingActions();
  const days = useTrainerAvailability(props.trainerUuid);

  if (days.length === 0) {
    return (
      <p className="text-muted-foreground text-sm">
        This trainer has no open hours.
      </p>
    );
  }

  return (
    <ToggleGroup
      className="w-full flex-wrap items-stretch"
      value={day ? [day] : []}
      onValueChange={(value) => value[0] && selectDay(value[0])}
    >
      {days.map((d) => {
        const date = new Date(d.date);
        const isToday = date.toDateString() === new Date().toDateString();
        return (
          <ToggleGroupItem
            key={d.date}
            value={d.date}
            className={cn(
              'w-22',
              'h-auto rounded-xl bg-transparent p-0 hover:bg-transparent aria-pressed:bg-transparent',
            )}
          >
            <DateBadge
              variant="outline"
              size="sm"
              isToday={isToday}
              weekday={weekday.format(date)}
              day={dayOfMonth.format(date)}
              month={month.format(date)}
              className={cn(
                'group-hover/toggle:border-foreground/30 group-aria-pressed/toggle:border-hour-booked group-aria-pressed/toggle:bg-hour-booked group-aria-pressed/toggle:text-hour-booked-foreground h-full w-full justify-center bg-white group-aria-pressed/toggle:[&_div]:text-inherit',
                // Keep today's number orange when selected
                isToday &&
                  'group-aria-pressed/toggle:[&_div.font-black]:text-primary',
              )}
            />
          </ToggleGroupItem>
        );
      })}
    </ToggleGroup>
  );
}
