import AppLayout from '@/components/layout/app-layout';

export default function Layout(props: LayoutProps<'/'>) {
  return <AppLayout>{props.children}</AppLayout>;
}
