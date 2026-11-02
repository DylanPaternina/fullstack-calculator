package api

import (
	"encoding/json"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestAPI(t *testing.T) {
	for _, tc := range []struct {
		name, method, path, body string
		status                   int
		code                     string
	}{
		{"success", "POST", "/api/calculate", `{"expression":"5+3*2"}`, 200, ""},
		{"missing", "POST", "/api/calculate", `{}`, 400, "invalid_request"},
		{"null", "POST", "/api/calculate", `null`, 400, "invalid_request"},
		{"null field", "POST", "/api/calculate", `{"expression":null}`, 400, "invalid_request"},
		{"number field", "POST", "/api/calculate", `{"expression":2}`, 400, "invalid_request"},
		{"array field", "POST", "/api/calculate", `{"expression":[]}`, 400, "invalid_request"},
		{"unknown field", "POST", "/api/calculate", `{"expression":"2","extra":true}`, 400, "invalid_request"},
		{"invalid json", "POST", "/api/calculate", `{`, 400, "invalid_request"},
		{"array", "POST", "/api/calculate", `[]`, 400, "invalid_request"},
		{"empty", "POST", "/api/calculate", ``, 400, "invalid_request"},
		{"multiple", "POST", "/api/calculate", `{"expression":"2"} {}`, 400, "invalid_request"},
		{"trailing garbage", "POST", "/api/calculate", `{"expression":"2"} x`, 400, "invalid_request"},
		{"whitespace tail", "POST", "/api/calculate", `{"expression":"5+3*2"} `, 200, ""},
		{"empty expression", "POST", "/api/calculate", `{"expression":""}`, 400, "invalid_expression"},
		{"invalid expression", "POST", "/api/calculate", `{"expression":"5+*2"}`, 400, "invalid_expression"},
		{"division", "POST", "/api/calculate", `{"expression":"5/0"}`, 400, "division_by_zero"},
		{"overflow", "POST", "/api/calculate", `{"expression":"1e308*2"}`, 400, "numeric_overflow"},
		{"expression limit", "POST", "/api/calculate", `{"expression":"` + strings.Repeat("1", 1025) + `"}`, 400, "expression_too_large"},
		{"body limit", "POST", "/api/calculate", strings.Repeat(" ", MaxBodyBytes+1), 413, "request_too_large"},
		{"tail body limit", "POST", "/api/calculate", `{"expression":"2"}` + strings.Repeat(" ", MaxBodyBytes), 413, "request_too_large"},
		{"method", "GET", "/api/calculate", "", 405, "method_not_allowed"},
		{"preflight", "OPTIONS", "/api/calculate", "", 405, "method_not_allowed"},
		{"route", "POST", "/missing", "", 404, "not_found"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			r := httptest.NewRequest(tc.method, tc.path, strings.NewReader(tc.body))
			NewHandler().ServeHTTP(w, r)
			if w.Code != tc.status {
				t.Fatalf("status %d, body %s", w.Code, w.Body.String())
			}
			if w.Header().Get("Content-Type") != "application/json" || w.Header().Get("Access-Control-Allow-Origin") != "" {
				t.Fatal("incorrect headers")
			}
			if tc.status == 405 && w.Header().Get("Allow") != "POST" {
				t.Fatal("missing Allow")
			}
			var response struct {
				Result *float64                        `json:"result"`
				Error  *struct{ Code, Message string } `json:"error"`
			}
			if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
				t.Fatal(err)
			}
			if tc.status == 200 {
				if response.Result == nil || *response.Result != 11 || response.Error != nil {
					t.Fatal("incorrect result")
				}
			} else {
				if response.Error == nil || response.Error.Code != tc.code || response.Error.Message == "" || response.Result != nil {
					t.Fatalf("incorrect error %s", w.Body.String())
				}
			}
		})
	}
}

func TestAPIZeroResult(t *testing.T) {
	for _, tc := range []struct {
		name, expression string
	}{
		{"zero division", "0/5"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			body, err := json.Marshal(map[string]string{"expression": tc.expression})
			if err != nil {
				t.Fatal(err)
			}
			w := httptest.NewRecorder()
			NewHandler().ServeHTTP(w, httptest.NewRequest("POST", "/api/calculate", strings.NewReader(string(body))))
			if w.Code != 200 {
				t.Fatalf("status %d, body %s", w.Code, w.Body.String())
			}
			var response struct {
				Result *float64        `json:"result"`
				Error  json.RawMessage `json:"error"`
			}
			if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
				t.Fatal(err)
			}
			if response.Result == nil || *response.Result != 0 || response.Error != nil {
				t.Fatalf("expected a present, non-null zero result with no error field: %s", w.Body.String())
			}
		})
	}
}
