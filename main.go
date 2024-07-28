package main

import (
	"encoding/json"
	"log"

	"github.com/carolineealdora/employee-hierarchy-app/repositories"
)

func main() {
	var repo = repositories.NewEmployeeRepository()

	populate, err := repo.PopulateEmployeeArrayData(1)

	if err != nil {
		log.Println(err)
	}

	executive, employeeRelations, err := repo.GenerateRelationMap(populate)

	if err != nil {
		log.Println(err)
	}

	log.Print(executive, employeeRelations, "relation map")

	empTree, err := repo.GenerateTree(executive, employeeRelations)

	if err != nil {
		log.Println(err)
	}

	for _, v := range empTree {
		a := *v
		log.Print(a.Employee, "yesy")
		log.Print(len(a.DirectReports))
	}

	j, _ := json.Marshal(empTree)
	log.Println(string(j), "result")

	log.Print(empTree, "empTree")
	log.Println(populate, "populate")
}
