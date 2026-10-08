import { CreditBadge } from './index';

import type { Meta, StoryObj } from '@storybook/react-vite';

const meta = {
  title: 'Components/CreditBadge',
  component: CreditBadge,
  tags: ['autodocs'],
  args: {
    credits: 5,
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
    credits: { control: 'number' },
    label: { control: 'text' },
  },
} satisfies Meta<typeof CreditBadge>;

export default meta;
type Story = StoryObj<typeof meta>;

export const Default: Story = {};

export const CustomLabel: Story = {
  args: { credits: 12, label: 'Credits' },
};
