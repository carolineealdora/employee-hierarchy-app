package entities

type EmployeeNode struct {
	Employee      *Employee
	DirectReports []*EmployeeNode
}
