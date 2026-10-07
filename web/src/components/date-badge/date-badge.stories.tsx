import { DateBadge } from './index';

import type { Meta, StoryObj } from '@storybook/react-vite';

const meta = {
  title: 'Components/DateBadge',
  component: DateBadge,
  tags: ['autodocs'],
  args: {
    weekday: 'Mon',
    day: '12',
    month: 'Oct',
  },
  // Canvas follows the toolbar theme (var(--background)); the rest of the page is untouched
  parameters: {
    backgrounds: {
      options: { theme: { name: 'theme', value: 'var(--background)' } },
    },
  },
  globals: { backgrounds: { value: 'theme' } },
  decorators: [
    (Story) => (
      <div className="text-foreground">
        <Story />
      </div>
    ),
  ],
  argTypes: {
    variant: { control: 'inline-radio', options: ['outline', 'transparent'] },
    size: { control: 'inline-radio', options: ['sm', 'md', 'lg'] },
    isToday: { control: 'boolean' },
  },
} satisfies Meta<typeof DateBadge>;

export default meta;
type Story = StoryObj<typeof meta>;

export const Outline: Story = {
  args: { variant: 'outline' },
};

export const Transparent: Story = {
  args: { variant: 'transparent', month: undefined },
  decorators: [
    (Story) => (
      <div className="w-24">
        <Story />
      </div>
    ),
  ],
};

export const Today: Story = {
  args: { variant: 'transparent', isToday: true, month: undefined },
  decorators: Transparent.decorators,
};

export const Sizes: Story = {
  args: { variant: 'outline' },
  render: (args) => (
    <div className="flex items-end gap-4">
      <DateBadge {...args} size="sm" />
      <DateBadge {...args} size="md" />
      <DateBadge {...args} size="lg" />
    </div>
  ),
};
