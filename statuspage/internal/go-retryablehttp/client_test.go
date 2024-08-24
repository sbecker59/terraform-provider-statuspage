package retryablehttp

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"
	"time"
)

func TestClient_DefaultBackoff(t *testing.T) {

	tests := []struct {
		name        string
		code        int
		retryHeader string
	}{
		{"http_429_seconds", http.StatusTooManyRequests, "2"},
		{"http_429_date", http.StatusTooManyRequests, "Fri, 31 Dec 1999 23:59:59 GMT"},
		{"http_420_seconds", 420, "2"},
		{"http_420_date", 420, "Fri, 31 Dec 1999 23:59:59 GMT"},
		{"http_503_seconds", http.StatusServiceUnavailable, "2"},
		{"http_503_date", http.StatusServiceUnavailable, "Fri, 31 Dec 1999 23:59:59 GMT"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Retry-After", test.retryHeader)
				http.Error(w, fmt.Sprintf("test_%d_body", test.code), test.code)
			}))
			defer ts.Close()

			client := NewClient()

			var retryAfter time.Duration
			retryable := false

			client.CheckRetry = func(_ context.Context, resp *http.Response, err error) (bool, error) {
				retryable, _ = DefaultRetryPolicy(context.Background(), resp, err)
				retryAfter = DefaultBackoff(client.RetryWaitMin, client.RetryWaitMax, 1, resp)
				return false, nil
			}

			_, err := client.Get(ts.URL)
			if err != nil {
				t.Fatalf("expected no errors since retryable")
			}

			if !retryable {
				t.Fatal("Since the error is recoverable, the default policy shall return true")
			}

			if retryAfter != 2*time.Second {
				t.Fatalf("The header Retry-After specified 2 seconds, and shall not be %d seconds", retryAfter/time.Second)
			}
		})
	}
}

func TestClient_DefaultRetryPolicy_TLS(t *testing.T) {
	ts := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
	}))
	defer ts.Close()

	attempts := 0
	client := NewClient()
	client.CheckRetry = func(_ context.Context, resp *http.Response, err error) (bool, error) {
		attempts++
		return DefaultRetryPolicy(context.TODO(), resp, err)
	}

	_, err := client.Get(ts.URL)
	if err == nil {
		t.Fatalf("expected x509 error, got nil")
	}
	if attempts != 1 {
		t.Fatalf("expected 1 attempt, got %d", attempts)
	}
}

func TestClient_DefaultRetryPolicy_redirects(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/", http.StatusFound)
	}))
	defer ts.Close()

	attempts := 0
	client := NewClient()
	client.CheckRetry = func(_ context.Context, resp *http.Response, err error) (bool, error) {
		attempts++
		return DefaultRetryPolicy(context.TODO(), resp, err)
	}

	_, err := client.Get(ts.URL)
	if err == nil {
		t.Fatalf("expected redirect error, got nil")
	}
	if attempts != 1 {
		t.Fatalf("expected 1 attempt, got %d", attempts)
	}
}

func TestClient_DefaultRetryPolicy_invalidscheme(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
	}))
	defer ts.Close()

	attempts := 0
	client := NewClient()
	client.CheckRetry = func(_ context.Context, resp *http.Response, err error) (bool, error) {
		attempts++
		return DefaultRetryPolicy(context.TODO(), resp, err)
	}

	url := strings.Replace(ts.URL, "http", "ftp", 1)
	_, err := client.Get(url)
	if err == nil {
		t.Fatalf("expected scheme error, got nil")
	}
	if attempts != 1 {
		t.Fatalf("expected 1 attempt, got %d", attempts)
	}
}

func TestClient_DefaultRetryPolicy_invalidheadername(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
	}))
	defer ts.Close()

	attempts := 0
	client := NewClient()
	client.CheckRetry = func(_ context.Context, resp *http.Response, err error) (bool, error) {
		attempts++
		return DefaultRetryPolicy(context.TODO(), resp, err)
	}

	req, err := http.NewRequest(http.MethodGet, ts.URL, nil)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	req.Header.Set("Header-Name-\033", "header value")
	_, err = client.StandardClient().Do(req)
	if err == nil {
		t.Fatalf("expected header error, got nil")
	}
	if attempts != 1 {
		t.Fatalf("expected 1 attempt, got %d", attempts)
	}
}

func TestClient_DefaultRetryPolicy_invalidheadervalue(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
	}))
	defer ts.Close()

	attempts := 0
	client := NewClient()
	client.CheckRetry = func(_ context.Context, resp *http.Response, err error) (bool, error) {
		attempts++
		return DefaultRetryPolicy(context.TODO(), resp, err)
	}

	req, err := http.NewRequest(http.MethodGet, ts.URL, nil)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	req.Header.Set("Header-Name", "bad header value \033")
	_, err = client.StandardClient().Do(req)
	if err == nil {
		t.Fatalf("expected header value error, got nil")
	}
	if attempts != 1 {
		t.Fatalf("expected 1 attempt, got %d", attempts)
	}
}

func TestClient_CheckRetryStop(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "test_500_body", http.StatusInternalServerError)
	}))
	defer ts.Close()

	client := NewClient()

	// Verify that this stops retries on the first try, with no errors from the client.
	called := 0
	client.CheckRetry = func(_ context.Context, resp *http.Response, err error) (bool, error) {
		called++
		return false, nil
	}

	_, err := client.Get(ts.URL)

	if called != 1 {
		t.Fatalf("CheckRetry called %d times, expected 1", called)
	}

	if err != nil {
		t.Fatalf("Expected no error, got:%v", err)
	}
}

func TestBackoff(t *testing.T) {
	type tcase struct {
		min    time.Duration
		max    time.Duration
		i      int
		expect time.Duration
	}
	cases := []tcase{
		{
			time.Second,
			5 * time.Minute,
			0,
			time.Second,
		},
		{
			time.Second,
			5 * time.Minute,
			1,
			2 * time.Second,
		},
		{
			time.Second,
			5 * time.Minute,
			2,
			4 * time.Second,
		},
		{
			time.Second,
			5 * time.Minute,
			3,
			8 * time.Second,
		},
		{
			time.Second,
			5 * time.Minute,
			63,
			5 * time.Minute,
		},
		{
			time.Second,
			5 * time.Minute,
			128,
			5 * time.Minute,
		},
	}

	for _, tc := range cases {
		if v := DefaultBackoff(tc.min, tc.max, tc.i, nil); v != tc.expect {
			t.Fatalf("bad: %#v -> %s", tc, v)
		}
	}
}

func TestClient_BackoffCustom(t *testing.T) {
	var retries int32

	client := NewClient()
	client.Backoff = func(min, max time.Duration, attemptNum int, resp *http.Response) time.Duration {
		atomic.AddInt32(&retries, 1)
		return time.Millisecond * 1
	}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if atomic.LoadInt32(&retries) == int32(client.RetryMax) {
			w.WriteHeader(200)
			return
		}
		w.WriteHeader(500)
	}))
	defer ts.Close()

	// Make the request.
	resp, err := client.Get(ts.URL + "/foo/bar")
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	resp.Body.Close()
	if retries != int32(client.RetryMax) {
		t.Fatalf("expected retries: %d != %d", client.RetryMax, retries)
	}
}
