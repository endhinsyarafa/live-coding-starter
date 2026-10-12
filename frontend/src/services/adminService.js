import api from '../api/axios'

export async function getDashboardStats() {
  const response = await api.get('/admin/dashboard')
  return response.data.data
}

export async function getUsers(params = {}) {
  const response = await api.get('/admin/users', { params })
  return {
    data: response.data.data,
    meta: response.data.meta,
  }
}

export async function getUser(id) {
  const response = await api.get(`/admin/users/${id}`)
  return response.data.data
}

export async function createUser(data) {
  const response = await api.post('/admin/users', data)
  return response.data.data
}

export async function updateUser(id, data) {
  const response = await api.patch(`/admin/users/${id}`, data)
  return response.data.data
}

export async function updateUserStatus(id, isActive) {
  const response = await api.patch(`/admin/users/${id}/status`, {
    is_active: isActive,
  })
  return response.data.data
}
