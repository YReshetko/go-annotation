// Code generated by Mapper annotation processor. DO NOT EDIT.
// versions:
//		go: go1.22.4
//		go-annotation: 0.1.0
//		Mapper: 0.0.1-alpha

package mappers

import (
	_imp_4 "github.com/YReshetko/go-annotation/examples/mapper/internal/models/common"
	_imp_3 "github.com/YReshetko/go-annotation/examples/mapper/internal/models/input"
	_imp_2 "github.com/YReshetko/go-annotation/examples/mapper/internal/models/output"
	"strconv"
	"strings"
	_imp_1 "time"
)

var _ PrimitivesMapper = (*PrimitivesMapperImpl)(nil)

type PrimitivesMapperImpl struct{}

func (_this_ PrimitivesMapperImpl) Primitives(in *_imp_3.Primitives) _imp_2.Primitives {
	out := _imp_2.Primitives{}
	if in != nil {
		out.Bool = in.Bool
		out.Byte = in.Byte
		out.Int = in.Int
		out.Int8 = in.Int8
		out.Int16 = in.Int16
		out.Int32 = in.Int32
		out.Int64 = in.Int64
		out.Float32 = in.Float32
		out.Float64 = in.Float64
		out.Complex64 = in.Complex64
		out.Complex128 = in.Complex128
		out.String = in.String
		out.Uintptr = in.Uintptr
		out.Rune = in.Rune
		out.PtrBool = in.PtrBool
		out.PtrUnt = in.PtrUnt
		out.PtrUnt8 = in.PtrUnt8
		out.PtrUnt16 = in.PtrUnt16
		out.PtrUnt32 = in.PtrUnt32
		out.PtrUnt64 = in.PtrUnt64
		out.PtrByte = in.PtrByte
		out.PtrInt = in.PtrInt
		out.PtrInt8 = in.PtrInt8
		out.PtrInt16 = in.PtrInt16
		out.PtrInt32 = in.PtrInt32
		out.PtrInt64 = in.PtrInt64
		out.PtrFloat32 = in.PtrFloat32
		out.PtrFloat64 = in.PtrFloat64
		out.PtrComplex64 = in.PtrComplex64
		out.PtrComplex128 = in.PtrComplex128
		out.PtrString = in.PtrString
		out.PtrUintptr = in.PtrUintptr
		out.PtrRune = in.PtrRune
		out.CreateAt = &in.CreateAt
		if in.MAP != nil {
			out.MAP = *in.MAP
		}
	}

	return out
}

func (_this_ PrimitivesMapperImpl) ConvertMethod(in0 *Inner, in1 ThisFileStruct) _imp_2.Primitives {
	out0 := _imp_2.Primitives{}
	out0.Bool = func(v string) bool {
		res := v == "1" || strings.ToLower(v) == "true"
		return res
	}(in1.Bool)
	out0.String = in1.String
	out0.PtrUnt64 = func(v string) *uint64 {
		res := func() uint64 { o, _ := strconv.Atoi(v); return uint64(o) }()
		return &res
	}(in1.PtrUnt64)
	out0.PtrInt = &in1.PtrInt
	out0.PtrComplex64 = func(v string) *complex64 {
		res := func() complex64 { o, _ := strconv.ParseComplex(v, 64); return complex64(o) }()
		return &res
	}(in1.PtrComplex64)
	if in0 != nil {
		out0.PtrBool = in0.PtrBool
	}
	if in1.PtrFloat64 != nil {
		out0.PtrFloat64 = func(v string) *float64 {
			res := func() float64 { o, _ := strconv.ParseFloat(v, 64); return o }()
			return &res
		}(*in1.PtrFloat64)
	}
	if in1.PtrComplex128 != nil {
		out0.PtrComplex128 = func(v string) *complex128 {
			res := func() complex128 { o, _ := strconv.ParseComplex(v, 128); return o }()
			return &res
		}(*in1.PtrComplex128)
	}
	if in1.Int != nil {
		out0.Int = *in1.Int
	}
	if in1.Float32 != nil {
		out0.Float32 = func(v int) float32 {
			res := float32(v)
			return res
		}(*in1.Float32)
	}

	return out0
}

func (_this_ PrimitivesMapperImpl) ConvertMethod2(in0 Inner, in1 *ThisFileStruct) _imp_2.Primitives {
	out0 := _imp_2.Primitives{}
	out0.PtrBool = in0.PtrBool
	if in1 != nil {
		out0.Bool = func(v string) bool {
			res := v == "1" || strings.ToLower(v) == "true"
			return res
		}(in1.Bool)
		out0.String = in1.String
		out0.PtrUnt64 = func(v string) *uint64 {
			res := func() uint64 { o, _ := strconv.Atoi(v); return uint64(o) }()
			return &res
		}(in1.PtrUnt64)
		out0.PtrInt = &in1.PtrInt
		out0.PtrComplex64 = func(v string) *complex64 {
			res := func() complex64 { o, _ := strconv.ParseComplex(v, 64); return complex64(o) }()
			return &res
		}(in1.PtrComplex64)
		if in1.PtrFloat64 != nil {
			out0.PtrFloat64 = func(v string) *float64 {
				res := func() float64 { o, _ := strconv.ParseFloat(v, 64); return o }()
				return &res
			}(*in1.PtrFloat64)
		}
		if in1.PtrComplex128 != nil {
			out0.PtrComplex128 = func(v string) *complex128 {
				res := func() complex128 { o, _ := strconv.ParseComplex(v, 128); return o }()
				return &res
			}(*in1.PtrComplex128)
		}
		if in1.Int != nil {
			out0.Int = *in1.Int
		}
		if in1.Float32 != nil {
			out0.Float32 = func(v int) float32 {
				res := float32(v)
				return res
			}(*in1.Float32)
		}
	}

	return out0
}

var _ ConstantMapper = (*ConstantMapperImpl)(nil)

type ConstantMapperImpl struct{}

func (_this_ ConstantMapperImpl) Primitives() _imp_2.Primitives {
	out0 := _imp_2.Primitives{}
	out0.Float32 = func(v string) float32 {
		res := func() float32 { o, _ := strconv.ParseFloat(v, 32); return float32(o) }()
		return res
	}("3.14")
	out0.String = "Hello"
	out0.PtrFloat32 = func(v string) *float32 {
		res := func() float32 { o, _ := strconv.ParseFloat(v, 32); return float32(o) }()
		return &res
	}("3.14")
	_var_0 := "World"
	out0.PtrString = &_var_0

	return out0
}

var _ UUIDMapper = (*UUIDMapperImpl)(nil)

type UUIDMapperImpl struct{}

func (_this_ UUIDMapperImpl) FromUUID(uuid _imp_4.WithUUID) _imp_4.StringForUUID {
	out0 := _imp_4.StringForUUID{}
	out0.UID = uuidToString(uuid.UID)

	return out0
}

var _ BaseStructuresMapper = (*BaseStructuresMapperImpl)(nil)

type BaseStructuresMapperImpl struct{}

func (_this_ BaseStructuresMapperImpl) Structures1(in *_imp_3.StructuresMapping) _imp_2.StructuresMapping {
	out0 := _imp_2.StructuresMapping{}
	if in != nil {
		out0.Field1 = in.Field1
		out0.Field2 = &in.Field2
		out0.Field4 = in.Field4
		if in.Field3 != nil {
			out0.Field3 = *in.Field3
		}
	}

	return out0
}

func (_this_ BaseStructuresMapperImpl) Structures2(in *_imp_3.StructuresMapping2) _imp_2.StructuresMapping2 {
	out0 := _imp_2.StructuresMapping2{}
	if in != nil {
		out0.AnotherField2 = &in.Field2.Field1.Field1
		out0.AnotherField4 = &in.Field2.Field1.Field1
		if in.Field1.Field1.Field2 != nil {
			out0.AnotherField1 = *in.Field1.Field1.Field2
		}
		if in.Field1.Field2 != nil && in.Field1.Field2.Field2 != nil {
			out0.AnotherField3 = *in.Field1.Field2.Field2
		}
	}

	return out0
}

func (_this_ BaseStructuresMapperImpl) Structures3(in1 *_imp_4.Common, in2 _imp_4.Common) _imp_2.StructuresMapping2 {
	out0 := _imp_2.StructuresMapping2{}
	out0.AnotherField4 = &in2
	if in1 != nil {
		out0.AnotherField3 = *in1
	}

	return out0
}

func (_this_ BaseStructuresMapperImpl) Structures4(in1 string, in2 *complex128, in3 *uint64) _imp_2.Primitives {
	out0 := _imp_2.Primitives{}
	out0.String = in1
	out0.PtrString = &in1
	if in3 != nil {
		out0.Uint64 = *in3
	}
	if in2 != nil {
		out0.Complex128 = *in2
		out0.PtrComplex128 = in2
	}

	return out0
}

func (_this_ BaseStructuresMapperImpl) TimeMapping(t _imp_1.Time) _imp_1.Time {
	out0 := _imp_1.Time{}

	return out0
}

var _ FunctionMapper = (*FunctionMapperImpl)(nil)

type FunctionMapperImpl struct{}

func (_this_ FunctionMapperImpl) Function1(in *_imp_3.StructuresMapping2) _imp_2.StructuresMapping2 {
	out0 := _imp_2.StructuresMapping2{}
	if in != nil && in.Field1.Field1.Field2 != nil && in.Field1.Field2 != nil && in.Field1.Field2.Field2 != nil {
		out0.AnotherField1 = _this_.fieldToField(
			in.Field1.Field1.Field1,
			in.Field1.Field1.Field2,
			in.Field1.Field2.Field1,
			in.Field1.Field2.Field2)
	}

	return out0
}

func (_this_ FunctionMapperImpl) Function2(in *_imp_3.StructuresMapping2) _imp_2.StructuresMapping2 {
	out0 := _imp_2.StructuresMapping2{}
	if in != nil && in.Field1.Field1.Field2 != nil && in.Field1.Field2 != nil && in.Field1.Field2.Field2 != nil {
		out0.AnotherField1 = fieldToField(
			in.Field1.Field1.Field1,
			in.Field1.Field1.Field2,
			in.Field1.Field2.Field1,
			in.Field1.Field2.Field2)
	}

	return out0
}

func (_this_ FunctionMapperImpl) fieldToField(in0 _imp_4.Common, in1 *_imp_4.Common, in2 _imp_4.Common, in3 *_imp_4.Common) _imp_4.Common {
	out0 := _imp_4.Common{}
	out0.Field1 = in0.Field1
	out0.Field2 = in0.Field2
	out0.Field3 = in0.Field3
	out0.Slice = in0.Slice

	return out0
}

var _ SliceMapping = (*SliceMappingImpl)(nil)

type SliceMappingImpl struct{}

func (_this_ SliceMappingImpl) Function1(in *_imp_3.SliceStruct) _imp_2.SliceStruct {
	out0 := _imp_2.SliceStruct{}
	if in != nil && in.Slice != nil {
		out0.Slice = sliceInOut(in.Slice)
	}

	return out0
}

func (_this_ SliceMappingImpl) Function2(in *_imp_3.SliceStruct) _imp_2.SliceStruct {
	out0 := _imp_2.SliceStruct{}
	if in != nil && in.Slice != nil {

		_var_0 := *in.Slice
		_var_1 := make([]_imp_2.StructuresMapping2, len(_var_0), len(_var_0))
		for _var_2, _var_3 := range _var_0 {
			_var_1[_var_2] = _this_.genMapper(_var_3)
		}
		out0.Slice = _var_1

	}

	return out0
}

func (_this_ SliceMappingImpl) genMapper(in _imp_3.Local2) _imp_2.StructuresMapping2 {
	out0 := _imp_2.StructuresMapping2{}

	return out0
}

func (_this_ SliceMappingImpl) Function3(in *_imp_3.SliceStruct2) _imp_2.SliceStruct2 {
	out0 := _imp_2.SliceStruct2{}
	if in != nil {

		_var_0 := in.Slice
		_var_1 := make([]_imp_4.Common2, len(_var_0), len(_var_0))
		for _var_2, _var_3 := range _var_0 {
			_var_1[_var_2] = _this_.genMapper2(_var_3)
		}
		out0.Slice = _var_1

	}

	return out0
}

func (_this_ SliceMappingImpl) genMapper2(in _imp_4.Common) _imp_4.Common2 {
	out0 := _imp_4.Common2{}
	out0.Field1 = in.Field1
	out0.Field2 = in.Field2
	out0.Field3 = in.Field3
	out0.Slice = in.Slice

	return out0
}

var _ MapMapping = (*MapMappingImpl)(nil)

type MapMappingImpl struct{}

func (_this_ MapMappingImpl) Function1(in _imp_3.MapStruct) _imp_2.MapStruct {
	out0 := _imp_2.MapStruct{}

	_var_0 := in.Map
	_var_1 := make(map[_imp_4.MapKey]_imp_4.Common2, len(_var_0))
	for _var_2, _var_3 := range _var_0 {
		_var_4, _var_5 := _this_.genMapper(_var_2, _var_3)
		_var_1[_var_4] = _var_5
	}
	out0.Map = _var_1

	return out0
}

func (_this_ MapMappingImpl) genMapper(k _imp_4.MapKey, v _imp_4.Common2) (_imp_4.MapKey, _imp_4.Common2) {
	out0 := _imp_4.MapKey{}
	out1 := _imp_4.Common2{}
	out0.Field1 = k.Field1
	out0.Field2 = k.Field2
	out0.Field3 = k.Field3
	out1.Field1 = k.Field1
	out1.Field2 = k.Field2
	out1.Field3 = k.Field3
	out1.Slice = v.Slice

	return out0, out1
}

func (_this_ MapMappingImpl) Function2(mapStruct _imp_3.MapStruct2) _imp_2.MapStruct2 {
	out0 := _imp_2.MapStruct2{}
	out0.Map = mapStruct.Map

	return out0
}
