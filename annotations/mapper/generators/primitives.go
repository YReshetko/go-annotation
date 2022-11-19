package generators

import (
	"fmt"
	"github.com/YReshetko/go-annotation/annotations/mapper/generators/nodes"
	"github.com/YReshetko/go-annotation/annotations/mapper/templates"
)

func isBothPrimitives(f1, f2 nodes.Type) bool {
	_, ok1 := f1.(*nodes.PrimitiveType)
	_, ok2 := f2.(*nodes.PrimitiveType)
	return ok1 && ok2
}

func mapPrimitives(toName, fromName string, toField, fromField *nodes.PrimitiveType, fromPrefix []string, c *cache, ic importCache) error {
	if toField.Name() == fromField.Name() {
		if toField.IsPointer() == fromField.IsPointer() {
			c.addIfClause(fromPrefix, fmt.Sprintf("%s = %s", toName, fromName))
			return nil
		}
		if toField.IsPointer() {
			c.addIfClause(fromPrefix, fmt.Sprintf("%s = &%s", toName, fromName))
			return nil
		}

		c.addIfClause(append(fromPrefix, fromName), fmt.Sprintf("%s = *%s", toName, fromName))
		return nil
	}

	cnv, imp := getPrimitiveConverter(toField.Name(), fromField.Name())
	if len(cnv) == 0 {
		return nil
	}
	if len(imp) > 0 {
		ic.AddImport(imp)
	}

	cnvData, err := templates.Execute(templates.PrimitiveConverterFuncTemplate, map[string]interface{}{
		"ReceiverName":      toName,
		"SourceType":        fromField.Name(),
		"IsPointerReceiver": toField.IsPointer(),
		"ReceiverType":      toField.Name(),
		"MappingLine":       cnv,
		"IsPointerSource":   fromField.IsPointer(),
		"SourceName":        fromName,
	})

	if err != nil {
		return fmt.Errorf("unable to build mapper template: %w", err)
	}

	if fromField.IsPointer() {
		c.addIfClause(append(fromPrefix, fromName), string(cnvData))
		return nil
	}
	c.addIfClause(fromPrefix, string(cnvData))
	return nil
}

func mapConstant(toName, toType, value string, IsPointerReceiver bool, c *cache, ic importCache) error {
	if toType == "string" {
		if IsPointerReceiver {
			interimVar := c.nextVar()
			c.addCodeLine(fmt.Sprintf(`%s := "%s"`, interimVar, value))
			c.addCodeLine(fmt.Sprintf(`%s = &%s`, toName, interimVar))
			return nil
		}
		c.addCodeLine(fmt.Sprintf(`%s = "%s"`, toName, value))
		return nil
	}

	cnv, imp := getPrimitiveConverter(toType, "string")
	if len(cnv) == 0 {
		return nil
	}
	if len(imp) > 0 {
		ic.AddImport(imp)
	}

	cnvData, err := templates.Execute(templates.PrimitiveConverterFuncTemplate, map[string]interface{}{
		"ReceiverName":      toName,
		"SourceType":        "string",
		"IsPointerReceiver": IsPointerReceiver,
		"ReceiverType":      toType,
		"MappingLine":       cnv,
		"IsPointerSource":   false,
		"SourceName":        "\"" + value + "\"",
	})

	if err != nil {
		return fmt.Errorf("unable to build mapper template: %w", err)
	}

	c.addCodeLine(string(cnvData))
	return nil
}

func getPrimitiveConverter(toType, fromType string) (string, string) {
	_, to := numeric[toType]
	_, from := numeric[fromType]
	if to && from {
		return toType + "(v)", ""
	}

	cnv, _ := converters[toType+"_"+fromType]
	return cnv[0], cnv[1]

}

// For numeric to numeric we can use to_type(value) (for example uint8(v))
var numeric = map[string]struct{}{
	"uint":    {},
	"uint8":   {},
	"uint16":  {},
	"uint32":  {},
	"uint64":  {},
	"byte":    {},
	"int":     {},
	"int8":    {},
	"int16":   {},
	"int32":   {},
	"int64":   {},
	"float32": {},
	"float64": {},
	"uintptr": {},
	"rune":    {},
}

// Schema to_from + related imports
var converters = map[string][2]string{
	"bool_uint":            {"v == 1", ""},
	"bool_uint8":           {"v == 1", ""},
	"bool_uint16":          {"v == 1", ""},
	"bool_uint32":          {"v == 1", ""},
	"bool_uint64":          {"v == 1", ""},
	"bool_byte":            {"v == 1", ""},
	"bool_int":             {"v == 1", ""},
	"bool_int8":            {"v == 1", ""},
	"bool_int16":           {"v == 1", ""},
	"bool_int32":           {"v == 1", ""},
	"bool_int64":           {"v == 1", ""},
	"bool_float32":         {"v == 1", ""},
	"bool_float64":         {"v == 1", ""},
	"bool_complex64":       {"v == 1", ""},
	"bool_complex128":      {"v == 1", ""},
	"bool_string":          {"v == \"1\" || strings.ToLower(v) == \"true\"", "strings"},
	"bool_uintptr":         {"v == 1", ""},
	"bool_rune":            {"v == 1", ""},
	"uint_bool":            {"func() uint {if v {return 1}; return 0}()", ""},
	"uint_complex64":       {"", ""},
	"uint_complex128":      {"", ""},
	"uint_string":          {"func() uint {o, _ := strconv.Atoi(v); return uint(o)}()", "strconv"},
	"uint8_bool":           {"func() uint8 {if v {return 1}; return 0}()", ""},
	"uint8_complex64":      {"", ""},
	"uint8_complex128":     {"", ""},
	"uint8_string":         {"func() uint8 {o, _ := strconv.Atoi(v); return uint8(o)}()", "strconv"},
	"uint16_bool":          {"func() uint16 {if v {return 1}; return 0}()", ""},
	"uint16_complex64":     {"", ""},
	"uint16_complex128":    {"", ""},
	"uint16_string":        {"func() uint16 {o, _ := strconv.Atoi(v); return uint16(o)}()", "strconv"},
	"uint32_bool":          {"func() uint32 {if v {return 1}; return 0}()", ""},
	"uint32_complex64":     {"", ""},
	"uint32_complex128":    {"", ""},
	"uint32_string":        {"func() uint32 {o, _ := strconv.Atoi(v); return uint32(o)}()", "strconv"},
	"uint64_bool":          {"func() uint64 {if v {return 1}; return 0}()", ""},
	"uint64_complex64":     {"", ""},
	"uint64_complex128":    {"", ""},
	"uint64_string":        {"func() uint64 {o, _ := strconv.Atoi(v); return uint64(o)}()", "strconv"},
	"byte_bool":            {"func() byte {if v {return 1}; return 0}()", ""},
	"byte_complex64":       {"", ""},
	"byte_complex128":      {"", ""},
	"byte_string":          {"func() byte {o, _ := strconv.Atoi(v); return byte(o)}()", "strconv"},
	"int_bool":             {"func() int {if v {return 1}; return 0}()", ""},
	"int_complex64":        {"", ""},
	"int_complex128":       {"", ""},
	"int_string":           {"func() int {o, _ := strconv.Atoi(v); return int(o)}()", "strconv"},
	"int8_bool":            {"func() int8 {if v {return 1}; return 0}()", ""},
	"int8_complex64":       {"", ""},
	"int8_complex128":      {"", ""},
	"int8_string":          {"func() int8 {o, _ := strconv.Atoi(v); return int8(o)}()", "strconv"},
	"int16_bool":           {"func() int16 {if v {return 1}; return 0}()", ""},
	"int16_complex64":      {"", ""},
	"int16_complex128":     {"", ""},
	"int16_string":         {"func() int16 {o, _ := strconv.Atoi(v); return int16(o)}()", "strconv"},
	"int32_bool":           {"func() int32 {if v {return 1}; return 0}()", ""},
	"int32_complex64":      {"", ""},
	"int32_complex128":     {"", ""},
	"int32_string":         {"func() int32 {o, _ := strconv.Atoi(v); return int32(o)}()", "strconv"},
	"int64_bool":           {"func() int64 {if v {return 1}; return 0}()", ""},
	"int64_complex64":      {"", ""},
	"int64_complex128":     {"", ""},
	"int64_string":         {"func() int64 {o, _ := strconv.Atoi(v); return int64(o)}()", "strconv"},
	"float32_bool":         {"func() float32 {if v {return 1}; return 0}()", ""},
	"float32_complex64":    {"", ""},
	"float32_complex128":   {"", ""},
	"float32_string":       {"func() float32 { o, _ := strconv.ParseFloat(v, 32); return float32(o) }()", "strconv"},
	"float64_bool":         {"func() float64 {if v {return 1}; return 0}()", ""},
	"float64_complex64":    {"", ""},
	"float64_complex128":   {"", ""},
	"float64_string":       {"func() float64 {o,_:=strconv.ParseFloat(v, 64); return o}()", "strconv"},
	"complex64_bool":       {"func() complex64 {if v {return 1}; return 0}()", ""},
	"complex64_uint":       {"complex64(complex(float64(v), 0))", ""},
	"complex64_uint8":      {"complex64(complex(float64(v), 0))", ""},
	"complex64_uint16":     {"complex64(complex(float64(v), 0))", ""},
	"complex64_uint32":     {"complex64(complex(float64(v), 0))", ""},
	"complex64_uint64":     {"complex64(complex(float64(v), 0))", ""},
	"complex64_byte":       {"complex64(complex(float64(v), 0))", ""},
	"complex64_int":        {"complex64(complex(float64(v), 0))", ""},
	"complex64_int8":       {"complex64(complex(float64(v), 0))", ""},
	"complex64_int16":      {"complex64(complex(float64(v), 0))", ""},
	"complex64_int32":      {"complex64(complex(float64(v), 0))", ""},
	"complex64_int64":      {"complex64(complex(float64(v), 0))", ""},
	"complex64_float32":    {"complex64(complex(float64(v), 0))", ""},
	"complex64_float64":    {"complex64(complex(v, 0))", ""},
	"complex64_complex128": {"complex64(v)", ""},
	"complex64_string":     {" func() complex64 { o, _ := strconv.ParseComplex(v, 64); return complex64(o) }()", "strconv"},
	"complex64_uintptr":    {"complex64(complex(float64(v), 0))", ""},
	"complex64_rune":       {"complex64(complex(float64(v), 0))", ""},
	"complex128_bool":      {"func() complex128 {if v {return 1}; return 0}()", ""},
	"complex128_uint":      {"complex(float64(v), 0)", ""},
	"complex128_uint8":     {"complex(float64(v), 0)", ""},
	"complex128_uint16":    {"complex(float64(v), 0)", ""},
	"complex128_uint32":    {"complex(float64(v), 0)", ""},
	"complex128_uint64":    {"complex(float64(v), 0)", ""},
	"complex128_byte":      {"complex(float64(v), 0)", ""},
	"complex128_int":       {"complex(float64(v), 0)", ""},
	"complex128_int8":      {"complex(float64(v), 0)", ""},
	"complex128_int16":     {"complex(float64(v), 0)", ""},
	"complex128_int32":     {"complex(float64(v), 0)", ""},
	"complex128_int64":     {"complex(float64(v), 0)", ""},
	"complex128_float32":   {"complex(float64(v), 0)", ""},
	"complex128_float64":   {"complex(v, 0)", ""},
	"complex128_complex64": {"complex128(v)", ""},
	"complex128_string":    {"func() complex128 { o, _ := strconv.ParseComplex(v, 128); return o }()", "strconv"},
	"complex128_uintptr":   {"complex(float64(v), 0)", ""},
	"complex128_rune":      {"complex(float64(v), 0)", ""},
	"string_bool":          {"func() string {if v {return \"true\"}; return \"false\"}()", ""},
	"string_uint":          {"strconv.Itoa(int(v))", "strconv"},
	"string_uint8":         {"strconv.Itoa(int(v))", "strconv"},
	"string_uint16":        {"strconv.Itoa(int(v))", "strconv"},
	"string_uint32":        {"strconv.Itoa(int(v))", "strconv"},
	"string_uint64":        {"strconv.Itoa(int(v))", "strconv"},
	"string_byte":          {"strconv.Itoa(int(v))", "strconv"},
	"string_int":           {"strconv.Itoa(v)", "strconv"},
	"string_int8":          {"strconv.Itoa(int(v))", "strconv"},
	"string_int16":         {"strconv.Itoa(int(v))", "strconv"},
	"string_int32":         {"strconv.Itoa(int(v))", "strconv"},
	"string_int64":         {"strconv.Itoa(int(v))", "strconv"},
	"string_float32":       {"strconv.FormatFloat(float64(v), 'g', 6, 32)", "strconv"},
	"string_float64":       {"strconv.FormatFloat(v, 'g', 6, 64)", "strconv"},
	"string_complex64":     {"strconv.FormatComplex(complex128(v), 'g', 6, 64)", "strconv"},
	"string_complex128":    {"strconv.FormatComplex(v, 'g', 6, 128)", "strconv"},
	"string_uintptr":       {"strconv.Itoa(int(v))", "strconv"},
	"string_rune":          {"strconv.Itoa(int(v))", "strconv"},
	"uintptr_bool":         {"func() uintptr {if v {return 1}; return 0}()", ""},
	"uintptr_complex64":    {"", ""},
	"uintptr_complex128":   {"", ""},
	"uintptr_string":       {"func() uintptr {o, _ := strconv.Atoi(v); return uintptr(o)}()", "strconv"},
	"rune_bool":            {"func() rune {if v {return 1}; return 0}()", ""},
	"rune_complex64":       {"", ""},
	"rune_complex128":      {"", ""},
	"rune_string":          {"func() rune {o, _ := strconv.Atoi(v); return rune(o)}()", "strconv"},
}

func isPrimitive(t string) bool {
	_, ok := map[string]struct{}{
		"bool":       {},
		"uint":       {},
		"uint8":      {},
		"uint16":     {},
		"uint32":     {},
		"uint64":     {},
		"byte":       {},
		"int":        {},
		"int8":       {},
		"int16":      {},
		"int32":      {},
		"int64":      {},
		"float32":    {},
		"float64":    {},
		"complex64":  {},
		"complex128": {},
		"string":     {},
		"uintptr":    {},
		"rune":       {},
	}[t]
	return ok
}
