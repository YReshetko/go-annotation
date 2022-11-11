package generators

import (
	"fmt"
	"go/ast"

	annotation "github.com/YReshetko/go-annotation/pkg"
)

type typeGenerator interface {
}

// @Builder(constructor="newPrimitiveTypeGeneratorBuilder", build="set{{ .FieldName }}", terminator="build", type="pointer")
type primitiveTypeGenerator struct {
	name string
}

// @PostConstruct
func (ptg *primitiveTypeGenerator) postConstruct() {
	//fmt.Println("created primitive field:", ptg.name)
}

// @Builder(constructor="newStructureTypeGeneratorBuilder", build="set{{ .FieldName }}", terminator="build", type="pointer")
type structureTypeGenerator struct {
	node   annotation.Node
	name   string
	fields []*ast.Field

	fieldGenerators []*fieldGenerator // @Init
}

// @PostConstruct
func (stg *structureTypeGenerator) buildFields() {
	//fmt.Println("created structure field:", stg.name)
	for i, field := range stg.fields {
		if len(field.Names) == 0 {
			stg.fieldGenerators = append(stg.fieldGenerators, buildFiledGenerator(field, fmt.Sprintf("in%d", i), stg.node))
		} else {
			for _, name := range field.Names {
				stg.fieldGenerators = append(stg.fieldGenerators, buildFiledGenerator(field, name.String(), stg.node))
			}
		}
	}
}

func isPrimitive(t string) bool {
	_, ok := map[string]struct{}{
		"bool":       {},
		"uint":       {},
		"uint8":      {},
		"uint16":     {},
		"uint32":     {},
		"uint64":     {},
		"byte":       {},
		"int":        {},
		"int8":       {},
		"int16":      {},
		"int32":      {},
		"int64":      {},
		"float32":    {},
		"float64":    {},
		"complex64":  {},
		"complex128": {},
		"string":     {},
		"uintptr":    {},
		"rune":       {},
	}[t]
	return ok
}
