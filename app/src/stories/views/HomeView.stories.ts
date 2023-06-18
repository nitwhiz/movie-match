import { Meta, StoryObj } from '@storybook/vue3';
import HomeView from '../../views/HomeView.vue';

const meta = {
  title: 'views/Home',
  component: HomeView,
} satisfies Meta<typeof HomeView>;

export default meta;

type Story = StoryObj<typeof meta>;

export const Default: Story = {};
