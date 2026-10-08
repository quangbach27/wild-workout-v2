import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';

import TrainingsList from '../trainings-list';

jest.mock('@/lib/env', () => ({ env: {} }));

jest.mock('@/routes/_app/trainings', () => ({
  TRAININGS_URL: '/_app/trainings',
}));

let mockPage = 1;
const mockNavigate = jest.fn();
jest.mock('@tanstack/react-router', () => ({
  Link: (props: { children: React.ReactNode }) => <a>{props.children}</a>,
  getRouteApi: () => ({
    useSearch: () => ({ page: mockPage }),
    useNavigate: () => mockNavigate,
  }),
}));

let mockItems: unknown[] = [];
let mockTotal: number | undefined;
jest.mock('../../api/get-trainings', () => ({
  getTrainingsOptions: (params: { page: number }) => ({
    queryKey: ['trainings', params],
    queryFn: async () => ({
      items: mockItems,
      pagination: {
        page: params.page,
        pageSize: 10,
        total: mockTotal ?? mockItems.length,
      },
    }),
  }),
}));

const mockCancel = jest.fn();
jest.mock('../../api/cancel-training', () => ({
  useCancelTraining: () => ({ mutate: mockCancel, isPending: false }),
}));

const training = (uuid: string, canceled = false) => ({
  uuid,
  hour: '2030-01-02T14:00:00.000Z',
  notes: '',
  attendeeUuid: 'a',
  attendeeUsername: 'a',
  trainerUuid: 't',
  trainerUsername: 't',
  canceled,
});

function renderList() {
  render(
    <QueryClientProvider client={new QueryClient()}>
      <TrainingsList />
    </QueryClientProvider>,
  );
}

describe('TrainingsList', () => {
  beforeEach(() => {
    mockPage = 1;
    mockTotal = undefined;
    mockNavigate.mockClear();
    mockCancel.mockClear();
    window.scrollTo = jest.fn();
  });

  it('renders a card with actions per upcoming training', async () => {
    mockItems = [training('1'), training('2')];
    renderList();

    expect(await screen.findAllByText('Move')).toHaveLength(2);
    expect(screen.getAllByText('Cancel')).toHaveLength(2);
  });

  it('cancels a training', async () => {
    mockItems = [training('1')];
    renderList();

    await userEvent.click(await screen.findByText('Cancel'));
    expect(mockCancel).toHaveBeenCalledWith('1');
  });

  it('hides the actions of canceled trainings', async () => {
    mockItems = [training('1', true)];
    renderList();

    expect(await screen.findByText('Canceled')).toBeInTheDocument();
    expect(screen.queryByText('Move')).not.toBeInTheDocument();
  });

  it('shows the empty state without trainings', async () => {
    mockItems = [];
    renderList();

    expect(
      await screen.findByText('NOTHING ON THE BAR YET'),
    ).toBeInTheDocument();
  });

  it('hides the pagination with a single page', async () => {
    mockItems = [training('1')];
    renderList();

    await screen.findByText('Move');
    expect(
      screen.queryByRole('button', { name: 'Next page' }),
    ).not.toBeInTheDocument();
  });

  it('goes to the next page', async () => {
    mockItems = [training('1')];
    mockTotal = 25;
    renderList();

    expect(await screen.findByText('Page 1 of 3')).toBeInTheDocument();
    expect(
      screen.getByRole('button', { name: 'Previous page' }),
    ).toBeDisabled();

    await userEvent.click(screen.getByRole('button', { name: 'Next page' }));
    expect(mockNavigate).toHaveBeenCalledWith({ search: { page: 2 } });
  });

  it('disables Next on the last page', async () => {
    mockPage = 3;
    mockItems = [training('1')];
    mockTotal = 25;
    renderList();

    expect(await screen.findByText('Page 3 of 3')).toBeInTheDocument();
    expect(screen.getByRole('button', { name: 'Next page' })).toBeDisabled();
  });

  it('moves to the last page when the current one is out of range', async () => {
    mockPage = 5;
    mockItems = [];
    mockTotal = 25;
    renderList();

    await screen.findByText('Page 5 of 3');
    expect(mockNavigate).toHaveBeenCalledWith({
      search: { page: 3 },
      replace: true,
    });
  });
});
