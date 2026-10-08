import {
  createContext,
  useCallback,
  useContext,
  useEffect,
  useMemo,
  useState,
  startTransition,
} from 'react';

import { getRouteApi } from '@tanstack/react-router';

import { buildWeeks, type Week } from '@/features/set-schedule/lib/weeks';
import { SET_SCHEDULE_URL } from '@/routes/_app/set-schedule';

const route = getRouteApi(SET_SCHEDULE_URL);

type WeeksState = {
  weeks: Week[];
  activeWeek: Week;
};

type WeeksActions = {
  selectWeek: (label: string | null) => void;
};

// State and actions live in separate contexts: consumers that only need the
// (stable) actions don't re-render when the active week changes.
const WeeksStateContext = createContext<WeeksState | null>(null);
const WeeksActionsContext = createContext<WeeksActions | null>(null);

export function WeeksProvider(props: { children: React.ReactNode }) {
  const weeks = useMemo(() => buildWeeks(), []);
  const navigate = route.useNavigate();

  // The URL is only read on the initial load, to pick the first selected week.
  // After that it is only written to.
  const search = route.useSearch();
  const [initial] = useState(() => ({
    week: weeks.find(
      (w) =>
        w.value.fromDate === search.fromDate &&
        w.value.toDate === search.toDate,
    ),
    hasSearch: search.fromDate !== undefined || search.toDate !== undefined,
  }));

  const [activeWeek, setActiveWeek] = useState(() => initial.week ?? weeks[0]);

  // Search params that don't match any week are stale: drop them.
  useEffect(() => {
    if (!initial.hasSearch || initial.week) return;

    navigate({
      search: (prev) => ({ ...prev, fromDate: undefined, toDate: undefined }),
      replace: true,
      resetScroll: false,
    });
  }, [initial, navigate]);

  const selectWeek = useCallback(
    (label: string | null) => {
      const selected = weeks.find((w) => w.label === label);
      if (!selected) return;

      // Keep showing the current grid while the new week loads
      startTransition(() => setActiveWeek(selected));
      navigate({
        search: (prev) => ({
          ...prev,
          fromDate: selected.value.fromDate,
          toDate: selected.value.toDate,
        }),
        replace: true,
        resetScroll: false,
      });
    },
    [weeks, navigate],
  );

  const state = useMemo(() => ({ weeks, activeWeek }), [weeks, activeWeek]);
  const actions = useMemo(() => ({ selectWeek }), [selectWeek]);

  return (
    <WeeksActionsContext value={actions}>
      <WeeksStateContext value={state}>{props.children}</WeeksStateContext>
    </WeeksActionsContext>
  );
}

export function useWeeksState() {
  const ctx = useContext(WeeksStateContext);
  if (!ctx) throw new Error('useWeeksState must be used within WeeksProvider');
  return ctx;
}

export function useWeeksActions() {
  const ctx = useContext(WeeksActionsContext);
  if (!ctx) {
    throw new Error('useWeeksActions must be used within WeeksProvider');
  }
  return ctx;
}
