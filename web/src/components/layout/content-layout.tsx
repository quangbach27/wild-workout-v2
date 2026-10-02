import React from 'react';

import { cn } from 'cn';

export default function ContentLayout(props: {
  title: string;
  description: string;
  className?: string;
  additionalContent?: React.ReactNode;
  children: React.ReactNode;
}) {
  return (
    <div className={cn('container', props.className)}>
      <div className="flex flex-col justify-between gap-4 pt-6 pb-3 md:flex-row md:items-end md:pt-10 md:pb-6">
        <div>
          <div className="text-4xl font-bold uppercase">{props.title}</div>
          <div>{props.description}</div>
        </div>
        {!!props.additionalContent && <div>{props.additionalContent}</div>}
      </div>
      <div className="text-muted-foreground flex flex-wrap items-center gap-x-6 gap-y-2 text-sm">
        {props.children}
      </div>
    </div>
  );
}
