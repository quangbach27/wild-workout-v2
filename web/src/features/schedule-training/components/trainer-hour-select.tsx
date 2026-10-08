import { cn } from 'cn';

import { ToggleGroup, ToggleGroupItem } from '@/components/ui/toggle-group';
import { formatTime } from '@/lib/format-time';

import { useTrainerAvailability } from '../hook/use-trainer-availability';

import { useBookingActions, useBookingState } from './booking-provider';

// Room for the usual rows of hours (h-11 each + gap), so the form doesn't
// jump: three on mobile, two from md.
const HEIGHT = 'min-h-37 md:min-h-24';

export default function TrainerHourSelect(props: { trainerUuid: string }) {
  const { day, hour } = useBookingState();
  const { selectHour } = useBookingActions();
  const days = useTrainerAvailability(props.trainerUuid);

  const selectedDay = days.find((d) => d.date === day);
  if (!selectedDay) {
    return (
      <p className={cn(HEIGHT, 'text-muted-foreground text-sm')}>
        Pick a day first.
      </p>
    );
  }

  return (
    <ToggleGroup
      variant="outline"
      className={cn(HEIGHT, 'w-full flex-wrap content-start')}
      value={hour ? [hour] : []}
      onValueChange={(value) => value[0] && selectHour(value[0])}
    >
      {selectedDay.hours.map((h) => (
        <ToggleGroupItem
          key={h.hour}
          value={h.hour}
          className={cn(
            'w-22',
            'aria-pressed:border-hour-booked aria-pressed:bg-hour-booked aria-pressed:text-hour-booked-foreground aria-pressed:hover:bg-hour-booked/90 h-11 text-base font-semibold',
          )}
        >
          {formatTime(h.hour)}
        </ToggleGroupItem>
      ))}
    </ToggleGroup>
  );
}
