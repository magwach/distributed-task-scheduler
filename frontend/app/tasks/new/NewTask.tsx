"use client";

import { useState } from "react";
import { useRouter } from "next/navigation";
import Link from "next/link";
import { createTask } from "@/services/api";
import { toast } from "sonner";

const CRON_EXAMPLES = [
  { label: "Every minute", value: "* * * * *" },
  { label: "Every 5 minutes", value: "*/5 * * * *" },
  { label: "Every hour", value: "0 * * * *" },
  { label: "Daily at 8AM", value: "0 8 * * *" },
  { label: "Every Monday", value: "0 9 * * 1" },
];

export default function CreateTaskPage() {
  const router = useRouter();
  const [form, setForm] = useState({ name: "", description: "", schedule: "" });
  const [errors, setErrors] = useState<Record<string, string>>({});
  const [submitting, setSubmitting] = useState(false);

  const validate = () => {
    const e: Record<string, string> = {};

    if (!form.name.trim()) e.name = "Task name is required";

    if (!form.schedule.trim()) {
      e.schedule = "Schedule is required";
    } else if (!isValidCron(form.schedule)) {
      e.schedule = "Invalid cron format";
    }

    return e;
  };

  function isValidCron(cron: string): boolean {
    const parts = cron.trim().split(/\s+/);

    // Must have exactly 5 parts
    if (parts.length !== 5) return false;

    const cronPart = /^(\*|\d+|\d+-\d+|\*\/\d+|\d+(,\d+)*)$/; // basic support

    return parts.every((part) => cronPart.test(part));
  }

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    const errs = validate();
    if (Object.keys(errs).length > 0) {
      setErrors(errs);
      return;
    }
    setSubmitting(true);
    try {
      await createTask({
        title: form.name.trim(),
        description: form.description.trim() || undefined,
        schedule: form.schedule.trim(),
      });
      toast.success("Task created successfully.");
      router.push("/dashboard");
    } catch (err: unknown) {
      toast.error(
        err instanceof Error ? err.message : "Failed to create task.",
      );
    } finally {
      setSubmitting(false);
    }
  };

  const handleChange = (field: string, value: string) => {
    setForm((prev) => ({ ...prev, [field]: value }));
    if (errors[field]) setErrors((prev) => ({ ...prev, [field]: "" }));
  };

  return (
    <>
      <div className="page-header">
        <div>
          <div className="breadcrumb">
            <Link href="/dashboard">Dashboard</Link>
            <span className="breadcrumb-sep">/</span>
            <span>New Task</span>
          </div>
          <h1 className="page-title">
            Create <span>Task</span>
          </h1>
          <p className="page-subtitle">// define a new scheduled task</p>
        </div>
      </div>

      <form onSubmit={handleSubmit}>
        <div className="form-card">
          <div className="form-group">
            <label className="form-label">
              Task Name <span>*</span>
            </label>
            <input
              type="text"
              className="form-input"
              placeholder="e.g. Send daily report"
              value={form.name}
              onChange={(e) => handleChange("name", e.target.value)}
            />
            {errors.name && <div className="form-error">⚠ {errors.name}</div>}
          </div>

          <div className="form-group">
            <label className="form-label">Description</label>
            <textarea
              className="form-textarea"
              placeholder="Optional description of what this task does..."
              value={form.description}
              onChange={(e) => handleChange("description", e.target.value)}
            />
          </div>

          <div className="form-group">
            <label className="form-label">
              Cron Schedule <span>*</span>
            </label>
            <input
              type="text"
              className="form-input"
              placeholder="e.g. */5 * * * *"
              value={form.schedule}
              onChange={(e) => handleChange("schedule", e.target.value)}
            />
            {errors.schedule && (
              <div className="form-error">⚠ {errors.schedule}</div>
            )}
            <div className="form-hint">
              Format: <code>minute hour day month weekday</code>
            </div>

            <div
              style={{
                marginTop: 14,
                display: "flex",
                flexWrap: "wrap",
                gap: 8,
              }}
            >
              {CRON_EXAMPLES.map((ex) => (
                <button
                  key={ex.value}
                  type="button"
                  onClick={() => handleChange("schedule", ex.value)}
                  style={{
                    background:
                      form.schedule === ex.value
                        ? "var(--accent-cyan-dim)"
                        : "var(--bg-secondary)",
                    border: `1px solid ${form.schedule === ex.value ? "rgba(0,212,255,0.3)" : "var(--border)"}`,
                    borderRadius: "var(--radius-sm)",
                    padding: "6px 12px",
                    fontSize: "11px",
                    fontFamily: "var(--font-mono)",
                    color:
                      form.schedule === ex.value
                        ? "var(--accent-cyan)"
                        : "var(--text-secondary)",
                    cursor: "pointer",
                    transition: "all var(--transition)",
                  }}
                >
                  {ex.label} — <span style={{ opacity: 0.7 }}>{ex.value}</span>
                </button>
              ))}
            </div>
          </div>

          <div className="form-actions">
            <button
              type="submit"
              className="btn btn-primary"
              disabled={submitting}
            >
              {submitting ? "⟳ Creating..." : "+ Create Task"}
            </button>
            <Link href="/dashboard" className="btn btn-secondary">
              Cancel
            </Link>
          </div>
        </div>
      </form>
    </>
  );
}
