package models

import "time"

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

type LogStackEntry struct {
	Timestamp time.Time `json:"timestamp" gorm:"autoCreateTime" bson:"timestamp"`
	Type      string    `json:"type" gorm:"type:varchar(50)" bson:"type"` // e.g., INFO, ERROR, DEBUG
	Message   string    `json:"message" gorm:"type:text" bson:"message"`
}
