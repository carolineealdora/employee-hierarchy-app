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

	t.Run("should return error internal server when failed on populate array data of employee", func(t *testing.T) {
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

		gin.SetMode(gin.TestMode)
		r := httptest.NewRequest(http.MethodPost, "/search-employee", strings.NewReader(string(reqJson)))
		c.Request = r

		mockRepo.On("GetDataSetEmployee", c).Return(falseDataPath, nil)
		mockRepo.On("PopulateEmployeeArrayData", falseDataPath).Return(nil, expResObj)
		dep.GenerateEmployeeData(c, req.DataSetType)

		mockRepo.AssertExpectations(t)
		assert.Equal(t, http.StatusInternalServerError, w.Result().StatusCode)
	})
}

// t.Run("should return error response when form required binding criteria not fulfilled", func(t *testing.T) {
// 	req := dtos.SearchEmployeeReq{
// 		EmployeeName: "kacie",
// 	}
// 	reqJson, _ := json.Marshal(req)
// 	expResObj := []apperror.ValidatorError{
// 		*apperror.NewValidatorError("DataSetType", "field DataSetType is required"),
// 	}
// 	expectedResString := "{\"message\":\"invalid request. please fill all field(s) with the correct data\",\"data\":[{\"field\":\"DataSetType\",\"message\":\"field DataSetType is required\"}]}"

// 	c.AbortWithStatusJSON(http.StatusBadRequest, utils.GenerateResponse(apperror.ErrValidator.Error(), expResObj))
// 	r := httptest.NewRequest(http.MethodPost, "/search-employee", strings.NewReader(string(reqJson)))
// 	c.Request = r

// 	dep.GetEmployee(c)

// 	assert.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
// 	assert.Equal(t, expectedResString, w.Body.String())
// })

// t.Run("should return error response when data set given is false", func(t *testing.T) {
// 	req := dtos.SearchEmployeeReq{
// 		EmployeeName: "kacie",
// 		DataSetType:  4,
// 	}
// 	reqJson, _ := json.Marshal(req)

// 	customErr := apperror.RetrieveDataError("data_set_type")
// 	resObj := apperror.NewError(
// 		customErr,
// 		filename,
// 		"employeeService.GenerateEmployeeData",
// 	)
// 	expResObj := errors.New(resObj.Error())

// 	expectedResString := "{\"message\":\"failed while retrieving data\",\"data\":\"data_set_type\"}"

// 	gin.SetMode(gin.TestMode)
// 	r := httptest.NewRequest(http.MethodPost, "/search-employee", strings.NewReader(string(reqJson)))
// 	c.Request = r

// 	mockService.On("GetEmployeeByName", c, req).Return(nil, expResObj)
// 	c.AbortWithStatusJSON(http.StatusBadRequest, utils.GenerateResponse(customErr.Message, customErr.Details))
// 	dep.GetEmployee(c)

// 	assert.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
// 	assert.Equal(t, expectedResString, w.Body.String())
// })

// t.Run("should return error response when data employee duplicated", func(t *testing.T) {
// 	req := dtos.SearchEmployeeReq{
// 		EmployeeName: "kacie",
// 		DataSetType:  3,
// 	}
// 	reqJson, _ := json.Marshal(req)

// 	customErr := apperror.FailedOnGeneratingTreeError()
// 	resObj := apperror.NewError(
// 		customErr,
// 		filename,
// 		"employeeService.GenerateTree",
// 	)
// 	expResObj := errors.New(resObj.Error())

// 	expectedResString := "{\"message\":\"Unable to process employeee hierarchy. Employee(s) has multiple managers\",\"data\":{\"employee_name\":[\"linton\"]}}"

// 	gin.SetMode(gin.TestMode)
// 	r := httptest.NewRequest(http.MethodPost, "/search-employee", strings.NewReader(string(reqJson)))
// 	c.Request = r

// 	mockService.On("GetEmployeeByName", c, req).Return(nil, expResObj)
// 	c.AbortWithStatusJSON(http.StatusBadRequest, utils.GenerateResponse(customErr.Message, customErr.Details))
// 	dep.GetEmployee(c)

// 	assert.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
// 	assert.Equal(t, expectedResString, w.Body.String())
// })
