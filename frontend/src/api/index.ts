import axios from 'axios'
import type { Session, Task, Project, Usage, DashboardStats } from '@/types'

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

// Session API
export const sessionApi = {
  list: () => api.get<{ ok: boolean; data: Session[] }>('/sessions').then(res => res.data),
  get: (id: string) => api.get<{ ok: boolean; data: Session }>(`/sessions/${id}`).then(res => res.data),
  status: () => api.get<{ ok: boolean; data: DashboardStats }>('/status').then(res => res.data),
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

export default api
