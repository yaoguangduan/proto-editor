package protos_animal

import "github.com/yaoguangduan/proto-editor/pbeditor"
import "iter"
import "math"
import "slices"
import "github.com/yaoguangduan/proto-editor/protodef"
import "google.golang.org/protobuf/encoding/protowire"
import "google.golang.org/protobuf/proto"

type I64AnimalEditorMap struct {
	parent      pbeditor.Editor
	keyInParent any
	value       map[int64]*AnimalEditor
	dirtied     map[int64]struct{}
	deleted     map[int64]struct{}
}

func NewI64AnimalEditorMap() *I64AnimalEditorMap {
	return &I64AnimalEditorMap{value: make(map[int64]*AnimalEditor), dirtied: make(map[int64]struct{}), deleted: make(map[int64]struct{})}
}

func (ms *I64AnimalEditorMap) SetParent(editor pbeditor.Editor, key any) {
	ms.parent = editor
	ms.keyInParent = key
}

func (ms *I64AnimalEditorMap) SetParentDirty() {
	if ms.parent != nil {
		ms.parent.SetDirty(ms.keyInParent, true, ms)
	}
}

func (ms *I64AnimalEditorMap) Entries() map[int64]*AnimalEditor {
	return ms.value
}

func (ms *I64AnimalEditorMap) SetDirty(key any, dirty bool, editor pbeditor.Editor) {
	if dirty {
		ms.dirtied[key.(int64)] = struct{}{}
		ms.SetParentDirty()
	} else {
		delete(ms.dirtied, key.(int64))
	}

}

func (ms *I64AnimalEditorMap) FlushDirty(dirty bool) {
	if !dirty {
		ms.dirtied = map[int64]struct{}{}
		ms.deleted = map[int64]struct{}{}
	}
	for _, v := range ms.value {
		v.FlushDirty(dirty)
	}
}

func (ms *I64AnimalEditorMap) Put(k int64, v *AnimalEditor) *I64AnimalEditorMap {
	v.SetParent(ms, k)
	old, exist := ms.value[k]
	if exist {
		if old != v {
			old.SetParent(nil, nil)
			ms.deleted[k] = struct{}{}
		}
	}
	v.FlushDirty(true)
	ms.value[k] = v
	ms.dirtied[k] = struct{}{}
	ms.SetParentDirty()
	return ms
}

func (ms *I64AnimalEditorMap) PutAll(kv map[int64]*AnimalEditor) *I64AnimalEditorMap {
	for k, v := range kv {
		ms.Put(k, v)
	}
	return ms
}

func (ms *I64AnimalEditorMap) Len() int {
	return len(ms.value)
}

func (ms *I64AnimalEditorMap) Clear() *I64AnimalEditorMap {
	if ms.Len() <= 0 {
		return ms
	}
	for k, v := range ms.value {
		v.SetParent(nil, nil)
		ms.deleted[k] = struct{}{}
	}
	ms.value = map[int64]*AnimalEditor{}
	ms.dirtied = map[int64]struct{}{}
	ms.SetParentDirty()
	return ms
}

func (ms *I64AnimalEditorMap) Get(k int64) *AnimalEditor {
	v := ms.value[k]
	return v
}

func (ms *I64AnimalEditorMap) Remove(k int64) *AnimalEditor {
	v, exist := ms.value[k]
	if !exist {
		return v
	}
	v.SetParent(nil, nil)
	delete(ms.value, k)
	delete(ms.dirtied, k)
	ms.deleted[k] = struct{}{}
	ms.SetParentDirty()
	return v
}

func (ms *I64AnimalEditorMap) RemoveAll(k []int64) {
	for _, kk := range k {
		ms.Remove(kk)
	}
}

func (ms *I64AnimalEditorMap) Iter() iter.Seq2[int64, *AnimalEditor] {
	return func(yield func(int64, *AnimalEditor) bool) {
		for k, v := range ms.value {
			if !yield(k, v) {
				return
			}
		}
	}
}

func (ms *I64AnimalEditorMap) Dirtied() map[int64]struct{} {
	return ms.dirtied
}

func (ms *I64AnimalEditorMap) Deleted() map[int64]struct{} {
	return ms.deleted
}

func (ms *I64AnimalEditorMap) ContainDirtied(kk int64) bool {
	for k := range ms.dirtied {
		if k == kk {
			return true
		}
	}
	return false
}

func (ms *I64AnimalEditorMap) ContainDeleted(kk int64) bool {
	for k := range ms.deleted {
		if k == kk {
			return true
		}
	}
	return false
}

func (ms *I64AnimalEditorMap) MergeDirtyFromPb(delList []int64, dirties map[int64]*Animal) {
	ms.RemoveAll(delList)
	for k, v := range dirties {
		var tmp = ms.Get(k)
		if tmp == nil {
			tmp = NewAnimalEditor()
			tmp.MergeDirtyFromPb(v)
			ms.Put(k, tmp)
		} else {
			tmp.MergeDirtyFromPb(v)
		}
	}
}

type AnimalEditor struct {
	type_        AnimalType
	count        uint32
	foods        *pbeditor.StrList
	_dirties     map[int]struct{}
	_parent      pbeditor.Editor
	_keyInParent any
}

func NewAnimalEditor() *AnimalEditor {
	return &AnimalEditor{
		_dirties: make(map[int]struct{}),
	}
}
func (x *AnimalEditor) SetDirty(key any, dirty bool, sync pbeditor.Editor) {

	if dirty {
		x._dirties[key.(int)] = struct{}{}
		x.SetParentDirty()
	} else {
		delete(x._dirties, key.(int))
	}
}
func (x *AnimalEditor) SetParentDirty() {
	if x._parent != nil {
		x._parent.SetDirty(x._keyInParent, true, x)
	}
}
func (x *AnimalEditor) SetParent(editor pbeditor.Editor, key any) {
	x._parent = editor
	x._keyInParent = key
}
func (x *AnimalEditor) FlushDirty(dirty bool) {
	if dirty || x.isTypeDirty() {
		x.setTypeDirty(dirty, true)
	}
	if dirty || x.isCountDirty() {
		x.setCountDirty(dirty, true)
	}
	if dirty || x.isFoodsDirty() {
		x.setFoodsDirty(dirty, true)
	}
}
func (x *AnimalEditor) setTypeDirty(dirty bool, recur bool) {
	x.SetDirty(2, dirty, x)
}
func (x *AnimalEditor) isTypeDirty() bool {
	_, ok := x._dirties[2]
	return ok
}
func (x *AnimalEditor) setCountDirty(dirty bool, recur bool) {
	x.SetDirty(3, dirty, x)
}
func (x *AnimalEditor) isCountDirty() bool {
	_, ok := x._dirties[3]
	return ok
}
func (x *AnimalEditor) setFoodsDirty(dirty bool, recur bool) {
	x.SetDirty(4, dirty, x)
	if recur && x.foods != nil {
		x.foods.FlushDirty(dirty)
	}
}
func (x *AnimalEditor) isFoodsDirty() bool {
	_, ok := x._dirties[4]
	return ok
}
func (x *AnimalEditor) Key() any {
	return x._keyInParent
}
func (x *AnimalEditor) SetKey(v any) {
	x._keyInParent = v
}
func (x *AnimalEditor) Clear() *AnimalEditor {
	x.SetType(AnimalType_AT_Dog)
	x.SetCount(0)
	if x.foods != nil {
		x.foods.Clear()
	}
	return x
}
func (x *AnimalEditor) CopyTo(r *Animal) *AnimalEditor {
	r.SetType(x.type_)
	r.SetCount(x.count)
	if x.foods != nil && x.foods.Len() > 0 {
		r.SetFoods(x.foods.Values())
	}
	return x
}
func (x *AnimalEditor) CopyFrom(r *Animal) *AnimalEditor {
	if r.Type != nil {
		x.SetType(*r.Type)
	}
	if r.Count != nil {
		x.SetCount(*r.Count)
	}
	if r.Foods != nil {
		x.GetFoods().AddAll(r.Foods)
	}
	return x
}
func (x *AnimalEditor) MergeDirtyFromPb(r *Animal) {
	var mutState = pbeditor.UnmarshalMutState(r.ProtoReflect().GetUnknown())
	if r.Type != nil {
		x.SetType(*r.Type)
	}
	if r.Count != nil {
		x.SetCount(*r.Count)
	}
	x.GetFoods().MergeDirtyFromPb(mutState.ListOpLog(4), r.Foods)
}
func (x *AnimalEditor) MergeDirtyFromBytes(buf []byte) *AnimalEditor {
	fds := pbeditor.ToRawMessage(buf)
	for _, rawFields := range fds.RawFields {
		if len(rawFields) <= 0 {
			continue
		}
		rawF := rawFields[0]
		switch rawF.Number {
		case 2:
			x.SetType(AnimalType(int32(rawF.Variant)))
		case 3:
			x.SetCount(uint32(rawF.Variant))
		case 4:
			sl := make([]string, 0)
			for _, field := range rawFields {
				sl = append(sl, pbeditor.Bys2Str(field.Bytes))
			}
			x.GetFoods().MergeDirtyFromPb(fds.MutState.ListOpLog(rawF.Number), sl)
		}
	}
	return x
}
func (x *AnimalEditor) GatherDirtyToPb(r *Animal) {
	mutState := pbeditor.UnmarshalMutState(r.unknownFields)
	if x.isTypeDirty() {
		r.SetType(x.type_)
	}
	if x.isCountDirty() {
		r.SetCount(x.count)
	}
	if x.isFoodsDirty() {
		ori := r.Foods
		lenOfFoods := len(ori)
		logsOfFoods := x.foods.ChangeLog()
		dirtiesOfFoods := x.foods.MonitorDirtyList()
		for _, val := range dirtiesOfFoods {
			r.Foods = append(r.Foods, val)
		}
		for _, log := range logsOfFoods {
			log.Pos = int32(lenOfFoods) + log.Pos
			mutState.ListAppendLog(4, log)
		}
	}
	r.ProtoReflect().SetUnknown(mutState.MarshalToFieldBytes())
}
func (x *AnimalEditor) GatherDirtyToBytes() []byte {
	mutState := pbeditor.CreateMutStateInfo()
	var buf []byte
	if x.isTypeDirty() {
		buf = protowire.AppendTag(buf, 2, 0)
		buf = protowire.AppendVarint(buf, uint64(x.type_))
	}
	if x.isCountDirty() {
		buf = protowire.AppendTag(buf, 3, 0)
		buf = protowire.AppendVarint(buf, uint64(x.count))
	}
	if x.isFoodsDirty() {
		for _, log := range x.foods.ChangeLog() {
			mutState.ListAppendLog(4, log)
		}
		for _, s := range x.foods.MonitorDirtyList() {
			buf = protowire.AppendTag(buf, 4, protowire.BytesType)
			buf = protowire.AppendString(buf, s)
		}
	}
	buf = append(buf, mutState.MarshalToFieldBytes()...)
	return buf
}
func (x *AnimalEditor) GetType() AnimalType {
	return x.type_
}
func (x *AnimalEditor) SetType(v AnimalType) *AnimalEditor {
	if x.type_ == v {
		return x
	}
	x.type_ = v
	x.setTypeDirty(true, false)
	return x
}
func (x *AnimalEditor) GetCount() uint32 {
	return x.count
}
func (x *AnimalEditor) SetCount(v uint32) *AnimalEditor {
	if x.count == v {
		return x
	}
	x.count = v
	x.setCountDirty(true, false)
	return x
}
func (x *AnimalEditor) GetFoods() *pbeditor.StrList {
	if x.foods == nil {
		x.foods = pbeditor.NewStrList()
		x.foods.SetParent(x, 4)
	}
	return x.foods
}
func (x *AnimalEditor) AddFoods(v string) *AnimalEditor {
	if x.foods == nil {
		x.foods = pbeditor.NewStrList()
		x.foods.SetParent(x, 4)
	}
	x.foods.Add(v)
	return x
}
func (xs *Animal) SetType(v AnimalType) {
	xs.Type = &v
}
func (xs *Animal) SetCount(v uint32) {
	xs.Count = &v
}
func (xs *Animal) SetFoods(v []string) {
	xs.Foods = v
}
func (xs *Animal) AddFoods(v string) {
	xs.Foods = append(xs.Foods, v)
}

func unused() {
	_ = math.MinInt
	_ = slices.Min([]int32{})
	_ = iter.Seq[int32](func(yield func(int32) bool) {

	})
	_ = protodef.ListOpType_Delete
	_ = pbeditor.U32List{}
	_ = protowire.Number(1)
	_ = proto.Error
}
