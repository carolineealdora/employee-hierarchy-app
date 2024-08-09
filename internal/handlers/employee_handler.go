package handlers

import (
	"net/http"
	"strconv"

	"github.com/carolineealdora/employee-hierarchy-app/internal/constants"
	"github.com/carolineealdora/employee-hierarchy-app/internal/dtos"
	"github.com/carolineealdora/employee-hierarchy-app/internal/pkg/apperror"
	"github.com/carolineealdora/employee-hierarchy-app/internal/pkg/utils"
	"github.com/carolineealdora/employee-hierarchy-app/internal/services"
	"github.com/gin-gonic/gin"
)

type EmployeeHandler struct {
	employeeService services.EmployeeService
}

func NewEmployeeHandler(se services.EmployeeService) *EmployeeHandler {
	return &EmployeeHandler{
		employeeService: se,
	}
}

func (h *EmployeeHandler) GetEmployee(c *gin.Context) {
	const methodName = "employeeHandler.GetEmployee"
	var reqData dtos.SearchEmployeeReq
	
	if c.Query("name") == "" || c.Query("data_set_type") == "" {
		var requiredField []string
		if c.Query("name") == "" {
			requiredField = append(requiredField, "name")
		}

		if c.Query("data_set_type") == "" {
			requiredField = append(requiredField, "data_set_type")
		}

		apperror.NewError(
			apperror.RequriedFieldError(requiredField),
			constants.EmployeeServFile,
			methodName,
		)
		return
	}

	dataSetType, err := strconv.Atoi(c.Query("data_set_type"))
	empName := c.Query("name")
	reqData = dtos.SearchEmployeeReq{
		DataSetType: dataSetType,
		EmployeeName: empName,
	}

	if err != nil {
		c.Error(err)
		return
	}

	dataEmployee, err := h.employeeService.GetEmployeeByName(c, reqData)

	if err != nil {
		c.Error(err)
		return
	}

	resp := utils.GenerateResponse(constants.ResponseMsgDataRetrieved, dataEmployee)
	c.JSON(http.StatusOK, resp)
}
