import { ErrorBoundary as ReactErrorBoundary } from 'react-error-boundary';

import { QueryErrorResetBoundary } from '@tanstack/react-query';

import { isApiError } from '@/api/client';
import { Button } from '@/components/ui/button';

// Shows a retry prompt when a suspense query inside it fails.
export default function ErrorBoundary(props: {
  children: React.ReactNode;
  fallbackMessage?: string;
}) {
  return (
    <QueryErrorResetBoundary>
      {({ reset }) => (
        <ReactErrorBoundary
          onReset={reset}
          fallbackRender={({ error, resetErrorBoundary }) => (
            <div className="flex items-center gap-3 text-sm">
              <span>
                {isApiError(error)
                  ? error.message
                  : (props.fallbackMessage ?? 'Something went wrong.')}
              </span>
              <Button variant="outline" onClick={resetErrorBoundary}>
                Retry
              </Button>
            </div>
          )}
        >
          {props.children}
        </ReactErrorBoundary>
      )}
    </QueryErrorResetBoundary>
  );
}
