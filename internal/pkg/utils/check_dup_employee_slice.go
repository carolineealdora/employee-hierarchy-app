package utils

import (

	"github.com/carolineealdora/employee-hierarchy-app/internal/entities"
)

func CheckDuplicateOnEmployeeSlice(employees []*entities.Employee) (bool, map[string][]string) {
	nameCount := make(map[string]int)
	empMultipleEntry := make(map[string][]string)

	for _, emp := range employees {
		nameCount[emp.Name]++
	}

	for name, count := range nameCount {
		if count > 1 {
			for _, emp := range employees {
				if emp.Name == name {
					manager := FindEmpNameByID(employees, emp.ManagerId)
					empMultipleEntry[name] = append(empMultipleEntry[name], manager)
				}
			}
		}
	}

	if len(empMultipleEntry) > 0 {
		return true, empMultipleEntry
	}

	return false, nil
}

func FindEmpNameByID(employees []*entities.Employee, id int) string {
	for _, emp := range employees {
		if emp.Id == id {
			return emp.Name
		}
	}
	return ""
}
