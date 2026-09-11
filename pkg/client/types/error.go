package types

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type APIErrorResponse struct {
	Errors []struct {
		Code   string `json:"code"`
		Status string `json:"status"`
		Detail string `json:"detail"`
		Source struct {
			Field string `json:"field"`
		} `json:"source"`
	} `json:"errors"`
}

type APIError struct {
	StatusCode int
	Message    string
}

func (e *APIError) Error() string {
	return fmt.Sprintf("api error (status %d): %s", e.StatusCode, e.Message)
}

func CheckResponse(resp *http.Response) error {
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		return nil
	}

	body, _ := io.ReadAll(resp.Body)

	var errResp APIErrorResponse
	if err := json.Unmarshal(body, &errResp); err == nil && len(errResp.Errors) > 0 {
		var msgs []string
		for _, e := range errResp.Errors {
			if e.Source.Field != "" {
				msgs = append(msgs, fmt.Sprintf("%s (%s: %s)", e.Detail, e.Source.Field, e.Code))
			} else {
				msgs = append(msgs, fmt.Sprintf("%s (%s)", e.Detail, e.Code))
			}
		}
		return &APIError{
			StatusCode: resp.StatusCode,
			Message:    strings.Join(msgs, "; "),
		}
	}

	return &APIError{
		StatusCode: resp.StatusCode,
		Message:    fallbackStatusMessage(resp.StatusCode),
	}
}

func fallbackStatusMessage(code int) string {
	switch code {
	case http.StatusBadRequest:
		return "invalid input data"
	case http.StatusUnauthorized:
		return "invalid API key"
	case http.StatusForbidden:
		return "insufficient permissions"
	case http.StatusNotFound:
		return "server does not exist"
	case http.StatusUnprocessableEntity:
		return "invalid field values"
	case http.StatusTooManyRequests:
		return "rate limit exceeded"
	default:
		return http.StatusText(code)
	}
}
