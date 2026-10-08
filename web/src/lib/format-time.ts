const timeFormat = new Intl.DateTimeFormat('en-GB', {
  hour: '2-digit',
  minute: '2-digit',
  hour12: false,
});

export const formatTime = (iso: string | Date) =>
  timeFormat.format(typeof iso === 'string' ? new Date(iso) : iso);

// Trainings last one hour.
export const formatTimeRange = (iso: string) => {
  const start = new Date(iso);
  return `${formatTime(start)} – ${formatTime(new Date(start.getTime() + 3_600_000))}`;
};
