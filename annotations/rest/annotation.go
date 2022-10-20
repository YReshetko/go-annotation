package rest

type Rest struct {
	Method string `annotation:"name=method,default=GET"`
	Path   string `annotation:"name=path,default=/"`
}
