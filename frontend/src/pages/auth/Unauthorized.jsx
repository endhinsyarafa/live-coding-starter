import { Link } from 'react-router-dom'

function Unauthorized() {
  return (
    <main className="container">
      <section className="state state--error">
        <h1>403</h1>
        <p>Anda tidak memiliki akses ke halaman ini.</p>
        <Link to="/products" className="button">
          Back
        </Link>
      </section>
    </main>
  )
}

export default Unauthorized
