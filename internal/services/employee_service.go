package services

import (
	"context"
	"encoding/json"

	"github.com/carolineealdora/employee-hierarchy-app/internal/constants"
	"github.com/carolineealdora/employee-hierarchy-app/internal/dtos"
	"github.com/carolineealdora/employee-hierarchy-app/internal/entities"
	"github.com/carolineealdora/employee-hierarchy-app/internal/pkg/apperror"
	"github.com/carolineealdora/employee-hierarchy-app/internal/pkg/utils"
	"github.com/carolineealdora/employee-hierarchy-app/internal/repositories"
)

type EmployeeService interface {
	GetEmployeeByName(ctx context.Context, reqData dtos.SearchEmployeeReq) (*dtos.FindEmployee, error)
}

type employeeService struct {
	employeeRepository repositories.EmployeeRepository
}

func NewEmployeeService(er repositories.EmployeeRepository) *employeeService {
	return &employeeService{
		employeeRepository: er,
	}
}

func (s *employeeService) GenerateEmployeeData(ctx context.Context, dataSetType int) ([]*entities.Employee, error) {
	const methodName = "employeeService.GenerateEmployeeData"

	dataSet := s.employeeRepository.GetDataSetEmployee(ctx)

	setDataPath, ok := dataSet[dataSetType]
	log.Println(ok, "setDataPath", setDataPath)
	if !ok {
		return nil, apperror.NewError(
			apperror.RetrieveDataError("data_set_type"),
			constants.EmployeeServFile,
			methodName,
		)
	}

	setData, err := s.employeeRepository.PopulateEmployeeArrayData(ctx, setDataPath)
	if err != nil {
		return nil, err
	}
	log.Println("setData", setData)
	return setData, nil
}

func (s *employeeService) GenerateRelationMap(ctx context.Context, employees []*entities.Employee) (*entities.Employee, map[*entities.Employee][]*entities.Employee, error) {
	const methodName = "employeeService.GenerateRelationMap"
	relations := make(map[*entities.Employee][]*entities.Employee)
	empWithNoManager := []*entities.Employee{}

	isDuplicated, duplicatedData := utils.CheckDuplicateOnEmployeeSlice(employees)

	if isDuplicated {
		return nil, nil, apperror.NewError(
			apperror.MultipleManagerEmployeeError(duplicatedData),
			constants.EmployeeServFile,
			methodName,
		)
	}

	for _, e := range employees {
		if e.ManagerId == 0 {
			parent := e
			empWithNoManager = append(empWithNoManager, e)
			relations[parent] = append(relations[parent], nil)
			continue
		}

		if e.Id == e.ManagerId{
			return nil, nil, apperror.NewError(
				apperror.SelfManagerError(e.Name),
				constants.EmployeeServFile,
				methodName,
			)
		}

		manager, err := s.employeeRepository.FindEmployeeByIdOnArrayData(ctx, e.ManagerId, employees)
		if err != nil {
			return nil, nil, err
		}

		child, parent := e, manager
		relations[parent] = append(relations[parent], child)
	}

	executive, err := s.FindExecutive(ctx, empWithNoManager, relations)

	if err != nil {
		return nil, nil, err
	}

	return executive, relations, nil
}

func (s *employeeService) FindExecutive(ctx context.Context, empWithNoManager []*entities.Employee, relations map[*entities.Employee][]*entities.Employee) (*entities.Employee, error) {
	const methodName = "employeeService.FindExecutive"
	var empWithNoHierarchy []string
	var executive *entities.Employee
	for _, d := range empWithNoManager {
		_, ok := relations[d]

		if ok && executive != nil {
			return nil, apperror.NewError(
				apperror.MaximumExecutiveError(),
				constants.EmployeeServFile,
				methodName,
			)
		}

		if ok && executive == nil {
			executive = d
		}

		// Case : no executive exists
		if !ok && len(empWithNoManager) == 1 {
			return nil, apperror.NewError(
				apperror.NoExecutiveFoundError(),
				constants.EmployeeServFile,
				methodName,
			)
		}

		// Case : employee does not have hierarchy
		if !ok && len(empWithNoManager) > 1 {
			empWithNoHierarchy = append(empWithNoHierarchy, d.Name)
		}
	}

	// Exit early. Case : employee with no hierarchy
	if len(empWithNoHierarchy) > 0 {
		return nil, apperror.NewError(
			apperror.NoHierarchyEmployeeError(empWithNoHierarchy),
			constants.EmployeeServFile,
			methodName,
		)
	}

	return executive, nil
}

func (s *employeeService) BuildNode(ctx context.Context, parent *entities.Employee, employeeRelations map[*entities.Employee][]*entities.Employee) *entities.EmployeeNode {

	var node = entities.EmployeeNode{}
	node.Employee = parent

	if _, ok := employeeRelations[parent]; !ok {
		return nil
	}

	for _, c := range employeeRelations[parent] {
		if c == nil {
			continue
		}

		var newNode = entities.EmployeeNode{
			Employee: c,
		}

		if _, ok := employeeRelations[c]; ok {
			childNode := s.BuildNode(ctx, c, employeeRelations)
			if childNode != nil {
				node.DirectReports = append(node.DirectReports, childNode)
			}
			continue
		}
		node.DirectReports = append(node.DirectReports, &newNode)
	}
	return &node
}

func (s *employeeService) GenerateTree(ctx context.Context, empData []*entities.Employee) (*entities.EmployeeNode, error) {
	const methodName = "employeeService.GenerateTree"

	root, employeeRelations, err := s.GenerateRelationMap(ctx, empData)

	if err != nil {
		return nil, err
	}

	empTree := s.BuildNode(ctx, root, employeeRelations)

	if empTree == nil {
		return nil, apperror.NewError(
			apperror.FailedOnGeneratingTreeError(),
			constants.EmployeeServFile,
			methodName,
		)
	}

	return empTree, nil
}

func (s *employeeService) SearchEmployee(ctx context.Context, empName string, empTree *entities.EmployeeNode) (*dtos.FindEmployee, error) {
	const methodName = "employeeService.SearchEmployee"

	var listManagers []string

	foundEmp := s.findEmployeeByNameOnTree(ctx, empName, empTree)

	if foundEmp == nil {
		return nil, apperror.NewError(
			apperror.EmployeeNotFoundError(empName),
			constants.EmployeeServFile,
			methodName,
		)
	}

	CountIndirectReports := s.CountIndirectReports(ctx, foundEmp, 0)

	listManagers = s.SearchForManagers(ctx, empTree, listManagers, foundEmp)

	result := &dtos.FindEmployee{
		EmployeeName:         empName,
		Managers:             listManagers,
		CountDirectReports:   len(foundEmp.DirectReports),
		CountIndirectReports: CountIndirectReports,
	}

	return result, nil
}

func (s *employeeService) SearchForManagers(ctx context.Context, node *entities.EmployeeNode, listManagers []string, empToFind *entities.EmployeeNode) []string {
	if node != nil {
		if node.Employee.Name == empToFind.Employee.Name {
			return listManagers
		}

		listManagers = append(listManagers, node.Employee.Name)

		if len(node.DirectReports) > 0 {
			for _, c := range node.DirectReports {
				list := s.SearchForManagers(ctx, c, listManagers, empToFind)
				if list != nil {
					return list
				}
			}
		}
	}
	return nil
}

func (s *employeeService) findEmployeeByNameOnTree(ctx context.Context, empName string, node *entities.EmployeeNode) *entities.EmployeeNode {
	if node != nil {
		if node.Employee.Name == empName {
			return node
		}

		if len(node.DirectReports) > 0 {
			for _, c := range node.DirectReports {
				empFound := s.findEmployeeByNameOnTree(ctx, empName, c)
				if empFound != nil {
					return empFound
				}
			}
		}
	}
	return nil
}

func (s *employeeService) CountIndirectReports(ctx context.Context, parentNode *entities.EmployeeNode, count int) int {
	if parentNode != nil {
		if len(parentNode.DirectReports) > 0 {
			for _, c := range parentNode.DirectReports {
				count += len(c.DirectReports)
				s.CountIndirectReports(ctx, c, count)
			}
		}
	}
	return count
}

func (s *employeeService) GetEmployeeByName(ctx context.Context, reqData dtos.SearchEmployeeReq) (*dtos.FindEmployee, error) {
	const methodName = "employeeService.GetEmployeeByName"
	
	dataSet, err := s.GenerateEmployeeData(ctx, reqData.DataSetType)

	if err != nil {
		return nil, err
	}

	if len(dataSet) == 0 {
		return nil, apperror.NewError(
			apperror.EmptyDataSetError(),
			constants.EmployeeRepoFile,
			methodName,
		)
	}

	empTree, err := s.GenerateTree(ctx, dataSet)

	if err != nil {
		return nil, err
	}

	emp, err := s.SearchEmployee(ctx, reqData.EmployeeName, empTree)

	if err != nil {
		return nil, err
	}

	return emp, nil
}
