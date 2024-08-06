package apperror

type ValidatorError struct{
	Field string `json:"field"`
	Message string `json:"message"`
}

func NewValidatorError(f, m string) *ValidatorError {
	return &ValidatorError{
		Field : f,
		Message: m,
	}
}

func (e *ValidatorError) Error() string{
	return e.Message
}