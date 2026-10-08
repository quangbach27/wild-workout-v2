import { Suspense } from 'react';

import { cn } from 'cn';

import ErrorBoundaryFallback from '@/components/query-boundary';
import { Card, CardContent } from '@/components/ui/card';
import { FieldLegend, FieldSet } from '@/components/ui/field';
import { Skeleton } from '@/components/ui/skeleton';
import { Textarea } from '@/components/ui/textarea';

import { useBookingActions, useBookingState } from './booking-provider';
import TrainerDaySelect from './trainer-day-select';
import TrainerHourSelect from './trainer-hour-select';
import TrainerSelect from './trainer-select';

const legendClass =
  'text-muted-foreground mb-3 text-xs font-semibold tracking-widest uppercase';

export default function BookingForm(props: { className?: string }) {
  const { trainerUuid, notes } = useBookingState();
  const { setNotes } = useBookingActions();

  return (
    <Card
      className={cn(
        'border-border min-h-144 border bg-white ring-0',
        props.className,
      )}
    >
      <CardContent className="space-y-6">
        <FieldSet>
          <FieldLegend className={legendClass}>
            1 · Choose a trainer
          </FieldLegend>
          <ErrorBoundaryFallback>
            <Suspense fallback={<Skeleton className="h-10 w-full sm:w-72" />}>
              <TrainerSelect />
            </Suspense>
          </ErrorBoundaryFallback>
        </FieldSet>

        {trainerUuid ? (
          <FieldSet>
            <FieldLegend className={legendClass}>2 · Pick a day</FieldLegend>
            <ErrorBoundaryFallback>
              <Suspense fallback={<Skeleton className="h-20 w-full" />}>
                <TrainerDaySelect trainerUuid={trainerUuid} />
              </Suspense>
            </ErrorBoundaryFallback>
          </FieldSet>
        ) : (
          <p className="text-muted-foreground text-sm">
            Choose a trainer to see their open days and hours.
          </p>
        )}

        {trainerUuid ? (
          <FieldSet>
            <FieldLegend className={legendClass}>3 · Pick a time</FieldLegend>
            <ErrorBoundaryFallback>
              <Suspense fallback={<Skeleton className="h-24 w-full" />}>
                <TrainerHourSelect trainerUuid={trainerUuid} />
              </Suspense>
            </ErrorBoundaryFallback>
          </FieldSet>
        ) : (
          <p className="text-muted-foreground text-sm">Pick a day first.</p>
        )}

        {trainerUuid && (
          <FieldSet>
            <FieldLegend className={legendClass}>
              4 · Notes for the trainer{' '}
              <span className="font-normal tracking-normal normal-case">
                (optional)
              </span>
            </FieldLegend>
            <Textarea
              className="min-h-24 resize-y"
              value={notes ?? ''}
              onChange={(e) => setNotes(e.target.value)}
              placeholder="Anything the trainer should know before the session…"
            />
          </FieldSet>
        )}
      </CardContent>
    </Card>
  );
}
