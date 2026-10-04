import { createFileRoute } from '@tanstack/react-router';

export const Route = createFileRoute('/_app/trainings')({
  component: TrainingsPage,
});

function TrainingsPage() {
  return <div>Trainings</div>;
}
