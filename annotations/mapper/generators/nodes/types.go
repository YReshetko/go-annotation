package nodes

import (
	"github.com/YReshetko/go-annotation/annotations/mapper/generators/cache"
	"go/ast"

	annotation "github.com/YReshetko/go-annotation/pkg"
)

var _ Type = (*PrimitiveType)(nil)
var _ Type = (*StructType)(nil)
var _ Type = (*ArrayType)(nil)
var _ Type = (*MapType)(nil)

type Type interface {
	DeclaredType() string
	IsPointer() bool
}

// @Builder(type="pointer")
type PrimitiveType struct {
	isPointer bool
	name      string
}

func (p *PrimitiveType) IsPointer() bool {
	return p.isPointer
}

func (p *PrimitiveType) Name() string {
	return p.name
}

func (p *PrimitiveType) DeclaredType() string {
	prefix := ""
	if p.isPointer {
		prefix = "*"
	}
	return prefix + p.name
}

// @Builder(type="pointer")
type StructType struct {
	impCache     *cache.ImportCache
	node         annotation.Node
	name         string
	isPointer    bool
	alias        string
	importPath   string
	parentImport string
	astStruct    *ast.StructType

	fields      []*Field // @Init
	importToken string   //@Exclude
}

func (s *StructType) IsPointer() bool {
	return s.isPointer
}

// @PostConstruct
func (s *StructType) buildFields() {
	if s.astStruct.Fields == nil {
		return
	}

	for _, field := range s.astStruct.Fields.List {
		if len(field.Names) == 0 {
			// TODO try to support embedded structure
			continue
		}
		for _, name := range field.Names {
			s.fields = append(s.fields, NewFieldBuilder().
				ParentImport(s.importPath).
				ImpCache(s.impCache).
				Name(name.String()).
				Node(s.node).
				AstExpr(field.Type).
				Build(),
			)
		}
	}
}

// @PostConstruct
func (s *StructType) buildImportToken() {
	s.importToken = s.impCache.StoreImport(s.importPath)
}

func (s *StructType) DeclaredType() string {
	prefix := ""
	if s.isPointer {
		prefix = "*"
	}
	if len(s.importToken) > 0 {
		prefix += "{{ ." + s.importToken + " }}."
	}

	return prefix + s.name
}

func (s *StructType) Fields() []*Field {
	return s.fields
}

func (s *StructType) Equal(s1 *StructType) bool {
	return s.importPath == s1.importPath && s.name == s1.name
}

// @Builder(type="pointer")
type ArrayType struct {
	impCache     *cache.ImportCache
	node         annotation.Node
	isPointer    bool
	astArray     *ast.ArrayType
	parentImport string

	isEllipsis bool // @Exclude
	ofType     Type // @Exclude
}

func (a *ArrayType) IsPointer() bool {
	return a.isPointer
}

// @PostConstruct
func (a *ArrayType) buildType() {
	if a.astArray.Len != nil {
		a.isEllipsis = true
	}

	a.ofType = newType(a.astArray.Elt, a.node, a.impCache, a.parentImport)
}

func (a *ArrayType) DeclaredType() string {
	prefix := ""
	if a.isPointer {
		prefix = "*"
	}
	arrPrefix := "[]"
	if a.isEllipsis {
		arrPrefix = "..."
	}

	return prefix + arrPrefix + a.ofType.DeclaredType()
}

func (a *ArrayType) Equal(a1 *ArrayType) bool {
	return a.DeclaredType() == a1.DeclaredType()
}

// @Builder(type="pointer")
type MapType struct {
	impCache     *cache.ImportCache
	node         annotation.Node
	isPointer    bool
	astKeyExpr   ast.Expr
	astValueExpr ast.Expr
	parentImport string

	keyType   Type // @Exclude
	valueType Type // @Exclude
}

func (m *MapType) IsPointer() bool {
	return m.isPointer
}

// @PostConstruct
func (m *MapType) buildTypes() {
	m.keyType = newType(m.astKeyExpr, m.node, m.impCache, m.parentImport)
	m.valueType = newType(m.astValueExpr, m.node, m.impCache, m.parentImport)
}

func (m *MapType) DeclaredType() string {
	prefix := ""
	if m.isPointer {
		prefix = "*"
	}

	return prefix + "map[" + m.keyType.DeclaredType() + "]" + m.valueType.DeclaredType()
}

func (m *MapType) Equal(m1 *MapType) bool {
	return m.keyType.DeclaredType() == m1.keyType.DeclaredType() && m.valueType.DeclaredType() == m1.valueType.DeclaredType()
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
