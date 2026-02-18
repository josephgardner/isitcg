import { defineConfig } from 'vitest/config'

export default defineConfig({
  test: {
    include: ['docs/src/**/*.test.js'],
  },
})
