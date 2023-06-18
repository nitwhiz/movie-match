import { Meta, StoryObj } from '@storybook/vue3';
import FeedView from '../../views/FeedView.vue';

const meta = {
  title: 'views/Feed',
  component: FeedView,
  parameters: {
    layout: 'fullscreen',
    viewport: {
      defaultViewport: 'pixel',
    },
  },
} satisfies Meta<typeof FeedView>;

export default meta;

type Story = StoryObj<typeof meta>;

export const Default: Story = {};
