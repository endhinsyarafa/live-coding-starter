import { Navigate, Outlet, useLocation } from 'react-router-dom'
import Loading from './Loading'
import { useAuth } from '../context/AuthContext'

function ProtectedRoute({ adminOnly = false }) {
  const { isAdmin, isAuthenticated, loading } = useAuth()
  const location = useLocation()

  if (loading) {
    return <Loading message="Checking session..." />
  }

  if (!isAuthenticated) {
    return <Navigate to="/login" replace state={{ from: location }} />
  }

  if (adminOnly && !isAdmin) {
    return <Navigate to="/unauthorized" replace />
  }

  return <Outlet />
}

export default ProtectedRoute
