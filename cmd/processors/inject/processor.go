package inject

import (
	"fmt"
	"github.com/YReshetko/go-annotation/pkg"
)

type Processor struct{}

func init() {
	pkg.Register(Inject{}, &Processor{})
}

func (p *Processor) Process(annotation pkg.Annotation, node pkg.Node) error {
	a := pkg.CastAnnotation[Inject](annotation)
	fmt.Println(a, node)
	return nil
}
func (p *Processor) Output() map[pkg.Path]pkg.Data {
	return nil
}
