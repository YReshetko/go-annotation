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
	node        annotation.Node
	name        string
	parentAlias string
	fields      []*ast.Field

	fieldGenerators []*fieldGenerator // @Init
}

// @PostConstruct
func (stg *structureTypeGenerator) buildFields() {
	//fmt.Println("created structure field:", stg.name)
	for i, field := range stg.fields {
		if len(field.Names) == 0 {
			stg.fieldGenerators = append(stg.fieldGenerators, buildFiledGenerator(field, fmt.Sprintf("in%d", i), stg.node, stg.parentAlias))
		} else {
			for _, name := range field.Names {
				stg.fieldGenerators = append(stg.fieldGenerators, buildFiledGenerator(field, name.String(), stg.node, stg.parentAlias))
			}
		}
	}
}

// @Builder(constructor="newSliceTypeGeneratorBuilder", build="set{{ .FieldName }}", terminator="build", type="pointer")
type sliceTypeGenerator struct {
	node        annotation.Node
	parentAlias string
	sliceType   ast.Node

	fieldGen *fieldGenerator //@Exclude
}

// @PostConstruct
func (stg *sliceTypeGenerator) buildFields() {
	stg.fieldGen = buildFiledGenerator(&ast.Field{Type: stg.sliceType.(ast.Expr)}, "", stg.node, stg.parentAlias)
}

func (stg *sliceTypeGenerator) buildType() string {
	if stg.fieldGen == nil {
		return ""
	}
	prefix := ""
	if stg.fieldGen.isPointer {
		prefix = "*"
	}
	if len(stg.fieldGen.alias) == 0 && len(stg.fieldGen.parentAlias) != 0 {
		prefix += stg.fieldGen.parentAlias + "."
	}

	return prefix + stg.fieldGen.buildArgType()
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
