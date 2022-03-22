package model

type DnsCheckResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Data    struct {
		Query struct {
			Type string `json:"type"`
			Name string `json:"name"`
		} `json:"query"`
		Answer []string `json:"answer"`
	} `json:"data"`
}
