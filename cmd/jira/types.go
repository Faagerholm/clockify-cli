package jira

type requestBody struct {
	From       string   `json:"from"`
	To         string   `json:"to"`
	ProjectKey []string `json:"projectKey"`
}

type worklogEntry struct {
	BillableSeconds int `json:"billableSeconds"`
	// Comment         string `json:"comment"`
	Date  string `json:"started"`
	Issue struct {
		ProjectKey string `json:"projectKey"`
	} `json:"issue"`
}
