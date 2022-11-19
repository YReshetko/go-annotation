package mock

import (
	"bytes"
	"fmt"
	"go/ast"
	"path/filepath"
	"text/template"

	"github.com/maxbrunsfeld/counterfeiter/v6/generator"

	annotation "github.com/YReshetko/go-annotation/pkg"
)

func init() {
	annotation.Register[Mock](&Processor{
		out:     map[string][]byte{},
		cache:   &generator.Cache{},
		toRerun: []string{},
	})
}

var _ annotation.AnnotationProcessor = (*Processor)(nil)

type Processor struct {
	out     map[string][]byte
	cache   generator.Cacher
	toRerun []string
}

func (p *Processor) Process(node annotation.Node) error {
	annotations := annotation.FindAnnotations[Mock](node.Annotations())
	if len(annotations) == 0 {
		return nil
	}

	if len(annotations) > 1 {
		return fmt.Errorf("expected 1 mock annotation, but got: %d", len(annotations))
	}

	a := annotations[0]
	typeName, mod, err := extractTypeName(node)
	if err != nil {
		return err
	}

	mockName, err := createMockInterfaceName(a.Name, typeName)
	if err != nil {
		return err
	}

	f, err := generator.NewFake(
		mod,
		typeName,
		node.Meta().Dir(),
		mockName,
		a.SubPackage,
		"",
		node.Meta().Root(),
		p.cache,
	)
	if err != nil {
		return fmt.Errorf("unable to prepare generator for %s: %w", typeName, err)
	}

	data, err := f.Generate(true)
	if err != nil {
		return fmt.Errorf("unable to generate mock for %s: %w", typeName, err)
	}

	outputFile := filepath.Join(node.Meta().Dir(), a.SubPackage, toSnakeCase(mockName)+".gen.go")
	if _, ok := p.out[outputFile]; ok {
		return fmt.Errorf("attemption to save same file with different data twice: %s", outputFile)
	}

	if mod == generator.Package {
		data, err = enrichMockAnnotation(data)
		if err != nil {
			return fmt.Errorf("unable to ingest generated code for %s: %w", mockName, err)
		}
		toFile, err := filepath.Rel(node.Meta().Root(), outputFile)
		if err != nil {
			return fmt.Errorf("unable to build relativ path %s: %w", outputFile, err)
		}
		p.toRerun = append(p.toRerun, toFile)
	}

	p.out[outputFile] = data

	return nil
}

func (p *Processor) Output() map[string][]byte {
	return p.out
}

func (p *Processor) Version() string {
	return "0.0.1"
}

func (p *Processor) Name() string {
	return "Mock"
}

func (p *Processor) ToRerun() []string {
	return p.toRerun
}

func (p *Processor) Clear() {
	p.toRerun = []string{}
	p.out = map[string][]byte{}
}
func extractTypeName(node annotation.Node) (string, generator.FakeMode, error) {
	var nameIdent *ast.Ident
	mod := generator.InterfaceOrFunction
	switch n := node.ASTNode().(type) {
	case *ast.TypeSpec:
		nameIdent = n.Name
		switch n.Type.(type) {
		case *ast.InterfaceType, *ast.FuncType, *ast.IndexListExpr, *ast.IndexExpr, *ast.Ident:
		default:
			return "", mod, fmt.Errorf("expected mocked type one of [*ast.InterfaceType, *ast.FuncType, *ast.IndexListExpr, *ast.IndexExpr, *ast.Ident], but got %T for %s", n.Type, nameIdent.String())
		}

	case *ast.FuncDecl:
		nameIdent = n.Name
	case *ast.File:
		nameIdent = n.Name
		mod = generator.Package
	default:
		return "", mod, fmt.Errorf("expected mocked type one of [ast.TypeSpec, *ast.TypeDecl, *ast.File], but got %T", node.ASTNode())
	}
	if nameIdent.String() == "" {
		return "", mod, fmt.Errorf("unable to prepare mock for interface in %s", node.Meta().Dir())
	}
	return nameIdent.String(), mod, nil
}

func createMockInterfaceName(nameTemplate, interfaceName string) (string, error) {
	tmpl, err := template.New("mock_name").Parse(nameTemplate)
	if err != nil {
		return "", fmt.Errorf("unable to parse Mock annotation name tmplate %s: %w", interfaceName, err)
	}
	buf := bytes.NewBufferString("")
	if err = tmpl.Execute(buf, map[string]string{"TypeName": interfaceName}); err != nil {
		return "", fmt.Errorf("unable to prepare mock name %s: %w", interfaceName, err)
	}
	return buf.String(), nil
}

func toSnakeCase(s string) string {
	size := 'a' - 'A'
	out := ""
	for i, ch := range s {
		if ch >= 'A' && ch <= 'Z' {
			if i == 0 {
				out += string(ch + size)
			} else {
				out += "_" + string(ch+size)
			}
		} else {
			out += string(ch)
		}

	}
	return out
}
