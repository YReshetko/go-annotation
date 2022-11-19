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
	slice
	dictionary
)

type mapping struct {
	source      string
	this        string
	function    string
	constant    string
	mappingType mappingType
}

func (m mapping) funcOrThisLine() string {
	if len(m.this) != 0 {
		return "_this_." + m.this
	}
	return m.function
}

// @Constructor(name="newOverloading", type="pointer", exported="false")
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
		source:      source,
		this:        this,
		function:    function,
		constant:    constant,
		mappingType: o.defineMappingType(source, this, function, constant),
	}
	return nil
}

func (o *overloading) AddSlice(target, source, this, function string) error {
	if len(this) > 0 && len(function) > 0 {
		return fmt.Errorf("invalid Mapping annotation for %s, expected exactly one of [this, func] non empty, but got %d", target, 2)
	}

	o.mappings[target] = mapping{
		source:      source,
		this:        this,
		function:    function,
		mappingType: slice,
	}
	return nil
}

func (o *overloading) AddMap(target, source, this, function string) error {
	if len(this) > 0 && len(function) > 0 {
		return fmt.Errorf("invalid Mapping annotation for %s, expected exactly one of [this, func] non empty, but got %d", target, 2)
	}

	o.mappings[target] = mapping{
		source:      source,
		this:        this,
		function:    function,
		mappingType: dictionary,
	}
	return nil
}

func (o *overloading) defineMappingType(s, t, f, c string) mappingType {
	switch {
	case len(s) != 0:
		return source
	case len(t) != 0 || len(f) != 0:
		return function
	case len(c) != 0:
		return constant
	}
	return none
}

func (o *overloading) find(field string) *mapping {
	v, ok := o.mappings[field]
	if ok {
		return &v
	}
	if !strings.Contains(field, ".") {
		return nil
	}

	v, ok = o.mappings[field[strings.Index(field, ".")+1:]]
	if ok {
		return &v
	}
	return nil
}
