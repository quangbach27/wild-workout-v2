import type React from 'react';

export default function ContentLayout(props: {
  title: string;
  description: string;
  additionalContent?: React.ReactNode;
  children: React.ReactNode;
}) {
  return (
    <div className="flex min-h-0 flex-1 flex-col pt-5 md:pt-10">
      <div className="flex shrink-0 flex-col gap-2 pb-3 md:flex-row md:items-end md:justify-between md:pb-5">
        <div>
          <h1 className="text-2xl font-bold uppercase md:text-4xl">
            {props.title}
          </h1>
          <span className="text-muted-foreground">{props.description}</span>
        </div>

        {props.additionalContent && <div>{props.additionalContent}</div>}
      </div>

      <div className="min-h-0 flex-1 overflow-auto">{props.children}</div>
    </div>
  );
}
