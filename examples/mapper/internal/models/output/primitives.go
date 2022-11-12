package output

import (
	"github.com/YReshetko/go-annotation/examples/mapper/internal/models/common"
	"time"
)

type Primitives struct {
	Bool          bool
	Uint          uint
	Uint8         uint8
	Uint16        uint16
	Uint32        uint32
	Uint64        uint64
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
	MAP           rune
	CreateAt      *time.Time // TODO: nil pointer as modules does not support loading default Go libs
}

type PkgSpecific struct {
	Field1 string
}

type StructuresMapping struct {
	Field1 common.Common
	Field2 *common.Common
	Field3 common.Common
	Field4 *common.Common

	PkgSpecific1 PkgSpecific
	PkgSpecific2 *PkgSpecific
	PkgSpecific3 PkgSpecific
	PkgSpecific4 *PkgSpecific
}

type StructuresMapping2 struct {
	AnotherField1 common.Common
	AnotherField2 *common.Common
	AnotherField3 common.Common
	AnotherField4 *common.Common
}

type SliceStruct struct {
	Slice []StructuresMapping2
}

type SliceStruct2 struct {
	Slice []common.Common2
}
