import { useEffect, useState } from 'react'
import { Layout } from '@/components/Layout'
import { StatCard, Loading, Error } from '@/components/common'
import { dashboardApi, taskApi, projectApi, usageApi, healthApi, cronApi, CronJob } from '@/api'
import { useAppStore } from '@/store/useAppStore'
import type { Session, Task, Project, Usage, HealthStatus } from '@/types'

export function Dashboard() {
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<Error | null>(null)
  const [sessions, setSessions] = useState<Session[]>([])
  const [tasks, setTasks] = useState<Task[]>([])
  const [projects, setProjects] = useState<Project[]>([])
  const [usage, setUsage] = useState<Usage | null>(null)
  const [cronJobs, setCronJobs] = useState<CronJob[]>([])
  const [stats, setStats] = useState({ sessions: 0, running: 0, tasks: 0, projects: 0 })
  const { setHealthStatus, setLastUpdate, autoRefresh, refreshInterval, recentSessionsCount } = useAppStore()

  const fetchData = async () => {
    setLoading(true)
    setError(null)
    try {
      // 分别调用每个 API，单独处理错误
      let dashboardRes, tasksRes, projectsRes, usageRes, healthRes, cronRes

      try { dashboardRes = await dashboardApi.stats() } catch(e) { console.error('dashboardApi error:', e) }
      try { tasksRes = await taskApi.list() } catch(e) { console.error('taskApi error:', e) }
      try { projectsRes = await projectApi.list() } catch(e) { console.error('projectApi error:', e) }
      try { usageRes = await usageApi.get() } catch(e) { console.error('usageApi error:', e) }
      try { healthRes = await healthApi.check() } catch(e) { console.error('healthApi error:', e) }
      try { cronRes = await cronApi.list(); console.log('cronRes:', cronRes) } catch(e) { console.error('cronApi error:', e) }

      if (dashboardRes?.sessions !== undefined) {
        setSessions(dashboardRes.data || [])
        setStats({
          sessions: dashboardRes.sessions,
          running: dashboardRes.running,
          tasks: dashboardRes.tasks,
          projects: dashboardRes.projects,
        })
      }

      setTasks(tasksRes?.data || [])
      setProjects(projectsRes?.data || [])
      setUsage(usageRes?.data || null)

      // 确保 cronRes 正确处理
      if (cronRes?.data && Array.isArray(cronRes.data)) {
        setCronJobs(cronRes.data)
      } else {
        setCronJobs([])
      }

      if (healthRes?.ok) {
        setHealthStatus(healthRes.status as HealthStatus, healthRes.message)
      }

      setLastUpdate(new Date().toLocaleString('zh-CN'))
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

  if (loading && !sessions.length) {
    return (
      <Layout title="系统总览" onRefresh={fetchData} loading={loading}>
        <Loading />
      </Layout>
    )
  }

  if (error) {
    return (
      <Layout title="系统总览" onRefresh={fetchData}>
        <Error message={error.message} onRetry={fetchData} />
      </Layout>
    )
  }

  return (
    <Layout title="系统总览" onRefresh={fetchData} loading={loading}>
      <div className="space-y-6 animate-fade-in">
        {/* 统计卡片 */}
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-2 gap-4">
          <StatCard icon="💬" label="活跃会话" value={stats.sessions} />
          <StatCard icon="▶️" label="运行中" value={stats.running} />
        </div>

        {/* 最近会话 */}
        <div className="bg-white rounded-lg shadow-sm p-6">
          <h3 className="text-lg font-semibold mb-4">最近会话 ({Math.min(recentSessionsCount, sessions.length)})</h3>
          <div className="overflow-x-auto">
            <table className="w-full">
              <thead>
                <tr className="text-left text-sm text-gray-500 border-b">
                  <th className="pb-3">会话ID</th>
                  <th className="pb-3">标签</th>
                  <th className="pb-3">Agent</th>
                  <th className="pb-3">状态</th>
                  <th className="pb-3">费用</th>
                  <th className="pb-3">最后消息</th>
                </tr>
              </thead>
              <tbody>
                {sessions.slice(0, recentSessionsCount).map((session) => (
                  <tr key={session.sessionKey} className="border-b last:border-0 hover:bg-gray-50">
                    <td className="py-3 text-sm">{session.sessionKey.slice(0, 12)}...</td>
                    <td className="py-3 text-sm">{session.label || '-'}</td>
                    <td className="py-3 text-sm">{session.agentId}</td>
                    <td className="py-3">
                      <span className={`inline-flex items-center px-2 py-1 rounded text-xs font-medium ${
                        session.state === 'running' ? 'bg-green-100 text-green-800' :
                        session.state === 'error' ? 'bg-red-100 text-red-800' :
                        'bg-gray-100 text-gray-800'
                      }`}>
                        {session.state}
                      </span>
                    </td>
                    <td className="py-3 text-sm">${session.cost.toFixed(2)}</td>
                    <td className="py-3 text-sm text-gray-500">{session.lastMessageAt}</td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        </div>

        {/* Cron 定时任务 */}
        <div className="bg-white rounded-lg shadow-sm p-6">
          <h3 className="text-lg font-semibold mb-4">定时任务 (Cron Jobs)</h3>
          <div className="overflow-x-auto">
            <table className="w-full">
              <thead>
                <tr className="text-left text-sm text-gray-500 border-b">
                  <th className="pb-3">任务ID</th>
                  <th className="pb-3">名称</th>
                  <th className="pb-3">状态</th>
                  <th className="pb-3">健康状态</th>
                  <th className="pb-3">下次运行</th>
                </tr>
              </thead>
              <tbody>
                {cronJobs.slice(0, 5).map((job) => (
                  <tr key={job.jobId} className="border-b last:border-0 hover:bg-gray-50">
                    <td className="py-3 text-sm">{job.jobId.slice(0, 12)}...</td>
                    <td className="py-3 text-sm">{job.name || '-'}</td>
                    <td className="py-3">
                      <span className={`inline-flex items-center px-2 py-1 rounded text-xs font-medium ${
                        job.enabled ? 'bg-green-100 text-green-800' : 'bg-gray-100 text-gray-800'
                      }`}>
                        {job.enabled ? '已启用' : '已禁用'}
                      </span>
                    </td>
                    <td className="py-3">
                      <span className={`inline-flex items-center px-2 py-1 rounded text-xs font-medium ${
                        job.health === 'scheduled' ? 'bg-blue-100 text-blue-800' :
                        job.health === 'due' ? 'bg-yellow-100 text-yellow-800' :
                        job.health === 'late' ? 'bg-red-100 text-red-800' :
                        job.health === 'disabled' ? 'bg-gray-100 text-gray-800' :
                        'bg-gray-100 text-gray-800'
                      }`}>
                        {job.health}
                      </span>
                    </td>
                    <td className="py-3 text-sm text-gray-500">{job.nextRunAt || '-'}</td>
                  </tr>
                ))}
                {cronJobs.length === 0 && (
                  <tr>
                    <td colSpan={5} className="py-8 text-center text-gray-500">暂无定时任务</td>
                  </tr>
                )}
              </tbody>
            </table>
          </div>
        </div>

        {/* 今日用量 */}
        {usage && (
          <div className="bg-white rounded-lg shadow-sm p-6">
            <h3 className="text-lg font-semibold mb-4">今日用量</h3>
            <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
              <div>
                <p className="text-sm text-gray-500">输入 Token</p>
                <p className="text-xl font-semibold">{usage.today.tokensIn.toLocaleString()}</p>
              </div>
              <div>
                <p className="text-sm text-gray-500">输出 Token</p>
                <p className="text-xl font-semibold">{usage.today.tokensOut.toLocaleString()}</p>
              </div>
              <div>
                <p className="text-sm text-gray-500">总 Token</p>
                <p className="text-xl font-semibold">{usage.today.totalTokens.toLocaleString()}</p>
              </div>
              <div>
                <p className="text-sm text-gray-500">费用</p>
                <p className="text-xl font-semibold text-green-600">${usage.today.cost.toFixed(2)}</p>
              </div>
            </div>
          </div>
        )}
      </div>
    </Layout>
  )
}
