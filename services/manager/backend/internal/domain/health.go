package domain

// HealthStatus is a simple domain object for service liveness.
type HealthStatus struct {
	Service string `json:"service"`
	Status  string `json:"status"`
	Time    string `json:"time"`
}
