import { useState } from 'react';

import { useSuspenseQuery } from '@tanstack/react-query';

import { buildWeeks } from '@/features/set-schedule/lib/weeks';

import { getAvailableTrainerHoursOptions } from '../api/get-available-trainer-hours';

// Days with open hours for a trainer across the whole bookable window. The day
// and hour steps both call this: the query key is the same, so the request is
// shared through the query cache.
export function useTrainerAvailability(trainerUuid: string) {
  // From the first week's start to the last week's end
  const [range] = useState(() => {
    const weeks = buildWeeks();
    return {
      fromDate: weeks[0].value.fromDate,
      toDate: weeks[weeks.length - 1].value.toDate,
    };
  });

  const { data } = useSuspenseQuery(
    getAvailableTrainerHoursOptions({ trainerUuid, ...range }),
  );

  return data.filter((d) => d.hasFreeHours);
}
