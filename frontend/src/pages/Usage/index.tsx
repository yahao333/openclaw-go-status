import { useEffect, useState } from 'react'
import { Layout } from '@/components/Layout'
import { Loading, Error } from '@/components/common'
import { usageApi } from '@/api'
import { useAppStore } from '@/store/useAppStore'
import type { Usage } from '@/types'

export function Usage() {
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<Error | null>(null)
  const [usage, setUsage] = useState<Usage | null>(null)
  const { setHealthStatus, setLastUpdate, autoRefresh, refreshInterval } = useAppStore()

  const fetchData = async () => {
    setLoading(true)
    setError(null)
    try {
      const res = await usageApi.get()
      setUsage(res?.data || null)
      setLastUpdate(new Date().toLocaleString('zh-CN'))
      setHealthStatus('healthy', 'Gateway 连接正常')
    } catch (err) {
      setError(err as Error)
      setHealthStatus('unhealthy', '数据加载失败')
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    fetchData()

    if (autoRefresh && refreshInterval > 0) {
      const timer = setInterval(fetchData, refreshInterval * 1000)
      return () => clearInterval(timer)
    }
  }, [autoRefresh, refreshInterval])

  if (loading && !usage) {
    return (
      <Layout title="用量统计" onRefresh={fetchData} loading={loading}>
        <Loading />
      </Layout>
    )
  }

  if (error) {
    return (
      <Layout title="用量统计" onRefresh={fetchData}>
        <Error message={error.message} onRetry={fetchData} />
      </Layout>
    )
  }

  if (!usage) {
    return (
      <Layout title="用量统计" onRefresh={fetchData}>
        <div className="text-center py-8 text-gray-500">暂无用量数据</div>
      </Layout>
    )
  }

  const today = usage.today
  const totalTokens = today.totalTokens || (today.tokensIn + today.tokensOut)

  return (
    <Layout title="用量统计" onRefresh={fetchData} loading={loading}>
      <div className="space-y-6 animate-fade-in">
        <div className="bg-white rounded-lg shadow-sm p-6">
          <h3 className="text-lg font-semibold mb-6">今日用量</h3>
          <div className="grid grid-cols-2 md:grid-cols-4 gap-6">
            <div className="text-center">
              <p className="text-sm text-gray-500 mb-1">输入 Token</p>
              <p className="text-2xl font-semibold text-primary-600">
                {today.tokensIn.toLocaleString()}
              </p>
            </div>
            <div className="text-center">
              <p className="text-sm text-gray-500 mb-1">输出 Token</p>
              <p className="text-2xl font-semibold text-primary-600">
                {today.tokensOut.toLocaleString()}
              </p>
            </div>
            <div className="text-center">
              <p className="text-sm text-gray-500 mb-1">总 Token</p>
              <p className="text-2xl font-semibold text-primary-600">
                {totalTokens.toLocaleString()}
              </p>
            </div>
            <div className="text-center">
              <p className="text-sm text-gray-500 mb-1">费用</p>
              <p className="text-2xl font-semibold text-green-600">
                ${today.cost.toFixed(2)}
              </p>
            </div>
          </div>
        </div>

        {/* Token 使用占比 */}
        <div className="bg-white rounded-lg shadow-sm p-6">
          <h3 className="text-lg font-semibold mb-6">使用占比</h3>
          <div className="relative pt-1">
            <div className="flex mb-2 items-center justify-between">
              <div>
                <span className="text-xs font-semibold inline-block py-1 px-2 uppercase rounded-full text-primary-600 bg-primary-200">
                  输入
                </span>
              </div>
              <div className="text-right">
                <span className="text-xs font-semibold inline-block text-primary-600">
                  {totalTokens > 0 ? ((today.tokensIn / totalTokens) * 100).toFixed(1) : 0}%
                </span>
              </div>
            </div>
            <div className="overflow-hidden h-2 mb-4 text-xs flex rounded bg-primary-200">
              <div
                style={{ width: `${totalTokens > 0 ? (today.tokensIn / totalTokens) * 100 : 0}%` }}
                className="shadow-none flex flex-col text-center whitespace-nowrap text-white justify-center bg-primary-500 transition-all duration-500"
              ></div>
            </div>
          </div>
          <div className="relative pt-1">
            <div className="flex mb-2 items-center justify-between">
              <div>
                <span className="text-xs font-semibold inline-block py-1 px-2 uppercase rounded-full text-green-600 bg-green-200">
                  输出
                </span>
              </div>
              <div className="text-right">
                <span className="text-xs font-semibold inline-block text-green-600">
                  {totalTokens > 0 ? ((today.tokensOut / totalTokens) * 100).toFixed(1) : 0}%
                </span>
              </div>
            </div>
            <div className="overflow-hidden h-2 text-xs flex rounded bg-green-200">
              <div
                style={{ width: `${totalTokens > 0 ? (today.tokensOut / totalTokens) * 100 : 0}%` }}
                className="shadow-none flex flex-col text-center whitespace-nowrap text-white justify-center bg-green-500 transition-all duration-500"
              ></div>
            </div>
          </div>
        </div>
      </div>
    </Layout>
  )
}
