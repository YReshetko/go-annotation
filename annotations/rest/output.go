package rest

import (
	"bytes"
	"fmt"
	"net/http"
	"text/template"
)

const (
	handlerPackage = `package %s
`
	handlerImports = `
import (
	"net/http"
)
`
	handlerTemplate = `
func (s *{{ .Receiver }}) Handlers() map[string]http.HandlerFunc {
	return map[string]http.HandlerFunc{
		{{range .Handlers}}{{ .Method }}: s.{{ .Handler }},
		{{end}}}
}

func (s *{{ .Receiver }}) Path() string {
	return "{{ .Path }}"
}
`
)

type outputKey struct {
	pkg string
	dir string
}

type handlerTemplateMethods struct {
	Method  string
	Handler string
}
type handlerTemplateParams struct {
	Receiver string
	Path     string
	Handlers []handlerTemplateMethods
}

type output struct {
	files map[outputKey]*bytes.Buffer
	tmp   *template.Template
}

func newOutput() (*output, error) {
	tmp, err := template.New("handler").Parse(handlerTemplate)
	if err != nil {
		return nil, err
	}
	return &output{
		files: make(map[outputKey]*bytes.Buffer),
		tmp:   tmp,
	}, nil

}

func (o *output) append(meta handlerMetadata, handlers handlerMapping) error {
	outKey := outputKey{
		pkg: meta.pkg,
		dir: meta.dir,
	}
	b := o.upsertBuffer(outKey)

	htm := []handlerTemplateMethods{}
	for m, h := range handlers.mapping {
		htm = append(htm, handlerTemplateMethods{Method: httpMethodConstant(m), Handler: h})
	}
	toTemplate := handlerTemplateParams{
		Receiver: meta.structName,
		Path:     handlers.path,
		Handlers: htm,
	}

	err := o.tmp.Execute(b, toTemplate)
	if err != nil {
		return err
	}
	return nil
}

func (o *output) get() map[string][]byte {
	out := map[string][]byte{}
	for k, buf := range o.files {
		out[k.dir+"/"+"rest.handlers.gen.go"] = buf.Bytes()
	}
	return out
}

func (o *output) upsertBuffer(k outputKey) *bytes.Buffer {
	if b, ok := o.files[k]; ok {
		return b
	}
	b := bytes.NewBufferString(fmt.Sprintf(handlerPackage, k.pkg) + handlerImports + "\n")
	o.files[k] = b
	return b
}

func httpMethodConstant(m string) string {
	c, _ := map[string]string{
		http.MethodGet:     "http.MethodGet",
		http.MethodHead:    "http.MethodHead",
		http.MethodPost:    "http.MethodPost",
		http.MethodPut:     "http.MethodPut",
		http.MethodPatch:   "http.MethodPatch",
		http.MethodDelete:  "http.MethodDelete",
		http.MethodConnect: "http.MethodConnect",
		http.MethodOptions: "http.MethodOptions",
		http.MethodTrace:   "http.MethodTrace",
	}[m]
	return c

}
