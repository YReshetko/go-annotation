package mappers

import (
	"github.com/YReshetko/go-annotation/examples/mapper/internal/models/common"
	"github.com/YReshetko/go-annotation/examples/mapper/internal/models/input"
	"github.com/YReshetko/go-annotation/examples/mapper/internal/models/output"
	"github.com/google/uuid"
	"time"
)

type LinkToThisFileStruct ThisFileStruct
type LinkToLinkToThisFileStruct LinkToThisFileStruct

// PrimitivesMapper example of primitives mapper
// @Mapper
type PrimitivesMapper interface {
	Primitives(in *input.Primitives) (out output.Primitives)
	ConvertMethod(*Inner, ThisFileStruct) output.Primitives
	ConvertMethod2(Inner, *ThisFileStruct) output.Primitives
}

// ConstantMapper example for constant mapping
// @Mapper
type ConstantMapper interface {
	// Primitives
	// @Mapping(target="PtrFloat32", const="3.14")
	// @Mapping(target="Float32", const="3.14")
	// @Mapping(target="String", const="Hello")
	// @Mapping(target="PtrString", const="World")
	Primitives() output.Primitives
}

// UUIDMapper example for arrays
// @Mapper
type UUIDMapper interface {
	// @Mapping(target="UID", func="uuidToString(uuid.UID)")
	FromUUID(uuid common.WithUUID) common.StringForUUID
}

func uuidToString(uuid *uuid.UUID) string {
	return uuid.String()
}

// BaseStructuresMapper example of base structures mapper
// @Mapper
type BaseStructuresMapper interface {
	Structures1(in *input.StructuresMapping) output.StructuresMapping

	// Structures2
	// @Mapping(target="AnotherField1", source="in.Field1.Field1.Field2")
	// @Mapping(target="AnotherField2", source="in.Field2.Field1.Field1")
	// @Mapping(target="AnotherField3", source="in.Field1.Field2.Field2")
	// @Mapping(target="AnotherField4", source="in.Field2.Field1.Field1")
	Structures2(in *input.StructuresMapping2) output.StructuresMapping2

	// @Mapping(target="AnotherField3", source="in1")
	// @Mapping(target="AnotherField4", source="in2")
	Structures3(in1 *common.Common, in2 common.Common) output.StructuresMapping2

	// @Mapping(target="String", source="in1")
	// @Mapping(target="PtrString", source="in1")
	// @Mapping(target="Complex128", source="in2")
	// @Mapping(target="PtrComplex128", source="in2")
	// @Mapping(target="Uint64", source="in3")
	// @Mapping(target="PtrUint64", source="in3")
	Structures4(in1 string, in2 *complex128, in3 *uint64) output.Primitives

	TimeMapping(t time.Time) time.Time
}

// @Mapper
type FunctionMapper interface {
	// @Mapping(target="AnotherField1", this="fieldToField(
	//		in.Field1.Field1.Field1,
	//		in.Field1.Field1.Field2,
	//		in.Field1.Field2.Field1,
	//		in.Field1.Field2.Field2)")
	Function1(in *input.StructuresMapping2) output.StructuresMapping2

	// @Mapping(target="AnotherField1", func="fieldToField(
	//		in.Field1.Field1.Field1,
	//		in.Field1.Field1.Field2,
	//		in.Field1.Field2.Field1,
	//		in.Field1.Field2.Field2)")
	Function2(in *input.StructuresMapping2) output.StructuresMapping2

	fieldToField(common.Common, *common.Common, common.Common, *common.Common) common.Common
}

// @Mapper
type SliceMapping interface {
	// @Mapping(target="Slice", func="sliceInOut(in.Slice)")
	Function1(in *input.SliceStruct) output.SliceStruct

	// @SliceMapping(target="Slice", source="in.Slice", this="genMapper")
	Function2(in *input.SliceStruct) output.SliceStruct
	genMapper(in input.Local2) output.StructuresMapping2

	// @SliceMapping(target="Slice", source="in.Slice", this="genMapper2")
	Function3(in *input.SliceStruct2) output.SliceStruct2
	genMapper2(in common.Common) common.Common2
}

// @Mapper
type MapMapping interface {
	// @MapMapping(target="Map", source="in.Map", this="genMapper")
	Function1(in input.MapStruct) output.MapStruct
	genMapper(k common.MapKey, v common.Common2) (common.MapKey, common.Common2)

	Function2(mapStruct input.MapStruct2) output.MapStruct2
}

func sliceInOut(s *[]input.Local2) []output.StructuresMapping2 {
	return nil
}

func fieldToField(in1 common.Common, in2 *common.Common, in3 common.Common, in4 *common.Common) common.Common {
	return common.Common{}
}

type AnotherStruct struct {
	f ThisFileStruct
}

type ThisFileStruct struct {
	Field         string
	String        string
	PtrUnt64      string
	PtrFloat64    *string
	PtrComplex128 *string
	PtrComplex64  string
	PtrInt        int
	Int           *int
	Float32       *int
	Bool          string
}

func some(pm PrimitivesMapper) {
	//pm.Primitives(&Inner{}, LinkToLinkToThisFileStruct(ThisFileStruct{}))
}
