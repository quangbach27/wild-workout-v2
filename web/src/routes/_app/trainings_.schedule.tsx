import { createFileRoute } from '@tanstack/react-router';

import ContentLayout from '@/components/layout/content-layout';
import ScheduleTraining from '@/features/schedule-training/components/schedule-training';

export const Route = createFileRoute('/_app/trainings_/schedule')({
  component: SchedulePage,
});

function SchedulePage() {
  return (
    <ContentLayout
      title="Schedule a training"
      description="Pick any open hour in your trainer's week"
    >
      <ScheduleTraining />
    </ContentLayout>
  );
}
