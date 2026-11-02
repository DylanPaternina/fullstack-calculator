package api

import (
	"encoding/json"
	"math"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestAdvancedExpressions(t *testing.T) {
	for _, tc := range []struct {
		expression string
		result     float64
		code       string
	}{
		{"2^3+sqrt(16)", 12, ""}, {"2^3^2", 512, ""}, {"-2^2", -4, ""}, {"(-2)^2", 4, ""},
		{"200+10%", 220, ""}, {"200-10%", 180, ""}, {"200*10%", 20, ""}, {"200/10%", 2000, ""},
		{"200+10%+5%", 231, ""}, {"sqrt(25)*20%", 1, ""}, {"(2+3)*4", 20, ""},
		{"sqrt(-1)", 0, "domain_error"}, {"(-2)^0.5", 0, "domain_error"}, {"0^-1", 0, "division_by_zero"},
		{"10^400", 0, "numeric_overflow"}, {"sqrt()", 0, "invalid_expression"}, {"2^", 0, "invalid_expression"},
		{strings.Repeat("(", 65) + "1" + strings.Repeat(")", 65), 0, "expression_too_complex"},
	} {
		t.Run(tc.expression, func(t *testing.T) {
			body, err := json.Marshal(map[string]string{"expression": tc.expression})
			if err != nil {
				t.Fatal(err)
			}
			w := httptest.NewRecorder()
			NewHandler().ServeHTTP(w, httptest.NewRequest("POST", "/api/calculate", strings.NewReader(string(body))))
			var response struct {
				Result *float64                        `json:"result"`
				Error  *struct{ Code, Message string } `json:"error"`
			}
			if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
				t.Fatal(err)
			}
			if tc.code != "" {
				if w.Code != 400 || response.Error == nil || response.Error.Code != tc.code || response.Error.Message == "" || response.Result != nil {
					t.Fatalf("unexpected error response: %d %s", w.Code, w.Body.String())
				}
			} else if w.Code != 200 || response.Result == nil || math.IsNaN(*response.Result) || math.IsInf(*response.Result, 0) || *response.Result != tc.result || response.Error != nil {
				t.Fatalf("unexpected success response: %d %s", w.Code, w.Body.String())
			}
		})
	}
}
