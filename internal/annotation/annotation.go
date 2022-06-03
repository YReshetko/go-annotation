package annotation

type Annotation struct {
	name   string
	params map[string]string
}

func (a Annotation) Name() string {
	return a.name
}

func (a Annotation) Params() map[string]string {
	return a.params
}
