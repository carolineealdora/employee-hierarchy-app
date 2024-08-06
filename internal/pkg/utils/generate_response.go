package utils

import "github.com/carolineealdora/employee-hierarchy-app/internal/dtos"

func GenerateResponse(msg string, data any) *dtos.ResponseMessage {
	return &dtos.ResponseMessage{
		Message: msg,
		Data:    data,
	}
}
