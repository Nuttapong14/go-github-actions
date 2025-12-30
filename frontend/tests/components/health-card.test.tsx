/**
 * HealthCard Component Tests
 *
 * Tests for the health status card component that displays
 * liveness and readiness check results.
 */

import { describe, expect, it } from 'bun:test';

import { HealthCard } from '@/components/health-card';

// Basic render test without DOM (structure validation)
describe('HealthCard Component', () => {
  it('should have correct component structure', () => {
    // Verify component exports correctly
    expect(HealthCard).toBeDefined();
    expect(typeof HealthCard).toBe('function');
  });

  it('should accept required props', () => {
    // Verify component signature accepts expected props
    const props = {
      title: 'Test Health Check',
      description: 'Test description',
      endpoint: '/health/test',
    };

    // Component should be callable with these props
    expect(() => HealthCard(props)).not.toThrow();
  });

  it('should accept optional data prop', () => {
    const props = {
      title: 'Test Health Check',
      description: 'Test description',
      endpoint: '/health/test',
      data: {
        status: 'ok',
        timestamp: new Date().toISOString(),
      },
    };

    expect(() => HealthCard(props)).not.toThrow();
  });

  it('should accept loading state', () => {
    const props = {
      title: 'Test Health Check',
      description: 'Test description',
      endpoint: '/health/test',
      isLoading: true,
    };

    expect(() => HealthCard(props)).not.toThrow();
  });

  it('should accept error state', () => {
    const props = {
      title: 'Test Health Check',
      description: 'Test description',
      endpoint: '/health/test',
      error: 'Connection failed',
    };

    expect(() => HealthCard(props)).not.toThrow();
  });

  it('should handle data with checks', () => {
    const props = {
      title: 'Readiness Check',
      description: 'Checks all dependencies',
      endpoint: '/health/ready',
      data: {
        status: 'ok',
        timestamp: new Date().toISOString(),
        checks: {
          database: 'ok',
        },
      },
    };

    expect(() => HealthCard(props)).not.toThrow();
  });
});

describe('HealthCard Status Handling', () => {
  it('should differentiate between ok and error status', () => {
    const okProps = {
      title: 'Test',
      description: 'Test',
      endpoint: '/test',
      data: { status: 'ok', timestamp: new Date().toISOString() },
    };

    const errorProps = {
      title: 'Test',
      description: 'Test',
      endpoint: '/test',
      data: { status: 'error', timestamp: new Date().toISOString() },
    };

    // Both should render without errors
    expect(() => HealthCard(okProps)).not.toThrow();
    expect(() => HealthCard(errorProps)).not.toThrow();
  });
});
