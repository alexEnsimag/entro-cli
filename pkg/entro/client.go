package entro

import (
	"alex/entro-cli/pkg/report"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

const urlBase = "http://localhost:8090"

// Client is an entro client
type Client struct {
	accessKeyID     string
	secretAccessKey string
	sessionToken    string
	region          string

	httpClient http.Client
	storage    report.Storage
}

func NewClient(accessKeyID, secretAccessKey, sessionToken, region string) *Client {
	return &Client{
		accessKeyID:     accessKeyID,
		secretAccessKey: secretAccessKey,
		sessionToken:    sessionToken,
		region:          region,
		httpClient: http.Client{
			Timeout: 10 * time.Second,
		},
		storage: report.Storage{},
	}
}

// CreateReport requests the generation of a report, and returns its ID
func (c *Client) CreateReport() (report.ID, error) {
	reqBody := struct {
		AWSAccessKeyID     string `json:"awsAccessKeyID"`
		AWSSecretAccessKey string `json:"awsSecretAccessKey"`
		AWSSessionToken    string `json:"awsSessionToken"`
		AWSRegion          string `json:"awsRegion"`
	}{
		AWSAccessKeyID:     c.accessKeyID,
		AWSSecretAccessKey: c.secretAccessKey,
		AWSSessionToken:    c.sessionToken,
		AWSRegion:          c.region,
	}
	var respBody struct {
		ReportID report.ID `json:"reportID"`
	}
	if err := c.request(http.MethodPost, urlBase+"/create", reqBody, &respBody); err != nil {
		return "", fmt.Errorf("failed to request report: %w", err)
	}
	return respBody.ReportID, nil
}

func (c *Client) GetReportStatus(reportID report.ID) (report.Status, error) {
	urlValues := url.Values{}
	urlValues.Add("reportID", string(reportID))
	u := fmt.Sprintf("%s/status?%s", urlBase, urlValues.Encode())

	var respBody struct {
		ReportStatus report.Status `json:"status"`
	}
	if err := c.request(http.MethodGet, u, nil, &respBody); err != nil {
		return "", fmt.Errorf("failed to query report status: %w", err)
	}
	return respBody.ReportStatus, nil
}

func (c *Client) GetReport(reportID report.ID) ([]byte, error) {
	urlValues := url.Values{}
	urlValues.Add("reportID", string(reportID))
	u := fmt.Sprintf("%s/filePath?%s", urlBase, urlValues.Encode())

	var respBody struct {
		Path string `json:"path"`
	}
	if err := c.request(http.MethodGet, u, nil, &respBody); err != nil {
		return nil, fmt.Errorf("failed to query report path: %w", err)
	}

	data, err := c.storage.ReadRaw(respBody.Path)
	if err != nil {
		return nil, fmt.Errorf("failed to get report: %w", err)
	}
	return data, nil
}

func (c *Client) request(method string, url string, reqBody, respBody any) error {
	//create body if exists
	var request *http.Request
	if reqBody != nil {
		body, err := json.Marshal(reqBody)
		if err != nil {
			panic("failed to marshal create body: " + err.Error())
		}

		req, err := http.NewRequest(method, url, bytes.NewBuffer(body))
		if err != nil {
			panic("failed to create http request: " + err.Error())
		}
		request = req
	} else {
		req, err := http.NewRequest(method, url, nil)
		if err != nil {
			panic("failed to create http request: " + err.Error())
		}
		request = req
	}

	// query entro
	resp, err := c.httpClient.Do(request)
	if err != nil {
		return fmt.Errorf("failed to query entro: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		// parse body if exists
		if respBody != nil {
			respBodyBytes, err := io.ReadAll(resp.Body)
			if err != nil {
				return fmt.Errorf("failed to read body: %w", err)
			}
			return json.Unmarshal(respBodyBytes, respBody)
		}
		return nil
	}
	return fmt.Errorf("request failed with status %d", resp.StatusCode)
}
