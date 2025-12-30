/**
 * Test Utilities
 *
 * Helper functions and utilities for testing React components
 * in the CI/CD Training Frontend.
 */

import React, { ReactElement } from 'react';

/**
 * Create a mock health data response
 */
export function createMockHealthData(
  status: 'ok' | 'error' = 'ok',
  checks?: Record<string, string>
) {
  return {
    status,
    timestamp: new Date().toISOString(),
    ...(checks && { checks }),
  };
}

/**
 * Create a mock fetch response
 */
export function createMockFetchResponse<T>(
  data: T,
  options: { ok?: boolean; status?: number } = {}
) {
  const { ok = true, status = 200 } = options;

  return Promise.resolve({
    ok,
    status,
    statusText: ok ? 'OK' : 'Error',
    json: () => Promise.resolve(data),
    text: () => Promise.resolve(JSON.stringify(data)),
  } as Response);
}

/**
 * Create a mock fetch error
 */
export function createMockFetchError(message = 'Network error') {
  return Promise.reject(new Error(message));
}

/**
 * Wait for a specified time (useful for async tests)
 */
export function wait(ms: number): Promise<void> {
  return new Promise((resolve) => setTimeout(resolve, ms));
}

/**
 * Mock API URL for testing
 */
export const MOCK_API_URL = 'http://localhost:8080';

/**
 * Common test props for HealthCard
 */
export const defaultHealthCardProps = {
  title: 'Test Health Check',
  description: 'A test health check component',
  endpoint: '/health/test',
};

/**
 * Test IDs for component selection
 */
export const testIds = {
  healthCard: 'health-card',
  healthStatus: 'health-status',
  healthEndpoint: 'health-endpoint',
  healthTimestamp: 'health-timestamp',
  healthError: 'health-error',
  loadingIndicator: 'loading-indicator',
};

/**
 * Wrapper component for providing context in tests
 */
export function TestWrapper({ children }: { children: React.ReactNode }) {
  return <>{children}</>;
}

/**
 * Type for render result (simplified for bun test)
 */
export interface RenderResult {
  container: HTMLElement;
  unmount: () => void;
}
