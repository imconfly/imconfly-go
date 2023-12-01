package resolver_errors

type Error struct {
	HTTPCode int
	Err      error
}

func New(httpCode int, err error) *Error {
	return &Error{
		HTTPCode: httpCode,
		Err:      err,
	}
}

func (r *Error) Error() string {
	return r.Err.Error()
}
