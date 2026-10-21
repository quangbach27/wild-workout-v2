import { useEffect, useState } from 'react';

const HOUR_MS = 3_600_000;

const startOfHour = (t = Date.now()) => Math.floor(t / HOUR_MS) * HOUR_MS;

/** Returns the current hour's start as a timestamp. Changes only when the hour changes. */
export function useCurrentHour(): number {
  const [hour, setHour] = useState(() => startOfHour());

  useEffect(() => {
    const msUntilNextHour = hour + HOUR_MS - Date.now();
    const id = setTimeout(() => setHour(startOfHour()), msUntilNextHour);
    return () => clearTimeout(id);
  }, [hour]);

  return hour;
}
