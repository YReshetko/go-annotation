package parser

type annotation struct {
	name   string
	params map[string]string
}

func (a annotation) Name() string {
	return a.name
}

func (a annotation) Params() map[string]string {
	return a.params
}
