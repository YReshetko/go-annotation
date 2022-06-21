package mapper

type Mapper struct {
	Name string `annotation:"defaultValue=*Impl,name=name"`
}

type Mapping struct {
	Target string `annotation:"required=true,name=target"`
	Source string `annotation:"required=true,name=source"`
}
