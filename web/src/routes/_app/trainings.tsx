import { createFileRoute } from '@tanstack/react-router';

import ContentLayout from '@/components/layout/content-layout';
import TrainingsList from '@/features/trainings/components/trainings-list';

export const Route = createFileRoute('/_app/trainings')({
  component: TrainingsPage,
});

function TrainingsPage() {
  return (
    <ContentLayout
      label="Your trainings"
      description="Upcomming sessions with your attendees"
    >
      <TrainingsList />
    </ContentLayout>
  );
}
