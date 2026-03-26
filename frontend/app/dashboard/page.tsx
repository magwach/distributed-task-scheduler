"use client";

import { useEffect, useState, useCallback } from "react";
import Link from "next/link";
import { getTasks, deleteTask } from "@/services/api";
import { Task, TaskUpdateEvent } from "@/services/types";
import StatusBadge from "@/components/StatusBadge";
import ConfirmDialog from "@/components/ConfirmDialog";
import { toast } from "sonner";
import { useWebSocket } from "@/hooks/useWebSockets";

function formatDate(dateStr: string) {
  return new Date(dateStr).toLocaleString("en-US", {
    month: "short",
    day: "numeric",
    hour: "2-digit",
    minute: "2-digit",
  });
}

export default function DashboardPage() {
  const [tasks, setTasks] = useState<Task[]>([]);
  const [loading, setLoading] = useState(true);
  const [deletingId, setDeletingId] = useState<string | null>(null);
  const [confirmId, setConfirmId] = useState<string | null>(null);

  const updateTask = (update: TaskUpdateEvent) => {
    setTasks((prev) =>
      prev.map((task) =>
        task.id === update.task_id
          ? {
              ...task,
              status: update.status,
              updated_at: update.updated_at,
              next_run_at: update.next_run_at ?? task.next_run_at,
              retry_count: update.retry_count ?? task.retry_count,
            }
          : task,
      ),
    );
  };

  useWebSocket({ onMessage: updateTask });

  const fetchTasks = useCallback(async () => {
    try {
      const data = await getTasks();
      setTasks(data ?? []);
    } catch {
      toast.error("Failed to load tasks. Is the API running?");
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => {
    fetchTasks();
  }, [fetchTasks]);

  const handleDeleteConfirmed = async () => {
    if (!confirmId) return;
    const id = confirmId;
    setConfirmId(null);
    setDeletingId(id);
    try {
      await deleteTask(id);
      setTasks((prev) => prev.filter((t) => t.id !== id));
      toast.success("Task deleted successfully.");
    } catch {
      toast.error("Failed to delete task.");
    } finally {
      setDeletingId(null);
    }
  };

  const counts = {
    total: tasks.length,
    pending: tasks.filter((t) => t.status === "pending").length,
    running: tasks.filter((t) => t.status === "running").length,
    failed: tasks.filter((t) => t.status === "failed").length,
  };

  return (
    <>
      <ConfirmDialog
        open={!!confirmId}
        title="Delete Task"
        description="This will permanently delete the task and all its execution history. This action cannot be undone."
        confirmLabel="Delete Task"
        cancelLabel="Cancel"
        dangerous
        onConfirm={handleDeleteConfirmed}
        onCancel={() => setConfirmId(null)}
      />

      <div className="page-header">
        <div>
          <h1 className="page-title">
            Task <span>Dashboard</span>
          </h1>
          <p className="page-subtitle">
            // monitoring {counts.total} registered tasks
          </p>
        </div>
        <div style={{ display: "flex", gap: "10px" }}>
          <button className="btn btn-secondary" onClick={fetchTasks}>
            ↻ Refresh
          </button>
          <Link href="/tasks/new" className="btn btn-primary">
            + New Task
          </Link>
        </div>
      </div>

      <div className="stats-grid">
        <div
          className="stat-card"
          style={
            { "--stat-color": "var(--accent-cyan)" } as React.CSSProperties
          }
        >
          <div className="stat-label">Total Tasks</div>
          <div className="stat-value cyan">{counts.total}</div>
          <div className="stat-icon">◈</div>
        </div>
        <div
          className="stat-card"
          style={
            { "--stat-color": "var(--status-pending)" } as React.CSSProperties
          }
        >
          <div className="stat-label">Pending</div>
          <div className="stat-value yellow">{counts.pending}</div>
          <div className="stat-icon">◷</div>
        </div>
        <div
          className="stat-card"
          style={
            { "--stat-color": "var(--status-running)" } as React.CSSProperties
          }
        >
          <div className="stat-label">Running</div>
          <div className="stat-value cyan">{counts.running}</div>
          <div className="stat-icon">⟳</div>
        </div>
        <div
          className="stat-card"
          style={
            { "--stat-color": "var(--status-failed)" } as React.CSSProperties
          }
        >
          <div className="stat-label">Failed</div>
          <div className="stat-value red">{counts.failed}</div>
          <div className="stat-icon">✕</div>
        </div>
      </div>

      <div className="section-header">
        <span className="section-title">All Tasks</span>
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
              <th>Task</th>
              <th>Schedule</th>
              <th>Status</th>
              <th>Created</th>
              <th>Actions</th>
            </tr>
          </thead>
          <tbody>
            {loading ? (
              <tr className="loading-row">
                <td colSpan={5}>
                  <div className="spinner" />
                  <div
                    style={{
                      color: "var(--text-muted)",
                      fontSize: "12px",
                      fontFamily: "var(--font-mono)",
                    }}
                  >
                    Loading tasks...
                  </div>
                </td>
              </tr>
            ) : tasks.length === 0 ? (
              <tr>
                <td colSpan={5}>
                  <div className="empty-state">
                    <span className="empty-icon">◈</span>
                    <div className="empty-title">No tasks yet</div>
                    <div className="empty-desc">
                      // create your first scheduled task to get started
                    </div>
                    <Link href="/tasks/new" className="btn btn-primary">
                      + Create Task
                    </Link>
                  </div>
                </td>
              </tr>
            ) : (
              tasks.map((task) => (
                <tr key={task.id}>
                  <td>
                    <div className="task-name">{task.title}</div>
                    {task.description && (
                      <div className="task-description">{task.description}</div>
                    )}
                  </td>
                  <td>
                    <span className="task-schedule">{task.schedule}</span>
                  </td>
                  <td>
                    <StatusBadge status={task.status} />
                  </td>
                  <td>
                    <span className="task-time">
                      {formatDate(task.created_at)}
                    </span>
                  </td>
                  <td>
                    <div style={{ display: "flex", gap: "8px" }}>
                      <Link
                        href={`/tasks/${task.id}`}
                        className="btn btn-secondary btn-sm"
                      >
                        View
                      </Link>
                      <button
                        className="btn btn-danger btn-sm"
                        onClick={() => setConfirmId(task.id)}
                        disabled={deletingId === task.id}
                      >
                        {deletingId === task.id ? "..." : "Delete"}
                      </button>
                    </div>
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
