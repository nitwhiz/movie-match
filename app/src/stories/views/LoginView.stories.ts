import { Meta, StoryObj } from '@storybook/vue3';
import LoginView from '../../views/LoginView.vue';

const meta = {
  title: 'views/Login',
  component: LoginView,
} satisfies Meta<typeof LoginView>;

export default meta;

type Story = StoryObj<typeof meta>;

export const Default: Story = {};
