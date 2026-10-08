import { cn } from 'cn';
import { ClockIcon } from 'lucide-react';

import type { Training } from '@/api/trainings/types';
import { DateBadge } from '@/components/date-badge';
import { Button } from '@/components/ui/button';
import { Card } from '@/components/ui/card';
import { Separator } from '@/components/ui/separator';
import { formatTimeRange } from '@/lib/format-time';

import { useCancelTraining } from '../api/cancel-training';

const weekday = new Intl.DateTimeFormat('en-US', { weekday: 'short' });
const dayOfMonth = new Intl.DateTimeFormat('en-US', { day: '2-digit' });
const month = new Intl.DateTimeFormat('en-US', { month: 'short' });

export default function TrainingCard(props: { training: Training }) {
  const { uuid, hour, canceled } = props.training;
  const { mutate: cancel, isPending } = useCancelTraining();
  const date = new Date(hour);

  return (
    <Card
      className={cn(
        'border-border flex-row flex-wrap items-center gap-x-4 gap-y-3 border bg-white px-4 py-3 ring-0',
        canceled && 'opacity-60',
      )}
    >
      <DateBadge
        size="sm"
        weekday={weekday.format(date)}
        day={dayOfMonth.format(date)}
        month={month.format(date)}
      />
      <Separator orientation="vertical" className="h-12" />
      <div className="flex items-center gap-2 text-lg font-semibold md:text-xl">
        <ClockIcon className="size-5" aria-hidden />
        <span>{formatTimeRange(hour)}</span>
        {canceled && (
          <span className="text-muted-foreground text-xs font-semibold tracking-widest uppercase">
            Canceled
          </span>
        )}
      </div>
      {!canceled && (
        <div className="ml-auto flex gap-2">
          <Button variant="outline" size="lg">
            Move
          </Button>
          <Button
            variant="outline"
            size="lg"
            className="text-destructive"
            disabled={isPending}
            onClick={() => cancel(uuid)}
          >
            Cancel
          </Button>
        </div>
      )}
    </Card>
  );
}
