import { createFileRoute } from '@tanstack/react-router';
import * as z from 'zod';

import ContentLayout from '@/components/layout/content-layout';
import WeeksPicker from '@/features/set-schedule/components/weeks-picker';
import { WeeksProvider } from '@/features/set-schedule/components/weeks-provider';
import WeeksScheduleList from '@/features/set-schedule/components/weeks-schedule-list';

export const SET_SCHEDULE_URL = '/_app/set-schedule';

// The router plugin needs a string literal here, so it can't use SET_SCHEDULE_URL
export const Route = createFileRoute('/_app/set-schedule')({
  component: SetSchedulePage,
  validateSearch: z.object({
    fromDate: z.string().optional(),
    toDate: z.string().optional(),
  }),
});

function SetSchedulePage() {
  return (
    <WeeksProvider>
      <ContentLayout
        title="Set your availability"
        description="Tap an hour to open or close it -- changes save instantly"
        additionalContent={<WeeksPicker />}
      >
        <WeeksScheduleList />
      </ContentLayout>
    </WeeksProvider>
  );
}
