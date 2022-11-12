package generators

import "github.com/YReshetko/go-annotation/annotations/constructor/templates"

type FileValues struct {
	PackageName string
	HasImports  bool
	Imports     []Import
	Data        string
}

type Data struct {
	Data    []byte
	Imports []Import
}

func Generate(pkgName string, data []Data) []byte {
	di := newDistinctImports()
	var out []byte
	for _, d := range data {
		out = append(out, d.Data...)
		di.mergeSlice(d.Imports)
	}

	fv := FileValues{
		PackageName: pkgName,
		HasImports:  !di.isEmpty(),
		Imports:     di.toSlice(),
		Data:        string(out),
	}

	return templates.Must(templates.Execute(templates.FileTpl, fv))
}
