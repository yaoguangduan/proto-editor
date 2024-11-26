package protos_enum

import "github.com/yaoguangduan/proto-editor/pbeditor"
import "iter"
import "math"
import "slices"
import "github.com/yaoguangduan/proto-editor/protodef"
import "google.golang.org/protobuf/encoding/protowire"
import "google.golang.org/protobuf/proto"

type StateEditorList struct {
	parent   pbeditor.Editor
	idxInPar any
	value    []State
	monitor  []State
	change   []*protodef.ListOpLog
}

func NewStateEditorList() *StateEditorList {
	return &StateEditorList{value: make([]State, 0)}
}

func (ms *StateEditorList) SetParent(editor pbeditor.Editor, key any) {
	ms.parent = editor
	ms.idxInPar = key
}
func (ms *StateEditorList) SetParentDirty() {
	if ms.parent != nil {
		ms.parent.SetDirty(ms.idxInPar, true, ms)
	}
}
func (ms *StateEditorList) SetDirty(key any, dirty bool, sync pbeditor.Editor) {
	var pos = key.(int)
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
func (ms *StateEditorList) FlushDirty(dirty bool) {
	if !dirty {
		ms.monitor = nil
		ms.change = nil
	}
}

func (ms *StateEditorList) Add(t State) *StateEditorList {
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

func (ms *StateEditorList) AddAll(t []State) *StateEditorList {
	for _, v := range t {
		ms.Add(v)
	}
	return ms
}

func (ms *StateEditorList) Insert(pos int, v State) *StateEditorList {
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
func (ms *StateEditorList) Set(pos int, v State) *StateEditorList {
	ms.value[pos] = v
	ms.SetDirty(pos, true, ms)
	return ms
}
func (ms *StateEditorList) Iter() iter.Seq2[int, State] {
	return func(yield func(int, State) bool) {
		for i := range ms.value {
			if !yield(i, ms.value[i]) {
				return
			}
		}
	}
}

func (ms *StateEditorList) Clear() {
	ms.value = make([]State, 0)
	ms.monitor = make([]State, 0)
	ms.change = nil
	ms.change = append(ms.change, &protodef.ListOpLog{
		OpType: protodef.ListOpType_Clear,
		Pos:    0,
		RefIdx: 0,
	})
	ms.SetParentDirty()
}
func (ms *StateEditorList) Len() int {
	return len(ms.value)
}

func (ms *StateEditorList) Remove(v State) {
	ms.RemoveByIdx(slices.Index(ms.value, v))
}

func (ms *StateEditorList) RemoveByIdx(pos int) {
	ms.value = slices.Delete(ms.value, pos, pos+1)
	ms.change = append(ms.change, &protodef.ListOpLog{
		OpType: protodef.ListOpType_Delete,
		Pos:    int32(pos),
	})
	ms.SetParentDirty()
}

func (ms *StateEditorList) ChangeLog() []*protodef.ListOpLog {
	return ms.change
}
func (ms *StateEditorList) MonitorDirtyList() []State {
	return ms.monitor
}
func (ms *StateEditorList) Values() []State {
	vv := make([]State, ms.Len())
	copy(vv, ms.value)
	return vv
}

func (ms *StateEditorList) MergeDirtyFromPb(logs []*protodef.ListOpLog, dirties []State) {
	for _, log := range logs {
		switch log.OpType {
		case protodef.ListOpType_Insert:
			v := dirties[log.RefIdx]
			if int(log.Pos) >= ms.Len() {
				ms.Add(v)
			} else {
				ms.Insert(int(log.Pos), v)
			}
		case protodef.ListOpType_Modify:
			ms.Set(int(log.Pos), dirties[log.RefIdx])
		case protodef.ListOpType_Delete:
			ms.RemoveByIdx(int(log.Pos))
		case protodef.ListOpType_Clear:
			ms.Clear()
		}
	}
}

func (ms *StateEditorList) Get(pos int32) State {
	return ms.value[pos]
}

type StrWorldTypeEditorMap struct {
	parent      pbeditor.Editor
	keyInParent any
	value       map[string]WorldType
	dirtied     map[string]struct{}
	deleted     map[string]struct{}
}

func NewStrWorldTypeEditorMap() *StrWorldTypeEditorMap {
	return &StrWorldTypeEditorMap{value: make(map[string]WorldType), dirtied: make(map[string]struct{}), deleted: make(map[string]struct{})}
}

func (ms *StrWorldTypeEditorMap) SetParent(editor pbeditor.Editor, key any) {
	ms.parent = editor
	ms.keyInParent = key
}

func (ms *StrWorldTypeEditorMap) SetParentDirty() {
	if ms.parent != nil {
		ms.parent.SetDirty(ms.keyInParent, true, ms)
	}
}

func (ms *StrWorldTypeEditorMap) Entries() map[string]WorldType {
	return ms.value
}

func (ms *StrWorldTypeEditorMap) SetDirty(key any, dirty bool, editor pbeditor.Editor) {
	if dirty {
		ms.dirtied[key.(string)] = struct{}{}
		ms.SetParentDirty()
	} else {
		delete(ms.dirtied, key.(string))
	}

}

func (ms *StrWorldTypeEditorMap) FlushDirty(dirty bool) {
	if !dirty {
		ms.dirtied = map[string]struct{}{}
		ms.deleted = map[string]struct{}{}
	}
}

func (ms *StrWorldTypeEditorMap) Put(k string, v WorldType) *StrWorldTypeEditorMap {
	old, exist := ms.value[k]
	if exist {
		if old != v {
			ms.deleted[k] = struct{}{}
		}
	}
	ms.value[k] = v
	ms.dirtied[k] = struct{}{}
	ms.SetParentDirty()
	return ms
}

func (ms *StrWorldTypeEditorMap) PutAll(kv map[string]WorldType) *StrWorldTypeEditorMap {
	for k, v := range kv {
		ms.Put(k, v)
	}
	return ms
}

func (ms *StrWorldTypeEditorMap) Len() int {
	return len(ms.value)
}

func (ms *StrWorldTypeEditorMap) Clear() *StrWorldTypeEditorMap {
	if ms.Len() <= 0 {
		return ms
	}
	for k, _ := range ms.value {
		ms.deleted[k] = struct{}{}
	}
	ms.value = map[string]WorldType{}
	ms.dirtied = map[string]struct{}{}
	ms.SetParentDirty()
	return ms
}

func (ms *StrWorldTypeEditorMap) Get(k string) (WorldType, bool) {
	v, ok := ms.value[k]
	return v, ok
}

func (ms *StrWorldTypeEditorMap) Remove(k string) (WorldType, bool) {
	v, exist := ms.value[k]
	if !exist {
		return v, exist
	}
	delete(ms.value, k)
	delete(ms.dirtied, k)
	ms.deleted[k] = struct{}{}
	ms.SetParentDirty()
	return v, exist
}

func (ms *StrWorldTypeEditorMap) RemoveAll(k []string) {
	for _, kk := range k {
		ms.Remove(kk)
	}
}

func (ms *StrWorldTypeEditorMap) Iter() iter.Seq2[string, WorldType] {
	return func(yield func(string, WorldType) bool) {
		for k, v := range ms.value {
			if !yield(k, v) {
				return
			}
		}
	}
}

func (ms *StrWorldTypeEditorMap) Dirtied() map[string]struct{} {
	return ms.dirtied
}

func (ms *StrWorldTypeEditorMap) Deleted() map[string]struct{} {
	return ms.deleted
}

func (ms *StrWorldTypeEditorMap) ContainDirtied(kk string) bool {
	for k := range ms.dirtied {
		if k == kk {
			return true
		}
	}
	return false
}

func (ms *StrWorldTypeEditorMap) ContainDeleted(kk string) bool {
	for k := range ms.deleted {
		if k == kk {
			return true
		}
	}
	return false
}

func (ms *StrWorldTypeEditorMap) MergeDirtyFromPb(delList []string, dirties map[string]WorldType) {
	ms.RemoveAll(delList)
	for k, v := range dirties {
		ms.Put(k, v)
	}
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
