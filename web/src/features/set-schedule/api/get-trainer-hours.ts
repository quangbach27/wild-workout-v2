import { queryOptions } from '@tanstack/react-query';

import { http } from '@/api/client';
import type { TrainerDay } from '@/api/trainers/types';

export function getTrainerHours(range: { fromDate: string; toDate: string }) {
  return http.get<TrainerDay[]>('/api/v1/trainer/hours', {
    params: { from: range.fromDate, to: range.toDate },
  });
}

// Invalidate ['trainer-hours'] after a mutation to refetch the current week.
export const getTrainerHoursOptions = ({
  fromDate,
  toDate,
}: {
  fromDate: string;
  toDate: string;
}) =>
  queryOptions({
    queryKey: ['trainer-hours', fromDate, toDate],
    queryFn: () => getTrainerHours({ fromDate, toDate }),
  });
