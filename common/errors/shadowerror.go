package errors

type ShadowError struct {
	VisibleMessage string
	Origin         error
}

// New returns an error that formats as the given text.
func New(message string, err error) error {
	return &ShadowError{
		VisibleMessage: message,
		Origin:         err,
	}
}

func (s *ShadowError) Error() string {
	return s.Origin.Error()
}

func (s *ShadowError) Message() string {
	return s.VisibleMessage
}
