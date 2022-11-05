package generators

import (
	"fmt"
	"github.com/YReshetko/go-annotation/annotations/constructor/annotations"
	annotation "github.com/YReshetko/go-annotation/pkg"
	"go/ast"
	"sort"
)

type PostConstructValues struct {
	Annotation annotations.PostConstruct
	MethodName string
}

func PostConstructReceiverName(node annotation.Node) (string, PostConstructValues, error) {
	ans := annotation.FindAnnotations[annotations.PostConstruct](node.Annotations())
	if len(ans) == 0 {
		return "", PostConstructValues{}, nil
	}

	if len(ans) > 1 {
		return "", PostConstructValues{}, fmt.Errorf("expected 1 PostConstruct annotation, but got: %d", len(ans))
	}

	fd, ok := annotation.CastNode[*ast.FuncDecl](node)
	if !ok {
		return "", PostConstructValues{}, fmt.Errorf("expected ast.FuncDecl for PostConstruct annotation, but got %t", node.Node())
	}

	if fd.Recv == nil {
		return "", PostConstructValues{}, fmt.Errorf("PostConstruct should be declared for methods only: %s", fd.Name.Name)
	}

	if fd.Type != nil && fd.Type.Params != nil && len(fd.Type.Params.List) > 0 {
		return "", PostConstructValues{}, fmt.Errorf("expected PostConstruct method should not have any arguments for %s", fd.Name.Name)
	}

	receiverName := ""
	switch rcv := fd.Recv.List[0].Type.(type) {
	case *ast.StarExpr:
		switch ile := rcv.X.(type) {
		case *ast.IndexListExpr:
			receiverName = ile.X.(*ast.Ident).Name
		case *ast.Ident:
			receiverName = ile.Name
		default:
			return "", PostConstructValues{}, fmt.Errorf("invalid method declaration for PostConstruct: %s, expected *ast.IndexListExpr, but got %T", fd.Name.Name, rcv.X)
		}
	case *ast.IndexListExpr:
		receiverName = rcv.X.(*ast.Ident).Name
	case *ast.Ident:
		receiverName = rcv.Name
	default:
		return "", PostConstructValues{}, fmt.Errorf("unexpected node for PostConstruct %T", fd.Recv.List[0].Type)
	}
	return receiverName, PostConstructValues{
		Annotation: ans[0],
		MethodName: fd.Name.Name,
	}, nil
}

func sortPostConstructs(pcvs []PostConstructValues) []PostConstructValues {
	if len(pcvs) == 0 {
		return pcvs
	}
	sort.Slice(pcvs, func(i, j int) bool {
		return pcvs[i].Annotation.Priority < pcvs[j].Annotation.Priority
	})
	return pcvs
}

func postConstructMethods(pcvs []PostConstructValues) []string {
	pcvs = sortPostConstructs(pcvs)
	if len(pcvs) == 0 {
		return []string{}
	}
	out := make([]string, len(pcvs))
	for i, pcv := range pcvs {
		out[i] = pcv.MethodName
	}
	return out
}
