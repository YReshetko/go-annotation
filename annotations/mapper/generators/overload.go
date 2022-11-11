package generators

import (
	"fmt"
	"strings"
)

type mappingType byte

const (
	none mappingType = iota
	source
	function
	constant
)

type mapping struct {
	source   string
	this     string
	function string
	constant string
}

func (m mapping) getMapping() (string, mappingType) {
	switch {
	case len(m.source) != 0:
		return m.source, source
	case len(m.this) != 0:
		return "_this_." + m.this, function
	case len(m.function) != 0:
		return m.function, function
	case len(m.constant) != 0:
		return m.constant, constant
	}
	return "", none
}

// @Constructor(name="newOverloading", type="pointer")
type overloading struct {
	isIgnoreDefault bool
	mappings        map[string]mapping // @Init
}

func (o *overloading) Add(target, source, this, function, constant string) error {
	var nonEmpty int
	for _, s := range []string{source, this, function, constant} {
		if len(s) != 0 {
			nonEmpty++
		}
	}
	if nonEmpty != 1 {
		return fmt.Errorf("invalid Mapping annotation for %s, expected exactly one of [source, this, func, const] non empty, but got %d", target, nonEmpty)
	}
	o.mappings[target] = mapping{
		source:   source,
		this:     this,
		function: function,
		constant: constant,
	}
	return nil
}

func (o *overloading) find(field string) (string, mappingType) {
	v, ok := o.mappings[field]
	if ok {
		return v.getMapping()
	}
	if !strings.Contains(field, ".") {
		return "", none
	}

	v, ok = o.mappings[field[strings.Index(field, ".")+1:]]
	if ok {
		return v.getMapping()
	}
	return "", none
}
