package input

import "github.com/YReshetko/go-annotation/examples/mapper/internal/models/common"

type Primitives struct {
	Bool          bool
	Unt           uint
	Unt8          uint8
	Unt16         uint16
	Unt32         uint32
	Unt64         uint64
	Byte          byte
	Int           int
	Int8          int8
	Int16         int16
	Int32         int32
	Int64         int64
	Float32       float32
	Float64       float64
	Complex64     complex64
	Complex128    complex128
	String        string
	Uintptr       uintptr
	Rune          rune
	PtrBool       *bool
	PtrUnt        *uint
	PtrUnt8       *uint8
	PtrUnt16      *uint16
	PtrUnt32      *uint32
	PtrUnt64      *uint64
	PtrByte       *byte
	PtrInt        *int
	PtrInt8       *int8
	PtrInt16      *int16
	PtrInt32      *int32
	PtrInt64      *int64
	PtrFloat32    *float32
	PtrFloat64    *float64
	PtrComplex64  *complex64
	PtrComplex128 *complex128
	PtrString     *string
	PtrUintptr    *uintptr
	PtrRune       *rune
	MAP           *rune
}

type PkgSpecific struct {
	Field1 string
}

type StructuresMapping struct {
	Field1 common.Common
	Field2 common.Common
	Field3 *common.Common
	Field4 *common.Common

	PkgSpecific1 PkgSpecific
	PkgSpecific2 PkgSpecific
	PkgSpecific3 *PkgSpecific
	PkgSpecific4 *PkgSpecific
}

type StructuresMapping2 struct {
	Field1 Local
	Field2 Local
	Field3 *Local
	Field4 *Local
}

type Local struct {
	Field1 Local2
	Field2 *Local2
}

type Local2 struct {
	Field1 common.Common
	Field2 *common.Common
}
