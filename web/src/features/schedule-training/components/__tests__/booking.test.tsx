import { Suspense, useEffect } from 'react';

import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';

import { BookingProvider, useBookingActions } from '../booking-provider';
import BookingSummary from '../booking-summary';
import TrainerDaySelect from '../trainer-day-select';
import TrainerHourSelect from '../trainer-hour-select';

jest.mock('@/lib/env', () => ({ env: { VITE_FUTURE_WEEKS: 2 } }));

const mockMutate = jest.fn();
jest.mock('../../api/schedule-training', () => ({
  useScheduleTraining: () => ({ mutate: mockMutate, isPending: false }),
}));

const HOUR = '2030-01-02T14:00:00.000Z';
jest.mock('../../api/get-available-trainer-hours', () => ({
  getAvailableTrainerHoursOptions: () => ({
    queryKey: ['hours'],
    queryFn: async () => [
      {
        date: '2030-01-02T00:00:00.000Z',
        hasFreeHours: true,
        hours: [{ hour: HOUR, status: 'availability' }],
      },
      { date: '2030-01-03T00:00:00.000Z', hasFreeHours: false, hours: [] },
    ],
  }),
}));

function SelectTrainer() {
  const { selectTrainer } = useBookingActions();
  useEffect(() => {
    selectTrainer({ trainerUuid: 'trainer-1', trainerName: 'Alex Johnson' });
  }, [selectTrainer]);
  return null;
}

describe('booking flow', () => {
  it('lists days, then hours, then schedules the training', async () => {
    const user = userEvent.setup();
    render(
      <QueryClientProvider client={new QueryClient()}>
        <BookingProvider>
          <SelectTrainer />
          <Suspense fallback="loading">
            <TrainerDaySelect trainerUuid="trainer-1" />
            <TrainerHourSelect trainerUuid="trainer-1" />
          </Suspense>
          <BookingSummary />
        </BookingProvider>
      </QueryClientProvider>,
    );

    // Only the day with free hours is offered, and no hours until it is picked
    const toggles = () => screen.getAllByRole('button', { pressed: false });
    await screen.findByText('Pick a day first.');
    expect(toggles()).toHaveLength(1);
    const confirm = screen.getByRole('button', { name: 'Confirm training' });
    expect(confirm).toBeDisabled();

    await user.click(toggles()[0]);
    await user.click(toggles()[0]);
    expect(confirm).toBeEnabled();

    await user.click(confirm);
    expect(mockMutate).toHaveBeenCalledWith(
      {
        trainerUuid: 'trainer-1',
        trainerUsername: 'Alex Johnson',
        hour: HOUR,
        notes: '',
      },
      expect.any(Object),
    );
  });
});
