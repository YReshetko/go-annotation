package internal

// TaggedExample @Validator(require="all")
type TaggedExample struct {
	Name     string `gav:"ignore=true,validator=functionName,range=1..100,enum=typeName" bson:"name"`
	Surname  string
	Age      int
	Route    float32
	Consumer OAuthCredentials
}
