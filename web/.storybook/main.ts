import type { StorybookConfig } from '@storybook/react-vite';

const config: StorybookConfig = {
  framework: '@storybook/react-vite',
  stories: ['../src/**/*.stories.@(ts|tsx)'],
  addons: ['@storybook/addon-docs'],
  viteFinal: (viteConfig) => {
    // The router plugin would regenerate routeTree.gen.ts; not needed in Storybook.
    const isRouterPlugin = (plugin: unknown): boolean => {
      if (Array.isArray(plugin)) return plugin.some(isRouterPlugin);
      const name = (plugin as { name?: string } | null)?.name;
      return !!name && name.includes('tanstack');
    };
    viteConfig.plugins = (viteConfig.plugins ?? []).filter(
      (plugin) => !isRouterPlugin(plugin),
    );
    return viteConfig;
  },
};

export default config;
