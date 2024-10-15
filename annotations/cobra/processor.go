package cobra

import (
	"errors"
	"fmt"
	"go/ast"
	"path/filepath"
	"regexp"
	"slices"
	"strings"

	"github.com/YReshetko/go-annotation/annotations/cobra/annotations"
	"github.com/YReshetko/go-annotation/annotations/cobra/cache"
	"github.com/YReshetko/go-annotation/annotations/cobra/templates"
	annotation "github.com/YReshetko/go-annotation/pkg"
)

func init() {
	p := &Processor{cache: cache.NewCache()}
	annotation.Register[annotations.CobraOutput](p)
	annotation.Register[annotations.Cobra](p)
	annotation.Register[annotations.CobraPersistPreRun](p)
	annotation.Register[annotations.CobraPreRun](p)
	annotation.Register[annotations.CobraRun](p)
	annotation.Register[annotations.CobraPostRun](p)
	annotation.Register[annotations.CobraPersistPostRun](p)
}

var _ annotation.AnnotationProcessor = (*Processor)(nil)

type Processor struct {
	outputRoot string
	cache      *cache.Cache
}

func (p *Processor) Process(node annotation.Node) error {
	cobraOutputs := annotation.FindAnnotations[annotations.CobraOutput](node.Annotations())
	if len(cobraOutputs) != 0 {
		p.outputRoot = node.Meta().Dir()
	}

	return errors.Join(
		addMethod[annotations.CobraPersistPreRun](node, p.cache),
		addMethod[annotations.CobraPreRun](node, p.cache),
		addMethod[annotations.CobraRun](node, p.cache),
		addMethod[annotations.CobraPostRun](node, p.cache),
		addMethod[annotations.CobraPersistPostRun](node, p.cache),
		addCobraCommand(node, p.cache),
	)
}

type methodAnnotation interface {
	IsPersistRun() bool
	IsPreRun() bool
	IsPostRun() bool
}

func addMethod[T methodAnnotation](node annotation.Node, cache *cache.Cache) error {
	runs := annotation.FindAnnotations[T](node.Annotations())
	if len(runs) == 0 {
		return nil
	}
	if len(runs) > 1 {
		return fmt.Errorf("ambigouse @%T marker", runs[0])
	}
	funcDecl, ok := annotation.CastNode[*ast.FuncDecl](node)
	if !ok || funcDecl.Recv == nil || len(funcDecl.Recv.List) == 0 {
		return fmt.Errorf("@%T can be used withing methods only", runs[0])
	}
	funcReceiverIdent, ok := funcDecl.Recv.List[0].Type.(*ast.Ident)
	if !ok {
		return fmt.Errorf("expected type identity as a receiver, but got %T", funcDecl.Recv.List[0].Type)
	}

	hasReturn := false
	if funcDecl.Type.Results != nil && len(funcDecl.Type.Results.List) > 0 {
		if len(funcDecl.Type.Results.List) > 1 {
			return fmt.Errorf("run function should have no more than one value, but got %d", len(funcDecl.Type.Results.List))
		}
		retTypeIdent, ok := funcDecl.Type.Results.List[0].Type.(*ast.Ident)
		if !ok {
			return fmt.Errorf("return type should be *ast.Ident, but got %T", funcDecl.Type.Results.List[0].Type)
		}

		if retTypeIdent.Name != "error" {
			return fmt.Errorf("expected return type is 'error', but got '%s'", retTypeIdent.Name)
		}
		hasReturn = true
	}

	handler := templates.Handler{
		MethodName:           funcDecl.Name.Name,
		ExecutorPackageAlias: node.Meta().PackageName(),
		ExecutorTypeName:     funcReceiverIdent.Name,
		IsPreRun:             runs[0].IsPreRun(),
		IsPostRun:            runs[0].IsPostRun(),
		IsPersistentRun:      runs[0].IsPersistRun(),
		HasReturn:            hasReturn,
	}
	cache.AddHandler(node.Meta().LocalPackage(), funcReceiverIdent.Name, handler)

	return nil
}

func addCobraCommand(node annotation.Node, cache *cache.Cache) error {
	cobra := annotation.FindAnnotations[annotations.Cobra](node.Annotations())
	if len(cobra) == 0 {
		return nil
	}
	typeSpec, ok := annotation.CastNode[*ast.TypeSpec](node)
	if !ok {
		return fmt.Errorf("expected @Cobra annotation on *ast.TypeSpec, but got %T", node.ASTNode())
	}
	structType, ok := typeSpec.Type.(*ast.StructType)
	if !ok {
		return fmt.Errorf("expected @Cobra annotation on *ast.StructType, but got %T", typeSpec.Type)
	}
	if typeSpec.Name == nil {
		return fmt.Errorf("expected @Cobar annotation on named structure type")
	}
	cache.AddCommandAnnotation(node.Meta().LocalPackage(), typeSpec.Name.Name, cobra[0])

	if structType.Fields == nil {
		return nil
	}
	for _, field := range structType.Fields.List {
		if field.Tag == nil || len(field.Tag.Value) == 0 {
			continue
		}
		flag, ok := getFlag(field)
		if !ok {
			continue
		}
		cache.AddFlag(node.Meta().LocalPackage(), typeSpec.Name.Name, flag)
	}

	//ast.Print(token.NewFileSet(), structType)
	return nil
}

func cutFlagTags(tags string) []string {
	exps := []*regexp.Regexp{
		regexp.MustCompile(`flag:".*?"`),
		regexp.MustCompile(`short:".*?"`),
		regexp.MustCompile(`default:".*?"`),
		regexp.MustCompile(`description:".*?"`),
	}
	out := make([]string, 0, 4)
	for _, exp := range exps {
		s := exp.FindString(tags)
		if len(s) > 0 {
			out = append(out, s)
		}
	}
	return out
}

func getFlag(field *ast.Field) (templates.Flag, bool) {
	tagValues := map[string]string{}
	for _, value := range cutFlagTags(strings.Trim(field.Tag.Value, "`")) {
		kv := strings.Split(value, ":")
		if len(kv) != 2 {
			continue
		}
		tagValues[kv[0]] = strings.Trim(kv[1], `"`)
	}

	flagNameValue, ok := tagValues["flag"]
	if !ok {
		return templates.Flag{}, false
	}

	flagNameValues := strings.Split(flagNameValue, ",")
	for i, value := range flagNameValues {
		flagNameValues[i] = strings.TrimSpace(value)
	}
	if len(flagNameValues) == 0 {
		return templates.Flag{}, false
	}

	if slices.Contains(flagNameValues, "inherited") {
		return templates.Flag{}, false
	}

	fieldTypeName := "string"
	fieldTypeIdent, ok := field.Type.(*ast.Ident)
	if ok {
		fieldTypeName = fieldTypeIdent.Name
	}

	flagType := templates.FlagTypeByASTPrimitive(fieldTypeName)
	defaultValue := tagValues["default"]
	if len(defaultValue) == 0 {
		defaultValue = templates.FlagTypeDefaultValue(flagType)
	}

	return templates.Flag{
		Type:         flagType,
		Name:         flagNameValues[0],
		Shorthand:    tagValues["short"],
		DefaultValue: defaultValue,
		Description:  tagValues["description"],
		IsRequired:   slices.Contains(flagNameValues, "required"),
		IsPersistent: slices.Contains(flagNameValues, "persist"),
	}, true
}

func (p *Processor) Output() map[string][]byte {
	initCommands, err := p.cache.GetInitCommands()
	if err != nil {
		panic(err)
	}

	out := map[string][]byte{}
	for buildTagName, commands := range initCommands {
		path := filepath.Join(p.outputRoot, buildTagName+".gen.go")
		data, err := templates.Execute(templates.InitCommandsTpl, commands)
		if err != nil {
			panic(err)
		}
		out[path] = data
	}

	path := filepath.Join(p.outputRoot, "main.gen.go")
	data, err := templates.Execute(templates.MainFileTpl, struct{}{})
	if err != nil {
		panic(err)
	}

	out[path] = data
	return out
}

func (p *Processor) Version() string {
	return "1.0.0"
}

func (p *Processor) Name() string {
	return "Cobra"
}
