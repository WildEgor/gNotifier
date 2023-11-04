package dtos

// StatusDto holds the information for the server's status
type StatusDto struct {
	Status      string `json:"status"`
	Version     string `json:"version"`
	Environment string `json:"environment"`
}
