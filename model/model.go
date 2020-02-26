package model

type Redirect struct {
	ID  string `bson:"_id"`
	URL string `bson:"url"`
}

const (
	IDLength     = 4
	MaxURLLength = 2048
)
