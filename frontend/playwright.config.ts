import { defineConfig } from '@playwright/test'

// Start the real Go API and Vite server before running this suite.
export default defineConfig({
  testDir: './e2e', fullyParallel: false, workers: 1,
  use: { baseURL: 'http://127.0.0.1:5173', browserName: 'chromium', trace: 'retain-on-failure' },
  reporter: 'list',
})
