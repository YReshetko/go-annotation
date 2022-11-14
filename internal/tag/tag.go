package tag

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

type setter func(v reflect.Value, t string) error

var primitives = map[reflect.Kind]setter{
	reflect.Bool:    setBool,
	reflect.Int:     setInt,
	reflect.Int8:    setInt,
	reflect.Int16:   setInt,
	reflect.Int32:   setInt,
	reflect.Int64:   setInt,
	reflect.Uint:    setUint,
	reflect.Uint8:   setUint,
	reflect.Uint16:  setUint,
	reflect.Uint32:  setUint,
	reflect.Uint64:  setUint,
	reflect.Float32: setFloat,
	reflect.Float64: setFloat,
	reflect.String:  setString,
}

func setInt(v reflect.Value, value string) error {
	i, err := strconv.ParseInt(value, 0, 0)
	if err != nil {
		return err
	}
	v.SetInt(i)
	return nil
}
func setUint(v reflect.Value, value string) error {
	i, err := strconv.ParseUint(value, 0, 0)
	if err != nil {
		return err
	}
	v.SetUint(i)
	return nil
}
func setFloat(v reflect.Value, value string) error {
	i, err := strconv.ParseFloat(value, 0)
	if err != nil {
		return err
	}
	v.SetFloat(i)
	return nil
}
func setBool(v reflect.Value, value string) error {
	b, err := strconv.ParseBool(value)
	if err != nil {
		return err
	}
	v.SetBool(b)
	return nil
}
func setString(v reflect.Value, value string) error {
	v.SetString(value)
	return nil
}

type tag struct {
	defaultValue string
	required     bool
	name         string // name of field in annotation
	oneOf        map[string]struct{}
}

func parse(target any, params map[string]string) (any, error) {
	t := reflect.TypeOf(target)
	if t.Kind() != reflect.Struct {
		return nil, fmt.Errorf("annotation %s must be a structure", t.Name())
	}
	iter := newIterator(t)

	v := reflect.ValueOf(&target).Elem()
	tmp := reflect.New(v.Elem().Type()).Elem()

	for iter.hasNext() {
		f, i := iter.next()

		nt, ok := newTag(f.Tag)
		if !ok {
			continue
		}

		value, err := nt.value(f.Name, params)
		if err != nil {
			return nil, fmt.Errorf("unable to get value for annotation %s: %w", t.Name(), err)
		}

		fn, ok := primitives[f.Type.Kind()]
		if !ok {
			continue
		}
		err = fn(tmp.Field(i), value)

		if err != nil {
			return nil, fmt.Errorf("unable set primitive by tag for field %s: %s", f.Name, err)
		}
	}
	v.Set(tmp)
	return target, nil
}

func toCamelCase(s string) string {
	return strings.ToLower(string(s[0])) + s[1:]
}

const (
	annotationName = "annotation"

	paramDefault  = "default"
	paramRequired = "required"
	paramName     = "name"
	oneOf         = "oneOf"

	equalSign  = "="
	falseValue = "false"
	trueValue  = "true"
)

func newTag(st reflect.StructTag) (tag, bool) {
	t, ok := st.Lookup(annotationName)
	if !ok {
		return tag{}, false
	}
	if len(t) == 0 {
		return tag{}, true
	}

	fields := strings.Split(t, ",")
	m := map[string]any{
		paramDefault:  "",
		paramRequired: falseValue,
		paramName:     "",
		oneOf:         map[string]struct{}{},
	}

	for _, f := range fields {
		k, v := tagParamKeyValue(f)
		if _, ok := m[k]; !ok {
			panic(errors.New(fmt.Sprintf("tag field %s is not allowed", k)))
		}
		m[k] = v
	}

	return tag{
		required:     m[paramRequired].(string) == trueValue,
		defaultValue: m[paramDefault].(string),
		name:         m[paramName].(string),
		oneOf:        m[oneOf].(map[string]struct{}),
	}, true
}

func (t tag) paramName(fieldName string) string {
	if len(t.name) > 0 {
		return t.name
	}
	return fieldName
}

func (t tag) value(fieldName string, params map[string]string) (string, error) {
	fn := toCamelCase(t.paramName(fieldName))
	fieldValue, ok := params[fn]
	if !ok {
		fieldValue = t.defaultValue
	}
	if len(fieldValue) == 0 && t.required {
		return "", fmt.Errorf("parameter %s is required, but not set", fn)
	}

	if len(t.oneOf) == 0 {
		return fieldValue, nil
	}
	if _, ok := t.oneOf[fieldValue]; !ok {
		return "", fmt.Errorf(`field "%s" is restricted to [%v] values, but got "%s"`, fieldName, strings.Join(mapToSlice(t.oneOf), ", "), fieldValue)
	}

	return fieldValue, nil
}

func mapToSlice(m map[string]struct{}) []string {
	o := make([]string, len(m), len(m))
	i := 0
	for k, _ := range m {
		o[i] = k
		i++
	}
	return o
}

func tagParamKeyValue(s string) (string, any) {
	if !strings.Contains(s, equalSign) {
		if s == paramRequired {
			return paramRequired, trueValue
		}
		panic(errors.New("tag has only one option without values: '" + paramRequired + "'"))

	}

	fv := strings.Split(s, equalSign)
	if len(fv) != 2 {
		panic(errors.New("tag violates key=value format"))
	}
	if fv[0] == oneOf {
		rv := strings.Split(fv[1], ";")
		m := map[string]struct{}{}
		for _, v := range rv {
			m[v] = struct{}{}
		}
		return fv[0], m
	}
	return fv[0], fv[1]
}
