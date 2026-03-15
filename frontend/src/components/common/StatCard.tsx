interface StatCardProps {
  icon: string
  label: string
  value: number | string
  subLabel?: string
}

export function StatCard({ icon, label, value, subLabel }: StatCardProps) {
  return (
    <div className="bg-white rounded-lg shadow-sm p-6 flex items-center space-x-4">
      <div className="text-3xl">{icon}</div>
      <div>
        <p className="text-sm text-gray-500">{label}</p>
        <p className="text-2xl font-semibold text-gray-900">{value}</p>
        {subLabel && <p className="text-xs text-gray-400">{subLabel}</p>}
      </div>
    </div>
  )
}
