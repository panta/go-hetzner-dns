package hetzner_dns

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"os"
	"time"

	"github.com/google/go-querystring/query"
	"github.com/pkg/errors"
)

const (
	BASE_URL        = "https://dns.hetzner.com/api/v1"
	DEFAULT_TIMEOUT = time.Second * 30
)

var (
	ErrAPIKeyNotSet = errors.New("hetzner_dns: API key has not been set")
	ErrMissingID    = errors.New("hetzner_dns: missing record ID")
)

// Client is the API service client structure.
type Client struct {
	BaseURL string
	ApiKey  string
	Debug   bool

	HttpClient *http.Client
}

func (client *Client) Perform(ctx context.Context, method string, endpoint string, queryParams, bodyParams, v interface{}) error {
	if client.HttpClient == nil {
		client.HttpClient = &http.Client{
			Timeout: DEFAULT_TIMEOUT,
		}
	}
	if client.BaseURL == "" {
		client.BaseURL = BASE_URL
	}

	// Create request
	endpointUrl := fmt.Sprintf("%s%s", client.BaseURL, endpoint)
	finalUrl := endpointUrl
	if queryParams != nil {
		v, err := query.Values(queryParams)
		if err != nil {
			return errors.Wrap(err, "can't process query params")
		}
		q := v.Encode()
		finalUrl += "?" + q
	}

	body := new(bytes.Buffer)
	if bodyParams != nil {
		if err := json.NewEncoder(body).Encode(bodyParams); err != nil {
			return errors.Wrap(err, "can't encode body params")
		}
	}
	// req, err := http.NewRequest(method, finalUrl, body)
	req, err := http.NewRequestWithContext(ctx, method, finalUrl, body)
	if err != nil {
		return errors.Wrap(err, "can't create http request")
	}

	// Headers
	apiKey := client.ApiKey
	if apiKey == "" {
		apiKey = os.Getenv("HETZNER_API_KEY")
	}
	if apiKey == "" {
		return ErrAPIKeyNotSet
	}
	req.Header.Add("Auth-API-Token", apiKey)

	if client.Debug {
		requestDump, _ := httputil.DumpRequestOut(req, true)
		log.Println(string(requestDump))
		log.Fatal()
	}

	// Perform Request
	resp, err := client.HttpClient.Do(req)
	if err != nil {
		return errors.Wrap(err, "can't perform http request")
	}
	defer resp.Body.Close()

	if client.Debug {
		responseDump, _ := httputil.DumpResponse(resp, true)
		log.Println(string(responseDump))
	}

	// Read Response Body
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return errors.Wrap(err, "can't read body")
	}

	if !isHttpSuccess(resp.StatusCode) {
		return errors.Wrap(err, "can't read body")
	}

	if v != nil {
		err = json.Unmarshal(respBody, v)
		if err != nil {
			return errors.Wrap(err, "can't parse response")
		}
	}
	return nil
}

func (client *Client) GetZones(ctx context.Context, name string, searchName string, page int, perPage int) (*ZonesResponse, error) {
	zonesResponse := ZonesResponse{}
	err := client.Perform(ctx, http.MethodGet, "/zones", struct {
		Name       string `url:"string"`
		Page       int    `url:"page"`
		PerPage    int    `url:"per_page"`
		SearchName string `url:"search_name"`
	}{
		Name:       name,
		SearchName: searchName,
		Page:       page,
		PerPage:    perPage,
	}, nil, &zonesResponse)
	return &zonesResponse, err
}

func (client *Client) GetRecords(ctx context.Context, zone_id string, page int, perPage int) (*RecordsResponse, error) {
	recordsResponse := RecordsResponse{}
	var params interface{}
	if (page > 0) && (perPage > 0) {
		params = struct {
			Page    int    `url:"page"`
			PerPage int    `url:"per_page"`
			ZoneId  string `url:"zone_id"`
		}{
			Page:    page,
			PerPage: perPage,
			ZoneId:  zone_id,
		}
	} else {
		params = struct {
			Page   int    `url:"page"`
			ZoneId string `url:"zone_id"`
		}{
			Page:   1,
			ZoneId: zone_id,
		}
	}
	err := client.Perform(ctx, http.MethodGet, "/records", params, nil, &recordsResponse)
	return &recordsResponse, err
}

func (client *Client) CreateRecord(ctx context.Context, record RecordRequest) (*RecordResponse, error) {
	recordResponse := RecordResponse{}
	err := client.Perform(ctx, http.MethodPost, "/records", nil, &record, &recordResponse)
	return &recordResponse, err
}

func (client *Client) GetRecord(ctx context.Context, recordId string) (*RecordResponse, error) {
	if recordId == "" {
		return nil, ErrMissingID
	}
	recordResponse := RecordResponse{}
	endpoint := fmt.Sprintf("/records/%v", recordId)
	err := client.Perform(ctx, http.MethodGet, endpoint, nil, nil, &recordResponse)
	return &recordResponse, err
}

func (client *Client) UpdateRecord(ctx context.Context, record RecordRequest) (*RecordResponse, error) {
	if record.ID == "" {
		return nil, ErrMissingID
	}
	recordResponse := RecordResponse{}
	endpoint := fmt.Sprintf("/records/%v", record.ID)
	err := client.Perform(ctx, http.MethodPut, endpoint, nil, nil, &recordResponse)
	return &recordResponse, err
}

func (client *Client) CreateOrUpdateRecord(ctx context.Context, record RecordRequest) (*RecordResponse, error) {
	if record.ID != "" {
		return client.UpdateRecord(ctx, record)
	}

	zoneId := record.ZoneID
	allRecords, err := client.GetRecords(ctx, zoneId, 0, 0)
	if err != nil {
		return nil, err
	}
	var foundRecord *Record
	for _, item := range allRecords.Records {
		if (record.ID != "") && (item.ID == record.ID) {
			foundRecord = &item
			break
		} else if (item.ZoneID == zoneId) && (item.Type == record.Type) && (item.Name == record.Name) {
			foundRecord = &item
			break
		}
	}
	if foundRecord != nil {
		return client.UpdateRecord(ctx, record)
	}
	return client.CreateRecord(ctx, record)
}

func (client *Client) DeleteRecord(ctx context.Context, recordId string) error {
	if recordId == "" {
		return ErrMissingID
	}
	endpoint := fmt.Sprintf("/records/%v", recordId)
	return client.Perform(ctx, http.MethodDelete, endpoint, nil, nil, nil)
}

func (client *Client) BulkCreateRecords(ctx context.Context, bulkRecordsRequest *BulkRecordRequest) (*BulkRecordResponse, error) {
	bulkRecordResponse := BulkRecordResponse{}
	err := client.Perform(ctx, http.MethodPost, "/records/bulk", nil, bulkRecordsRequest, &bulkRecordResponse)
	return &bulkRecordResponse, err
}

func (client *Client) BulkUpdateRecords(ctx context.Context, bulkRecordsRequest *BulkRecordRequest) (*BulkRecordResponse, error) {
	bulkRecordResponse := BulkRecordResponse{}
	err := client.Perform(ctx, http.MethodPut, "/records/bulk", nil, bulkRecordsRequest, &bulkRecordResponse)
	return &bulkRecordResponse, err
}

// // isHttpInformational returns true if HTTP status code is 1xx.
// func isHttpInformational(code int) bool {
// 	return (code >= 100) && (code <= 199)
// }

// isHttpSuccess returns true if HTTP status code is 2xx.
func isHttpSuccess(code int) bool {
	return (code >= 200) && (code <= 299)
}

// // isHttpRedirection returns true if HTTP status code is 3xx.
// func isHttpRedirection(code int) bool {
// 	return (code >= 300) && (code <= 399)
// }
//
// // isHttpClientError returns true if HTTP status code is 4xx.
// func isHttpClientError(code int) bool {
// 	return (code >= 400) && (code <= 499)
// }
//
// // isHttpServerError returns true if HTTP status code is 5xx.
// func isHttpServerError(code int) bool {
// 	return (code >= 500) && (code <= 599)
// }
