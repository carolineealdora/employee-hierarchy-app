package utils

import (

	"github.com/carolineealdora/employee-hierarchy-app/internal/entities"
)

func CheckDuplicateOnEmployeeSlice(data []*entities.Employee) (bool, map[string][]string) {
	nameCount := make(map[string]int)
	empMultipleEntry := make(map[string][]string)

	for _, d := range data {
		nameCount[d.Name]++
	}

	for name, count := range nameCount {
		if count > 1 {
			for _, d := range data {
				if d.Name == name {
					manager := FindEmpNameByID(data, d.ManagerId)
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

func FindEmpNameByID(data []*entities.Employee, id int) string {
	for _, d := range data {
		if d.Id == id {
			return d.Name
		}
	}
	return ""
}
