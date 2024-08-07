package repositories

import (
	"context"
	"encoding/json"
	"os"

	"github.com/carolineealdora/employee-hierarchy-app/internal/entities"
	"github.com/carolineealdora/employee-hierarchy-app/internal/pkg/apperror"

	"github.com/carolineealdora/employee-hierarchy-app/internal/constants"
)

type EmployeeRepository interface {
	GetDataSetEmployee(ctx context.Context) map[int]string 
	PopulateEmployeeArrayData(ctx context.Context, filePath string) ([]*entities.Employee, error) 
	FindEmployeeByIdOnArrayData(ctx context.Context, id int, employees []*entities.Employee) (*entities.Employee, error) 
}

type employeeRepository struct {
}

func NewEmployeeRepository() *employeeRepository {
	return &employeeRepository{}
}

func (r *employeeRepository) PopulateEmployeeArrayData(ctx context.Context, filePath string) ([]*entities.Employee, error) {
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
	err = json.Unmarshal(employeeDataJson, &employees)

	if err != nil {
		return nil, apperror.NewError(
			apperror.InternalServerError(),
			constants.EmployeeRepoFile,
			methodName,
		)
	}
	
	return employees, nil
}

func (r *employeeRepository) FindEmployeeByIdOnArrayData(ctx context.Context, id int, employees []*entities.Employee) (*entities.Employee, error) {
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

func (r *employeeRepository) GetDataSetEmployee(ctx context.Context) map[int]string {
	dataSetEmployee := map[int]string{
		1: "./internal/json_data/correct-employees.json",
		2: "./internal/json_data/faulty-employees-1.json",
		3: "./internal/json_data/faulty-employees-2.json",
		4: "./internal/json_data/faulty-employees-empty-data-set.json",
		5: "./internal/json_data/faulty-employees-looped-relations.json",
		6: "./internal/json_data/faulty-employee-their-own-manager.json",
	}

	return dataSetEmployee
}
