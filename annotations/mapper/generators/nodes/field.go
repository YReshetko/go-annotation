package nodes

import (
	"errors"
	"fmt"
	"github.com/YReshetko/go-annotation/annotations/mapper/generators/cache"
	"go/ast"

	annotation "github.com/YReshetko/go-annotation/pkg"
)

// @Builder(type="pointer")
type Field struct {
	impCache     *cache.ImportCache
	node         annotation.Node
	name         string
	astExpr      ast.Expr
	parentImport string

	nodeType Type // @Exclude
}

// @PostConstruct
func (f *Field) buildType() {
	f.nodeType = newType(f.astExpr, f.node, f.impCache, f.parentImport)
}

func (f *Field) Name() string {
	return f.name
}

func (f *Field) DeclaredType() string {
	return f.nodeType.DeclaredType()
}

func (f *Field) Declaration() string {
	decType := f.nodeType.DeclaredType()
	if decType[0] == '*' {
		decType = "&" + decType[1:]
	}
	return decType + "{}"
}

func (f *Field) Type() Type {
	return f.nodeType
}

func (f *Field) IsExported() bool {
	return len(f.name) == 0 || (f.name[0] >= 'A' && f.name[0] <= 'Z')
}

var WrongPathErr = errors.New("wrong name path")

func (f *Field) FindNilChecks(prefix string, names []string) ([]string, bool, error) {
	if len(names) == 0 {
		return nil, false, fmt.Errorf("names are empty")
	}
	if names[0] != f.name {
		return nil, false, WrongPathErr
	}
	if len(names) == 1 {
		if f.Type().IsPointer() {
			return []string{VariableNameJoin(prefix, f.name)}, true, nil
		}
		return nil, false, nil
	}

	st, ok := f.nodeType.(*StructType)
	if !ok {
		return nil, false, fmt.Errorf("unable to get path for non structure field")
	}

	var out []string
	if f.Type().IsPointer() {
		out = append(out, VariableNameJoin(prefix, f.name))
	}

	for _, field := range st.fields {
		ch, isPointer, err := field.FindNilChecks(VariableNameJoin(prefix, f.name), names[1:])
		if errors.Is(err, WrongPathErr) {
			continue
		}
		if err != nil {
			return nil, false, err
		}
		out = append(out, ch...)
		return out, isPointer, nil
	}
	return nil, false, WrongPathErr
}

func VariableNameJoin(prefix, name string) string {
	if len(prefix) == 0 {
		return name
	}
	return prefix + "." + name
}

func newType(astExpr ast.Expr, node annotation.Node, impCache *cache.ImportCache, parentImport string) Type {
	// Unsupported pointer to pointer expressions
	unwrappedStart, isPointer := unwrapStarExpression(astExpr)
	unwrappedSelector, alias := unwrapSelectorExpression(unwrappedStart)

	switch expr := unwrappedSelector.(type) {
	case *ast.Ident:
		return newTypeByIdent(expr, node, alias, isPointer, impCache, parentImport)
	case *ast.ArrayType:
		return NewArrayTypeBuilder().
			Node(node).
			ImpCache(impCache).
			ParentImport(parentImport).
			AstArray(expr).
			IsPointer(isPointer).
			Build()
	case *ast.MapType:
		return NewMapTypeBuilder().
			Node(node).
			ImpCache(impCache).
			ParentImport(parentImport).
			AstKeyExpr(expr.Key).
			AstValueExpr(expr.Value).
			IsPointer(isPointer).
			Build()
	default:
		fmt.Printf("unsupported field type: %T\n", expr)
		//ast.Print(token.NewFileSet(), expr)
	}
	return nil
}

func newTypeByIdent(ident *ast.Ident, node annotation.Node, alias string, isPointer bool, impCache *cache.ImportCache, parentImport string) Type {
	if ident.Obj != nil {
		// TODO support embedded structures and declared in the same file
		if ident.Obj.Kind != ast.Typ {
			fmt.Println("TODO: unsupported object declaration:", ident.Obj.Kind)
			//ast.Print(token.NewFileSet(), ident.Obj)
			return nil
		}
		astNode, ok := ident.Obj.Decl.(ast.Node)
		if !ok {
			fmt.Printf("TODO: unsupported object declaration: %t\n", ident.Obj.Decl)
			return nil
		}

		return newNonPrimitiveType(astNode, node, isPointer, alias, parentImport, impCache, parentImport)
	}

	if len(alias) == 0 && isPrimitive(ident.String()) {
		return NewPrimitiveTypeBuilder().IsPointer(isPointer).Name(ident.String()).Build()
	}

	newNode, importPath, err := node.Lookup().FindNodeByAlias(alias, ident.String())
	if err != nil {
		fmt.Println("WARN:", err.Error())
		return nil
	}

	if len(parentImport) != 0 && len(importPath) == 0 {
		fmt.Printf("replacin import: %s\n", parentImport)
		importPath = parentImport
	}

	return newNonPrimitiveType(newNode.ASTNode(), newNode, isPointer, alias, importPath, impCache, parentImport)
}

func newNonPrimitiveType(astNode ast.Node, node annotation.Node, isPointer bool, alias, importPath string, impCache *cache.ImportCache, parentImport string) Type {
	switch nnt := astNode.(type) {
	case *ast.TypeSpec:
		name := nnt.Name.String()
		switch innt := nnt.Type.(type) {
		case *ast.StructType:
			return NewStructTypeBuilder().
				ImpCache(impCache).
				ParentImport(parentImport).
				Node(node).
				Name(name).
				IsPointer(isPointer).
				Alias(alias).
				ImportPath(importPath).
				AstStruct(innt).
				Build()
		default:
			fmt.Printf("unsupported internal loaded type %T\n", nnt.Type)
			//ast.Print(token.NewFileSet(), astNode)
		}
	default:
		fmt.Printf("unsupported internal type %T\n", astNode)
		//ast.Print(token.NewFileSet(), astNode)
	}
	return nil
}

func unwrapSelectorExpression(e ast.Node) (ast.Node, string) {
	o, ok := e.(*ast.SelectorExpr)
	if ok {
		return o.Sel, o.X.(*ast.Ident).String()
	}
	return e, ""
}

func unwrapStarExpression(e ast.Node) (ast.Node, bool) {
	o, ok := e.(*ast.StarExpr)
	if ok {
		return o.X, true
	}
	return e, false
}
