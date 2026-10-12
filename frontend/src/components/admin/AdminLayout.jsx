import { useState } from 'react'
import { NavLink, Outlet, useNavigate } from 'react-router-dom'
import { useAuth } from '../../context/AuthContext'

function AdminLayout() {
  const { logout, user } = useAuth()
  const navigate = useNavigate()
  const [sidebarOpen, setSidebarOpen] = useState(false)

  async function handleLogout() {
    await logout()
    navigate('/login', { replace: true })
  }

  return (
    <div className="admin-shell">
      <aside className={`admin-sidebar ${sidebarOpen ? 'is-open' : ''}`}>
        <div className="admin-sidebar__brand">Admin</div>
        <nav className="admin-sidebar__nav">
          <NavLink to="/admin" end onClick={() => setSidebarOpen(false)}>
            Dashboard
          </NavLink>
          <NavLink to="/admin/users" onClick={() => setSidebarOpen(false)}>
            User Management
          </NavLink>
        </nav>
        <button type="button" className="admin-sidebar__logout" onClick={handleLogout}>
          Logout
        </button>
      </aside>

      <div className="admin-main">
        <header className="admin-topbar">
          <button
            type="button"
            className="admin-menu-button"
            onClick={() => setSidebarOpen((current) => !current)}
          >
            Menu
          </button>
          <div className="admin-profile">
            <div className="admin-avatar">{getInitial(user?.name)}</div>
            <div>
              <strong>{user?.name}</strong>
              <span>{user?.role}</span>
            </div>
            <button type="button" className="button button--small" onClick={handleLogout}>
              Logout
            </button>
          </div>
        </header>

        <main className="admin-content">
          <Outlet />
        </main>
      </div>
    </div>
  )
}

function getInitial(name = '') {
  return name.trim().charAt(0).toUpperCase() || 'A'
}

export default AdminLayout
