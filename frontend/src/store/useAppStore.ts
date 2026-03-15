import { create } from 'zustand'
import { persist } from 'zustand/middleware'
import type { HealthStatus } from '@/types'

interface AppState {
  // UI 状态
  sidebarCollapsed: boolean
  toggleSidebar: () => void

  // 健康状态
  healthStatus: HealthStatus
  healthMessage: string
  setHealthStatus: (status: HealthStatus, message: string) => void

  // 刷新控制
  lastUpdate: string
  setLastUpdate: (time: string) => void
  autoRefresh: boolean
  setAutoRefresh: (enabled: boolean) => void

  // 主题
  theme: 'light' | 'dark'
  setTheme: (theme: 'light' | 'dark') => void
}

export const useAppStore = create<AppState>()(
  persist(
    (set) => ({
      // UI 状态
      sidebarCollapsed: false,
      toggleSidebar: () => set((state) => ({ sidebarCollapsed: !state.sidebarCollapsed })),

      // 健康状态
      healthStatus: 'healthy',
      healthMessage: 'Gateway 连接正常',
      setHealthStatus: (status, message) => set({ healthStatus: status, healthMessage: message }),

      // 刷新控制
      lastUpdate: '',
      setLastUpdate: (time) => set({ lastUpdate: time }),
      autoRefresh: false,
      setAutoRefresh: (enabled) => set({ autoRefresh: enabled }),

      // 主题
      theme: 'light',
      setTheme: (theme) => set({ theme }),
    }),
    {
      name: 'openclaw-storage',
    }
  )
)
