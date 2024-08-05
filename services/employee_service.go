package services

import (
	"log"

	"github.com/carolineealdora/employee-hierarchy-app/apperror"
	"github.com/carolineealdora/employee-hierarchy-app/constants"
	"github.com/carolineealdora/employee-hierarchy-app/dtos"
	"github.com/carolineealdora/employee-hierarchy-app/entities"
	"github.com/carolineealdora/employee-hierarchy-app/repositories"
	"github.com/carolineealdora/employee-hierarchy-app/utils"
)

type EmployeeService interface {
	GenerateTree(employees []*entities.Employee) (*entities.EmployeeNode, error)
	SearchEmployee(empName string, empTree *entities.EmployeeNode) *dtos.FindEmployee
}

type employeeService struct {
	employeeRepository repositories.EmployeeRepository
}

func NewEmployeeService(er repositories.EmployeeRepository) *employeeService {
	return &employeeService{
		employeeRepository: er,
	}
}

func (s *employeeService) GenerateRelationMap(employees []*entities.Employee) (*entities.Employee, map[*entities.Employee][]*entities.Employee, error) {
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

		manager, err := s.employeeRepository.FindEmployeeByIdOnArrayData(e.ManagerId, employees)
		if err != nil {
			return nil, nil, err
		}

		child, parent := e, manager
		relations[parent] = append(relations[parent], child)
	}

	executive, err := s.FindExecutive(empWithNoManager, relations)

	if err != nil {
		return nil, nil, err
	}

	return executive, relations, nil
}

func (s *employeeService) FindExecutive(empWithNoManager []*entities.Employee, relations map[*entities.Employee][]*entities.Employee) (*entities.Employee, error) {

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

func (s *employeeService) BuildNode(parent *entities.Employee, employeeRelations map[*entities.Employee][]*entities.Employee) *entities.EmployeeNode {

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
			childNode := s.BuildNode(c, employeeRelations)
			if childNode != nil {
				node.DirectReports = append(node.DirectReports, childNode)
			}
			continue
		}
		node.DirectReports = append(node.DirectReports, &newNode)
	}
	return &node
}

func (s *employeeService) GenerateTree(employees []*entities.Employee) (*entities.EmployeeNode, error) {
	const methodName = "employeeService.GenerateTree"

	root, employeeRelations, err := s.GenerateRelationMap(employees)

	if err != nil {
		return nil, err
	}

	empTree := s.BuildNode(root, employeeRelations)

	return empTree, nil
}

func (s *employeeService) SearchEmployee(empName string, empTree *entities.EmployeeNode) (*dtos.FindEmployee, error) {
	const methodName = "employeeService.SearchEmployee"

	var queue []*entities.EmployeeNode
	queue = append(queue, empTree)

	// isEmpFound := false
	var foundEmp *entities.EmployeeNode
	var listManagers []string
	// var CountDirectReports int
	var CountIndirectReports int

	// for len(queue) > 0 {
	// 	node := queue[0]
	// 	queue = queue[1:]

	// 	if isEmpFound && foundEmp.Employee.ManagerId != node.Employee.ManagerId {
	// 		CountIndirectReports += len(node.DirectReports)
	// 	} // bisa aja pake yang lain juga

	// 	if node.Employee.Name == empName {
	// 		isEmpFound = true
	// 		foundEmp = *node
	// 		CountDirectReports = len(node.DirectReports)
	// 		queue = node.DirectReports
	// 	}

	// 	queue = append(queue, node.DirectReports...)
	// }

	foundEmp = s.SearchForEmployee(empName, empTree)

	if foundEmp == nil {
		return nil, apperror.NewError(
			apperror.EmployeeNotFoundError(empName),
			constants.EmployeeServFile,
			methodName,
		)
	}

	CountIndirectReports = s.CountIndirectReports(foundEmp, 0)

	listManagers = s.SearchForManagers(empTree, listManagers, foundEmp)
	log.Println(listManagers, "list")

	result := &dtos.FindEmployee{
		EmployeeName:         empName,
		Managers:             listManagers,
		CountDirectReports:   len(foundEmp.DirectReports),
		CountIndirectReports: CountIndirectReports,
	}

	return result, nil
}

func (s *employeeService) SearchForManagers(node *entities.EmployeeNode, listManagers []string, empToFind *entities.EmployeeNode) []string{
	if node != nil {
		if node.Employee.Name == empToFind.Employee.Name {
			return listManagers
		}

		listManagers = append(listManagers, node.Employee.Name)

		if len(node.DirectReports) > 0 {
			for _, c := range node.DirectReports {
				list := s.SearchForManagers(c, listManagers, empToFind)
				if list != nil {
					return list
				}
			}
		}
	}
	return nil
}

func (s *employeeService) SearchForEmployee(empName string, node *entities.EmployeeNode) *entities.EmployeeNode{
	if node != nil {
		if node.Employee.Name == empName {
			return node
		}

		if len(node.DirectReports) > 0 {
			for _, c := range node.DirectReports {
				empFound := s.SearchForEmployee(empName, c)
				if empFound != nil {
					return empFound
				}
			}
		}
	}
	return nil
}

func (s *employeeService) CountIndirectReports(parentNode *entities.EmployeeNode, count int) int{
	if parentNode != nil {
		if len(parentNode.DirectReports) > 0 {
			for _, c := range parentNode.DirectReports {
				count += len(c.DirectReports)
				s.CountIndirectReports(c, count)
			}
		}
	}
	return count
}