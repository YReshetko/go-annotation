package internal

// @Constructor(name="NewSomeStructureOverride")
type SomeStructure struct {
	a int
	b float64
	c *bool
	d **complex128
}

// @Constructor(name="NewAnotherStructOverride", type="pointer")
type AnotherStruct struct {
	a    SomeStructure
	b    *SomeStructure
	c, d int
	fn   func(**SomeStructure) AnotherStruct
}
