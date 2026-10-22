import { useState } from 'react';

import { cn } from 'cn';
import { CalendarIcon, ClockIcon, XIcon } from 'lucide-react';

import type { Training } from '@/api/trainings/types';
import { DateBadge } from '@/components/date-badge';
import {
  AlertDialog,
  AlertDialogAction,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
  AlertDialogTrigger,
} from '@/components/ui/alert-dialog';
import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import { Card } from '@/components/ui/card';
import { Separator } from '@/components/ui/separator';
import { Spinner } from '@/components/ui/spinner';
import { formatTimeRange } from '@/lib/format-time';

import { useCancelTraining } from '../api/cancel-training';

const weekday = new Intl.DateTimeFormat('en-US', { weekday: 'short' });
const dayOfMonth = new Intl.DateTimeFormat('en-US', { day: '2-digit' });
const month = new Intl.DateTimeFormat('en-US', { month: 'short' });
const fullDate = new Intl.DateTimeFormat('en-US', {
  weekday: 'long',
  day: 'numeric',
  month: 'long',
});

export default function TrainingCard(props: { training: Training }) {
  const [open, setOpen] = useState(false);

  const { uuid, hour, canceled } = props.training;
  const { mutate: cancel, isPending } = useCancelTraining();
  const date = new Date(hour);
  const timeRange = formatTimeRange(hour);

  return (
    <Card
      className={cn(
        'border-border flex-row flex-wrap items-center gap-x-4 gap-y-3 border bg-white px-4 py-3 ring-0',
        canceled && 'bg-muted/40',
      )}
    >
      <div className={cn(canceled && 'opacity-50')}>
        <DateBadge
          size="sm"
          weekday={weekday.format(date)}
          day={dayOfMonth.format(date)}
          month={month.format(date)}
        />
      </div>

      <Separator orientation="vertical" className="h-12" />

      <div className="flex items-center gap-3">
        <div
          className={cn(
            'flex items-center gap-2 text-lg font-semibold md:text-xl',
            canceled && 'text-muted-foreground line-through',
          )}
        >
          <ClockIcon className="size-5" aria-hidden />
          <span>{timeRange}</span>
        </div>
        {canceled && <Badge variant="destructive">Canceled</Badge>}
      </div>

      {!canceled && (
        <div className="ml-auto">
          <AlertDialog open={open} onOpenChange={setOpen}>
            <AlertDialogTrigger
              render={
                <Button
                  variant="outline"
                  size="lg"
                  className="text-destructive hover:bg-destructive/10 hover:text-destructive"
                />
              }
            >
              Cancel
            </AlertDialogTrigger>
            <AlertDialogContent>
              <AlertDialogCancel
                render={<Button variant="ghost" size="icon-sm" />}
                nativeButton
                disabled={isPending}
                aria-label="Close"
                className="absolute top-3 right-3"
              >
                <XIcon className="size-4" />
              </AlertDialogCancel>
              <AlertDialogHeader>
                <AlertDialogTitle>Cancel this training?</AlertDialogTitle>
                <AlertDialogDescription>
                  This can't be undone.
                </AlertDialogDescription>
              </AlertDialogHeader>

              <div className="bg-muted/50 flex items-center gap-3 rounded-lg border px-3 py-2.5">
                <CalendarIcon
                  className="text-muted-foreground size-4"
                  aria-hidden
                />
                <div className="text-sm">
                  <p className="font-medium">{fullDate.format(date)}</p>
                  <p className="text-muted-foreground">{timeRange}</p>
                </div>
              </div>

              <AlertDialogFooter>
                <AlertDialogCancel disabled={isPending}>
                  Keep training
                </AlertDialogCancel>
                <AlertDialogAction
                  variant="destructive"
                  disabled={isPending}
                  onClick={() =>
                    cancel(uuid, { onSuccess: () => setOpen(false) })
                  }
                >
                  {isPending && <Spinner className="size-4 text-current" />}
                  Cancel training
                </AlertDialogAction>
              </AlertDialogFooter>
            </AlertDialogContent>
          </AlertDialog>
        </div>
      )}
    </Card>
  );
}
