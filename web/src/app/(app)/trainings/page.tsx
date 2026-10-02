import ContentLayout from '@/components/layout/content-layout';

export default async function Page() {
  return (
    <ContentLayout
      className="px-2"
      title="Your Trainings"
      description="Upcoming session with your attendees"
    >
      <div>Content</div>
    </ContentLayout>
  );
}
