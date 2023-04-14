package internal

// User @Validator(require="all")
type User struct {
	Name     string
	Surname  string
	Age      int
	Route    float32
	Consumer OAuthCredentials
}

type OAuthCredentials struct {
	Key    string
	Secret string
}
