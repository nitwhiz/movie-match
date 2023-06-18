import type { Preview } from '@storybook/vue3';
import { INITIAL_VIEWPORTS } from '@storybook/addon-viewport';

import '../src/assets/styles/index.scss';
import Poster from '../src/common/Poster';
import { vueRouter } from 'storybook-vue3-router';
import router from '../src/router';

Poster.startGC();

const preview: Preview = {
  decorators: [
    (story) => ({
      components: {
        story,
      },
      template: '<Suspense><story /></Suspense>',
    }),
    vueRouter(router.getRoutes()),
  ],
  parameters: {
    viewport: {
      viewports: INITIAL_VIEWPORTS,
    },
    layout: 'centered',
    actions: { argTypesRegex: '^on[A-Z].*' },
    controls: {
      matchers: {
        color: /(background|color)$/i,
        date: /Date$/,
      },
    },
  },
};

export default preview;
