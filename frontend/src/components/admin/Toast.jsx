function Toast({ message, type = 'success', onClose }) {
  if (!message) {
    return null
  }

  return (
    <div className={`toast toast--${type}`}>
      <span>{message}</span>
      <button type="button" onClick={onClose} aria-label="Close notification">
        x
      </button>
    </div>
  )
}

export default Toast
