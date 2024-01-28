package fetcher

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

var ErrNotFound = errors.New("404 Not Found")

type Fetcher struct{}

func New() Fetcher {
	return Fetcher{}
}

func (f Fetcher) Fetch(ctx context.Context, header http.Header, rawUrl string) ([]byte, error) {

	parsedUrl, err := url.Parse(rawUrl)
	if err != nil {
		return nil, fmt.Errorf("failed parse URL: %w", err)
	}

	if parsedUrl.Scheme == "" {
		parsedUrl.Scheme = "http"
	}

	if parsedUrl.Scheme != "http" {
		return nil, fmt.Errorf("failed scheme URL: %w", err)
	}

	request, err := http.NewRequestWithContext(ctx, http.MethodGet, parsedUrl.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed prepare request: %w", err)
	}

	request.Header = header

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return nil, fmt.Errorf("failed request: %w", err)
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(response.Body)

	if response.StatusCode == http.StatusNotFound {
		return nil, ErrNotFound
	}

	result, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("failed reading response body: %w", err)
	}

	return result, nil
}
