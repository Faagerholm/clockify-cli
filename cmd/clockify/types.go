package clockify

/* const */
const (
	week_seconds = 27000 // 37.5 * 60 * 60
)

/* Request type */
type Report struct {
	Start         string         `json:"dateRangeStart"`
	End           string         `json:"dateRangeEnd"`
	SummaryFilter *report_filter `json:"summaryFilter,omitempty"`
	SortOrder     string         `json:"sortOrder"`
	Users         *report_user   `json:"users,omitempty"`
}

type project struct {
	Name string `json:"name"`
	Id   string `json:"id"`
}

type report_filter struct {
	Groups []string `json:"groups"`
}

type report_user struct {
	Ids      []string `json:"ids"`
	Contains string   `json:"contains"`
	Status   string   `json:"status"`
}

type LogEntry struct {
	Start       string `json:"start"`
	End         string `json:"end,omitempty"`
	ProjectId   string `json:"projectId"`
	Description string `json:"description,omitempty"`
}

/* Response types */
type ClockifyReport struct {
	Entries []struct {
		User     string `json:"name"`
		Total    int    `json:"duration"`
		Children []struct {
			Date     string `json:"name"`
			Duration int
			Children []struct {
				Project  string `json:"name"`
				Duration int
				Amount   int
			} `json:"children"`
		} `json:"children"`
	} `json:"groupOne"`
}

type partTime []struct {
	startDate string
	endDate   string
	capacity  float64
}
