import { ErrorBoundary } from 'react-error-boundary';

import { QueryErrorResetBoundary } from '@tanstack/react-query';

import { isApiError } from '@/api/client';
import { Button } from '@/components/ui/button';

// Shows a retry prompt when a suspense query inside it fails.
export default function ErrorBoundaryFallback(props: {
  children: React.ReactNode;
}) {
  return (
    <QueryErrorResetBoundary>
      {({ reset }) => (
        <ErrorBoundary
          onReset={reset}
          fallbackRender={({ error, resetErrorBoundary }) => (
            <div className="flex items-center gap-3 text-sm">
              <span>
                {isApiError(error) ? error.message : 'Something went wrong.'}
              </span>
              <Button variant="outline" onClick={resetErrorBoundary}>
                Retry
              </Button>
            </div>
          )}
        >
          {props.children}
        </ErrorBoundary>
      )}
    </QueryErrorResetBoundary>
  );
}
