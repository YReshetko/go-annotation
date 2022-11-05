package annotations

const (
	structType  = "struct"
	pointerType = "pointer"
)

type Constructor struct {
	Name string `annotation:"name=name,default=New{{.TypeName}}"`
	Type string `annotation:"name=type,default=struct"` // defines return structType or pointerType
}

type Optional struct {
	Name            string `annotation:"name=name,default={{.TypeName}}Option"`
	ConstructorName string `annotation:"name=constructor,default=New{{.TypeName}}"`
	WithPattern     string `annotation:"name=with,default=With{{.FieldName}}"`
	Type            string `annotation:"name=type,default=struct"` // defines return structType or pointerType
}

type Builder struct {
	StructureName   string `annotation:"name=name,default={{.TypeName}}Builder"`
	ConstructorName string `annotation:"name=constructor,default=New{{.TypeName}}Builder"`
	BuildPattern    string `annotation:"name=build,default={{.FieldName}}"`
	BuilderName     string `annotation:"name=terminator,default=Build"`
	Type            string `annotation:"name=type,default=struct"` // defines return structType or pointerType
}

// Init is used for fields initialisation such as slice, map, chan
// If Init.Len and Init.Cap then the values are set by default (chan is non-buffered)
type Init struct {
	Len int `annotation:"name=len,default=-1"`
	Cap int `annotation:"name=cap,default=-1"`
}

type Exclude struct{} // Excludes structure field from constructors
