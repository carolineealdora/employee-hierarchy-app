package utils

import "github.com/carolineealdora/employee-hierarchy-app/internal/entities"

func LoopDetector(employee *entities.Employee, relations map[*entities.Employee][]*entities.Employee, visited map[*entities.Employee]bool, visitedStacks map[*entities.Employee]bool) bool {
	if !visited[employee] {
		visited[employee] = true
		visitedStacks[employee] = true

		directReport, isEmpRelationExists := relations[employee]
		for _, employee := range directReport{
			if isEmpRelationExists {
				if !visited[employee] && LoopDetector(employee, relations, visited, visitedStacks) {
					return true
				}
				if visitedStacks[employee] {
					return true
				}
			}
		}
	}

	delete(visitedStacks, employee)
	return false
}