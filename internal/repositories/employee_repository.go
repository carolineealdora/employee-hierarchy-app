package repositories

import (
	"encoding/json"
	"os"

	"github.com/carolineealdora/employee-hierarchy-app/internal/entities"
	"github.com/carolineealdora/employee-hierarchy-app/internal/apperror"
	"github.com/carolineealdora/employee-hierarchy-app/internal/constants"
)

type EmployeeRepository interface {
	GetDataSetEmployee() map[int]string
	PopulateEmployeeArrayData(filePath string) ([]*entities.Employee, error)
	FindEmployeeByIdOnArrayData(id int, employees []*entities.Employee) (*entities.Employee, error)
}

type employeeRepository struct {
}

func NewEmployeeRepository() *employeeRepository {
	return &employeeRepository{}
}

func (r *employeeRepository) PopulateEmployeeArrayData(filePath string) ([]*entities.Employee, error) {
	const methodName = "employeeRepository.PopulateEmployeeData"

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

func (r *employeeRepository) GetDataSetEmployee() map[int]string {
	dataSetEmployee := map[int]string{
		1: "./json_data/correct-employees.json",
		2: "./json_data/faulty-employees-1.json",
		3: "./json_data/faulty-employees-2.json",
	}

	return dataSetEmployee
}
