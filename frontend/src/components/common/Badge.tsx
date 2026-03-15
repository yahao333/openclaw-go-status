import type { ReactNode } from 'react'

interface BadgeProps {
  variant?: 'default' | 'success' | 'warning' | 'error' | 'info'
  children: ReactNode
}

const variantClasses = {
  default: 'bg-gray-100 text-gray-800',
  success: 'bg-green-100 text-green-800',
  warning: 'bg-yellow-100 text-yellow-800',
  error: 'bg-red-100 text-red-800',
  info: 'bg-blue-100 text-blue-800',
}

export function Badge({ variant = 'default', children }: BadgeProps) {
  return (
    <span className={`inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium ${variantClasses[variant]}`}>
      {children}
    </span>
  )
}

// 根据状态获取 Badge 变种
export function getStateBadgeVariant(state: string): BadgeProps['variant'] {
  switch (state) {
    case 'running':
    case 'completed':
    case 'healthy':
      return 'success'
    case 'idle':
    case 'pending':
    case 'in_progress':
      return 'info'
    case 'blocked':
    case 'waiting_approval':
      return 'warning'
    case 'error':
    case 'failed':
    case 'unhealthy':
      return 'error'
    default:
      return 'default'
  }
}
