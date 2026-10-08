import type { TrainerDay } from '@/api/trainers/types';
import { DateBadge } from '@/components/date-badge';
import { Button } from '@/components/ui/button';

import { useMakeHoursAvailable } from '../api/make-hours-available';

import HourCell from './hour-cell';

const weekdayFormat = new Intl.DateTimeFormat('en-US', { weekday: 'short' });
const dayFormat = new Intl.DateTimeFormat('en-US', { day: '2-digit' });

export default function DayColumn(props: { day: TrainerDay }) {
  const makeAvailable = useMakeHoursAvailable();

  const date = new Date(props.day.date);
  const isToday = date.toDateString() === new Date().toDateString();
  const closedHours = props.day.hours
    .filter((hour) => hour.status === 'not-availability')
    .map((hour) => hour.hour);

  return (
    <div className="flex flex-col gap-2 md:gap-3">
      <div className="bg-background sticky top-0 z-10">
        {/* Reserve the "Today" pill row so every column's hours line up */}
        {!isToday && <div aria-hidden className="h-5" />}
        <DateBadge
          variant="transparent"
          day={dayFormat.format(date)}
          isToday={isToday}
          weekday={weekdayFormat.format(date)}
        >
          <Button
            variant="link"
            className="text-hour-action font-semibold"
            disabled={closedHours.length === 0 || makeAvailable.isPending}
            onClick={() => makeAvailable.mutate(closedHours)}
          >
            Open all
          </Button>
        </DateBadge>
      </div>
      {props.day.hours.map((hour) => (
        <HourCell key={hour.hour} hour={hour} />
      ))}
    </div>
  );
}
