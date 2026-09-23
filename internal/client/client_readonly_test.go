//go:build readonly

package client

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestReadOnlyMode_Enabled(t *testing.T) {
	if !ReadOnlyMode {
		t.Skip("skipping readonly test: binary compiled without -tags readonly")
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	cfg := testConfig(server.URL)
	c, err := New(cfg)
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name   string
		fn     func() error
		method string
	}{
		{"POST blocked", func() error { return c.Post(context.Background(), "/test", nil, nil) }, "POST"},
		{"PUT blocked", func() error { return c.Put(context.Background(), "/test", nil, nil) }, "PUT"},
		{"DELETE blocked", func() error { return c.Delete(context.Background(), "/test") }, "DELETE"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.fn()
			if err == nil {
				t.Fatalf("%s should be blocked in readonly mode", tt.method)
			}
			if !strings.Contains(err.Error(), "read-only mode") {
				t.Errorf("error should mention read-only mode, got: %v", err)
			}
		})
	}

	if err := c.Get(context.Background(), "/test", nil); err != nil {
		t.Errorf("GET should still work in readonly mode, got: %v", err)
	}
}
