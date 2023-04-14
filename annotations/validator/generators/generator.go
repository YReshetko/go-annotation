package generators

import (
	"github.com/YReshetko/go-annotation/annotations/validator/annotations"
	"github.com/YReshetko/go-annotation/annotations/validator/model"
	annotation "github.com/YReshetko/go-annotation/pkg"
	"go/ast"
	"go/token"
)

// Generator
// @Constructor
type Generator struct {
	ts   *ast.TypeSpec
	va   annotations.Validator
	st   model.Validatable // @Exclude
	node annotation.Node
}

// @PostConstruct
func (g *Generator) postConstruct() {
	structName := g.ts.Name.String()
	st := g.ts.Type.(*ast.StructType)
	g.st = model.Validatable{
		Name:          structName,
		ValidatorName: g.va.BuildName(structName),
		Fields:        g.fields(st, g.node, false),
	}
}
func (g *Generator) fields(st *ast.StructType, node annotation.Node, checkUnexported bool) []model.Field {
	var vFileds []model.Field
	for i := 0; i < st.Fields.NumFields(); i++ {
		f := st.Fields.List[i]
		if g.va.Required == annotations.Tagged && f.Tag == nil {
			continue
		}

		fType, innerValidator := g.fieldTypeWithValidator(f, node, checkUnexported)
		if fType == "" {
			continue
		}

		var tag *model.Tag
		if f.Tag != nil {
			tag = buildTag(f.Tag)
		}
		if len(f.Names) == 0 {
			name := ""
			if innerValidator != nil {
				name = innerValidator.ValidatorName
			}
			vFileds = append(vFileds, model.Field{
				Name:        name,
				FType:       fType,
				Tag:         tag,
				TypeNode:    f.Type,
				Validatable: innerValidator,
			})
		}

		for _, name := range f.Names {
			if !isExported(name.String()) && checkUnexported {
				continue
			}
			vFileds = append(vFileds, model.Field{
				Name:        name.String(),
				FType:       fType,
				Tag:         tag,
				TypeNode:    f.Type,
				Validatable: innerValidator,
			})
		}
	}
	return vFileds
}

func (g *Generator) fieldTypeWithValidator(field *ast.Field, node annotation.Node, checkUnexported bool) (string, *model.Validatable) {
	//ast.Print(token.NewFileSet(), field)
	switch expr := field.Type.(type) {
	case *ast.Ident:
		name, st, ok := extractStructFromIdent(expr)
		if !ok {
			return expr.Name, nil
		}
		return "static", &model.Validatable{
			ValidatorName: name,
			Name:          "",
			Fields:        g.fields(st, node, checkUnexported),
		}
	case *ast.StarExpr:
		n, _ := unwrapStarExpression(expr)
		name, st, ok := extractStructFromIdent(n)
		if !ok {
			return "pointer", nil
		}
		return "pointer", &model.Validatable{
			ValidatorName: name,
			Name:          "",
			Fields:        g.fields(st, node, checkUnexported),
		}
	case *ast.ArrayType, *ast.MapType, *ast.InterfaceType, *ast.FuncType:
		return "pointer", nil
	case *ast.StructType:
		return "static", &model.Validatable{
			ValidatorName: "",
			Name:          "",
			Fields:        g.fields(expr, node, checkUnexported),
		}
	case *ast.SelectorExpr:
		alias, name, ok := extractExternalType(expr)
		if !ok {
			return "", nil
		}
		newNode, _, err := node.Lookup().FindNodeByAlias(alias, name)
		if err != nil {
			panic(err)
		}
		typeName, st, ok := extractStructFromTypeSpec(newNode.ASTNode())
		return "static", &model.Validatable{
			ValidatorName: typeName,
			Name:          "",
			Fields:        g.fields(st, newNode, true),
		}
		return "", nil
	default:
		ast.Print(token.NewFileSet(), field.Type)
		panic("unable process non [ast.Ident, ast.StarExpr, ast.ArrayType, ast.MapType] field type")
	}
}

func (g *Generator) Validatable() model.Validatable {
	return g.st
}

func unwrapStarExpression(e ast.Node) (ast.Node, bool) {
	o, ok := e.(*ast.StarExpr)
	if ok {
		return o.X, true
	}
	return e, false
}

func extractStructFromIdent(n ast.Node) (string, *ast.StructType, bool) {
	ident, ok := n.(*ast.Ident)
	if !ok || ident.Obj == nil || ident.Obj.Kind != ast.Typ || ident.Obj.Decl == nil {
		return "", nil, false
	}

	ts, ok := ident.Obj.Decl.(*ast.TypeSpec)
	if !ok {
		return "", nil, false
	}
	st, ok := ts.Type.(*ast.StructType)
	return ident.String(), st, ok
}

func extractStructFromTypeSpec(n ast.Node) (string, *ast.StructType, bool) {
	ts, ok := n.(*ast.TypeSpec)
	if !ok {
		return "", nil, false
	}

	st, ok := ts.Type.(*ast.StructType)
	if !ok {
		return "", nil, false
	}
	return ts.Name.String(), st, true
}

func extractExternalType(se *ast.SelectorExpr) (string, string, bool) {
	pkg, ok := se.X.(*ast.Ident)
	if !ok {
		return "", "", false
	}
	return pkg.String(), se.Sel.String(), true
}

func isExported(v string) bool {
	return v[0] >= 'A' && v[0] <= 'Z'
}
