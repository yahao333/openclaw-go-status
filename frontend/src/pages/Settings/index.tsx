import { useState } from 'react'
import { Layout } from '@/components/Layout'
import { useAppStore } from '@/store/useAppStore'

export function Settings() {
  const { autoRefresh, setAutoRefresh, refreshInterval, setRefreshInterval } = useAppStore()
  const [intervalInput, setIntervalInput] = useState(refreshInterval.toString())
  const [saved, setSaved] = useState(false)

  const handleIntervalChange = (value: string) => {
    setIntervalInput(value)
    const num = parseInt(value, 10)
    if (!isNaN(num) && num >= 5 && num <= 300) {
      setRefreshInterval(num)
      setSaved(true)
      setTimeout(() => setSaved(false), 2000)
    }
  }

  return (
    <Layout title="系统设置">
      <div className="max-w-2xl space-y-6">
        {/* 自动刷新设置 */}
        <div className="bg-white rounded-lg shadow-sm p-6">
          <h3 className="text-lg font-semibold mb-4">自动刷新设置</h3>

          <div className="space-y-4">
            {/* 开关控制 */}
            <div className="flex items-center justify-between">
              <div>
                <p className="font-medium">启用自动刷新</p>
                <p className="text-sm text-gray-500">开启后页面将自动刷新数据</p>
              </div>
              <button
                onClick={() => setAutoRefresh(!autoRefresh)}
                className={`relative inline-flex h-6 w-11 items-center rounded-full transition-colors ${
                  autoRefresh ? 'bg-primary-600' : 'bg-gray-200'
                }`}
              >
                <span
                  className={`inline-block h-4 w-4 transform rounded-full bg-white transition-transform ${
                    autoRefresh ? 'translate-x-6' : 'translate-x-1'
                  }`}
                />
              </button>
            </div>

            {/* 刷新间隔 */}
            <div className="flex items-center justify-between">
              <div>
                <p className="font-medium">刷新间隔（秒）</p>
                <p className="text-sm text-gray-500">范围：5-300 秒</p>
              </div>
              <input
                type="number"
                min="5"
                max="300"
                value={intervalInput}
                onChange={(e) => handleIntervalChange(e.target.value)}
                disabled={!autoRefresh}
                className="w-24 px-3 py-2 border rounded-lg disabled:bg-gray-100 disabled:cursor-not-allowed focus:ring-2 focus:ring-primary-500 focus:border-primary-500"
              />
            </div>
          </div>

          {saved && (
            <p className="mt-4 text-sm text-green-600">设置已保存</p>
          )}
        </div>

        {/* 当前设置预览 */}
        <div className="bg-gray-50 rounded-lg p-4 text-sm text-gray-600">
          <p>当前设置：自动刷新 {autoRefresh ? '已启用' : '已禁用'}</p>
          {autoRefresh && <p>刷新间隔：{refreshInterval} 秒</p>}
        </div>
      </div>
    </Layout>
  )
}
