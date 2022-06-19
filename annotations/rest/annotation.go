package rest

type Rest struct {
	Method string `annotation:"name=method,defaultValue=GET"`
	Path   string `annotation:"name=path,defaultValue=/"`
}
