package dto


type (
	AcctDimensionData struct {
		ProjectID    *int64 `json:"project_id" required:"false"`
		CostCenterID *int64 `json:"cost_center_id" required:"false"`
	}
	AccountingDimensionDto struct {
		Project        *string `json:"project"`
		ProjectID      *int64  `json:"project_id"`
		ProjectUUID    *string `json:"project_uuid"`
		CostCenter     *string `json:"cost_center"`
		CostCenterID   *int64  `json:"cost_center_id"`
		CostCenterUUID *string `json:"cost_center_uuid"`
	}
)