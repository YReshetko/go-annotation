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
	"uint":    UintFlagType,
	"uint8":   Uint8FlagType,
	"uint16":  Uint16FlagType,
	"uint32":  Uint32FlagType,
	"uint64":  Uint64FlagType,
	"float32": Float32FlagType,
	"float64": Float64FlagType,
}

var flagTypeStringDefaultValues = map[FlagType]string{
	BoolFlagType:    "false",
	IntFlagType:     "0",
	Int8FlagType:    "0",
	Int16FlagType:   "0",
	Int32FlagType:   "0",
	Int64FlagType:   "0",
	UintFlagType:    "0",
	Uint8FlagType:   "0",
	Uint16FlagType:  "0",
	Uint32FlagType:  "0",
	Uint64FlagType:  "0",
	Float32FlagType: "0",
	Float64FlagType: "0",
}

func FlagTypeByASTPrimitive(astType string) FlagType {
	t, ok := astTypesToFlagType[astType]
	if !ok {
		return StringFlagType
	}
	return t
}

func FlagTypeDefaultValue(flagType FlagType) string {
	v, ok := flagTypeStringDefaultValues[flagType]
	if !ok {
		return ""
	}
	return v
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
