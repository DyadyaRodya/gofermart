package accrualsdk

import (
	"context"
	"crypto/tls"
	"fmt"
	"golang.org/x/time/rate"
	"io"
	"net/http"
	"net/url"
	"path"
	"strconv"
	"strings"
)

const apiPrefix = "api"

var ErrAccrualClientRetryLater = fmt.Errorf("retry later")

type rateLimitTransport struct {
	limiter *rate.Limiter
	xport   http.RoundTripper
}

var _ http.RoundTripper = &rateLimitTransport{}

func newRateLimitTransport(r float64, xport http.RoundTripper) http.RoundTripper {
	return &rateLimitTransport{
		limiter: rate.NewLimiter(rate.Limit(r), 1),
		xport:   xport,
	}
}

func (t *rateLimitTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	t.limiter.Wait(r.Context())
	return t.xport.RoundTrip(r)
}

type AccrualClient struct {
	BaseURL          *url.URL
	defaultTransport *http.Transport
	httpClient       *http.Client
	skipTLS          bool
}

func NewAccrualClient(baseURL *url.URL, skipTLS bool) *AccrualClient {
	skipTLSTransport := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: skipTLS, //nolint:gosec // allow to override skip tls
		},
	}
	return &AccrualClient{
		BaseURL:          baseURL,
		defaultTransport: skipTLSTransport,
		httpClient: &http.Client{
			Transport: skipTLSTransport,
		},
		skipTLS: skipTLS,
	}
}

func (c *AccrualClient) limitFromBody(content []byte) (float64, error) {
	// Convert byte slice to string
	str := string(content)

	// Split the string by spaces
	parts := strings.Fields(str)

	// Iterate through parts to find the numeric value
	for _, part := range parts {
		if num, err := strconv.ParseFloat(part, 64); err == nil {
			return num, nil
		}
	}

	return 0, fmt.Errorf("no numbers found in the string: %s", str)
}

func (c *AccrualClient) updateRateLimitPerMinute(limitPerMinute float64) {
	limit := limitPerMinute * 60 // rps
	c.httpClient.Transport = newRateLimitTransport(limit, c.defaultTransport)
}

func (c *AccrualClient) makeLimitedGetRequest(
	ctx context.Context, requestPath string, queries map[string]string,
) ([]byte, int, error) {
	requestURL := *c.BaseURL
	requestURL.Path = path.Join(requestURL.Path, apiPrefix, requestPath)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, requestURL.String(), http.NoBody)
	if err != nil {
		return nil, 0, fmt.Errorf("http.NewRequestWithContext: %w", err)
	}

	q := req.URL.Query()
	for k, v := range queries {
		q.Add(k, v)
	}
	req.URL.RawQuery = q.Encode()

	httpResp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, 0, fmt.Errorf("c.httpClient.Do: %w", err)
	}
	defer func() {
		_ = httpResp.Body.Close()
	}()

	if httpResp.StatusCode == http.StatusNoContent {
		return nil, httpResp.StatusCode, nil
	}

	content, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return nil, httpResp.StatusCode, fmt.Errorf("io.ReadAll: %w", err)
	}

	if httpResp.StatusCode == http.StatusTooManyRequests {
		limitPerMinute, err := c.limitFromBody(content)
		if err != nil {
			return nil, httpResp.StatusCode, fmt.Errorf("AccrualClient.limitFromBody: %w", err)
		}
		c.updateRateLimitPerMinute(limitPerMinute)
		return content, httpResp.StatusCode, ErrAccrualClientRetryLater
	}

	return content, httpResp.StatusCode, nil
}
