package alloha

import (
	"bytes"
	"compress/flate"
	"compress/gzip"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

// HttpClient provides an interface for executing HTTP requests
type HttpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// AllohaClient is the structure of the client API
type AllohaClient struct {
	apiToken string
	baseURL  string
	client   HttpClient
}

//region - Constructor

// NewAllohaClient creates a new AllohaClient instance
func NewAllohaClient(httpClient HttpClient, apiToken, baseApiURL string) (*AllohaClient, error) {
	buildURL, err := buildApiURL(apiToken, baseApiURL)
	if err != nil {
		return nil, err
	}

	client := &AllohaClient{
		apiToken: apiToken,
		baseURL:  buildURL,
		client:   httpClient,
	}

	return client, nil
}

//endregion

//region - Public Methods

// FindByKPId finds a movie by its KP ID
func (c *AllohaClient) FindByKPId(ctx context.Context, kpId int) (*FindOneResponse, error) {
	if kpId <= 0 {
		return nil, InvalidKPIdParameterError
	}

	var err error
	var bodyBytes []byte
	var parsedBaseURL *url.URL
	var statusCode = 0
	var response *FindOneResponse

	parsedBaseURL, err = url.Parse(c.baseURL)
	if err != nil {
		return nil, err
	}

	queryValues := parsedBaseURL.Query()
	queryValues.Set("kp", strconv.Itoa(kpId))
	parsedBaseURL.RawQuery = queryValues.Encode()

	bodyBytes, statusCode, err = c.doApiRequest(ctx, http.MethodGet, parsedBaseURL.String(), nil)
	if err != nil {
		return nil, err
	}

	if statusCode != 200 {
		return nil, fmt.Errorf("unexpected server response with a status code: %d", statusCode)
	}

	err = json.Unmarshal(bodyBytes, &response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

// SetApiToken sets a new API token
func (c *AllohaClient) SetApiToken(apiToken string) error {
	buildURL, err := buildApiURL(apiToken, c.baseURL)
	if err != nil {
		return err
	}

	c.apiToken = apiToken
	c.baseURL = buildURL

	return nil
}

// SetBaseApiUrl sets a new base API URL
func (c *AllohaClient) SetBaseApiUrl(baseApiURL string) error {
	buildURL, err := buildApiURL(c.apiToken, baseApiURL)
	if err != nil {
		return err
	}

	c.baseURL = buildURL

	return nil
}

//endregion

//region - Private Methods

// buildApiURL builds the API URL
func buildApiURL(apiToken, baseApiUrl string) (string, error) {
	if apiToken == "" { //  TODO len(apiToken) <= 0
		return "", ApiTokenEmptyError
	}

	if baseApiUrl == "" {
		return "", BaseApiUrlEmptyError
	}

	parsedURL, err := url.Parse(baseApiUrl)
	if err != nil {
		return "", err
	}

	if parsedURL.Host == "" {
		return "", BaseApiUrlInvalidHostError
	}

	parsedURL.Path = "/"

	queryValues := parsedURL.Query()
	queryValues.Set("token", apiToken)

	parsedURL.RawQuery = queryValues.Encode()

	return parsedURL.String(), nil
}

// doApiRequest executes the specified HTTP request to the specified URL with the specified request body and returns
// the response body, the response code, and the error, if any.
func (c *AllohaClient) doApiRequest(ctx context.Context, method, endpointApiUrl string, requestBody []byte) ([]byte, int, error) {
	var bodyBytes []byte
	var err error
	var req *http.Request
	var resp *http.Response
	var respReader io.ReadCloser
	var statusCode = 0

	if len(endpointApiUrl) <= 0 {
		return nil, statusCode, EmptyEndpointApiURLError
	}

	if requestBody == nil || len(requestBody) <= 0 {
		req, err = http.NewRequestWithContext(ctx, method, endpointApiUrl, nil)
	} else {
		req, err = http.NewRequestWithContext(ctx, method, endpointApiUrl, bytes.NewBuffer(requestBody))
		if req != nil {
			req.Header.Set("Content-Type", "application/json; charset=utf-8")
		}
	}
	if err != nil {
		return nil, statusCode, err
	}
	if req == nil {
		return nil, statusCode, FailedCreateRequestError
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Accept-Encoding", "gzip, deflate")
	req.Header.Add("Accept-Language", "ru-RU,ru;q=0.9,en-US;q=0.8,en;q=0.7")
	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/89.0.4389.90 Safari/537.36")

	resp, err = c.client.Do(req)
	if err != nil {
		return nil, statusCode, err
	}

	statusCode = resp.StatusCode

	// Check if the response is encoded in deflate or gzip format
	respContentEncoding := strings.ToLower(resp.Header.Get("Content-Encoding"))
	switch respContentEncoding {
	case "deflate":
		respReader = flate.NewReader(resp.Body)
	case "gzip":
		respReader, err = gzip.NewReader(resp.Body)
		if err != nil {
			return nil, statusCode, err
		}
	default:
		respReader = resp.Body
	}

	defer func() {
		if closeErr := respReader.Close(); closeErr != nil {
			log.Printf("failed to close response reader: %s", closeErr.Error())
		}
	}()

	bodyBytes, err = io.ReadAll(resp.Body)
	if err != nil {
		return nil, statusCode, err
	}

	return bodyBytes, statusCode, nil
}

//endregion
