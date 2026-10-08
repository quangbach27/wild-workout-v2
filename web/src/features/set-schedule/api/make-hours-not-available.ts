import { useMutation, useQueryClient } from '@tanstack/react-query';
import { toast } from 'sonner';

import { http } from '@/api/client';
import type { HoursRequest } from '@/api/trainers/types';

import { TRAINER_HOURS_KEY } from './get-trainer-hours';

export function makeHoursNotAvailable(hours: string[]) {
  return http.put<void, HoursRequest>('/api/v1/trainer/hours/not-available', {
    hours,
  });
}

export function useMakeHoursNotAvailable() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: makeHoursNotAvailable,
    onSuccess: () => toast.success('Hour is now not available'),
    onError: (error) => toast.error(error.message),
    onSettled: () =>
      queryClient.invalidateQueries({ queryKey: TRAINER_HOURS_KEY }),
  });
}
