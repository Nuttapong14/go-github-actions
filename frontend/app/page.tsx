"use client";

import { useCallback, useEffect, useState } from "react";

import { HealthCard } from "@/components/health-card";
import { Badge } from "@/components/ui/badge";
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

export default function DashboardPage() {
  const [livenessData, setLivenessData] = useState<HealthData | null>(null);
  const [readinessData, setReadinessData] = useState<HealthData | null>(null);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const apiUrl = process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080";

  const fetchHealthData = useCallback(async () => {
    setIsLoading(true);
    setError(null);

    try {
      const [livenessRes, readinessRes] = await Promise.all([
        fetch(`${apiUrl}/health/live`),
        fetch(`${apiUrl}/health/ready`),
      ]);

      if (livenessRes.ok) {
        const data = await livenessRes.json();
        setLivenessData(data);
      }

      if (readinessRes.ok) {
        const data = await readinessRes.json();
        setReadinessData(data);
      }
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to fetch health data");
    } finally {
      setIsLoading(false);
    }
  }, [apiUrl]);

  useEffect(() => {
    fetchHealthData();
    const interval = setInterval(fetchHealthData, 30000); // Refresh every 30 seconds
    return () => clearInterval(interval);
  }, [fetchHealthData]);

  const getOverallStatus = () => {
    if (isLoading) return "loading";
    if (error) return "error";
    if (livenessData?.status === "ok" && readinessData?.status === "ok") {
      return "healthy";
    }
    return "unhealthy";
  };

  const overallStatus = getOverallStatus();

  return (
    <div className="space-y-8">
      {/* Page Header */}
      <div>
        <h1 className="text-3xl font-bold">Dashboard</h1>
        <p className="text-muted-foreground mt-2">
          Monitor your CI/CD Training Application
        </p>
      </div>

      {/* Overall Status */}
      <Card>
        <CardHeader>
          <div className="flex items-center justify-between">
            <div>
              <CardTitle>System Status</CardTitle>
              <CardDescription>
                Overall health of the CI/CD Training Application
              </CardDescription>
            </div>
            <Badge
              variant={
                overallStatus === "healthy"
                  ? "default"
                  : overallStatus === "loading"
                  ? "secondary"
                  : "destructive"
              }
              className="text-sm"
            >
              {overallStatus === "healthy" && (
                <span className="w-2 h-2 rounded-full mr-2 bg-green-500" />
              )}
              {overallStatus === "loading" && (
                <span className="w-2 h-2 rounded-full mr-2 bg-yellow-500 animate-pulse" />
              )}
              {(overallStatus === "error" || overallStatus === "unhealthy") && (
                <span className="w-2 h-2 rounded-full mr-2 bg-red-500" />
              )}
              {overallStatus.toUpperCase()}
            </Badge>
          </div>
        </CardHeader>
        <CardContent>
          <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
            <div className="flex items-center space-x-3 p-4 bg-muted rounded-lg">
              <div className="w-10 h-10 rounded-full bg-primary/10 flex items-center justify-center">
                <svg
                  className="w-5 h-5 text-primary"
                  fill="none"
                  stroke="currentColor"
                  viewBox="0 0 24 24"
                >
                  <path
                    strokeLinecap="round"
                    strokeLinejoin="round"
                    strokeWidth={2}
                    d="M5 12h14M12 5l7 7-7 7"
                  />
                </svg>
              </div>
              <div>
                <p className="text-sm text-muted-foreground">Backend</p>
                <p className="font-medium">
                  {livenessData?.status === "ok" ? "Running" : "Unknown"}
                </p>
              </div>
            </div>

            <div className="flex items-center space-x-3 p-4 bg-muted rounded-lg">
              <div className="w-10 h-10 rounded-full bg-primary/10 flex items-center justify-center">
                <svg
                  className="w-5 h-5 text-primary"
                  fill="none"
                  stroke="currentColor"
                  viewBox="0 0 24 24"
                >
                  <path
                    strokeLinecap="round"
                    strokeLinejoin="round"
                    strokeWidth={2}
                    d="M4 7v10c0 2.21 3.582 4 8 4s8-1.79 8-4V7M4 7c0 2.21 3.582 4 8 4s8-1.79 8-4M4 7c0-2.21 3.582-4 8-4s8 1.79 8 4"
                  />
                </svg>
              </div>
              <div>
                <p className="text-sm text-muted-foreground">Database</p>
                <p className="font-medium">
                  {readinessData?.checks?.database === "ok"
                    ? "Connected"
                    : "Unknown"}
                </p>
              </div>
            </div>

            <div className="flex items-center space-x-3 p-4 bg-muted rounded-lg">
              <div className="w-10 h-10 rounded-full bg-primary/10 flex items-center justify-center">
                <svg
                  className="w-5 h-5 text-primary"
                  fill="none"
                  stroke="currentColor"
                  viewBox="0 0 24 24"
                >
                  <path
                    strokeLinecap="round"
                    strokeLinejoin="round"
                    strokeWidth={2}
                    d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z"
                  />
                </svg>
              </div>
              <div>
                <p className="text-sm text-muted-foreground">Metrics</p>
                <p className="font-medium">Available</p>
              </div>
            </div>
          </div>
        </CardContent>
      </Card>

      {/* Health Checks */}
      <div>
        <h2 className="text-xl font-semibold mb-4">Health Checks</h2>
        <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
          <HealthCard
            title="Liveness Probe"
            description="Checks if the service is running"
            endpoint="/health/live"
            data={livenessData || undefined}
            isLoading={isLoading}
            error={error || undefined}
          />
          <HealthCard
            title="Readiness Probe"
            description="Checks if the service is ready to accept traffic"
            endpoint="/health/ready"
            data={readinessData || undefined}
            isLoading={isLoading}
            error={error || undefined}
          />
        </div>
      </div>

      {/* Quick Links */}
      <Card>
        <CardHeader>
          <CardTitle>Quick Links</CardTitle>
          <CardDescription>
            Useful endpoints and documentation
          </CardDescription>
        </CardHeader>
        <CardContent>
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
            <a
              href={`${apiUrl}/health/live`}
              target="_blank"
              rel="noopener noreferrer"
              className="p-4 bg-muted rounded-lg hover:bg-muted/80 transition-colors"
            >
              <h3 className="font-medium">Liveness Endpoint</h3>
              <p className="text-sm text-muted-foreground">/health/live</p>
            </a>
            <a
              href={`${apiUrl}/health/ready`}
              target="_blank"
              rel="noopener noreferrer"
              className="p-4 bg-muted rounded-lg hover:bg-muted/80 transition-colors"
            >
              <h3 className="font-medium">Readiness Endpoint</h3>
              <p className="text-sm text-muted-foreground">/health/ready</p>
            </a>
            <a
              href={`${apiUrl}/metrics`}
              target="_blank"
              rel="noopener noreferrer"
              className="p-4 bg-muted rounded-lg hover:bg-muted/80 transition-colors"
            >
              <h3 className="font-medium">Prometheus Metrics</h3>
              <p className="text-sm text-muted-foreground">/metrics</p>
            </a>
            <a
              href="/deployments"
              className="p-4 bg-muted rounded-lg hover:bg-muted/80 transition-colors"
            >
              <h3 className="font-medium">Deployments</h3>
              <p className="text-sm text-muted-foreground">View deployment history</p>
            </a>
          </div>
        </CardContent>
      </Card>
    </div>
  );
}
