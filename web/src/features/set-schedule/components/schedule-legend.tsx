import { cn } from '@/lib/utils';

const ITEMS = [
  { label: 'Closed', swatch: 'border-hour-closed-border bg-hour-closed' },
  {
    label: 'Open for booking',
    swatch: 'border-hour-open-border bg-hour-open',
  },
  {
    label: 'Training booked — locked',
    swatch: 'border-hour-booked bg-hour-booked',
  },
];

export default function ScheduleLegend() {
  return (
    <ul className="text-muted-foreground flex flex-col pb-2 md:flex-row md:flex-wrap md:gap-x-6 md:gap-y-2 md:pb-6">
      {ITEMS.map((item) => (
        <li key={item.label} className="flex items-center gap-2">
          <span
            aria-hidden="true"
            className={cn('size-5 rounded-md border', item.swatch)}
          />
          {item.label}
        </li>
      ))}
    </ul>
  );
}
