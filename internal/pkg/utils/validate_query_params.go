package utils

import (
	"github.com/carolineealdora/employee-hierarchy-app/internal/constants"
	"github.com/carolineealdora/employee-hierarchy-app/internal/pkg/apperror"
	"github.com/gin-gonic/gin"
)

func ValidateSearchForEmpQueryParams(c *gin.Context) error {
	const methodName = "validateQueryParamsUtils.ValidateSearchForEmpQueryParams"
	if c.Query("name") == "" || c.Query("data_set_type") == "" {
		var requiredField []string
		if c.Query("name") == "" {
			requiredField = append(requiredField, "name")
		}

		if c.Query("data_set_type") == "" {
			requiredField = append(requiredField, "data_set_type")
		}

		return apperror.NewError(
			apperror.RequriedFieldError(requiredField),
			constants.EmployeeServFile,
			methodName,
		)
	}
	return nil
}
