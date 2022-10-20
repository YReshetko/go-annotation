package tag

func Parse(target any, params map[string]string) (any, error) {
	return parse(target, params)
}
