package services

import "github.com/carolineealdora/employee-hierarchy-app/repositories"

type EmployeeService interface {

}

type employeeService struct {
	employeeRepository        repositories.EmployeeRepository
}

func NewEmployeeService(er repositories.EmployeeRepository) *employeeService {
	return &employeeService{
		employeeRepository: er,
	}
}