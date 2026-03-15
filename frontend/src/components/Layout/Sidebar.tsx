import { NavLink } from 'react-router-dom'
import { useAppStore } from '@/store/useAppStore'

const navItems = [
  { path: '/', icon: '📊', label: '总览' },
  { path: '/sessions', icon: '💬', label: '会话' },
  { path: '/tasks', icon: '📋', label: '任务' },
  { path: '/projects', icon: '📁', label: '项目' },
  { path: '/usage', icon: '📈', label: '用量' },
]

export function Sidebar() {
  const { sidebarCollapsed, toggleSidebar, healthStatus, healthMessage } = useAppStore()

  return (
    <aside className={`bg-gray-900 text-white transition-all duration-300 ${sidebarCollapsed ? 'w-16' : 'w-64'}`}>
      <div className="p-4">
        <div className="flex items-center justify-between mb-6">
          {!sidebarCollapsed && (
            <div>
              <h1 className="text-xl font-bold">🦀 OpenClaw</h1>
              <p className="text-xs text-gray-400">状态监控</p>
            </div>
          )}
          <button
            onClick={toggleSidebar}
            className="p-2 hover:bg-gray-800 rounded transition-colors"
          >
            {sidebarCollapsed ? '→' : '←'}
          </button>
        </div>

        <nav className="space-y-1">
          {navItems.map((item) => (
            <NavLink
              key={item.path}
              to={item.path}
              className={({ isActive }) =>
                `flex items-center space-x-3 px-3 py-2 rounded transition-colors ${
                  isActive
                    ? 'bg-primary-600 text-white'
                    : 'text-gray-300 hover:bg-gray-800'
                }`
              }
            >
              <span>{item.icon}</span>
              {!sidebarCollapsed && <span>{item.label}</span>}
            </NavLink>
          ))}
        </nav>
      </div>

      <div className="absolute bottom-0 left-0 right-0 p-4 border-t border-gray-800">
        <div className="flex items-center space-x-2">
          <span className={`w-2 h-2 rounded-full ${healthStatus === 'healthy' ? 'bg-green-500' : 'bg-red-500'}`}></span>
          {!sidebarCollapsed && (
            <span className="text-xs text-gray-400 truncate">{healthMessage}</span>
          )}
        </div>
      </div>
    </aside>
  )
}
