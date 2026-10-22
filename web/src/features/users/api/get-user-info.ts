import { queryOptions } from '@tanstack/react-query';

import { http } from '@/api/client';
import type { User } from '@/api/users/types';

type Params = {
  userId: string;
};
const getUserInfo = ({ userId: userID }: Params) => {
  return http.get<User>(`/api/v1/users/${userID}`);
};

export const USER_KEY = ['user'] as const;

export const getUserOptions = ({ userId }: Params) =>
  queryOptions({
    queryKey: [...USER_KEY, userId],
    queryFn: () => getUserInfo({ userId }),
  });
