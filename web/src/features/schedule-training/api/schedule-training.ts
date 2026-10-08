import { useMutation, useQueryClient } from '@tanstack/react-query';
import { toast } from 'sonner';

import { http } from '@/api/client';
import type {
  ScheduledTraining,
  ScheduleTrainingRequest,
} from '@/api/trainings/types';
import { TRAININGS_KEY } from '@/features/trainings/api/get-trainings';

import { AVAILABLE_TRAINER_HOURS_KEY } from './get-available-trainer-hours';

export function scheduleTraining(request: ScheduleTrainingRequest) {
  return http.post<ScheduledTraining, ScheduleTrainingRequest>(
    '/api/v1/trainings',
    request,
  );
}

export function useScheduleTraining() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: scheduleTraining,
    onSuccess: () => toast.success('Training scheduled'),
    onError: (error) => toast.error(error.message),
    onSettled: () =>
      Promise.all([
        queryClient.invalidateQueries({
          queryKey: AVAILABLE_TRAINER_HOURS_KEY,
        }),
        queryClient.invalidateQueries({ queryKey: TRAININGS_KEY }),
      ]),
  });
}
