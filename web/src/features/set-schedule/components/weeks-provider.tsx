import {
  createContext,
  useCallback,
  useContext,
  useEffect,
  useMemo,
  useState,
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
  selectWeekRange: (label: string | null) => void;
};

// State and actions live in separate contexts: consumers that only need the
// (stable) actions don't re-render when the active week changes.
const WeeksStateContext = createContext<WeeksState | null>(null);
const WeeksActionsContext = createContext<WeeksActions | null>(null);

export function WeeksProvider(props: { children: React.ReactNode }) {
  const navigate = route.useNavigate();
  const { fromDate, toDate } = route.useSearch();

  const weeks = useMemo(() => buildWeeks(), []);

  const [initial] = useState(() => {
    if (!fromDate || !toDate) {
      return null;
    }

    return {
      week: weeks.find(
        (w) => w.value.fromDate === fromDate && w.value.toDate === toDate,
      ),
    };
  });

  const [activeWeek, setActiveWeek] = useState(() => initial?.week ?? weeks[0]);

  // On first load, check the fromDate and toDate search params. If they don't match any week, remove them.
  useEffect(() => {
    if (fromDate && toDate && initial?.week) return;

    navigate({
      search: (prev) => ({ ...prev, fromDate: undefined, toDate: undefined }),
      replace: true,
      resetScroll: false,
    });
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [navigate]);

  const selectWeekRange = useCallback(
    (label: string | null) => {
      const selected = weeks.find((w) => w.label === label);
      if (!selected) return;

      setActiveWeek(selected);
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
  const actions = useMemo(() => ({ selectWeekRange }), [selectWeekRange]);

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
