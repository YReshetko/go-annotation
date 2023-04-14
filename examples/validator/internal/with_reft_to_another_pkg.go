package internal

import "github.com/YReshetko/go-annotation/examples/validator/internal/pkg"

type myType pkg.SomeStructInPkg

// SomeStructWithRef @Validator
type SomeStructWithRef struct {
	ssip  pkg.SomeStructInPkg
	ssip2 *pkg.SomeStructInPkg
}

// SomeStructWithEmbeddingRef @Validator
type SomeStructWithEmbeddingRef struct {
	pkg.SomeStructInPkg
	ssip2 *pkg.SomeStructInPkg
}

// SomeStructWithMyTypeRef @Validator
type SomeStructWithMyTypeRef struct {
	ssip  myType
	ssip2 *myType
}
