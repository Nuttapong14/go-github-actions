"use client";

import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";

export interface Deployment {
  id: string;
  environment_id: string;
  release_id: string;
  status: "pending" | "running" | "success" | "failed" | "rolled_back";
  started_at?: string;
  completed_at?: string;
  triggered_by: string;
  trigger_type: "manual" | "automatic" | "hotfix" | "rollback";
  rollback_of_id?: string;
  error_message?: string;
  created_at: string;
  environment?: {
    id: string;
    name: string;
    url: string;
    is_production: boolean;
  };
  release?: {
    id: string;
    version: string;
    git_sha_short: string;
  };
}

interface DeploymentListProps {
  deployments: Deployment[];
  isLoading?: boolean;
  onRollback?: (deploymentId: string) => void;
}

export function DeploymentList({
  deployments,
  isLoading,
  onRollback,
}: DeploymentListProps) {
  const getStatusVariant = (
    status: Deployment["status"]
  ): "default" | "secondary" | "destructive" | "outline" => {
    switch (status) {
      case "success":
        return "default";
      case "pending":
      case "running":
        return "secondary";
      case "failed":
      case "rolled_back":
        return "destructive";
      default:
        return "outline";
    }
  };

  const getStatusColor = (status: Deployment["status"]) => {
    switch (status) {
      case "success":
        return "bg-green-500";
      case "pending":
        return "bg-yellow-500";
      case "running":
        return "bg-blue-500 animate-pulse";
      case "failed":
        return "bg-red-500";
      case "rolled_back":
        return "bg-orange-500";
      default:
        return "bg-gray-500";
    }
  };

  const getTriggerVariant = (
    triggerType: Deployment["trigger_type"]
  ): "default" | "secondary" | "destructive" | "outline" => {
    switch (triggerType) {
      case "automatic":
        return "default";
      case "manual":
        return "secondary";
      case "hotfix":
        return "destructive";
      case "rollback":
        return "outline";
      default:
        return "secondary";
    }
  };

  const formatDate = (dateString?: string) => {
    if (!dateString) return "-";
    return new Date(dateString).toLocaleString();
  };

  const getDuration = (startedAt?: string, completedAt?: string) => {
    if (!startedAt) return "-";
    const start = new Date(startedAt);
    const end = completedAt ? new Date(completedAt) : new Date();
    const diffMs = end.getTime() - start.getTime();
    const diffSec = Math.floor(diffMs / 1000);
    const diffMin = Math.floor(diffSec / 60);
    if (diffMin > 0) {
      return `${diffMin}m ${diffSec % 60}s`;
    }
    return `${diffSec}s`;
  };

  if (isLoading) {
    return (
      <div className="w-full p-8 text-center text-muted-foreground">
        <div className="animate-spin w-8 h-8 border-2 border-primary border-t-transparent rounded-full mx-auto mb-4" />
        Loading deployments...
      </div>
    );
  }

  if (deployments.length === 0) {
    return (
      <div className="w-full p-8 text-center text-muted-foreground">
        No deployments found.
      </div>
    );
  }

  return (
    <div className="rounded-md border">
      <Table>
        <TableHeader>
          <TableRow>
            <TableHead className="w-[100px]">Status</TableHead>
            <TableHead>Environment</TableHead>
            <TableHead>Version</TableHead>
            <TableHead>Trigger</TableHead>
            <TableHead>Started</TableHead>
            <TableHead>Duration</TableHead>
            <TableHead>Triggered By</TableHead>
            <TableHead className="text-right">Actions</TableHead>
          </TableRow>
        </TableHeader>
        <TableBody>
          {deployments.map((deployment) => (
            <TableRow key={deployment.id}>
              <TableCell>
                <Badge variant={getStatusVariant(deployment.status)}>
                  <span
                    className={`w-2 h-2 rounded-full mr-2 ${getStatusColor(
                      deployment.status
                    )}`}
                  />
                  {deployment.status.toUpperCase().replace("_", " ")}
                </Badge>
              </TableCell>
              <TableCell>
                <div className="flex flex-col">
                  <span className="font-medium">
                    {deployment.environment?.name || "Unknown"}
                  </span>
                  {deployment.environment?.is_production && (
                    <Badge variant="destructive" className="w-fit text-xs mt-1">
                      PROD
                    </Badge>
                  )}
                </div>
              </TableCell>
              <TableCell>
                <div className="flex flex-col">
                  <span className="font-mono text-sm">
                    {deployment.release?.version || "Unknown"}
                  </span>
                  <span className="text-xs text-muted-foreground font-mono">
                    {deployment.release?.git_sha_short || ""}
                  </span>
                </div>
              </TableCell>
              <TableCell>
                <Badge variant={getTriggerVariant(deployment.trigger_type)}>
                  {deployment.trigger_type.toUpperCase()}
                </Badge>
                {deployment.rollback_of_id && (
                  <div className="text-xs text-muted-foreground mt-1">
                    Rollback of previous
                  </div>
                )}
              </TableCell>
              <TableCell className="text-sm">
                {formatDate(deployment.started_at)}
              </TableCell>
              <TableCell className="font-mono text-sm">
                {getDuration(deployment.started_at, deployment.completed_at)}
              </TableCell>
              <TableCell className="text-sm">
                {deployment.triggered_by}
              </TableCell>
              <TableCell className="text-right">
                {deployment.status === "success" && onRollback && (
                  <Button
                    variant="outline"
                    size="sm"
                    onClick={() => onRollback(deployment.id)}
                  >
                    Rollback
                  </Button>
                )}
                {deployment.status === "failed" && deployment.error_message && (
                  <span
                    className="text-xs text-destructive"
                    title={deployment.error_message}
                  >
                    View Error
                  </span>
                )}
              </TableCell>
            </TableRow>
          ))}
        </TableBody>
      </Table>
    </div>
  );
}
