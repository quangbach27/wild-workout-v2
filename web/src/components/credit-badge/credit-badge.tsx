import { cn } from 'cn';

type CreditBadgeProps = React.ComponentProps<'div'> & {
  credits: number;
  label?: string;
};

export default function CreditBadge({
  credits,
  label = 'Trainings left',
  className,
  ...props
}: CreditBadgeProps) {
  return (
    <div
      data-slot="credit-badge"
      className={cn(
        'bg-card inline-flex h-11 items-center gap-3 rounded-lg border px-4',
        className,
      )}
      {...props}
    >
      <span className="text-primary text-2xl leading-none font-black">
        {credits}
      </span>
      <span className="text-muted-foreground text-[10px] leading-tight font-semibold tracking-widest uppercase">
        {label}
      </span>
    </div>
  );
}
