import { queryOptions } from '@tanstack/react-query';

import { http } from '@/api/client';
import type { TrainerDay } from '@/api/trainers/types';

type Params = { trainerUuid: string; fromDate: string; toDate: string };

export function getAvailableTrainerHours({
  trainerUuid,
  fromDate,
  toDate,
}: Params) {
  return http.get<TrainerDay[]>(`/api/v1/trainers/${trainerUuid}/hours`, {
    params: { from: fromDate, to: toDate, status: 'availability' },
  });
}

// Invalidate this key after a training is scheduled to refetch the picker.
export const AVAILABLE_TRAINER_HOURS_KEY = ['trainer-available-hours'] as const;

export const getAvailableTrainerHoursOptions = ({
  trainerUuid,
  fromDate,
  toDate,
}: Params) =>
  queryOptions({
    queryKey: [...AVAILABLE_TRAINER_HOURS_KEY, trainerUuid, fromDate, toDate],
    queryFn: () => getAvailableTrainerHours({ trainerUuid, fromDate, toDate }),
  });
