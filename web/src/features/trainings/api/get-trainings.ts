import { queryOptions } from '@tanstack/react-query';

import { http } from '@/api/client';
import type { TrainingsPage } from '@/api/trainings/types';

type Params = { page?: number; pageSize?: number };

export function getTrainings(params: Params = {}) {
  return http.get<TrainingsPage>('/api/v1/trainings', { params });
}

// Invalidate this key after a training is scheduled or canceled.
export const TRAININGS_KEY = ['trainings'] as const;

export const getTrainingsOptions = (params: Params = {}) =>
  queryOptions({
    queryKey: [...TRAININGS_KEY, params],
    queryFn: () => getTrainings(params),
  });
