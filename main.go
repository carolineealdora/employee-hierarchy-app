package main

import (
	"log"

	"github.com/carolineealdora/employee-hierarchy-app/repositories"
	"github.com/carolineealdora/employee-hierarchy-app/services"
)

func main() {
	var repo = repositories.NewEmployeeRepository()
	var serv = services.NewEmployeeService(repo)

	populate, err := repo.PopulateEmployeeArrayData(1)

	if err != nil {
		log.Println(err)
	}

	// executive, employeeRelations, err := serv.GenerateRelationMap(populate)

	if err != nil {
		log.Println(err)
	}

	empTree, err := serv.GenerateTree(populate)

	if err != nil {
		log.Println(err)
	}

	res, err := serv.SearchEmployee("eveleen", empTree)

	if err != nil {
		log.Println(err)
	}

	log.Println(empTree, "tree")
	log.Println(res, "searched")

	// j, _ := json.Marshal(empTree)
	// log.Println(string(j), "result")
}
