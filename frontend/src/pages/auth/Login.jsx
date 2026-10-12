import { useState } from 'react'
import { Link, Navigate, useLocation, useNavigate } from 'react-router-dom'
import ErrorMessage from '../../components/ErrorMessage'
import { useAuth } from '../../context/AuthContext'

function Login() {
  const { authMessage, isAdmin, isAuthenticated, login } = useAuth()
  const navigate = useNavigate()
  const location = useLocation()

  const [form, setForm] = useState({ email: '', password: '' })
  const [showPassword, setShowPassword] = useState(false)
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState(authMessage)

  if (isAuthenticated) {
    return <Navigate to={isAdmin ? '/admin' : '/products'} replace />
  }

  function handleChange(event) {
    const { name, value } = event.target
    setForm((currentForm) => ({ ...currentForm, [name]: value }))
  }

  async function handleSubmit(event) {
    event.preventDefault()

    if (!form.email.trim() || !form.password) {
      setError('Email dan password wajib diisi.')
      return
    }

    try {
      setLoading(true)
      setError('')
      const user = await login({
        email: form.email.trim(),
        password: form.password,
      })
      const fallbackPath = user.role === 'admin' ? '/admin' : '/products'
      navigate(location.state?.from?.pathname || fallbackPath, { replace: true })
    } catch (err) {
      setError(err?.response?.data?.message || 'Login gagal.')
    } finally {
      setLoading(false)
    }
  }

  return (
    <main className="auth-page">
      <section className="auth-panel">
        <div>
          <p className="eyebrow">Admin Management</p>
          <h1>Login</h1>
          <p className="auth-copy">Masuk untuk mengelola dashboard dan data user.</p>
        </div>

        <ErrorMessage message={error} />

        <form className="form auth-form" onSubmit={handleSubmit}>
          <div className="form-group">
            <label htmlFor="email">Email</label>
            <input
              id="email"
              name="email"
              type="email"
              value={form.email}
              onChange={handleChange}
              placeholder="admin@example.com"
            />
          </div>

          <div className="form-group">
            <label htmlFor="password">Password</label>
            <div className="password-field">
              <input
                id="password"
                name="password"
                type={showPassword ? 'text' : 'password'}
                value={form.password}
                onChange={handleChange}
                placeholder="Password"
              />
              <button
                type="button"
                className="button button--small"
                onClick={() => setShowPassword((current) => !current)}
              >
                {showPassword ? 'Hide' : 'Show'}
              </button>
            </div>
          </div>

          <button type="submit" className="button button--primary" disabled={loading}>
            {loading ? 'Logging in...' : 'Login'}
          </button>
        </form>

        <Link to="/products" className="auth-link">
          Back to products
        </Link>
      </section>
    </main>
  )
}

export default Login
