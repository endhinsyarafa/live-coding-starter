import { useEffect, useState } from 'react'
import ErrorMessage from '../../components/ErrorMessage'
import Loading from '../../components/Loading'
import StatCard from '../../components/admin/StatCard'
import { getDashboardStats } from '../../services/adminService'

function AdminDashboard() {
  const [stats, setStats] = useState(null)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState('')

  useEffect(() => {
    loadStats()
  }, [])

  async function loadStats() {
    try {
      setLoading(true)
      setError('')
      const data = await getDashboardStats()
      setStats(data)
    } catch (err) {
      setError(err?.response?.data?.message || 'Gagal memuat dashboard.')
    } finally {
      setLoading(false)
    }
  }

  if (loading) {
    return <Loading message="Loading dashboard..." />
  }

  return (
    <section className="admin-page">
      <div className="page-header">
        <div>
          <p className="eyebrow">Admin Dashboard</p>
          <h1>Dashboard</h1>
        </div>
      </div>

      <ErrorMessage message={error} />

      <div className="stats-grid">
        <StatCard label="Total Users" value={stats?.total_users ?? 0} />
        <StatCard label="Active Users" value={stats?.active_users ?? 0} />
        <StatCard label="Inactive Users" value={stats?.inactive_users ?? 0} />
        <StatCard label="New Users" value={stats?.new_users ?? 0} />
      </div>
    </section>
  )
}

export default AdminDashboard
