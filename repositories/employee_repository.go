package repositories

import (
	"encoding/json"
	"log"
	"os"

	"github.com/carolineealdora/employee-hierarchy-app/apperror"
	"github.com/carolineealdora/employee-hierarchy-app/constants"
	"github.com/carolineealdora/employee-hierarchy-app/entities"
)

type EmployeeRepository interface {
	PopulateEmployeeArrayData(dataSetType int) ([]*entities.Employee, error) 
	FindEmployeeByIdOnArrayData(id int, employees []*entities.Employee) (*entities.Employee, error) 
	BuildNode(parent *entities.Employee, employeeRelations map[*entities.Employee][]*entities.Employee) (*entities.EmployeeNode, error) 
	FindExecutive(employees []*entities.Employee) (*entities.Employee, error)
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

func (r *employeeRepository) GenerateRelationMap(employees []*entities.Employee) (*entities.Employee, map[*entities.Employee][]*entities.Employee, error) {
	const methodName = "employeeRepository.GenerateRelationMap"
	relations := make(map[*entities.Employee][]*entities.Employee)
	empWithNoManager := []*entities.Employee{}

	for _, e := range employees {
		if e.ManagerId == 0 {
			parent := e
			empWithNoManager = append(empWithNoManager, e)
			relations[parent] = append(relations[parent],  nil)
			continue
		}

		manager, err := r.FindEmployeeByIdOnArrayData(e.ManagerId, employees)
		if err != nil {
			return nil, nil, err
		}

		child, parent := e, manager
		relations[parent] = append(relations[parent], child)
	}

	var empWithNoHierarchy []string
	var executive []*entities.Employee
	for _, d := range empWithNoManager{
		_, ok := relations[d]

		if ok {
			executive = append(executive, d)
		}

		if !ok && len(empWithNoManager) == 1 {
			return nil, nil, apperror.NewError(
				apperror.NoExecutiveFoundError(),
				constants.EmployeeRepoFile,
				methodName,
			)
		}

		if !ok && len(empWithNoManager) > 1 {
			empWithNoHierarchy = append(empWithNoHierarchy, d.Name)
		}
	}

	if len(empWithNoHierarchy) > 0 {
		return nil, nil, apperror.NewError(
			apperror.NoHierarchyEmployeeError(empWithNoHierarchy),
			constants.EmployeeRepoFile,
			methodName,
		)
	}

	if len(executive) > 1 {
		return nil, nil, apperror.NewError(
			apperror.MaximumExecutiveError(),
			constants.EmployeeRepoFile,
			methodName,
		)
	}

	// dup emp with same manager
	
	return executive[0], relations, nil
}

func(r *employeeRepository) GenerateTree(root *entities.Employee, employeeRelations map[*entities.Employee][]*entities.Employee) (entities.EmployeeNode, error){
	const methodName = "employeeRepository.GenerateTree"
	
	empTree, err := r.BuildNode(root, employeeRelations)

	if err!=nil{
		log.Print("failed")
	}

	return *empTree, nil
}

func (r *employeeRepository) BuildNode(parent *entities.Employee, employeeRelations map[*entities.Employee][]*entities.Employee) *entities.EmployeeNode {

	var node = entities.EmployeeNode{}
	
	for _, c := range employeeRelations[parent] {
		node.Employee = parent
		if c != nil {
			var childNode = entities.EmployeeNode{
				Employee: c,
			}

			if _, ok := employeeRelations[c]; ok {
				childNode := r.BuildNode(c, employeeRelations)

				node.DirectReports = append(node.DirectReports, childNode)
			}
			node.DirectReports = append(node.DirectReports, &childNode)
		}
	}
	return &node
}
