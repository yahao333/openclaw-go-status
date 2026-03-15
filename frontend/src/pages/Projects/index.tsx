import { useEffect, useState } from 'react'
import { Layout } from '@/components/Layout'
import { Loading, Error } from '@/components/common'
import { projectApi } from '@/api'
import { useAppStore } from '@/store/useAppStore'
import type { Project } from '@/types'

export function Projects() {
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<Error | null>(null)
  const [projects, setProjects] = useState<Project[]>([])
  const { setHealthStatus, setLastUpdate } = useAppStore()

  const fetchData = async () => {
    setLoading(true)
    setError(null)
    try {
      const res = await projectApi.list()
      setProjects(res?.data || [])
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

  if (loading && !projects.length) {
    return (
      <Layout title="项目列表" onRefresh={fetchData} loading={loading}>
        <Loading />
      </Layout>
    )
  }

  if (error) {
    return (
      <Layout title="项目列表" onRefresh={fetchData}>
        <Error message={error.message} onRetry={fetchData} />
      </Layout>
    )
  }

  return (
    <Layout title="项目列表" onRefresh={fetchData} loading={loading}>
      <div className="space-y-4 animate-fade-in">
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
          {projects.map((project) => (
            <div key={project.projectID} className="bg-white rounded-lg shadow-sm p-6 hover:shadow-md transition-shadow">
              <h3 className="text-lg font-semibold mb-2">{project.name}</h3>
              <p className="text-sm text-gray-600 mb-4 line-clamp-2">
                {project.description || '暂无描述'}
              </p>
              <div className="flex items-center justify-between text-sm text-gray-500">
                <span>会话数: {project.sessionCount}</span>
                <span>{project.createdAt}</span>
              </div>
            </div>
          ))}
          {projects.length === 0 && (
            <div className="col-span-full text-center py-8 text-gray-500">
              暂无项目数据
            </div>
          )}
        </div>
      </div>
    </Layout>
  )
}
