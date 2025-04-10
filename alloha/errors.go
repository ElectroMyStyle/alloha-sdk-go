package alloha

import (
	"errors"
	"fmt"
)

var (
	ApiTokenEmptyError              = errors.New("api token is empty")
	BaseApiUrlEmptyError            = errors.New("base api url is empty")
	BaseApiUrlInvalidHostError      = errors.New("base api url host is invalid")
	EmptyEndpointApiURLError        = errors.New("endpoint api url is empty")
	EmptyIMDbIdParameterError       = errors.New("imdb id param is empty")
	EmptyHttpMethodError            = errors.New("http method param is empty")
	EmptyMovieNameParameterError    = errors.New("movie name param is empty")
	FailedCreateRequestError        = errors.New("failed to create a request object")
	InvalidKPIdParameterError       = errors.New("kp id param is invalid")
	InvalidTMDbIdParameterError     = errors.New("tmdb id param is invalid")
	InvalidPageNumberParameterError = errors.New("page number param is invalid")
)

// EmptyResponseBodyError represents an error when the response body is empty
type EmptyResponseBodyError struct {
	StatusCode int
}

// Error implements the error interface
func (e *EmptyResponseBodyError) Error() string {
	return fmt.Sprintf("empty request response body, with StatusCode: %d", e.StatusCode)
}
