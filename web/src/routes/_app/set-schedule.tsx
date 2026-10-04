import { createFileRoute } from '@tanstack/react-router';

import ContentLayout from '@/components/layout/content-layout';

export const Route = createFileRoute('/_app/set-schedule')({
  component: SetSchedulePage,
});

function SetSchedulePage() {
  return (
    <ContentLayout
      label="Set your availability"
      description="Tap an hour to open or close it -- changes save instantly"
      additionalContent="[][][]"
    >
      Schedule page
    </ContentLayout>
  );
}
