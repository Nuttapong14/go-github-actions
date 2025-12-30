/**
 * Test Setup Configuration
 *
 * This file configures the testing environment for the CI/CD Training Frontend.
 * It sets up testing utilities and mocks for React components.
 */

import { mock, beforeEach, afterEach, beforeAll, afterAll } from "bun:test";
import '@testing-library/jest-dom';

// Mock Next.js router
const mockRouter = {
  push: () => Promise.resolve(true),
  replace: () => Promise.resolve(true),
  prefetch: () => Promise.resolve(),
  back: () => {},
  forward: () => {},
  refresh: () => {},
  pathname: '/',
  route: '/',
  query: {},
  asPath: '/',
  basePath: '',
  isLocaleDomain: false,
  isReady: true,
  isPreview: false,
};

// Mock useRouter hook
mock.module('next/navigation', () => ({
  useRouter: () => mockRouter,
  usePathname: () => '/',
  useSearchParams: () => new URLSearchParams(),
}));

// Mock fetch for API calls
global.fetch = mock(() =>
  Promise.resolve({
    ok: true,
    json: () => Promise.resolve({}),
    text: () => Promise.resolve(''),
    status: 200,
    statusText: 'OK',
  } as Response)
);

// Reset mocks before each test
beforeEach(() => {
  // jest.clearAllMocks();
});

// Clean up after tests
afterEach(() => {
  // jest.restoreAllMocks();
});

// Suppress console errors during tests (optional)
const originalError = console.error;
beforeAll(() => {
  console.error = (...args: unknown[]) => {
    if (
      typeof args[0] === 'string' &&
      args[0].includes('Warning: ReactDOM.render')
    ) {
      return;
    }
    originalError.call(console, ...args);
  };
});

afterAll(() => {
  console.error = originalError;
});

export { mockRouter };
