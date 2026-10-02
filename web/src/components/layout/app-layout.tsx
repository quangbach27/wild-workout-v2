import React from 'react';

import Footer from './footer';
import Header from './header';

export default function AppLayout(props: { children: React.ReactNode }) {
  return (
    <div className="flex h-dvh flex-col">
      <Header />
      <main className="flex-1 overflow-y-auto">{props.children}</main>
      <Footer />
    </div>
  );
}
