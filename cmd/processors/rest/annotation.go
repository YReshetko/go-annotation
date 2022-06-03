package endpoint

type RestEndpoint struct {
	Method string `annotation:"defaultValue=GET"`
	Path   string `annotation:"required=true"`
}
