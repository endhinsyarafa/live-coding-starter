import api from '../api/axios'

export async function login(credentials) {
  const response = await api.post('/auth/login', credentials)
  return response.data.data
}

export async function logout() {
  const response = await api.post('/auth/logout')
  return response.data
}

export async function getMe() {
  const response = await api.get('/auth/me')
  return response.data.data
}
