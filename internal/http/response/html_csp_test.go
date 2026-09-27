// SPDX-FileCopyrightText: Copyright The Miniflux Authors. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package response // import "miniflux.app/v2/internal/http/response"

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHTMLWithCSPResponse(t *testing.T) {
	r, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	w := httptest.NewRecorder()

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		HTMLWithCSP(w, r, "Some HTML", "default-src 'none'; script-src 'nonce-abc';")
	})

	handler.ServeHTTP(w, r)
	resp := w.Result()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf(`Unexpected status code, got %d instead of %d`, resp.StatusCode, http.StatusOK)
	}

	if actualBody := w.Body.String(); actualBody != `Some HTML` {
		t.Fatalf(`Unexpected body, got %s instead of %s`, actualBody, `Some HTML`)
	}

	headers := map[string]string{
		"Content-Type":            "text/html; charset=utf-8",
		"Cache-Control":           "no-cache, max-age=0, must-revalidate, no-store",
		"Content-Security-Policy": "default-src 'none'; script-src 'nonce-abc';",
	}

	for header, expected := range headers {
		if actual := resp.Header.Get(header); actual != expected {
			t.Fatalf(`Unexpected header value, got %q instead of %q`, actual, expected)
		}
	}
}
