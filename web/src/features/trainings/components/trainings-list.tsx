import EmptyTrainings from './empty-trainings';

export default function TrainingsList() {
  const isTrainingsEmpty = true;

  if (isTrainingsEmpty) {
    return <EmptyTrainings />;
  }

  return <div>Training List</div>;
}
