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

import BarbellIcon from './barbell-icon';

export default function EmptyTrainings() {
  return (
    <Empty className="border border-dashed bg-white">
      <EmptyHeader>
        <EmptyMedia>
          <BarbellIcon />
        </EmptyMedia>
        <EmptyTitle>NOTHING ON THE BAR YET</EmptyTitle>
        <EmptyDescription>
          No sessions booked so for. Open some hours so attendees can book you
        </EmptyDescription>
      </EmptyHeader>
      <EmptyContent>
        <Button
          variant="outline"
          render={<Link to="/set-schedule">Set you availability</Link>}
        />
      </EmptyContent>
    </Empty>
  );
}
