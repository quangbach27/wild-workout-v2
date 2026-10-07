import { cva, type VariantProps } from 'class-variance-authority';
import { cn } from 'cn';

const dateBadgeVariants = cva('flex flex-col items-center text-center', {
  variants: {
    variant: {
      outline: 'bg-accent border-border rounded-xl border',
      transparent: 'rounded-xl bg-transparent',
    },
    size: {
      sm: 'min-w-14 px-2 py-1.5',
      md: 'min-w-16 px-3 py-2',
      lg: 'min-w-20 px-4 py-3',
    },
  },
  defaultVariants: {
    variant: 'transparent',
    size: 'md',
  },
});

const labelSizes = { sm: 'text-xs', md: 'text-sm', lg: 'text-base' } as const;
const daySizes = { sm: 'text-2xl', md: 'text-3xl', lg: 'text-4xl' } as const;

type DateBadgeProps = React.ComponentProps<'div'> &
  VariantProps<typeof dateBadgeVariants> & {
    weekday: string;
    day: string;
    month?: string;
    isToday?: boolean;
  };

export default function DateBadge({
  variant,
  size = 'md',
  weekday,
  day,
  month,
  isToday = false,
  className,
  children,
  ...props
}: DateBadgeProps) {
  const labelClass = cn(
    'text-muted-foreground font-semibold tracking-widest uppercase',
    labelSizes[size ?? 'md'],
  );

  return (
    <div
      data-slot="date-badge"
      data-variant={variant}
      data-today={isToday || undefined}
      className={cn('group', dateBadgeVariants({ variant, size }), className)}
      {...props}
    >
      {/* Reserve the pill's height only in the grid header so columns stay aligned */}
      <div className="flex h-5 items-center justify-center">
        {isToday && (
          <span className="bg-primary rounded-full px-2 text-xs font-semibold text-white">
            Today
          </span>
        )}
      </div>
      <div className={labelClass}>{weekday}</div>
      <div
        className={cn(
          'group-data-today:text-primary leading-tight font-black',
          daySizes[size ?? 'md'],
        )}
      >
        {day}
      </div>
      {month && <div className={labelClass}>{month}</div>}
      {children}
    </div>
  );
}
