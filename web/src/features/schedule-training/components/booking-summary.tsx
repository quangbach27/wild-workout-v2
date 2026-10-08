import { cn } from 'cn';

import { Button } from '@/components/ui/button';
import {
  Card,
  CardContent,
  CardDescription,
  CardFooter,
  CardHeader,
} from '@/components/ui/card';
import { formatTimeRange } from '@/lib/format-time';

import { useScheduleTraining } from '../api/schedule-training';

import { useBookingActions, useBookingState } from './booking-provider';

const dayFormat = new Intl.DateTimeFormat('en-US', {
  weekday: 'short',
  day: '2-digit',
  month: 'short',
});

// A summary row that stays blurred and faded until its value is picked.
function SummaryField(props: {
  label: string;
  value: string | null;
  placeholder: string;
  valueClassName?: string;
}) {
  const selected = props.value !== null;

  return (
    <div>
      <div className="text-muted-foreground text-xs font-semibold tracking-widest uppercase">
        {props.label}
      </div>
      <div
        aria-hidden={!selected}
        className={cn(
          'mt-1 text-xl font-bold tracking-tight transition-[filter,opacity] duration-200 md:text-2xl',
          props.valueClassName,
          !selected && 'opacity-40 blur-sm select-none',
        )}
      >
        {props.value ?? props.placeholder}
      </div>
    </div>
  );
}

export default function BookingSummary() {
  const { trainerUuid, trainerName, day, hour, notes } = useBookingState();
  const { reset } = useBookingActions();
  const { mutate, isPending } = useScheduleTraining();

  const confirm = () => {
    if (!trainerUuid || !trainerName || !hour) return;
    mutate(
      {
        trainerUuid,
        // The trainers API exposes no username, so the display name is sent
        trainerUsername: trainerName,
        hour,
        notes: notes ?? '',
      },
      { onSuccess: reset },
    );
  };

  return (
    <Card className="dark w-full md:sticky md:top-4 md:w-80 md:self-start">
      <CardHeader className="text-muted-foreground text-xs font-semibold tracking-widest uppercase">
        Your Session
      </CardHeader>
      <CardContent className="space-y-4">
        <SummaryField
          label="Trainer"
          value={trainerName}
          placeholder="Alex Johnson"
        />
        <SummaryField
          label="Date"
          value={day ? dayFormat.format(new Date(day)).toUpperCase() : null}
          placeholder="THU 08 OCT"
        />
        <SummaryField
          label="Time"
          value={hour ? formatTimeRange(hour) : null}
          placeholder="21:00 - 22:00"
          valueClassName="text-primary"
        />
        <dl className="border-border space-y-2 border-t pt-4">
          <div className="flex justify-between">
            <dt className="text-muted-foreground">Duration</dt>
            <dd className="font-medium">60 min</dd>
          </div>
          <div className="flex justify-between">
            <dt className="text-muted-foreground">Cost</dt>
            <dd className="font-medium">1 credit</dd>
          </div>
        </dl>
      </CardContent>
      <CardFooter className="flex-col gap-3">
        <Button
          className="w-full"
          size="lg"
          disabled={!trainerUuid || !hour || isPending}
          onClick={confirm}
        >
          Confirm training
        </Button>
        <CardDescription>
          * Free cancellation until 24h before the session.
        </CardDescription>
      </CardFooter>
    </Card>
  );
}
