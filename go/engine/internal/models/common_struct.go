package models

type ApiResponse struct {
	Status     string      `json:"status"`
	Message    string      `json:"message"`
	Data       interface{} `json:"data,omitempty"`
	Pagination Pagination  `json:"pagination,omitempty"`
}

type Pagination struct {
	TotalRecords int64 `json:"total_records"`
	TotalPages   int   `json:"total_pages"`
	CurrentPage  int   `json:"current_page"`
	PageSize     int   `json:"page_size"`
}

type IdVersionMapping struct {
	NIMB_ID     string `json:"nimbId" gorm:"not null" bson:"nimb_id"`
	NextVersion int    `json:"nextVersion" gorm:"not null" bson:"next_version"`
}
