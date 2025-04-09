package alloha

import "errors"

var (
	ApiTokenEmptyError         = errors.New("api token is empty")
	BaseApiUrlEmptyError       = errors.New("base api url is empty")
	BaseApiUrlInvalidHostError = errors.New("base api url host is invalid")
	EmptyEndpointApiURLError   = errors.New("endpoint api url is empty")
	FailedCreateRequestError   = errors.New("failed to create a request object")
	InvalidKPIdParameterError  = errors.New("kp id param is invalid")
)
