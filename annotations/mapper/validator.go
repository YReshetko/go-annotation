package mapper

import (
	"fmt"
	"github.com/YReshetko/go-annotation/annotations/mapper/annotations"
	annotation "github.com/YReshetko/go-annotation/pkg"
	"go/ast"
)

func validateAndGetMapperWithTypeSpec(n annotation.Node) (annotations.Mapper, *ast.TypeSpec, error) {
	ans := annotation.FindAnnotations[annotations.Mapper](n.Annotations())
	if len(ans) == 0 {
		return annotations.Mapper{}, nil, nil
	}

	if len(ans) > 1 {
		return annotations.Mapper{}, nil, fmt.Errorf("expected Mapper annotation number is 1, but got %d", len(ans))
	}

	ts, ok := n.ASTNode().(*ast.TypeSpec)
	if !ok {
		return annotations.Mapper{}, nil, fmt.Errorf("expected Mapper is declared on *ast.TypeSpec, but got %T", n.ASTNode())
	}

	if ts.TypeParams != nil {
		return annotations.Mapper{}, nil, fmt.Errorf("parametrized mappers are not supported: %s", ts.Name.String())
	}

	_, ok = ts.Type.(*ast.InterfaceType)
	if !ok {
		return annotations.Mapper{}, nil, fmt.Errorf("expacted Mapper annotation on *ast.InterfaceType, but got : %T", ts.Type)
	}

	return ans[0], ts, nil
}
