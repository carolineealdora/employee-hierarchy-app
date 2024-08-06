package handlers

import (
	"net/http"

	"github.com/carolineealdora/employee-hierarchy-app/internal/constants"
	"github.com/carolineealdora/employee-hierarchy-app/internal/dtos"
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
	var reqData dtos.SearchEmployeeReq
	err := c.ShouldBindJSON(&reqData)
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
