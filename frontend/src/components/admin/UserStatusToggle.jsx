function UserStatusToggle({ disabled, isActive, onChange, title }) {
  return (
    <button
      type="button"
      className={`switch ${isActive ? 'is-on' : ''}`}
      onClick={() => onChange(!isActive)}
      disabled={disabled}
      title={title}
      aria-label={isActive ? 'Deactivate user' : 'Activate user'}
    >
      <span>{isActive ? 'ON' : 'OFF'}</span>
    </button>
  )
}

export default UserStatusToggle
