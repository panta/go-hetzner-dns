package hetzner_dns

import "time"

type Pagination struct {
	Page         int `json:"page"`
	PerPage      int `json:"per_page"`
	LastPage     int `json:"last_page"`
	TotalEntries int `json:"total_entries"`
}

type Meta struct {
	Pagination Pagination `json:"pagination"`
}

type Zone struct {
	ID              string    `json:"id"`
	Created         time.Time `json:"created"`
	Modified        time.Time `json:"modified"`
	LegacyDNSHost   string    `json:"legacy_dns_host"`
	LegacyNs        []string  `json:"legacy_ns"`
	Name            string    `json:"name"`
	Ns              []string  `json:"ns"`
	Owner           string    `json:"owner"`
	Paused          bool      `json:"paused"`
	Permission      string    `json:"permission"`
	Project         string    `json:"project"`
	Registrar       string    `json:"registrar"`
	Status          string    `json:"status"`
	TTL             int       `json:"ttl"`
	Verified        time.Time `json:"verified"`
	RecordsCount    int       `json:"records_count"`
	IsSecondaryDNS  bool      `json:"is_secondary_dns"`
	TxtVerification struct {
		Name  string `json:"name"`
		Token string `json:"token"`
	} `json:"txt_verification"`
}

type ZonesResponse struct {
	Zones []Zone `json:"zones"`
	Meta  Meta   `json:"meta"`
}

type ZoneRequest struct {
	Name string `json:"name"`
	TTL  int    `json:"ttl"`
}
