import BookingForm from './booking-form';
import { BookingProvider } from './booking-provider';
import BookingSummary from './booking-summary';

export default function ScheduleTraining() {
  return (
    <BookingProvider>
      <div className="flex flex-col gap-2 md:flex-row md:gap-4">
        <BookingForm className="flex-1" />
        <BookingSummary />
      </div>
    </BookingProvider>
  );
}
