package pkg

import (
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
