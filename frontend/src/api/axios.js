import axios from 'axios'

const baseURL = import.meta.env.VITE_API_BASE_URL

const api = axios.create({
  baseURL,
  headers: {
    'Content-Type': 'application/json',
  },
})

api.interceptors.request.use((config) => {
  if (!baseURL) {
    throw new Error('Missing VITE_API_BASE_URL. Create frontend/.env from .env.example.')
  }

  const token = localStorage.getItem('auth_token')
  if (token) {
    config.headers.Authorization = `Bearer ${token}`
  }

  return config
})

api.interceptors.response.use(
  (response) => response,
  (error) => {
    const status = error?.response?.status
    const message = error?.response?.data?.message || ''

    if (
      status === 401 ||
      (status === 403 && message.toLowerCase().includes('dinonaktifkan'))
    ) {
      window.dispatchEvent(
        new CustomEvent('auth:unauthorized', {
          detail: { message },
        }),
      )
    }

    return Promise.reject(error)
  },
)

export default api
