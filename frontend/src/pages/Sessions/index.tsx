import { useEffect, useState } from 'react'
import { Layout } from '@/components/Layout'
import { Loading, Error, getStateBadgeVariant, Badge } from '@/components/common'
import { sessionApi } from '@/api'
import { useAppStore } from '@/store/useAppStore'
import type { Session, SessionState } from '@/types'

export function Sessions() {
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<Error | null>(null)
  const [sessions, setSessions] = useState<Session[]>([])
  const [filter, setFilter] = useState<SessionState | ''>('')
  const { setHealthStatus, setLastUpdate, autoRefresh, refreshInterval } = useAppStore()

  const fetchData = async () => {
    setLoading(true)
    setError(null)
    try {
      const res = await sessionApi.list()
      setSessions(res?.data || [])
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

  const filteredSessions = filter
    ? sessions.filter(s => s.state === filter)
    : sessions

  if (loading && !sessions.length) {
    return (
      <Layout title="会话列表" onRefresh={fetchData} loading={loading}>
        <Loading />
      </Layout>
    )
  }

  if (error) {
    return (
      <Layout title="会话列表" onRefresh={fetchData}>
        <Error message={error.message} onRetry={fetchData} />
      </Layout>
    )
  }

  return (
    <Layout title="会话列表" onRefresh={fetchData} loading={loading}>
      <div className="space-y-4 animate-fade-in">
        {/* 筛选 */}
        <div className="flex items-center space-x-4">
          <select
            value={filter}
            onChange={(e) => setFilter(e.target.value as SessionState | '')}
            className="px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-primary-500 focus:border-transparent"
          >
            <option value="">全部状态</option>
            <option value="running">运行中</option>
            <option value="idle">空闲</option>
            <option value="blocked">阻塞</option>
            <option value="waiting_approval">等待审批</option>
            <option value="error">错误</option>
          </select>
          <span className="text-sm text-gray-500">
            共 {filteredSessions.length} 条记录
          </span>
        </div>

        {/* 表格 */}
        <div className="bg-white rounded-lg shadow-sm overflow-hidden">
          <div className="overflow-x-auto">
            <table className="w-full">
              <thead className="bg-gray-50">
                <tr className="text-left text-sm text-gray-500">
                  <th className="px-6 py-3">会话ID</th>
                  <th className="px-6 py-3">标签</th>
                  <th className="px-6 py-3">Agent ID</th>
                  <th className="px-6 py-3">状态</th>
                  <th className="px-6 py-3">Token 输入</th>
                  <th className="px-6 py-3">Token 输出</th>
                  <th className="px-6 py-3">费用</th>
                  <th className="px-6 py-3">最后消息</th>
                </tr>
              </thead>
              <tbody className="divide-y divide-gray-200">
                {filteredSessions.map((session) => (
                  <tr key={session.sessionKey} className="hover:bg-gray-50">
                    <td className="px-6 py-4 text-sm font-mono text-primary-600">
                      {session.sessionKey}
                    </td>
                    <td className="px-6 py-4 text-sm">{session.label || '-'}</td>
                    <td className="px-6 py-4 text-sm">{session.agentId}</td>
                    <td className="px-6 py-4">
                      <Badge variant={getStateBadgeVariant(session.state)}>
                        {session.state}
                      </Badge>
                    </td>
                    <td className="px-6 py-4 text-sm">{session.tokensIn.toLocaleString()}</td>
                    <td className="px-6 py-4 text-sm">{session.tokensOut.toLocaleString()}</td>
                    <td className="px-6 py-4 text-sm">${session.cost.toFixed(2)}</td>
                    <td className="px-6 py-4 text-sm text-gray-500">{session.lastMessageAt}</td>
                  </tr>
                ))}
                {filteredSessions.length === 0 && (
                  <tr>
                    <td colSpan={8} className="px-6 py-8 text-center text-gray-500">
                      暂无会话数据
                    </td>
                  </tr>
                )}
              </tbody>
            </table>
          </div>
        </div>
      </div>
    </Layout>
  )
}
