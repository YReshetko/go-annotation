package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"golang.org/x/exp/constraints"
	"log/slog"
	"os"
	"reflect"
	"strings"
)

// Generate by default
var root *cobra.Command

// Generate by default
func main() {
	if err := root.Execute(); err != nil {
		fatal(err)
	}
}

// Generate by default
func fatal(err error) {
	slog.Default().Error(err.Error())
	os.Exit(1)
}

// Generate by default
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

// Generate by default
func noopConverter[T any](v T) T {
	return v
}

func intConverter[T constraints.Signed](v T) int64 {
	return int64(v)
}
func uintConverter[T constraints.Unsigned](v T) uint64 {
	return uint64(v)
}
func floatConverter[T constraints.Float](v T) float64 {
	return float64(v)
}

// Generate by default
func fieldSetter[T any, V any](extract func(string) (T, error), converter func(T) V, setter func(V), key string) error {
	value, err := extract(key)
	if err != nil {
		return err
	}
	setter(converter(value))
	return nil
}

// Generate by default
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
		flagName, isInherited, ok := flagName(field.Tag, "flag")
		if !ok {
			continue
		}

		if !target.Field(i).CanSet() {
			return fmt.Errorf("unable to set '%s' flag to '%T.%s', the field should be addressible and exported", flagName, executor, field.Name)
		}

		setter, ok := flagTypeSetters[field.Type.Kind()]
		if !ok {
			return fmt.Errorf("no type setter fot '%s' flag to '%T.%s' please develop one", flagName, executor, field.Name)
		}
		flagSet := cmd.Flags()
		if isInherited {
			flagSet = cmd.InheritedFlags()
		}

		if err := setter(flagSet, target.Field(i), flagName); err != nil {
			return fmt.Errorf("unable to set value to '%T.%s': %w", executor, field.Name, err)
		}
	}

	return nil
}

// Generate by default
func flagName(tag reflect.StructTag, key string) (string, bool, bool) {
	isPersistent := false
	tagValue, ok := tag.Lookup(key)
	if !ok || len(tagValue) == 0 {
		return "", false, false
	}
	values := strings.Split(tagValue, ",")
	for _, v := range values {
		if v == "inherited" {
			isPersistent = true
			break
		}
	}

	return strings.TrimSpace(values[0]), isPersistent, true
}
