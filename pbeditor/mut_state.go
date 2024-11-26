package pbeditor

import (
	"github.com/yaoguangduan/proto-editor/protodef"
	"google.golang.org/protobuf/encoding/protowire"
	"google.golang.org/protobuf/proto"
	"log"
)

type MutStateInfo struct {
	*protodef.MutState
}

func (ms *MutStateInfo) MarshalToFieldBytes() []byte {
	var buf []byte
	buf = protowire.AppendTag(buf, 20000, protowire.BytesType)
	marshal, err := proto.Marshal(ms)
	if err != nil {
		log.Println("marshal err:", err)
	}
	buf = protowire.AppendBytes(buf, marshal)
	return buf
}
func CreateMutStateInfo() *MutStateInfo {
	return &MutStateInfo{&protodef.MutState{}}
}

func UnmarshalMutState(raw []byte) MutStateInfo {
	ms := MutStateInfo{&protodef.MutState{}}

	for len(raw) > 0 {
		num, typ, n := protowire.ConsumeTag(raw)
		if num == 20000 && typ == protowire.BytesType {
			raw = raw[n:]
			bys, _ := protowire.ConsumeBytes(raw)
			err := proto.Unmarshal(bys, ms)
			if err != nil {
				return ms
			}
			return ms
		}
		num, typ, n = protowire.ConsumeField(raw)
		if n <= 0 {
			return ms
		}
		raw = raw[n:]
	}
	return ms
}
func (ms *MutStateInfo) ListOpLog(number int32) []*protodef.ListOpLog {
	if ms.Array == nil {
		ms.Array = make(map[int32]*protodef.ListMutState)
	}
	listMut, ok := ms.Array[number]
	if !ok {
		return nil
	}
	return listMut.GetLog()
}
func (ms *MutStateInfo) ListAppendLog(number int32, log *protodef.ListOpLog) {
	if ms.Array == nil {
		ms.Array = make(map[int32]*protodef.ListMutState)
	}
	listMut, ok := ms.Array[number]
	if !ok {
		ms.Array[number] = &protodef.ListMutState{}
		listMut = ms.Array[number]
		listMut.Log = append(listMut.Log, log)
	} else {
		listMut.Log = append(listMut.Log, log)
	}
}
func (ms *MutStateInfo) AddStr(number int32, v string) {
	if ms.Map == nil {
		ms.Map = make(map[int32]*protodef.MapMutState)
	}
	mapMut, ok := ms.Map[number]
	if !ok {
		ms.Map[number] = &protodef.MapMutState{}
		mapMut = ms.Map[number]
	}
	mapMut.StrList = append(mapMut.StrList, v)
}
func (ms *MutStateInfo) AddBool(number int32, v bool) {
	if ms.Map == nil {
		ms.Map = make(map[int32]*protodef.MapMutState)
	}
	mapMut, ok := ms.Map[number]
	if !ok {
		ms.Map[number] = &protodef.MapMutState{}
		mapMut = ms.Map[number]
	}
	mapMut.BoolList = append(mapMut.BoolList, v)
}
func (ms *MutStateInfo) AddI32(number int32, v int32) {
	if ms.Map == nil {
		ms.Map = make(map[int32]*protodef.MapMutState)
	}
	mapMut, ok := ms.Map[number]
	if !ok {
		ms.Map[number] = &protodef.MapMutState{}
		mapMut = ms.Map[number]
	}
	mapMut.I32List = append(mapMut.I32List, v)
}
func (ms *MutStateInfo) AddI64(number int32, v int64) {
	if ms.Map == nil {
		ms.Map = make(map[int32]*protodef.MapMutState)
	}
	mapMut, ok := ms.Map[number]
	if !ok {
		ms.Map[number] = &protodef.MapMutState{}
		mapMut = ms.Map[number]
	}
	mapMut.I64List = append(mapMut.I64List, v)
}
func (ms *MutStateInfo) AddU32(number int32, v uint32) {
	if ms.Map == nil {
		ms.Map = make(map[int32]*protodef.MapMutState)
	}
	mapMut, ok := ms.Map[number]
	if !ok {
		ms.Map[number] = &protodef.MapMutState{}
		mapMut = ms.Map[number]
	}
	mapMut.U32List = append(mapMut.U32List, v)
}
func (ms *MutStateInfo) AddU64(number int32, v uint64) {
	if ms.Map == nil {
		ms.Map = make(map[int32]*protodef.MapMutState)
	}
	mapMut, ok := ms.Map[number]
	if !ok {
		ms.Map[number] = &protodef.MapMutState{}
		mapMut = ms.Map[number]
	}
	mapMut.U64List = append(mapMut.U64List, v)
}
func (ms *MutStateInfo) StrList(number int32) []string {
	if ms.Map == nil {
		ms.Map = make(map[int32]*protodef.MapMutState)
	}
	mapMut, ok := ms.Map[number]
	if ok {
		return mapMut.StrList
	}
	return nil
}

func (ms *MutStateInfo) I32List(number int32) []int32 {
	if ms.Map == nil {
		ms.Map = make(map[int32]*protodef.MapMutState)
	}
	mapMut, ok := ms.Map[number]
	if ok {
		return mapMut.I32List
	}
	return nil
}

func (ms *MutStateInfo) I64List(number int32) []int64 {
	if ms.Map == nil {
		ms.Map = make(map[int32]*protodef.MapMutState)
	}
	mapMut, ok := ms.Map[number]
	if ok {
		return mapMut.I64List
	}
	return nil
}

func (ms *MutStateInfo) U32List(number int32) []uint32 {
	if ms.Map == nil {
		ms.Map = make(map[int32]*protodef.MapMutState)
	}
	mapMut, ok := ms.Map[number]
	if ok {
		return mapMut.U32List
	}
	return nil
}

func (ms *MutStateInfo) U64List(number int32) []uint64 {
	if ms.Map == nil {
		ms.Map = make(map[int32]*protodef.MapMutState)
	}
	mapMut, ok := ms.Map[number]
	if ok {
		return mapMut.U64List
	}
	return nil
}

func (ms *MutStateInfo) BoolList(number int32) []bool {
	if ms.Map == nil {
		ms.Map = make(map[int32]*protodef.MapMutState)
	}
	mapMut, ok := ms.Map[number]
	if ok {
		return mapMut.BoolList
	}
	return nil
}
