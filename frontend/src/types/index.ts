// Session 类型
export interface Session {
  sessionKey: string
  label: string
  agentId: string
  state: SessionState
  tokensIn: number
  tokensOut: number
  cost: number
  lastMessageAt: string
}

export type SessionState = 'idle' | 'running' | 'blocked' | 'waiting_approval' | 'error'

// Task 类型
export interface Task {
  taskID: string
  title: string
  owner: string
  status: TaskStatus
  updatedAt: string
}

export type TaskStatus = 'pending' | 'in_progress' | 'completed' | 'failed'

// Project 类型
export interface Project {
  projectID: string
  name: string
  description: string
  createdAt: string
  sessionCount: number
}

// Usage 类型
export interface Usage {
  today: UsageSnapshot
  week7: UsageSnapshot[]
  month30: UsageSnapshot[] | null
  total: UsageSnapshot
}

export interface UsageSnapshot {
  date: string
  tokensIn: number
  tokensOut: number
  totalTokens: number
  cost: number
}

// Dashboard 统计数据
export interface DashboardStats {
  sessions: number
  running: number
  tasks: number
  projects: number
  data: Session[]
}

// API 响应类型
export interface ApiResponse<T> {
  ok: boolean
  data?: T
  error?: string
}

// 健康状态
export type HealthStatus = 'healthy' | 'degraded' | 'unhealthy'
