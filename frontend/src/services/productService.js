import api from '../api/axios'

export async function getProducts() {
  const response = await api.get('/products')
  const data = getResponseData(response)

  return Array.isArray(data) ? data : []
}

export async function getProduct(id) {
  const response = await api.get(`/products/${id}`)
  return getResponseData(response)
}

export async function createProduct(data) {
  const response = await api.post('/products', data)
  return getResponseData(response)
}

export async function updateProduct(id, data) {
  const response = await api.put(`/products/${id}`, data)
  return getResponseData(response)
}

export async function deleteProduct(id) {
  const response = await api.delete(`/products/${id}`)
  return getResponseData(response)
}

function getResponseData(response) {
  if (!response?.data || typeof response.data !== 'object') {
    throw new Error('Invalid API response.')
  }

  return response.data.data
}
