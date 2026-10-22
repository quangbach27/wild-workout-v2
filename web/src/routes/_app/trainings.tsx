import { useAuth } from '@clerk/clerk-react';
import { useQuery } from '@tanstack/react-query';
import { createFileRoute, Link } from '@tanstack/react-router';
import * as z from 'zod';

import { CreditBadge } from '@/components/credit-badge';
import ContentLayout from '@/components/layout/content-layout';
import { Button } from '@/components/ui/button';
import TrainingsList from '@/features/trainings/components/trainings-list';
import { getUserOptions } from '@/features/users/api/get-user-info';
import { useRole } from '@/hooks/use-role';

export const TRAININGS_URL = '/_app/trainings';

// The router plugin needs a string literal here, so it can't use TRAININGS_URL
export const Route = createFileRoute('/_app/trainings')({
  component: TrainingsPage,
  validateSearch: z.object({
    page: z.number().int().min(1).catch(1).default(1),
  }),
});

function TrainingsPage() {
  const { isAttendee } = useRole();

  return (
    <ContentLayout
      title="Your trainings"
      description="Upcoming sessions with your trainer"
      additionalContent={isAttendee ? <AdditionalContent /> : undefined}
    >
      <TrainingsList />
    </ContentLayout>
  );
}

const AdditionalContent = () => {
  const { userId } = useAuth();
  const { data } = useQuery(getUserOptions({ userId: userId! }));

  return (
    <div className="flex items-center gap-2">
      <CreditBadge credits={data?.balance ?? 0} />
      <Button
        variant="default"
        size="lg"
        nativeButton={false}
        render={<Link to="/trainings/schedule">+ Schedule training</Link>}
      />
    </div>
  );
};
