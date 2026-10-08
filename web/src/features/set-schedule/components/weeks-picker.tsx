import { Fragment } from 'react';

import { CalendarDaysIcon } from 'lucide-react';

import {
  Select,
  SelectContent,
  SelectGroup,
  SelectItem,
  SelectSeparator,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select';

import { useWeeksActions, useWeeksState } from './weeks-provider';

function ThisWeekTag(props: { children: string }) {
  return (
    <span className="bg-hour-action/10 text-hour-action rounded-full px-2 py-0.5 text-xs font-semibold">
      {props.children}
    </span>
  );
}

export default function WeeksPicker() {
  const { weeks, activeWeek } = useWeeksState();
  const { selectWeek } = useWeeksActions();

  return (
    <Select
      items={weeks.map((w) => ({ value: w.label, label: w.label }))}
      value={activeWeek.label}
      onValueChange={selectWeek}
    >
      <SelectTrigger className="h-10 w-full min-w-52 bg-white md:w-fit">
        <CalendarDaysIcon className="text-muted-foreground" />
        <SelectValue>{activeWeek.label}</SelectValue>
        {activeWeek.tag && <ThisWeekTag>{activeWeek.tag}</ThisWeekTag>}
      </SelectTrigger>
      <SelectContent
        alignItemWithTrigger={false}
        align="end"
        sideOffset={6}
        className="min-w-56"
      >
        <SelectGroup>
          {weeks.map((item) => (
            <Fragment key={item.label}>
              <SelectItem value={item.label} className="py-2 pr-10">
                {item.label}
                {item.tag && <ThisWeekTag>{item.tag}</ThisWeekTag>}
              </SelectItem>
              {item.tag && <SelectSeparator />}
            </Fragment>
          ))}
        </SelectGroup>
      </SelectContent>
    </Select>
  );
}
