package generators

import (
	"fmt"
	"github.com/YReshetko/go-annotation/annotations/mapper/templates"
	"go/ast"
	"go/token"
	"strings"

	annotation "github.com/YReshetko/go-annotation/pkg"
)

// @Builder(constructor="newFieldGeneratorBuilder", build="set{{ .FieldName }}", terminator="build", type="pointer")
type fieldGenerator struct {
	node      annotation.Node
	ast       ast.Node
	name      string
	alias     string
	importPkg string
	isPointer bool

	structGen    *structureTypeGenerator // @Exclude
	primitiveGen *primitiveTypeGenerator // @Exclude
}

// @PostConstruct
func (fg *fieldGenerator) buildFields() {
	//fmt.Println("Data:", fg.isPointer, fg.alias, fg.name)
	//ast.Print(token.NewFileSet(), fg.ast)
	if fg.ast == nil {
		return
	}

	switch astType := fg.ast.(type) {
	case *ast.Ident:
		fg.buildByIdentity(astType)
	default:
		fmt.Println("UNSUPPORTED FIELD TYPE")
		ast.Print(token.NewFileSet(), fg.ast)
	}
}

func (fg *fieldGenerator) buildByIdentity(ident *ast.Ident) {
	if ident.Obj != nil {
		// TODO support embedded structures and declared in the same file
		if ident.Obj.Kind != ast.Typ {
			fmt.Println("TODO: UNSUPPORTED OBJECT DECLARATION", ident.Obj.Kind)
			ast.Print(token.NewFileSet(), ident.Obj)
			return
		}
		fg.buildNonPrimitiveType(fg.node, ident.Obj.Decl.(ast.Node))
		return
	}

	if len(fg.alias) == 0 && isPrimitive(ident.String()) {
		fg.primitiveGen = newPrimitiveTypeGeneratorBuilder().setName(ident.String()).build()
		return
	}

	newNode, importPath, err := fg.node.FindNodeByAlias(fg.alias, ident.String())
	if err != nil {
		panic(err)
	}

	fg.importPkg = importPath
	fg.buildNonPrimitiveType(newNode, newNode.ASTNode())
}

func (fg *fieldGenerator) buildNonPrimitiveType(node annotation.Node, astNode ast.Node) {
	switch nnt := astNode.(type) {
	case *ast.TypeSpec:
		name := nnt.Name.String()
		switch innt := nnt.Type.(type) {
		case *ast.StructType:
			var fields []*ast.Field
			if innt.Fields != nil {
				fields = innt.Fields.List
			}
			fg.structGen = newStructureTypeGeneratorBuilder().setName(name).setNode(node).setFields(fields).build()
			return
		default:
			fmt.Printf("UNSUPPORTED INTERNAL LOADED TYPE %T\n", nnt.Type)
			ast.Print(token.NewFileSet(), astNode)
		}
	default:
		fmt.Printf("UNSUPPORTED LOADED TYPE %T\n", astNode)
		ast.Print(token.NewFileSet(), astNode)
	}
}

func (fg *fieldGenerator) argumentNameAndType(c *cache) ([]byte, error) {
	fg.appendImport(c)
	return templates.Execute(templates.ArgumentNameAndTypeTemplate, map[string]interface{}{
		"Name":      fg.name,
		"Type":      fg.buildArgType(),
		"IsPointer": fg.isPointer,
	})
}

func (fg *fieldGenerator) appendImport(c *cache) {
	if len(fg.alias) != 0 {
		c.addImport(Import{
			Alias:  fg.alias,
			Import: fg.importPkg,
		})
	}
}

func (fg *fieldGenerator) buildArgType() string {
	var argType string
	switch {
	case fg.primitiveGen != nil:
		argType = fg.primitiveGen.name
	case fg.structGen != nil:
		argType = fg.structGen.name
	}

	if len(fg.alias) > 0 {
		return strings.Join([]string{fg.alias, argType}, ".")
	}
	return argType
}

func (fg *fieldGenerator) buildReturnType() string {
	tpy := fg.buildArgType()
	if fg.isPointer {
		tpy = "*" + tpy
	}
	return tpy
}

func (fg *fieldGenerator) buildReturnValue() string {
	tpy := fg.buildArgType() + "{}"
	if fg.isPointer {
		tpy = "&" + tpy
	}
	return tpy
}

func (fg *fieldGenerator) generate(name string, in []*fieldGenerator, c *cache, o *overloading) error {
	switch {
	case fg.primitiveGen != nil:
		// TODO generate mapping for primitives
	case fg.structGen != nil:
		return fg.generateStructMapping(name, in, c, o)
	}
	return nil
}

func (fg *fieldGenerator) generateStructMapping(name string, in []*fieldGenerator, c *cache, o *overloading) error {
	varDeclaration, err := templates.Execute(templates.NewVariableTemplate, map[string]interface{}{
		"VariableName": name,
		"IsPointer":    fg.isPointer,
		"Type":         fg.buildArgType(),
	})

	if err != nil {
		return fmt.Errorf("unable to declare field %s: %w", name, err)
	}
	c.addCodeLine(string(varDeclaration))

	for _, field := range fg.structGen.fieldGenerators {
		if !isExported(field) {
			continue
		}
		mapperLine, mt := o.find(name + "." + field.name)
		if mt != none {
			switch mt {
			case source:
				err := fg.generateOverloadedSource(name+"."+field.name, mapperLine, in, field, c)
				if err != nil {
					return fmt.Errorf("unable to generate overloaded source %s: %w", mapperLine, err)
				}
			case function:
				err := fg.generateOverloadedFunction(name+"."+field.name, mapperLine, in, c)
				if err != nil {
					return fmt.Errorf("unable to generate overloaded function %s: %w", mapperLine, err)
				}
			case constant:
				if field.primitiveGen == nil {
					return fmt.Errorf("unable to set constant for non-primitive %s | %s", name+"."+field.name, mapperLine)
				}
				err := mapConstant(name+"."+field.name, field.primitiveGen.name, mapperLine, field.isPointer, c)
				if err != nil {
					return fmt.Errorf("unable to build constant mapping %s: %w", name+"."+field.name, err)
				}
			}
			continue
		}

		if o.isIgnoreDefault {
			continue
		}

		potentialFields := findFieldsByName(field.name, in)
		if len(potentialFields) == 0 {
			continue
		}

		if len(potentialFields) > 1 {
			// TODO support selection by appropriate type
		}

		var fromPrefix []string
		if potentialFields[0].root.isPointer {
			fromPrefix = append(fromPrefix, potentialFields[0].root.name)
		}
		if isBothPrimitives(field, potentialFields[0].field) {
			err := mapPrimitives(name+"."+field.name, potentialFields[0].name, field, potentialFields[0].field, fromPrefix, c)
			if err != nil {
				return fmt.Errorf("unable to build primitive mapping: %w", err)
			}
			continue
		}

		if isBothStructures(field, potentialFields[0].field) {
			err := mapStructures(name+"."+field.name, potentialFields[0].name, field, potentialFields[0].field, fromPrefix, c)
			if err != nil {
				return fmt.Errorf("unable to build structures mapping: %w", err)
			}
			continue
		}

		// TODO Slice mapping, map mapping should be processed here
		// Other cases primitives -> struct, map, slice should be overridden above
	}

	return nil
}

func (fg *fieldGenerator) generateOverloadedSource(toName, fromLine string, in []*fieldGenerator, to *fieldGenerator, c *cache) error {
	names := strings.Split(fromLine, ".")
	buff, isFromPointer, err := fg.findPointersInPath("", names, in, []string{})
	if err != nil {
		return fmt.Errorf("unable to build mapping for %s: %w", toName, err)
	}

	assignLine := toName + " = " + fromLine
	if to.isPointer != isFromPointer {
		if to.isPointer {
			assignLine = toName + " = &" + fromLine
		} else {
			assignLine = toName + " = *" + fromLine
		}
	}

	if len(buff) > 0 {
		c.addIfClause(buff, assignLine)
		return nil
	}
	c.addCodeLine(assignLine)
	return nil
}

func (fg *fieldGenerator) generateOverloadedFunction(toName, mappingLine string, in []*fieldGenerator, c *cache) error {
	if strings.Index(mappingLine, ")") < strings.Index(mappingLine, ")") {
		return fmt.Errorf("invalid function call %s", mappingLine)
	}

	args := strings.Split(mappingLine[strings.Index(mappingLine, "(")+1:strings.Index(mappingLine, ")")], ",")
	splitArgs := make([][]string, len(args))
	for i, arg := range args {
		splitArgs[i] = strings.Split(strings.TrimSpace(arg), ".")
	}
	var nilableLines []string
	for _, arg := range splitArgs {
		buff, _, err := fg.findPointersInPath("", arg, in, []string{})
		if err != nil {
			return fmt.Errorf("unable to build mapping for %s=%s: %w", toName, mappingLine, err)
		}
		if len(buff) > 0 {
			nilableLines = append(nilableLines, buff...)
		}
	}

	assignLine := toName + " = " + mappingLine
	if len(nilableLines) > 0 {
		var ls []string
		distinct := map[string]struct{}{}
		for _, line := range nilableLines {
			if _, ok := distinct[line]; !ok {
				ls = append(ls, line)
				distinct[line] = struct{}{}
			}
		}
		c.addIfClause(ls, assignLine)
		return nil
	}
	c.addCodeLine(assignLine)
	return nil
}

func (fg *fieldGenerator) findPointersInPath(prefix string, rest []string, in []*fieldGenerator, buffer []string) ([]string, bool, error) {
	name := rest[0]
	var fromFieldGen *fieldGenerator
	for _, generator := range in {
		if generator.name == name {
			fromFieldGen = generator
			break
		}
	}
	if fromFieldGen == nil {
		return nil, false, fmt.Errorf("unable to find field %s.%s", prefix, name)
	}

	if fromFieldGen.primitiveGen != nil && len(rest) > 1 {
		return nil, false, fmt.Errorf("dot notation is not allowed for pimitive %s.%s", prefix, name)
	}
	if len(prefix) == 0 {
		prefix = name
	} else {
		prefix = prefix + "." + name
	}

	if fromFieldGen.isPointer {
		buffer = append(buffer, prefix)
	}
	if len(rest) == 1 {
		return buffer, fromFieldGen.isPointer, nil
	}
	var newIn []*fieldGenerator
	if fromFieldGen.structGen != nil {
		newIn = fromFieldGen.structGen.fieldGenerators
	}
	return fg.findPointersInPath(prefix, rest[1:], newIn, buffer)
}

type nameFieldPair struct {
	name  string
	field *fieldGenerator
	root  *fieldGenerator
}

func findFieldsByName(name string, in []*fieldGenerator) []nameFieldPair {
	var out []nameFieldPair
	for _, generator := range in {
		if generator.primitiveGen != nil {
			// TODO support primitive mapping via annotations
			continue
		}
		if generator.structGen == nil {
			continue
		}

		for _, fg := range generator.structGen.fieldGenerators {
			if !isExported(fg) {
				continue
			}
			if fg.name == name {
				out = append(out, nameFieldPair{
					name:  generator.name + "." + fg.name,
					field: fg,
					root:  generator,
				})
			}
		}
	}
	return out
}

func buildFiledGenerator(f *ast.Field, name string, n annotation.Node) *fieldGenerator {
	field, isPointer := unwrapStarExpression(f.Type)
	internalType, alias := unwrapSelectorExpression(field)

	return newFieldGeneratorBuilder().
		setName(name).
		setNode(n).
		setAlias(alias).
		setAst(internalType).
		setIsPointer(isPointer).
		build()
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

func isExported(f *fieldGenerator) bool {
	if len(f.name) == 0 {
		return false
	}
	return f.name[0] >= 'A' && f.name[0] <= 'Z'
}
