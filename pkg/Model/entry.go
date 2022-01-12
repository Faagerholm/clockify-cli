package model

type Entry struct {
	ID           string `json:"id"`
	ProjectID    string `json:"projectId"`
	Description  string `json:"description"`
	TaskID       string `json:"taskId"`
	Billable     bool   `json:"billable"`
	TimeInterval struct {
		Start string `json:"start"`
		End   string `json:"end"`
	}
	TagIDs []string `json:"tagIds"`
}

type UpdateEntry struct {
	Start       string   `json:"start"`
	Billable    bool     `json:"billable"`
	Description string   `json:"description"`
	ProjectID   string   `json:"projectId"`
	TaskID      string   `json:"taskId"`
	End         string   `json:"end"`
	TagIDs      []string `json:"tagIds"`
}
