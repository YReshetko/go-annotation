package generators

import (
	"errors"
	"fmt"
	"strings"

	"github.com/YReshetko/go-annotation/annotations/mapper/generators/nodes"
)

func generate(to *nodes.Field, in []*nodes.Field, c *cache, o *overloading, imp importCache) error {
	switch toType := to.Type().(type) {
	case *nodes.PrimitiveType:
		// TODO generate mapping for primitives
	case *nodes.StructType:
		return generateStructMapping(to.Name(), toType, in, c, o, imp)
	case *nodes.ArrayType:
		// TODO support slice output mapping
	case *nodes.MapType:
		// TODO support map output mapping
	}
	return nil
}

func generateStructMapping(name string, toType *nodes.StructType, in []*nodes.Field, c *cache, o *overloading, imp importCache) error {
	for _, field := range toType.Fields() {
		if !field.IsExported() {
			continue
		}
		longFieldName := name + "." + field.Name()
		m := o.find(longFieldName)
		if m != nil && m.mappingType != none {
			err := override(m, longFieldName, field.Type(), in, c, imp)
			if err != nil {
				return fmt.Errorf("unable to build override for %s: %w", longFieldName, err)
			}
			continue
		}

		if o.isIgnoreDefault {
			continue
		}

		potentialFields := fieldsDefaultMapping(field.Name(), in)
		if len(potentialFields) == 0 {
			continue
		}

		if len(potentialFields) > 1 {
			// TODO support selection by appropriate type
		}

		var fromPrefix []string
		if potentialFields[0].fromRootType.IsPointer() {
			fromPrefix = append(fromPrefix, potentialFields[0].fromRootName)
		}

		fromName := nodes.VariableNameJoin(potentialFields[0].fromRootName, potentialFields[0].fromName)
		if isBothPrimitives(field.Type(), potentialFields[0].fromType) {
			err := mapPrimitives(longFieldName, fromName, field.Type().(*nodes.PrimitiveType), potentialFields[0].fromType.(*nodes.PrimitiveType), fromPrefix, c, imp)
			if err != nil {
				return fmt.Errorf("unable to build default primitive mapping: %w", err)
			}
			continue
		}

		if isBothStructures(field.Type(), potentialFields[0].fromType) {
			err := mapStructures(longFieldName, fromName, field.Type().(*nodes.StructType), potentialFields[0].fromType.(*nodes.StructType), fromPrefix, c)
			if err != nil {
				return fmt.Errorf("unable to build default structures mapping: %w", err)
			}
			continue
		}
		if isEqualSlices(field.Type(), potentialFields[0].fromType) {
			err := assignSlice(longFieldName, fromName, field.Type().(*nodes.ArrayType), potentialFields[0].fromType.(*nodes.ArrayType), fromPrefix, c)
			if err != nil {
				return fmt.Errorf("unable to build default slice mapping: %w", err)
			}
			continue
		}
		if isEqualMaps(field.Type(), potentialFields[0].fromType) {
			err := assignMap(longFieldName, fromName, field.Type().(*nodes.MapType), potentialFields[0].fromType.(*nodes.MapType), fromPrefix, c)
			if err != nil {
				return fmt.Errorf("unable to build default map mapping: %w", err)
			}
			continue
		}
	}

	return nil
}

type defaultFieldsMapping struct {
	fromRootName string
	fromRootType nodes.Type
	fromName     string
	fromType     nodes.Type
}

func fieldsDefaultMapping(name string, in []*nodes.Field) []defaultFieldsMapping {
	var out []defaultFieldsMapping
	for _, field := range in {
		switch ft := field.Type().(type) {
		case *nodes.PrimitiveType:
			if name == ft.Name() {
				out = append(out, defaultFieldsMapping{
					fromName: field.Name(),
					fromType: ft,
				})
			}
		case *nodes.StructType:
			for _, f := range ft.Fields() {
				if !f.IsExported() {
					continue
				}
				if f.Name() == name {
					out = append(out, defaultFieldsMapping{
						fromRootName: field.Name(),
						fromRootType: field.Type(),
						fromName:     f.Name(),
						fromType:     f.Type(),
					})
				}
			}
		}
	}
	return out
}

func extractNilCheck(fromLine string, in []*nodes.Field) ([]string, bool, error) {
	names := strings.Split(fromLine, ".")
	for _, field := range in {
		buff, isFromPointer, err := field.FindNilChecks("", names)
		if err == nil {
			return buff, isFromPointer, nil
		}
		if errors.Is(err, nodes.WrongPathErr) {
			continue
		}
		return nil, false, fmt.Errorf("unable to find field by path: %w", err)
	}
	return nil, false, errors.New("unable to find variable path")
}
