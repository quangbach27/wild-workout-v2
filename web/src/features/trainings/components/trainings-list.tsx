import { Suspense, startTransition, useEffect } from 'react';

import { useSuspenseQuery } from '@tanstack/react-query';
import { getRouteApi } from '@tanstack/react-router';

import ErrorBoundaryFallback from '@/components/query-boundary';
import { Skeleton } from '@/components/ui/skeleton';
import { TRAININGS_URL } from '@/routes/_app/trainings';

import { getTrainingsOptions } from '../api/get-trainings';

import EmptyTrainings from './empty-trainings';
import TrainingCard from './training-card';
import TrainingsPagination from './trainings-pagination';

const PAGE_SIZE = 10;

const route = getRouteApi(TRAININGS_URL);

function Trainings() {
  const { page } = route.useSearch();
  const navigate = route.useNavigate();
  const { data } = useSuspenseQuery(
    getTrainingsOptions({ page, pageSize: PAGE_SIZE }),
  );

  const totalPages = Math.max(1, Math.ceil(data.pagination.total / PAGE_SIZE));

  const goToPage = (next: number) => {
    window.scrollTo({ top: 0, behavior: 'smooth' });
    startTransition(() => {
      void navigate({ search: { page: next } });
    });
  };

  // Trainings can disappear (canceled, past) so the current page may no longer exist.
  const outOfRange = page > totalPages;
  useEffect(() => {
    if (outOfRange) {
      void navigate({ search: { page: totalPages }, replace: true });
    }
  }, [outOfRange, totalPages, navigate]);

  if (data.pagination.total === 0) {
    return <EmptyTrainings />;
  }

  return (
    <>
      <ul className="flex flex-col gap-3 pb-6">
        {data.items.map((training) => (
          <li key={training.uuid}>
            <TrainingCard training={training} />
          </li>
        ))}
      </ul>
      <TrainingsPagination
        page={page}
        pageSize={PAGE_SIZE}
        total={data.pagination.total}
        onPageChange={goToPage}
      />
    </>
  );
}

export default function TrainingsList() {
  return (
    <ErrorBoundaryFallback>
      <Suspense
        fallback={
          <div className="flex flex-col gap-3" role="status" aria-busy="true">
            {Array.from({ length: 3 }, (_, i) => (
              <Skeleton key={i} className="bg-hour-closed-border h-24 w-full" />
            ))}
          </div>
        }
      >
        <Trainings />
      </Suspense>
    </ErrorBoundaryFallback>
  );
}
