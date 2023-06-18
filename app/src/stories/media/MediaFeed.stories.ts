import { Meta, StoryObj } from '@storybook/vue3';

import recommendedMedia from '../assets/fixtures/recommended-media.json';
import MediaFeed from '../../components/media/MediaFeed.vue';
import { RecommendedMedia } from '../../model/Media';

const meta = {
  title: 'media/MediaFeed',
  component: MediaFeed,
  parameters: {
    layout: 'fullscreen',
    viewport: {
      defaultViewport: 'pixel',
    },
  },
} satisfies Meta<typeof MediaFeed>;

export default meta;

type Story = StoryObj<typeof meta>;

export const Default: Story = {
  args: {
    media: recommendedMedia as unknown as RecommendedMedia[],
  },
};
