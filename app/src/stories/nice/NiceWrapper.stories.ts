import type { Meta, StoryObj } from '@storybook/vue3';
import NiceWrapper from '../../components/nice/NiceWrapper.vue';

const colorOptions = [
  {
    label: 'Aqua',
    colors: ['rgb(148, 55, 255)', 'rgb(101, 229, 255)'],
  },
  {
    label: 'Ice Cream',
    colors: ['rgb(255, 55, 140)', 'rgb(187, 255, 101)'],
  },
  {
    label: 'RGB',
    colors: ['red', 'green', 'blue'],
  },
];

const meta = {
  title: 'nice/NiceWrapper',
  component: NiceWrapper,
  argTypes: {
    colors: {
      options: colorOptions.map((c) => c.label),
      mapping: colorOptions.reduce(
        (a, c) => ({
          ...a,
          [c.label]: c.colors,
        }),
        {}
      ),
    },
  },
} satisfies Meta<typeof NiceWrapper>;

export default meta;

type Story = StoryObj<typeof meta>;

// @ts-ignore
const baseStory = (): Story => ({
  render: (args) => ({
    components: {
      NiceWrapper,
    },
    setup: () => ({ args }),
    template: '<NiceWrapper v-bind="args">Hello World!</NiceWrapper>',
  }),
});

// @ts-ignore
export const Default: Story = {
  ...baseStory(),
};

// @ts-ignore
export const WithInput: Story = {
  render: (args) => ({
    components: {
      NiceWrapper,
    },
    setup: () => ({ args }),
    template:
      '<NiceWrapper v-bind="args"><input placeholder="Type Here" type="text" /></NiceWrapper>',
  }),
};
