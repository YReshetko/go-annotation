package pkg

type SomeStructInPkg struct {
	Name      string
	Name2     *string
	name2     *string
	Surname   string
	Age       int
	age       int
	Route     float32
	route     float32
	Consumer  AnotherStructInPkg
	consumer  AnotherStructInPkg
	Consumer2 *AnotherStructInPkg
	consumer2 *AnotherStructInPkg
}

type AnotherStructInPkg struct {
	Name    string
	Name2   *string
	Surname string
	Age     []int
	Route   float32
	Asip    AnotherStructInPkg2
	Asip2   *AnotherStructInPkg2
	asip2   *AnotherStructInPkg2
}

type AnotherStructInPkg2 struct {
	ame    string
	ame2   *string
	urname string
	ge     []int
	Route  float32
}
