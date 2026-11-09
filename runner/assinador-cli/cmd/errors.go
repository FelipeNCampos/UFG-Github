package cmd

type ExitError struct {
	Code    int
	Message string
}

func (e *ExitError) Error() string {
	if e == nil {
		return ""
	}
	if e.Message == "" {
		return "command failed"
	}
	return e.Message
}
