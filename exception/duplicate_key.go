package exception

type DuplicateKeyError struct {
	Error string
}

func NewDuplicateKeyError(error string) DuplicateKeyError {
	return DuplicateKeyError{Error: error}
}
