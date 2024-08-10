package services

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/carolineealdora/employee-hierarchy-app/internal/constants"
	"github.com/carolineealdora/employee-hierarchy-app/internal/dtos"
	"github.com/carolineealdora/employee-hierarchy-app/internal/pkg/apperror"
	"github.com/carolineealdora/employee-hierarchy-app/internal/pkg/utils"
	"github.com/carolineealdora/employee-hierarchy-app/mocks"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestNewEmployeeService(t *testing.T) {
	dep := NewEmployeeService(nil)
	assert.NotNil(t, dep)
}

func TestGenerateEmployeeData(t *testing.T) {
	mockRepo := mocks.NewEmployeeRepository(t)
	dep := NewEmployeeService(mockRepo)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	t.Run("should return nil for the array of employee data with error when get data failed caused by invalid set data given by user", func(t *testing.T) {
		req := dtos.SearchEmployeeReq{
			EmployeeName: "kacie",
			DataSetType:  4,
		}
		reqJson, _ := json.Marshal(req)

		customErr := apperror.RetrieveDataError("data_set_type")
		resObj := apperror.NewError(
			customErr,
			constants.EmployeeServFile,
			"employeeService.GenerateEmployeeData",
		)
		expResObj := errors.New(resObj.Error())

		expectedResString := "{\"message\":\"failed while retrieving data\",\"data\":\"data_set_type\"}"

		gin.SetMode(gin.TestMode)
		r := httptest.NewRequest(http.MethodPost, "/search-employee", strings.NewReader(string(reqJson)))
		c.Request = r

		mockRepo.On("GetDataSetEmployee", c).Return(nil, expResObj)
		c.AbortWithStatusJSON(http.StatusBadRequest, utils.GenerateResponse(customErr.Message, customErr.Details))
		dep.GenerateEmployeeData(c, req.DataSetType)

		mockRepo.AssertExpectations(t)
		assert.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
		assert.Equal(t, expectedResString, w.Body.String())
	})

	t.Run("should return error when failed on populating array data of employee", func(t *testing.T) {
		req := dtos.SearchEmployeeReq{
			EmployeeName: "kacie",
			DataSetType:  1,
		}
		reqJson, _ := json.Marshal(req)

		customErr := apperror.InternalServerError()
		resObj := apperror.NewError(
			customErr,
			constants.EmployeeServFile,
			"employeeService.GenerateEmployeeData",
		)
		expResObj := errors.New(resObj.Error())
		falseDataPath := "./??"
		expDataSetEmp := map[int]string{
			1: falseDataPath,
		}

		gin.SetMode(gin.TestMode)
		r := httptest.NewRequest(http.MethodPost, "/search-employee", strings.NewReader(string(reqJson)))
		c.Request = r

		mockRepo.On("GetDataSetEmployee", c).Return(expDataSetEmp, nil)
		mockRepo.On("PopulateEmployeeArrayData", c, falseDataPath).Return(nil, expResObj)
		dep.GenerateEmployeeData(c, req.DataSetType)

		mockRepo.AssertNumberOfCalls(t, "GetDataSetEmployee", 1)
		mockRepo.AssertNumberOfCalls(t, "PopulateEmployeeArrayData", 1)
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

func TestGetEmployeeByName(t *testing.T) {
	mockRepo := mocks.NewEmployeeRepository(t)
	mockServ := mocks.EmployeeService(*mockRepo)
	dep := NewEmployeeService(mockRepo)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	t.Run("should return error when generating data failed caused by invalid data set", func(t *testing.T) {
		req := dtos.SearchEmployeeReq{
			EmployeeName: "kacie",
			DataSetType:  10,
		}
		reqJson, _ := json.Marshal(req)

		customErr := apperror.RetrieveDataError("data_set_type")
		resObj := apperror.NewError(
			customErr,
			constants.EmployeeServFile,
			"employeeService.GenerateEmployeeData",
		)
		expResObj := errors.New(resObj.Error())

		gin.SetMode(gin.TestMode)
		r := httptest.NewRequest(http.MethodPost, "/search-employee", strings.NewReader(string(reqJson)))
		c.Request = r

		mockRepo.On("GetDataSetEmployee", c).Return(nil, expResObj)
		mockServ.On("GenerateEmployeeData", c, req.DataSetType).Return(nil, expResObj)
		c.AbortWithStatusJSON(http.StatusBadRequest, utils.GenerateResponse(customErr.Message, customErr.Details))
		dep.GetEmployeeByName(c, req)

		mockRepo.AssertExpectations(t)
		assert.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
	})
}