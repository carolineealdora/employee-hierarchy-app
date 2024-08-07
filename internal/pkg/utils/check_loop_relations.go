package utils

import "github.com/carolineealdora/employee-hierarchy-app/internal/entities"

func LoopDetector(employee *entities.Employee, relations map[*entities.Employee][]*entities.Employee, visited map[*entities.Employee]bool, visitedStacks map[*entities.Employee]bool) bool {
	if !visited[employee] {
		visited[employee] = true
		visitedStacks[employee] = true

		directReport, ok := relations[employee]
		for _, employee := range directReport{
			if ok {
				if !visited[employee] && LoopDetector(employee, relations, visited, visitedStacks) {
					return true
				}
				if visitedStacks[employee] {
					return true
				}
			}
		}
	}

	
	// visitedStacks[employee] = false
	delete(visitedStacks, employee)
	return false
}