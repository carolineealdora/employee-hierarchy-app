package repositories

import (
	"encoding/json"
	"os"

	"github.com/carolineealdora/employee-hierarchy-app/apperror"
	"github.com/carolineealdora/employee-hierarchy-app/constants"
	"github.com/carolineealdora/employee-hierarchy-app/entities"
)

type EmployeeRepository interface {
	PopulateEmployeeArrayData(dataSetType int) ([]*entities.Employee, error) 
	FindEmployeeByIdOnArrayData(id int, employees []*entities.Employee) (*entities.Employee, error) 
}

type employeeRepository struct {
}

func NewEmployeeRepository() *employeeRepository {
	return &employeeRepository{}
}

func (r *employeeRepository) PopulateEmployeeArrayData(dataSetType int) ([]*entities.Employee, error) {
	const methodName = "employeeRepository.PopulateEmployeeData"

	var filePath string
	switch dataSetType {
	case 1:
		filePath = "./json_data/correct-employees.json"
	case 2:
		filePath = "./json_data/faulty-employees-1.json"
	case 3:
		filePath = "./json_data/faulty-employees-2.json"
	default:
		filePath = "./json_data/correct-employees.json"
	}

	employeeDataJson, err := os.ReadFile(filePath)
	if err != nil {
		return nil, apperror.NewError(
			apperror.InternalServerError(),
			constants.EmployeeRepoFile,
			methodName,
		)
	}
	var employees []*entities.Employee
	if err := json.Unmarshal(employeeDataJson, &employees); err != nil {
		return nil, apperror.NewError(
			apperror.InternalServerError(),
			constants.EmployeeRepoFile,
			methodName,
		)
	}

	return employees, nil
}

func (r *employeeRepository) FindEmployeeByIdOnArrayData(id int, employees []*entities.Employee) (*entities.Employee, error) {
	const methodName = "employeeRepository.FindEmployeeByIdOnArrayData"
	
	for _, d := range employees {
		if id == d.Id {
			return d, nil
		}
	}
	return nil, apperror.NewError(
		apperror.InternalServerError(),
		constants.EmployeeRepoFile,
		methodName,
	)
}