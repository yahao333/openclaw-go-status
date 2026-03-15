import { useEffect, useState } from 'react'
import { Layout } from '@/components/Layout'
import { Loading, Error, Badge, getStateBadgeVariant } from '@/components/common'
import { taskApi } from '@/api'
import { useAppStore } from '@/store/useAppStore'
import type { Task, TaskStatus } from '@/types'

export function Tasks() {
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<Error | null>(null)
  const [tasks, setTasks] = useState<Task[]>([])
  const [filter, setFilter] = useState<TaskStatus | ''>('')
  const { setHealthStatus, setLastUpdate } = useAppStore()

  const fetchData = async () => {
    setLoading(true)
    setError(null)
    try {
      const res = await taskApi.list()
      setTasks(res?.data || [])
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
  }, [])

  const filteredTasks = filter
    ? tasks.filter(t => t.status === filter)
    : tasks

  if (loading && !tasks.length) {
    return (
      <Layout title="任务列表" onRefresh={fetchData} loading={loading}>
        <Loading />
      </Layout>
    )
  }

  if (error) {
    return (
      <Layout title="任务列表" onRefresh={fetchData}>
        <Error message={error.message} onRetry={fetchData} />
      </Layout>
    )
  }

  return (
    <Layout title="任务列表" onRefresh={fetchData} loading={loading}>
      <div className="space-y-4 animate-fade-in">
        {/* 筛选 */}
        <div className="flex items-center space-x-4">
          <select
            value={filter}
            onChange={(e) => setFilter(e.target.value as TaskStatus | '')}
            className="px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-primary-500 focus:border-transparent"
          >
            <option value="">全部状态</option>
            <option value="pending">待处理</option>
            <option value="in_progress">进行中</option>
            <option value="completed">已完成</option>
            <option value="failed">失败</option>
          </select>
          <span className="text-sm text-gray-500">
            共 {filteredTasks.length} 条记录
          </span>
        </div>

        {/* 表格 */}
        <div className="bg-white rounded-lg shadow-sm overflow-hidden">
          <div className="overflow-x-auto">
            <table className="w-full">
              <thead className="bg-gray-50">
                <tr className="text-left text-sm text-gray-500">
                  <th className="px-6 py-3">任务ID</th>
                  <th className="px-6 py-3">标题</th>
                  <th className="px-6 py-3">负责人</th>
                  <th className="px-6 py-3">状态</th>
                  <th className="px-6 py-3">更新时间</th>
                </tr>
              </thead>
              <tbody className="divide-y divide-gray-200">
                {filteredTasks.map((task) => (
                  <tr key={task.taskID} className="hover:bg-gray-50">
                    <td className="px-6 py-4 text-sm font-mono text-primary-600">
                      {task.taskID}
                    </td>
                    <td className="px-6 py-4 text-sm">{task.title}</td>
                    <td className="px-6 py-4 text-sm">{task.owner}</td>
                    <td className="px-6 py-4">
                      <Badge variant={getStateBadgeVariant(task.status)}>
                        {task.status === 'in_progress' ? '进行中' :
                         task.status === 'completed' ? '已完成' :
                         task.status === 'failed' ? '失败' : '待处理'}
                      </Badge>
                    </td>
                    <td className="px-6 py-4 text-sm text-gray-500">{task.updatedAt}</td>
                  </tr>
                ))}
                {filteredTasks.length === 0 && (
                  <tr>
                    <td colSpan={5} className="px-6 py-8 text-center text-gray-500">
                      暂无任务数据
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
