package domain

type Report struct {
	Start         string         `json:"dateRangeStart"`
	End           string         `json:"dateRangeEnd"`
	SummaryFilter *Report_filter `json:"summaryFilter,omitempty"`
	SortOrder     string         `json:"sortOrder"`
	Users         *Report_user   `json:"users,omitempty"`
}
type Report_filter struct {
	Groups []string `json:"groups"`
}
type Report_user struct {
	Ids      []string `json:"ids"`
	Contains string   `json:"contains"`
	Status   string   `json:"status"`
}

type Result struct {
	Entries []ResultUser `json:"groupOne"`
}

type ResultUser struct {
	Name     string
	Duration int64
	Entries  []ReportEntry `json:"children"`
}

type ReportEntry struct {
	Duration int
	Date     string `json:"name"`
}

type PartTime struct {
	Start    string
	End      string
	Capacity int64 // 0...100 (percent)
}
