import { useAppStore } from '@/store/useAppStore'

interface HeaderProps {
  title: string
  onRefresh?: () => void
  loading?: boolean
}

export function Header({ title, onRefresh, loading }: HeaderProps) {
  const { lastUpdate } = useAppStore()

  return (
    <header className="bg-white shadow-sm px-6 py-4 flex items-center justify-between">
      <h2 className="text-xl font-semibold text-gray-800">{title}</h2>
      <div className="flex items-center space-x-4">
        {lastUpdate && (
          <span className="text-sm text-gray-500">最后更新: {lastUpdate}</span>
        )}
        {onRefresh && (
          <button
            onClick={onRefresh}
            disabled={loading}
            className="px-4 py-2 bg-gray-100 hover:bg-gray-200 rounded transition-colors disabled:opacity-50"
          >
            {loading ? '刷新中...' : '🔄 刷新'}
          </button>
        )}
      </div>
    </header>
  )
}
