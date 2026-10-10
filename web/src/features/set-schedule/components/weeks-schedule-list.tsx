import { Suspense } from 'react';

import { ErrorBoundary } from '@/components/error-boundary';

import ScheduleGrid from './schedule-grid';
import ScheduleLegend from './schedule-legend';
import ScheduleSkeleton from './schedule-skeleton';

export default function WeeksScheduleList() {
  return (
    <>
      <ScheduleLegend />
      <ErrorBoundary fallbackMessage="Could not load your hours.">
        <Suspense fallback={<ScheduleSkeleton />}>
          <ScheduleGrid />
        </Suspense>
      </ErrorBoundary>
    </>
  );
}
