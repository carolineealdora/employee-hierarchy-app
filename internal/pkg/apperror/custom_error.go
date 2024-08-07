package apperror

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CustomError struct {
	Code    int    `json:"-"`
	Name    string `json:"name"`
	Message string `json:"message"`
	Details any    `json:"details,omitempty"`
}

func (e *CustomError) Error() string {
	errorString := fmt.Sprintf("%s: %s.", e.Name, e.Message)
	if e.Details != nil {
		errorString += fmt.Sprintf("Details: %s", e.Details)
	}
	return errorString
}

func NotFoundError(field string, id int64) *CustomError {
	return &CustomError{
		Code:    http.StatusNotFound,
		Name:    "ErrNotFound",
		Message: "resource not found",
		Details: gin.H{
			field: id,
		},
	}
}

func NoExecutiveFoundError() *CustomError {
	return &CustomError{
		Code:    http.StatusBadRequest,
		Name:    "ErrNoExecutiveFound",
		Message: "Unable to process employeee hierarchy. Employee with executive role not found",
	}
}

func MaximumExecutiveError() *CustomError {
	return &CustomError{
		Code:    http.StatusBadRequest,
		Name:    "ErrMaximumExecutive",
		Message: "Executive reached maximum number(1)",
	}
}

func NoHierarchyEmployeeError(empName []string) *CustomError {
	return &CustomError{
		Code:    http.StatusInternalServerError,
		Name:    "ErrNoHierarchyEmployee",
		Message: "Unable to process employeee hierarchy. Employee(s) does not have any hierarchy",
		Details: gin.H{
			"list_employee_name": empName,
		},
	}
}

func MultipleManagerEmployeeError(empName []string) *CustomError {
	return &CustomError{
		Code:    http.StatusInternalServerError,
		Name:    "ErrMultipleManager",
		Message: "Unable to process employeee hierarchy. Employee(s) has multiple managers",
		Details: gin.H{
			"employee_name": empName,
		},
	}
}

func EmployeeNotFoundError(empName string) *CustomError {
	return &CustomError{
		Code:    http.StatusBadRequest,
		Name:    "ErrEmployeeNotFound",
		Message: "Employee with given name not found",
		Details: gin.H{
			"employee_name": empName,
		},
	}
}

func SelfManagerError(empName string) *CustomError {
	return &CustomError{
		Code:    http.StatusBadRequest,
		Name:    "ErrSelfManager",
		Message: "Failed on generating tree. Employee is their own manager",
		Details: gin.H{
			"employee_name": empName,
		},
	}
}

func EmptyDataSetError() *CustomError {
	return &CustomError{
		Code:    http.StatusBadRequest,
		Name:    "ErrEmptyDataSet",
		Message: "Failed on generating tree. Data set is empty",
	}
}

func LoopRelationError() *CustomError {
	return &CustomError{
		Code:    http.StatusBadRequest,
		Name:    "ErrLoopRelation",
		Message: "Failed on generating tree. Loop detected in data relation",
	}
}

func FailedOnGeneratingTreeError() *CustomError {
	return &CustomError{
		Code:    http.StatusBadRequest,
		Name:    "ErrEmployeeNotFound",
		Message: "Failed on generating tree. Please try again.",
	}
}

func RetrieveDataError(dataName string) *CustomError {
	return &CustomError{
		Code:    http.StatusInternalServerError,
		Name:    "ErrFailedRetrieveData",
		Message: "failed while retrieving data",
		Details: dataName,
	}
}

func ErrPageNotFoundError() *CustomError {
	return &CustomError{
		Code:    http.StatusNotFound,
		Name:    "ErrPageNotFound",
		Message: "page not found",
	}
}

func InternalServerError() *CustomError {
	return &CustomError{
		Code:    http.StatusInternalServerError,
		Name:    "ErrInternalServer",
		Message: "our server encounter error. please try again.",
	}
}
