package templates

// Simply debug, improve and copy ./main/main.go file to the constant below
var mainFile = `
package main

import (
	"fmt"
	"log/slog"
	"os"
	"reflect"
	"strings"
	"time"
	"unsafe"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"golang.org/x/exp/constraints"
)


var root *cobra.Command


func main() {
	if err := root.Execute(); err != nil {
		fatal(err)
	}
}


func fatal(err error) {
	slog.Default().Error(err.Error())
	os.Exit(1)
}


var flagTypeSetters = map[reflect.Kind]func(*pflag.FlagSet, reflect.Value, string) error{
	reflect.Int: func(flagSet *pflag.FlagSet, field reflect.Value, flagName string) error {
		return fieldSetter(flagSet.GetInt, intConverter[int], field.SetInt, flagName)
	},
	reflect.Int8: func(flagSet *pflag.FlagSet, field reflect.Value, flagName string) error {
		return fieldSetter(flagSet.GetInt8, intConverter[int8], field.SetInt, flagName)
	},
	reflect.Int16: func(flagSet *pflag.FlagSet, field reflect.Value, flagName string) error {
		return fieldSetter(flagSet.GetInt16, intConverter[int16], field.SetInt, flagName)
	},
	reflect.Int32: func(flagSet *pflag.FlagSet, field reflect.Value, flagName string) error {
		return fieldSetter(flagSet.GetInt32, intConverter[int32], field.SetInt, flagName)
	},
	reflect.Int64: func(flagSet *pflag.FlagSet, field reflect.Value, flagName string) error {
		return fieldSetter(flagSet.GetInt64, noopConverter[int64], field.SetInt, flagName)
	},
	reflect.Uint: func(flagSet *pflag.FlagSet, field reflect.Value, flagName string) error {
		return fieldSetter(flagSet.GetUint, uintConverter[uint], field.SetUint, flagName)
	},
	reflect.Uint8: func(flagSet *pflag.FlagSet, field reflect.Value, flagName string) error {
		return fieldSetter(flagSet.GetUint8, uintConverter[uint8], field.SetUint, flagName)
	},
	reflect.Uint16: func(flagSet *pflag.FlagSet, field reflect.Value, flagName string) error {
		return fieldSetter(flagSet.GetUint16, uintConverter[uint16], field.SetUint, flagName)
	},
	reflect.Uint32: func(flagSet *pflag.FlagSet, field reflect.Value, flagName string) error {
		return fieldSetter(flagSet.GetUint32, uintConverter[uint32], field.SetUint, flagName)
	},
	reflect.Uint64: func(flagSet *pflag.FlagSet, field reflect.Value, flagName string) error {
		return fieldSetter(flagSet.GetUint64, uintConverter[uint64], field.SetUint, flagName)
	},
	reflect.Float32: func(flagSet *pflag.FlagSet, field reflect.Value, flagName string) error {
		return fieldSetter(flagSet.GetFloat32, floatConverter[float32], field.SetFloat, flagName)
	},
	reflect.Float64: func(flagSet *pflag.FlagSet, field reflect.Value, flagName string) error {
		return fieldSetter(flagSet.GetFloat64, floatConverter[float64], field.SetFloat, flagName)
	},
	reflect.String: func(flagSet *pflag.FlagSet, field reflect.Value, flagName string) error {
		return fieldSetter(flagSet.GetString, noopConverter[string], field.SetString, flagName)
	},
	reflect.Bool: func(flagSet *pflag.FlagSet, field reflect.Value, flagName string) error {
		return fieldSetter(flagSet.GetBool, noopConverter[bool], field.SetBool, flagName)
	},
}

var flagNameTypeSetters = map[string]func(*pflag.FlagSet, reflect.Value, string) error{
	"Duration": func(flagSet *pflag.FlagSet, field reflect.Value, flagName string) error {
		return fieldSetter(flagSet.GetString, durationConverter, field.SetInt, flagName)
	},
}


func noopConverter[T any](v T) (T, error) {
	return v, nil
}

func intConverter[T constraints.Signed](v T) (int64, error) {
	return int64(v), nil
}
func uintConverter[T constraints.Unsigned](v T) (uint64, error) {
	return uint64(v), nil
}
func floatConverter[T constraints.Float](v T) (float64, error) {
	return float64(v), nil
}

func durationConverter(s string) (int64, error) {
	v, err := time.ParseDuration(s)
	return int64(v), err
}

func customMarshallingSetter(flagSet *pflag.FlagSet, field reflect.Value, flagName string) error {
	flagValue, err := flagSet.GetString(flagName)
	if err != nil {
		return err
	}

	method, ok := reflect.PointerTo(field.Type()).MethodByName("MarshalFlag")
	if !ok {
		return fmt.Errorf("unexpectedly method MarshalFlag is not found on type")
	}

	fieldPointer := reflect.NewAt(field.Type(), unsafe.Pointer(field.UnsafeAddr()))
	res := method.Func.Call([]reflect.Value{fieldPointer, reflect.ValueOf(flagValue)})
	if len(res) != 1 {
		return fmt.Errorf("expected single return value for MarshalFlag method")
	}
	if res[0].IsNil() {
		return nil
	}
	return fmt.Errorf("unable to marshal value: %v", res[0])
}

type flagsMarshaller interface {
	MarshalFlag(string) error
}

func resolveSetter(field reflect.StructField) (func(*pflag.FlagSet, reflect.Value, string) error, bool) {
	marshallerType := reflect.TypeOf((*flagsMarshaller)(nil)).Elem()
	if reflect.PointerTo(field.Type).Implements(marshallerType) {
		return customMarshallingSetter, true
	}
	setter, ok := flagNameTypeSetters[field.Type.Name()]
	if ok {
		return setter, true
	}
	setter, ok = flagTypeSetters[field.Type.Kind()]
	if ok {
		return setter, true
	}

	return nil, false
}


func fieldSetter[T any, V any](extract func(string) (T, error), converter func(T) (V, error), setter func(V), key string) error {
	value, err := extract(key)
	if err != nil {
		return err
	}
	convertedValue, err := converter(value)
	if err != nil {
		return err
	}
	setter(convertedValue)
	return nil
}


func parseFlags(cmd *cobra.Command, executor any) error {
	if reflect.TypeOf(executor).Kind() != reflect.Pointer {
		return fmt.Errorf("expected pointer to a structure, but got '%s' of '%T'", reflect.TypeOf(executor).Kind().String(), executor)
	}
	if reflect.TypeOf(executor).Elem().Kind() != reflect.Struct {
		return fmt.Errorf("expected executor shuld be a structure type, but got  '%s' of '%T'", reflect.TypeOf(executor).Elem().Kind(), executor)
	}

	target := reflect.ValueOf(executor).Elem()
	sourceType := reflect.TypeOf(executor).Elem()

	for i := 0; i < sourceType.NumField(); i++ {
		field := sourceType.Field(i)
		flagName, isPersistent, isInherited, ok := flagName(field.Tag, "flag")
		if !ok {
			continue
		}

		if !target.Field(i).CanSet() {
			return fmt.Errorf("unable to set '%s' flag to '%T.%s', the field should be addressible and exported", flagName, executor, field.Name)
		}

		flagSet := cmd.Flags()
		if isInherited {
			flagSet = cmd.InheritedFlags()
		}
		if isPersistent {
			flagSet = cmd.PersistentFlags()
		}

		setter, ok := resolveSetter(field)
		if !ok {
			return fmt.Errorf("the receiver type '%T.%s' is not primitive, 'time.Duration', or implements 'MarshalFlag(string) error' method", executor, field.Name)
		}

		if err := setter(flagSet, target.Field(i), flagName); err != nil {
			return fmt.Errorf("unable to set value to '%T.%s': %w", executor, field.Name, err)
		}
	}

	return nil
}


func flagName(tag reflect.StructTag, key string) (string, bool, bool, bool) {
	isPersistent := false
	isInherited := false
	tagValue, ok := tag.Lookup(key)
	if !ok || len(tagValue) == 0 {
		return "", false, false, false
	}
	values := strings.Split(tagValue, ",")
	for _, v := range values {
		if v == "inherited" {
			isInherited = true
			break
		}
		if v == "persist" {
			isPersistent = true
			break
		}
	}

	return strings.TrimSpace(values[0]), isPersistent, isInherited, true
}
`
