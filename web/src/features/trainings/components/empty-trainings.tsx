import { Link } from '@tanstack/react-router';

import { Button } from '@/components/ui/button';
import {
  Empty,
  EmptyContent,
  EmptyDescription,
  EmptyHeader,
  EmptyMedia,
  EmptyTitle,
} from '@/components/ui/empty';
import { useRole } from '@/lib/role';

import BarbellIcon from './barbell-icon';

export default function EmptyTrainings() {
  const { isTrainer } = useRole();

  return (
    <Empty className="border border-dashed bg-white">
      <EmptyHeader>
        <EmptyMedia>
          <BarbellIcon />
        </EmptyMedia>
        <EmptyTitle>NOTHING ON THE BAR YET</EmptyTitle>
        <EmptyDescription>
          {isTrainer
            ? 'No sessions booked so far. Open some hours so attendees can book you'
            : 'No sessions booked yet. Find a trainer and book your first one'}
        </EmptyDescription>
      </EmptyHeader>
      <EmptyContent>
        {isTrainer ? (
          <Button
            variant="outline"
            render={<Link to="/set-schedule">Set your availability</Link>}
          />
        ) : (
          <Button
            variant="outline"
            render={<Link to="/trainings/schedule">Schedule training</Link>}
          />
        )}
      </EmptyContent>
    </Empty>
  );
}
