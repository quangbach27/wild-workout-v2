import { Skeleton } from '@/components/ui/skeleton';

const DAYS = 7;
const HOURS = 9;

// The default `bg-muted` is almost the same colour as the page background, so
// the placeholders use the hour-cell tokens to stay visible.
export default function ScheduleSkeleton() {
  return (
    <div
      className="grid min-w-4xl gap-3 pb-6"
      style={{ gridTemplateColumns: `repeat(${DAYS}, minmax(0, 1fr))` }}
      role="status"
      aria-busy="true"
      aria-label="Loading your hours"
    >
      {Array.from({ length: DAYS }, (_, d) => (
        <div key={d} className="flex flex-col gap-3">
          <div className="flex flex-col items-center gap-2 pb-3">
            <Skeleton className="bg-hour-closed-border h-3.5 w-10" />
            <Skeleton className="bg-hour-closed-border h-8 w-14" />
            <Skeleton className="bg-hour-closed-border h-4 w-16" />
          </div>
          {Array.from({ length: HOURS }, (_, h) => (
            <Skeleton
              key={h}
              className="bg-hour-closed-border h-18 rounded-xl"
              style={{ animationDelay: `${(d + h) * 60}ms` }}
            />
          ))}
        </div>
      ))}
    </div>
  );
}
