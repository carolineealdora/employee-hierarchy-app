// Code generated by mockery v2.14.0. DO NOT EDIT.

package mocks

import (
	context "context"

	dtos "github.com/carolineealdora/employee-hierarchy-app/internal/dtos"
	mock "github.com/stretchr/testify/mock"
)

// EmployeeService is an autogenerated mock type for the EmployeeService type
type EmployeeService struct {
	mock.Mock
}

// GetEmployeeByName provides a mock function with given fields: ctx, reqData
func (_m *EmployeeService) GetEmployeeByName(ctx context.Context, reqData dtos.SearchEmployeeReq) (*dtos.FindEmployee, error) {
	ret := _m.Called(ctx, reqData)

	var r0 *dtos.FindEmployee
	if rf, ok := ret.Get(0).(func(context.Context, dtos.SearchEmployeeReq) *dtos.FindEmployee); ok {
		r0 = rf(ctx, reqData)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*dtos.FindEmployee)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, dtos.SearchEmployeeReq) error); ok {
		r1 = rf(ctx, reqData)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTNewEmployeeService interface {
	mock.TestingT
	Cleanup(func())
}

// NewEmployeeService creates a new instance of EmployeeService. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewEmployeeService(t mockConstructorTestingTNewEmployeeService) *EmployeeService {
	mock := &EmployeeService{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
