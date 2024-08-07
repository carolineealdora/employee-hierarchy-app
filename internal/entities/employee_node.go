package entities

type EmployeeNode struct{
	Employee *Employee
	// Manager *EmployeeNode
	DirectReports []*EmployeeNode
}
