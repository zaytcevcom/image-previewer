package fetcher

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFetcher(t *testing.T) {

	handler := func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte("Hello, World!"))
		if err != nil {
			return
		}
	}

	tests := []struct {
		name    string
		handler http.HandlerFunc
		rawUrl  string
		wantErr bool
	}{
		{
			name:    "success",
			handler: handler,
			rawUrl:  "/",
			wantErr: false,
		},
		{
			name:    "invalid URL",
			handler: handler,
			rawUrl:  ":foo",
			wantErr: true,
		},
		{
			name:    "success with scheme",
			handler: handler,
			rawUrl:  "https://example.com",
			wantErr: true,
		},
		{
			name:    "success without scheme",
			handler: handler,
			rawUrl:  "example.com",
			wantErr: true,
		},
		{
			name:    "invalid scheme",
			handler: handler,
			rawUrl:  "ftp://example.com",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ts := httptest.NewServer(tt.handler)
			defer ts.Close()

			fetcher := New()
			_, err := fetcher.Fetch(context.Background(), nil, ts.URL+tt.rawUrl)
			hasErr := err != nil

			assert.Equal(t, tt.wantErr, hasErr)
		})
	}
}

func TestNew(t *testing.T) {
	fetcher := New()
	assert.NotNil(t, fetcher)
}
