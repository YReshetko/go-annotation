package pkg

import (
	"fmt"
	"go/ast"
	"strings"
)

func CastAnnotation[T Annotation](a Annotation) T {
	t, ok := TryCastAnnotation[T](a)
	if !ok {
		panic("unable parse annotation")
	}
	return t
}

func TryCastAnnotation[T Annotation](a Annotation) (T, bool) {
	t, ok := a.(T)
	return t, ok
}

func MethodReceiver(decl *ast.FuncDecl) string {
	if decl.Recv == nil {
		return ""
	}

	for _, v := range decl.Recv.List {
		switch rv := v.Type.(type) {
		case *ast.StarExpr:
			return rv.X.(*ast.Ident).Name
		case *ast.UnaryExpr:
			return rv.X.(*ast.Ident).Name
		}
	}
	return ""
}

//==================== Function parser ======================

func FunctionSignature(n ast.Node) (FuncSignature, error) {
	f, ok := n.(*ast.Field)
	if !ok {
		return FuncSignature{}, fmt.Errorf("expected *ast.Field, but got: %T", n)
	}
	ft, ok := f.Type.(*ast.FuncType)
	if !ok {
		return FuncSignature{}, fmt.Errorf("expected *ast.Field, but got: %T", f.Type)
	}
	if len(f.Names) != 1 {
		return FuncSignature{}, fmt.Errorf("unexpected interface method signature")
	}
	return FuncSignature{
		Name:    f.Names[0].Name,
		Params:  ParamList(ft.Params),
		Results: ParamList(ft.Results),
	}, nil
}

func ParamList(pl *ast.FieldList) []NodeField {
	var out []NodeField
	for _, field := range pl.List {
		out = append(out, Param(field)...)
	}
	return out
}

func Param(p *ast.Field) []NodeField {
	var out []NodeField
	n := names(p.Names)
	for _, name := range n {
		fp := NodeField{
			Exp:  p.Type,
			Name: name,
		}
		out = append(out, funcParam(p.Type, fp))
	}
	return out
}

func names(n []*ast.Ident) []string {
	if len(n) == 0 {
		return []string{""}
	}
	out := make([]string, len(n))
	for i, ident := range n {
		out[i] = ident.Name
	}
	return out
}

func funcParam(expr ast.Expr, fp NodeField) NodeField {
	switch e := expr.(type) {
	case *ast.StarExpr:
		fp.IsPointer = true
		fp = funcParam(e.X, fp)
	case *ast.SliceExpr:
		fp.FieldType = ArrayFieldType
		v := funcParam(e.X, NodeField{})
		fp.Value = &v
	case *ast.ArrayType:
		fp.FieldType = ArrayFieldType
		v := funcParam(e.Elt, NodeField{})
		fp.Value = &v
	case *ast.MapType:
		fp.FieldType = MapFieldType
		key := funcParam(e.Key, NodeField{})
		value := funcParam(e.Value, NodeField{})
		fp.Key, fp.Value = &key, &value
	case *ast.SelectorExpr:
		if e.X != nil {
			si, ok := e.X.(*ast.Ident)
			if ok {
				fp.Selector = si.Name
			}
		}
		if e.Sel != nil {
			_, ok := primitives[e.Sel.Name]
			if ok {
				fp.FieldType = BasicFieldType
			} else {
				fp.FieldType = SelectorFieldType
			}
			fp.TypeName = e.Sel.Name
		}
	case *ast.Ident:
		_, ok := primitives[e.Name]
		if ok {
			fp.FieldType = BasicFieldType
		} else {
			fp.FieldType = SelectorFieldType
		}
		fp.TypeName = e.Name
	case *ast.StructType:
		fp.FieldType = StructureFieldType
	case *ast.InterfaceType:
		fp.FieldType = InterfaceFieldType
	case *ast.FuncType:
		fp.FieldType = FunctionFieldType
	}
	return fp
}

type FuncSignature struct {
	Name    string
	Params  []NodeField
	Results []NodeField
}

type NodeFieldType string

const (
	InterfaceFieldType = "interface"
	StructureFieldType = "struct"
	FunctionFieldType  = "function"
	ArrayFieldType     = "array"
	MapFieldType       = "map"
	BasicFieldType     = "basic"
	SelectorFieldType  = "selector"
)

type NodeField struct {
	Name      string
	IsPointer bool
	Selector  string
	TypeName  string
	FieldType NodeFieldType
	Exp       ast.Expr
	Key       *NodeField
	Value     *NodeField
}

var primitives = map[string]struct{}{
	"bool":       {},
	"uint8":      {},
	"uint16":     {},
	"uint32":     {},
	"uint64":     {},
	"int8":       {},
	"int16":      {},
	"int32":      {},
	"int64":      {},
	"float32":    {},
	"float64":    {},
	"complex64":  {},
	"complex128": {},
	"string":     {},
	"int":        {},
	"uint":       {},
	"uintptr":    {},
	"byte":       {},
	"rune":       {},
	"any":        {},
}

//========================== Import =================

func FindImport(f *ast.File, alias string) string {
	var found bool
	var out string
	ast.Inspect(f, func(node ast.Node) bool {
		if found {
			return false
		}
		imp, ok := node.(*ast.ImportSpec)
		if !ok {
			return true
		}
		impPath := unquote(imp.Path.Value)
		if imp.Name != nil && imp.Name.Name == alias {
			found = true
			out = impPath
			return false
		}
		if strings.HasSuffix(impPath, "/"+alias) || strings.HasSuffix(strings.ReplaceAll(impPath, "-", "_"), "/"+alias) {
			found = true
			out = impPath
			return false
		}
		return true
	})

	return out
}

func unquote(s string) string {
	out := strings.TrimSpace(s)
	if out[0] == '"' && out[len(out)-1] == '"' {
		return out[1 : len(out)-1]
	}
	return out
}
