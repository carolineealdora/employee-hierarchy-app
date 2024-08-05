package dtos

type FindEmployee struct {
	EmployeeName string `json:"name"`
	Managers []string `json:"managers"`
	CountDirectReports int `json:"total_direct_reports"`
	CountIndirectReports int `json:"total_indirect_reports"`
}

type SearchEmployeeReq struct {
	EmployeeName string `json:"name" required:"binding"`
	DataSetType int `json:"data_set_type" required:"binding"`
}