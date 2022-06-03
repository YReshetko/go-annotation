package intpkg

import "fmt"

var (
	// test var group 1
	// @Inject(value="15")
	testVarGroup1 int
	// @Inject(value="35")
	// test var group 2
	testVarGroup2 int
)

// testSingleVar test single var
// @Inject(value="10")
var testSingleVar int

// TestStruct test struct
// @Schema(src="hello", value="world")
// @   RestEndpoint   (
//			path="/api/v1/post"    ,
//			method="GET"
//	)
type TestStruct struct {
	// @Inject(package="test", value="10")
	innerField int
}

// @Endpoint
// TestInterface test interface
type TestInterface interface {
}

/* Print internal print

@Multiline()
*/
func Print() {
	fmt.Println("RUN TEST")
}
