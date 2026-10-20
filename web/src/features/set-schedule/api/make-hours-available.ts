import {
  useMutation,
  useQueryClient,
  type UseMutationOptions,
} from '@tanstack/react-query';
import { toast } from 'sonner';

import { http } from '@/api/client';
import type { HoursRequest } from '@/api/trainers/types';

import { TRAINER_HOURS_KEY } from './get-trainer-hours';

export function makeHoursAvailable(hours: string[]) {
  return http.put<void, HoursRequest>('/api/v1/trainer/hours/available', {
    hours,
  });
}

type MakeHoursAvailableOptions = Pick<
  UseMutationOptions<void, Error, string[]>,
  'onSuccess' | 'onError' | 'onSettled'
>;

export function useMakeHoursAvailable(options?: MakeHoursAvailableOptions) {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: makeHoursAvailable,
    onSuccess: (data, hours, onMutateResult, context) => {
      toast.success(
        hours.length === 1
          ? 'Hour is now available'
          : `${hours.length} hours are now available`,
      );
      return options?.onSuccess?.(data, hours, onMutateResult, context);
    },
    onError: (error, hours, onMutateResult, context) => {
      toast.error(error.message);
      return options?.onError?.(error, hours, onMutateResult, context);
    },
    onSettled: (data, error, hours, onMutateResult, context) => {
      queryClient.invalidateQueries({ queryKey: TRAINER_HOURS_KEY });
      return options?.onSettled?.(data, error, hours, onMutateResult, context);
    },
  });
}
