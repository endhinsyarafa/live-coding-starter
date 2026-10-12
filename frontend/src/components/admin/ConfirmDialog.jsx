function ConfirmDialog({ confirmLabel, description, onCancel, onConfirm, title }) {
  if (!title) {
    return null
  }

  return (
    <div className="dialog-backdrop" role="presentation">
      <section className="dialog" role="dialog" aria-modal="true">
        <h2>{title}</h2>
        <p>{description}</p>
        <div className="form-actions">
          <button type="button" className="button" onClick={onCancel}>
            Cancel
          </button>
          <button type="button" className="button button--danger" onClick={onConfirm}>
            {confirmLabel}
          </button>
        </div>
      </section>
    </div>
  )
}

export default ConfirmDialog
