import { render, screen } from '@testing-library/react';

import DateBadge from '../date-badge';

describe('DateBadge', () => {
  it('renders weekday, day and month', () => {
    render(<DateBadge weekday="Mon" day="12" month="Oct" />);

    expect(screen.getByText('Mon')).toBeInTheDocument();
    expect(screen.getByText('12')).toBeInTheDocument();
    expect(screen.getByText('Oct')).toBeInTheDocument();
  });

  it('shows the Today pill only when isToday in the transparent variant', () => {
    const { rerender } = render(
      <DateBadge variant="transparent" weekday="Mon" day="12" isToday />,
    );
    expect(screen.getByText('Today')).toBeInTheDocument();

    rerender(<DateBadge variant="transparent" weekday="Mon" day="12" />);
    expect(screen.queryByText('Today')).not.toBeInTheDocument();
  });

  it('sets data attributes for variant and today', () => {
    const { container } = render(
      <DateBadge variant="outline" weekday="Mon" day="12" isToday />,
    );
    const badge = container.querySelector('[data-slot="date-badge"]');

    expect(badge).toHaveAttribute('data-variant', 'outline');
    expect(badge).toHaveAttribute('data-today', 'true');
  });
});
