package internal

type ToEmbed struct {
	a int
	b float32
	s *string
}

// SomeStructWithEmbedding @Validator
type SomeStructWithEmbedding struct {
	ToEmbed
}

// AnotherStructWithEmbedding @Validator
type AnotherStructWithEmbedding struct {
	*ToEmbed
}
