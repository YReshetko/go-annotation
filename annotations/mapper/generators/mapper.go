package generators

import (
	"fmt"
	"github.com/YReshetko/go-annotation/annotations/mapper/annotations"
	cache2 "github.com/YReshetko/go-annotation/annotations/mapper/generators/cache"
	"github.com/YReshetko/go-annotation/annotations/mapper/templates"
	annotation "github.com/YReshetko/go-annotation/pkg"
	"go/ast"
)

// @Builder(type="pointer")
type MapperGenerator struct {
	impCache   *cache2.ImportCache
	node       annotation.Node
	intType    *ast.InterfaceType
	intName    string
	structName string

	mgs []*methodGenerator // @Init
}

// @PostConstruct
func (mg *MapperGenerator) buildMethodGenerators() {
	if mg.intType == nil {
		return
	}
	methods := mg.intType.Methods
	if methods == nil {
		return
	}

	for _, method := range mg.intType.Methods.List {
		ft, ok := method.Type.(*ast.FuncType)
		if !ok {
			continue
		}

		o, err := getOverloading(mg.node.AnnotatedNode(method))
		if err != nil {
			panic(err)
		}

		methodGen := newMethodGeneratorBuilder().
			setImpCache(mg.impCache).
			setNode(mg.node).
			setName(methodName(method.Names)).
			setInput(fieldsFromFiledList(ft.Params)).
			setOutput(fieldsFromFiledList(ft.Results)).
			setOverloading(o).
			build()

		mg.mgs = append(mg.mgs, methodGen)
	}
}

func (mg *MapperGenerator) Generate() ([]byte, []Import, error) {
	var methods []string
	imports := make(map[Import]struct{})
	for _, generator := range mg.mgs {
		m, err := generator.generate(mg.structName)
		if err != nil {
			return nil, nil, fmt.Errorf("unable to generte method %s: %w", generator.name, err)
		}
		methods = append(methods, string(m))
	}

	data, err := templates.Execute(templates.MapperStructureTemplate, map[string]interface{}{
		"TypeName":      mg.structName,
		"InterfaceName": mg.intName,
		"Methods":       methods,
	})

	return data, getImports(imports), err
}

func fieldsFromFiledList(fl *ast.FieldList) []*ast.Field {
	if fl == nil {
		return nil
	}
	return fl.List
}

func methodName(names []*ast.Ident) string {
	if len(names) == 0 {
		return ""
	}
	return names[0].String()
}

func getImports(imports map[Import]struct{}) []Import {
	imps := make([]Import, len(imports))
	var ind int
	for i, _ := range imports {
		imps[ind] = i
		ind++
	}
	return imps
}

func getOverloading(node annotation.Node) (*overloading, error) {
	idm := annotation.FindAnnotations[annotations.IgnoreDefaultMapping](node.Annotations())
	o := newOverloading(len(idm) > 0)
	for _, a := range annotation.FindAnnotations[annotations.Mapping](node.Annotations()) {
		err := o.Add(a.Target, a.Source, a.This, a.Func, a.Constant)
		if err != nil {
			return nil, fmt.Errorf("unable to prepare overloading config: %w", err)
		}
	}
	for _, sm := range annotation.FindAnnotations[annotations.SliceMapping](node.Annotations()) {
		err := o.AddSlice(sm.Target, sm.Source, sm.This, sm.Func)
		if err != nil {
			return nil, fmt.Errorf("unable to prepare overloading config: %w", err)
		}
	}
	for _, sm := range annotation.FindAnnotations[annotations.MapMapping](node.Annotations()) {
		err := o.AddMap(sm.Target, sm.Source, sm.This, sm.Func)
		if err != nil {
			return nil, fmt.Errorf("unable to prepare overloading config: %w", err)
		}
	}

	return o, nil
}
