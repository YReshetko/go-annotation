package nodes

import (
	"go/ast"
	"go/parser"
	"go/token"
	"io/fs"
	"path/filepath"
	"strings"

	"github.com/YReshetko/go-annotation/internal/annotation"
)

func ReadProject(path string) ([]Node, error) {
	nodes := []Node{}
	err := filepath.Walk(path, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if isGoFile(info) {
			fset := token.NewFileSet()
			fileSpec, err := parser.ParseFile(fset, path, nil, parser.ParseComments)
			if err != nil {
				return err
			}

			ast.Inspect(fileSpec, func(node ast.Node) bool {
				processedNodes, proceed := processNode(node)
				for _, n := range processedNodes {
					n.Metadata.FileSpec = fileSpec
					n.Metadata.Dir = dirByPath(path, info.Name())
					n.Metadata.FileName = info.Name()
					nodes = append(nodes, n)
				}

				return proceed
			})

		}
		return nil
	})

	return nodes, err
}

func dirByPath(fullFilePath, fileName string) string {
	dir := strings.TrimRight(fullFilePath, fileName)
	if dir[len(dir)-1] == '/' {
		return dir[:len(dir)-1]
	}
	return dir

}

func isGoFile(info fs.FileInfo) bool {
	return !info.IsDir() && info.Name()[len(info.Name())-2:] == "go"
}

func processNode(node ast.Node) ([]Node, bool) {
	switch v := node.(type) {
	case *ast.FuncDecl:
		a, ok := annotation.Parse(v.Doc.Text())
		if !ok {
			return nil, false
		}
		return []Node{
			{
				Metadata: Metadata{
					Name: v.Name.Name,
					Type: Function,
				},
				Annotations: a,
				GoNode:      v,
			},
		}, false

	case *ast.GenDecl:
		nodes := []Node{}
		for _, spec := range v.Specs {
			node := processSpec(v, spec)
			if node != nil {
				nodes = append(nodes, *node)
			}
		}
		return nodes, false

	default:
		//fmt.Printf("%T, %+v\n", v, v)
	}
	return nil, true
}

func processSpec(n *ast.GenDecl, spec ast.Spec) *Node {
	switch v := spec.(type) {
	case *ast.TypeSpec:
		switch t := v.Type.(type) {
		case *ast.StructType:
			a, ok := annotation.Parse(n.Doc.Text())
			fields := processFields(t.Fields)
			if len(fields) == 0 && !ok {
				return nil
			}
			return &Node{
				Metadata: Metadata{
					Name: v.Name.Name,
					Type: Structure,
				},
				Annotations: a,
				GoNode:      v,
				Inner:       fields,
			}

		case *ast.InterfaceType:
			a, ok := annotation.Parse(n.Doc.Text())
			if !ok {
				return nil
			}
			return &Node{
				Metadata: Metadata{
					Name: v.Name.Name,
					Type: Interface,
				},
				Annotations: a,
				GoNode:      v,
			}
		}
	case *ast.ValueSpec:
		doc := n.Doc.Text()
		if len(doc) == 0 {
			doc = v.Doc.Text()
		}
		a, ok := annotation.Parse(doc)
		if !ok {
			return nil
		}
		return &Node{
			Metadata: Metadata{
				Name: v.Names[0].Name,
				Type: Variable,
			},
			Annotations: a,
			GoNode:      v,
		}
	}
	return nil
}

func processFields(params *ast.FieldList) []Node {
	if params == nil {
		return nil
	}

	if len(params.List) == 0 {
		return nil
	}

	nodes := []Node{}
	for _, field := range params.List {
		a, ok := annotation.Parse(field.Doc.Text())
		if !ok {
			continue
		}
		nodes = append(nodes, Node{
			Metadata: Metadata{
				Name: field.Names[0].Name,
				Type: Field,
			},
			Annotations: a,
			GoNode:      field,
		})
	}
	return nodes
}
