package apperror

import "fmt"

type ErrorWrap struct {
	Err       error
	File      string
	Method    string
}

func (e *ErrorWrap) Error() string {
	errorString := fmt.Sprintf("error on %s: %s", e.File, e.Method)
	errorString += fmt.Sprintf(": %s", e.Err.Error())

	return errorString
}

func NewError(e error, fileName string, method string) *ErrorWrap {
	return &ErrorWrap{
		Err:       e,
		File:      fileName,
		Method:    method,
	}
}