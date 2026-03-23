"use client";

import { useEffect, useState, useCallback } from "react";
import { useParams } from "next/navigation";
import Link from "next/link";
import { getTask } from "@/services/api";
import { Task, TaskExecution } from "@/services/types";
import StatusBadge from "@/components/StatusBadge";
import { toast } from "sonner";

function formatDate(dateStr: string | null) {
  if (!dateStr) return "—";
  return new Date(dateStr).toLocaleString("en-US", {
    month: "short",
    day: "numeric",
    hour: "2-digit",
    minute: "2-digit",
    second: "2-digit",
  });
}

function duration(start: string, end: string | null) {
  if (!end) return "running...";
  const ms = new Date(end).getTime() - new Date(start).getTime();
  if (ms < 1000) return `${ms}ms`;
  return `${(ms / 1000).toFixed(1)}s`;
}

export default function TaskDetailPage() {
  const { id } = useParams<{ id: string }>();
  const [task, setTask] = useState<Task | null>(null);
  const [loading, setLoading] = useState(true);

  const fetchTask = useCallback(async () => {
    try {
      const data = await getTask(id);
      setTask(data);
    } catch {
      toast.error("Failed to load task details.");
    } finally {
      setLoading(false);
    }
  }, [id]);

  useEffect(() => {
    fetchTask();
    const interval = setInterval(fetchTask, 5000);
    return () => clearInterval(interval);
  }, [fetchTask]);

  if (loading) {
    return (
      <div
        style={{ display: "flex", justifyContent: "center", paddingTop: 120 }}
      >
        <div>
          <div className="spinner" />
          <div
            style={{
              color: "var(--text-muted)",
              fontSize: "12px",
              fontFamily: "var(--font-mono)",
              textAlign: "center",
              marginTop: 12,
            }}
          >
            Loading task...
          </div>
        </div>
      </div>
    );
  }

  if (!task) {
    return (
      <>
        <div className="alert alert-error">Task not found.</div>
        <Link href="/dashboard" className="btn btn-secondary">
          ← Back to Dashboard
        </Link>
      </>
    );
  }

  return (
    <>
      <div className="page-header">
        <div>
          <div className="breadcrumb">
            <Link href="/dashboard">Dashboard</Link>
            <span className="breadcrumb-sep">/</span>
            <Link href="/tasks">Tasks</Link>
            <span className="breadcrumb-sep">/</span>
            <span>{task.title}</span>
          </div>
          <h1 className="page-title">
            Task <span>Detail</span>
          </h1>
          <p className="page-subtitle">// {task.id}</p>
        </div>
        <Link href="/dashboard" className="btn btn-secondary">
          ← Back
        </Link>
      </div>

      {/* Task Info Card */}
      <div
        className="form-card"
        style={{
          maxWidth: "100%",
          marginBottom: 32,
          display: "grid",
          gridTemplateColumns: "repeat(3, 1fr)",
          gap: 32,
        }}
      >
        <div>
          <div className="form-label">Task Name</div>
          <div
            style={{
              fontSize: 18,
              fontWeight: 700,
              color: "var(--text-primary)",
            }}
          >
            {task.title}
          </div>
          {task.description && (
            <div
              style={{
                marginTop: 6,
                fontSize: 13,
                color: "var(--text-secondary)",
                fontFamily: "var(--font-mono)",
              }}
            >
              {task.description}
            </div>
          )}
        </div>
        <div>
          <div className="form-label">Schedule</div>
          <span className="task-schedule" style={{ fontSize: 14 }}>
            {task.schedule}
          </span>
        </div>
        <div>
          <div className="form-label">Current Status</div>
          <StatusBadge status={task.status} />
        </div>
        <div>
          <div className="form-label">Created At</div>
          <div
            style={{
              fontFamily: "var(--font-mono)",
              fontSize: 13,
              color: "var(--text-secondary)",
            }}
          >
            {formatDate(task.created_at)}
          </div>
        </div>
        <div>
          <div className="form-label">Last Updated</div>
          <div
            style={{
              fontFamily: "var(--font-mono)",
              fontSize: 13,
              color: "var(--text-secondary)",
            }}
          >
            {formatDate(task.updated_at)}
          </div>
        </div>
        <div>
          <div className="form-label">Executions</div>
          <div
            style={{
              fontSize: 24,
              fontWeight: 800,
              color: "var(--accent-cyan)",
            }}
          >
            {task.executions?.length ?? 0}
          </div>
        </div>
      </div>

      {/* Executions */}
      <div className="section-header">
        <span className="section-title">Execution History</span>
        <span
          style={{
            fontFamily: "var(--font-mono)",
            fontSize: "11px",
            color: "var(--text-muted)",
          }}
        >
          auto-refresh every 5s
        </span>
      </div>

      <div className="table-container">
        <table className="tasks-table">
          <thead>
            <tr>
              <th>Execution ID</th>
              <th>Status</th>
              <th>Started At</th>
              <th>Finished At</th>
              <th>Duration</th>
              <th>Error</th>
            </tr>
          </thead>
          <tbody>
            {!task.executions || task.executions.length === 0 ? (
              <tr>
                <td colSpan={6}>
                  <div className="empty-state">
                    <span className="empty-icon">◷</span>
                    <div className="empty-title">No executions yet</div>
                    <div className="empty-desc">
                      // waiting for scheduler to pick up this task
                    </div>
                  </div>
                </td>
              </tr>
            ) : (
              [...task.executions].reverse().map((exec: TaskExecution) => (
                <tr key={exec.id}>
                  <td>
                    <span
                      style={{
                        fontFamily: "var(--font-mono)",
                        fontSize: 11,
                        color: "var(--text-muted)",
                      }}
                    >
                      {exec.id.slice(0, 8)}...
                    </span>
                  </td>
                  <td>
                    <StatusBadge status={exec.status} />
                  </td>
                  <td>
                    <span className="task-time">
                      {formatDate(exec.started_at)}
                    </span>
                  </td>
                  <td>
                    <span className="task-time">
                      {formatDate(exec.finished_at)}
                    </span>
                  </td>
                  <td>
                    <span
                      style={{
                        fontFamily: "var(--font-mono)",
                        fontSize: 12,
                        color: "var(--text-secondary)",
                      }}
                    >
                      {duration(exec.started_at, exec.finished_at)}
                    </span>
                  </td>
                  <td>
                    {exec.error_message ? (
                      <span
                        style={{
                          fontFamily: "var(--font-mono)",
                          fontSize: 11,
                          color: "var(--status-failed)",
                        }}
                      >
                        {exec.error_message}
                      </span>
                    ) : (
                      <span
                        style={{ color: "var(--text-muted)", fontSize: 12 }}
                      >
                        —
                      </span>
                    )}
                  </td>
                </tr>
              ))
            )}
          </tbody>
        </table>
      </div>
    </>
  );
}
