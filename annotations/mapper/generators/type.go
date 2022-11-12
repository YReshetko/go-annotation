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
	return buildType(stg.fieldGen)
}

// @Builder(constructor="newMapTypeGeneratorBuilder", build="set{{ .FieldName }}", terminator="build", type="pointer")
type mapTypeGenerator struct {
	node        annotation.Node
	parentAlias string
	keyType     ast.Node
	valueType   ast.Node

	key   *fieldGenerator
	value *fieldGenerator
}

// @PostConstruct
func (mtg *mapTypeGenerator) buildFields() {
	mtg.key = buildFiledGenerator(&ast.Field{Type: mtg.keyType.(ast.Expr)}, "", mtg.node, mtg.parentAlias)
	mtg.value = buildFiledGenerator(&ast.Field{Type: mtg.valueType.(ast.Expr)}, "", mtg.node, mtg.parentAlias)
}

func (mtg *mapTypeGenerator) buildType() string {
	if mtg.key == nil || mtg.value == nil {
		return ""
	}
	key := buildType(mtg.key)
	value := buildType(mtg.value)
	return fmt.Sprintf("map[%s]%s", key, value)
}

func buildType(f *fieldGenerator) string {
	prefix := ""
	if f.isPointer {
		prefix = "*"
	}
	if len(f.alias) == 0 && len(f.parentAlias) != 0 {
		prefix += f.parentAlias + "."
	}
	return prefix + f.buildArgType()
}
