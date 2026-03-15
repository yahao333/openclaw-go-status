import { useEffect, useState } from 'react'
import { Layout } from '@/components/Layout'
import { StatCard, Loading, Error } from '@/components/common'
import { dashboardApi, taskApi, projectApi, usageApi, healthApi } from '@/api'
import { useAppStore } from '@/store/useAppStore'
import type { Session, Task, Project, Usage, HealthStatus } from '@/types'

export function Dashboard() {
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<Error | null>(null)
  const [sessions, setSessions] = useState<Session[]>([])
  const [tasks, setTasks] = useState<Task[]>([])
  const [projects, setProjects] = useState<Project[]>([])
  const [usage, setUsage] = useState<Usage | null>(null)
  const [stats, setStats] = useState({ sessions: 0, running: 0, tasks: 0, projects: 0 })
  const { setHealthStatus, setLastUpdate } = useAppStore()

  const fetchData = async () => {
    setLoading(true)
    setError(null)
    try {
      const [dashboardRes, tasksRes, projectsRes, usageRes, healthRes] = await Promise.all([
        dashboardApi.stats(),
        taskApi.list(),
        projectApi.list(),
        usageApi.get(),
        healthApi.check(),
      ])

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
  }, [])

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
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
          <StatCard icon="💬" label="活跃会话" value={stats.sessions} />
          <StatCard icon="▶️" label="运行中" value={stats.running} />
          <StatCard icon="📋" label="任务总数" value={stats.tasks} />
          <StatCard icon="📁" label="项目总数" value={stats.projects} />
        </div>

        {/* 最近会话 */}
        <div className="bg-white rounded-lg shadow-sm p-6">
          <h3 className="text-lg font-semibold mb-4">最近会话</h3>
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
                {sessions.slice(0, 5).map((session) => (
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

        {/* 进行中的任务 */}
        <div className="bg-white rounded-lg shadow-sm p-6">
          <h3 className="text-lg font-semibold mb-4">进行中的任务</h3>
          <div className="overflow-x-auto">
            <table className="w-full">
              <thead>
                <tr className="text-left text-sm text-gray-500 border-b">
                  <th className="pb-3">任务ID</th>
                  <th className="pb-3">标题</th>
                  <th className="pb-3">负责人</th>
                  <th className="pb-3">状态</th>
                  <th className="pb-3">更新时间</th>
                </tr>
              </thead>
              <tbody>
                {tasks.slice(0, 5).map((task) => (
                  <tr key={task.taskID} className="border-b last:border-0 hover:bg-gray-50">
                    <td className="py-3 text-sm">{task.taskID.slice(0, 12)}...</td>
                    <td className="py-3 text-sm">{task.title}</td>
                    <td className="py-3 text-sm">{task.owner}</td>
                    <td className="py-3">
                      <span className={`inline-flex items-center px-2 py-1 rounded text-xs font-medium ${
                        task.status === 'in_progress' ? 'bg-blue-100 text-blue-800' :
                        task.status === 'completed' ? 'bg-green-100 text-green-800' :
                        task.status === 'failed' ? 'bg-red-100 text-red-800' :
                        'bg-gray-100 text-gray-800'
                      }`}>
                        {task.status}
                      </span>
                    </td>
                    <td className="py-3 text-sm text-gray-500">{task.updatedAt}</td>
                  </tr>
                ))}
                {tasks.length === 0 && (
                  <tr>
                    <td colSpan={5} className="py-8 text-center text-gray-500">暂无任务数据</td>
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
