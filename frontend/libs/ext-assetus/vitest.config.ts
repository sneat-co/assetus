import { defineConfig } from 'vitest/config';

// A service spec imports @sneat/api, which transitively pulls @ionic/angular.
// @ionic/angular does a bare directory import of '@ionic/core/components' that
// Node's ESM resolver rejects; redirect it to the package's index.js. Inlining
// the @ionic/@sneat packages keeps the rest of the chain bundled by Vitest.
export default defineConfig({
  resolve: {
    alias: [
      {
        find: /^@ionic\/core\/components$/,
        replacement: '@ionic/core/components/index.js',
      },
    ],
  },
  test: {
    server: { deps: { inline: [/@ionic/, /ionicons/, /@sneat/] } },
  },
});
