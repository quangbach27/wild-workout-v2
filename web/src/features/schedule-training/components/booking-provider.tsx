import { createContext, useContext, useMemo, useReducer } from 'react';

type BookingState = {
  trainerUuid: string | null;
  trainerName: string | null;
  day: string | null;
  hour: string | null;
  notes: string | null;
};

type BookingActions = {
  selectTrainer: (trainer: {
    trainerUuid: string;
    trainerName: string;
  }) => void;
  selectDay: (day: string) => void;
  selectHour: (hour: string) => void;
  setNotes: (notes: string) => void;
  reset: () => void;
};

const INITIAL_BOOKING: BookingState = {
  trainerUuid: null,
  trainerName: null,
  day: null,
  hour: null,
  notes: null,
};

// State and actions live in separate contexts: consumers that only need the
// (stable) actions don't re-render when the draft booking changes.
const BookingStateContext = createContext<BookingState | null>(null);
const BookingActionsContext = createContext<BookingActions | null>(null);

type BookingAction =
  | {
      type: 'selectTrainer';
      payload: { trainerUuid: string; trainerName: string };
    }
  | { type: 'selectDay'; payload: { day: string } }
  | { type: 'selectHour'; payload: { hour: string } }
  | { type: 'setNotes'; payload: { notes: string } }
  | { type: 'reset' };

function bookingReducer(
  state: BookingState,
  action: BookingAction,
): BookingState {
  switch (action.type) {
    case 'selectTrainer': {
      const { trainerUuid, trainerName } = action.payload;
      return state.trainerUuid === trainerUuid
        ? state
        : { ...state, trainerUuid, trainerName, day: null, hour: null };
    }
    case 'selectDay': {
      return state.day === action.payload.day
        ? state
        : { ...state, day: action.payload.day, hour: null };
    }
    case 'selectHour':
      return { ...state, hour: action.payload.hour };
    case 'setNotes':
      return { ...state, notes: action.payload.notes };
    case 'reset':
      return INITIAL_BOOKING;
  }
}

// Holds the booking the user is composing, before it is confirmed.
export function BookingProvider(props: { children: React.ReactNode }) {
  const [booking, dispatch] = useReducer(bookingReducer, INITIAL_BOOKING);

  const actions = useMemo<BookingActions>(
    () => ({
      selectTrainer: ({ trainerUuid, trainerName }) =>
        dispatch({
          type: 'selectTrainer',
          payload: { trainerUuid, trainerName },
        }),
      selectDay: (day) => dispatch({ type: 'selectDay', payload: { day } }),
      selectHour: (hour) => dispatch({ type: 'selectHour', payload: { hour } }),
      setNotes: (notes) => dispatch({ type: 'setNotes', payload: { notes } }),
      reset: () => dispatch({ type: 'reset' }),
    }),
    [],
  );

  return (
    <BookingActionsContext value={actions}>
      <BookingStateContext value={booking}>
        {props.children}
      </BookingStateContext>
    </BookingActionsContext>
  );
}

export function useBookingState() {
  const ctx = useContext(BookingStateContext);
  if (!ctx) {
    throw new Error('useBookingState must be used within BookingProvider');
  }
  return ctx;
}

export function useBookingActions() {
  const ctx = useContext(BookingActionsContext);
  if (!ctx) {
    throw new Error('useBookingActions must be used within BookingProvider');
  }
  return ctx;
}
