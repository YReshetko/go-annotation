package fixtures

import (
	"fmt"
	mlog "log"

	. "github.com/davecgh/go-spew/spew"
	_ "github.com/davecgh/go-spew/spew"

	"github.com/davecgh/go-spew/spew"
	
	"github.com/YReshetko/go-annotation/internal/lookup/fixtures/dashed-package"
)

func SomeTestFunction() {
	fmt.Println("msg")
	mlog.Println("msg")

	_ = Config

	_ = anythingelse.Exported{}

	_ = spew.Config
}
