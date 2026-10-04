import { createFileRoute, Link, Outlet } from '@tanstack/react-router';

export const Route = createFileRoute('/_app')({
  component: AppLayout,
});

function AppLayout() {
  return (
    <>
      <nav>
        <Link to="/trainings">Trainings</Link>
        <Link to="/set-schedule">Set schedule</Link>
      </nav>
      <Outlet />
    </>
  );
}
