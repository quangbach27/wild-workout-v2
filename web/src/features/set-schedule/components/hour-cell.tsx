import type { HourStatus, TrainerHour } from '@/api/trainers/types';
import { cn } from '@/lib/utils';

type CellStatus = 'closed' | 'open' | 'booked';

const STATUS_VIEW: Record<HourStatus, CellStatus> = {
  'not-availability': 'closed',
  availability: 'open',
  'training-scheduled': 'booked',
};

const CELL_STYLES: Record<CellStatus, string> = {
  closed:
    'border-hour-closed-border bg-hour-closed text-hour-closed-foreground',
  open: 'border-hour-open-border bg-hour-open text-hour-open-foreground font-bold',
  booked: 'border-hour-booked bg-hour-booked text-hour-booked-foreground',
};

const hourFormat = new Intl.DateTimeFormat('en-GB', {
  hour: '2-digit',
  minute: '2-digit',
});

export default function HourCell(props: { hour: TrainerHour }) {
  const status = STATUS_VIEW[props.hour.status];

  return (
    <button
      type="button"
      data-status={status}
      disabled={status === 'booked'}
      className={cn(
        'flex h-[4.5rem] w-full flex-col items-center justify-center rounded-xl border text-lg font-medium disabled:cursor-not-allowed',
        CELL_STYLES[status],
      )}
    >
      {hourFormat.format(new Date(props.hour.hour))}
      {status === 'booked' && (
        <span className="text-xs font-normal">booked</span>
      )}
    </button>
  );
}
