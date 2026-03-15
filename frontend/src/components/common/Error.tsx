interface ErrorProps {
  message?: string
  onRetry?: () => void
}

export function Error({ message = '加载失败', onRetry }: ErrorProps) {
  return (
    <div className="flex flex-col items-center justify-center h-64">
      <div className="text-red-500 text-lg mb-2">⚠️</div>
      <p className="text-gray-600 mb-4">{message}</p>
      {onRetry && (
        <button
          onClick={onRetry}
          className="px-4 py-2 bg-primary-600 text-white rounded hover:bg-primary-700 transition-colors"
        >
          重试
        </button>
      )}
    </div>
  )
}
