"use client";

import { Badge } from "@/components/ui/badge";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";

interface HealthCheck {
  database: string;
}

interface HealthData {
  status: string;
  timestamp: string;
  checks?: HealthCheck;
}

interface HealthCardProps {
  title: string;
  description: string;
  endpoint: string;
  data?: HealthData;
  isLoading?: boolean;
  error?: string;
}

export function HealthCard({
  title,
  description,
  endpoint,
  data,
  isLoading,
  error,
}: HealthCardProps) {
  const getStatusVariant = (status: string) => {
    switch (status) {
      case "ok":
        return "default";
      case "error":
        return "destructive";
      default:
        return "secondary";
    }
  };

  const getStatusColor = (status: string) => {
    switch (status) {
      case "ok":
        return "bg-green-500";
      case "error":
        return "bg-red-500";
      default:
        return "bg-yellow-500";
    }
  };

  return (
    <Card className="w-full">
      <CardHeader>
        <div className="flex items-center justify-between">
          <div>
            <CardTitle className="text-lg">{title}</CardTitle>
            <CardDescription>{description}</CardDescription>
          </div>
          {!isLoading && !error && data && (
            <Badge variant={getStatusVariant(data.status)}>
              <span
                className={`w-2 h-2 rounded-full mr-2 ${getStatusColor(
                  data.status
                )}`}
              />
              {data.status.toUpperCase()}
            </Badge>
          )}
          {isLoading && (
            <Badge variant="secondary">
              <span className="w-2 h-2 rounded-full mr-2 bg-yellow-500 animate-pulse" />
              LOADING
            </Badge>
          )}
          {error && (
            <Badge variant="destructive">
              <span className="w-2 h-2 rounded-full mr-2 bg-red-500" />
              ERROR
            </Badge>
          )}
        </div>
      </CardHeader>
      <CardContent>
        <div className="space-y-2">
          <div className="flex items-center text-sm text-muted-foreground">
            <span className="font-medium">Endpoint:</span>
            <code className="ml-2 px-2 py-1 bg-muted rounded text-xs">
              {endpoint}
            </code>
          </div>

          {error && (
            <div className="text-sm text-destructive">
              <span className="font-medium">Error:</span> {error}
            </div>
          )}

          {data && (
            <>
              <div className="text-sm text-muted-foreground">
                <span className="font-medium">Last Checked:</span>{" "}
                {new Date(data.timestamp).toLocaleString()}
              </div>

              {data.checks && (
                <div className="mt-4">
                  <h4 className="text-sm font-medium mb-2">Dependency Checks</h4>
                  <div className="grid gap-2">
                    {Object.entries(data.checks).map(([key, value]) => (
                      <div
                        key={key}
                        className="flex items-center justify-between px-3 py-2 bg-muted rounded-md"
                      >
                        <span className="text-sm capitalize">{key}</span>
                        <Badge
                          variant={value === "ok" ? "default" : "destructive"}
                        >
                          {value.toUpperCase()}
                        </Badge>
                      </div>
                    ))}
                  </div>
                </div>
              )}
            </>
          )}
        </div>
      </CardContent>
    </Card>
  );
}
