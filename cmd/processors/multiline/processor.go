package endpoint

import (
	"fmt"
	"github.com/YReshetko/go-annotation/pkg"
)

type Processor struct{}

func init() {
	pkg.Register(Multiline{}, &Processor{})
}

func (p *Processor) Process(annotation pkg.Annotation, node pkg.Node) error {
	a := pkg.CastAnnotation[Multiline](annotation)
	fmt.Println(a, node)
	return nil
}
func (p *Processor) Output() map[pkg.Path]pkg.Data {
	return nil
}
