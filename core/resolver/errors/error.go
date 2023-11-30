package errors

type ResolverError struct {
	HTTPCode int
	Err      error
}

func (r *ResolverError) Error() string {
	return r.Err.Error()
}
