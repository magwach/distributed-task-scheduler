"use client";

import { logout } from "@/services/api";
import Link from "next/link";
import { usePathname } from "next/navigation";
import { useState } from "react";
import { toast } from "sonner";
import { useUser } from "./ProtectedRoute";

const navItems = [
  {
    section: "Overview",
    links: [{ href: "/dashboard", label: "Dashboard", icon: "⬡" }],
  },
  {
    section: "Tasks",
    links: [
      { href: "/tasks", label: "All Tasks", icon: "◈" },
      { href: "/tasks/new", label: "Create Task", icon: "+" },
    ],
  },
  {
    section: "System",
    links: [{ href: "/workers", label: "Workers", icon: "⬡" }],
  },
];

export default function Sidebar() {
  const pathname = usePathname();
  const [loggingOut, setLoggingOut] = useState(false);

  const { user } = useUser();

  const handleLogout = async () => {
    setLoggingOut(true);
    try {
      await logout();
      window.location.href = "/login";
    } catch {
      toast.error("Failed to sign out.");
    } finally {
      setLoggingOut(false);
    }
  };

  const initials = user?.email ? user.email.slice(0, 2).toUpperCase() : "??";

  if (!user) {
    return null;
  }

  return (
    <aside className="sidebar">
      <div className="sidebar-logo">
        <div className="logo-mark">
          <div className="logo-icon">TS</div>
          <div>
            <div className="logo-text">Task Scheduler</div>
            <div className="logo-sub">Distributed System</div>
          </div>
        </div>
      </div>

      <nav className="sidebar-nav">
        {navItems.map((section) => (
          <div key={section.section}>
            <div className="nav-section-label">{section.section}</div>
            {section.links.map((link) => (
              <Link
                key={link.href}
                href={link.href}
                className={`nav-link ${
                  pathname === link.href || pathname.startsWith(link.href + "/")
                    ? "active"
                    : ""
                }`}
              >
                <span
                  style={{ fontFamily: "var(--font-mono)", fontSize: "14px" }}
                >
                  {link.icon}
                </span>
                {link.label}
              </Link>
            ))}
          </div>
        ))}
      </nav>

      <div className="sidebar-footer">
        {/* User info */}
        {user && (
          <div
            style={{
              display: "flex",
              alignItems: "center",
              gap: 10,
              padding: "12px 0",
              marginBottom: 12,
              borderBottom: "1px solid var(--border)",
            }}
          >
            <div
              style={{
                width: 32,
                height: 32,
                borderRadius: "50%",
                background:
                  "linear-gradient(135deg, var(--accent-cyan), var(--accent-purple))",
                display: "flex",
                alignItems: "center",
                justifyContent: "center",
                fontSize: 11,
                fontWeight: 700,
                color: "#fff",
                flexShrink: 0,
              }}
            >
              {initials}
            </div>
            <div style={{ overflow: "hidden", flex: 1 }}>
              <div
                style={{
                  fontSize: 12,
                  fontWeight: 600,
                  color: "var(--text-primary)",
                  whiteSpace: "nowrap",
                  overflow: "hidden",
                  textOverflow: "ellipsis",
                }}
              >
                {user.email}
              </div>
              <div
                style={{
                  fontSize: 10,
                  fontFamily: "var(--font-mono)",
                  color: "var(--text-muted)",
                  textTransform: "uppercase",
                  letterSpacing: "0.08em",
                }}
              >
                {user.role}
              </div>
            </div>
          </div>
        )}

        {/* System status */}
        <div className="system-status" style={{ marginBottom: 12 }}>
          <div className="status-dot" />
          System Online
        </div>

        {/* Logout */}
        <button
          onClick={handleLogout}
          disabled={loggingOut}
          style={{
            display: "flex",
            alignItems: "center",
            gap: 8,
            width: "100%",
            padding: "8px 0",
            background: "none",
            border: "none",
            color: "var(--text-muted)",
            fontFamily: "var(--font-mono)",
            fontSize: 11,
            letterSpacing: "0.08em",
            cursor: "pointer",
            transition: "color var(--transition)",
          }}
          onMouseEnter={(e) =>
            (e.currentTarget.style.color = "var(--status-failed)")
          }
          onMouseLeave={(e) =>
            (e.currentTarget.style.color = "var(--text-muted)")
          }
        >
          <span style={{ fontSize: 13 }}>→</span>
          {loggingOut ? "Signing out..." : "Sign Out"}
        </button>
      </div>
    </aside>
  );
}
