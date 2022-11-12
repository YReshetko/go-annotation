package output

import (
	"github.com/YReshetko/go-annotation/examples/mapper/internal/models/common"
	"github.com/YReshetko/go-annotation/examples/mapper/internal/models/unused"
)

type MapStruct struct {
	Map  map[common.MapKey]common.Common2
	None unused.Unused
}

type MapStruct2 struct {
	Map map[common.MapKey]*common.Common2
}
