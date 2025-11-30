package status

// Status represents lifecycle states for LCP-managed resources.
type Status string

const (
	StatusActive   Status = "active"
	StatusInactive Status = "inactive"
)
