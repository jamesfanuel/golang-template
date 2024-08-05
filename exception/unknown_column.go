package exception

type UnknownColumnError struct {
	Error string
}

func NewUnknownColumnError(error string) UnknownColumnError {
	return UnknownColumnError{Error: error}
}
