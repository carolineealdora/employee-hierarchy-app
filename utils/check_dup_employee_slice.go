package utils

import "github.com/carolineealdora/employee-hierarchy-app/entities"

func CheckDuplicateOnEmployeeSlice(data []*entities.Employee) (bool, []string){

	nameMapped := make(map[string]int)
	var empMultipleEntry []string
	for _, d := range data {
		_, ok := nameMapped[d.Name]

		if !ok{
			nameMapped[d.Name] = 1
			continue
		}

		if ok {
			empMultipleEntry = append(empMultipleEntry, d.Name)
		}
	}

	if len(empMultipleEntry) > 0 {
		return true, empMultipleEntry
	}
	return false, nil
}