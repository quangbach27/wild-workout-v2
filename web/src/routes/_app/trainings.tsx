import { createFileRoute, Link } from '@tanstack/react-router';
import * as z from 'zod';

import { CreditBadge } from '@/components/credit-badge';
import ContentLayout from '@/components/layout/content-layout';
import { Button } from '@/components/ui/button';
import TrainingsList from '@/features/trainings/components/trainings-list';

export const TRAININGS_URL = '/_app/trainings';

// The router plugin needs a string literal here, so it can't use TRAININGS_URL
export const Route = createFileRoute('/_app/trainings')({
  component: TrainingsPage,
  validateSearch: z.object({
    page: z.number().int().min(1).catch(1).default(1),
  }),
});

function TrainingsPage() {
  return (
    <ContentLayout
      title="Your trainings"
      description="Upcoming sessions with your trainer"
      additionalContent={<AdditionalContent />}
    >
      <TrainingsList />
    </ContentLayout>
  );
}

const AdditionalContent = () => {
  return (
    <div className="flex items-center gap-2">
      <CreditBadge credits={41} />
      <Button
        variant="default"
        size="lg"
        render={<Link to="/trainings/schedule">+ Schedule training</Link>}
      />
    </div>
  );
};
