package pkg

import (
	"fmt"
	"go/ast"
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

func ParamList(pl *ast.FieldList) []FuncParam {
	var out []FuncParam
	for _, field := range pl.List {
		n := names(field.Names)
		for _, name := range n {
			fp := FuncParam{
				Exp:  field.Type,
				Name: name,
			}
			out = append(out, funcParam(field.Type, fp))
		}
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

func funcParam(expr ast.Expr, fp FuncParam) FuncParam {
	switch e := expr.(type) {
	case *ast.StarExpr:
		fp.IsPointer = true
		fp = funcParam(e.X, fp)
	case *ast.SelectorExpr:
		if e.X != nil {
			si, ok := e.X.(*ast.Ident)
			if ok {
				fp.Selector = si.Name
			}
		}
		if e.Sel != nil {
			_, ok := primitives[e.Sel.Name]
			fp.TypeName = e.Sel.Name
			fp.IsBasicType = ok
		}
	case *ast.Ident:
		_, ok := primitives[e.Name]
		fp.TypeName = e.Name
		fp.IsBasicType = ok
	case *ast.StructType:
		fp.TypeName = "struct"
	case *ast.InterfaceType:
		fp.TypeName = "interface"
	case *ast.FuncType:
		fp.TypeName = "function"
	}
	return fp
}

type FuncSignature struct {
	Name    string
	Params  []FuncParam
	Results []FuncParam
}

type FuncParam struct {
	Name        string
	IsPointer   bool
	IsBasicType bool
	Selector    string
	TypeName    string
	Exp         ast.Expr
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
