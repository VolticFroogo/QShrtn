package model

type Redirect struct {
	ID, URL string
}

const (
	ResponseSuccess = iota
	ResponseInternalServerError
	ResponseForbiddenDomain
	ResponseIDTaken
	ResponseInvalidURL
)

const (
	IDLength = 4
)

// Code is the struct used in JSON responses with just a code.
type Code struct {
	Code int
}
