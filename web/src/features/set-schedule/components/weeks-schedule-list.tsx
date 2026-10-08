import { Suspense } from 'react';

import { ErrorBoundary } from 'react-error-boundary';

import { QueryErrorResetBoundary } from '@tanstack/react-query';

import { isApiError } from '@/api/client';
import { Button } from '@/components/ui/button';

import ScheduleGrid from './schedule-grid';
import ScheduleLegend from './schedule-legend';
import ScheduleSkeleton from './schedule-skeleton';

export default function WeeksScheduleList() {
  return (
    <>
      <ScheduleLegend />
      <QueryErrorResetBoundary>
        {({ reset }) => (
          <ErrorBoundary
            onReset={reset}
            fallbackRender={({ error, resetErrorBoundary }) => (
              <div className="flex items-center gap-3">
                <span>
                  {isApiError(error)
                    ? error.message
                    : 'Could not load your hours.'}
                </span>
                <Button variant="outline" onClick={resetErrorBoundary}>
                  Retry
                </Button>
              </div>
            )}
          >
            <Suspense fallback={<ScheduleSkeleton />}>
              <ScheduleGrid />
            </Suspense>
          </ErrorBoundary>
        )}
      </QueryErrorResetBoundary>
    </>
  );
}
