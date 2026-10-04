import { createFileRoute } from '@tanstack/react-router';

import { Button } from '@/components/ui/button';

export const Route = createFileRoute('/_app/set-schedule')({
  component: SetSchedulePage,
});

function SetSchedulePage() {
  return <Button>Hello</Button>;
}
