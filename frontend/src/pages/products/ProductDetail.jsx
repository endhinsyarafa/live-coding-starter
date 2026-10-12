import { useEffect, useState } from 'react'
import { Link, useParams } from 'react-router-dom'
import ErrorMessage from '../../components/ErrorMessage'
import Loading from '../../components/Loading'
import { getProduct } from '../../services/productService'

function ProductDetail() {
  const { id } = useParams()
  const [product, setProduct] = useState(null)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState('')

  useEffect(() => {
    async function loadProduct() {
      try {
        setLoading(true)
        setError('')
        const data = await getProduct(id)
        setProduct(data)
      } catch (err) {
        setError(getErrorMessage(err, 'Failed to load product.'))
      } finally {
        setLoading(false)
      }
    }

    loadProduct()
  }, [id])

  if (loading) {
    return <Loading message="Loading product..." />
  }

  if (error) {
    return (
      <section className="page-section page-section--narrow">
        <ErrorMessage message={error} />
        <Link to="/products" className="button">
          Back
        </Link>
      </section>
    )
  }

  return (
    <section className="page-section page-section--narrow">
      <div className="page-header">
        <div>
          <p className="eyebrow">Product Detail</p>
          <h1>{product.name}</h1>
        </div>
        <Link to={`/products/${product.id}/edit`} className="button button--primary">
          Edit
        </Link>
      </div>

      <dl className="detail-list">
        <div>
          <dt>ID</dt>
          <dd>{product.id}</dd>
        </div>
        <div>
          <dt>Name</dt>
          <dd>{product.name}</dd>
        </div>
        <div>
          <dt>Description</dt>
          <dd>{product.description || '-'}</dd>
        </div>
        <div>
          <dt>Price</dt>
          <dd>{formatCurrency(product.price)}</dd>
        </div>
        <div>
          <dt>Stock</dt>
          <dd>{product.stock}</dd>
        </div>
        <div>
          <dt>Created At</dt>
          <dd>{formatDate(product.created_at)}</dd>
        </div>
        <div>
          <dt>Updated At</dt>
          <dd>{formatDate(product.updated_at)}</dd>
        </div>
      </dl>

      <div className="form-actions">
        <Link to="/products" className="button">
          Back
        </Link>
      </div>
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

function formatDate(value) {
  return new Intl.DateTimeFormat('id-ID', {
    dateStyle: 'medium',
    timeStyle: 'short',
  }).format(new Date(value))
}

function getErrorMessage(err, fallback) {
  return err?.response?.data?.message || fallback
}

export default ProductDetail
