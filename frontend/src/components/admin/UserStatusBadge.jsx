function UserStatusBadge({ isActive }) {
  return (
    <span className={`status-badge ${isActive ? 'status-badge--active' : 'status-badge--inactive'}`}>
      {isActive ? 'Active' : 'Inactive'}
    </span>
  )
}

export default UserStatusBadge
