package fixtures

import (
	"fmt"
	"github.com/YReshetko/go-annotation/internal/lookup/fixtures/dashed-package"
	mlog "log"

	. "github.com/davecgh/go-spew/spew"
	_ "github.com/davecgh/go-spew/spew"

	"github.com/davecgh/go-spew/spew"
)

func SomeTestFunction() {
	fmt.Println("msg")
	mlog.Println("msg")

	_ = Config

	_ = anythingelse.Exported{}

	_ = spew.Config
}
