import { Meta, StoryObj } from '@storybook/vue3';
import MediaCard from '../../components/media/MediaCard.vue';

import recommendedMedia from '../assets/fixtures/recommended-media.json';
import { RecommendedMedia } from '../../model/Media';

const meta = {
  title: 'media/MediaCard',
  component: MediaCard,
  parameters: {
    layout: 'fullscreen',
    viewport: {
      defaultViewport: 'pixel',
    },
  },
  argTypes: {
    media: {
      options: recommendedMedia.map((m) => m.title),
      mapping: recommendedMedia.reduce(
        (a, m) => ({
          ...a,
          [m.title]: m,
        }),
        {}
      ),
    },
  },
} satisfies Meta<typeof MediaCard>;

export default meta;

type Story = StoryObj<typeof meta>;

export const Default: Story = {
  args: {
    media: recommendedMedia[0] as unknown as RecommendedMedia,
  },
};
