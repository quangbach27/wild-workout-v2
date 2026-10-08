import { useSuspenseQuery } from '@tanstack/react-query';

import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select';

import { getTrainersOptions } from '../api/get-trainers';

import { useBookingActions, useBookingState } from './booking-provider';

export default function TrainerSelect() {
  const { trainerUuid } = useBookingState();
  const { selectTrainer } = useBookingActions();
  const { data } = useSuspenseQuery(getTrainersOptions());

  const items = data.items.map((t) => ({
    value: t.uuid,
    label: t.displayName,
  }));

  return (
    <Select
      value={trainerUuid}
      onValueChange={(value) => {
        const trainer = items.find((t) => t.value === value);
        if (trainer) {
          selectTrainer({
            trainerUuid: trainer.value,
            trainerName: trainer.label,
          });
        }
      }}
      items={items}
    >
      <SelectTrigger className="h-10 w-full sm:w-72">
        <SelectValue placeholder="Select a trainer" />
      </SelectTrigger>
      <SelectContent>
        {items.map((t) => (
          <SelectItem key={t.value} value={t.value}>
            {t.label}
          </SelectItem>
        ))}
      </SelectContent>
    </Select>
  );
}
