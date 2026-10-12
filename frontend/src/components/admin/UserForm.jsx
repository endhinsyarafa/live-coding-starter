import { useEffect, useState } from 'react'

const defaultForm = {
  name: '',
  email: '',
  password: '',
  confirm_password: '',
  role: 'user',
  is_active: true,
}

function UserForm({ mode = 'create', onCancel, onSubmit, submitting, user }) {
  const [form, setForm] = useState(defaultForm)
  const [errors, setErrors] = useState({})

  useEffect(() => {
    if (user) {
      setForm({
        name: user.name || '',
        email: user.email || '',
        password: '',
        confirm_password: '',
        role: user.role || 'user',
        is_active: Boolean(user.is_active),
      })
    } else {
      setForm(defaultForm)
    }
  }, [user])

  function handleChange(event) {
    const { checked, name, type, value } = event.target
    setForm((currentForm) => ({
      ...currentForm,
      [name]: type === 'checkbox' ? checked : value,
    }))
  }

  function handleSubmit(event) {
    event.preventDefault()
    const validationErrors = validate(form, mode)
    setErrors(validationErrors)
    if (Object.keys(validationErrors).length > 0) {
      return
    }

    const payload = {
      name: form.name.trim(),
      email: form.email.trim(),
      role: form.role,
      is_active: form.is_active,
    }

    if (mode === 'create') {
      payload.password = form.password
      payload.confirm_password = form.confirm_password
    }

    onSubmit(payload)
  }

  return (
    <form className="form admin-user-form" onSubmit={handleSubmit}>
      <div className="form-grid">
        <div className="form-group">
          <label htmlFor="user-name">Name</label>
          <input id="user-name" name="name" value={form.name} onChange={handleChange} />
          {errors.name && <span className="field-error">{errors.name}</span>}
        </div>
        <div className="form-group">
          <label htmlFor="user-email">Email</label>
          <input id="user-email" name="email" type="email" value={form.email} onChange={handleChange} />
          {errors.email && <span className="field-error">{errors.email}</span>}
        </div>
      </div>

      {mode === 'create' && (
        <div className="form-grid">
          <div className="form-group">
            <label htmlFor="user-password">Password</label>
            <input
              id="user-password"
              name="password"
              type="password"
              value={form.password}
              onChange={handleChange}
            />
            {errors.password && <span className="field-error">{errors.password}</span>}
          </div>
          <div className="form-group">
            <label htmlFor="user-confirm-password">Confirm Password</label>
            <input
              id="user-confirm-password"
              name="confirm_password"
              type="password"
              value={form.confirm_password}
              onChange={handleChange}
            />
            {errors.confirm_password && (
              <span className="field-error">{errors.confirm_password}</span>
            )}
          </div>
        </div>
      )}

      <div className="form-grid">
        <div className="form-group">
          <label htmlFor="user-role">Role</label>
          <select id="user-role" name="role" value={form.role} onChange={handleChange}>
            <option value="admin">Admin</option>
            <option value="user">User</option>
          </select>
          {errors.role && <span className="field-error">{errors.role}</span>}
        </div>
        <label className="checkbox-row">
          <input
            name="is_active"
            type="checkbox"
            checked={form.is_active}
            onChange={handleChange}
          />
          Active
        </label>
      </div>

      <div className="form-actions">
        <button type="button" className="button" onClick={onCancel}>
          Cancel
        </button>
        <button type="submit" className="button button--primary" disabled={submitting}>
          {submitting ? 'Saving...' : mode === 'create' ? 'Add User' : 'Save Changes'}
        </button>
      </div>
    </form>
  )
}

function validate(form, mode) {
  const errors = {}
  if (!form.name.trim()) {
    errors.name = 'Name is required.'
  }
  if (!/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(form.email)) {
    errors.email = 'Valid email is required.'
  }
  if (!['admin', 'user'].includes(form.role)) {
    errors.role = 'Role is required.'
  }
  if (mode === 'create') {
    if (form.password.length < 8) {
      errors.password = 'Password must be at least 8 characters.'
    }
    if (form.password !== form.confirm_password) {
      errors.confirm_password = 'Password confirmation must match.'
    }
  }

  return errors
}

export default UserForm
