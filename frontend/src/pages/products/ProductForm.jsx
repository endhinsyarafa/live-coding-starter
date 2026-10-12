import { useEffect, useState } from 'react'
import { Link, useNavigate, useParams } from 'react-router-dom'
import ErrorMessage from '../../components/ErrorMessage'
import Loading from '../../components/Loading'
import { createProduct, getProduct, updateProduct } from '../../services/productService'

const initialForm = {
  name: '',
  description: '',
  price: '',
  stock: '',
}

function ProductForm() {
  const { id } = useParams()
  const navigate = useNavigate()
  const isEdit = Boolean(id)

  const [form, setForm] = useState(initialForm)
  const [loading, setLoading] = useState(isEdit)
  const [submitting, setSubmitting] = useState(false)
  const [error, setError] = useState('')
  const [validationErrors, setValidationErrors] = useState({})

  useEffect(() => {
    if (!isEdit) {
      return
    }

    async function loadProduct() {
      try {
        setLoading(true)
        setError('')
        const product = await getProduct(id)
        setForm({
          name: product.name,
          description: product.description || '',
          price: String(product.price),
          stock: String(product.stock),
        })
      } catch (err) {
        setError(getErrorMessage(err, 'Failed to load product.'))
      } finally {
        setLoading(false)
      }
    }

    loadProduct()
  }, [id, isEdit])

  function handleChange(event) {
    const { name, value } = event.target
    setForm((currentForm) => ({
      ...currentForm,
      [name]: value,
    }))
  }

  async function handleSubmit(event) {
    event.preventDefault()

    const errors = validateForm(form)
    setValidationErrors(errors)
    if (Object.keys(errors).length > 0) {
      return
    }

    const payload = {
      name: form.name.trim(),
      description: form.description.trim(),
      price: Number(form.price),
      stock: Number(form.stock),
    }

    try {
      setSubmitting(true)
      setError('')

      if (isEdit) {
        await updateProduct(id, payload)
      } else {
        await createProduct(payload)
      }

      navigate('/products')
    } catch (err) {
      setError(getErrorMessage(err, 'Failed to save product.'))
    } finally {
      setSubmitting(false)
    }
  }

  if (loading) {
    return <Loading message="Loading product..." />
  }

  return (
    <section className="page-section page-section--narrow">
      <div className="page-header">
        <div>
          <p className="eyebrow">Products</p>
          <h1>{isEdit ? 'Edit Product' : 'Add Product'}</h1>
        </div>
      </div>

      <ErrorMessage message={error} />

      <form className="form" onSubmit={handleSubmit}>
        <div className="form-group">
          <label htmlFor="name">Name</label>
          <input
            id="name"
            name="name"
            value={form.name}
            onChange={handleChange}
            placeholder="Keyboard Mechanical"
          />
          {validationErrors.name && (
            <span className="field-error">{validationErrors.name}</span>
          )}
        </div>

        <div className="form-group">
          <label htmlFor="description">Description</label>
          <textarea
            id="description"
            name="description"
            value={form.description}
            onChange={handleChange}
            placeholder="Mechanical keyboard untuk programmer"
            rows="4"
          />
        </div>

        <div className="form-grid">
          <div className="form-group">
            <label htmlFor="price">Price</label>
            <input
              id="price"
              name="price"
              type="number"
              min="0"
              value={form.price}
              onChange={handleChange}
              placeholder="750000"
            />
            {validationErrors.price && (
              <span className="field-error">{validationErrors.price}</span>
            )}
          </div>

          <div className="form-group">
            <label htmlFor="stock">Stock</label>
            <input
              id="stock"
              name="stock"
              type="number"
              min="0"
              value={form.stock}
              onChange={handleChange}
              placeholder="10"
            />
            {validationErrors.stock && (
              <span className="field-error">{validationErrors.stock}</span>
            )}
          </div>
        </div>

        <div className="form-actions">
          <Link to="/products" className="button">
            Cancel
          </Link>
          <button type="submit" className="button button--primary" disabled={submitting}>
            {submitting ? 'Saving...' : 'Save Product'}
          </button>
        </div>
      </form>
    </section>
  )
}

function validateForm(form) {
  const errors = {}
  const price = Number(form.price)
  const stock = Number(form.stock)

  if (!form.name.trim()) {
    errors.name = 'Name is required.'
  }

  if (form.price === '' || Number.isNaN(price) || price < 0) {
    errors.price = 'Price must be zero or greater.'
  }

  if (form.stock === '' || Number.isNaN(stock) || stock < 0) {
    errors.stock = 'Stock must be zero or greater.'
  }

  return errors
}

function getErrorMessage(err, fallback) {
  return err?.response?.data?.message || fallback
}

export default ProductForm
