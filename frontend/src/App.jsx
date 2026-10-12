import { Navigate, Route, Routes } from 'react-router-dom'
import AdminLayout from './components/admin/AdminLayout'
import Layout from './components/Layout'
import ProtectedRoute from './components/ProtectedRoute'
import AdminDashboard from './pages/admin/AdminDashboard'
import AdminUsers from './pages/admin/AdminUsers'
import Login from './pages/auth/Login'
import Unauthorized from './pages/auth/Unauthorized'
import ProductDetail from './pages/products/ProductDetail'
import ProductForm from './pages/products/ProductForm'
import ProductList from './pages/products/ProductList'

function App() {
  return (
    <Routes>
      <Route path="/login" element={<Login />} />
      <Route path="/unauthorized" element={<Unauthorized />} />

      <Route element={<Layout />}>
        <Route path="/" element={<Navigate to="/products" replace />} />
        <Route path="/products" element={<ProductList />} />
        <Route path="/products/new" element={<ProductForm />} />
        <Route path="/products/:id" element={<ProductDetail />} />
        <Route path="/products/:id/edit" element={<ProductForm />} />
      </Route>

      <Route element={<ProtectedRoute adminOnly />}>
        <Route element={<AdminLayout />}>
          <Route path="/admin" element={<AdminDashboard />} />
          <Route path="/admin/users" element={<AdminUsers />} />
        </Route>
      </Route>
    </Routes>
  )
}

export default App
