import { queryOptions } from '@tanstack/react-query';

import { http } from '@/api/client';
import type { TrainersPage } from '@/api/users/types';

type Params = { page?: number; pageSize?: number };

export function getTrainers(params: Params = {}) {
  return http.get<TrainersPage>('/api/v1/trainers', { params });
}

export const TRAINERS_KEY = ['trainers'] as const;

export const getTrainersOptions = (params: Params = {}) =>
  queryOptions({
    queryKey: [...TRAINERS_KEY, params],
    queryFn: () => getTrainers(params),
  });
