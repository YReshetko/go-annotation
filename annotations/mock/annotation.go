package mock

type Mock struct {
	Name       string `annotation:"name=name,default={{ .TypeName }}Mock"`
	SubPackage string `annotation:"name=sub,default=mocks"`
}
