package domain

type Err string

func (e Err) Error() string {
	return string(e)
}

const (
	// Validation Errors
	ErrMissingEmail    = Err("email required")
	ErrExistingContact = Err("email must be unique")
)
