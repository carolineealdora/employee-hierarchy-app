package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/carolineealdora/employee-hierarchy-app/internal/dtos"
	"github.com/carolineealdora/employee-hierarchy-app/internal/pkg/apperror"
	"github.com/carolineealdora/employee-hierarchy-app/internal/pkg/utils"

	// "github.com/carolineealdora/employee-hierarchy-app/internal/pkg/apperror"
	"github.com/carolineealdora/employee-hierarchy-app/mocks"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestNewEmployeeHandler(t *testing.T) {
	dep := NewEmployeeHandler(nil)
	assert.NotNil(t, dep)
}

func TestGetEmployee(t *testing.T) {
	mockService := mocks.NewEmployeeService(t)
	dep := NewEmployeeHandler(mockService)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	t.Run("should return status OK when get data successful", func(t *testing.T) {
		req := dtos.SearchEmployeeReq{
			EmployeeName: "kacie",
			DataSetType:  1,
		}
		reqJson, _ := json.Marshal(req)
		expectedResObj := &dtos.FindEmployee{
			EmployeeName:         "kacie",
			Managers:             []string{"raelynn"},
			CountDirectReports:   3,
			CountIndirectReports: 1,
		}
		expectedResString := "{\"message\":\"successfully retrieved data\",\"data\":{\"employee_name\":\"kacie\",\"managers\":[\"raelynn\"],\"total_direct_reports\":3,\"total_indirect_reports\":1}}"

		gin.SetMode(gin.TestMode)
		r := httptest.NewRequest(http.MethodPost, "/search-employee", strings.NewReader(string(reqJson)))
		c.Request = r

		mockService.On("GetEmployeeByName", c, req).Return(expectedResObj, nil)
		dep.GetEmployee(c)

		assert.Equal(t, http.StatusOK, w.Result().StatusCode)
		assert.Equal(t, expectedResString, w.Body.String())
	})

	t.Run("should return error response when form required binding criteria not fulfilled", func(t *testing.T) {
		req := dtos.SearchEmployeeReq{
			EmployeeName: "kacie",
		}
		reqJson, _ := json.Marshal(req)
		expResObj := []apperror.ValidatorError{
			*apperror.NewValidatorError("DataSetType", "field DataSetType is required"),
		}
		expectedResString := "{\"message\":\"invalid request. please fill all field(s) with the correct data\",\"data\":[{\"field\":\"DataSetType\",\"message\":\"field DataSetType is required\"}]}"

		c.AbortWithStatusJSON(http.StatusBadRequest, utils.GenerateResponse(apperror.ErrValidator.Error(), expResObj))
		r := httptest.NewRequest(http.MethodPost, "/search-employee", strings.NewReader(string(reqJson)))
		c.Request = r

		dep.GetEmployee(c)

		assert.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
		assert.Equal(t, expectedResString, w.Body.String())
	})
}