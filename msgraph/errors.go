package msgraph

import (
	"context"
	"errors"
	"net/http"
	"net/url"
	"time"
)

type AuthError struct {
	StatusCode int    `json:"-"` // http status code
	Status     string `json:"-"` // http status
	Code       string `json:"error,omitempty"`
	Message    string `json:"error_description,omitempty"`
	RetryAfter time.Duration
}

func AsAuthError(err error) (ae *AuthError, ok bool) {
	ok = errors.As(err, &ae)
	return
}

func IsAuthError(err error) bool {
	_, ok := AsAuthError(err)
	return ok
}

func newAuthError(res *http.Response) *AuthError {
	return &AuthError{
		StatusCode: res.StatusCode,
		Status:     res.Status,
	}
}

func (ae *AuthError) GetRetryAfter() time.Duration {
	return ae.RetryAfter
}

func (ae *AuthError) Error() string {
	es := ae.Status

	if ae.RetryAfter > 0 {
		es += " (Retry After " + ae.RetryAfter.String() + ")"
	}

	es += " - " + ae.Code + ": " + ae.Message

	return es
}

type DetailError struct {
	Code    string `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
}

func (de *DetailError) Error() string {
	return de.Code + ": " + de.Message
}

type ResultError struct {
	Method     string       `json:"-"` // http request method
	URL        *url.URL     `json:"-"` // http request URL
	StatusCode int          `json:"-"` // http status code
	Status     string       `json:"-"` // http status
	Detail     *DetailError `json:"error,omitempty"`
	RetryAfter time.Duration
}

func AsResultError(err error) (re *ResultError, ok bool) {
	ok = errors.As(err, &re)
	return
}

func IsResultError(err error) bool {
	_, ok := AsResultError(err)
	return ok
}

func newResultError(res *http.Response) *ResultError {
	return &ResultError{
		Method:     res.Request.Method,
		URL:        res.Request.URL,
		StatusCode: res.StatusCode,
		Status:     res.Status,
	}
}

func (re *ResultError) GetRetryAfter() time.Duration {
	return re.RetryAfter
}

func (re *ResultError) Error() string {
	es := re.Status

	if re.RetryAfter > 0 {
		es += " (Retry After " + re.RetryAfter.String() + ")"
	}

	es += " (" + re.Method + " " + re.URL.String() + ")"

	if re.Detail != nil {
		es += " - " + re.Detail.Error()
	}

	return es
}

func shouldRetry(err error) bool {
	if re, ok := AsResultError(err); ok {
		return re.StatusCode == http.StatusTooManyRequests || (re.StatusCode >= 500 && re.StatusCode <= 599)
	}
	return !errors.Is(err, context.Canceled)
}
