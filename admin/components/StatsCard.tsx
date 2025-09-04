import { ReactNode } from 'react'
import { TrendingUp, TrendingDown } from 'lucide-react'

interface StatsCardProps {
  title: string
  value: string
  icon: ReactNode
  trend?: string
  trendUp?: boolean
}

export default function StatsCard({ title, value, icon, trend, trendUp }: StatsCardProps) {
  return (
    <div className="card p-6">
      <div className="flex items-center justify-between">
        <div>
          <p className="text-sm font-medium text-gray-600">{title}</p>
          <p className="text-3xl font-bold text-gray-900 mt-2">{value}</p>
          {trend && (
            <div className={`flex items-center mt-2 text-sm ${
              trendUp ? 'text-green-600' : 'text-red-600'
            }`}>
              {trendUp ? (
                <TrendingUp className="w-4 h-4 mr-1" />
              ) : (
                <TrendingDown className="w-4 h-4 mr-1" />
              )}
              {trend}
            </div>
          )}
        </div>
        <div className={`p-3 rounded-full ${
          trendUp === false ? 'bg-red-100 text-red-600' : 'bg-primary-100 text-primary-600'
        }`}>
          {icon}
        </div>
      </div>
    </div>
  )
}