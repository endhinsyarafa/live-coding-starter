import { useEffect, useState } from 'react'
import { Link } from 'react-router-dom'
import ErrorMessage from '../../components/ErrorMessage'
import Loading from '../../components/Loading'
import { deleteProduct, getProducts } from '../../services/productService'

function ProductList() {
  const [products, setProducts] = useState([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState('')
  const [deletingId, setDeletingId] = useState(null)

  async function loadProducts() {
    try {
      setLoading(true)
      setError('')
      const data = await getProducts()
      setProducts(Array.isArray(data) ? data : [])
    } catch (err) {
      setProducts([])
      setError(getErrorMessage(err, 'Failed to load products.'))
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    loadProducts()
  }, [])

  async function handleDelete(product) {
    const confirmed = window.confirm(`Delete product "${product.name}"?`)
    if (!confirmed) {
      return
    }

    try {
      setDeletingId(product.id)
      setError('')
      await deleteProduct(product.id)
      setProducts((currentProducts) =>
        Array.isArray(currentProducts)
          ? currentProducts.filter((item) => item.id !== product.id)
          : [],
      )
    } catch (err) {
      setError(getErrorMessage(err, 'Failed to delete product.'))
    } finally {
      setDeletingId(null)
    }
  }

  if (loading) {
    return <Loading message="Loading products..." />
  }

  const safeProducts = Array.isArray(products) ? products : []

  return (
    <section className="page-section">
      <div className="page-header">
        <div>
          <p className="eyebrow">Products</p>
          <h1>Product List</h1>
        </div>
        <Link to="/products/new" className="button button--primary">
          Add Product
        </Link>
      </div>

      <ErrorMessage message={error} />

      {safeProducts.length === 0 ? (
        <div className="empty-state">
          <h2>No products yet</h2>
          <p>Create your first product to test the CRUD flow.</p>
          <Link to="/products/new" className="button button--primary">
            Add Product
          </Link>
        </div>
      ) : (
        <div className="table-wrapper">
          <table>
            <thead>
              <tr>
                <th>Name</th>
                <th>Description</th>
                <th>Price</th>
                <th>Stock</th>
                <th>Actions</th>
              </tr>
            </thead>
            <tbody>
              {safeProducts.map((product) => (
                <tr key={product.id}>
                  <td data-label="Name">{product.name}</td>
                  <td data-label="Description">{product.description || '-'}</td>
                  <td data-label="Price">{formatCurrency(product.price)}</td>
                  <td data-label="Stock">{product.stock}</td>
                  <td data-label="Actions">
                    <div className="table-actions">
                      <Link to={`/products/${product.id}`} className="button button--small">
                        Detail
                      </Link>
                      <Link
                        to={`/products/${product.id}/edit`}
                        className="button button--small"
                      >
                        Edit
                      </Link>
                      <button
                        type="button"
                        className="button button--small button--danger"
                        onClick={() => handleDelete(product)}
                        disabled={deletingId === product.id}
                      >
                        {deletingId === product.id ? 'Deleting...' : 'Delete'}
                      </button>
                    </div>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      )}
    </section>
  )
}

function formatCurrency(value) {
  return new Intl.NumberFormat('id-ID', {
    style: 'currency',
    currency: 'IDR',
    maximumFractionDigits: 0,
  }).format(value)
}

function getErrorMessage(err, fallback) {
  return err?.response?.data?.message || err?.message || fallback
}

export default ProductList
