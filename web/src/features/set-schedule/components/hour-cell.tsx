import type { HourStatus, TrainerHour } from '@/api/trainers/types';
import { Button } from '@/components/ui/button';
import {
  Tooltip,
  TooltipContent,
  TooltipTrigger,
} from '@/components/ui/tooltip';
import { cn } from '@/lib/utils';

import { useMakeHoursAvailable } from '../api/make-hours-available';
import { useMakeHoursNotAvailable } from '../api/make-hours-not-available';

type CellStatus = 'closed' | 'open' | 'booked';

const STATUS_VIEW: Record<HourStatus, CellStatus> = {
  'not-availability': 'closed',
  availability: 'open',
  'training-scheduled': 'booked',
};

const CELL_STYLES: Record<CellStatus, string> = {
  closed:
    'border-hour-closed-border bg-hour-closed text-hour-closed-foreground hover:bg-hour-closed/80',
  open: 'border-hour-open-border bg-hour-open text-hour-open-foreground font-bold hover:bg-hour-open/80',
  booked: 'border-hour-booked bg-hour-booked text-hour-booked-foreground',
};

const TOOLTIP_TEXT: Record<Exclude<CellStatus, 'booked'>, string> = {
  closed: 'Tap to make available',
  open: 'Tap to make not available',
};

// TODO: point to the trainings page once it exists.
const CANCEL_TRAINING_URL = '#';

const hourFormat = new Intl.DateTimeFormat('en-GB', {
  hour: '2-digit',
  minute: '2-digit',
});

export default function HourCell(props: { hour: TrainerHour }) {
  const status = STATUS_VIEW[props.hour.status];
  const makeAvailable = useMakeHoursAvailable();
  const makeNotAvailable = useMakeHoursNotAvailable();

  // const isHourPast =

  const handleToggleHourCell = () => {
    if (status === 'open') makeNotAvailable.mutate([props.hour.hour]);
    else if (status === 'closed') makeAvailable.mutate([props.hour.hour]);
  };

  return (
    <Tooltip>
      {/* A disabled button emits no hover events, so the span is the trigger. */}
      <TooltipTrigger render={<span className="block" />}>
        <Button
          variant="ghost"
          data-status={status}
          disabled={
            status === 'booked' ||
            makeAvailable.isPending ||
            makeNotAvailable.isPending
          }
          onClick={handleToggleHourCell}
          className={cn(
            'flex h-18 w-full flex-col items-center justify-center rounded-xl border text-lg font-medium disabled:opacity-100',
            CELL_STYLES[status],
          )}
        >
          {hourFormat.format(new Date(props.hour.hour))}
          {status === 'booked' && (
            <span className="text-xs font-normal">booked</span>
          )}
        </Button>
      </TooltipTrigger>
      <TooltipContent className="flex-col items-start gap-0.5">
        {status === 'booked' ? (
          <>
            <span className="font-semibold">Training scheduled</span>
            <span className="text-background/70">
              This hour can&apos;t be changed. Cancel it in{' '}
              <a
                href={CANCEL_TRAINING_URL}
                className="text-background font-medium underline underline-offset-2"
              >
                Trainings
              </a>
              .
            </span>
          </>
        ) : (
          <span className="font-medium">{TOOLTIP_TEXT[status]}</span>
        )}
      </TooltipContent>
    </Tooltip>
  );
}
