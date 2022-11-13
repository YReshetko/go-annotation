package fixtures

import "fmt"

// SingleInterface single line comment
type SingleInterface interface {
	SingleConst()
}

/*
multiline comment
new line
*/
type SingleStruct struct {
	SingleVar int
}

/*Adding some long comments
  Adding some long comments
  Adding some long comments
  Adding some long comments
  Adding some long comments
  Adding some long comments
*/
type ExploreParentStructs struct {

	/*Adding some long comments
	  Adding some long comments
	  Adding some long comments
	  Adding some long comments
	  Adding some long comments
	  Adding some long comments
	*/
	F struct {
		/*Adding some long comments
		  Adding some long comments
		  Adding some long comments
		  Adding some long comments
		  Adding some long comments
		  Adding some long comments
		*/
		F struct {
			/*Adding some long comments
			  Adding some long comments
			  Adding some long comments
			  Adding some long comments
			  Adding some long comments
			  Adding some long comments
			*/
			LookingFor int
		}
	}
}

type (
	// @SomeMeta() single line comment
	GroupInterface1 interface {
		GroupConst1()
		GroupConst2()
	}

	// Several single line comments
	// Several single line comments
	GroupStruct1 struct {
		GroupVar1 int
	}

	//
	GroupInterface2 interface {
		GroupConst1()
		// Non top lavel node comment
		GroupConst2()
	}

	GroupStruct2 struct {
		/*
			Multiline comment on
			GroupStruct2.GroupVar2
		*/
		GroupVar2 int
	}
)

// Single line comment on constant
const SingleConst = "some string"

/*
Multiline comment on
Variables
*/
var SeveralVars1, SeveralVars2 string

const (
	/*
		Multiline comment GroupConst1
		Variables
	*/
	GroupConst1 = 1
	// Single line comment on GroupConst2
	GroupConst2 = "str"
)

var SingleVar string

var (
	// Single line comment on GroupVar1
	GroupVar1 int
	GroupVar2 bool
	/*
	   Multiline comment GroupSeveralVars1 and GroupSeveralVars2
	   Variables
	*/
	GroupSeveralVars1, GroupSeveralVars2 bool
)

// Single line comment on SomeFunction
func SomeFunction(a int) error {
	var SingleVariable int
	fmt.Println(SingleVariable)

	InternalFunction := func(a int) {

	}
	InternalFunction(a)
	fmt.Println(SingleConst)
	fmt.Println(GroupConst1)
	fmt.Println(GroupConst2)
	return nil
}

/*
Multiline comment on SomeMethod
Method
*/
func (GroupStruct2) SomeMethod(b int) error {
	const SingleConstant = 10
	fmt.Println(SingleConstant)
	fmt.Println(SingleVar)
	fmt.Println(GroupVar1)
	fmt.Println(GroupVar2)
	fmt.Println(SeveralVars1)
	fmt.Println(SeveralVars2)
	fmt.Println(GroupSeveralVars1)
	fmt.Println(GroupSeveralVars2)
	return nil
}

type ExploreParents interface {
	InternalMethod()
}
