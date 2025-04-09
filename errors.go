package alloha

import "errors"

var (
	ApiTokenEmptyError          = errors.New("api token is empty")
	BaseApiUrlEmptyError        = errors.New("base api url is empty")
	BaseApiUrlInvalidHostError  = errors.New("base api url host is invalid")
	EmptyEndpointApiURLError    = errors.New("endpoint api url is empty")
	FailedCreateRequestError    = errors.New("failed to create a request object")
	EmptyIMDbIdParameterError   = errors.New("imdb id param is empty")
	InvalidKPIdParameterError   = errors.New("kp id param is invalid")
	InvalidTMDbIdParameterError = errors.New("tmdb id param is invalid")
)
