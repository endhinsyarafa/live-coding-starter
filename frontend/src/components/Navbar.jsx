import { NavLink } from 'react-router-dom'
import { useAuth } from '../context/AuthContext'

function Navbar() {
  const { isAdmin, isAuthenticated, logout, user } = useAuth()

  return (
    <header className="navbar">
      <div className="navbar__inner">
        <NavLink to="/products" className="navbar__brand">
          Live Coding Starter
        </NavLink>
        <nav className="navbar__links">
          <NavLink to="/products">Products</NavLink>
          {isAdmin && <NavLink to="/admin">Admin</NavLink>}
          {isAuthenticated ? (
            <button type="button" className="nav-button" onClick={logout}>
              Logout {user?.name ? `(${user.name})` : ''}
            </button>
          ) : (
            <NavLink to="/login">Login</NavLink>
          )}
        </nav>
      </div>
    </header>
  )
}

export default Navbar
