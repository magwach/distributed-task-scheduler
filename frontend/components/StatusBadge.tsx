type Status = 'pending' | 'running' | 'success' | 'failed'

interface StatusBadgeProps {
  status: Status | string
}

export default function StatusBadge({ status }: StatusBadgeProps) {
  const normalized = status.toLowerCase() as Status
  return (
    <span className={`status-badge ${normalized}`}>
      {normalized}
    </span>
  )
}