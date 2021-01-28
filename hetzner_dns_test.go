package hetzner_dns_test

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/pkg/errors"

	hetzner_dns "github.com/panta/go-hetzner-dns"
)

func TestClient_GetZones(t *testing.T) {
	handler := http.NotFound
	hs := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		handler(rw, req)
	}))
	defer hs.Close()
	c := hetzner_dns.Client{
		BaseURL: hs.URL,
	}

	handler = func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.Path != "/zones" {
			t.Error("Bad path!")
		}
		_, _ = io.WriteString(rw, `{
  "zones": [
    {
      "id": "sample-id",
      "created": "2021-01-28T14:23:31Z",
      "modified": "2021-01-28T14:23:31Z",
      "legacy_dns_host": "legacy-dns-host",
      "legacy_ns": [
        "legacy-ns-1"
      ],
      "name": "sample-name",
      "ns": [
        "sample-ns-1"
      ],
      "owner": "owner",
      "paused": true,
      "permission": "string",
      "project": "sample-project",
      "registrar": "sample-registrar",
      "status": "verified",
      "ttl": 0,
      "verified": "2021-01-28T14:23:31Z",
      "records_count": 0,
      "is_secondary_dns": true,
      "txt_verification": {
        "name": "string",
        "token": "string"
      }
    }
  ],
  "meta": {
    "pagination": {
      "page": 1,
      "per_page": 1,
      "last_page": 1,
      "total_entries": 0
    }
  }
}`)
	}

	_, err := c.GetZones(context.Background(), "name", "search-name", 1, 1)
	if err == nil {
		t.Error("Expected error to be non-nil")
	}
	if !errors.Is(err, hetzner_dns.ErrAPIKeyNotSet) {
		t.Error("Expected ErrAPIKeyNotSet")
	}

	c.ApiKey = "dummy"
	zonesResponse, err := c.GetZones(context.Background(), "name", "search-name", 1, 1)
	if err != nil {
		log.Println(err)
		t.Fatal("Got error performing request")
	}
	if zonesResponse == nil {
		t.Fatal("Did not get a response!")
	}
	if zonesResponse.Meta.Pagination.Page != 1 {
		t.Error("Wrong value for Page")
	}
	if zonesResponse.Meta.Pagination.PerPage != 1 {
		t.Error("Wrong value for PerPage")
	}
	if len(zonesResponse.Zones) != 1 {
		t.Error("Wrong number of zones")
	}
	if zonesResponse.Zones[0].ID != "sample-id" {
		t.Error("Wrong id for zone")
	}
	fmt.Println(zonesResponse)
}

func TestClient_GetRecords(t *testing.T) {
	handler := http.NotFound
	hs := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		handler(rw, req)
	}))
	defer hs.Close()
	c := hetzner_dns.Client{
		BaseURL: hs.URL,
	}

	handler = func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.Path != "/records" {
			t.Error("Bad path!")
		}
		_, _ = io.WriteString(rw, `{
  "records": [
    {
      "type": "A",
      "id": "string",
      "created": "2021-01-28T14:23:31Z",
      "modified": "2021-01-28T14:23:31Z",
      "zone_id": "sample-id",
      "name": "sample-name-1",
      "value": "sample-value-1",
      "ttl": 0
    },
    {
      "type": "A",
      "id": "string",
      "created": "2021-01-28T14:23:31Z",
      "modified": "2021-01-28T14:23:31Z",
      "zone_id": "sample-id",
      "name": "sample-name-2",
      "value": "sample-value-2",
      "ttl": 0
    },
    {
      "type": "A",
      "id": "string",
      "created": "2021-01-28T14:23:31Z",
      "modified": "2021-01-28T14:23:31Z",
      "zone_id": "sample-id",
      "name": "sample-name-3",
      "value": "sample-value-3",
      "ttl": 0
    }
  ]
}`)
	}

	_, err := c.GetRecords(context.Background(), "sample-id", 1, 10)
	if err == nil {
		t.Error("Expected error to be non-nil")
	}
	if !errors.Is(err, hetzner_dns.ErrAPIKeyNotSet) {
		t.Error("Expected ErrAPIKeyNotSet")
	}

	c.ApiKey = "dummy"
	recordsResponse, err := c.GetRecords(context.Background(), "sample-id", 1, 10)
	if err != nil {
		log.Println(err)
		t.Fatal("Got error performing request")
	}
	if recordsResponse == nil {
		t.Fatal("Did not get a response!")
	}
	if len(recordsResponse.Records) != 3 {
		t.Error("Wrong number of records")
	}
	if recordsResponse.Records[2].Value != "sample-value-3" {
		t.Error("Wrong value for record")
	}
	fmt.Println(recordsResponse)
}

func TestClient_CreateRecord(t *testing.T) {
	handler := http.NotFound
	hs := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		handler(rw, req)
	}))
	defer hs.Close()
	c := hetzner_dns.Client{
		BaseURL: hs.URL,
	}

	handler = func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.Path != "/records" {
			t.Error("Bad path!")
		}
		_, _ = io.WriteString(rw, `{
  "record": {
    "type": "A",
    "id": "sample-id",
    "created": "2021-01-28T14:23:31Z",
    "modified": "2021-01-28T14:23:31Z",
    "zone_id": "sample-zone",
    "name": "sample-name",
    "value": "sample-value",
    "ttl": 0
  }
}`)
	}

	_, err := c.CreateRecord(context.Background(), hetzner_dns.RecordRequest{
		ZoneID: "sample-zone",
		Type:   "A",
		Name:   "sample-name",
		Value:  "sample-value",
		TTL:    0,
	})
	if err == nil {
		t.Error("Expected error to be non-nil")
	}
	if !errors.Is(err, hetzner_dns.ErrAPIKeyNotSet) {
		t.Error("Expected ErrAPIKeyNotSet")
	}

	c.ApiKey = "dummy"
	recordResponse, err := c.CreateRecord(context.Background(), hetzner_dns.RecordRequest{
		ZoneID: "sample-zone",
		Type:   "A",
		Name:   "sample-name",
		Value:  "sample-value",
		TTL:    0,
	})
	if err != nil {
		log.Println(err)
		t.Fatal("Got error performing request")
	}
	if recordResponse == nil {
		t.Fatal("Did not get a response!")
	}
	if recordResponse.Record.Type != "A" {
		t.Error("Wrong record type")
	}
	fmt.Println(recordResponse)
}

func TestClient_GetRecord(t *testing.T) {
	handler := http.NotFound
	hs := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		handler(rw, req)
	}))
	defer hs.Close()
	c := hetzner_dns.Client{
		BaseURL: hs.URL,
	}

	handler = func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.Path != "/records/sample-id" {
			t.Error("Bad path!")
		}
		_, _ = io.WriteString(rw, `{
  "record": {
    "type": "A",
    "id": "sample-id",
    "created": "2021-01-28T14:23:31Z",
    "modified": "2021-01-28T14:23:31Z",
    "zone_id": "sample-zone",
    "name": "sample-name",
    "value": "sample-value",
    "ttl": 0
  }
}`)
	}

	_, err := c.GetRecord(context.Background(), "sample-id")
	if err == nil {
		t.Error("Expected error to be non-nil")
	}
	if !errors.Is(err, hetzner_dns.ErrAPIKeyNotSet) {
		t.Error("Expected ErrAPIKeyNotSet")
	}

	c.ApiKey = "dummy"
	recordResponse, err := c.GetRecord(context.Background(), "sample-id")
	if err != nil {
		log.Println(err)
		t.Fatal("Got error performing request")
	}
	if recordResponse == nil {
		t.Fatal("Did not get a response!")
	}
	if recordResponse.Record.Type != "A" {
		t.Error("Wrong record type")
	}
	fmt.Println(recordResponse)
}

func TestClient_UpdateRecord(t *testing.T) {
	handler := http.NotFound
	hs := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		handler(rw, req)
	}))
	defer hs.Close()
	c := hetzner_dns.Client{
		BaseURL: hs.URL,
	}

	handler = func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.Path != "/records/sample-id" {
			t.Error("Bad path!")
		}
		_, _ = io.WriteString(rw, `{
  "record": {
    "type": "A",
    "id": "sample-id",
    "created": "2021-01-28T14:23:31Z",
    "modified": "2021-01-28T14:23:31Z",
    "zone_id": "sample-zone",
    "name": "sample-name",
    "value": "sample-value",
    "ttl": 0
  }
}`)
	}

	_, err := c.UpdateRecord(context.Background(), hetzner_dns.RecordRequest{
		ID:     "sample-id",
		ZoneID: "sample-zone",
		Type:   "A",
		Name:   "sample-name",
		Value:  "sample-value",
		TTL:    0,
	})
	if err == nil {
		t.Error("Expected error to be non-nil")
	}
	if !errors.Is(err, hetzner_dns.ErrAPIKeyNotSet) {
		t.Error("Expected ErrAPIKeyNotSet")
	}

	_, err = c.UpdateRecord(context.Background(), hetzner_dns.RecordRequest{
		ZoneID: "sample-zone",
		Type:   "A",
		Name:   "sample-name",
		Value:  "sample-value",
		TTL:    0,
	})
	if err == nil {
		t.Error("Expected error to be non-nil")
	}
	if !errors.Is(err, hetzner_dns.ErrMissingID) {
		t.Error("Expected ErrMissingID")
	}

	c.ApiKey = "dummy"
	recordResponse, err := c.UpdateRecord(context.Background(), hetzner_dns.RecordRequest{
		ID:     "sample-id",
		ZoneID: "sample-zone",
		Type:   "A",
		Name:   "sample-name",
		Value:  "sample-value",
		TTL:    0,
	})
	if err != nil {
		log.Println(err)
		t.Fatal("Got error performing request")
	}
	if recordResponse == nil {
		t.Fatal("Did not get a response!")
	}
	if recordResponse.Record.Type != "A" {
		t.Error("Wrong record type")
	}
	fmt.Println(recordResponse)
}

func TestClient_DeleteRecord(t *testing.T) {
	handler := http.NotFound
	hs := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		handler(rw, req)
	}))
	defer hs.Close()
	c := hetzner_dns.Client{
		BaseURL: hs.URL,
	}

	handler = func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.Path != "/records/sample-id" {
			t.Error("Bad path!")
		}
		_, _ = io.WriteString(rw, `{
  "record": {
    "type": "A",
    "id": "sample-id",
    "created": "2021-01-28T14:23:31Z",
    "modified": "2021-01-28T14:23:31Z",
    "zone_id": "sample-zone",
    "name": "sample-name",
    "value": "sample-value",
    "ttl": 0
  }
}`)
	}

	err := c.DeleteRecord(context.Background(), "sample-id")
	if err == nil {
		t.Error("Expected error to be non-nil")
	}
	if !errors.Is(err, hetzner_dns.ErrAPIKeyNotSet) {
		t.Error("Expected ErrAPIKeyNotSet")
	}

	err = c.DeleteRecord(context.Background(), "")
	if err == nil {
		t.Error("Expected error to be non-nil")
	}
	if !errors.Is(err, hetzner_dns.ErrMissingID) {
		t.Error("Expected ErrMissingID")
	}

	c.ApiKey = "dummy"
	err = c.DeleteRecord(context.Background(), "sample-id")
	if err != nil {
		log.Println(err)
		t.Fatal("Got error performing request")
	}
}

func TestClient_BulkCreateRecords(t *testing.T) {
	handler := http.NotFound
	hs := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		handler(rw, req)
	}))
	defer hs.Close()
	c := hetzner_dns.Client{
		BaseURL: hs.URL,
	}

	handler = func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.Path != "/records/bulk" {
			t.Error("Bad path!")
		}
		_, _ = io.WriteString(rw, `{
  "records": [
    {
      "type": "A",
      "id": "sample-id",
      "created": "2021-01-28T14:23:31Z",
      "modified": "2021-01-28T14:23:31Z",
      "zone_id": "sample-zone",
      "name": "sample-name",
      "value": "sample-value",
      "ttl": 0
    }
  ],
  "valid_records": [
    {
      "zone_id": "sample-zone",
      "type": "A",
      "name": "sample-name",
      "value": "sample-value",
      "ttl": 0
    }
  ],
  "invalid_records": [
    {
      "zone_id": "sample-zone",
      "type": "A",
      "name": "sample-name-invalid",
      "value": "sample-value-invalid",
      "ttl": 0
    }
  ]
}`)
	}

	_, err := c.BulkCreateRecords(context.Background(), &hetzner_dns.BulkRecordRequest{Records: []hetzner_dns.RecordRequest{
		{
			ZoneID: "sample-zone",
			Type:   "A",
			Name:   "sample-name",
			Value:  "sample-value",
			TTL:    0,
		},
	}})
	if err == nil {
		t.Error("Expected error to be non-nil")
	}
	if !errors.Is(err, hetzner_dns.ErrAPIKeyNotSet) {
		t.Error("Expected ErrAPIKeyNotSet")
	}

	c.ApiKey = "dummy"
	bulkRecordsResponse, err := c.BulkCreateRecords(context.Background(), &hetzner_dns.BulkRecordRequest{Records: []hetzner_dns.RecordRequest{
		{
			ZoneID: "sample-zone",
			Type:   "A",
			Name:   "sample-name",
			Value:  "sample-value",
			TTL:    0,
		},
	}})
	if err != nil {
		log.Println(err)
		t.Fatal("Got error performing request")
	}
	if bulkRecordsResponse == nil {
		t.Fatal("Did not get a response!")
	}
	if len(bulkRecordsResponse.Records) != 1 {
		t.Error("Wrong # of records")
	}
	if len(bulkRecordsResponse.ValidRecords) != 1 {
		t.Error("Wrong # of valid records")
	}
	if len(bulkRecordsResponse.InvalidRecords) != 1 {
		t.Error("Wrong # of invalid records")
	}
	fmt.Println(bulkRecordsResponse)
}

func TestClient_BulkUpdateRecords(t *testing.T) {
	handler := http.NotFound
	hs := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		handler(rw, req)
	}))
	defer hs.Close()
	c := hetzner_dns.Client{
		BaseURL: hs.URL,
	}

	handler = func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.Path != "/records/bulk" {
			t.Error("Bad path!")
		}
		_, _ = io.WriteString(rw, `{
  "records": [
    {
      "type": "A",
      "id": "string",
      "created": "2021-01-28T14:23:31Z",
      "modified": "2021-01-28T14:23:31Z",
      "zone_id": "string",
      "name": "string",
      "value": "string",
      "ttl": 0
    }
  ],
  "failed_records": [
    {
      "zone_id": "string",
      "type": "A",
      "name": "string",
      "value": "string",
      "ttl": 0
    }
  ]
}`)
	}

	_, err := c.BulkUpdateRecords(context.Background(), &hetzner_dns.BulkRecordRequest{Records: []hetzner_dns.RecordRequest{
		{
			ID:     "sample-id",
			ZoneID: "sample-zone",
			Type:   "A",
			Name:   "sample-name",
			Value:  "sample-value",
			TTL:    0,
		},
	}})
	if err == nil {
		t.Error("Expected error to be non-nil")
	}
	if !errors.Is(err, hetzner_dns.ErrAPIKeyNotSet) {
		t.Error("Expected ErrAPIKeyNotSet")
	}

	c.ApiKey = "dummy"
	bulkRecordsResponse, err := c.BulkCreateRecords(context.Background(), &hetzner_dns.BulkRecordRequest{Records: []hetzner_dns.RecordRequest{
		{
			ID:     "sample-id",
			ZoneID: "sample-zone",
			Type:   "A",
			Name:   "sample-name",
			Value:  "sample-value",
			TTL:    0,
		},
	}})
	if err != nil {
		log.Println(err)
		t.Fatal("Got error performing request")
	}
	if bulkRecordsResponse == nil {
		t.Fatal("Did not get a response!")
	}
	if len(bulkRecordsResponse.Records) != 1 {
		t.Error("Wrong # of records")
	}
	if len(bulkRecordsResponse.ValidRecords) != 0 {
		t.Error("Found unexpected valid records (should not be present)")
	}
	if len(bulkRecordsResponse.InvalidRecords) != 0 {
		t.Error("Found unexpected invalid records (should not be present)")
	}
	if len(bulkRecordsResponse.FailedRecords) != 1 {
		t.Error("Wrong # of failed records")
	}
	fmt.Println(bulkRecordsResponse)
}
