package alloha

import (
	"context"
	"io"
	"log"
	"net/http"
	"net/url"
	"time"
)

const (
	defaultTimeout = 10 * time.Second
)

// Client is the structure of the client API
type Client struct {
	apiToken string
	baseURL  string
	client   *http.Client
	useProxy bool
}

//region - Constructor

// NewClient creates a new client instance
func NewClient(apiToken, baseApiURL string) (*Client, error) {
	buildURL, err := buildApiURL(baseApiURL, apiToken)
	if err != nil {
		return nil, err
	}

	client := &Client{
		apiToken: apiToken,
		baseURL:  buildURL,
		client: &http.Client{
			Timeout: defaultTimeout,
		},
		useProxy: false,
	}

	return client, nil
}

//endregion

//region - Public Methods

// SetApiToken sets a new API token
func (c *Client) SetApiToken(apiToken string) error {
	buildURL, err := buildApiURL(c.baseURL, apiToken)
	if err != nil {
		return err
	}

	c.apiToken = apiToken
	c.baseURL = buildURL

	return nil
}

// SetBaseApiUrl sets a new base API URL
func (c *Client) SetBaseApiUrl(baseApiURL string) error {
	buildURL, err := buildApiURL(baseApiURL, c.apiToken)
	if err != nil {
		return err
	}

	c.baseURL = buildURL

	return nil
}

// SetHttpTimeout sets a new HTTP timeout
func (c *Client) SetHttpTimeout(timeout time.Duration) {
	c.client.Timeout = timeout
}

//endregion

//region - Private Methods

// buildApiURL builds the API URL
func buildApiURL(apiToken, baseApiUrl string) (string, error) {
	if apiToken == "" {
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

// doGetApiRequest sends a GET request to the API.
func (c *Client) doGetApiRequest(ctx context.Context, endpointApi string) ([]byte, int, error) {
	var bodyBytes []byte
	var err error
	var PTransport *http.Transport
	var req *http.Request
	var resp *http.Response
	var statusCode = 0

	if c.useProxy {
		PTransport = &http.Transport{
			Proxy: http.ProxyFromEnvironment,
			/* Proxy: http.ProxyURL(&url.URL{
			     Scheme: "http",
			     User:   url.UserPassword("login", "password"),
			     Host:   "IP:PORT",
			   }),
			*/
		}
	}

	client := &http.Client{
		Transport: PTransport,
	}

	if endpointApi == "" {
		return nil, statusCode, EmptyEndpointApiURLError
	}

	req, err = http.NewRequestWithContext(ctx, http.MethodGet, endpointApi, nil)
	if err != nil {
		return nil, statusCode, err
	}

	req.Header.Add("Accept", "application/json, text/plain, */*")
	req.Header.Add("Accept-Encoding", "gzip, deflate, br")
	req.Header.Add("Accept-Language", "ru-RU,ru;q=0.9,en-US;q=0.8,en;q=0.7")
	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/89.0.4389.90 Safari/537.36")

	resp, err = client.Do(req)
	if err != nil {
		return nil, statusCode, err
	}

	defer func() {
		if respErr := resp.Body.Close(); respErr != nil {
			log.Printf("failed to close response body: %v", respErr)
		}
	}()

	statusCode = resp.StatusCode

	bodyBytes, err = io.ReadAll(resp.Body)
	if err != nil {
		return nil, statusCode, err
	}

	return bodyBytes, statusCode, nil
}

//endregion
