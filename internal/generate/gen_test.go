package generate

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/yaoguangduan/proto-editor/pbeditor"
	"github.com/yaoguangduan/proto-editor/pbgen/protos_world"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"log"
	"maps"
	"math"
	"slices"
	"testing"
)

func mockTestEditorData() *protos_world.TestEditor {
	testEditor := protos_world.NewTestEditor()
	testEditor.SetI32(32).SetI64(64).SetU32(23).SetU64(46).
		SetStr("str").SetB(true).SetF32(32.32).SetF64(64.64).
		SetBys(pbeditor.Str2Bys("bytes")).SetTe(protos_world.TestEnum_Test1).
		SetMsg(protos_world.NewTestMsgEditor().SetField1(31).SetField2("msg").SetInnerMsg(protos_world.NewPlainMsgEditor().SetI32(64).SetB(true).SetStr("inner")))
	testEditor.AddF32List(32.32).AddF32List(32.23).AddF64List(64.4).AddF64List(64.6).AddBList(true).AddBList(false).AddStrList("list1").AddStrList("list2")
	testEditor.AddI32List(3).AddI32List(2).AddI64List(5).AddI64List(4)
	testEditor.AddU32List(23).AddU32List(33).AddU64List(66).AddU64List(66)
	testEditor.AddMsgList(protos_world.NewTestMsgEditor().SetField1(11)).AddMsgList(protos_world.NewTestMsgEditor().SetField1(22).SetField2("fffff"))
	testEditor.GetMsgMap().Put(111, protos_world.NewTestMsgEditor().SetField1(222))
	testEditor.GetMsgMap().Put(333, protos_world.NewTestMsgEditor().SetField1(444))
	testEditor.GetEnumMap().Put(555, protos_world.TestEnum_Test0)
	testEditor.GetEnumMap().Put(666, protos_world.TestEnum_Test1)
	testEditor.GetEnumMap().Put(666, protos_world.TestEnum_Test2)
	testEditor.GetU64F64Map().Put(64, 64.64).Put(46, 46.46)
	testEditor.GetStrF32Map().Put("f32", 32.32).Put("f64", 64.64)
	testEditor.GetI32StrMap().Put(32, "32i").Put(23, "23i")
	return testEditor
}

func TestMapBytes(t *testing.T) {
	dirty := &protos_world.Test{}
	dirty.I32StrMap = map[int32]string{}
	dirty.I32StrMap[32] = "32i"
	dirty.I32StrMap[23] = "23i"
	marshal, err := proto.Marshal(dirty)
	if err != nil {
		return
	}
	fmt.Println((marshal))
}

func TestBytesOp(t *testing.T) {
	editor := mockTestEditorData()

	bytes := editor.GatherDirtyToBytes()
	fmt.Println(bytes)
	dirty := &protos_world.Test{}
	err := proto.Unmarshal(bytes, dirty)
	if err != nil {
		panic(err)
	}
	assert.Nil(t, err)
	mutState := pbeditor.UnmarshalMutState(dirty.ProtoReflect().GetUnknown())

	dirtyPb := &protos_world.Test{}
	editor.GatherDirtyToPb(dirtyPb)
	mutStateFromPb := pbeditor.UnmarshalMutState(dirtyPb.ProtoReflect().GetUnknown())

	dirtyPb.ProtoReflect().SetUnknown(nil)
	dirty.ProtoReflect().SetUnknown(nil)
	assert.Equal(t, dirtyPb, dirty)
	assert.Equal(t, mutState, mutStateFromPb)

	t.Log(protojson.MarshalOptions{Indent: ""}.Format(dirty))
	t.Log(protojson.MarshalOptions{Indent: ""}.Format(dirtyPb))

	empty := protos_world.NewTestEditor()
	empty.MergeDirtyFromBytes(bytes)
	gathered := &protos_world.Test{}
	empty.GatherDirtyToPb(gathered)
	gathered.ProtoReflect().SetUnknown(nil)
	bys, err := proto.Marshal(gathered)
	assert.Nil(t, err)
	newg := &protos_world.Test{}
	err = proto.Unmarshal(bys, newg)
	assert.Nil(t, err)
	t.Log(protojson.MarshalOptions{Indent: ""}.Format(dirty))
	t.Log(protojson.MarshalOptions{Indent: ""}.Format(gathered))

}

func TestCopy(t *testing.T) {
	testEditor := mockTestEditorData()
	t.Log(maps.Collect(testEditor.GetEnumMap().Iter()))
	test := &protos_world.Test{}
	testEditor.CopyTo(test)
	testEditorEmpty := protos_world.NewTestEditor()
	testEditorEmpty.CopyFrom(test)
	testEditor.FlushDirty(false)
	testEditorEmpty.FlushDirty(false)
	testCopy := &protos_world.Test{}
	testEditorEmpty.CopyTo(testCopy)
	assert.Equal(t, testEditor, testEditorEmpty)
	assert.Equal(t, test, testCopy)
}

func TestDirtyPb(t *testing.T) {
	testEditor := mockTestEditorData()
	test := &protos_world.Test{}
	testEditor.GatherDirtyToPb(test)

	testEditorAnother := protos_world.NewTestEditor()
	testEditorAnother.MergeDirtyFromPb(test)

	testAnother := &protos_world.Test{}
	testEditorAnother.GatherDirtyToPb(testAnother)

	unknown1 := pbeditor.UnmarshalMutState(test.ProtoReflect().GetUnknown())
	unknown2 := pbeditor.UnmarshalMutState(test.ProtoReflect().GetUnknown())
	t.Log(unknown1)
	t.Log(unknown2)
	assert.Equal(t, unknown1, unknown2)
	test.ProtoReflect().SetUnknown(nil)
	testAnother.ProtoReflect().SetUnknown(nil)
	assert.Equal(t, test, testAnother)
}

func TestRealWorldOp(t *testing.T) {
	testEditor := mockTestEditorData()
	test := &protos_world.Test{}
	testEditor.CopyTo(test)

	editor := protos_world.NewTestEditor()
	editor.CopyFrom(test)
	editor.FlushDirty(false)

	incr := &protos_world.Test{}
	editor.GatherDirtyToPb(incr)
	log.Println(protojson.Format(incr))
	assert.Equal(t, protojson.Format(incr), "{}")

	mirror := protos_world.NewTestEditor()
	mirror.CopyFrom(test)
	mirror.FlushDirty(false)

	//i32
	var i32 int32 = math.MaxInt32
	editor.SetI32(i32)
	incr = &protos_world.Test{}
	editor.GatherDirtyToPb(incr)
	assert.Equal(t, *incr.I32, i32)
	mirror.MergeDirtyFromPb(incr)
	assert.Equal(t, mirror.GetI32(), i32)
	//i64
	var i64 int64 = math.MaxInt64
	editor.SetI64(i64)
	incr = &protos_world.Test{}
	editor.GatherDirtyToPb(incr)
	assert.Equal(t, *incr.I64, i64)
	mirror.MergeDirtyFromPb(incr)
	assert.Equal(t, mirror.GetI64(), i64)

	//str
	var str = "new str"
	editor.SetStr(str)
	incr = &protos_world.Test{}
	editor.GatherDirtyToPb(incr)
	assert.Equal(t, *incr.Str, str)
	mirror.MergeDirtyFromPb(incr)
	assert.Equal(t, mirror.GetStr(), str)
	assert.Equal(t, mirror.GetStr(), editor.GetStr())
	//b
	editor.SetB(!editor.GetB())
	incr = &protos_world.Test{}
	editor.GatherDirtyToPb(incr)
	assert.Equal(t, *incr.B, editor.GetB())
	mirror.MergeDirtyFromPb(incr)
	assert.Equal(t, mirror.GetB(), editor.GetB())
	assert.Equal(t, mirror.GetB(), editor.GetB())
	//f
	editor.SetF32(math.MaxFloat32)
	incr = &protos_world.Test{}
	editor.GatherDirtyToPb(incr)
	assert.Equal(t, *incr.F32, editor.GetF32())
	mirror.MergeDirtyFromPb(incr)
	assert.Equal(t, mirror.GetF32(), float32(math.MaxFloat32))
	assert.Equal(t, mirror.GetF32(), editor.GetF32())

	editor.FlushDirty(false)
	editor.AddI32List(math.MaxInt32)
	incr = &protos_world.Test{}
	editor.GatherDirtyToPb(incr)
	assert.Equal(t, incr.GetI32List(), []int32{math.MaxInt32})
	mirror.MergeDirtyFromPb(incr)
	assert.Equal(t, mirror.GetI32List().Values(), editor.GetI32List().Values())

	editor.FlushDirty(false)
	editor.GetI32List().RemoveByIdx(0)
	incr = &protos_world.Test{}
	editor.GatherDirtyToPb(incr)
	assert.Equal(t, incr.I32List, []int32(nil))
	mirror.MergeDirtyFromPb(incr)
	assert.Equal(t, mirror.GetI32List().Values(), editor.GetI32List().Values())

	editor.FlushDirty(false)
	editor.GetI32List().Set(0, math.MinInt32)
	incr = &protos_world.Test{}
	editor.GatherDirtyToPb(incr)
	assert.Equal(t, incr.I32List, []int32{math.MinInt32})
	mirror.MergeDirtyFromPb(incr)
	assert.Equal(t, mirror.GetI32List().Values(), editor.GetI32List().Values())
	t.Log(mirror.GetI32List().Values(), editor.GetI32List().Values())

	editor.FlushDirty(false)
	editor.GetI32List().Clear()
	incr = &protos_world.Test{}
	editor.GatherDirtyToPb(incr)
	assert.Equal(t, incr.I32List, []int32(nil))
	mirror.MergeDirtyFromPb(incr)
	assert.Equal(t, mirror.GetI32List().Values(), editor.GetI32List().Values())
	t.Log(mirror.GetI32List().Values(), editor.GetI32List().Values())

	editor.FlushDirty(false)
	editor.GetStrList().Add("newstritem")
	incr = &protos_world.Test{}
	editor.GatherDirtyToPb(incr)
	assert.Equal(t, incr.StrList, []string{"newstritem"})
	mirror.MergeDirtyFromPb(incr)
	assert.Equal(t, mirror.GetStrList().Values(), editor.GetStrList().Values())
	t.Log(mirror.GetStrList().Values(), editor.GetStrList().Values())

	editor.FlushDirty(false)
	editor.GetStrList().RemoveByIdx(0)
	incr = &protos_world.Test{}
	editor.GatherDirtyToPb(incr)
	assert.Equal(t, incr.StrList, []string(nil))
	mirror.MergeDirtyFromPb(incr)
	assert.Equal(t, mirror.GetStrList().Values(), editor.GetStrList().Values())
	t.Log(mirror.GetStrList().Values(), editor.GetStrList().Values())

	editor.FlushDirty(false)
	editor.GetStrList().Set(0, "newset")
	incr = &protos_world.Test{}
	editor.GatherDirtyToPb(incr)
	assert.Equal(t, incr.StrList, []string{"newset"})
	mirror.MergeDirtyFromPb(incr)
	assert.Equal(t, mirror.GetStrList().Values(), editor.GetStrList().Values())
	t.Log(mirror.GetStrList().Values(), editor.GetStrList().Values())

	editor.FlushDirty(false)
	editor.GetStrList().Clear()
	incr = &protos_world.Test{}
	editor.GatherDirtyToPb(incr)
	assert.Equal(t, incr.StrList, []string(nil))
	mirror.MergeDirtyFromPb(incr)
	assert.Equal(t, mirror.GetStrList().Values(), editor.GetStrList().Values())
	t.Log(mirror.GetStrList().Values(), editor.GetStrList().Values())

	editor.FlushDirty(false)
	editor.GetF64List().Add(math.MaxFloat64)
	incr = &protos_world.Test{}
	editor.GatherDirtyToPb(incr)
	assert.Equal(t, incr.F64List, []float64{math.MaxFloat64})
	mirror.MergeDirtyFromPb(incr)
	assert.Equal(t, mirror.GetF64List().Values(), editor.GetF64List().Values())
	t.Log(mirror.GetF64List().Values(), editor.GetF64List().Values())

	editor.FlushDirty(false)
	editor.GetF64List().RemoveByIdx(2)
	incr = &protos_world.Test{}
	editor.GatherDirtyToPb(incr)
	assert.Equal(t, incr.F64List, []float64(nil))
	mirror.MergeDirtyFromPb(incr)
	assert.Equal(t, mirror.GetF64List().Values(), editor.GetF64List().Values())
	t.Log(mirror.GetF64List().Values(), editor.GetF64List().Values())

	editor.FlushDirty(false)
	editor.GetF64List().Set(1, math.MaxFloat64)
	incr = &protos_world.Test{}
	editor.GatherDirtyToPb(incr)
	assert.Equal(t, incr.F64List, []float64{math.MaxFloat64})
	mirror.MergeDirtyFromPb(incr)
	assert.Equal(t, mirror.GetF64List().Values(), editor.GetF64List().Values())
	t.Log(mirror.GetF64List().Values(), editor.GetF64List().Values())

	editor.FlushDirty(false)
	editor.GetF64List().Clear()
	incr = &protos_world.Test{}
	editor.GatherDirtyToPb(incr)
	assert.Equal(t, incr.F64List, []float64(nil))
	mirror.MergeDirtyFromPb(incr)
	assert.Equal(t, mirror.GetF64List().Values(), editor.GetF64List().Values())
	t.Log(mirror.GetF64List().Values(), editor.GetF64List().Values())

	editor.FlushDirty(false)
	editor.GetMsgList().Add(protos_world.NewTestMsgEditor().SetInnerMsg(protos_world.NewPlainMsgEditor().SetI32(32222)).SetField1(333).SetField2("strfield2"))
	incr = &protos_world.Test{}
	editor.GatherDirtyToPb(incr)
	assert.Equal(t, len(incr.MsgList), 1)
	mirror.MergeDirtyFromPb(incr)
	assert.Equal(t, len(mirror.GetMsgList().Values()), len(editor.GetMsgList().Values()))

	editor.FlushDirty(false)
	editor.GetMsgList().Get(2).GetInnerMsg().SetI32(64).SetB(true).SetStr("inner")
	editor.GetMsgList().Get(2).SetField2("newfield2").SetField1(math.MinInt32)
	incr = &protos_world.Test{}
	editor.GatherDirtyToPb(incr)
	assert.Equal(t, len(incr.MsgList), 1)
	assert.Equal(t, *incr.MsgList[0].Field1, int32(math.MinInt32))
	assert.Equal(t, *incr.MsgList[0].Field2, "newfield2")
	assert.Equal(t, (*incr.MsgList[0].InnerMsg).GetI32(), int32(64))
	assert.Equal(t, (*incr.MsgList[0].InnerMsg).GetStr(), "inner")
	mirror.MergeDirtyFromPb(incr)
	assert.Equal(t, len(mirror.GetMsgList().Values()), len(editor.GetMsgList().Values()))
	assert.Equal(t, mirror.GetMsgList().Get(2).GetInnerMsg().GetI32(), editor.GetMsgList().Get(2).GetInnerMsg().GetI32())
	assert.Equal(t, mirror.GetMsgList().Get(2).GetInnerMsg().GetB(), editor.GetMsgList().Get(2).GetInnerMsg().GetB())
	assert.Equal(t, mirror.GetMsgList().Get(2).GetField1(), editor.GetMsgList().Get(2).GetField1())
	assert.Equal(t, mirror.GetMsgList().Get(2).GetField2(), editor.GetMsgList().Get(2).GetField2())

	editor.FlushDirty(false)
	editor.GetMsgList().Set(1, protos_world.NewTestMsgEditor().SetField1(999999))
	incr = &protos_world.Test{}
	editor.GatherDirtyToPb(incr)
	assert.Equal(t, len(incr.MsgList), 1)
	assert.Equal(t, *incr.MsgList[0].Field1, int32(999999))
	mirror.MergeDirtyFromPb(incr)
	assert.Equal(t, len(mirror.GetMsgList().Values()), len(editor.GetMsgList().Values()))
	assert.Equal(t, mirror.GetMsgList().Get(1).GetField1(), editor.GetMsgList().Get(1).GetField1())
	t.Log(protojson.MarshalOptions{Indent: ""}.Format(incr), protojson.MarshalOptions{Indent: ""}.Format(pbeditor.UnmarshalMutState(incr.ProtoReflect().GetUnknown())))
	t.Log(mirror.GetMsgList().Get(1).GetField2())

	editor.FlushDirty(false)
	editor.GetMsgList().RemoveByIdx(1)
	incr = &protos_world.Test{}
	editor.GatherDirtyToPb(incr)
	assert.Equal(t, len(incr.MsgList), 0)
	mirror.MergeDirtyFromPb(incr)
	assert.Equal(t, len(mirror.GetMsgList().Values()), len(editor.GetMsgList().Values()))
	assert.Equal(t, mirror.GetMsgList().Get(1).GetField1(), editor.GetMsgList().Get(1).GetField1())
	t.Log(mirror.GetMsgList().Get(1).GetField2())

	editor.FlushDirty(false)
	editor.GetMsgList().Clear()
	incr = &protos_world.Test{}
	editor.GatherDirtyToPb(incr)
	assert.Equal(t, len(incr.MsgList), 0)
	mirror.MergeDirtyFromPb(incr)
	assert.Equal(t, len(mirror.GetMsgList().Values()), len(editor.GetMsgList().Values()))
	assert.Equal(t, len(mirror.GetMsgList().Values()), 0)

	//map
	editor.FlushDirty(false)
	editor.GetMsgMap().Get(111).SetField1(555).SetField2("map_str_v").GetInnerMsg().SetI32(32)
	incr = &protos_world.Test{}
	editor.GatherDirtyToPb(incr)
	assert.Equal(t, len(incr.MsgMap), 1)
	assert.Equal(t, incr.MsgMap[111].GetField1(), int32(555))
	assert.Equal(t, incr.MsgMap[111].GetField2(), "map_str_v")
	assert.Equal(t, incr.MsgMap[111].GetInnerMsg().GetI32(), int32(32))
	mirror.MergeDirtyFromPb(incr)
	assert.Equal(t, mirror.GetMsgMap().Get(111).GetField1(), int32(555))
	assert.Equal(t, mirror.GetMsgMap().Get(111).GetField2(), "map_str_v")
	assert.Equal(t, mirror.GetMsgMap().Get(111).GetInnerMsg().GetI32(), int32(32))

	editor.FlushDirty(false)
	var msg = protos_world.NewTestMsgEditor().SetField1(555).SetField2("map_str_v")
	msg.GetInnerMsg().SetI32(32)
	editor.GetMsgMap().Put(999, msg)
	incr = &protos_world.Test{}
	editor.GatherDirtyToPb(incr)
	assert.Equal(t, len(incr.MsgMap), 1)
	assert.Equal(t, incr.MsgMap[999].GetField1(), int32(555))
	assert.Equal(t, incr.MsgMap[999].GetField2(), "map_str_v")
	assert.Equal(t, incr.MsgMap[999].GetInnerMsg().GetI32(), int32(32))
	mirror.MergeDirtyFromPb(incr)
	assert.Equal(t, mirror.GetMsgMap().Get(999).GetField1(), int32(555))
	assert.Equal(t, mirror.GetMsgMap().Get(999).GetField2(), "map_str_v")
	assert.Equal(t, mirror.GetMsgMap().Get(999).GetInnerMsg().GetI32(), int32(32))

	editor.FlushDirty(false)
	editor.GetMsgMap().Remove(333)
	incr = &protos_world.Test{}
	editor.GatherDirtyToPb(incr)
	mirror.MergeDirtyFromPb(incr)
	assert.Equal(t, len(mirror.GetMsgMap().Entries()), 2)
	mutState := pbeditor.UnmarshalMutState(incr.ProtoReflect().GetUnknown())
	assert.Equal(t, mutState.GetMap()[59].I32List, []int32{333})

	editor.FlushDirty(false)
	editor.GetMsgMap().Put(1024, protos_world.NewTestMsgEditor().SetField1(1024).SetField2("newv"))
	incr = &protos_world.Test{}
	editor.GatherDirtyToPb(incr)
	mirror.MergeDirtyFromPb(incr)
	assert.Equal(t, len(mirror.GetMsgMap().Entries()), 3)
	mutState = pbeditor.UnmarshalMutState(incr.ProtoReflect().GetUnknown())
	assert.Equal(t, slices.Collect(maps.Keys(incr.MsgMap)), []int32{1024})
	assert.Equal(t, mirror.GetMsgMap().Get(1024).GetField1(), int32(1024))
	assert.Equal(t, mirror.GetMsgMap().Get(1024).GetField2(), "newv")

	editor.FlushDirty(false)
	editor.GetMsgMap().Clear()
	incr = &protos_world.Test{}
	editor.GatherDirtyToPb(incr)
	mirror.MergeDirtyFromPb(incr)
	assert.Equal(t, len(mirror.GetMsgMap().Entries()), 0)
	mutState = pbeditor.UnmarshalMutState(incr.ProtoReflect().GetUnknown())
	assert.Equal(t, slices.Sorted(slices.Values(mutState.GetMap()[59].I32List)), []int32{111, 999, 1024})

	editor.GetStrF32Map().Put("s1", math.MaxFloat32)
	editor.GetStrF32Map().Put("s2", -2234.3212)
	incr = &protos_world.Test{}
	editor.GatherDirtyToPb(incr)
	mirror.MergeDirtyFromPb(incr)
	assert.Equal(t, len(mirror.GetStrF32Map().Entries()), 4)
	assert.True(t, slices.Contains(slices.Collect(maps.Keys(mirror.GetStrF32Map().Entries())), "s1"))
	assert.True(t, slices.Contains(slices.Collect(maps.Keys(mirror.GetStrF32Map().Entries())), "s2"))
	f, b := mirror.GetStrF32Map().Get("s1")
	assert.True(t, b)
	assert.Equal(t, f, float32(math.MaxFloat32))
	f, b = mirror.GetStrF32Map().Get("s2")
	assert.True(t, b)
	assert.Equal(t, f, float32(-2234.3212))

	editor.GetStrF32Map().Put("s1", 99.99)
	editor.GetStrF32Map().Remove("s2")
	incr = &protos_world.Test{}
	editor.GatherDirtyToPb(incr)
	mirror.MergeDirtyFromPb(incr)
	assert.Equal(t, len(mirror.GetStrF32Map().Entries()), 3)
	assert.True(t, slices.Contains(slices.Collect(maps.Keys(mirror.GetStrF32Map().Entries())), "s1"))
	assert.False(t, slices.Contains(slices.Collect(maps.Keys(mirror.GetStrF32Map().Entries())), "s2"))
	f, b = mirror.GetStrF32Map().Get("s1")
	assert.True(t, b)
	assert.Equal(t, f, float32(99.99))
	f, b = mirror.GetStrF32Map().Get("s2")
	assert.False(t, b)

}
