import { useState, useEffect } from 'react'
import Layout from '@/components/Layout'
import StatsCard from '@/components/StatsCard'
import { Users, FileText, DollarSign, TrendingUp } from 'lucide-react'

interface DashboardStats {
  totalUsers: number
  totalApplications: number
  totalAmountRequested: number
  approvalRate: number
}

export default function Dashboard() {
  const [stats, setStats] = useState<DashboardStats>({
    totalUsers: 0,
    totalApplications: 0,
    totalAmountRequested: 0,
    approvalRate: 0
  })
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    // Fetch dashboard statistics from API
    const fetchStats = async () => {
      try {
        // This would call the actual API endpoint
        // const response = await fetch('/api/v1/admin/loan-applications/stats')
        // const data = await response.json()
        
        // Mock data for now
        setTimeout(() => {
          setStats({
            totalUsers: 1250,
            totalApplications: 856,
            totalAmountRequested: 12450000,
            approvalRate: 75.5
          })
          setLoading(false)
        }, 1000)
      } catch (error) {
        console.error('Failed to fetch stats:', error)
        setLoading(false)
      }
    }

    fetchStats()
  }, [])

  return (
    <Layout>
      <div className="p-6">
        <h1 className="text-3xl font-bold text-gray-900 mb-6">Dashboard Overview</h1>
        
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6 mb-8">
          <StatsCard
            title="Total Users"
            value={loading ? '...' : stats.totalUsers.toLocaleString()}
            icon={<Users className="w-6 h-6" />}
            trend="+12% from last month"
            trendUp={true}
          />
          
          <StatsCard
            title="Loan Applications"
            value={loading ? '...' : stats.totalApplications.toLocaleString()}
            icon={<FileText className="w-6 h-6" />}
            trend="+8% from last month"
            trendUp={true}
          />
          
          <StatsCard
            title="Total Requested"
            value={loading ? '...' : `$${(stats.totalAmountRequested / 1000000).toFixed(1)}M`}
            icon={<DollarSign className="w-6 h-6" />}
            trend="+15% from last month"
            trendUp={true}
          />
          
          <StatsCard
            title="Approval Rate"
            value={loading ? '...' : `${stats.approvalRate}%`}
            icon={<TrendingUp className="w-6 h-6" />}
            trend="-2% from last month"
            trendUp={false}
          />
        </div>

        <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
          <div className="card p-6">
            <h2 className="text-xl font-semibold mb-4">Recent Activity</h2>
            <div className="space-y-3">
              {loading ? (
                <div className="animate-pulse space-y-3">
                  <div className="h-4 bg-gray-200 rounded w-3/4"></div>
                  <div className="h-4 bg-gray-200 rounded w-1/2"></div>
                  <div className="h-4 bg-gray-200 rounded w-2/3"></div>
                </div>
              ) : (
                <>
                  <div className="flex justify-between items-center py-2">
                    <span className="text-sm text-gray-600">New loan application submitted</span>
                    <span className="text-xs text-gray-400">2 hours ago</span>
                  </div>
                  <div className="flex justify-between items-center py-2">
                    <span className="text-sm text-gray-600">User KYC approved</span>
                    <span className="text-xs text-gray-400">4 hours ago</span>
                  </div>
                  <div className="flex justify-between items-center py-2">
                    <span className="text-sm text-gray-600">Invoice processed successfully</span>
                    <span className="text-xs text-gray-400">6 hours ago</span>
                  </div>
                </>
              )}
            </div>
          </div>

          <div className="card p-6">
            <h2 className="text-xl font-semibold mb-4">Quick Actions</h2>
            <div className="grid grid-cols-2 gap-3">
              <button className="btn btn-primary">
                Review Applications
              </button>
              <button className="btn btn-secondary">
                Manage Users
              </button>
              <button className="btn btn-secondary">
                Generate Reports
              </button>
              <button className="btn btn-secondary">
                View Analytics
              </button>
            </div>
          </div>
        </div>
      </div>
    </Layout>
  )
}