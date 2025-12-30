"use client";

import { useEffect, useState, useCallback } from "react";

import { HealthCard } from "@/components/health-card";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";

interface HealthData {
  status: string;
  timestamp: string;
  checks?: {
    database: string;
  };
}

interface HealthHistory {
  timestamp: Date;
  liveness: string;
  readiness: string;
}

export default function HealthPage() {
  const [livenessData, setLivenessData] = useState<HealthData | null>(null);
  const [readinessData, setReadinessData] = useState<HealthData | null>(null);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [lastRefresh, setLastRefresh] = useState<Date>(new Date());
  const [autoRefresh, setAutoRefresh] = useState(true);
  const [healthHistory, setHealthHistory] = useState<HealthHistory[]>([]);

  const apiUrl = process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080";

  const fetchHealthData = useCallback(async () => {
    setIsLoading(true);
    setError(null);

    try {
      const [livenessRes, readinessRes] = await Promise.all([
        fetch(`${apiUrl}/health/live`),
        fetch(`${apiUrl}/health/ready`),
      ]);

      let livenessStatus = "error";
      let readinessStatus = "error";

      if (livenessRes.ok) {
        const data = await livenessRes.json();
        setLivenessData(data);
        livenessStatus = data.status;
      }

      if (readinessRes.ok) {
        const data = await readinessRes.json();
        setReadinessData(data);
        readinessStatus = data.status;
      }

      // Add to history
      setHealthHistory((prev) => {
        const newEntry = {
          timestamp: new Date(),
          liveness: livenessStatus,
          readiness: readinessStatus,
        };
        return [newEntry, ...prev].slice(0, 10); // Keep last 10 entries
      });

      setLastRefresh(new Date());
    } catch (err) {
      setError(
        err instanceof Error ? err.message : "Failed to fetch health data"
      );
    } finally {
      setIsLoading(false);
    }
  }, [apiUrl]);

  useEffect(() => {
    fetchHealthData();

    if (autoRefresh) {
      const interval = setInterval(fetchHealthData, 10000); // Refresh every 10 seconds
      return () => clearInterval(interval);
    }
  }, [autoRefresh, fetchHealthData]);

  const handleRefresh = () => {
    fetchHealthData();
  };

  return (
    <div className="space-y-8">
      {/* Page Header */}
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-bold">Health Status</h1>
          <p className="text-muted-foreground mt-2">
            Monitor backend service health and dependencies
          </p>
        </div>
        <div className="flex items-center space-x-4">
          <div className="flex items-center space-x-2">
            <input
              type="checkbox"
              id="autoRefresh"
              checked={autoRefresh}
              onChange={(e) => setAutoRefresh(e.target.checked)}
              className="rounded border-gray-300"
            />
            <label htmlFor="autoRefresh" className="text-sm">
              Auto-refresh (10s)
            </label>
          </div>
          <Button onClick={handleRefresh} disabled={isLoading}>
            {isLoading ? "Refreshing..." : "Refresh Now"}
          </Button>
        </div>
      </div>

      {/* Last Refresh */}
      <div className="text-sm text-muted-foreground">
        Last refresh: {lastRefresh.toLocaleTimeString()}
      </div>

      {/* Health Checks */}
      <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
        <HealthCard
          title="Liveness Probe"
          description="Kubernetes liveness check - Is the service running?"
          endpoint="/health/live"
          data={livenessData || undefined}
          isLoading={isLoading}
          error={error || undefined}
        />
        <HealthCard
          title="Readiness Probe"
          description="Kubernetes readiness check - Is the service ready for traffic?"
          endpoint="/health/ready"
          data={readinessData || undefined}
          isLoading={isLoading}
          error={error || undefined}
        />
      </div>

      {/* Endpoint Details */}
      <Card>
        <CardHeader>
          <CardTitle>Endpoint Details</CardTitle>
          <CardDescription>
            Raw response data from health endpoints
          </CardDescription>
        </CardHeader>
        <CardContent>
          <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
            <div>
              <h3 className="text-sm font-medium mb-2">Liveness Response</h3>
              <pre className="p-4 bg-muted rounded-lg text-xs overflow-auto">
                {livenessData
                  ? JSON.stringify(livenessData, null, 2)
                  : "No data"}
              </pre>
            </div>
            <div>
              <h3 className="text-sm font-medium mb-2">Readiness Response</h3>
              <pre className="p-4 bg-muted rounded-lg text-xs overflow-auto">
                {readinessData
                  ? JSON.stringify(readinessData, null, 2)
                  : "No data"}
              </pre>
            </div>
          </div>
        </CardContent>
      </Card>

      {/* Health History */}
      <Card>
        <CardHeader>
          <CardTitle>Health History</CardTitle>
          <CardDescription>Recent health check results</CardDescription>
        </CardHeader>
        <CardContent>
          {healthHistory.length === 0 ? (
            <p className="text-sm text-muted-foreground">
              No health history yet
            </p>
          ) : (
            <div className="space-y-2">
              {healthHistory.map((entry, index) => (
                <div
                  key={index}
                  className="flex items-center justify-between p-3 bg-muted rounded-lg"
                >
                  <span className="text-sm text-muted-foreground">
                    {entry.timestamp.toLocaleTimeString()}
                  </span>
                  <div className="flex items-center space-x-4">
                    <div className="flex items-center space-x-2">
                      <span className="text-sm">Liveness:</span>
                      <Badge
                        variant={
                          entry.liveness === "ok" ? "default" : "destructive"
                        }
                      >
                        {entry.liveness.toUpperCase()}
                      </Badge>
                    </div>
                    <div className="flex items-center space-x-2">
                      <span className="text-sm">Readiness:</span>
                      <Badge
                        variant={
                          entry.readiness === "ok" ? "default" : "destructive"
                        }
                      >
                        {entry.readiness.toUpperCase()}
                      </Badge>
                    </div>
                  </div>
                </div>
              ))}
            </div>
          )}
        </CardContent>
      </Card>

      {/* Metrics Link */}
      <Card>
        <CardHeader>
          <CardTitle>Prometheus Metrics</CardTitle>
          <CardDescription>
            View raw Prometheus metrics from the backend
          </CardDescription>
        </CardHeader>
        <CardContent>
          <a
            href={`${apiUrl}/metrics`}
            target="_blank"
            rel="noopener noreferrer"
            className="inline-flex items-center space-x-2 text-primary hover:underline"
          >
            <span>Open Metrics Endpoint</span>
            <svg
              className="w-4 h-4"
              fill="none"
              stroke="currentColor"
              viewBox="0 0 24 24"
            >
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                strokeWidth={2}
                d="M10 6H6a2 2 0 00-2 2v10a2 2 0 002 2h10a2 2 0 002-2v-4M14 4h6m0 0v6m0-6L10 14"
              />
            </svg>
          </a>
          <p className="mt-2 text-sm text-muted-foreground">
            Custom metrics available: deployments_total,
            health_check_duration_seconds
          </p>
        </CardContent>
      </Card>
    </div>
  );
}
