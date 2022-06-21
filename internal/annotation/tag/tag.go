package tag

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/YReshetko/go-annotation/internal/annotation"
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
}

func Parse(target any, annotation annotation.Annotation) any {
	t := reflect.TypeOf(target)
	if t.Kind() != reflect.Struct {
		fmt.Printf("%s\n", t.Kind())
		panic(errors.New("annotation must be a structure"))
	}

	v := reflect.ValueOf(&target).Elem()
	tmp := reflect.New(v.Elem().Type()).Elem()

	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)

		exported := f.PkgPath == ""
		if !exported {
			continue
		}

		if fn, ok := primitives[f.Type.Kind()]; ok {
			nt, err := newTag(f.Tag.Get("annotation"))
			if err != nil {
				panic(fmt.Sprintf("tag was not built for @%s.%s: %s", t.Name(), f.Name, err))
			}

			fieldName := nt.name
			if len(fieldName) == 0 {
				fieldName = f.Name
			}
			fieldName = toCamelCase(fieldName)

			fieldValue, ok := annotation.Params()[fieldName]
			if !ok {
				fieldValue = nt.defaultValue
			}
			if len(fieldValue) == 0 && nt.required {
				panic(fmt.Sprintf("annotation %s requires %s parameter set", t.Name(), fieldName))
			}

			err = fn(tmp.Field(i), fieldValue)

			if err != nil {
				panic(fmt.Sprintf("unable set primitive by tag for field %s: %s", f.Name, err))
			}
		}

	}
	v.Set(tmp)
	return target
}

func toCamelCase(s string) string {
	return strings.ToLower(string(s[0])) + s[1:]
}

func newTag(t string) (tag, error) {
	if len(t) == 0 {
		return tag{}, nil
	}

	fields := strings.Split(t, ",")
	m := map[string]string{
		"defaultValue": "",
		"required":     "false",
		"name":         "",
	}

	for _, f := range fields {
		fv := strings.Split(f, "=")
		if len(fv) != 2 {
			return tag{}, errors.New("tag violates key=value format")
		}
		if _, ok := m[fv[0]]; !ok {
			return tag{}, errors.New(fmt.Sprintf("tag field %s is not allowed", fv[0]))
		}
		m[fv[0]] = fv[1]
	}

	return tag{
		required:     strings.ToLower(m["required"]) == "true",
		defaultValue: m["defaultValue"],
		name:         m["name"],
	}, nil
}
