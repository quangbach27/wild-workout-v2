import ContentLayout from '@/components/layout/content-layout';

export default async function Page() {
  return (
    <ContentLayout
      title="Set you availability"
      description="Tap an hour to open or close it - changes save instantly"
      className="px-2"
    >
      <div className="text-muted-foreground flex flex-wrap items-center gap-x-6 gap-y-2 text-sm">
        <div className="flex items-center gap-2">
          <span className="bg-hour-closed border-hour-closed-border size-4 rounded-[5px] border" />
          <span>Closed</span>
        </div>

        <div className="flex items-center gap-2">
          <span className="bg-hour-open border-hour-open-border size-4 rounded-[5px] border" />
          <span>Open for booking</span>
        </div>

        <div className="flex items-center gap-2">
          <span className="bg-hour-booked border-hour-booked size-4 rounded-[5px] border" />
          <span>Training booked - locked</span>
        </div>
      </div>
    </ContentLayout>
  );
}
