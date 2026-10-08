import { ChevronLeft, ChevronRight } from 'lucide-react';

import { Button } from '@/components/ui/button';

type Props = {
  page: number;
  pageSize: number;
  total: number;
  onPageChange: (page: number) => void;
};

export default function TrainingsPagination({
  page,
  pageSize,
  total,
  onPageChange,
}: Props) {
  const totalPages = Math.ceil(total / pageSize);
  if (totalPages <= 1) {
    return null;
  }

  const from = (page - 1) * pageSize + 1;
  const to = Math.min(page * pageSize, total);

  return (
    <nav
      aria-label="Trainings pagination"
      className="flex flex-col items-center gap-2 pb-6"
    >
      <div className="flex items-center gap-3">
        <Button
          variant="outline"
          size="lg"
          aria-label="Previous page"
          disabled={page <= 1}
          onClick={() => onPageChange(page - 1)}
        >
          <ChevronLeft data-icon="inline-start" />
          Previous
        </Button>
        <span className="min-w-24 text-center text-sm font-medium">
          Page {page} of {totalPages}
        </span>
        <Button
          variant="outline"
          size="lg"
          aria-label="Next page"
          disabled={page >= totalPages}
          onClick={() => onPageChange(page + 1)}
        >
          Next
          <ChevronRight data-icon="inline-end" />
        </Button>
      </div>
      <p className="text-muted-foreground text-xs">
        Showing {from}–{to} of {total}
      </p>
    </nav>
  );
}
