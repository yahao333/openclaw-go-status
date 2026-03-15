import { useState, useEffect, useCallback } from 'react'

interface UseFetchOptions<T> {
  immediate?: boolean
  refreshInterval?: number
}

export function useFetch<T>(
  fetchFn: () => Promise<T>,
  options: UseFetchOptions<T> = {}
) {
  const [data, setData] = useState<T | null>(null)
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState<Error | null>(null)

  const execute = useCallback(async () => {
    setLoading(true)
    setError(null)
    try {
      const result = await fetchFn()
      setData(result)
    } catch (err) {
      setError(err as Error)
    } finally {
      setLoading(false)
    }
  }, [fetchFn])

  useEffect(() => {
    if (options.immediate !== false) {
      execute()
    }

    if (options.refreshInterval) {
      const interval = setInterval(execute, options.refreshInterval)
      return () => clearInterval(interval)
    }
  }, [execute, options.immediate, options.refreshInterval])

  return { data, loading, error, refetch: execute }
}
