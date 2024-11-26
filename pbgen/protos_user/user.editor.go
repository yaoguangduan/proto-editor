package protos_user

import (
	protos_animal "github.com/yaoguangduan/proto-editor/pbgen/protos_animal"
)

import "github.com/yaoguangduan/proto-editor/pbeditor"
import "iter"
import "math"
import "slices"
import "github.com/yaoguangduan/proto-editor/protodef"
import "google.golang.org/protobuf/encoding/protowire"
import "google.golang.org/protobuf/proto"

type StrUserEditorMap struct {
	parent      pbeditor.Editor
	keyInParent any
	value       map[string]*UserEditor
	dirtied     map[string]struct{}
	deleted     map[string]struct{}
}

func NewStrUserEditorMap() *StrUserEditorMap {
	return &StrUserEditorMap{value: make(map[string]*UserEditor), dirtied: make(map[string]struct{}), deleted: make(map[string]struct{})}
}

func (ms *StrUserEditorMap) SetParent(editor pbeditor.Editor, key any) {
	ms.parent = editor
	ms.keyInParent = key
}

func (ms *StrUserEditorMap) SetParentDirty() {
	if ms.parent != nil {
		ms.parent.SetDirty(ms.keyInParent, true, ms)
	}
}

func (ms *StrUserEditorMap) Entries() map[string]*UserEditor {
	return ms.value
}

func (ms *StrUserEditorMap) SetDirty(key any, dirty bool, editor pbeditor.Editor) {
	if dirty {
		ms.dirtied[key.(string)] = struct{}{}
		ms.SetParentDirty()
	} else {
		delete(ms.dirtied, key.(string))
	}

}

func (ms *StrUserEditorMap) FlushDirty(dirty bool) {
	if !dirty {
		ms.dirtied = map[string]struct{}{}
		ms.deleted = map[string]struct{}{}
	}
	for _, v := range ms.value {
		v.FlushDirty(dirty)
	}
}

func (ms *StrUserEditorMap) Put(k string, v *UserEditor) *StrUserEditorMap {
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

func (ms *StrUserEditorMap) PutAll(kv map[string]*UserEditor) *StrUserEditorMap {
	for k, v := range kv {
		ms.Put(k, v)
	}
	return ms
}

func (ms *StrUserEditorMap) Len() int {
	return len(ms.value)
}

func (ms *StrUserEditorMap) Clear() *StrUserEditorMap {
	if ms.Len() <= 0 {
		return ms
	}
	for k, v := range ms.value {
		v.SetParent(nil, nil)
		ms.deleted[k] = struct{}{}
	}
	ms.value = map[string]*UserEditor{}
	ms.dirtied = map[string]struct{}{}
	ms.SetParentDirty()
	return ms
}

func (ms *StrUserEditorMap) Get(k string) *UserEditor {
	v := ms.value[k]
	return v
}

func (ms *StrUserEditorMap) Remove(k string) *UserEditor {
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

func (ms *StrUserEditorMap) RemoveAll(k []string) {
	for _, kk := range k {
		ms.Remove(kk)
	}
}

func (ms *StrUserEditorMap) Iter() iter.Seq2[string, *UserEditor] {
	return func(yield func(string, *UserEditor) bool) {
		for k, v := range ms.value {
			if !yield(k, v) {
				return
			}
		}
	}
}

func (ms *StrUserEditorMap) Dirtied() map[string]struct{} {
	return ms.dirtied
}

func (ms *StrUserEditorMap) Deleted() map[string]struct{} {
	return ms.deleted
}

func (ms *StrUserEditorMap) ContainDirtied(kk string) bool {
	for k := range ms.dirtied {
		if k == kk {
			return true
		}
	}
	return false
}

func (ms *StrUserEditorMap) ContainDeleted(kk string) bool {
	for k := range ms.deleted {
		if k == kk {
			return true
		}
	}
	return false
}

func (ms *StrUserEditorMap) MergeDirtyFromPb(delList []string, dirties map[string]*User) {
	ms.RemoveAll(delList)
	for k, v := range dirties {
		var tmp = ms.Get(k)
		if tmp == nil {
			tmp = NewUserEditor()
			tmp.MergeDirtyFromPb(v)
			ms.Put(k, tmp)
		} else {
			tmp.MergeDirtyFromPb(v)
		}
	}
}

type UserEditor struct {
	name         string
	age          int32
	pet          *protos_animal.I64AnimalEditorMap
	friends      *FriendEditorList
	tempChange   *pbeditor.F32List
	_dirties     map[int]struct{}
	_parent      pbeditor.Editor
	_keyInParent any
}

func NewUserEditor() *UserEditor {
	return &UserEditor{
		_dirties: make(map[int]struct{}),
	}
}
func (x *UserEditor) SetDirty(key any, dirty bool, sync pbeditor.Editor) {

	if dirty {
		x._dirties[key.(int)] = struct{}{}
		x.SetParentDirty()
	} else {
		delete(x._dirties, key.(int))
	}
}
func (x *UserEditor) SetParentDirty() {
	if x._parent != nil {
		x._parent.SetDirty(x._keyInParent, true, x)
	}
}
func (x *UserEditor) SetParent(editor pbeditor.Editor, key any) {
	x._parent = editor
	x._keyInParent = key
}
func (x *UserEditor) FlushDirty(dirty bool) {
	if dirty || x.isNameDirty() {
		x.setNameDirty(dirty, true)
	}
	if dirty || x.isAgeDirty() {
		x.setAgeDirty(dirty, true)
	}
	if dirty || x.isPetDirty() {
		x.setPetDirty(dirty, true)
	}
	if dirty || x.isFriendsDirty() {
		x.setFriendsDirty(dirty, true)
	}
	if dirty || x.isTempChangeDirty() {
		x.setTempChangeDirty(dirty, true)
	}
}
func (x *UserEditor) setNameDirty(dirty bool, recur bool) {
	x.SetDirty(1, dirty, x)
}
func (x *UserEditor) isNameDirty() bool {
	_, ok := x._dirties[1]
	return ok
}
func (x *UserEditor) setAgeDirty(dirty bool, recur bool) {
	x.SetDirty(2, dirty, x)
}
func (x *UserEditor) isAgeDirty() bool {
	_, ok := x._dirties[2]
	return ok
}
func (x *UserEditor) setPetDirty(dirty bool, recur bool) {
	x.SetDirty(3, dirty, x)
	if recur && x.pet != nil {
		x.pet.FlushDirty(dirty)
	}
}
func (x *UserEditor) isPetDirty() bool {
	_, ok := x._dirties[3]
	return ok
}
func (x *UserEditor) setFriendsDirty(dirty bool, recur bool) {
	x.SetDirty(4, dirty, x)
	if recur && x.friends != nil {
		x.friends.FlushDirty(dirty)
	}
}
func (x *UserEditor) isFriendsDirty() bool {
	_, ok := x._dirties[4]
	return ok
}
func (x *UserEditor) setTempChangeDirty(dirty bool, recur bool) {
	x.SetDirty(5, dirty, x)
	if recur && x.tempChange != nil {
		x.tempChange.FlushDirty(dirty)
	}
}
func (x *UserEditor) isTempChangeDirty() bool {
	_, ok := x._dirties[5]
	return ok
}
func (x *UserEditor) Key() any {
	return x._keyInParent
}
func (x *UserEditor) SetKey(v any) {
	x._keyInParent = v
}
func (x *UserEditor) Clear() *UserEditor {
	x.SetName("")
	x.SetAge(0)
	if x.pet != nil {
		x.pet.Clear()
	}
	if x.friends != nil {
		x.friends.Clear()
	}
	if x.tempChange != nil {
		x.tempChange.Clear()
	}
	return x
}
func (x *UserEditor) CopyTo(r *User) *UserEditor {
	r.SetName(x.name)
	r.SetAge(x.age)
	if x.pet != nil && x.pet.Len() > 0 {
		tmp := make(map[int64]*protos_animal.Animal)
		for k, v := range x.pet.Iter() {
			tmpV := &protos_animal.Animal{}
			v.CopyTo(tmpV)
			tmp[k] = tmpV
		}
		r.SetPet(tmp)
	}
	if x.friends != nil && x.friends.Len() > 0 {

		for _, v := range x.friends.Iter() {
			if v == nil {
				r.AddFriends(nil)
			} else {
				tmp := &Friend{}
				v.CopyTo(tmp)
				r.AddFriends(tmp)
			}
		}

	}
	if x.tempChange != nil && x.tempChange.Len() > 0 {
		r.SetTempChange(x.tempChange.Values())
	}
	return x
}
func (x *UserEditor) CopyFrom(r *User) *UserEditor {
	if r.Name != nil {
		x.SetName(*r.Name)
	}
	if r.Age != nil {
		x.SetAge(*r.Age)
	}
	for k, v := range r.Pet {
		if v != nil {
			vv := protos_animal.NewAnimalEditor()
			vv.CopyFrom(v)
			x.GetPet().Put(k, vv)
		}
	}
	if r.Friends != nil {

		for _, v := range r.Friends {
			if v == nil {
				x.GetFriends().Add(nil)
			} else {
				tmp := NewFriendEditor()
				tmp.CopyFrom(v)
				x.GetFriends().Add(tmp)
			}
		}

	}
	if r.TempChange != nil {
		x.GetTempChange().AddAll(r.TempChange)
	}
	return x
}
func (x *UserEditor) MergeDirtyFromPb(r *User) {
	var mutState = pbeditor.UnmarshalMutState(r.ProtoReflect().GetUnknown())
	if r.Name != nil {
		x.SetName(*r.Name)
	}
	if r.Age != nil {
		x.SetAge(*r.Age)
	}
	x.GetPet().MergeDirtyFromPb(mutState.I64List(3), r.Pet)
	x.GetFriends().MergeDirtyFromPb(mutState.ListOpLog(4), r.Friends)
	x.GetTempChange().MergeDirtyFromPb(mutState.ListOpLog(5), r.TempChange)
}
func (x *UserEditor) MergeDirtyFromBytes(buf []byte) *UserEditor {
	fds := pbeditor.ToRawMessage(buf)
	if x.pet != nil {
		x.pet.RemoveAll(fds.MutState.I64List(3))
	}
	for _, rawFields := range fds.RawFields {
		if len(rawFields) <= 0 {
			continue
		}
		rawF := rawFields[0]
		switch rawF.Number {
		case 1:
			x.SetName(pbeditor.Bys2Str(rawF.Bytes))
		case 2:
			x.SetAge(int32(rawF.Variant))
		case 3:
			r := make(map[int64]*protos_animal.Animal)
			for _, field := range rawFields {
				tmp := pbeditor.ParseMap(field.Bytes)
				k := tmp.KeyI64()
				v := &protos_animal.Animal{}
				proto.Unmarshal(field.Bytes, v)
				r[k] = v
			}
			x.GetPet().MergeDirtyFromPb(fds.MutState.I64List(rawF.Number), r)
		case 4:
			sl := make([]*Friend, 0)
			for _, field := range rawFields {
				tmp := &Friend{}
				proto.Unmarshal(field.Bytes, tmp)
				sl = append(sl, tmp)
			}
			x.GetFriends().MergeDirtyFromPb(fds.MutState.ListOpLog(rawF.Number), sl)
		case 5:
			tmp := rawF.Bytes
			dl := make([]float32, 0)
			for len(tmp) > 0 {
				val, n := protowire.ConsumeFixed32(tmp)
				if n < 0 {
					panic(n)
				}
				tmp = tmp[n:]
				dl = append(dl, math.Float32frombits(val))
			}
			x.GetTempChange().MergeDirtyFromPb(fds.MutState.ListOpLog(rawF.Number), dl)
		}
	}
	return x
}
func (x *UserEditor) GatherDirtyToPb(r *User) {
	mutState := pbeditor.UnmarshalMutState(r.unknownFields)
	if x.isNameDirty() {
		r.SetName(x.name)
	}
	if x.isAgeDirty() {
		r.SetAge(x.age)
	}
	if x.isPetDirty() {
		updated := make([]int64, 0)
		if r.Pet != nil {
			for k, v := range r.Pet {
				if x.pet.ContainDeleted(k) {
					delete(r.Pet, k)
					continue
				}
				if x.pet.ContainDirtied(k) {
					updated = append(updated, k)
					tmp := x.pet.Get(k)
					if tmp == nil {
						r.Pet[k] = nil
						continue
					}
					if v == nil {
						v = &protos_animal.Animal{}
						tmp.GatherDirtyToPb(v)
						r.Pet[k] = v

					} else {
						tmp.GatherDirtyToPb(v)
					}
				}
			}
		} else {
			r.Pet = make(map[int64]*protos_animal.Animal)
		}
		for k := range x.pet.Dirtied() {
			if !slices.Contains(updated, k) {
				tmp := x.pet.Get(k)
				if tmp == nil {
					r.Pet[k] = nil
				} else {
					v := &protos_animal.Animal{}
					tmp.GatherDirtyToPb(v)
					r.Pet[k] = v
				}
			}
		}
		deleted := mutState.I64List(3)
		for k := range x.pet.Deleted() {
			if !slices.Contains(deleted, k) {
				mutState.AddI64(3, k)
			}
		}
	}
	if x.isFriendsDirty() {
		ori := r.Friends
		lenOfFriends := len(ori)
		logsOfFriends := x.friends.ChangeLog()
		dirtiesOfFriends := x.friends.MonitorDirtyList()
		for _, val := range dirtiesOfFriends {
			o := &Friend{}
			val.GatherDirtyToPb(o)
			r.Friends = append(r.Friends, o)
		}
		for _, log := range logsOfFriends {
			log.Pos = int32(lenOfFriends) + log.Pos
			mutState.ListAppendLog(4, log)
		}
	}
	if x.isTempChangeDirty() {
		ori := r.TempChange
		lenOfTempChange := len(ori)
		logsOfTempChange := x.tempChange.ChangeLog()
		dirtiesOfTempChange := x.tempChange.MonitorDirtyList()
		for _, val := range dirtiesOfTempChange {
			r.TempChange = append(r.TempChange, val)
		}
		for _, log := range logsOfTempChange {
			log.Pos = int32(lenOfTempChange) + log.Pos
			mutState.ListAppendLog(5, log)
		}
	}
	r.ProtoReflect().SetUnknown(mutState.MarshalToFieldBytes())
}
func (x *UserEditor) GatherDirtyToBytes() []byte {
	mutState := pbeditor.CreateMutStateInfo()
	var buf []byte
	if x.isNameDirty() {
		buf = protowire.AppendTag(buf, 1, 2)
		buf = protowire.AppendString(buf, x.name)
	}
	if x.isAgeDirty() {
		buf = protowire.AppendTag(buf, 2, 0)
		buf = protowire.AppendVarint(buf, uint64(x.age))
	}
	if x.isPetDirty() {
		if len(x.pet.Deleted()) > 0 {
			deleted := x.pet.Deleted()
			for v := range deleted {
				mutState.AddI64(3, v)
			}
		}
		if x.pet.Len() > 0 {
			for k, v := range x.pet.Iter() {
				if !x.pet.ContainDirtied(k) {
					continue
				}
				var tmp []byte
				tmp = protowire.AppendTag(tmp, 1, protowire.VarintType)
				tmp = protowire.AppendVarint(tmp, uint64(k))
				vb := v.GatherDirtyToBytes()
				tmp = protowire.AppendTag(tmp, 2, protowire.BytesType)
				tmp = protowire.AppendBytes(tmp, vb)
				buf = protowire.AppendTag(buf, 3, protowire.BytesType)
				buf = protowire.AppendBytes(buf, tmp)
			}
		}
	}
	if x.isFriendsDirty() {
		for _, log := range x.friends.ChangeLog() {
			mutState.ListAppendLog(4, log)
		}
		for _, e := range x.friends.MonitorDirtyList() {
			buf = protowire.AppendTag(buf, 4, protowire.BytesType)
			buf = protowire.AppendBytes(buf, e.GatherDirtyToBytes())
		}
	}
	if x.isTempChangeDirty() {
		for _, log := range x.tempChange.ChangeLog() {
			mutState.ListAppendLog(5, log)
		}
		tmp := make([]byte, 0)
		for _, f := range x.tempChange.MonitorDirtyList() {
			tmp = protowire.AppendFixed32(tmp, math.Float32bits(f))
		}
		buf = protowire.AppendTag(buf, 5, protowire.BytesType)
		buf = protowire.AppendBytes(buf, tmp)
	}
	buf = append(buf, mutState.MarshalToFieldBytes()...)
	return buf
}
func (x *UserEditor) GetName() string {
	return x.name
}
func (x *UserEditor) SetName(v string) *UserEditor {
	if x.name == v {
		return x
	}
	x.name = v
	x.setNameDirty(true, false)
	return x
}
func (x *UserEditor) GetAge() int32 {
	return x.age
}
func (x *UserEditor) SetAge(v int32) *UserEditor {
	if x.age == v {
		return x
	}
	x.age = v
	x.setAgeDirty(true, false)
	return x
}
func (x *UserEditor) GetPet() *protos_animal.I64AnimalEditorMap {
	if x.pet == nil {
		x.pet = protos_animal.NewI64AnimalEditorMap()
		x.pet.SetParent(x, 3)
	}
	return x.pet
}
func (x *UserEditor) GetFriends() *FriendEditorList {
	if x.friends == nil {
		x.friends = NewFriendEditorList()
		x.friends.SetParent(x, 4)
	}
	return x.friends
}
func (x *UserEditor) AddFriends(v *FriendEditor) *UserEditor {
	if x.friends == nil {
		x.friends = NewFriendEditorList()
		x.friends.SetParent(x, 4)
	}
	x.friends.Add(v)
	return x
}
func (x *UserEditor) GetTempChange() *pbeditor.F32List {
	if x.tempChange == nil {
		x.tempChange = pbeditor.NewF32List()
		x.tempChange.SetParent(x, 5)
	}
	return x.tempChange
}
func (x *UserEditor) AddTempChange(v float32) *UserEditor {
	if x.tempChange == nil {
		x.tempChange = pbeditor.NewF32List()
		x.tempChange.SetParent(x, 5)
	}
	x.tempChange.Add(v)
	return x
}
func (xs *User) SetName(v string) {
	xs.Name = &v
}
func (xs *User) SetAge(v int32) {
	xs.Age = &v
}
func (xs *User) SetPet(v map[int64]*protos_animal.Animal) {
	xs.Pet = v
}
func (xs *User) SetFriends(v []*Friend) {
	xs.Friends = v
}
func (xs *User) AddFriends(v *Friend) {
	xs.Friends = append(xs.Friends, v)
}
func (xs *User) SetTempChange(v []float32) {
	xs.TempChange = v
}
func (xs *User) AddTempChange(v float32) {
	xs.TempChange = append(xs.TempChange, v)
}

type FriendEditorList struct {
	parent   pbeditor.Editor
	idxInPar any
	value    []*FriendEditor
	monitor  []*FriendEditor
	change   []*protodef.ListOpLog
}

func NewFriendEditorList() *FriendEditorList {
	return &FriendEditorList{value: make([]*FriendEditor, 0)}
}

func (ms *FriendEditorList) SetParent(editor pbeditor.Editor, key any) {
	ms.parent = editor
	ms.idxInPar = key
}
func (ms *FriendEditorList) SetParentDirty() {
	if ms.parent != nil {
		ms.parent.SetDirty(ms.idxInPar, true, ms)
	}
}
func (ms *FriendEditorList) SetDirty(key any, dirty bool, sync pbeditor.Editor) {
	var pos = -1
	if key != nil {
		pos = key.(int)
	} else {
		pos = slices.Index(ms.value, sync.(*FriendEditor))
	}
	changed := ms.value[pos]
	idx := slices.Index(ms.monitor, changed)
	if idx == -1 {
		ms.monitor = append(ms.monitor, changed)
		idx = len(ms.monitor) - 1
	}
	ms.change = append(ms.change, &protodef.ListOpLog{
		OpType: protodef.ListOpType_Modify,
		Pos:    int32(pos),
		RefIdx: int32(idx),
	})
	ms.SetParentDirty()
}
func (ms *FriendEditorList) FlushDirty(dirty bool) {
	if !dirty {
		ms.monitor = nil
		ms.change = nil
	}
	for _, v := range ms.value {
		v.FlushDirty(dirty)
	}
}

func (ms *FriendEditorList) Add(t *FriendEditor) *FriendEditorList {
	t.SetParent(ms, nil)
	ms.value = append(ms.value, t)
	pos := len(ms.value) - 1
	idx := slices.Index(ms.monitor, t)
	if idx == -1 {
		ms.monitor = append(ms.monitor, t)
		idx = len(ms.monitor) - 1
	}
	ms.change = append(ms.change, &protodef.ListOpLog{
		OpType: protodef.ListOpType_Insert,
		Pos:    int32(pos),
		RefIdx: int32(idx),
	})
	ms.SetParentDirty()
	return ms
}

func (ms *FriendEditorList) AddAll(t []*FriendEditor) *FriendEditorList {
	for _, v := range t {
		ms.Add(v)
	}
	return ms
}

func (ms *FriendEditorList) Insert(pos int, v *FriendEditor) *FriendEditorList {
	v.SetParent(ms, nil)
	ms.value = slices.Insert(ms.value, pos, v)
	idx := slices.Index(ms.monitor, v)
	if idx == -1 {
		ms.monitor = append(ms.monitor, v)
		idx = len(ms.monitor) - 1
	}
	ms.change = append(ms.change, &protodef.ListOpLog{
		OpType: protodef.ListOpType_Insert,
		Pos:    int32(pos),
		RefIdx: int32(idx),
	})
	return ms
}
func (ms *FriendEditorList) Set(pos int, v *FriendEditor) *FriendEditorList {
	v.SetParent(ms, nil)
	old := ms.value[pos]
	if old != nil {
		if old != v {
			old.SetParent(nil, nil)
		}
	}
	ms.value[pos] = v
	ms.change = append(ms.change, &protodef.ListOpLog{
		OpType: protodef.ListOpType_Delete,
		Pos:    int32(pos),
	})
	ms.monitor = append(ms.monitor, v)
	ms.change = append(ms.change, &protodef.ListOpLog{
		OpType: protodef.ListOpType_Insert,
		Pos:    int32(pos),
		RefIdx: int32(len(ms.monitor) - 1),
	})
	v.SetParentDirty()
	return ms
}
func (ms *FriendEditorList) Iter() iter.Seq2[int, *FriendEditor] {
	return func(yield func(int, *FriendEditor) bool) {
		for i := range ms.value {
			if !yield(i, ms.value[i]) {
				return
			}
		}
	}
}

func (ms *FriendEditorList) Clear() {
	for _, v := range ms.value {
		v.SetParent(nil, nil)
	}
	ms.value = make([]*FriendEditor, 0)
	ms.monitor = make([]*FriendEditor, 0)
	ms.change = nil
	ms.change = append(ms.change, &protodef.ListOpLog{
		OpType: protodef.ListOpType_Clear,
		Pos:    0,
		RefIdx: 0,
	})
	ms.SetParentDirty()
}
func (ms *FriendEditorList) Len() int {
	return len(ms.value)
}

func (ms *FriendEditorList) Remove(v *FriendEditor) {
	ms.RemoveByIdx(slices.Index(ms.value, v))
}

func (ms *FriendEditorList) RemoveByIdx(pos int) {
	v := ms.value[pos]
	v.SetParent(nil, nil)
	ms.value = slices.Delete(ms.value, pos, pos+1)
	ms.change = append(ms.change, &protodef.ListOpLog{
		OpType: protodef.ListOpType_Delete,
		Pos:    int32(pos),
	})
	ms.SetParentDirty()
}

func (ms *FriendEditorList) ChangeLog() []*protodef.ListOpLog {
	return ms.change
}
func (ms *FriendEditorList) MonitorDirtyList() []*FriendEditor {
	return ms.monitor
}
func (ms *FriendEditorList) Values() []*FriendEditor {
	vv := make([]*FriendEditor, ms.Len())
	copy(vv, ms.value)
	return vv
}

func (ms *FriendEditorList) MergeDirtyFromPb(logs []*protodef.ListOpLog, dirties []*Friend) {
	for _, log := range logs {
		switch log.OpType {
		case protodef.ListOpType_Insert:
			v := dirties[log.RefIdx]
			tmp := NewFriendEditor().CopyFrom(v)
			if int(log.Pos) >= ms.Len() {
				ms.Add(tmp)
			} else {
				ms.Insert(int(log.Pos), tmp)
			}
		case protodef.ListOpType_Modify:
			ms.Get(log.Pos).MergeDirtyFromPb(dirties[log.RefIdx])
		case protodef.ListOpType_Delete:
			ms.RemoveByIdx(int(log.Pos))
		case protodef.ListOpType_Clear:
			ms.Clear()
		}
	}
}

func (ms *FriendEditorList) Get(pos int32) *FriendEditor {
	return ms.value[pos]
}

type FriendEditor struct {
	name         string
	relation     float32
	_dirties     map[int]struct{}
	_parent      pbeditor.Editor
	_keyInParent any
}

func NewFriendEditor() *FriendEditor {
	return &FriendEditor{
		_dirties: make(map[int]struct{}),
	}
}
func (x *FriendEditor) SetDirty(key any, dirty bool, sync pbeditor.Editor) {

	if dirty {
		x._dirties[key.(int)] = struct{}{}
		x.SetParentDirty()
	} else {
		delete(x._dirties, key.(int))
	}
}
func (x *FriendEditor) SetParentDirty() {
	if x._parent != nil {
		x._parent.SetDirty(x._keyInParent, true, x)
	}
}
func (x *FriendEditor) SetParent(editor pbeditor.Editor, key any) {
	x._parent = editor
	x._keyInParent = key
}
func (x *FriendEditor) FlushDirty(dirty bool) {
	if dirty || x.isNameDirty() {
		x.setNameDirty(dirty, true)
	}
	if dirty || x.isRelationDirty() {
		x.setRelationDirty(dirty, true)
	}
}
func (x *FriendEditor) setNameDirty(dirty bool, recur bool) {
	x.SetDirty(1, dirty, x)
}
func (x *FriendEditor) isNameDirty() bool {
	_, ok := x._dirties[1]
	return ok
}
func (x *FriendEditor) setRelationDirty(dirty bool, recur bool) {
	x.SetDirty(2, dirty, x)
}
func (x *FriendEditor) isRelationDirty() bool {
	_, ok := x._dirties[2]
	return ok
}
func (x *FriendEditor) Key() any {
	return x._keyInParent
}
func (x *FriendEditor) SetKey(v any) {
	x._keyInParent = v
}
func (x *FriendEditor) Clear() *FriendEditor {
	x.SetName("")
	x.SetRelation(0)
	return x
}
func (x *FriendEditor) CopyTo(r *Friend) *FriendEditor {
	r.SetName(x.name)
	r.SetRelation(x.relation)
	return x
}
func (x *FriendEditor) CopyFrom(r *Friend) *FriendEditor {
	if r.Name != nil {
		x.SetName(*r.Name)
	}
	if r.Relation != nil {
		x.SetRelation(*r.Relation)
	}
	return x
}
func (x *FriendEditor) MergeDirtyFromPb(r *Friend) {
	if r.Name != nil {
		x.SetName(*r.Name)
	}
	if r.Relation != nil {
		x.SetRelation(*r.Relation)
	}
}
func (x *FriendEditor) MergeDirtyFromBytes(buf []byte) *FriendEditor {
	fds := pbeditor.ToRawMessage(buf)
	for _, rawFields := range fds.RawFields {
		if len(rawFields) <= 0 {
			continue
		}
		rawF := rawFields[0]
		switch rawF.Number {
		case 1:
			x.SetName(pbeditor.Bys2Str(rawF.Bytes))
		case 2:
			x.SetRelation(rawF.F32)
		}
	}
	return x
}
func (x *FriendEditor) GatherDirtyToPb(r *Friend) {
	if x.isNameDirty() {
		r.SetName(x.name)
	}
	if x.isRelationDirty() {
		r.SetRelation(x.relation)
	}
}
func (x *FriendEditor) GatherDirtyToBytes() []byte {
	var buf []byte
	if x.isNameDirty() {
		buf = protowire.AppendTag(buf, 1, 2)
		buf = protowire.AppendString(buf, x.name)
	}
	if x.isRelationDirty() {
		buf = protowire.AppendTag(buf, 2, 5)
		buf = protowire.AppendFixed32(buf, math.Float32bits(x.relation))
	}
	return buf
}
func (x *FriendEditor) GetName() string {
	return x.name
}
func (x *FriendEditor) SetName(v string) *FriendEditor {
	if x.name == v {
		return x
	}
	x.name = v
	x.setNameDirty(true, false)
	return x
}
func (x *FriendEditor) GetRelation() float32 {
	return x.relation
}
func (x *FriendEditor) SetRelation(v float32) *FriendEditor {
	if x.relation == v {
		return x
	}
	x.relation = v
	x.setRelationDirty(true, false)
	return x
}
func (xs *Friend) SetName(v string) {
	xs.Name = &v
}
func (xs *Friend) SetRelation(v float32) {
	xs.Relation = &v
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
