import { useMutation, useQueryClient } from '@tanstack/react-query';
import { toast } from 'sonner';

import { http } from '@/api/client';
import type { HoursRequest } from '@/api/trainers/types';

import { TRAINER_HOURS_KEY } from './get-trainer-hours';

export function makeHoursAvailable(hours: string[]) {
  return http.put<void, HoursRequest>('/api/v1/trainer/hours/available', {
    hours,
  });
}

export function useMakeHoursAvailable() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: makeHoursAvailable,
    onSuccess: (_, hours) =>
      toast.success(
        hours.length === 1
          ? 'Hour is now available'
          : `${hours.length} hours are now available`,
      ),
    onError: (error) => toast.error(error.message),
    onSettled: () =>
      queryClient.invalidateQueries({ queryKey: TRAINER_HOURS_KEY }),
  });
}
