package protos

import (
	"fmt"
	"github.com/yaoguangduan/proto-editor/pbeditor"
	"github.com/yaoguangduan/proto-editor/pbgen/protos_user"
	"github.com/yaoguangduan/proto-editor/protodef"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/encoding/protowire"
	"google.golang.org/protobuf/proto"
	"math/rand/v2"
	"slices"
	"testing"
)

func TestMut(t *testing.T) {
	u := &protos_user.User{}
	var buf = make([]byte, 0)
	buf = protowire.AppendTag(buf, 20000, protowire.BytesType)
	ms := &protodef.MutState{Map: map[int32]*protodef.MapMutState{
		12: &protodef.MapMutState{StrList: make([]string, 0)},
	}}
	marshal, err := proto.Marshal(ms)
	if err != nil {
		panic(err)
		protowire.ConsumeField()
	}
	buf = protowire.AppendBytes(buf, marshal)
	err = proto.Unmarshal(buf, u)
	if err != nil {
		panic(err)
	}
	fmt.Println(u.ProtoReflect().GetUnknown())
	msn, err := pbeditor.GetMutState(u.ProtoReflect().GetUnknown())

	if err != nil {
		panic(err)
	}
	fmt.Println(msn)
}

func TestList(t *testing.T) {
	original := []int32{4, 2, 6, 4, 8, 7}
	il := pbeditor.NewI32List()
	il.AddAll(original)
	cp := pbeditor.NewI32List()
	for _, v := range il.Values() {
		cp.Add(v)
	}
	cp.FlushDirty(false)
	il.FlushDirty(false)
	for i := 0; i < 1000; i++ {
		rn := rand.IntN(4) + 1
		pos := rand.IntN(il.Len()+1) - 1
		if pos == -1 {
			pos = 0
		}
		switch rn {
		case 1:
			il.Add(int32(i))
			break
		case 2:
			if il.Len() <= 0 {
				il.Clear()
			}
			break
		case 3:
			if il.Len() <= 0 {
				il.Add(int32(i))
			}
			break
		case 4:
			il.Clear()
			break
		}
	}
	cp.MergeDirtyFromPb(il.ChangeLog(), il.MonitorDirtyList())
	fmt.Println(il.Values())
	fmt.Println(cp.Values())
}

func TestListError(t *testing.T) {
	original := []int32{4, 2, 6, 4, 8, 7}
	il := pbeditor.NewI32List()
	il.Clear()
	il.Add(222)
	dirtyList := il.MonitorDirtyList()
	for _, log := range il.ChangeLog() {
		fmt.Println(il.Values())
		fmt.Println(dirtyList)
		fmt.Println(protojson.Format(log))
		switch log.OpType {
		case protodef.ListOpType_Insert:
			if int(log.Pos) >= len(original) {
				original = append(original, dirtyList[log.RefIdx])
			} else {
				original = slices.Insert(original, int(log.Pos), dirtyList[log.RefIdx])
			}
		case protodef.ListOpType_Delete:
			original = slices.Delete(original, int(log.Pos), int(log.Pos+1))
		case protodef.ListOpType_Clear:
			original = nil
		case protodef.ListOpType_Modify:
			original[log.Pos] = dirtyList[log.RefIdx]
		}
	}
	fmt.Println(original)
}
