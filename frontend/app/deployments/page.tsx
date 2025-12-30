"use client";

import { useEffect, useState, useCallback } from "react";

import { DeploymentList, Deployment } from "@/components/deployment-list";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";

interface Environment {
  id: string;
  name: string;
  url: string;
  approval_required: boolean;
  is_production: boolean;
  soak_period_hours: number;
}

interface ListResponse<T> {
  data: T[];
  pagination?: {
    total: number;
    page: number;
    per_page: number;
  };
}

export default function DeploymentsPage() {
  const [deployments, setDeployments] = useState<Deployment[]>([]);
  const [environments, setEnvironments] = useState<Environment[]>([]);
  const [selectedEnvironment, setSelectedEnvironment] = useState<string>("all");
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const apiUrl = process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080";

  const fetchDeployments = useCallback(async () => {
    setIsLoading(true);
    setError(null);

    try {
      let url = `${apiUrl}/api/v1/deployments`;
      if (selectedEnvironment !== "all") {
        url += `?environment_id=${selectedEnvironment}`;
      }

      const response = await fetch(url);
      if (!response.ok) {
        throw new Error("Failed to fetch deployments");
      }

      const data: ListResponse<Deployment> = await response.json();
      setDeployments(data.data || []);
    } catch (err) {
      setError(
        err instanceof Error ? err.message : "Failed to fetch deployments"
      );
    } finally {
      setIsLoading(false);
    }
  }, [apiUrl, selectedEnvironment]);

  const fetchEnvironments = useCallback(async () => {
    try {
      const response = await fetch(`${apiUrl}/api/v1/environments`);
      if (!response.ok) {
        throw new Error("Failed to fetch environments");
      }

      const data: ListResponse<Environment> = await response.json();
      setEnvironments(data.data || []);
    } catch (err) {
      console.error("Failed to fetch environments:", err);
    }
  }, [apiUrl]);

  useEffect(() => {
    fetchEnvironments();
  }, [fetchEnvironments]);

  useEffect(() => {
    fetchDeployments();
    const interval = setInterval(fetchDeployments, 10000); // Refresh every 10 seconds
    return () => clearInterval(interval);
  }, [fetchDeployments]);

  const handleRollback = async (deploymentId: string) => {
    try {
      const response = await fetch(
        `${apiUrl}/api/v1/deployments/${deploymentId}/rollback`,
        {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
          },
          body: JSON.stringify({
            triggered_by: "web-ui",
          }),
        }
      );

      if (!response.ok) {
        throw new Error("Failed to initiate rollback");
      }

      // Refresh deployments after rollback
      await fetchDeployments();
    } catch (err) {
      setError(
        err instanceof Error ? err.message : "Failed to initiate rollback"
      );
    }
  };

  const getDeploymentStats = () => {
    const total = deployments.length;
    const successful = deployments.filter((d) => d.status === "success").length;
    const failed = deployments.filter((d) => d.status === "failed").length;
    const running = deployments.filter(
      (d) => d.status === "running" || d.status === "pending"
    ).length;

    return { total, successful, failed, running };
  };

  const stats = getDeploymentStats();

  return (
    <div className="space-y-8">
      {/* Page Header */}
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-bold">Deployments</h1>
          <p className="text-muted-foreground mt-2">
            Monitor and manage deployment history across environments
          </p>
        </div>
        <Button onClick={fetchDeployments} variant="outline">
          <svg
            className="w-4 h-4 mr-2"
            fill="none"
            stroke="currentColor"
            viewBox="0 0 24 24"
          >
            <path
              strokeLinecap="round"
              strokeLinejoin="round"
              strokeWidth={2}
              d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15"
            />
          </svg>
          Refresh
        </Button>
      </div>

      {/* Statistics Cards */}
      <div className="grid grid-cols-1 md:grid-cols-4 gap-4">
        <Card>
          <CardHeader className="pb-2">
            <CardDescription>Total Deployments</CardDescription>
            <CardTitle className="text-4xl">{stats.total}</CardTitle>
          </CardHeader>
        </Card>

        <Card>
          <CardHeader className="pb-2">
            <CardDescription>Successful</CardDescription>
            <CardTitle className="text-4xl text-green-600">
              {stats.successful}
            </CardTitle>
          </CardHeader>
        </Card>

        <Card>
          <CardHeader className="pb-2">
            <CardDescription>Failed</CardDescription>
            <CardTitle className="text-4xl text-red-600">
              {stats.failed}
            </CardTitle>
          </CardHeader>
        </Card>

        <Card>
          <CardHeader className="pb-2">
            <CardDescription>In Progress</CardDescription>
            <CardTitle className="text-4xl text-blue-600">
              {stats.running}
            </CardTitle>
          </CardHeader>
        </Card>
      </div>

      {/* Filters */}
      <Card>
        <CardHeader>
          <CardTitle>Filters</CardTitle>
        </CardHeader>
        <CardContent>
          <div className="flex items-center gap-4">
            <div className="flex items-center gap-2">
              <label className="text-sm font-medium">Environment:</label>
              <Select
                value={selectedEnvironment}
                onValueChange={setSelectedEnvironment}
              >
                <SelectTrigger className="w-[200px]">
                  <SelectValue placeholder="Select environment" />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="all">All Environments</SelectItem>
                  {environments.map((env) => (
                    <SelectItem key={env.id} value={env.id}>
                      <div className="flex items-center gap-2">
                        <span>{env.name}</span>
                        {env.is_production && (
                          <Badge variant="destructive" className="text-xs">
                            PROD
                          </Badge>
                        )}
                      </div>
                    </SelectItem>
                  ))}
                </SelectContent>
              </Select>
            </div>
          </div>
        </CardContent>
      </Card>

      {/* Error Message */}
      {error && (
        <Card className="border-destructive">
          <CardContent className="pt-6">
            <div className="flex items-center gap-2 text-destructive">
              <svg
                className="w-5 h-5"
                fill="none"
                stroke="currentColor"
                viewBox="0 0 24 24"
              >
                <path
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  strokeWidth={2}
                  d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"
                />
              </svg>
              <span>{error}</span>
            </div>
          </CardContent>
        </Card>
      )}

      {/* Deployment List */}
      <Card>
        <CardHeader>
          <CardTitle>Deployment History</CardTitle>
          <CardDescription>
            View and manage deployments across all environments
          </CardDescription>
        </CardHeader>
        <CardContent>
          <DeploymentList
            deployments={deployments}
            isLoading={isLoading}
            onRollback={handleRollback}
          />
        </CardContent>
      </Card>

      {/* Environment Status */}
      {environments.length > 0 && (
        <Card>
          <CardHeader>
            <CardTitle>Environments</CardTitle>
            <CardDescription>
              Available deployment targets and their configuration
            </CardDescription>
          </CardHeader>
          <CardContent>
            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
              {environments.map((env) => (
                <div
                  key={env.id}
                  className="p-4 bg-muted rounded-lg space-y-2"
                >
                  <div className="flex items-center justify-between">
                    <span className="font-medium">{env.name}</span>
                    {env.is_production && (
                      <Badge variant="destructive">PROD</Badge>
                    )}
                  </div>
                  <a
                    href={env.url}
                    target="_blank"
                    rel="noopener noreferrer"
                    className="text-sm text-muted-foreground hover:underline truncate block"
                  >
                    {env.url}
                  </a>
                  <div className="flex flex-wrap gap-1">
                    {env.approval_required && (
                      <Badge variant="outline" className="text-xs">
                        Approval Required
                      </Badge>
                    )}
                    {env.soak_period_hours > 0 && (
                      <Badge variant="outline" className="text-xs">
                        {env.soak_period_hours}h Soak
                      </Badge>
                    )}
                  </div>
                </div>
              ))}
            </div>
          </CardContent>
        </Card>
      )}
    </div>
  );
}
