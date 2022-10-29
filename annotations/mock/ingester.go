package mock

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/printer"

	"github.com/YReshetko/go-annotation/internal/utils/astutils"
)

func ingestMockAnnotation(data []byte) ([]byte, error) {
	fileSpec, fset, err := astutils.BytesToAST(data)
	if err != nil {
		return nil, fmt.Errorf("unable to convert bytes to AST: %w", err)
	}

	ast.Inspect(fileSpec, func(node ast.Node) bool {
		gd, ok := node.(*ast.GenDecl)
		if !ok {
			return true
		}

		ts, ok := gd.Specs[0].(*ast.TypeSpec)
		if !ok {
			return true
		}

		_, ok = ts.Type.(*ast.InterfaceType)
		if !ok {
			return true
		}

		if gd.Doc != nil {
			gd.Doc.List = append(gd.Doc.List, &ast.Comment{
				Text: "// Ingested annotations: @Mock(sub=\"mocksfakes\")",
			})
		} else {
			gd.Doc = &ast.CommentGroup{
				List: []*ast.Comment{{
					Text: "Ingested annotations: @Mock(sub=\"mocksfakes\")",
				}},
			}
		}
		return true
	})

	b := bytes.NewBufferString("")
	err = printer.Fprint(b, fset, fileSpec)
	if err != nil {
		return nil, fmt.Errorf("unable to prepare ingested AST: %w", err)
	}

	return b.Bytes(), nil
}
