import { useMutation, useQueryClient } from '@tanstack/react-query';
import { toast } from 'sonner';

import { http } from '@/api/client';
import { AVAILABLE_TRAINER_HOURS_KEY } from '@/features/schedule-training/api/get-available-trainer-hours';

import { TRAININGS_KEY } from './get-trainings';

export function cancelTraining(trainingUuid: string) {
  return http.post<void>(`/api/v1/trainings/${trainingUuid}/cancel`);
}

export function useCancelTraining() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: cancelTraining,
    onSuccess: () => toast.success('Training canceled'),
    onError: (error) => toast.error(error.message),
    onSettled: () =>
      Promise.all([
        queryClient.invalidateQueries({ queryKey: TRAININGS_KEY }),
        queryClient.invalidateQueries({
          queryKey: AVAILABLE_TRAINER_HOURS_KEY,
        }),
      ]),
  });
}
