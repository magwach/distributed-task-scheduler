import type { Metadata } from 'next'
import './globals.css'
import Sidebar from '@/components/Sidebar'
import { Toaster } from 'sonner'

export const metadata: Metadata = {
  title: 'Task Scheduler',
  description: 'Distributed Task Scheduling System',
}

export default function RootLayout({
  children,
}: {
  children: React.ReactNode
}) {
  return (
    <html lang="en">
      <body>
        <div className="grid-bg" />
        <div className="layout">
          <Sidebar />
          <main className="main-content">
            {children}
          </main>
        </div>
        <Toaster
          position="bottom-right"
          toastOptions={{
            style: {
              background: 'var(--bg-card)',
              border: '1px solid var(--border-bright)',
              color: 'var(--text-primary)',
              fontFamily: 'var(--font-mono)',
              fontSize: '13px',
              borderRadius: 'var(--radius-sm)',
            },
          }}
          theme="dark"
        />
      </body>
    </html>
  )
}