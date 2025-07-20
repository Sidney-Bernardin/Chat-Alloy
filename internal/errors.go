package internal

type DomainError struct {
	Type  string         `json:"type"`
	Msg   string         `json:"message,omitempty"`
	Attrs map[string]any `json:"attrs,omitempty"`
}
