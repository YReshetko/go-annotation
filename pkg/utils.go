package pkg

import (
	"go/ast"
)

func CastAnnotation[T Annotation](a Annotation) T {
	t, ok := a.(T)
	if !ok {
		panic("unable parse annotation")
	}
	return t
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
