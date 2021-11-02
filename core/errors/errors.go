package errors

type DestinationError struct {
	err error
}

func (r *DestinationError) Error() string {
	return r.err.Error()
}
func NewDestinationError(err error) *DestinationError {
	if err == nil {
		return nil
	}
	return &DestinationError{err: err}
}
func IsDestinationError(err error) bool {
	switch err.(type) {
	case *DestinationError:
		return true
	default:
		return false
	}
}
