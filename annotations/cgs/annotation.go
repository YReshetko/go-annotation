// cgs - abbreviation of CGS (Constructor - Getters - Setters)

package cgs

const (
	structType  = "struct"
	pointerType = "pointer"
)

type Constructor struct {
	Name string `annotation:"name=name,default=New{{.TypeName}}"`
	Type string `annotation:"name=type,default=struct"` // defines return structType or pointerType
}

type Exclude struct{} // Excludes structure field from constructor

type Getter struct {
	Name string `annotation:"name=name,default={{.FieldName}}"`
}

type Setter struct {
	Name string `annotation:"name=name,default={{.FieldName}}"`
}
