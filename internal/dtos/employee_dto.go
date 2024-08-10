package dtos

type FindEmployee struct {
	EmployeeName         string   `json:"employee_name"`
	Managers             []string `json:"managers"`
	CountDirectReports   int      `json:"total_direct_reports"`
	CountIndirectReports int      `json:"total_indirect_reports"`
}

type SearchEmployeeReq struct {
	EmployeeName string `json:"name" binding:"required"`
	DataSetType  int    `json:"data_set_type" binding:"required"`
}
