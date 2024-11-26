package pbeditor

import (
	"errors"
	"fmt"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/encoding/protowire"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"math"
	"unsafe"
)

var ErrParseRawFields = errors.New("parse raw fields error")

// RawField 内部的Field只会是bool或者repeated int32 int64 uint32 uint64 string bool
type RawField struct {
	Number  int32
	Type    protowire.Type
	Bytes   []byte
	Variant uint64
	F32     float32
	F64     float64
}

func Str2Bys(str string) []byte {
	return unsafe.Slice(unsafe.StringData(str), len(str))
}

func Bys2Str(bys []byte) string {
	return unsafe.String(unsafe.SliceData(bys), len(bys))
}

type RawMessage struct {
	RawFields map[int32][]*RawField
	MutState  *MutStateInfo
}

func (rm *RawMessage) Marshal() []byte {
	var buf []byte
	for _, f := range rm.RawFields {
		for _, field := range f {
			buf = protowire.AppendTag(buf, protowire.Number(field.Number), field.Type)
			if field.Type == protowire.VarintType {
				buf = protowire.AppendVarint(buf, uint64(1))
			} else {
				buf = protowire.AppendBytes(buf, field.Bytes)
			}

		}
	}

	return buf
}

func (rm *RawMessage) AddString(num int32, val string) {
	bys := []byte(val)
	rm.RawFields[num] = append(rm.RawFields[num], &RawField{Number: num, Type: protowire.BytesType, Bytes: bys})
}

func (rm *RawMessage) GetStringList(num int32) []string {
	_, exist := rm.RawFields[num]
	if exist {
		ret := make([]string, len(rm.RawFields[num]))
		for i, f := range rm.RawFields[num] {
			ret[i] = string(f.Bytes)
		}
		return ret
	} else {
		return make([]string, 0)
	}
}
func (rm *RawMessage) GetInt32List(num int32) []int32 {
	f, exist := rm.RawFields[num]
	if exist {
		ret := make([]int32, 0)
		tmp := f[0].Bytes
		for len(tmp) > 0 {
			val, n := protowire.ConsumeVarint(tmp)
			if n < 0 {
				panic(n)
			}
			tmp = tmp[n:]
			ret = append(ret, int32(val))
		}
		return ret
	} else {
		return make([]int32, 0)
	}
}
func (rm *RawMessage) GetUint32List(num int32) []uint32 {
	f, exist := rm.RawFields[num]
	if exist {
		ret := make([]uint32, 0)
		tmp := f[0].Bytes
		for len(tmp) > 0 {
			val, n := protowire.ConsumeVarint(tmp)
			if n < 0 {
				panic(n)
			}
			tmp = tmp[n:]
			ret = append(ret, uint32(val))
		}
		return ret
	} else {
		return make([]uint32, 0)
	}
}
func (rm *RawMessage) GetInt64List(num int32) []int64 {
	f, exist := rm.RawFields[num]
	if exist {
		ret := make([]int64, 0)
		tmp := f[0].Bytes
		for len(tmp) > 0 {
			val, n := protowire.ConsumeVarint(tmp)
			if n < 0 {
				panic(n)
			}
			tmp = tmp[n:]
			ret = append(ret, int64(val))
		}
		return ret
	} else {
		return make([]int64, 0)
	}
}
func (rm *RawMessage) GetUint64List(num int32) []uint64 {
	f, exist := rm.RawFields[num]
	if exist {
		ret := make([]uint64, 0)
		tmp := f[0].Bytes
		for len(tmp) > 0 {
			val, n := protowire.ConsumeVarint(tmp)
			if n < 0 {
				panic(n)
			}
			tmp = tmp[n:]
			ret = append(ret, val)
		}
		return ret
	} else {
		return make([]uint64, 0)
	}
}
func (rm *RawMessage) GetBoolList(num int32) []bool {
	f, exist := rm.RawFields[num]
	if exist {
		ret := make([]bool, 0)
		tmp := f[0].Bytes
		for len(tmp) > 0 {
			val, n := protowire.ConsumeVarint(tmp)
			if n < 0 {
				panic(n)
			}
			tmp = tmp[n:]
			ret = append(ret, val > 0)
		}
		return ret
	} else {
		return make([]bool, 0)
	}
}
func (rm *RawMessage) SetBool(num int32) {
	_, exist := rm.RawFields[num]
	if !exist {
		rm.RawFields[num] = append(rm.RawFields[num], &RawField{Number: num, Type: protowire.VarintType})
	}
}
func (rm *RawMessage) GetBool(num int32) bool {
	_, exist := rm.RawFields[num]
	return exist
}

func (rm *RawMessage) ClearBool(num int32) {
	delete(rm.RawFields, num)
}

func (rm *RawMessage) AddVarint(num int32, val uint64) {
	rm.RawFields[num] = append(rm.RawFields[num], &RawField{Number: num, Type: protowire.VarintType, Variant: val})
}

type MapKey interface {
	string | int32 | int64 | uint32 | uint64 | bool
}

func FieldTypeToWireType(field *protogen.Field) protowire.Type {
	if field.Desc.Kind() == protoreflect.EnumKind || field.Desc.Kind() == protoreflect.BoolKind || field.Desc.Kind() == protoreflect.Int32Kind || field.Desc.Kind() == protoreflect.Int64Kind ||
		field.Desc.Kind() == protoreflect.Uint64Kind || field.Desc.Kind() == protoreflect.Uint32Kind || field.Desc.Kind() == protoreflect.Uint64Kind || field.Desc.Kind() == protoreflect.Sint64Kind ||
		field.Desc.Kind() == protoreflect.Sint32Kind || field.Desc.Kind() == protoreflect.Sfixed32Kind || field.Desc.Kind() == protoreflect.Sfixed64Kind || field.Desc.Kind() == protoreflect.Fixed32Kind || field.Desc.Kind() == protoreflect.Fixed64Kind {
		return protowire.VarintType
	}
	if field.Desc.Kind() == protoreflect.DoubleKind {
		return protowire.Fixed64Type
	}
	if field.Desc.Kind() == protoreflect.FloatKind {
		return protowire.Fixed32Type
	}
	return protowire.BytesType
}

func ToRawMessage(buf []byte) RawMessage {
	rm := RawMessage{RawFields: make(map[int32][]*RawField)}
	for len(buf) > 0 {
		numb, typ, n := protowire.ConsumeTag(buf)
		if n < 0 {
			panic(fmt.Sprintf("parse error:consumeTag:%d", n))
		}
		num := int32(numb)
		buf = buf[n:]
		if num == 20000 && typ == protowire.BytesType {
			bys, n := protowire.ConsumeBytes(buf)
			if n < 0 {
				panic(fmt.Sprintf("parse error:ConsumeBytes:%d", n))
			}
			mutState := CreateMutStateInfo()
			_ = proto.Unmarshal(bys, mutState)
			rm.MutState = mutState
			buf = buf[n:]
			continue
		}
		switch typ {
		case protowire.VarintType:
			v, n := protowire.ConsumeVarint(buf)
			if n < 0 {
				panic(fmt.Sprintf("parse error:ConsumeVarint:%d", n))
			}
			buf = buf[n:]
			rm.RawFields[num] = append(rm.RawFields[num], &RawField{Number: num, Type: protowire.VarintType, Variant: v})
		case protowire.Fixed32Type:
			u32, n := protowire.ConsumeFixed32(buf)
			if n < 0 {
				panic(fmt.Sprintf("parse error:ConsumeFixed32:%d", n))
			}
			rm.RawFields[num] = append(rm.RawFields[num], &RawField{Number: num, Type: protowire.Fixed32Type, F32: math.Float32frombits(u32)})
			buf = buf[n:]
		case protowire.Fixed64Type:
			u64, n := protowire.ConsumeFixed64(buf)
			if n < 0 {
				panic(fmt.Sprintf("parse error:ConsumeFixed32:%d", n))
			}
			rm.RawFields[num] = append(rm.RawFields[num], &RawField{Number: num, Type: protowire.Fixed32Type, F64: math.Float64frombits(u64)})
			buf = buf[n:]
		case protowire.BytesType:
			bytes, n := protowire.ConsumeBytes(buf)
			if n < 0 {
				panic(fmt.Sprintf("parse error:ConsumeBytes:%d", n))
			}
			rm.RawFields[num] = append(rm.RawFields[num], &RawField{Number: num, Type: protowire.BytesType, Bytes: bytes})
			buf = buf[n:]
		default:
			panic(fmt.Sprintf("unsupported type: %d", typ))
		}
	}
	return rm
}

func VarintRange(tmp []byte, f func(v uint64)) {
	for len(tmp) > 0 {
		val, n := protowire.ConsumeVarint(tmp)
		if n < 0 {
			panic(ErrParseRawFields)
		}
		tmp = tmp[n:]
		f(val)
	}
}

func Fixed32Range(tmp []byte, f func(v float32)) {
	for len(tmp) > 0 {
		val, n := protowire.ConsumeFixed32(tmp)
		if n < 0 {
			panic(ErrParseRawFields)
		}
		tmp = tmp[n:]
		f(math.Float32frombits(val))
	}
}

func Fixed64Range(tmp []byte, f func(v float64)) {
	for len(tmp) > 0 {
		val, n := protowire.ConsumeFixed64(tmp)
		if n < 0 {
			panic(ErrParseRawFields)
		}
		tmp = tmp[n:]
		f(math.Float64frombits(val))
	}
}

type MapKV struct {
	keyVarint    uint64
	keyBytes     []byte
	valueVarint  uint64
	valueBytes   []byte
	valueFixed32 float32
	valueFixed64 float64
}

func (mkv *MapKV) KeyI32() int32 {
	return int32(mkv.keyVarint)
}
func (mkv *MapKV) KeyI64() int64 {
	return int64(mkv.keyVarint)
}
func (mkv *MapKV) KeyU32() uint32 {
	return uint32(mkv.keyVarint)
}
func (mkv *MapKV) KeyU64() uint64 {
	return uint64(mkv.keyVarint)
}
func (mkv *MapKV) KeyStr() string {
	return Bys2Str(mkv.keyBytes)
}
func (mkv *MapKV) KeyBool() bool {
	return mkv.keyVarint > 0
}

func (mkv *MapKV) ValI32() int32 {
	return int32(mkv.valueVarint)
}
func (mkv *MapKV) ValI64() int64 {
	return int64(mkv.valueVarint)
}
func (mkv *MapKV) ValU32() uint32 {
	return uint32(mkv.valueVarint)
}
func (mkv *MapKV) ValU64() uint64 {
	return uint64(mkv.valueVarint)
}
func (mkv *MapKV) ValStr() string {
	return Bys2Str(mkv.valueBytes)
}
func (mkv *MapKV) ValBytes() []byte {
	return mkv.valueBytes
}
func (mkv *MapKV) ValF32() float32 {
	return mkv.valueFixed32
}
func (mkv *MapKV) ValF64() float64 {
	return mkv.valueFixed64
}
func (mkv *MapKV) ValBool() bool {
	return mkv.valueVarint > 0
}

func ParseMap(tmp []byte) *MapKV {
	num, tp, n := protowire.ConsumeTag(tmp)
	if n < 0 || num != 1 {
		panic(ErrParseRawFields)
	}
	tmp = tmp[n:]
	mapKv := &MapKV{}
	if tp == protowire.BytesType {
		bys, n := protowire.ConsumeBytes(tmp)
		if n < 0 {
			panic(ErrParseRawFields)
		}
		tmp = tmp[n:]
		mapKv.keyBytes = bys
	} else {
		v, n := protowire.ConsumeVarint(tmp)
		if n < 0 {
			panic(ErrParseRawFields)
		}
		tmp = tmp[n:]
		mapKv.keyVarint = v
	}
	num, tp, n = protowire.ConsumeTag(tmp)
	if n < 0 || num != 2 {
		panic(ErrParseRawFields)
	}
	tmp = tmp[n:]
	if tp == protowire.BytesType {
		bys, n := protowire.ConsumeBytes(tmp)
		if n < 0 {
			panic(ErrParseRawFields)
		}
		tmp = tmp[n:]
		mapKv.valueBytes = bys
	} else if tp == protowire.Fixed32Type {
		u32, n := protowire.ConsumeFixed32(tmp)
		if n < 0 {
			panic(ErrParseRawFields)
		}
		tmp = tmp[n:]
		mapKv.valueFixed32 = math.Float32frombits(u32)
	} else if tp == protowire.Fixed64Type {
		u64, n := protowire.ConsumeFixed64(tmp)
		if n < 0 {
			panic(ErrParseRawFields)
		}
		tmp = tmp[n:]
		mapKv.valueFixed64 = math.Float64frombits(u64)
	} else {
		v, n := protowire.ConsumeVarint(tmp)
		if n < 0 {
			panic(ErrParseRawFields)
		}
		tmp = tmp[n:]
		mapKv.valueVarint = v
	}

	return mapKv
}
