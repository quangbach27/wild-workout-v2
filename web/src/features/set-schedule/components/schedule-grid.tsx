import { useSuspenseQuery } from '@tanstack/react-query';

import { getTrainerHoursOptions } from '../api/get-trainer-hours';

import DayColumn from './day-column';
import { useWeeksState } from './weeks-provider';

export default function ScheduleGrid() {
  const { activeWeek } = useWeeksState();
  const { data } = useSuspenseQuery(getTrainerHoursOptions(activeWeek.value));

  return (
    <div
      className="grid min-w-4xl gap-3 pb-6"
      style={{
        gridTemplateColumns: `repeat(${data.length}, minmax(0, 1fr))`,
      }}
    >
      {data.map((day) => (
        <DayColumn key={day.date} day={day} />
      ))}
    </div>
  );
}
