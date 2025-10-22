package cocobase

import "fmt"

type APIError struct {
	StatusCode int
	Method     string
	URL        string
	Body       string
	Suggestion string
}

func (e *APIError) Error() string {
	return fmt.Sprintf("API request failed: %s %s (status: %d)\nBody: %s\nSuggestion: %s",
		e.Method, e.URL, e.StatusCode, e.Body, e.Suggestion)
}

func getErrorSuggestion(status int, method string) string {
	switch status {
	case 401:
		return "Check if your API key is valid and properly set"
	case 403:
		return "You don't have permission to perform this action. Verify your access rights"
	case 404:
		return "The requested resource was not found. Verify the path and ID are correct"
	case 405:
		return fmt.Sprintf("The %s method is not allowed for this endpoint. Check the API documentation for supported methods", method)
	case 429:
		return "You've exceeded the rate limit. Please wait before making more requests"
	default:
		return "Check the API documentation and verify your request format"
	}
}
