package hetzner_dns

type RecordRequest struct {
	ID     string `json:"id,omitempty"`
	ZoneID string `json:"zone_id"`
	Type   string `json:"type"`
	Name   string `json:"name"`
	Value  string `json:"value"`
	TTL    int    `json:"ttl"`
}

type BulkRecordRequest struct {
	Records []RecordRequest `json:"records"`
}

type Record struct {
	Type     string      `json:"type"`
	ID       string      `json:"id"`
	Created  HetznerTime `json:"created"`
	Modified HetznerTime `json:"modified"`
	ZoneID   string      `json:"zone_id"`
	Name     string      `json:"name"`
	Value    string      `json:"value"`
	TTL      int         `json:"ttl"`
}

type RecordResponse struct {
	Record Record `json:"record"`
}

type RecordsResponse struct {
	Records []Record `json:"records"`
}

type BulkRecordResponse struct {
	Records        []Record        `json:"records"`
	ValidRecords   []RecordRequest `json:"valid_records,omitempty"`
	InvalidRecords []RecordRequest `json:"invalid_records,omitempty"`
	FailedRecords  []RecordRequest `json:"failed_records,omitempty"`
}
