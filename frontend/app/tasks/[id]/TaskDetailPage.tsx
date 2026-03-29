'use client'

import { useEffect, useState, useCallback } from 'react'
import { useParams } from 'next/navigation'
import Link from 'next/link'
import { getTask, retryTask } from '@/services/api'
import { Task, TaskExecution, TaskLog } from '@/services/types'
import StatusBadge from '@/components/StatusBadge'
import { toast } from 'sonner'
import { TaskUpdateEvent } from '@/services/types'
import { useWebSocket } from '@/hooks/useWebSockets'

type Tab = 'executions' | 'logs'

function formatDate(dateStr: string | null) {
  if (!dateStr) return '—'
  return new Date(dateStr).toLocaleString('en-US', {
    month: 'short',
    day: 'numeric',
    hour: '2-digit',
    minute: '2-digit',
    second: '2-digit',
  })
}

function duration(start: string, end: string | null) {
  if (!end) return 'running...'
  const ms = new Date(end).getTime() - new Date(start).getTime()
  if (ms < 1000) return `${ms}ms`
  return `${(ms / 1000).toFixed(1)}s`
}

function LogLevelBadge({ level }: { level: string }) {
  const styles: Record<string, React.CSSProperties> = {
    info: {
      background: 'var(--accent-cyan-dim)',
      color: 'var(--accent-cyan)',
      border: '1px solid rgba(0,212,255,0.2)',
    },
    warn: {
      background: 'var(--status-pending-bg)',
      color: 'var(--status-pending)',
      border: '1px solid rgba(245,158,11,0.2)',
    },
    warning: {
      background: 'var(--status-pending-bg)',
      color: 'var(--status-pending)',
      border: '1px solid rgba(245,158,11,0.2)',
    },
    error: {
      background: 'var(--status-failed-bg)',
      color: 'var(--status-failed)',
      border: '1px solid rgba(239,68,68,0.2)',
    },
  }

  const style = styles[level.toLowerCase()] ?? styles.info

  return (
    <span
      style={{
        ...style,
        padding: '2px 8px',
        borderRadius: '4px',
        fontSize: '10px',
        fontFamily: 'var(--font-mono)',
        fontWeight: 600,
        letterSpacing: '0.08em',
        textTransform: 'uppercase',
        whiteSpace: 'nowrap',
      }}
    >
      {level}
    </span>
  )
}

export default function TaskDetailPage() {
  const { id } = useParams<{ id: string }>()
  const [task, setTask] = useState<Task | null>(null)
  const [loading, setLoading] = useState(true)
  const [retrying, setRetrying] = useState(false)
  const [activeTab, setActiveTab] = useState<Tab>('executions')

  const fetchTask = useCallback(async () => {
    try {
      const data = await getTask(id)
      setTask(data)
    } catch {
      toast.error('Failed to load task details.')
    } finally {
      setLoading(false)
    }
  }, [id])

  useEffect(() => {
    fetchTask()
  }, [fetchTask])

  const handleWsMessage = useCallback((update: TaskUpdateEvent) => {
    if (update.task_id !== id) return
    setTask((prev) => {
      if (!prev) return prev
      return {
        ...prev,
        status: update.status,
        updated_at: update.updated_at,
        next_run_at: update.next_run_at ?? prev.next_run_at,
        retry_count: update.retry_count ?? prev.retry_count,
      }
    })
    fetchTask()
  }, [id, fetchTask])

  useWebSocket({ onMessage: handleWsMessage })

  const handleRetry = async () => {
    if (!task) return
    setRetrying(true)
    try {
      const result = await retryTask(task.id)
      if (result === 'success') {
        toast.success('Task queued for retry.')
        fetchTask()
      } else {
        toast.error('Failed to retry task.')
      }
    } finally {
      setRetrying(false)
    }
  }

  const allLogs: TaskLog[] = task?.executions
    ?.flatMap((e) => e.logs ?? [])
    .sort((a, b) => new Date(a.created_at).getTime() - new Date(b.created_at).getTime()) ?? []

  if (loading) {
    return (
      <div style={{ display: 'flex', justifyContent: 'center', paddingTop: 120 }}>
        <div>
          <div className="spinner" />
          <div style={{ color: 'var(--text-muted)', fontSize: '12px', fontFamily: 'var(--font-mono)', textAlign: 'center', marginTop: 12 }}>
            Loading task...
          </div>
        </div>
      </div>
    )
  }

  if (!task) {
    return (
      <>
        <div className="alert alert-error">Task not found.</div>
        <Link href="/dashboard" className="btn btn-secondary">← Back to Dashboard</Link>
      </>
    )
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
        <div style={{ display: 'flex', gap: 10, alignItems: 'center' }}>
          {task.status === 'failed' && (
            <button
              className="btn btn-primary"
              onClick={handleRetry}
              disabled={retrying}
            >
              {retrying ? '⟳ Retrying...' : '↺ Retry Task'}
            </button>
          )}
          <Link href="/dashboard" className="btn btn-secondary">← Back</Link>
        </div>
      </div>

      {/* Task Info Card */}
      <div
        className="form-card"
        style={{ maxWidth: '100%', marginBottom: 32, display: 'grid', gridTemplateColumns: 'repeat(3, 1fr)', gap: 32 }}
      >
        <div>
          <div className="form-label">Task Name</div>
          <div style={{ fontSize: 18, fontWeight: 700, color: 'var(--text-primary)' }}>{task.title}</div>
          {task.description && (
            <div style={{ marginTop: 6, fontSize: 13, color: 'var(--text-secondary)', fontFamily: 'var(--font-mono)' }}>
              {task.description}
            </div>
          )}
        </div>

        <div>
          <div className="form-label">Schedule</div>
          <span className="task-schedule" style={{ fontSize: 14 }}>{task.schedule}</span>
        </div>

        <div>
          <div className="form-label">Current Status</div>
          <StatusBadge status={task.status} />
        </div>

        <div>
          <div className="form-label">Next Run At</div>
          <div style={{ fontFamily: 'var(--font-mono)', fontSize: 13, color: 'var(--text-secondary)' }}>
            {formatDate(task.next_run_at)}
          </div>
        </div>

        <div>
          <div className="form-label">Last Run At</div>
          <div style={{ fontFamily: 'var(--font-mono)', fontSize: 13, color: 'var(--text-secondary)' }}>
            {formatDate(task.last_run_at)}
          </div>
        </div>

        <div>
          <div className="form-label">Retry Count</div>
          <div style={{ display: 'flex', alignItems: 'baseline', gap: 4 }}>
            <span style={{
              fontSize: 24,
              fontWeight: 800,
              color: task.retry_count > 0 ? 'var(--status-pending)' : 'var(--text-primary)',
            }}>
              {task.retry_count}
            </span>
            <span style={{ fontFamily: 'var(--font-mono)', fontSize: 12, color: 'var(--text-muted)' }}>
              / {task.max_retries} max
            </span>
          </div>
          {task.retry_count > 0 && (
            <div style={{ marginTop: 4 }}>
              <div style={{
                height: 4,
                borderRadius: 2,
                background: 'var(--border)',
                overflow: 'hidden',
                width: 80,
              }}>
                <div style={{
                  height: '100%',
                  width: `${(task.retry_count / task.max_retries) * 100}%`,
                  background: task.retry_count >= task.max_retries
                    ? 'var(--status-failed)'
                    : 'var(--status-pending)',
                  borderRadius: 2,
                  transition: 'width 0.3s ease',
                }} />
              </div>
            </div>
          )}
        </div>

        <div>
          <div className="form-label">Executions</div>
          <div style={{ fontSize: 24, fontWeight: 800, color: 'var(--accent-cyan)' }}>
            {task.executions?.length ?? 0}
          </div>
        </div>

        <div>
          <div className="form-label">Created At</div>
          <div style={{ fontFamily: 'var(--font-mono)', fontSize: 13, color: 'var(--text-secondary)' }}>
            {formatDate(task.created_at)}
          </div>
        </div>

        <div>
          <div className="form-label">Last Updated</div>
          <div style={{ fontFamily: 'var(--font-mono)', fontSize: 13, color: 'var(--text-secondary)' }}>
            {formatDate(task.updated_at)}
          </div>
        </div>
      </div>

      {/* Tabs */}
      <div style={{ display: 'flex', gap: 4, marginBottom: 20 }}>
        {(['executions', 'logs'] as Tab[]).map((tab) => (
          <button
            key={tab}
            onClick={() => setActiveTab(tab)}
            style={{
              padding: '8px 18px',
              borderRadius: 'var(--radius-sm)',
              border: '1px solid',
              borderColor: activeTab === tab ? 'rgba(0,212,255,0.3)' : 'var(--border)',
              background: activeTab === tab ? 'var(--accent-cyan-dim)' : 'transparent',
              color: activeTab === tab ? 'var(--accent-cyan)' : 'var(--text-secondary)',
              fontFamily: 'var(--font-mono)',
              fontSize: '12px',
              fontWeight: 600,
              letterSpacing: '0.08em',
              textTransform: 'uppercase',
              cursor: 'pointer',
              transition: 'all var(--transition)',
            }}
          >
            {tab === 'executions'
              ? `Executions (${task.executions?.length ?? 0})`
              : `Logs (${allLogs.length})`}
          </button>
        ))}
      </div>

      {/* Executions Tab */}
      {activeTab === 'executions' && (
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
                      <div className="empty-desc">// waiting for scheduler to pick up this task</div>
                    </div>
                  </td>
                </tr>
              ) : (
                [...task.executions].reverse().map((exec: TaskExecution) => (
                  <tr key={exec.id}>
                    <td>
                      <span style={{ fontFamily: 'var(--font-mono)', fontSize: 11, color: 'var(--text-muted)' }}>
                        {exec.id.slice(0, 8)}...
                      </span>
                    </td>
                    <td><StatusBadge status={exec.status} /></td>
                    <td><span className="task-time">{formatDate(exec.started_at)}</span></td>
                    <td><span className="task-time">{formatDate(exec.finished_at)}</span></td>
                    <td>
                      <span style={{ fontFamily: 'var(--font-mono)', fontSize: 12, color: 'var(--text-secondary)' }}>
                        {duration(exec.started_at, exec.finished_at)}
                      </span>
                    </td>
                    <td>
                      {exec.error_message ? (
                        <span style={{ fontFamily: 'var(--font-mono)', fontSize: 11, color: 'var(--status-failed)' }}>
                          {exec.error_message}
                        </span>
                      ) : (
                        <span style={{ color: 'var(--text-muted)', fontSize: 12 }}>—</span>
                      )}
                    </td>
                  </tr>
                ))
              )}
            </tbody>
          </table>
        </div>
      )}

      {/* Logs Tab */}
      {activeTab === 'logs' && (
        <div className="table-container">
          {allLogs.length === 0 ? (
            <div className="empty-state">
              <span className="empty-icon">◈</span>
              <div className="empty-title">No logs yet</div>
              <div className="empty-desc">// logs will appear here once the task executes</div>
            </div>
          ) : (
            <div style={{ padding: '8px 0' }}>
              {allLogs.map((log: TaskLog) => (
                <div
                  key={log.id}
                  style={{
                    display: 'grid',
                    gridTemplateColumns: '160px 64px 1fr',
                    alignItems: 'start',
                    gap: 16,
                    padding: '10px 20px',
                    borderBottom: '1px solid var(--border)',
                    transition: 'background var(--transition)',
                  }}
                  onMouseEnter={(e) => (e.currentTarget.style.background = 'var(--bg-card-hover)')}
                  onMouseLeave={(e) => (e.currentTarget.style.background = 'transparent')}
                >
                  <span style={{ fontFamily: 'var(--font-mono)', fontSize: 11, color: 'var(--text-muted)', paddingTop: 2 }}>
                    {formatDate(log.created_at)}
                  </span>
                  <LogLevelBadge level={log.level} />
                  <span style={{ fontFamily: 'var(--font-mono)', fontSize: 12, color: 'var(--text-primary)', lineHeight: 1.6 }}>
                    {log.message}
                  </span>
                </div>
              ))}
            </div>
          )}
        </div>
      )}
    </>
  )
}