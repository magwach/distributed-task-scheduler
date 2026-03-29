"use client";

import Link from "next/link";
import { usePathname } from "next/navigation";
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

  const { user } = useUser();


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
                className={`nav-link ${pathname === link.href || pathname.startsWith(link.href + "/") ? "active" : ""}`}
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
        <div className="system-status">
          <div className="status-dot" />
          System Online
        </div>
      </div>
    </aside>
  );
}
