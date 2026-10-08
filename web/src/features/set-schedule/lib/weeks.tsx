import { env } from '@/lib/env';

const DAYS_PER_WEEK = 7;

export type Week = {
  label: string;
  tag?: string;
  value: {
    fromDate: string;
    toDate: string;
  };
};

type buildWeeksConfig = {
  futureWeeks?: number;
};

export function buildWeeks(config: buildWeeksConfig = {}): Week[] {
  const { futureWeeks = env.VITE_FUTURE_WEEKS } = config;

  // Earliest bookable time: next full hour after now + 1h
  const startTime = new Date();
  startTime.setHours(startTime.getHours() + 2, 0, 0, 0);

  const year = startTime.getFullYear();
  const month = startTime.getMonth();
  const day = startTime.getDate();

  return Array.from({ length: futureWeeks }, (_, i) => {
    const offset = i * DAYS_PER_WEEK;

    const from = i === 0 ? startTime : new Date(year, month, day + offset);
    const lastDay = new Date(year, month, day + offset + DAYS_PER_WEEK - 1);
    const to = new Date(year, month, day + offset + DAYS_PER_WEEK);

    const sameMonth =
      from.getMonth() === lastDay.getMonth() &&
      from.getFullYear() === lastDay.getFullYear();

    return {
      label: `${monthDay.format(from)} – ${
        sameMonth ? dayOnly.format(lastDay) : monthDay.format(lastDay)
      }`,
      tag: i === 0 ? 'Current week' : undefined,
      value: {
        fromDate: from.toISOString(),
        toDate: to.toISOString(),
      },
    };
  });
}

const monthDay = new Intl.DateTimeFormat('en-US', {
  month: 'short',
  day: 'numeric',
});

const dayOnly = new Intl.DateTimeFormat('en-US', {
  day: 'numeric',
});
