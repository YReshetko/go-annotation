package generators

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/printer"
	"go/token"
	"strings"

	"github.com/YReshetko/go-annotation/annotations/constructor/annotations"
	annotation "github.com/YReshetko/go-annotation/pkg"
)

type arguments struct {
	incoming map[string]string
	toInit   map[string]string
}

func extractArguments(n *ast.TypeSpec, lookup importLookup, node annotation.Node) (arguments, distinctImports) {
	out := arguments{
		incoming: map[string]string{},
		toInit:   map[string]string{},
	}
	imps := newDistinctImports()

	strTpy, ok := n.Type.(*ast.StructType)
	if !ok {
		panic("not a *ast.StructType")
	}

	if strTpy.Fields == nil {
		fmt.Println("no TypeParams")
		return out, imps
	}
	fields := strTpy.Fields.List
	if len(fields) == 0 {
		fmt.Println("no TypeParams.List")
		return out, imps
	}

	for _, field := range fields {
		toCheck := node.AnnotatedNode(field)
		if len(annotation.FindAnnotations[annotations.Exclude](toCheck.Annotations())) > 0 {
			continue
		}
		buff := bytes.NewBufferString("")
		err := printer.Fprint(buff, &token.FileSet{}, field.Type)
		if err != nil {
			panic(err)
		}
		initAnnotations := annotation.FindAnnotations[annotations.Init](toCheck.Annotations())
		if len(initAnnotations) == 0 {
			for _, ident := range field.Names {
				out.incoming[ident.Name] = buff.String()
			}
		} else {
			result := ""
			origVal := buff.String()
			switch field.Type.(type) {
			case *ast.ArrayType:
				result = sliceInitialisation(initAnnotations[0], origVal)
			case *ast.MapType:
				result = mapInitialisation(initAnnotations[0], origVal)
			case *ast.ChanType:
				result = chanInitialisation(initAnnotations[0], origVal)
			default:
				continue
			}
			for _, ident := range field.Names {
				out.toInit[ident.Name] = result
				out.incoming[ident.Name] = origVal
			}
		}

		imps.merge(getImports(field.Type, lookup))

	}
	return out, imps
}

func sliceInitialisation(a annotations.Init, t string) string {
	switch {
	case a.Cap > -1 && a.Len > -1:
		return fmt.Sprintf("make(%s, %d, %d)", t, a.Len, a.Cap)
	case a.Len > -1:
		return fmt.Sprintf("make(%s, %d)", t, a.Len)
	case a.Cap > -1:
		return fmt.Sprintf("make(%s, 0, %d)", t, a.Cap)
	default:
		return fmt.Sprintf("%s{}", t)
	}
}

func mapInitialisation(a annotations.Init, t string) string {
	switch {
	case a.Cap > -1:
		return fmt.Sprintf("make(%s, %d)", t, a.Cap)
	default:
		return fmt.Sprintf("make(%s)", t)
	}
}

func chanInitialisation(a annotations.Init, t string) string {
	switch {
	case a.Cap > -1:
		return fmt.Sprintf("make(%s, %d)", t, a.Cap)
	default:
		return fmt.Sprintf("make(%s)", t)
	}
}

type parameters struct {
	constraints    string
	parameters     string
	isParametrised bool
}

func extractParameters(n *ast.TypeSpec, lookup importLookup) (parameters, distinctImports) {
	if n.TypeParams == nil || len(n.TypeParams.List) == 0 {
		return parameters{}, nil
	}
	var c []string
	var p []string
	imps := newDistinctImports()
	for _, field := range n.TypeParams.List {
		buff := bytes.NewBufferString("")
		err := printer.Fprint(buff, &token.FileSet{}, field.Type)
		if err != nil {
			panic(err)
		}
		for _, name := range field.Names {
			if name == nil || len(name.Name) == 0 {
				continue
			}
			c = append(c, name.Name+" "+buff.String())
			p = append(p, name.Name)
		}
		imps.merge(getImports(field.Type, lookup))
	}

	return parameters{
		constraints:    strings.Join(c, ","),
		parameters:     strings.Join(p, ","),
		isParametrised: true,
	}, imps
}
