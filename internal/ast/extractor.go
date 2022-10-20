package ast

import (
	"go/ast"
	"reflect"
	"strings"
)

func extractComment(n ast.Node) (string, bool) {
	v, ok := valueOf(n)
	if !ok {
		return "", false
	}
	cgs := extractCommentGroups(v)
	if len(cgs) == 0 {
		return "", false
	}

	sb := strings.Builder{}
	for _, cg := range cgs {
		sb.WriteString(cg.Text())
	}

	return strings.TrimRight(sb.String(), "\n"), true
}

func extractCommentGroups(v reflect.Value) []*ast.CommentGroup {
	var out []*ast.CommentGroup
	for i := 0; i < v.NumField(); i++ {
		f := v.Field(i)
		if f.Kind() != reflect.Pointer {
			continue
		}
		cg, ok := f.Interface().(*ast.CommentGroup)
		if !ok {
			continue
		}
		if cg != nil {
			out = append(out, cg)
		}
	}
	return out
}

func valueOf(n ast.Node) (reflect.Value, bool) {
	if n == nil {
		return reflect.Value{}, false
	}
	switch reflect.TypeOf(n).Kind() {
	case reflect.Pointer:
		e := reflect.ValueOf(n).Elem()
		if e.Kind() == reflect.Struct {
			return e, true
		}
	case reflect.Struct:
		return reflect.ValueOf(n), true

	}
	return reflect.Value{}, false
}
