package api

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"calculator/internal/evaluator"
)

const MaxBodyBytes = 8192

func writeJSON(w http.ResponseWriter, status int, value interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.Header().Set("Cache-Control", "no-store")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}

func fail(w http.ResponseWriter, status int, code, message string) {
	writeJSON(w, status, struct {
		Error *evaluator.Error `json:"error"`
	}{&evaluator.Error{Code: code, Message: message}})
}

// NewHandler returns an independently testable API with no CORS middleware.
func NewHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/calculate" {
			fail(w, 404, "not_found", "Endpoint not found")
			return
		}
		if r.Method != http.MethodPost {
			w.Header().Set("Allow", "POST")
			fail(w, 405, "method_not_allowed", "Use POST for calculations")
			return
		}
		r.Body = http.MaxBytesReader(w, r.Body, MaxBodyBytes)
		decoder := json.NewDecoder(r.Body)
		decoder.DisallowUnknownFields()
		var request struct {
			Expression *string `json:"expression"`
		}
		err := decoder.Decode(&request)
		if err == nil {
			var extra interface{}
			if tailErr := decoder.Decode(&extra); tailErr != io.EOF {
				if tailErr == nil {
					tailErr = errors.New("extra JSON")
				}
				err = tailErr
			}
		}
		if err != nil {
			var tooLarge *http.MaxBytesError
			if errors.As(err, &tooLarge) {
				fail(w, 413, "request_too_large", "Request body must be at most 8192 bytes")
			} else {
				fail(w, 400, "invalid_request", "Send one JSON object containing a string expression")
			}
			return
		}
		if request.Expression == nil {
			fail(w, 400, "invalid_request", "Expression is required and must be a string")
			return
		}
		result, err := evaluator.Evaluate(*request.Expression)
		if err != nil {
			var calculationError *evaluator.Error
			if errors.As(err, &calculationError) {
				fail(w, 400, calculationError.Code, calculationError.Message)
			} else {
				fail(w, 500, "internal_error", "Calculation failed")
			}
			return
		}
		writeJSON(w, 200, struct {
			Result float64 `json:"result"`
		}{result})
	})
}
