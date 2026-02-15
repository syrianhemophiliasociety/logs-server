import { defineConfig, devices } from "@playwright/test";

export default defineConfig({
  testDir: "./tests",
  fullyParallel: false,
  forbidOnly: !!process.env.CI,
  retries: process.env.CI ? 2 : 0,
  workers: process.env.CI ? 1 : 1,
  reporter: "dot",
  use: {
    baseURL: "http://localhost:3000",
    trace: "on-first-retry",
  },

  projects: [
    {
      name: "chromium",
      use: { ...devices["Desktop Chrome"] },
    },
  ],

  webServer: {
    env: {
      DB_NAME: "shsdb",
      DB_HOST: "localhost:3306",
      DB_USERNAME: "root",
      DB_PASSWORD: "previetcomrade",
      CACHE_HOST: "localhost:6379",
      CACHE_PASSWORD: "previetcomrade",
      SUPERADMIN_USERNAME: "b",
      SUPERADMIN_PASSWORD: "kurwamatch",
    },
    command: "make shs-server-local",
    url: "http://127.0.0.1:3000/v1/status",
    reuseExistingServer: true,
  },
});
