package app

// simpleResponse is just a handy way to collect errors to return in a
// JSON payload.
type simpleResponse struct {
	ObjID  string
	Errors []string
}

// newSimpleResponse creates a simpleResponse with its Errors slice already allocated.
func newSimpleResponse() *simpleResponse {
	return &simpleResponse{Errors: []string{}}
}

// AddError appends an error to a simpleResponses Errors slice.
func (s *simpleResponse) AddError(err error) {
	s.Errors = append(s.Errors, err.Error())
}
