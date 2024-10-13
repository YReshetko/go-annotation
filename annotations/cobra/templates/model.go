package templates

type InitCommands struct {
	BuildTag string
	Imports  []Import
	Commands []Command
}

type Import struct {
	Alias   string
	Package string
}

type Command struct {
	IsRoot        bool
	VarName       string
	ParentVarName string

	Use           string
	Example       string
	Short         string
	Long          string
	SilenceUsage  bool
	SilenceErrors bool

	Flags    []Flag
	Handlers []Handler
}

type FlagType string

const (
	StringFlagType  FlagType = "String"
	BoolFlagType    FlagType = "Bool"
	IntFlagType     FlagType = "Int"
	Int8FlagType    FlagType = "Int8"
	Int16FlagType   FlagType = "Int16"
	Int32FlagType   FlagType = "Int32"
	Int64FlagType   FlagType = "Int64"
	UintFlagType    FlagType = "Uint"
	Uint8FlagType   FlagType = "Uint8"
	Uint16FlagType  FlagType = "Uint16"
	Uint32FlagType  FlagType = "Uint32"
	Uint64FlagType  FlagType = "Uint64"
	Float32FlagType FlagType = "Float32"
	Float64FlagType FlagType = "Float64"
)

var astTypesToFlagType = map[string]FlagType{
	"string":  StringFlagType,
	"bool":    BoolFlagType,
	"int":     IntFlagType,
	"int8":    Int8FlagType,
	"int16":   Int16FlagType,
	"int32":   Int32FlagType,
	"int64":   Int64FlagType,
	"unt":     UintFlagType,
	"unt8":    Uint8FlagType,
	"unt16":   Uint16FlagType,
	"unt32":   Uint32FlagType,
	"unt64":   Uint64FlagType,
	"float32": Float32FlagType,
	"float64": Float64FlagType,
}

func FlagTypeByASTPrimitive(astType string) FlagType {
	t, ok := astTypesToFlagType[astType]
	if !ok {
		return StringFlagType
	}
	return t
}

type Flag struct {
	Type FlagType

	Name         string
	Shorthand    string
	DefaultValue string
	Description  string

	IsRequired   bool
	IsPersistent bool
}

type Handler struct {
	MethodName           string
	ExecutorPackageAlias string
	ExecutorTypeName     string

	IsPreRun        bool
	IsPostRun       bool
	IsPersistentRun bool

	HasReturn bool
}
