package model

type SchedulerExcelReport struct {
	Id           string      `bson:"_id"`
	ReportName   string      `json:"reportName"`
	CreatedAt    interface{} `json:"createdAt"`
	UpdatedAt    interface{} `json:"updatedAt"`
	Status       string      `json:"status"`
	User         string      `json:"user"`
	FileName     string      `json:"fileName"`
	GeneratedAt  interface{} `json:"generatedAt"`
	ErrorMessage string      `json:"errorMessage"`
}
