package asl

// Err wraps any error received back from the API
type Err struct {
	Status int
	Msg    string
}

func (e *Err) Error() string { return e.Msg }
func (e *Err) Unwrap() error { return nil }
func (e *Err) Is(err error) bool {
	_, is := err.(*Err)
	return is
}
