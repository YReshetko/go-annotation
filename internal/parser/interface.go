package parser

type Annotation interface {
	Name() string
	Params() map[string]string
}

func Parse(doc string) ([]Annotation, error) {
	a, err := parse(doc)
	if err != nil {
		return nil, err
	}

	out := make([]Annotation, len(a))
	for i, val := range a {
		out[i] = val
	}

	return out, nil
}
