import axios from 'axios'
import type { Task, Project, Usage } from '@/types'

const api = axios.create({
  baseURL: '/api',
  timeout: 30000,
})

// 响应拦截器
api.interceptors.response.use(
  (response) => response.data,
  (error) => {
    console.error('API Error:', error)
    return Promise.reject(error)
  }
)

// Dashboard API - 统一的统计数据端点
export const dashboardApi = {
  stats: () => api.get<{
    sessions: number
    running: number
    tasks: number
    projects: number
    data: any[]
  }>('/dashboard').then(res => res.data),
}

// Session API
export const sessionApi = {
  list: () => api.get<{ ok: boolean; data: any[] }>('/sessions').then(res => res.data),
  get: (id: string) => api.get<{ ok: boolean; data: any }>(`/sessions/${id}`).then(res => res.data),
  status: () => api.get<{ ok: boolean; data: any[] }>('/status').then(res => res.data),
}

// Task API
export const taskApi = {
  list: () => api.get<{ ok: boolean; data: Task[] }>('/tasks').then(res => res.data),
}

// Project API
export const projectApi = {
  list: () => api.get<{ ok: boolean; data: Project[] }>('/projects').then(res => res.data),
}

// Usage API
export const usageApi = {
  get: () => api.get<{ ok: boolean; data: Usage }>('/usage').then(res => res.data),
}

// Health API
export const healthApi = {
  check: () => api.get<{ ok: boolean; status: string; message: string }>('/health').then(res => res.data),
}

// Cron API
export interface CronJob {
  jobId: string
  name?: string
  enabled: boolean
  nextRunAt?: string
  health: string
}

export const cronApi = {
  list: () => api.get<{ ok: boolean; data: CronJob[] }>('/cron').then(res => res.data),
}

export default api
