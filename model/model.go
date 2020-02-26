package model

type Redirect struct {
	ID  string `bson:"_id"`
	URL string `bson:"url"`
}

const (
	ResponseSuccess = iota
	ResponseInternalServerError
	ResponseForbiddenDomain
	ResponseIDTaken
	ResponseInvalidURL
)

const (
	IDLength     = 4
	MaxURLLength = 2048
)

// Code is the struct used in JSON responses with just a code.
type Code struct {
	Code int
}
