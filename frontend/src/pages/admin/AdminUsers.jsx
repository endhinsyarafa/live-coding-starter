import { useEffect, useMemo, useState } from 'react'
import ConfirmDialog from '../../components/admin/ConfirmDialog'
import Toast from '../../components/admin/Toast'
import UserForm from '../../components/admin/UserForm'
import UserStatusBadge from '../../components/admin/UserStatusBadge'
import UserStatusToggle from '../../components/admin/UserStatusToggle'
import ErrorMessage from '../../components/ErrorMessage'
import Loading from '../../components/Loading'
import { useAuth } from '../../context/AuthContext'
import {
  createUser,
  getUsers,
  updateUser,
  updateUserStatus,
} from '../../services/adminService'

const initialMeta = {
  page: 1,
  limit: 10,
  total: 0,
  total_pages: 0,
}

function AdminUsers() {
  const { user: currentUser } = useAuth()

  const [users, setUsers] = useState([])
  const [meta, setMeta] = useState(initialMeta)
  const [search, setSearch] = useState('')
  const [debouncedSearch, setDebouncedSearch] = useState('')
  const [status, setStatus] = useState('')
  const [role, setRole] = useState('')
  const [page, setPage] = useState(1)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState('')
  const [toast, setToast] = useState(null)
  const [formMode, setFormMode] = useState(null)
  const [selectedUser, setSelectedUser] = useState(null)
  const [submitting, setSubmitting] = useState(false)
  const [confirmAction, setConfirmAction] = useState(null)

  const safeUsers = Array.isArray(users) ? users : []

  useEffect(() => {
    const timer = window.setTimeout(() => {
      setDebouncedSearch(search)
      setPage(1)
    }, 350)

    return () => window.clearTimeout(timer)
  }, [search])

  useEffect(() => {
    loadUsers()
  }, [debouncedSearch, status, role, page])

  async function loadUsers() {
    try {
      setLoading(true)
      setError('')
      const result = await getUsers({
        search: debouncedSearch || undefined,
        status: status || undefined,
        role: role || undefined,
        page,
        limit: 10,
      })
      setUsers(Array.isArray(result.data) ? result.data : [])
      setMeta(result.meta || initialMeta)
    } catch (err) {
      setUsers([])
      setError(err?.response?.data?.message || 'Gagal memuat data user.')
    } finally {
      setLoading(false)
    }
  }

  function openCreateForm() {
    setSelectedUser(null)
    setFormMode('create')
  }

  function openEditForm(user) {
    setSelectedUser(user)
    setFormMode('edit')
  }

  async function handleSubmit(payload) {
    try {
      setSubmitting(true)
      setError('')
      if (formMode === 'create') {
        const createdUser = await createUser(payload)
        setUsers((currentUsers) => [createdUser, ...currentUsers])
        showToast('User berhasil ditambahkan.')
      } else {
        const updatedUser = await updateUser(selectedUser.id, payload)
        replaceUser(updatedUser)
        showToast('Data user berhasil diperbarui.')
      }
      closeForm()
      loadUsers()
    } catch (err) {
      setError(err?.response?.data?.message || 'Gagal menyimpan data user.')
    } finally {
      setSubmitting(false)
    }
  }

  function closeForm() {
    setFormMode(null)
    setSelectedUser(null)
  }

  function requestStatusChange(targetUser, nextStatus) {
    if (targetUser.id === currentUser.id && !nextStatus) {
      showToast('Anda tidak dapat menonaktifkan akun sendiri.', 'error')
      return
    }

    setConfirmAction({
      user: targetUser,
      nextStatus,
    })
  }

  async function confirmStatusChange() {
    if (!confirmAction) {
      return
    }

    try {
      const updatedUser = await updateUserStatus(confirmAction.user.id, confirmAction.nextStatus)
      replaceUser(updatedUser)
      showToast(
        updatedUser.is_active ? 'User berhasil diaktifkan.' : 'User berhasil dinonaktifkan.',
      )
    } catch (err) {
      showToast(err?.response?.data?.message || 'Gagal memperbarui status user.', 'error')
    } finally {
      setConfirmAction(null)
    }
  }

  function replaceUser(updatedUser) {
    setUsers((currentUsers) =>
      Array.isArray(currentUsers)
        ? currentUsers.map((item) => (item.id === updatedUser.id ? updatedUser : item))
        : [],
    )
  }

  function showToast(message, type = 'success') {
    setToast({ message, type })
  }

  const emptyText = useMemo(() => {
    if (debouncedSearch || status || role) {
      return 'Tidak ada user yang sesuai dengan pencarian.'
    }

    return 'Belum ada user.'
  }, [debouncedSearch, role, status])

  return (
    <section className="admin-page">
      <div className="page-header">
        <div>
          <p className="eyebrow">Admin</p>
          <h1>User Management</h1>
        </div>
        <button type="button" className="button button--primary" onClick={openCreateForm}>
          + Add User
        </button>
      </div>

      <Toast
        message={toast?.message}
        type={toast?.type}
        onClose={() => setToast(null)}
      />

      <ErrorMessage message={error} />

      {formMode && (
        <UserForm
          mode={formMode}
          onCancel={closeForm}
          onSubmit={handleSubmit}
          submitting={submitting}
          user={selectedUser}
        />
      )}

      <div className="admin-filters">
        <input
          value={search}
          onChange={(event) => setSearch(event.target.value)}
          placeholder="Search name or email"
        />
        <select
          value={status}
          onChange={(event) => {
            setStatus(event.target.value)
            setPage(1)
          }}
        >
          <option value="">All Status</option>
          <option value="active">Active</option>
          <option value="inactive">Inactive</option>
        </select>
        <select
          value={role}
          onChange={(event) => {
            setRole(event.target.value)
            setPage(1)
          }}
        >
          <option value="">All Roles</option>
          <option value="admin">Admin</option>
          <option value="user">User</option>
        </select>
        <button type="button" className="button" onClick={loadUsers}>
          Retry
        </button>
      </div>

      {loading ? (
        <Loading message="Loading users..." />
      ) : safeUsers.length === 0 ? (
        <div className="empty-state">
          <h2>{emptyText}</h2>
          <p>User yang dibuat akan tampil di tabel ini.</p>
        </div>
      ) : (
        <div className="table-wrapper">
          <table>
            <thead>
              <tr>
                <th>#</th>
                <th>Name</th>
                <th>Email</th>
                <th>Role</th>
                <th>Status</th>
                <th>Last Login</th>
                <th>Created At</th>
                <th>Actions</th>
              </tr>
            </thead>
            <tbody>
              {safeUsers.map((user, index) => {
                const isCurrentUser = user.id === currentUser?.id
                return (
                  <tr key={user.id}>
                    <td data-label="#">{(meta.page - 1) * meta.limit + index + 1}</td>
                    <td data-label="Name">{user.name}</td>
                    <td data-label="Email">{user.email}</td>
                    <td data-label="Role">{formatRole(user.role)}</td>
                    <td data-label="Status">
                      <UserStatusBadge isActive={user.is_active} />
                    </td>
                    <td data-label="Last Login">{formatDate(user.last_login_at)}</td>
                    <td data-label="Created At">{formatDate(user.created_at)}</td>
                    <td data-label="Actions">
                      <div className="table-actions">
                        <UserStatusToggle
                          disabled={isCurrentUser}
                          isActive={user.is_active}
                          onChange={(nextStatus) => requestStatusChange(user, nextStatus)}
                          title={
                            isCurrentUser
                              ? 'Anda tidak dapat menonaktifkan akun sendiri.'
                              : undefined
                          }
                        />
                        <button
                          type="button"
                          className="button button--small"
                          onClick={() => openEditForm(user)}
                        >
                          Edit
                        </button>
                      </div>
                    </td>
                  </tr>
                )
              })}
            </tbody>
          </table>
        </div>
      )}

      <div className="pagination">
        <button
          type="button"
          className="button button--small"
          disabled={meta.page <= 1}
          onClick={() => setPage((currentPage) => Math.max(currentPage - 1, 1))}
        >
          Previous
        </button>
        <span>
          Page {meta.page} of {meta.total_pages || 1}
        </span>
        <button
          type="button"
          className="button button--small"
          disabled={meta.page >= meta.total_pages}
          onClick={() => setPage((currentPage) => currentPage + 1)}
        >
          Next
        </button>
      </div>

      <ConfirmDialog
        title={
          confirmAction
            ? confirmAction.nextStatus
              ? 'Aktifkan user ini?'
              : 'Nonaktifkan user ini?'
            : ''
        }
        description={
          confirmAction?.nextStatus
            ? 'User akan dapat login dan mengakses aplikasi kembali.'
            : 'User tidak akan dapat login atau mengakses aplikasi sampai diaktifkan kembali.'
        }
        confirmLabel={confirmAction?.nextStatus ? 'Activate User' : 'Deactivate User'}
        onCancel={() => setConfirmAction(null)}
        onConfirm={confirmStatusChange}
      />
    </section>
  )
}

function formatDate(value) {
  if (!value) {
    return '-'
  }

  return new Intl.DateTimeFormat('id-ID', {
    dateStyle: 'medium',
    timeStyle: 'short',
  }).format(new Date(value))
}

function formatRole(role) {
  return role === 'admin' ? 'Admin' : 'User'
}

export default AdminUsers
