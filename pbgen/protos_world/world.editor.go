package protos_world

import (
	"fmt"
	protos_animal "github.com/yaoguangduan/proto-editor/pbgen/protos_animal"
	protos_enum "github.com/yaoguangduan/proto-editor/pbgen/protos_enum"
	protos_user "github.com/yaoguangduan/proto-editor/pbgen/protos_user"
	_ "github.com/yaoguangduan/proto-editor/protodef"
)

import "github.com/yaoguangduan/proto-editor/pbeditor"
import "iter"
import "math"
import "slices"
import "github.com/yaoguangduan/proto-editor/protodef"
import "google.golang.org/protobuf/encoding/protowire"
import "google.golang.org/protobuf/proto"

type TestEnumEditorList struct {
	parent   pbeditor.Editor
	idxInPar any
	value    []TestEnum
	monitor  []TestEnum
	change   []*protodef.ListOpLog
}

func NewTestEnumEditorList() *TestEnumEditorList {
	return &TestEnumEditorList{value: make([]TestEnum, 0)}
}

func (ms *TestEnumEditorList) SetParent(editor pbeditor.Editor, key any) {
	ms.parent = editor
	ms.idxInPar = key
}
func (ms *TestEnumEditorList) SetParentDirty() {
	if ms.parent != nil {
		ms.parent.SetDirty(ms.idxInPar, true, ms)
	}
}
func (ms *TestEnumEditorList) SetDirty(key any, dirty bool, sync pbeditor.Editor) {
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
func (ms *TestEnumEditorList) FlushDirty(dirty bool) {
	if !dirty {
		ms.monitor = nil
		ms.change = nil
	}
}

func (ms *TestEnumEditorList) Add(t TestEnum) *TestEnumEditorList {
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

func (ms *TestEnumEditorList) AddAll(t []TestEnum) *TestEnumEditorList {
	for _, v := range t {
		ms.Add(v)
	}
	return ms
}

func (ms *TestEnumEditorList) Insert(pos int, v TestEnum) *TestEnumEditorList {
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
func (ms *TestEnumEditorList) Set(pos int, v TestEnum) *TestEnumEditorList {
	ms.value[pos] = v
	ms.SetDirty(pos, true, ms)
	return ms
}
func (ms *TestEnumEditorList) Iter() iter.Seq2[int, TestEnum] {
	return func(yield func(int, TestEnum) bool) {
		for i := range ms.value {
			if !yield(i, ms.value[i]) {
				return
			}
		}
	}
}

func (ms *TestEnumEditorList) Clear() {
	ms.value = make([]TestEnum, 0)
	ms.monitor = make([]TestEnum, 0)
	ms.change = nil
	ms.change = append(ms.change, &protodef.ListOpLog{
		OpType: protodef.ListOpType_Clear,
		Pos:    0,
		RefIdx: 0,
	})
	ms.SetParentDirty()
}
func (ms *TestEnumEditorList) Len() int {
	return len(ms.value)
}

func (ms *TestEnumEditorList) Remove(v TestEnum) {
	ms.RemoveByIdx(slices.Index(ms.value, v))
}

func (ms *TestEnumEditorList) RemoveByIdx(pos int) {
	ms.value = slices.Delete(ms.value, pos, pos+1)
	ms.change = append(ms.change, &protodef.ListOpLog{
		OpType: protodef.ListOpType_Delete,
		Pos:    int32(pos),
	})
	ms.SetParentDirty()
}

func (ms *TestEnumEditorList) ChangeLog() []*protodef.ListOpLog {
	return ms.change
}
func (ms *TestEnumEditorList) MonitorDirtyList() []TestEnum {
	return ms.monitor
}
func (ms *TestEnumEditorList) Values() []TestEnum {
	vv := make([]TestEnum, ms.Len())
	copy(vv, ms.value)
	return vv
}

func (ms *TestEnumEditorList) MergeDirtyFromPb(logs []*protodef.ListOpLog, dirties []TestEnum) {
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

func (ms *TestEnumEditorList) Get(pos int32) TestEnum {
	return ms.value[pos]
}

type U32TestEnumEditorMap struct {
	parent      pbeditor.Editor
	keyInParent any
	value       map[uint32]TestEnum
	dirtied     map[uint32]struct{}
	deleted     map[uint32]struct{}
}

func NewU32TestEnumEditorMap() *U32TestEnumEditorMap {
	return &U32TestEnumEditorMap{value: make(map[uint32]TestEnum), dirtied: make(map[uint32]struct{}), deleted: make(map[uint32]struct{})}
}

func (ms *U32TestEnumEditorMap) SetParent(editor pbeditor.Editor, key any) {
	ms.parent = editor
	ms.keyInParent = key
}

func (ms *U32TestEnumEditorMap) SetParentDirty() {
	if ms.parent != nil {
		ms.parent.SetDirty(ms.keyInParent, true, ms)
	}
}

func (ms *U32TestEnumEditorMap) Entries() map[uint32]TestEnum {
	return ms.value
}

func (ms *U32TestEnumEditorMap) SetDirty(key any, dirty bool, editor pbeditor.Editor) {
	if dirty {
		ms.dirtied[key.(uint32)] = struct{}{}
		ms.SetParentDirty()
	} else {
		delete(ms.dirtied, key.(uint32))
	}

}

func (ms *U32TestEnumEditorMap) FlushDirty(dirty bool) {
	if !dirty {
		ms.dirtied = map[uint32]struct{}{}
		ms.deleted = map[uint32]struct{}{}
	}
}

func (ms *U32TestEnumEditorMap) Put(k uint32, v TestEnum) *U32TestEnumEditorMap {
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

func (ms *U32TestEnumEditorMap) PutAll(kv map[uint32]TestEnum) *U32TestEnumEditorMap {
	for k, v := range kv {
		ms.Put(k, v)
	}
	return ms
}

func (ms *U32TestEnumEditorMap) Len() int {
	return len(ms.value)
}

func (ms *U32TestEnumEditorMap) Clear() *U32TestEnumEditorMap {
	if ms.Len() <= 0 {
		return ms
	}
	for k, _ := range ms.value {
		ms.deleted[k] = struct{}{}
	}
	ms.value = map[uint32]TestEnum{}
	ms.dirtied = map[uint32]struct{}{}
	ms.SetParentDirty()
	return ms
}

func (ms *U32TestEnumEditorMap) Get(k uint32) (TestEnum, bool) {
	v, ok := ms.value[k]
	return v, ok
}

func (ms *U32TestEnumEditorMap) Remove(k uint32) (TestEnum, bool) {
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

func (ms *U32TestEnumEditorMap) RemoveAll(k []uint32) {
	for _, kk := range k {
		ms.Remove(kk)
	}
}

func (ms *U32TestEnumEditorMap) Iter() iter.Seq2[uint32, TestEnum] {
	return func(yield func(uint32, TestEnum) bool) {
		for k, v := range ms.value {
			if !yield(k, v) {
				return
			}
		}
	}
}

func (ms *U32TestEnumEditorMap) Dirtied() map[uint32]struct{} {
	return ms.dirtied
}

func (ms *U32TestEnumEditorMap) Deleted() map[uint32]struct{} {
	return ms.deleted
}

func (ms *U32TestEnumEditorMap) ContainDirtied(kk uint32) bool {
	for k := range ms.dirtied {
		if k == kk {
			return true
		}
	}
	return false
}

func (ms *U32TestEnumEditorMap) ContainDeleted(kk uint32) bool {
	for k := range ms.deleted {
		if k == kk {
			return true
		}
	}
	return false
}

func (ms *U32TestEnumEditorMap) MergeDirtyFromPb(delList []uint32, dirties map[uint32]TestEnum) {
	ms.RemoveAll(delList)
	for k, v := range dirties {
		ms.Put(k, v)
	}
}

type WorldEditor struct {
	id                int32
	desc              string
	alive             int64
	location          float32
	weight            float64
	manager           *protos_user.UserEditor
	animals           *protos_animal.I64AnimalEditorMap
	cl                protos_enum.CivilLevel
	cost              uint32
	area              int32
	treeCntMap        *pbeditor.StrI32Map
	stateRecords      *protos_enum.StateEditorList
	logs              *pbeditor.StrList
	historyResults    *pbeditor.BoolList
	assistants        *protos_user.StrUserEditorMap
	resultCnt         *pbeditor.BoolU32Map
	tracks            *pbeditor.U64List
	health            bool
	neighborNameTypes *protos_enum.StrWorldTypeEditorMap
	_dirties          map[int]struct{}
	_parent           pbeditor.Editor
	_keyInParent      any
}

func NewWorldEditor() *WorldEditor {
	return &WorldEditor{
		_dirties: make(map[int]struct{}),
	}
}
func (x *WorldEditor) SetDirty(key any, dirty bool, sync pbeditor.Editor) {

	if dirty {
		x._dirties[key.(int)] = struct{}{}
		x.SetParentDirty()
	} else {
		delete(x._dirties, key.(int))
	}
}
func (x *WorldEditor) SetParentDirty() {
	if x._parent != nil {
		x._parent.SetDirty(x._keyInParent, true, x)
	}
}
func (x *WorldEditor) SetParent(editor pbeditor.Editor, key any) {
	x._parent = editor
	x._keyInParent = key
}
func (x *WorldEditor) FlushDirty(dirty bool) {
	if dirty || x.isIdDirty() {
		x.setIdDirty(dirty, true)
	}
	if dirty || x.isDescDirty() {
		x.setDescDirty(dirty, true)
	}
	if dirty || x.isAliveDirty() {
		x.setAliveDirty(dirty, true)
	}
	if dirty || x.isLocationDirty() {
		x.setLocationDirty(dirty, true)
	}
	if dirty || x.isWeightDirty() {
		x.setWeightDirty(dirty, true)
	}
	if dirty || x.isManagerDirty() {
		x.setManagerDirty(dirty, true)
	}
	if dirty || x.isAnimalsDirty() {
		x.setAnimalsDirty(dirty, true)
	}
	if dirty || x.isClDirty() {
		x.setClDirty(dirty, true)
	}
	if dirty || x.isCostDirty() {
		x.setCostDirty(dirty, true)
	}
	if dirty || x.isAreaDirty() {
		x.setAreaDirty(dirty, true)
	}
	if dirty || x.isTreeCntMapDirty() {
		x.setTreeCntMapDirty(dirty, true)
	}
	if dirty || x.isStateRecordsDirty() {
		x.setStateRecordsDirty(dirty, true)
	}
	if dirty || x.isLogsDirty() {
		x.setLogsDirty(dirty, true)
	}
	if dirty || x.isHistoryResultsDirty() {
		x.setHistoryResultsDirty(dirty, true)
	}
	if dirty || x.isAssistantsDirty() {
		x.setAssistantsDirty(dirty, true)
	}
	if dirty || x.isResultCntDirty() {
		x.setResultCntDirty(dirty, true)
	}
	if dirty || x.isTracksDirty() {
		x.setTracksDirty(dirty, true)
	}
	if dirty || x.isHealthDirty() {
		x.setHealthDirty(dirty, true)
	}
	if dirty || x.isNeighborNameTypesDirty() {
		x.setNeighborNameTypesDirty(dirty, true)
	}
}
func (x *WorldEditor) setIdDirty(dirty bool, recur bool) {
	x.SetDirty(1, dirty, x)
}
func (x *WorldEditor) isIdDirty() bool {
	_, ok := x._dirties[1]
	return ok
}
func (x *WorldEditor) setDescDirty(dirty bool, recur bool) {
	x.SetDirty(2, dirty, x)
}
func (x *WorldEditor) isDescDirty() bool {
	_, ok := x._dirties[2]
	return ok
}
func (x *WorldEditor) setAliveDirty(dirty bool, recur bool) {
	x.SetDirty(3, dirty, x)
}
func (x *WorldEditor) isAliveDirty() bool {
	_, ok := x._dirties[3]
	return ok
}
func (x *WorldEditor) setLocationDirty(dirty bool, recur bool) {
	x.SetDirty(4, dirty, x)
}
func (x *WorldEditor) isLocationDirty() bool {
	_, ok := x._dirties[4]
	return ok
}
func (x *WorldEditor) setWeightDirty(dirty bool, recur bool) {
	x.SetDirty(5, dirty, x)
}
func (x *WorldEditor) isWeightDirty() bool {
	_, ok := x._dirties[5]
	return ok
}
func (x *WorldEditor) setManagerDirty(dirty bool, recur bool) {
	x.SetDirty(6, dirty, x)
	if recur && x.manager != nil {
		x.manager.FlushDirty(dirty)
	}
}
func (x *WorldEditor) isManagerDirty() bool {
	_, ok := x._dirties[6]
	return ok
}
func (x *WorldEditor) setAnimalsDirty(dirty bool, recur bool) {
	x.SetDirty(7, dirty, x)
	if recur && x.animals != nil {
		x.animals.FlushDirty(dirty)
	}
}
func (x *WorldEditor) isAnimalsDirty() bool {
	_, ok := x._dirties[7]
	return ok
}
func (x *WorldEditor) setClDirty(dirty bool, recur bool) {
	x.SetDirty(8, dirty, x)
}
func (x *WorldEditor) isClDirty() bool {
	_, ok := x._dirties[8]
	return ok
}
func (x *WorldEditor) setCostDirty(dirty bool, recur bool) {
	x.SetDirty(9, dirty, x)
}
func (x *WorldEditor) isCostDirty() bool {
	_, ok := x._dirties[9]
	return ok
}
func (x *WorldEditor) setAreaDirty(dirty bool, recur bool) {
	x.SetDirty(10, dirty, x)
}
func (x *WorldEditor) isAreaDirty() bool {
	_, ok := x._dirties[10]
	return ok
}
func (x *WorldEditor) setTreeCntMapDirty(dirty bool, recur bool) {
	x.SetDirty(11, dirty, x)
	if recur && x.treeCntMap != nil {
		x.treeCntMap.FlushDirty(dirty)
	}
}
func (x *WorldEditor) isTreeCntMapDirty() bool {
	_, ok := x._dirties[11]
	return ok
}
func (x *WorldEditor) setStateRecordsDirty(dirty bool, recur bool) {
	x.SetDirty(12, dirty, x)
	if recur && x.stateRecords != nil {
		x.stateRecords.FlushDirty(dirty)
	}
}
func (x *WorldEditor) isStateRecordsDirty() bool {
	_, ok := x._dirties[12]
	return ok
}
func (x *WorldEditor) setLogsDirty(dirty bool, recur bool) {
	x.SetDirty(13, dirty, x)
	if recur && x.logs != nil {
		x.logs.FlushDirty(dirty)
	}
}
func (x *WorldEditor) isLogsDirty() bool {
	_, ok := x._dirties[13]
	return ok
}
func (x *WorldEditor) setHistoryResultsDirty(dirty bool, recur bool) {
	x.SetDirty(14, dirty, x)
	if recur && x.historyResults != nil {
		x.historyResults.FlushDirty(dirty)
	}
}
func (x *WorldEditor) isHistoryResultsDirty() bool {
	_, ok := x._dirties[14]
	return ok
}
func (x *WorldEditor) setAssistantsDirty(dirty bool, recur bool) {
	x.SetDirty(15, dirty, x)
	if recur && x.assistants != nil {
		x.assistants.FlushDirty(dirty)
	}
}
func (x *WorldEditor) isAssistantsDirty() bool {
	_, ok := x._dirties[15]
	return ok
}
func (x *WorldEditor) setResultCntDirty(dirty bool, recur bool) {
	x.SetDirty(16, dirty, x)
	if recur && x.resultCnt != nil {
		x.resultCnt.FlushDirty(dirty)
	}
}
func (x *WorldEditor) isResultCntDirty() bool {
	_, ok := x._dirties[16]
	return ok
}
func (x *WorldEditor) setTracksDirty(dirty bool, recur bool) {
	x.SetDirty(17, dirty, x)
	if recur && x.tracks != nil {
		x.tracks.FlushDirty(dirty)
	}
}
func (x *WorldEditor) isTracksDirty() bool {
	_, ok := x._dirties[17]
	return ok
}
func (x *WorldEditor) setHealthDirty(dirty bool, recur bool) {
	x.SetDirty(18, dirty, x)
}
func (x *WorldEditor) isHealthDirty() bool {
	_, ok := x._dirties[18]
	return ok
}
func (x *WorldEditor) setNeighborNameTypesDirty(dirty bool, recur bool) {
	x.SetDirty(19, dirty, x)
	if recur && x.neighborNameTypes != nil {
		x.neighborNameTypes.FlushDirty(dirty)
	}
}
func (x *WorldEditor) isNeighborNameTypesDirty() bool {
	_, ok := x._dirties[19]
	return ok
}
func (x *WorldEditor) Key() any {
	return x._keyInParent
}
func (x *WorldEditor) SetKey(v any) {
	x._keyInParent = v
}
func (x *WorldEditor) Clear() *WorldEditor {
	x.SetId(0)
	x.SetDesc("")
	x.SetAlive(0)
	x.SetLocation(0)
	x.SetWeight(0)
	if x.manager != nil {
		x.manager.Clear()
	}
	if x.animals != nil {
		x.animals.Clear()
	}
	x.SetCl(protos_enum.CivilLevel_CL_Bronze)
	x.SetCost(0)
	x.SetArea(0)
	if x.treeCntMap != nil {
		x.treeCntMap.Clear()
	}
	if x.stateRecords != nil {
		x.stateRecords.Clear()
	}
	if x.logs != nil {
		x.logs.Clear()
	}
	if x.historyResults != nil {
		x.historyResults.Clear()
	}
	if x.assistants != nil {
		x.assistants.Clear()
	}
	if x.resultCnt != nil {
		x.resultCnt.Clear()
	}
	if x.tracks != nil {
		x.tracks.Clear()
	}
	x.SetHealth(false)
	if x.neighborNameTypes != nil {
		x.neighborNameTypes.Clear()
	}
	return x
}
func (x *WorldEditor) CopyTo(r *World) *WorldEditor {
	r.SetId(x.id)
	r.SetDesc(x.desc)
	r.SetAlive(x.alive)
	r.SetLocation(x.location)
	r.SetWeight(x.weight)
	if x.manager != nil {
		tmp := &protos_user.User{}
		x.manager.CopyTo(tmp)
		r.SetManager(tmp)
	}
	if x.animals != nil && x.animals.Len() > 0 {
		tmp := make(map[int64]*protos_animal.Animal)
		for k, v := range x.animals.Iter() {
			tmpV := &protos_animal.Animal{}
			v.CopyTo(tmpV)
			tmp[k] = tmpV
		}
		r.SetAnimals(tmp)
	}
	r.SetCl(x.cl)
	r.SetCost(x.cost)
	r.SetArea(x.area)
	if x.treeCntMap != nil && x.treeCntMap.Len() > 0 {
		tmp := make(map[string]int32)
		for k, v := range x.treeCntMap.Iter() {
			tmp[k] = v
		}
		r.SetTreeCntMap(tmp)
	}
	if x.stateRecords != nil && x.stateRecords.Len() > 0 {
		r.SetStateRecords(x.stateRecords.Values())
	}
	if x.logs != nil && x.logs.Len() > 0 {
		r.SetLogs(x.logs.Values())
	}
	if x.historyResults != nil && x.historyResults.Len() > 0 {
		r.SetHistoryResults(x.historyResults.Values())
	}
	if x.assistants != nil && x.assistants.Len() > 0 {
		tmp := make(map[string]*protos_user.User)
		for k, v := range x.assistants.Iter() {
			tmpV := &protos_user.User{}
			v.CopyTo(tmpV)
			tmp[k] = tmpV
		}
		r.SetAssistants(tmp)
	}
	if x.resultCnt != nil && x.resultCnt.Len() > 0 {
		tmp := make(map[bool]uint32)
		for k, v := range x.resultCnt.Iter() {
			tmp[k] = v
		}
		r.SetResultCnt(tmp)
	}
	if x.tracks != nil && x.tracks.Len() > 0 {
		r.SetTracks(x.tracks.Values())
	}
	r.SetHealth(x.health)
	if x.neighborNameTypes != nil && x.neighborNameTypes.Len() > 0 {
		tmp := make(map[string]protos_enum.WorldType)
		for k, v := range x.neighborNameTypes.Iter() {
			tmp[k] = v
		}
		r.SetNeighborNameTypes(tmp)
	}
	return x
}
func (x *WorldEditor) CopyFrom(r *World) *WorldEditor {
	if r.Id != nil {
		x.SetId(*r.Id)
	}
	if r.Desc != nil {
		x.SetDesc(*r.Desc)
	}
	if r.Alive != nil {
		x.SetAlive(*r.Alive)
	}
	if r.Location != nil {
		x.SetLocation(*r.Location)
	}
	if r.Weight != nil {
		x.SetWeight(*r.Weight)
	}
	if r.Manager != nil {
		x.GetManager().CopyFrom(r.Manager)
	}
	for k, v := range r.Animals {
		if v != nil {
			vv := protos_animal.NewAnimalEditor()
			vv.CopyFrom(v)
			x.GetAnimals().Put(k, vv)
		}
	}
	if r.Cl != nil {
		x.SetCl(*r.Cl)
	}
	if r.Cost != nil {
		x.SetCost(*r.Cost)
	}
	if r.Area != nil {
		x.SetArea(*r.Area)
	}
	for k, v := range r.TreeCntMap {
		x.GetTreeCntMap().Put(k, v)
	}
	if r.StateRecords != nil {
		x.GetStateRecords().AddAll(r.StateRecords)
	}
	if r.Logs != nil {
		x.GetLogs().AddAll(r.Logs)
	}
	if r.HistoryResults != nil {
		x.GetHistoryResults().AddAll(r.HistoryResults)
	}
	for k, v := range r.Assistants {
		if v != nil {
			vv := protos_user.NewUserEditor()
			vv.CopyFrom(v)
			x.GetAssistants().Put(k, vv)
		}
	}
	for k, v := range r.ResultCnt {
		x.GetResultCnt().Put(k, v)
	}
	if r.Tracks != nil {
		x.GetTracks().AddAll(r.Tracks)
	}
	if r.Health != nil {
		x.SetHealth(*r.Health)
	}
	for k, v := range r.NeighborNameTypes {
		x.GetNeighborNameTypes().Put(k, v)
	}
	return x
}
func (x *WorldEditor) MergeDirtyFromPb(r *World) {
	var mutState = pbeditor.UnmarshalMutState(r.ProtoReflect().GetUnknown())
	if r.Id != nil {
		x.SetId(*r.Id)
	}
	if r.Desc != nil {
		x.SetDesc(*r.Desc)
	}
	if r.Alive != nil {
		x.SetAlive(*r.Alive)
	}
	if r.Location != nil {
		x.SetLocation(*r.Location)
	}
	if r.Weight != nil {
		x.SetWeight(*r.Weight)
	}
	if r.Manager != nil {
		x.GetManager().MergeDirtyFromPb(r.Manager)
	}
	x.GetAnimals().MergeDirtyFromPb(mutState.I64List(7), r.Animals)
	if r.Cl != nil {
		x.SetCl(*r.Cl)
	}
	if r.Cost != nil {
		x.SetCost(*r.Cost)
	}
	if r.Area != nil {
		x.SetArea(*r.Area)
	}
	x.GetTreeCntMap().MergeDirtyFromPb(mutState.StrList(11), r.TreeCntMap)
	x.GetStateRecords().MergeDirtyFromPb(mutState.ListOpLog(12), r.StateRecords)
	x.GetLogs().MergeDirtyFromPb(mutState.ListOpLog(13), r.Logs)
	x.GetHistoryResults().MergeDirtyFromPb(mutState.ListOpLog(14), r.HistoryResults)
	x.GetAssistants().MergeDirtyFromPb(mutState.StrList(15), r.Assistants)
	x.GetResultCnt().MergeDirtyFromPb(mutState.BoolList(16), r.ResultCnt)
	x.GetTracks().MergeDirtyFromPb(mutState.ListOpLog(17), r.Tracks)
	if r.Health != nil {
		x.SetHealth(*r.Health)
	}
	x.GetNeighborNameTypes().MergeDirtyFromPb(mutState.StrList(19), r.NeighborNameTypes)
}
func (x *WorldEditor) MergeDirtyFromBytes(buf []byte) *WorldEditor {
	fds := pbeditor.ToRawMessage(buf)
	if x.animals != nil {
		x.animals.RemoveAll(fds.MutState.I64List(7))
	}
	if x.treeCntMap != nil {
		x.treeCntMap.RemoveAll(fds.MutState.StrList(11))
	}
	if x.assistants != nil {
		x.assistants.RemoveAll(fds.MutState.StrList(15))
	}
	if x.resultCnt != nil {
		x.resultCnt.RemoveAll(fds.MutState.BoolList(16))
	}
	if x.neighborNameTypes != nil {
		x.neighborNameTypes.RemoveAll(fds.MutState.StrList(19))
	}
	for _, rawFields := range fds.RawFields {
		if len(rawFields) <= 0 {
			continue
		}
		rawF := rawFields[0]
		switch rawF.Number {
		case 1:
			x.SetId(int32(rawF.Variant))
		case 2:
			x.SetDesc(pbeditor.Bys2Str(rawF.Bytes))
		case 3:
			x.SetAlive(int64(rawF.Variant))
		case 4:
			x.SetLocation(rawF.F32)
		case 5:
			x.SetWeight(rawF.F64)
		case 6:
			x.GetManager().MergeDirtyFromBytes(rawF.Bytes)
		case 7:
			r := make(map[int64]*protos_animal.Animal)
			for _, field := range rawFields {
				tmp := pbeditor.ParseMap(field.Bytes)
				k := tmp.KeyI64()
				v := &protos_animal.Animal{}
				proto.Unmarshal(field.Bytes, v)
				r[k] = v
			}
			x.GetAnimals().MergeDirtyFromPb(fds.MutState.I64List(rawF.Number), r)
		case 8:
			x.SetCl(protos_enum.CivilLevel(int32(rawF.Variant)))
		case 9:
			x.SetCost(uint32(rawF.Variant))
		case 10:
			x.SetArea(int32(rawF.Variant))
		case 11:
			r := make(map[string]int32)
			for _, field := range rawFields {
				tmp := pbeditor.ParseMap(field.Bytes)
				k := tmp.KeyStr()
				v := tmp.ValI32()
				r[k] = v
			}
			x.GetTreeCntMap().MergeDirtyFromPb(fds.MutState.StrList(rawF.Number), r)
		case 12:
			tmp := rawF.Bytes
			dl := make([]protos_enum.State, 0)
			for len(tmp) > 0 {
				val, n := protowire.ConsumeVarint(tmp)
				if n < 0 {
					panic(n)
				}
				tmp = tmp[n:]
				dl = append(dl, protos_enum.State(val))
			}
			x.GetStateRecords().MergeDirtyFromPb(fds.MutState.ListOpLog(rawF.Number), dl)
		case 13:
			sl := make([]string, 0)
			for _, field := range rawFields {
				sl = append(sl, pbeditor.Bys2Str(field.Bytes))
			}
			x.GetLogs().MergeDirtyFromPb(fds.MutState.ListOpLog(rawF.Number), sl)
		case 14:
			tmp := rawF.Bytes
			dl := make([]bool, 0)
			for len(tmp) > 0 {
				val, n := protowire.ConsumeVarint(tmp)
				if n < 0 {
					panic(n)
				}
				tmp = tmp[n:]
				dl = append(dl, protowire.DecodeBool(val))
			}
			x.GetHistoryResults().MergeDirtyFromPb(fds.MutState.ListOpLog(rawF.Number), dl)
		case 15:
			r := make(map[string]*protos_user.User)
			for _, field := range rawFields {
				tmp := pbeditor.ParseMap(field.Bytes)
				k := tmp.KeyStr()
				v := &protos_user.User{}
				proto.Unmarshal(field.Bytes, v)
				r[k] = v
			}
			x.GetAssistants().MergeDirtyFromPb(fds.MutState.StrList(rawF.Number), r)
		case 16:
			r := make(map[bool]uint32)
			for _, field := range rawFields {
				tmp := pbeditor.ParseMap(field.Bytes)
				k := tmp.KeyBool()
				v := tmp.ValU32()
				r[k] = v
			}
			x.GetResultCnt().MergeDirtyFromPb(fds.MutState.BoolList(rawF.Number), r)
		case 17:
			tmp := rawF.Bytes
			dl := make([]uint64, 0)
			for len(tmp) > 0 {
				val, n := protowire.ConsumeVarint(tmp)
				if n < 0 {
					panic(n)
				}
				tmp = tmp[n:]
				dl = append(dl, uint64(val))
			}
			x.GetTracks().MergeDirtyFromPb(fds.MutState.ListOpLog(rawF.Number), dl)
		case 18:
			x.SetHealth(protowire.DecodeBool(rawF.Variant))
		case 19:
			r := make(map[string]protos_enum.WorldType)
			for _, field := range rawFields {
				tmp := pbeditor.ParseMap(field.Bytes)
				k := tmp.KeyStr()
				v := tmp.ValI32()
				r[k] = protos_enum.WorldType(v)
			}
			x.GetNeighborNameTypes().MergeDirtyFromPb(fds.MutState.StrList(rawF.Number), r)
		}
	}
	return x
}
func (x *WorldEditor) GatherDirtyToPb(r *World) {
	mutState := pbeditor.UnmarshalMutState(r.unknownFields)
	if x.isIdDirty() {
		r.SetId(x.id)
	}
	if x.isDescDirty() {
		r.SetDesc(x.desc)
	}
	if x.isAliveDirty() {
		r.SetAlive(x.alive)
	}
	if x.isLocationDirty() {
		r.SetLocation(x.location)
	}
	if x.isWeightDirty() {
		r.SetWeight(x.weight)
	}
	if x.isManagerDirty() && x.manager != nil {
		if r.Manager == nil {
			r.Manager = &protos_user.User{}
		}
		x.manager.GatherDirtyToPb(r.Manager)
	}
	if x.isAnimalsDirty() {
		updated := make([]int64, 0)
		if r.Animals != nil {
			for k, v := range r.Animals {
				if x.animals.ContainDeleted(k) {
					delete(r.Animals, k)
					continue
				}
				if x.animals.ContainDirtied(k) {
					updated = append(updated, k)
					tmp := x.animals.Get(k)
					if tmp == nil {
						r.Animals[k] = nil
						continue
					}
					if v == nil {
						v = &protos_animal.Animal{}
						tmp.GatherDirtyToPb(v)
						r.Animals[k] = v

					} else {
						tmp.GatherDirtyToPb(v)
					}
				}
			}
		} else {
			r.Animals = make(map[int64]*protos_animal.Animal)
		}
		for k := range x.animals.Dirtied() {
			if !slices.Contains(updated, k) {
				tmp := x.animals.Get(k)
				if tmp == nil {
					r.Animals[k] = nil
				} else {
					v := &protos_animal.Animal{}
					tmp.GatherDirtyToPb(v)
					r.Animals[k] = v
				}
			}
		}
		deleted := mutState.I64List(7)
		for k := range x.animals.Deleted() {
			if !slices.Contains(deleted, k) {
				mutState.AddI64(7, k)
			}
		}
	}
	if x.isClDirty() {
		r.SetCl(x.cl)
	}
	if x.isCostDirty() {
		r.SetCost(x.cost)
	}
	if x.isAreaDirty() {
		r.SetArea(x.area)
	}
	if x.isTreeCntMapDirty() {
		updated := make([]string, 0)
		if r.TreeCntMap != nil {
			for k, v := range r.TreeCntMap {
				if x.treeCntMap.ContainDeleted(k) {
					delete(r.TreeCntMap, k)
					continue
				}
				if x.treeCntMap.ContainDirtied(k) {
					updated = append(updated, k)
					v, _ = x.treeCntMap.Get(k)
					r.TreeCntMap[k] = v
				}
			}
		} else {
			r.TreeCntMap = make(map[string]int32)
		}
		for k := range x.treeCntMap.Dirtied() {
			if !slices.Contains(updated, k) {
				tmp, _ := x.treeCntMap.Get(k)
				r.TreeCntMap[k] = tmp
			}
		}
		deleted := mutState.StrList(11)
		for k := range x.treeCntMap.Deleted() {
			if !slices.Contains(deleted, k) {
				mutState.AddStr(11, k)
			}
		}
	}
	if x.isStateRecordsDirty() {
		ori := r.StateRecords
		lenOfStateRecords := len(ori)
		logsOfStateRecords := x.stateRecords.ChangeLog()
		dirtiesOfStateRecords := x.stateRecords.MonitorDirtyList()
		for _, val := range dirtiesOfStateRecords {
			r.StateRecords = append(r.StateRecords, val)
		}
		for _, log := range logsOfStateRecords {
			log.Pos = int32(lenOfStateRecords) + log.Pos
			mutState.ListAppendLog(12, log)
		}
	}
	if x.isLogsDirty() {
		ori := r.Logs
		lenOfLogs := len(ori)
		logsOfLogs := x.logs.ChangeLog()
		dirtiesOfLogs := x.logs.MonitorDirtyList()
		for _, val := range dirtiesOfLogs {
			r.Logs = append(r.Logs, val)
		}
		for _, log := range logsOfLogs {
			log.Pos = int32(lenOfLogs) + log.Pos
			mutState.ListAppendLog(13, log)
		}
	}
	if x.isHistoryResultsDirty() {
		ori := r.HistoryResults
		lenOfHistoryResults := len(ori)
		logsOfHistoryResults := x.historyResults.ChangeLog()
		dirtiesOfHistoryResults := x.historyResults.MonitorDirtyList()
		for _, val := range dirtiesOfHistoryResults {
			r.HistoryResults = append(r.HistoryResults, val)
		}
		for _, log := range logsOfHistoryResults {
			log.Pos = int32(lenOfHistoryResults) + log.Pos
			mutState.ListAppendLog(14, log)
		}
	}
	if x.isAssistantsDirty() {
		updated := make([]string, 0)
		if r.Assistants != nil {
			for k, v := range r.Assistants {
				if x.assistants.ContainDeleted(k) {
					delete(r.Assistants, k)
					continue
				}
				if x.assistants.ContainDirtied(k) {
					updated = append(updated, k)
					tmp := x.assistants.Get(k)
					if tmp == nil {
						r.Assistants[k] = nil
						continue
					}
					if v == nil {
						v = &protos_user.User{}
						tmp.GatherDirtyToPb(v)
						r.Assistants[k] = v

					} else {
						tmp.GatherDirtyToPb(v)
					}
				}
			}
		} else {
			r.Assistants = make(map[string]*protos_user.User)
		}
		for k := range x.assistants.Dirtied() {
			if !slices.Contains(updated, k) {
				tmp := x.assistants.Get(k)
				if tmp == nil {
					r.Assistants[k] = nil
				} else {
					v := &protos_user.User{}
					tmp.GatherDirtyToPb(v)
					r.Assistants[k] = v
				}
			}
		}
		deleted := mutState.StrList(15)
		for k := range x.assistants.Deleted() {
			if !slices.Contains(deleted, k) {
				mutState.AddStr(15, k)
			}
		}
	}
	if x.isResultCntDirty() {
		updated := make([]bool, 0)
		if r.ResultCnt != nil {
			for k, v := range r.ResultCnt {
				if x.resultCnt.ContainDeleted(k) {
					delete(r.ResultCnt, k)
					continue
				}
				if x.resultCnt.ContainDirtied(k) {
					updated = append(updated, k)
					v, _ = x.resultCnt.Get(k)
					r.ResultCnt[k] = v
				}
			}
		} else {
			r.ResultCnt = make(map[bool]uint32)
		}
		for k := range x.resultCnt.Dirtied() {
			if !slices.Contains(updated, k) {
				tmp, _ := x.resultCnt.Get(k)
				r.ResultCnt[k] = tmp
			}
		}
		deleted := mutState.BoolList(16)
		for k := range x.resultCnt.Deleted() {
			if !slices.Contains(deleted, k) {
				mutState.AddBool(16, k)
			}
		}
	}
	if x.isTracksDirty() {
		ori := r.Tracks
		lenOfTracks := len(ori)
		logsOfTracks := x.tracks.ChangeLog()
		dirtiesOfTracks := x.tracks.MonitorDirtyList()
		for _, val := range dirtiesOfTracks {
			r.Tracks = append(r.Tracks, val)
		}
		for _, log := range logsOfTracks {
			log.Pos = int32(lenOfTracks) + log.Pos
			mutState.ListAppendLog(17, log)
		}
	}
	if x.isHealthDirty() {
		r.SetHealth(x.health)
	}
	if x.isNeighborNameTypesDirty() {
		updated := make([]string, 0)
		if r.NeighborNameTypes != nil {
			for k, v := range r.NeighborNameTypes {
				if x.neighborNameTypes.ContainDeleted(k) {
					delete(r.NeighborNameTypes, k)
					continue
				}
				if x.neighborNameTypes.ContainDirtied(k) {
					updated = append(updated, k)
					v, _ = x.neighborNameTypes.Get(k)
					r.NeighborNameTypes[k] = v
				}
			}
		} else {
			r.NeighborNameTypes = make(map[string]protos_enum.WorldType)
		}
		for k := range x.neighborNameTypes.Dirtied() {
			if !slices.Contains(updated, k) {
				tmp, _ := x.neighborNameTypes.Get(k)
				r.NeighborNameTypes[k] = tmp
			}
		}
		deleted := mutState.StrList(19)
		for k := range x.neighborNameTypes.Deleted() {
			if !slices.Contains(deleted, k) {
				mutState.AddStr(19, k)
			}
		}
	}
	r.ProtoReflect().SetUnknown(mutState.MarshalToFieldBytes())
}
func (x *WorldEditor) GatherDirtyToBytes() []byte {
	mutState := pbeditor.CreateMutStateInfo()
	var buf []byte
	if x.isIdDirty() {
		buf = protowire.AppendTag(buf, 1, 0)
		buf = protowire.AppendVarint(buf, uint64(x.id))
	}
	if x.isDescDirty() {
		buf = protowire.AppendTag(buf, 2, 2)
		buf = protowire.AppendString(buf, x.desc)
	}
	if x.isAliveDirty() {
		buf = protowire.AppendTag(buf, 3, 0)
		buf = protowire.AppendVarint(buf, uint64(x.alive))
	}
	if x.isLocationDirty() {
		buf = protowire.AppendTag(buf, 4, 5)
		buf = protowire.AppendFixed32(buf, math.Float32bits(x.location))
	}
	if x.isWeightDirty() {
		buf = protowire.AppendTag(buf, 5, 1)
		buf = protowire.AppendFixed64(buf, math.Float64bits(x.weight))
	}
	if x.isManagerDirty() && x.manager != nil {
		buf = protowire.AppendTag(buf, 6, protowire.BytesType)
		buf = protowire.AppendBytes(buf, x.manager.GatherDirtyToBytes())
	}
	if x.isAnimalsDirty() {
		if len(x.animals.Deleted()) > 0 {
			deleted := x.animals.Deleted()
			for v := range deleted {
				mutState.AddI64(7, v)
			}
		}
		if x.animals.Len() > 0 {
			for k, v := range x.animals.Iter() {
				if !x.animals.ContainDirtied(k) {
					continue
				}
				var tmp []byte
				tmp = protowire.AppendTag(tmp, 1, protowire.VarintType)
				tmp = protowire.AppendVarint(tmp, uint64(k))
				vb := v.GatherDirtyToBytes()
				tmp = protowire.AppendTag(tmp, 2, protowire.BytesType)
				tmp = protowire.AppendBytes(tmp, vb)
				buf = protowire.AppendTag(buf, 7, protowire.BytesType)
				buf = protowire.AppendBytes(buf, tmp)
			}
		}
	}
	if x.isClDirty() {
		buf = protowire.AppendTag(buf, 8, 0)
		buf = protowire.AppendVarint(buf, uint64(x.cl))
	}
	if x.isCostDirty() {
		buf = protowire.AppendTag(buf, 9, 0)
		buf = protowire.AppendVarint(buf, uint64(x.cost))
	}
	if x.isAreaDirty() {
		buf = protowire.AppendTag(buf, 10, 0)
		buf = protowire.AppendVarint(buf, uint64(x.area))
	}
	if x.isTreeCntMapDirty() {
		if len(x.treeCntMap.Deleted()) > 0 {
			deleted := x.treeCntMap.Deleted()
			for v := range deleted {
				mutState.AddStr(11, v)
			}
		}
		if x.treeCntMap.Len() > 0 {
			for k, v := range x.treeCntMap.Iter() {
				if !x.treeCntMap.ContainDirtied(k) {
					continue
				}
				var tmp []byte
				tmp = protowire.AppendTag(tmp, 1, protowire.BytesType)
				tmp = protowire.AppendString(tmp, k)
				tmp = protowire.AppendTag(tmp, 2, protowire.VarintType)
				tmp = protowire.AppendVarint(tmp, uint64(v))
				buf = protowire.AppendTag(buf, 11, protowire.BytesType)
				buf = protowire.AppendBytes(buf, tmp)
			}
		}
	}
	if x.isStateRecordsDirty() {
		for _, log := range x.stateRecords.ChangeLog() {
			mutState.ListAppendLog(12, log)
		}
		tmp := make([]byte, 0)
		for _, v := range x.stateRecords.MonitorDirtyList() {
			tmp = protowire.AppendVarint(tmp, uint64(v))
		}
		buf = protowire.AppendTag(buf, 12, protowire.BytesType)
		buf = protowire.AppendBytes(buf, tmp)
	}
	if x.isLogsDirty() {
		for _, log := range x.logs.ChangeLog() {
			mutState.ListAppendLog(13, log)
		}
		for _, s := range x.logs.MonitorDirtyList() {
			buf = protowire.AppendTag(buf, 13, protowire.BytesType)
			buf = protowire.AppendString(buf, s)
		}
	}
	if x.isHistoryResultsDirty() {
		for _, log := range x.historyResults.ChangeLog() {
			mutState.ListAppendLog(14, log)
		}
		tmp := make([]byte, 0)
		for _, b := range x.historyResults.MonitorDirtyList() {
			var t uint64 = 0
			if b {
				t = 1
			}
			tmp = protowire.AppendVarint(tmp, t)
		}
		buf = protowire.AppendTag(buf, 14, protowire.BytesType)
		buf = protowire.AppendBytes(buf, tmp)
	}
	if x.isAssistantsDirty() {
		if len(x.assistants.Deleted()) > 0 {
			deleted := x.assistants.Deleted()
			for v := range deleted {
				mutState.AddStr(15, v)
			}
		}
		if x.assistants.Len() > 0 {
			for k, v := range x.assistants.Iter() {
				if !x.assistants.ContainDirtied(k) {
					continue
				}
				var tmp []byte
				tmp = protowire.AppendTag(tmp, 1, protowire.BytesType)
				tmp = protowire.AppendString(tmp, k)
				vb := v.GatherDirtyToBytes()
				tmp = protowire.AppendTag(tmp, 2, protowire.BytesType)
				tmp = protowire.AppendBytes(tmp, vb)
				buf = protowire.AppendTag(buf, 15, protowire.BytesType)
				buf = protowire.AppendBytes(buf, tmp)
			}
		}
	}
	if x.isResultCntDirty() {
		if len(x.resultCnt.Deleted()) > 0 {
			deleted := x.resultCnt.Deleted()
			for v := range deleted {
				mutState.AddBool(16, v)
			}
		}
		if x.resultCnt.Len() > 0 {
			for k, v := range x.resultCnt.Iter() {
				if !x.resultCnt.ContainDirtied(k) {
					continue
				}
				var tmp []byte
				var t uint64 = 0
				if k {
					t = 1
				}
				tmp = protowire.AppendTag(tmp, 1, protowire.VarintType)
				tmp = protowire.AppendVarint(tmp, t)
				tmp = protowire.AppendTag(tmp, 2, protowire.VarintType)
				tmp = protowire.AppendVarint(tmp, uint64(v))
				buf = protowire.AppendTag(buf, 16, protowire.BytesType)
				buf = protowire.AppendBytes(buf, tmp)
			}
		}
	}
	if x.isTracksDirty() {
		for _, log := range x.tracks.ChangeLog() {
			mutState.ListAppendLog(17, log)
		}
		tmp := make([]byte, 0)
		for _, v := range x.tracks.MonitorDirtyList() {
			tmp = protowire.AppendVarint(tmp, uint64(v))
		}
		buf = protowire.AppendTag(buf, 17, protowire.BytesType)
		buf = protowire.AppendBytes(buf, tmp)
	}
	if x.isHealthDirty() {
		buf = protowire.AppendTag(buf, 18, 0)
		var t uint64 = 0
		if x.health {
			t = 1
		}
		buf = protowire.AppendVarint(buf, t)
	}
	if x.isNeighborNameTypesDirty() {
		if len(x.neighborNameTypes.Deleted()) > 0 {
			deleted := x.neighborNameTypes.Deleted()
			for v := range deleted {
				mutState.AddStr(19, v)
			}
		}
		if x.neighborNameTypes.Len() > 0 {
			for k, v := range x.neighborNameTypes.Iter() {
				if !x.neighborNameTypes.ContainDirtied(k) {
					continue
				}
				var tmp []byte
				tmp = protowire.AppendTag(tmp, 1, protowire.BytesType)
				tmp = protowire.AppendString(tmp, k)
				tmp = protowire.AppendTag(tmp, 2, protowire.VarintType)
				tmp = protowire.AppendVarint(tmp, uint64(v))
				buf = protowire.AppendTag(buf, 19, protowire.BytesType)
				buf = protowire.AppendBytes(buf, tmp)
			}
		}
	}
	buf = append(buf, mutState.MarshalToFieldBytes()...)
	return buf
}
func (x *WorldEditor) GetId() int32 {
	return x.id
}
func (x *WorldEditor) SetId(v int32) *WorldEditor {
	if x.id == v {
		return x
	}
	x.id = v
	x.setIdDirty(true, false)
	return x
}
func (x *WorldEditor) GetDesc() string {
	return x.desc
}
func (x *WorldEditor) SetDesc(v string) *WorldEditor {
	if x.desc == v {
		return x
	}
	x.desc = v
	x.setDescDirty(true, false)
	return x
}
func (x *WorldEditor) GetAlive() int64 {
	return x.alive
}
func (x *WorldEditor) SetAlive(v int64) *WorldEditor {
	if x.alive == v {
		return x
	}
	x.alive = v
	x.setAliveDirty(true, false)
	return x
}
func (x *WorldEditor) GetLocation() float32 {
	return x.location
}
func (x *WorldEditor) SetLocation(v float32) *WorldEditor {
	if x.location == v {
		return x
	}
	x.location = v
	x.setLocationDirty(true, false)
	return x
}
func (x *WorldEditor) GetWeight() float64 {
	return x.weight
}
func (x *WorldEditor) SetWeight(v float64) *WorldEditor {
	if x.weight == v {
		return x
	}
	x.weight = v
	x.setWeightDirty(true, false)
	return x
}
func (x *WorldEditor) GetManager() *protos_user.UserEditor {
	if x.manager == nil {
		x.manager = protos_user.NewUserEditor()
		x.manager.SetParent(x, 6)
	}
	return x.manager
}
func (x *WorldEditor) SetManager(v *protos_user.UserEditor) *WorldEditor {
	if v != nil {
		v.SetParent(x, 6)
	}
	if x.manager != nil {
		x.manager.SetParent(nil, -1)
	}
	x.manager = v
	x.setManagerDirty(true, false)
	return x
}
func (x *WorldEditor) GetAnimals() *protos_animal.I64AnimalEditorMap {
	if x.animals == nil {
		x.animals = protos_animal.NewI64AnimalEditorMap()
		x.animals.SetParent(x, 7)
	}
	return x.animals
}
func (x *WorldEditor) GetCl() protos_enum.CivilLevel {
	return x.cl
}
func (x *WorldEditor) SetCl(v protos_enum.CivilLevel) *WorldEditor {
	if x.cl == v {
		return x
	}
	x.cl = v
	x.setClDirty(true, false)
	return x
}
func (x *WorldEditor) GetCost() uint32 {
	return x.cost
}
func (x *WorldEditor) SetCost(v uint32) *WorldEditor {
	if x.cost == v {
		return x
	}
	x.cost = v
	x.setCostDirty(true, false)
	return x
}
func (x *WorldEditor) GetArea() int32 {
	return x.area
}
func (x *WorldEditor) SetArea(v int32) *WorldEditor {
	if x.area == v {
		return x
	}
	x.area = v
	x.setAreaDirty(true, false)
	return x
}
func (x *WorldEditor) GetTreeCntMap() *pbeditor.StrI32Map {
	if x.treeCntMap == nil {
		x.treeCntMap = pbeditor.NewStrI32Map()
		x.treeCntMap.SetParent(x, 11)
	}
	return x.treeCntMap
}
func (x *WorldEditor) GetStateRecords() *protos_enum.StateEditorList {
	if x.stateRecords == nil {
		x.stateRecords = protos_enum.NewStateEditorList()
		x.stateRecords.SetParent(x, 12)
	}
	return x.stateRecords
}
func (x *WorldEditor) AddStateRecords(v protos_enum.State) *WorldEditor {
	if x.stateRecords == nil {
		x.stateRecords = protos_enum.NewStateEditorList()
		x.stateRecords.SetParent(x, 12)
	}
	x.stateRecords.Add(v)
	return x
}
func (x *WorldEditor) GetLogs() *pbeditor.StrList {
	if x.logs == nil {
		x.logs = pbeditor.NewStrList()
		x.logs.SetParent(x, 13)
	}
	return x.logs
}
func (x *WorldEditor) AddLogs(v string) *WorldEditor {
	if x.logs == nil {
		x.logs = pbeditor.NewStrList()
		x.logs.SetParent(x, 13)
	}
	x.logs.Add(v)
	return x
}
func (x *WorldEditor) GetHistoryResults() *pbeditor.BoolList {
	if x.historyResults == nil {
		x.historyResults = pbeditor.NewBoolList()
		x.historyResults.SetParent(x, 14)
	}
	return x.historyResults
}
func (x *WorldEditor) AddHistoryResults(v bool) *WorldEditor {
	if x.historyResults == nil {
		x.historyResults = pbeditor.NewBoolList()
		x.historyResults.SetParent(x, 14)
	}
	x.historyResults.Add(v)
	return x
}
func (x *WorldEditor) GetAssistants() *protos_user.StrUserEditorMap {
	if x.assistants == nil {
		x.assistants = protos_user.NewStrUserEditorMap()
		x.assistants.SetParent(x, 15)
	}
	return x.assistants
}
func (x *WorldEditor) GetResultCnt() *pbeditor.BoolU32Map {
	if x.resultCnt == nil {
		x.resultCnt = pbeditor.NewBoolU32Map()
		x.resultCnt.SetParent(x, 16)
	}
	return x.resultCnt
}
func (x *WorldEditor) GetTracks() *pbeditor.U64List {
	if x.tracks == nil {
		x.tracks = pbeditor.NewU64List()
		x.tracks.SetParent(x, 17)
	}
	return x.tracks
}
func (x *WorldEditor) AddTracks(v uint64) *WorldEditor {
	if x.tracks == nil {
		x.tracks = pbeditor.NewU64List()
		x.tracks.SetParent(x, 17)
	}
	x.tracks.Add(v)
	return x
}
func (x *WorldEditor) GetHealth() bool {
	return x.health
}
func (x *WorldEditor) SetHealth(v bool) *WorldEditor {
	if x.health == v {
		return x
	}
	x.health = v
	x.setHealthDirty(true, false)
	return x
}
func (x *WorldEditor) GetNeighborNameTypes() *protos_enum.StrWorldTypeEditorMap {
	if x.neighborNameTypes == nil {
		x.neighborNameTypes = protos_enum.NewStrWorldTypeEditorMap()
		x.neighborNameTypes.SetParent(x, 19)
	}
	return x.neighborNameTypes
}
func (xs *World) SetId(v int32) {
	xs.Id = &v
}
func (xs *World) SetDesc(v string) {
	xs.Desc = &v
}
func (xs *World) SetAlive(v int64) {
	xs.Alive = &v
}
func (xs *World) SetLocation(v float32) {
	xs.Location = &v
}
func (xs *World) SetWeight(v float64) {
	xs.Weight = &v
}
func (xs *World) SetManager(v *protos_user.User) {
	xs.Manager = v
}
func (xs *World) SetAnimals(v map[int64]*protos_animal.Animal) {
	xs.Animals = v
}
func (xs *World) SetCl(v protos_enum.CivilLevel) {
	xs.Cl = &v
}
func (xs *World) SetCost(v uint32) {
	xs.Cost = &v
}
func (xs *World) SetArea(v int32) {
	xs.Area = &v
}
func (xs *World) SetTreeCntMap(v map[string]int32) {
	xs.TreeCntMap = v
}
func (xs *World) SetStateRecords(v []protos_enum.State) {
	xs.StateRecords = v
}
func (xs *World) AddStateRecords(v protos_enum.State) {
	xs.StateRecords = append(xs.StateRecords, v)
}
func (xs *World) SetLogs(v []string) {
	xs.Logs = v
}
func (xs *World) AddLogs(v string) {
	xs.Logs = append(xs.Logs, v)
}
func (xs *World) SetHistoryResults(v []bool) {
	xs.HistoryResults = v
}
func (xs *World) AddHistoryResults(v bool) {
	xs.HistoryResults = append(xs.HistoryResults, v)
}
func (xs *World) SetAssistants(v map[string]*protos_user.User) {
	xs.Assistants = v
}
func (xs *World) SetResultCnt(v map[bool]uint32) {
	xs.ResultCnt = v
}
func (xs *World) SetTracks(v []uint64) {
	xs.Tracks = v
}
func (xs *World) AddTracks(v uint64) {
	xs.Tracks = append(xs.Tracks, v)
}
func (xs *World) SetHealth(v bool) {
	xs.Health = &v
}
func (xs *World) SetNeighborNameTypes(v map[string]protos_enum.WorldType) {
	xs.NeighborNameTypes = v
}

type TestEditor struct {
	i32          int32
	i64          int64
	u32          uint32
	u64          uint64
	b            bool
	str          string
	bys          []byte
	te           TestEnum
	f32          float32
	f64          float64
	msg          *TestMsgEditor
	i32List      *pbeditor.I32List
	i64List      *pbeditor.I64List
	u32List      *pbeditor.U32List
	u64List      *pbeditor.U64List
	bList        *pbeditor.BoolList
	strList      *pbeditor.StrList
	teList       *TestEnumEditorList
	f32List      *pbeditor.F32List
	f64List      *pbeditor.F64List
	msgList      *TestMsgEditorList
	i32StrMap    *pbeditor.I32StrMap
	strF32Map    *pbeditor.StrF32Map
	boolU64Map   *pbeditor.BoolU64Map
	u64F64Map    *pbeditor.U64F64Map
	msgMap       *I32TestMsgEditorMap
	strMsgMap    *StrPlainMsgEditorMap
	enumMap      *U32TestEnumEditorMap
	_dirties     map[int]struct{}
	_parent      pbeditor.Editor
	_keyInParent any
}

func NewTestEditor() *TestEditor {
	return &TestEditor{
		_dirties: make(map[int]struct{}),
	}
}
func (x *TestEditor) SetDirty(key any, dirty bool, sync pbeditor.Editor) {

	if dirty {
		x._dirties[key.(int)] = struct{}{}
		x.SetParentDirty()
	} else {
		delete(x._dirties, key.(int))
	}
}
func (x *TestEditor) SetParentDirty() {
	if x._parent != nil {
		x._parent.SetDirty(x._keyInParent, true, x)
	}
}
func (x *TestEditor) SetParent(editor pbeditor.Editor, key any) {
	x._parent = editor
	x._keyInParent = key
}
func (x *TestEditor) FlushDirty(dirty bool) {
	if dirty || x.isI32Dirty() {
		x.setI32Dirty(dirty, true)
	}
	if dirty || x.isI64Dirty() {
		x.setI64Dirty(dirty, true)
	}
	if dirty || x.isU32Dirty() {
		x.setU32Dirty(dirty, true)
	}
	if dirty || x.isU64Dirty() {
		x.setU64Dirty(dirty, true)
	}
	if dirty || x.isBDirty() {
		x.setBDirty(dirty, true)
	}
	if dirty || x.isStrDirty() {
		x.setStrDirty(dirty, true)
	}
	if dirty || x.isBysDirty() {
		x.setBysDirty(dirty, true)
	}
	if dirty || x.isTeDirty() {
		x.setTeDirty(dirty, true)
	}
	if dirty || x.isF32Dirty() {
		x.setF32Dirty(dirty, true)
	}
	if dirty || x.isF64Dirty() {
		x.setF64Dirty(dirty, true)
	}
	if dirty || x.isMsgDirty() {
		x.setMsgDirty(dirty, true)
	}
	if dirty || x.isI32ListDirty() {
		x.setI32ListDirty(dirty, true)
	}
	if dirty || x.isI64ListDirty() {
		x.setI64ListDirty(dirty, true)
	}
	if dirty || x.isU32ListDirty() {
		x.setU32ListDirty(dirty, true)
	}
	if dirty || x.isU64ListDirty() {
		x.setU64ListDirty(dirty, true)
	}
	if dirty || x.isBListDirty() {
		x.setBListDirty(dirty, true)
	}
	if dirty || x.isStrListDirty() {
		x.setStrListDirty(dirty, true)
	}
	if dirty || x.isTeListDirty() {
		x.setTeListDirty(dirty, true)
	}
	if dirty || x.isF32ListDirty() {
		x.setF32ListDirty(dirty, true)
	}
	if dirty || x.isF64ListDirty() {
		x.setF64ListDirty(dirty, true)
	}
	if dirty || x.isMsgListDirty() {
		x.setMsgListDirty(dirty, true)
	}
	if dirty || x.isI32StrMapDirty() {
		x.setI32StrMapDirty(dirty, true)
	}
	if dirty || x.isStrF32MapDirty() {
		x.setStrF32MapDirty(dirty, true)
	}
	if dirty || x.isBoolU64MapDirty() {
		x.setBoolU64MapDirty(dirty, true)
	}
	if dirty || x.isU64F64MapDirty() {
		x.setU64F64MapDirty(dirty, true)
	}
	if dirty || x.isMsgMapDirty() {
		x.setMsgMapDirty(dirty, true)
	}
	if dirty || x.isStrMsgMapDirty() {
		x.setStrMsgMapDirty(dirty, true)
	}
	if dirty || x.isEnumMapDirty() {
		x.setEnumMapDirty(dirty, true)
	}
}
func (x *TestEditor) setI32Dirty(dirty bool, recur bool) {
	x.SetDirty(1, dirty, x)
}
func (x *TestEditor) isI32Dirty() bool {
	_, ok := x._dirties[1]
	return ok
}
func (x *TestEditor) setI64Dirty(dirty bool, recur bool) {
	x.SetDirty(2, dirty, x)
}
func (x *TestEditor) isI64Dirty() bool {
	_, ok := x._dirties[2]
	return ok
}
func (x *TestEditor) setU32Dirty(dirty bool, recur bool) {
	x.SetDirty(3, dirty, x)
}
func (x *TestEditor) isU32Dirty() bool {
	_, ok := x._dirties[3]
	return ok
}
func (x *TestEditor) setU64Dirty(dirty bool, recur bool) {
	x.SetDirty(4, dirty, x)
}
func (x *TestEditor) isU64Dirty() bool {
	_, ok := x._dirties[4]
	return ok
}
func (x *TestEditor) setBDirty(dirty bool, recur bool) {
	x.SetDirty(11, dirty, x)
}
func (x *TestEditor) isBDirty() bool {
	_, ok := x._dirties[11]
	return ok
}
func (x *TestEditor) setStrDirty(dirty bool, recur bool) {
	x.SetDirty(12, dirty, x)
}
func (x *TestEditor) isStrDirty() bool {
	_, ok := x._dirties[12]
	return ok
}
func (x *TestEditor) setBysDirty(dirty bool, recur bool) {
	x.SetDirty(13, dirty, x)
}
func (x *TestEditor) isBysDirty() bool {
	_, ok := x._dirties[13]
	return ok
}
func (x *TestEditor) setTeDirty(dirty bool, recur bool) {
	x.SetDirty(14, dirty, x)
}
func (x *TestEditor) isTeDirty() bool {
	_, ok := x._dirties[14]
	return ok
}
func (x *TestEditor) setF32Dirty(dirty bool, recur bool) {
	x.SetDirty(15, dirty, x)
}
func (x *TestEditor) isF32Dirty() bool {
	_, ok := x._dirties[15]
	return ok
}
func (x *TestEditor) setF64Dirty(dirty bool, recur bool) {
	x.SetDirty(16, dirty, x)
}
func (x *TestEditor) isF64Dirty() bool {
	_, ok := x._dirties[16]
	return ok
}
func (x *TestEditor) setMsgDirty(dirty bool, recur bool) {
	x.SetDirty(17, dirty, x)
	if recur && x.msg != nil {
		x.msg.FlushDirty(dirty)
	}
}
func (x *TestEditor) isMsgDirty() bool {
	_, ok := x._dirties[17]
	return ok
}
func (x *TestEditor) setI32ListDirty(dirty bool, recur bool) {
	x.SetDirty(18, dirty, x)
	if recur && x.i32List != nil {
		x.i32List.FlushDirty(dirty)
	}
}
func (x *TestEditor) isI32ListDirty() bool {
	_, ok := x._dirties[18]
	return ok
}
func (x *TestEditor) setI64ListDirty(dirty bool, recur bool) {
	x.SetDirty(28, dirty, x)
	if recur && x.i64List != nil {
		x.i64List.FlushDirty(dirty)
	}
}
func (x *TestEditor) isI64ListDirty() bool {
	_, ok := x._dirties[28]
	return ok
}
func (x *TestEditor) setU32ListDirty(dirty bool, recur bool) {
	x.SetDirty(38, dirty, x)
	if recur && x.u32List != nil {
		x.u32List.FlushDirty(dirty)
	}
}
func (x *TestEditor) isU32ListDirty() bool {
	_, ok := x._dirties[38]
	return ok
}
func (x *TestEditor) setU64ListDirty(dirty bool, recur bool) {
	x.SetDirty(48, dirty, x)
	if recur && x.u64List != nil {
		x.u64List.FlushDirty(dirty)
	}
}
func (x *TestEditor) isU64ListDirty() bool {
	_, ok := x._dirties[48]
	return ok
}
func (x *TestEditor) setBListDirty(dirty bool, recur bool) {
	x.SetDirty(118, dirty, x)
	if recur && x.bList != nil {
		x.bList.FlushDirty(dirty)
	}
}
func (x *TestEditor) isBListDirty() bool {
	_, ok := x._dirties[118]
	return ok
}
func (x *TestEditor) setStrListDirty(dirty bool, recur bool) {
	x.SetDirty(128, dirty, x)
	if recur && x.strList != nil {
		x.strList.FlushDirty(dirty)
	}
}
func (x *TestEditor) isStrListDirty() bool {
	_, ok := x._dirties[128]
	return ok
}
func (x *TestEditor) setTeListDirty(dirty bool, recur bool) {
	x.SetDirty(148, dirty, x)
	if recur && x.teList != nil {
		x.teList.FlushDirty(dirty)
	}
}
func (x *TestEditor) isTeListDirty() bool {
	_, ok := x._dirties[148]
	return ok
}
func (x *TestEditor) setF32ListDirty(dirty bool, recur bool) {
	x.SetDirty(158, dirty, x)
	if recur && x.f32List != nil {
		x.f32List.FlushDirty(dirty)
	}
}
func (x *TestEditor) isF32ListDirty() bool {
	_, ok := x._dirties[158]
	return ok
}
func (x *TestEditor) setF64ListDirty(dirty bool, recur bool) {
	x.SetDirty(168, dirty, x)
	if recur && x.f64List != nil {
		x.f64List.FlushDirty(dirty)
	}
}
func (x *TestEditor) isF64ListDirty() bool {
	_, ok := x._dirties[168]
	return ok
}
func (x *TestEditor) setMsgListDirty(dirty bool, recur bool) {
	x.SetDirty(178, dirty, x)
	if recur && x.msgList != nil {
		x.msgList.FlushDirty(dirty)
	}
}
func (x *TestEditor) isMsgListDirty() bool {
	_, ok := x._dirties[178]
	return ok
}
func (x *TestEditor) setI32StrMapDirty(dirty bool, recur bool) {
	x.SetDirty(19, dirty, x)
	if recur && x.i32StrMap != nil {
		x.i32StrMap.FlushDirty(dirty)
	}
}
func (x *TestEditor) isI32StrMapDirty() bool {
	_, ok := x._dirties[19]
	return ok
}
func (x *TestEditor) setStrF32MapDirty(dirty bool, recur bool) {
	x.SetDirty(29, dirty, x)
	if recur && x.strF32Map != nil {
		x.strF32Map.FlushDirty(dirty)
	}
}
func (x *TestEditor) isStrF32MapDirty() bool {
	_, ok := x._dirties[29]
	return ok
}
func (x *TestEditor) setBoolU64MapDirty(dirty bool, recur bool) {
	x.SetDirty(39, dirty, x)
	if recur && x.boolU64Map != nil {
		x.boolU64Map.FlushDirty(dirty)
	}
}
func (x *TestEditor) isBoolU64MapDirty() bool {
	_, ok := x._dirties[39]
	return ok
}
func (x *TestEditor) setU64F64MapDirty(dirty bool, recur bool) {
	x.SetDirty(49, dirty, x)
	if recur && x.u64F64Map != nil {
		x.u64F64Map.FlushDirty(dirty)
	}
}
func (x *TestEditor) isU64F64MapDirty() bool {
	_, ok := x._dirties[49]
	return ok
}
func (x *TestEditor) setMsgMapDirty(dirty bool, recur bool) {
	x.SetDirty(59, dirty, x)
	if recur && x.msgMap != nil {
		x.msgMap.FlushDirty(dirty)
	}
}
func (x *TestEditor) isMsgMapDirty() bool {
	_, ok := x._dirties[59]
	return ok
}
func (x *TestEditor) setStrMsgMapDirty(dirty bool, recur bool) {
	x.SetDirty(69, dirty, x)
	if recur && x.strMsgMap != nil {
		x.strMsgMap.FlushDirty(dirty)
	}
}
func (x *TestEditor) isStrMsgMapDirty() bool {
	_, ok := x._dirties[69]
	return ok
}
func (x *TestEditor) setEnumMapDirty(dirty bool, recur bool) {
	x.SetDirty(79, dirty, x)
	if recur && x.enumMap != nil {
		x.enumMap.FlushDirty(dirty)
	}
}
func (x *TestEditor) isEnumMapDirty() bool {
	_, ok := x._dirties[79]
	return ok
}
func (x *TestEditor) Key() any {
	return x._keyInParent
}
func (x *TestEditor) SetKey(v any) {
	x._keyInParent = v
}
func (x *TestEditor) Clear() *TestEditor {
	x.SetI32(0)
	x.SetI64(0)
	x.SetU32(0)
	x.SetU64(0)
	x.SetB(false)
	x.SetStr("")
	x.SetBys(make([]byte, 0))
	x.SetTe(TestEnum_Test0)
	x.SetF32(0)
	x.SetF64(0)
	if x.msg != nil {
		x.msg.Clear()
	}
	if x.i32List != nil {
		x.i32List.Clear()
	}
	if x.i64List != nil {
		x.i64List.Clear()
	}
	if x.u32List != nil {
		x.u32List.Clear()
	}
	if x.u64List != nil {
		x.u64List.Clear()
	}
	if x.bList != nil {
		x.bList.Clear()
	}
	if x.strList != nil {
		x.strList.Clear()
	}
	if x.teList != nil {
		x.teList.Clear()
	}
	if x.f32List != nil {
		x.f32List.Clear()
	}
	if x.f64List != nil {
		x.f64List.Clear()
	}
	if x.msgList != nil {
		x.msgList.Clear()
	}
	if x.i32StrMap != nil {
		x.i32StrMap.Clear()
	}
	if x.strF32Map != nil {
		x.strF32Map.Clear()
	}
	if x.boolU64Map != nil {
		x.boolU64Map.Clear()
	}
	if x.u64F64Map != nil {
		x.u64F64Map.Clear()
	}
	if x.msgMap != nil {
		x.msgMap.Clear()
	}
	if x.strMsgMap != nil {
		x.strMsgMap.Clear()
	}
	if x.enumMap != nil {
		x.enumMap.Clear()
	}
	return x
}
func (x *TestEditor) CopyTo(r *Test) *TestEditor {
	r.SetI32(x.i32)
	r.SetI64(x.i64)
	r.SetU32(x.u32)
	r.SetU64(x.u64)
	r.SetB(x.b)
	r.SetStr(x.str)
	r.SetBys(x.bys)
	r.SetBys(slices.Clone(x.bys))
	r.SetTe(x.te)
	r.SetF32(x.f32)
	r.SetF64(x.f64)
	if x.msg != nil {
		tmp := &TestMsg{}
		x.msg.CopyTo(tmp)
		r.SetMsg(tmp)
	}
	if x.i32List != nil && x.i32List.Len() > 0 {
		r.SetI32List(x.i32List.Values())
	}
	if x.i64List != nil && x.i64List.Len() > 0 {
		r.SetI64List(x.i64List.Values())
	}
	if x.u32List != nil && x.u32List.Len() > 0 {
		r.SetU32List(x.u32List.Values())
	}
	if x.u64List != nil && x.u64List.Len() > 0 {
		r.SetU64List(x.u64List.Values())
	}
	if x.bList != nil && x.bList.Len() > 0 {
		r.SetBList(x.bList.Values())
	}
	if x.strList != nil && x.strList.Len() > 0 {
		r.SetStrList(x.strList.Values())
	}
	if x.teList != nil && x.teList.Len() > 0 {
		r.SetTeList(x.teList.Values())
	}
	if x.f32List != nil && x.f32List.Len() > 0 {
		r.SetF32List(x.f32List.Values())
	}
	if x.f64List != nil && x.f64List.Len() > 0 {
		r.SetF64List(x.f64List.Values())
	}
	if x.msgList != nil && x.msgList.Len() > 0 {

		for _, v := range x.msgList.Iter() {
			if v == nil {
				r.AddMsgList(nil)
			} else {
				tmp := &TestMsg{}
				v.CopyTo(tmp)
				r.AddMsgList(tmp)
			}
		}

	}
	if x.i32StrMap != nil && x.i32StrMap.Len() > 0 {
		tmp := make(map[int32]string)
		for k, v := range x.i32StrMap.Iter() {
			tmp[k] = v
		}
		r.SetI32StrMap(tmp)
	}
	if x.strF32Map != nil && x.strF32Map.Len() > 0 {
		tmp := make(map[string]float32)
		for k, v := range x.strF32Map.Iter() {
			tmp[k] = v
		}
		r.SetStrF32Map(tmp)
	}
	if x.boolU64Map != nil && x.boolU64Map.Len() > 0 {
		tmp := make(map[bool]uint64)
		for k, v := range x.boolU64Map.Iter() {
			tmp[k] = v
		}
		r.SetBoolU64Map(tmp)
	}
	if x.u64F64Map != nil && x.u64F64Map.Len() > 0 {
		tmp := make(map[uint64]float64)
		for k, v := range x.u64F64Map.Iter() {
			tmp[k] = v
		}
		r.SetU64F64Map(tmp)
	}
	if x.msgMap != nil && x.msgMap.Len() > 0 {
		tmp := make(map[int32]*TestMsg)
		for k, v := range x.msgMap.Iter() {
			tmpV := &TestMsg{}
			v.CopyTo(tmpV)
			tmp[k] = tmpV
		}
		r.SetMsgMap(tmp)
	}
	if x.strMsgMap != nil && x.strMsgMap.Len() > 0 {
		tmp := make(map[string]*PlainMsg)
		for k, v := range x.strMsgMap.Iter() {
			tmpV := &PlainMsg{}
			v.CopyTo(tmpV)
			tmp[k] = tmpV
		}
		r.SetStrMsgMap(tmp)
	}
	if x.enumMap != nil && x.enumMap.Len() > 0 {
		tmp := make(map[uint32]TestEnum)
		for k, v := range x.enumMap.Iter() {
			tmp[k] = v
		}
		r.SetEnumMap(tmp)
	}
	return x
}
func (x *TestEditor) CopyFrom(r *Test) *TestEditor {
	if r.I32 != nil {
		x.SetI32(*r.I32)
	}
	if r.I64 != nil {
		x.SetI64(*r.I64)
	}
	if r.U32 != nil {
		x.SetU32(*r.U32)
	}
	if r.U64 != nil {
		x.SetU64(*r.U64)
	}
	if r.B != nil {
		x.SetB(*r.B)
	}
	if r.Str != nil {
		x.SetStr(*r.Str)
	}
	if len(r.Bys) > 0 {
		x.SetBys(slices.Clone(r.Bys))
	}
	if r.Te != nil {
		x.SetTe(*r.Te)
	}
	if r.F32 != nil {
		x.SetF32(*r.F32)
	}
	if r.F64 != nil {
		x.SetF64(*r.F64)
	}
	if r.Msg != nil {
		x.GetMsg().CopyFrom(r.Msg)
	}
	if r.I32List != nil {
		x.GetI32List().AddAll(r.I32List)
	}
	if r.I64List != nil {
		x.GetI64List().AddAll(r.I64List)
	}
	if r.U32List != nil {
		x.GetU32List().AddAll(r.U32List)
	}
	if r.U64List != nil {
		x.GetU64List().AddAll(r.U64List)
	}
	if r.BList != nil {
		x.GetBList().AddAll(r.BList)
	}
	if r.StrList != nil {
		x.GetStrList().AddAll(r.StrList)
	}
	if r.TeList != nil {
		x.GetTeList().AddAll(r.TeList)
	}
	if r.F32List != nil {
		x.GetF32List().AddAll(r.F32List)
	}
	if r.F64List != nil {
		x.GetF64List().AddAll(r.F64List)
	}
	if r.MsgList != nil {

		for _, v := range r.MsgList {
			if v == nil {
				x.GetMsgList().Add(nil)
			} else {
				tmp := NewTestMsgEditor()
				tmp.CopyFrom(v)
				x.GetMsgList().Add(tmp)
			}
		}

	}
	for k, v := range r.I32StrMap {
		x.GetI32StrMap().Put(k, v)
	}
	for k, v := range r.StrF32Map {
		x.GetStrF32Map().Put(k, v)
	}
	for k, v := range r.BoolU64Map {
		x.GetBoolU64Map().Put(k, v)
	}
	for k, v := range r.U64F64Map {
		x.GetU64F64Map().Put(k, v)
	}
	for k, v := range r.MsgMap {
		if v != nil {
			vv := NewTestMsgEditor()
			vv.CopyFrom(v)
			x.GetMsgMap().Put(k, vv)
		}
	}
	for k, v := range r.StrMsgMap {
		if v != nil {
			vv := NewPlainMsgEditor()
			vv.CopyFrom(v)
			x.GetStrMsgMap().Put(k, vv)
		}
	}
	for k, v := range r.EnumMap {
		x.GetEnumMap().Put(k, v)
	}
	return x
}
func (x *TestEditor) MergeDirtyFromPb(r *Test) {
	var mutState = pbeditor.UnmarshalMutState(r.ProtoReflect().GetUnknown())
	if r.I32 != nil {
		x.SetI32(*r.I32)
	}
	if r.I64 != nil {
		x.SetI64(*r.I64)
	}
	if r.U32 != nil {
		x.SetU32(*r.U32)
	}
	if r.U64 != nil {
		x.SetU64(*r.U64)
	}
	if r.B != nil {
		x.SetB(*r.B)
	}
	if r.Str != nil {
		x.SetStr(*r.Str)
	}
	if len(r.Bys) > 0 {
		x.SetBys(slices.Clone(r.Bys))
	}
	if r.Te != nil {
		x.SetTe(*r.Te)
	}
	if r.F32 != nil {
		x.SetF32(*r.F32)
	}
	if r.F64 != nil {
		x.SetF64(*r.F64)
	}
	if r.Msg != nil {
		x.GetMsg().MergeDirtyFromPb(r.Msg)
	}
	x.GetI32List().MergeDirtyFromPb(mutState.ListOpLog(18), r.I32List)
	x.GetI64List().MergeDirtyFromPb(mutState.ListOpLog(28), r.I64List)
	x.GetU32List().MergeDirtyFromPb(mutState.ListOpLog(38), r.U32List)
	x.GetU64List().MergeDirtyFromPb(mutState.ListOpLog(48), r.U64List)
	x.GetBList().MergeDirtyFromPb(mutState.ListOpLog(118), r.BList)
	x.GetStrList().MergeDirtyFromPb(mutState.ListOpLog(128), r.StrList)
	x.GetTeList().MergeDirtyFromPb(mutState.ListOpLog(148), r.TeList)
	x.GetF32List().MergeDirtyFromPb(mutState.ListOpLog(158), r.F32List)
	x.GetF64List().MergeDirtyFromPb(mutState.ListOpLog(168), r.F64List)
	x.GetMsgList().MergeDirtyFromPb(mutState.ListOpLog(178), r.MsgList)
	x.GetI32StrMap().MergeDirtyFromPb(mutState.I32List(19), r.I32StrMap)
	x.GetStrF32Map().MergeDirtyFromPb(mutState.StrList(29), r.StrF32Map)
	x.GetBoolU64Map().MergeDirtyFromPb(mutState.BoolList(39), r.BoolU64Map)
	x.GetU64F64Map().MergeDirtyFromPb(mutState.U64List(49), r.U64F64Map)
	x.GetMsgMap().MergeDirtyFromPb(mutState.I32List(59), r.MsgMap)
	x.GetStrMsgMap().MergeDirtyFromPb(mutState.StrList(69), r.StrMsgMap)
	x.GetEnumMap().MergeDirtyFromPb(mutState.U32List(79), r.EnumMap)
}
func (x *TestEditor) MergeDirtyFromBytes(buf []byte) *TestEditor {
	fds := pbeditor.ToRawMessage(buf)
	if x.i32StrMap != nil {
		x.i32StrMap.RemoveAll(fds.MutState.I32List(19))
	}
	if x.strF32Map != nil {
		x.strF32Map.RemoveAll(fds.MutState.StrList(29))
	}
	if x.boolU64Map != nil {
		x.boolU64Map.RemoveAll(fds.MutState.BoolList(39))
	}
	if x.u64F64Map != nil {
		x.u64F64Map.RemoveAll(fds.MutState.U64List(49))
	}
	if x.msgMap != nil {
		x.msgMap.RemoveAll(fds.MutState.I32List(59))
	}
	if x.strMsgMap != nil {
		x.strMsgMap.RemoveAll(fds.MutState.StrList(69))
	}
	if x.enumMap != nil {
		x.enumMap.RemoveAll(fds.MutState.U32List(79))
	}
	for _, rawFields := range fds.RawFields {
		if len(rawFields) <= 0 {
			continue
		}
		rawF := rawFields[0]
		switch rawF.Number {
		case 1:
			x.SetI32(int32(rawF.Variant))
		case 2:
			x.SetI64(int64(rawF.Variant))
		case 3:
			x.SetU32(uint32(rawF.Variant))
		case 4:
			x.SetU64(uint64(rawF.Variant))
		case 11:
			x.SetB(protowire.DecodeBool(rawF.Variant))
		case 12:
			x.SetStr(pbeditor.Bys2Str(rawF.Bytes))
		case 13:
			x.SetBys(rawF.Bytes)
		case 14:
			x.SetTe(TestEnum(int32(rawF.Variant)))
		case 15:
			x.SetF32(rawF.F32)
		case 16:
			x.SetF64(rawF.F64)
		case 17:
			x.GetMsg().MergeDirtyFromBytes(rawF.Bytes)
		case 18:
			tmp := rawF.Bytes
			dl := make([]int32, 0)
			for len(tmp) > 0 {
				val, n := protowire.ConsumeVarint(tmp)
				if n < 0 {
					panic(n)
				}
				tmp = tmp[n:]
				dl = append(dl, int32(val))
			}
			x.GetI32List().MergeDirtyFromPb(fds.MutState.ListOpLog(rawF.Number), dl)
		case 28:
			tmp := rawF.Bytes
			dl := make([]int64, 0)
			for len(tmp) > 0 {
				val, n := protowire.ConsumeVarint(tmp)
				if n < 0 {
					panic(n)
				}
				tmp = tmp[n:]
				dl = append(dl, int64(val))
			}
			x.GetI64List().MergeDirtyFromPb(fds.MutState.ListOpLog(rawF.Number), dl)
		case 38:
			tmp := rawF.Bytes
			dl := make([]uint32, 0)
			for len(tmp) > 0 {
				val, n := protowire.ConsumeVarint(tmp)
				if n < 0 {
					panic(n)
				}
				tmp = tmp[n:]
				dl = append(dl, uint32(val))
			}
			x.GetU32List().MergeDirtyFromPb(fds.MutState.ListOpLog(rawF.Number), dl)
		case 48:
			tmp := rawF.Bytes
			dl := make([]uint64, 0)
			for len(tmp) > 0 {
				val, n := protowire.ConsumeVarint(tmp)
				if n < 0 {
					panic(n)
				}
				tmp = tmp[n:]
				dl = append(dl, uint64(val))
			}
			x.GetU64List().MergeDirtyFromPb(fds.MutState.ListOpLog(rawF.Number), dl)
		case 118:
			tmp := rawF.Bytes
			dl := make([]bool, 0)
			for len(tmp) > 0 {
				val, n := protowire.ConsumeVarint(tmp)
				if n < 0 {
					panic(n)
				}
				tmp = tmp[n:]
				dl = append(dl, protowire.DecodeBool(val))
			}
			x.GetBList().MergeDirtyFromPb(fds.MutState.ListOpLog(rawF.Number), dl)
		case 128:
			sl := make([]string, 0)
			for _, field := range rawFields {
				sl = append(sl, pbeditor.Bys2Str(field.Bytes))
			}
			x.GetStrList().MergeDirtyFromPb(fds.MutState.ListOpLog(rawF.Number), sl)
		case 148:
			tmp := rawF.Bytes
			dl := make([]TestEnum, 0)
			for len(tmp) > 0 {
				val, n := protowire.ConsumeVarint(tmp)
				if n < 0 {
					panic(n)
				}
				tmp = tmp[n:]
				dl = append(dl, TestEnum(val))
			}
			x.GetTeList().MergeDirtyFromPb(fds.MutState.ListOpLog(rawF.Number), dl)
		case 158:
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
			x.GetF32List().MergeDirtyFromPb(fds.MutState.ListOpLog(rawF.Number), dl)
		case 168:
			tmp := rawF.Bytes
			dl := make([]float64, 0)
			for len(tmp) > 0 {
				val, n := protowire.ConsumeFixed64(tmp)
				if n < 0 {
					panic(n)
				}
				tmp = tmp[n:]
				dl = append(dl, math.Float64frombits(val))
			}
			x.GetF64List().MergeDirtyFromPb(fds.MutState.ListOpLog(rawF.Number), dl)
		case 178:
			sl := make([]*TestMsg, 0)
			for _, field := range rawFields {
				tmp := &TestMsg{}
				proto.Unmarshal(field.Bytes, tmp)
				sl = append(sl, tmp)
			}
			x.GetMsgList().MergeDirtyFromPb(fds.MutState.ListOpLog(rawF.Number), sl)
		case 19:
			r := make(map[int32]string)
			for _, field := range rawFields {
				fmt.Println(field.Bytes)
				tmp := pbeditor.ParseMap(field.Bytes)
				k := tmp.KeyI32()
				v := tmp.ValStr()
				r[k] = v
			}
			x.GetI32StrMap().MergeDirtyFromPb(fds.MutState.I32List(rawF.Number), r)
		case 29:
			r := make(map[string]float32)
			for _, field := range rawFields {
				tmp := pbeditor.ParseMap(field.Bytes)
				k := tmp.KeyStr()
				v := tmp.ValF32()
				r[k] = v
			}
			x.GetStrF32Map().MergeDirtyFromPb(fds.MutState.StrList(rawF.Number), r)
		case 39:
			r := make(map[bool]uint64)
			for _, field := range rawFields {
				tmp := pbeditor.ParseMap(field.Bytes)
				k := tmp.KeyBool()
				v := tmp.ValU64()
				r[k] = v
			}
			x.GetBoolU64Map().MergeDirtyFromPb(fds.MutState.BoolList(rawF.Number), r)
		case 49:
			r := make(map[uint64]float64)
			for _, field := range rawFields {
				tmp := pbeditor.ParseMap(field.Bytes)
				k := tmp.KeyU64()
				v := tmp.ValF64()
				r[k] = v
			}
			x.GetU64F64Map().MergeDirtyFromPb(fds.MutState.U64List(rawF.Number), r)
		case 59:
			r := make(map[int32]*TestMsg)
			for _, field := range rawFields {
				tmp := pbeditor.ParseMap(field.Bytes)
				k := tmp.KeyI32()
				v := &TestMsg{}
				proto.Unmarshal(field.Bytes, v)
				r[k] = v
			}
			x.GetMsgMap().MergeDirtyFromPb(fds.MutState.I32List(rawF.Number), r)
		case 69:
			r := make(map[string]*PlainMsg)
			for _, field := range rawFields {
				tmp := pbeditor.ParseMap(field.Bytes)
				k := tmp.KeyStr()
				v := &PlainMsg{}
				proto.Unmarshal(field.Bytes, v)
				r[k] = v
			}
			x.GetStrMsgMap().MergeDirtyFromPb(fds.MutState.StrList(rawF.Number), r)
		case 79:
			r := make(map[uint32]TestEnum)
			for _, field := range rawFields {
				tmp := pbeditor.ParseMap(field.Bytes)
				k := tmp.KeyU32()
				v := tmp.ValI32()
				r[k] = TestEnum(v)
			}
			x.GetEnumMap().MergeDirtyFromPb(fds.MutState.U32List(rawF.Number), r)
		}
	}
	return x
}
func (x *TestEditor) GatherDirtyToPb(r *Test) {
	mutState := pbeditor.UnmarshalMutState(r.unknownFields)
	if x.isI32Dirty() {
		r.SetI32(x.i32)
	}
	if x.isI64Dirty() {
		r.SetI64(x.i64)
	}
	if x.isU32Dirty() {
		r.SetU32(x.u32)
	}
	if x.isU64Dirty() {
		r.SetU64(x.u64)
	}
	if x.isBDirty() {
		r.SetB(x.b)
	}
	if x.isStrDirty() {
		r.SetStr(x.str)
	}
	if x.isBysDirty() {
		r.SetBys(x.bys)
	}
	if x.isBysDirty() {
		r.SetBys(slices.Clone(x.bys))
	}
	if x.isTeDirty() {
		r.SetTe(x.te)
	}
	if x.isF32Dirty() {
		r.SetF32(x.f32)
	}
	if x.isF64Dirty() {
		r.SetF64(x.f64)
	}
	if x.isMsgDirty() && x.msg != nil {
		if r.Msg == nil {
			r.Msg = &TestMsg{}
		}
		x.msg.GatherDirtyToPb(r.Msg)
	}
	if x.isI32ListDirty() {
		ori := r.I32List
		lenOfI32List := len(ori)
		logsOfI32List := x.i32List.ChangeLog()
		dirtiesOfI32List := x.i32List.MonitorDirtyList()
		for _, val := range dirtiesOfI32List {
			r.I32List = append(r.I32List, val)
		}
		for _, log := range logsOfI32List {
			log.Pos = int32(lenOfI32List) + log.Pos
			mutState.ListAppendLog(18, log)
		}
	}
	if x.isI64ListDirty() {
		ori := r.I64List
		lenOfI64List := len(ori)
		logsOfI64List := x.i64List.ChangeLog()
		dirtiesOfI64List := x.i64List.MonitorDirtyList()
		for _, val := range dirtiesOfI64List {
			r.I64List = append(r.I64List, val)
		}
		for _, log := range logsOfI64List {
			log.Pos = int32(lenOfI64List) + log.Pos
			mutState.ListAppendLog(28, log)
		}
	}
	if x.isU32ListDirty() {
		ori := r.U32List
		lenOfU32List := len(ori)
		logsOfU32List := x.u32List.ChangeLog()
		dirtiesOfU32List := x.u32List.MonitorDirtyList()
		for _, val := range dirtiesOfU32List {
			r.U32List = append(r.U32List, val)
		}
		for _, log := range logsOfU32List {
			log.Pos = int32(lenOfU32List) + log.Pos
			mutState.ListAppendLog(38, log)
		}
	}
	if x.isU64ListDirty() {
		ori := r.U64List
		lenOfU64List := len(ori)
		logsOfU64List := x.u64List.ChangeLog()
		dirtiesOfU64List := x.u64List.MonitorDirtyList()
		for _, val := range dirtiesOfU64List {
			r.U64List = append(r.U64List, val)
		}
		for _, log := range logsOfU64List {
			log.Pos = int32(lenOfU64List) + log.Pos
			mutState.ListAppendLog(48, log)
		}
	}
	if x.isBListDirty() {
		ori := r.BList
		lenOfBList := len(ori)
		logsOfBList := x.bList.ChangeLog()
		dirtiesOfBList := x.bList.MonitorDirtyList()
		for _, val := range dirtiesOfBList {
			r.BList = append(r.BList, val)
		}
		for _, log := range logsOfBList {
			log.Pos = int32(lenOfBList) + log.Pos
			mutState.ListAppendLog(118, log)
		}
	}
	if x.isStrListDirty() {
		ori := r.StrList
		lenOfStrList := len(ori)
		logsOfStrList := x.strList.ChangeLog()
		dirtiesOfStrList := x.strList.MonitorDirtyList()
		for _, val := range dirtiesOfStrList {
			r.StrList = append(r.StrList, val)
		}
		for _, log := range logsOfStrList {
			log.Pos = int32(lenOfStrList) + log.Pos
			mutState.ListAppendLog(128, log)
		}
	}
	if x.isTeListDirty() {
		ori := r.TeList
		lenOfTeList := len(ori)
		logsOfTeList := x.teList.ChangeLog()
		dirtiesOfTeList := x.teList.MonitorDirtyList()
		for _, val := range dirtiesOfTeList {
			r.TeList = append(r.TeList, val)
		}
		for _, log := range logsOfTeList {
			log.Pos = int32(lenOfTeList) + log.Pos
			mutState.ListAppendLog(148, log)
		}
	}
	if x.isF32ListDirty() {
		ori := r.F32List
		lenOfF32List := len(ori)
		logsOfF32List := x.f32List.ChangeLog()
		dirtiesOfF32List := x.f32List.MonitorDirtyList()
		for _, val := range dirtiesOfF32List {
			r.F32List = append(r.F32List, val)
		}
		for _, log := range logsOfF32List {
			log.Pos = int32(lenOfF32List) + log.Pos
			mutState.ListAppendLog(158, log)
		}
	}
	if x.isF64ListDirty() {
		ori := r.F64List
		lenOfF64List := len(ori)
		logsOfF64List := x.f64List.ChangeLog()
		dirtiesOfF64List := x.f64List.MonitorDirtyList()
		for _, val := range dirtiesOfF64List {
			r.F64List = append(r.F64List, val)
		}
		for _, log := range logsOfF64List {
			log.Pos = int32(lenOfF64List) + log.Pos
			mutState.ListAppendLog(168, log)
		}
	}
	if x.isMsgListDirty() {
		ori := r.MsgList
		lenOfMsgList := len(ori)
		logsOfMsgList := x.msgList.ChangeLog()
		dirtiesOfMsgList := x.msgList.MonitorDirtyList()
		for _, val := range dirtiesOfMsgList {
			o := &TestMsg{}
			val.GatherDirtyToPb(o)
			r.MsgList = append(r.MsgList, o)
		}
		for _, log := range logsOfMsgList {
			log.Pos = int32(lenOfMsgList) + log.Pos
			mutState.ListAppendLog(178, log)
		}
	}
	if x.isI32StrMapDirty() {
		updated := make([]int32, 0)
		if r.I32StrMap != nil {
			for k, v := range r.I32StrMap {
				if x.i32StrMap.ContainDeleted(k) {
					delete(r.I32StrMap, k)
					continue
				}
				if x.i32StrMap.ContainDirtied(k) {
					updated = append(updated, k)
					v, _ = x.i32StrMap.Get(k)
					r.I32StrMap[k] = v
				}
			}
		} else {
			r.I32StrMap = make(map[int32]string)
		}
		for k := range x.i32StrMap.Dirtied() {
			if !slices.Contains(updated, k) {
				tmp, _ := x.i32StrMap.Get(k)
				r.I32StrMap[k] = tmp
			}
		}
		deleted := mutState.I32List(19)
		for k := range x.i32StrMap.Deleted() {
			if !slices.Contains(deleted, k) {
				mutState.AddI32(19, k)
			}
		}
	}
	if x.isStrF32MapDirty() {
		updated := make([]string, 0)
		if r.StrF32Map != nil {
			for k, v := range r.StrF32Map {
				if x.strF32Map.ContainDeleted(k) {
					delete(r.StrF32Map, k)
					continue
				}
				if x.strF32Map.ContainDirtied(k) {
					updated = append(updated, k)
					v, _ = x.strF32Map.Get(k)
					r.StrF32Map[k] = v
				}
			}
		} else {
			r.StrF32Map = make(map[string]float32)
		}
		for k := range x.strF32Map.Dirtied() {
			if !slices.Contains(updated, k) {
				tmp, _ := x.strF32Map.Get(k)
				r.StrF32Map[k] = tmp
			}
		}
		deleted := mutState.StrList(29)
		for k := range x.strF32Map.Deleted() {
			if !slices.Contains(deleted, k) {
				mutState.AddStr(29, k)
			}
		}
	}
	if x.isBoolU64MapDirty() {
		updated := make([]bool, 0)
		if r.BoolU64Map != nil {
			for k, v := range r.BoolU64Map {
				if x.boolU64Map.ContainDeleted(k) {
					delete(r.BoolU64Map, k)
					continue
				}
				if x.boolU64Map.ContainDirtied(k) {
					updated = append(updated, k)
					v, _ = x.boolU64Map.Get(k)
					r.BoolU64Map[k] = v
				}
			}
		} else {
			r.BoolU64Map = make(map[bool]uint64)
		}
		for k := range x.boolU64Map.Dirtied() {
			if !slices.Contains(updated, k) {
				tmp, _ := x.boolU64Map.Get(k)
				r.BoolU64Map[k] = tmp
			}
		}
		deleted := mutState.BoolList(39)
		for k := range x.boolU64Map.Deleted() {
			if !slices.Contains(deleted, k) {
				mutState.AddBool(39, k)
			}
		}
	}
	if x.isU64F64MapDirty() {
		updated := make([]uint64, 0)
		if r.U64F64Map != nil {
			for k, v := range r.U64F64Map {
				if x.u64F64Map.ContainDeleted(k) {
					delete(r.U64F64Map, k)
					continue
				}
				if x.u64F64Map.ContainDirtied(k) {
					updated = append(updated, k)
					v, _ = x.u64F64Map.Get(k)
					r.U64F64Map[k] = v
				}
			}
		} else {
			r.U64F64Map = make(map[uint64]float64)
		}
		for k := range x.u64F64Map.Dirtied() {
			if !slices.Contains(updated, k) {
				tmp, _ := x.u64F64Map.Get(k)
				r.U64F64Map[k] = tmp
			}
		}
		deleted := mutState.U64List(49)
		for k := range x.u64F64Map.Deleted() {
			if !slices.Contains(deleted, k) {
				mutState.AddU64(49, k)
			}
		}
	}
	if x.isMsgMapDirty() {
		updated := make([]int32, 0)
		if r.MsgMap != nil {
			for k, v := range r.MsgMap {
				if x.msgMap.ContainDeleted(k) {
					delete(r.MsgMap, k)
					continue
				}
				if x.msgMap.ContainDirtied(k) {
					updated = append(updated, k)
					tmp := x.msgMap.Get(k)
					if tmp == nil {
						r.MsgMap[k] = nil
						continue
					}
					if v == nil {
						v = &TestMsg{}
						tmp.GatherDirtyToPb(v)
						r.MsgMap[k] = v

					} else {
						tmp.GatherDirtyToPb(v)
					}
				}
			}
		} else {
			r.MsgMap = make(map[int32]*TestMsg)
		}
		for k := range x.msgMap.Dirtied() {
			if !slices.Contains(updated, k) {
				tmp := x.msgMap.Get(k)
				if tmp == nil {
					r.MsgMap[k] = nil
				} else {
					v := &TestMsg{}
					tmp.GatherDirtyToPb(v)
					r.MsgMap[k] = v
				}
			}
		}
		deleted := mutState.I32List(59)
		for k := range x.msgMap.Deleted() {
			if !slices.Contains(deleted, k) {
				mutState.AddI32(59, k)
			}
		}
	}
	if x.isStrMsgMapDirty() {
		updated := make([]string, 0)
		if r.StrMsgMap != nil {
			for k, v := range r.StrMsgMap {
				if x.strMsgMap.ContainDeleted(k) {
					delete(r.StrMsgMap, k)
					continue
				}
				if x.strMsgMap.ContainDirtied(k) {
					updated = append(updated, k)
					tmp := x.strMsgMap.Get(k)
					if tmp == nil {
						r.StrMsgMap[k] = nil
						continue
					}
					if v == nil {
						v = &PlainMsg{}
						tmp.GatherDirtyToPb(v)
						r.StrMsgMap[k] = v

					} else {
						tmp.GatherDirtyToPb(v)
					}
				}
			}
		} else {
			r.StrMsgMap = make(map[string]*PlainMsg)
		}
		for k := range x.strMsgMap.Dirtied() {
			if !slices.Contains(updated, k) {
				tmp := x.strMsgMap.Get(k)
				if tmp == nil {
					r.StrMsgMap[k] = nil
				} else {
					v := &PlainMsg{}
					tmp.GatherDirtyToPb(v)
					r.StrMsgMap[k] = v
				}
			}
		}
		deleted := mutState.StrList(69)
		for k := range x.strMsgMap.Deleted() {
			if !slices.Contains(deleted, k) {
				mutState.AddStr(69, k)
			}
		}
	}
	if x.isEnumMapDirty() {
		updated := make([]uint32, 0)
		if r.EnumMap != nil {
			for k, v := range r.EnumMap {
				if x.enumMap.ContainDeleted(k) {
					delete(r.EnumMap, k)
					continue
				}
				if x.enumMap.ContainDirtied(k) {
					updated = append(updated, k)
					v, _ = x.enumMap.Get(k)
					r.EnumMap[k] = v
				}
			}
		} else {
			r.EnumMap = make(map[uint32]TestEnum)
		}
		for k := range x.enumMap.Dirtied() {
			if !slices.Contains(updated, k) {
				tmp, _ := x.enumMap.Get(k)
				r.EnumMap[k] = tmp
			}
		}
		deleted := mutState.U32List(79)
		for k := range x.enumMap.Deleted() {
			if !slices.Contains(deleted, k) {
				mutState.AddU32(79, k)
			}
		}
	}
	r.ProtoReflect().SetUnknown(mutState.MarshalToFieldBytes())
}
func (x *TestEditor) GatherDirtyToBytes() []byte {
	mutState := pbeditor.CreateMutStateInfo()
	var buf []byte
	if x.isI32Dirty() {
		buf = protowire.AppendTag(buf, 1, 0)
		buf = protowire.AppendVarint(buf, uint64(x.i32))
	}
	if x.isI64Dirty() {
		buf = protowire.AppendTag(buf, 2, 0)
		buf = protowire.AppendVarint(buf, uint64(x.i64))
	}
	if x.isU32Dirty() {
		buf = protowire.AppendTag(buf, 3, 0)
		buf = protowire.AppendVarint(buf, uint64(x.u32))
	}
	if x.isU64Dirty() {
		buf = protowire.AppendTag(buf, 4, 0)
		buf = protowire.AppendVarint(buf, uint64(x.u64))
	}
	if x.isBDirty() {
		buf = protowire.AppendTag(buf, 11, 0)
		var t uint64 = 0
		if x.b {
			t = 1
		}
		buf = protowire.AppendVarint(buf, t)
	}
	if x.isStrDirty() {
		buf = protowire.AppendTag(buf, 12, 2)
		buf = protowire.AppendString(buf, x.str)
	}
	if x.isBysDirty() {
		buf = protowire.AppendTag(buf, 13, 2)
		buf = protowire.AppendBytes(buf, x.bys)
	}
	if x.isTeDirty() {
		buf = protowire.AppendTag(buf, 14, 0)
		buf = protowire.AppendVarint(buf, uint64(x.te))
	}
	if x.isF32Dirty() {
		buf = protowire.AppendTag(buf, 15, 5)
		buf = protowire.AppendFixed32(buf, math.Float32bits(x.f32))
	}
	if x.isF64Dirty() {
		buf = protowire.AppendTag(buf, 16, 1)
		buf = protowire.AppendFixed64(buf, math.Float64bits(x.f64))
	}
	if x.isMsgDirty() && x.msg != nil {
		buf = protowire.AppendTag(buf, 17, protowire.BytesType)
		buf = protowire.AppendBytes(buf, x.msg.GatherDirtyToBytes())
	}
	if x.isI32ListDirty() {
		for _, log := range x.i32List.ChangeLog() {
			mutState.ListAppendLog(18, log)
		}
		tmp := make([]byte, 0)
		for _, v := range x.i32List.MonitorDirtyList() {
			tmp = protowire.AppendVarint(tmp, uint64(v))
		}
		buf = protowire.AppendTag(buf, 18, protowire.BytesType)
		buf = protowire.AppendBytes(buf, tmp)
	}
	if x.isI64ListDirty() {
		for _, log := range x.i64List.ChangeLog() {
			mutState.ListAppendLog(28, log)
		}
		tmp := make([]byte, 0)
		for _, v := range x.i64List.MonitorDirtyList() {
			tmp = protowire.AppendVarint(tmp, uint64(v))
		}
		buf = protowire.AppendTag(buf, 28, protowire.BytesType)
		buf = protowire.AppendBytes(buf, tmp)
	}
	if x.isU32ListDirty() {
		for _, log := range x.u32List.ChangeLog() {
			mutState.ListAppendLog(38, log)
		}
		tmp := make([]byte, 0)
		for _, v := range x.u32List.MonitorDirtyList() {
			tmp = protowire.AppendVarint(tmp, uint64(v))
		}
		buf = protowire.AppendTag(buf, 38, protowire.BytesType)
		buf = protowire.AppendBytes(buf, tmp)
	}
	if x.isU64ListDirty() {
		for _, log := range x.u64List.ChangeLog() {
			mutState.ListAppendLog(48, log)
		}
		tmp := make([]byte, 0)
		for _, v := range x.u64List.MonitorDirtyList() {
			tmp = protowire.AppendVarint(tmp, uint64(v))
		}
		buf = protowire.AppendTag(buf, 48, protowire.BytesType)
		buf = protowire.AppendBytes(buf, tmp)
	}
	if x.isBListDirty() {
		for _, log := range x.bList.ChangeLog() {
			mutState.ListAppendLog(118, log)
		}
		tmp := make([]byte, 0)
		for _, b := range x.bList.MonitorDirtyList() {
			var t uint64 = 0
			if b {
				t = 1
			}
			tmp = protowire.AppendVarint(tmp, t)
		}
		buf = protowire.AppendTag(buf, 118, protowire.BytesType)
		buf = protowire.AppendBytes(buf, tmp)
	}
	if x.isStrListDirty() {
		for _, log := range x.strList.ChangeLog() {
			mutState.ListAppendLog(128, log)
		}
		for _, s := range x.strList.MonitorDirtyList() {
			buf = protowire.AppendTag(buf, 128, protowire.BytesType)
			buf = protowire.AppendString(buf, s)
		}
	}
	if x.isTeListDirty() {
		for _, log := range x.teList.ChangeLog() {
			mutState.ListAppendLog(148, log)
		}
		tmp := make([]byte, 0)
		for _, v := range x.teList.MonitorDirtyList() {
			tmp = protowire.AppendVarint(tmp, uint64(v))
		}
		buf = protowire.AppendTag(buf, 148, protowire.BytesType)
		buf = protowire.AppendBytes(buf, tmp)
	}
	if x.isF32ListDirty() {
		for _, log := range x.f32List.ChangeLog() {
			mutState.ListAppendLog(158, log)
		}
		tmp := make([]byte, 0)
		for _, f := range x.f32List.MonitorDirtyList() {
			tmp = protowire.AppendFixed32(tmp, math.Float32bits(f))
		}
		buf = protowire.AppendTag(buf, 158, protowire.BytesType)
		buf = protowire.AppendBytes(buf, tmp)
	}
	if x.isF64ListDirty() {
		for _, log := range x.f64List.ChangeLog() {
			mutState.ListAppendLog(168, log)
		}
		tmp := make([]byte, 0)
		for _, f := range x.f64List.MonitorDirtyList() {
			tmp = protowire.AppendFixed64(tmp, math.Float64bits(f))
		}
		buf = protowire.AppendTag(buf, 168, protowire.BytesType)
		buf = protowire.AppendBytes(buf, tmp)
	}
	if x.isMsgListDirty() {
		for _, log := range x.msgList.ChangeLog() {
			mutState.ListAppendLog(178, log)
		}
		for _, e := range x.msgList.MonitorDirtyList() {
			buf = protowire.AppendTag(buf, 178, protowire.BytesType)
			buf = protowire.AppendBytes(buf, e.GatherDirtyToBytes())
		}
	}
	if x.isI32StrMapDirty() {
		if len(x.i32StrMap.Deleted()) > 0 {
			deleted := x.i32StrMap.Deleted()
			for v := range deleted {
				mutState.AddI32(19, v)
			}
		}
		if x.i32StrMap.Len() > 0 {
			for k, v := range x.i32StrMap.Iter() {
				if !x.i32StrMap.ContainDirtied(k) {
					continue
				}
				var tmp []byte
				tmp = protowire.AppendTag(tmp, 1, protowire.VarintType)
				tmp = protowire.AppendVarint(tmp, uint64(k))
				tmp = protowire.AppendTag(tmp, 2, protowire.BytesType)
				tmp = protowire.AppendString(tmp, v)
				buf = protowire.AppendTag(buf, 19, protowire.BytesType)
				buf = protowire.AppendBytes(buf, tmp)
			}
		}
	}
	if x.isStrF32MapDirty() {
		if len(x.strF32Map.Deleted()) > 0 {
			deleted := x.strF32Map.Deleted()
			for v := range deleted {
				mutState.AddStr(29, v)
			}
		}
		if x.strF32Map.Len() > 0 {
			for k, v := range x.strF32Map.Iter() {
				if !x.strF32Map.ContainDirtied(k) {
					continue
				}
				var tmp []byte
				tmp = protowire.AppendTag(tmp, 1, protowire.BytesType)
				tmp = protowire.AppendString(tmp, k)
				tmp = protowire.AppendTag(tmp, 2, protowire.Fixed32Type)
				tmp = protowire.AppendFixed32(tmp, math.Float32bits(v))
				buf = protowire.AppendTag(buf, 29, protowire.BytesType)
				buf = protowire.AppendBytes(buf, tmp)
			}
		}
	}
	if x.isBoolU64MapDirty() {
		if len(x.boolU64Map.Deleted()) > 0 {
			deleted := x.boolU64Map.Deleted()
			for v := range deleted {
				mutState.AddBool(39, v)
			}
		}
		if x.boolU64Map.Len() > 0 {
			for k, v := range x.boolU64Map.Iter() {
				if !x.boolU64Map.ContainDirtied(k) {
					continue
				}
				var tmp []byte
				var t uint64 = 0
				if k {
					t = 1
				}
				tmp = protowire.AppendTag(tmp, 1, protowire.VarintType)
				tmp = protowire.AppendVarint(tmp, t)
				tmp = protowire.AppendTag(tmp, 2, protowire.VarintType)
				tmp = protowire.AppendVarint(tmp, uint64(v))
				buf = protowire.AppendTag(buf, 39, protowire.BytesType)
				buf = protowire.AppendBytes(buf, tmp)
			}
		}
	}
	if x.isU64F64MapDirty() {
		if len(x.u64F64Map.Deleted()) > 0 {
			deleted := x.u64F64Map.Deleted()
			for v := range deleted {
				mutState.AddU64(49, v)
			}
		}
		if x.u64F64Map.Len() > 0 {
			for k, v := range x.u64F64Map.Iter() {
				if !x.u64F64Map.ContainDirtied(k) {
					continue
				}
				var tmp []byte
				tmp = protowire.AppendTag(tmp, 1, protowire.VarintType)
				tmp = protowire.AppendVarint(tmp, uint64(k))
				tmp = protowire.AppendTag(tmp, 2, protowire.Fixed64Type)
				tmp = protowire.AppendFixed64(tmp, math.Float64bits(v))
				buf = protowire.AppendTag(buf, 49, protowire.BytesType)
				buf = protowire.AppendBytes(buf, tmp)
			}
		}
	}
	if x.isMsgMapDirty() {
		if len(x.msgMap.Deleted()) > 0 {
			deleted := x.msgMap.Deleted()
			for v := range deleted {
				mutState.AddI32(59, v)
			}
		}
		if x.msgMap.Len() > 0 {
			for k, v := range x.msgMap.Iter() {
				if !x.msgMap.ContainDirtied(k) {
					continue
				}
				var tmp []byte
				tmp = protowire.AppendTag(tmp, 1, protowire.VarintType)
				tmp = protowire.AppendVarint(tmp, uint64(k))
				vb := v.GatherDirtyToBytes()
				tmp = protowire.AppendTag(tmp, 2, protowire.BytesType)
				tmp = protowire.AppendBytes(tmp, vb)
				buf = protowire.AppendTag(buf, 59, protowire.BytesType)
				buf = protowire.AppendBytes(buf, tmp)
			}
		}
	}
	if x.isStrMsgMapDirty() {
		if len(x.strMsgMap.Deleted()) > 0 {
			deleted := x.strMsgMap.Deleted()
			for v := range deleted {
				mutState.AddStr(69, v)
			}
		}
		if x.strMsgMap.Len() > 0 {
			for k, v := range x.strMsgMap.Iter() {
				if !x.strMsgMap.ContainDirtied(k) {
					continue
				}
				var tmp []byte
				tmp = protowire.AppendTag(tmp, 1, protowire.BytesType)
				tmp = protowire.AppendString(tmp, k)
				vb := v.GatherDirtyToBytes()
				tmp = protowire.AppendTag(tmp, 2, protowire.BytesType)
				tmp = protowire.AppendBytes(tmp, vb)
				buf = protowire.AppendTag(buf, 69, protowire.BytesType)
				buf = protowire.AppendBytes(buf, tmp)
			}
		}
	}
	if x.isEnumMapDirty() {
		if len(x.enumMap.Deleted()) > 0 {
			deleted := x.enumMap.Deleted()
			for v := range deleted {
				mutState.AddU32(79, v)
			}
		}
		if x.enumMap.Len() > 0 {
			for k, v := range x.enumMap.Iter() {
				if !x.enumMap.ContainDirtied(k) {
					continue
				}
				var tmp []byte
				tmp = protowire.AppendTag(tmp, 1, protowire.VarintType)
				tmp = protowire.AppendVarint(tmp, uint64(k))
				tmp = protowire.AppendTag(tmp, 2, protowire.VarintType)
				tmp = protowire.AppendVarint(tmp, uint64(v))
				buf = protowire.AppendTag(buf, 79, protowire.BytesType)
				buf = protowire.AppendBytes(buf, tmp)
			}
		}
	}
	buf = append(buf, mutState.MarshalToFieldBytes()...)
	return buf
}
func (x *TestEditor) GetI32() int32 {
	return x.i32
}
func (x *TestEditor) SetI32(v int32) *TestEditor {
	if x.i32 == v {
		return x
	}
	x.i32 = v
	x.setI32Dirty(true, false)
	return x
}
func (x *TestEditor) GetI64() int64 {
	return x.i64
}
func (x *TestEditor) SetI64(v int64) *TestEditor {
	if x.i64 == v {
		return x
	}
	x.i64 = v
	x.setI64Dirty(true, false)
	return x
}
func (x *TestEditor) GetU32() uint32 {
	return x.u32
}
func (x *TestEditor) SetU32(v uint32) *TestEditor {
	if x.u32 == v {
		return x
	}
	x.u32 = v
	x.setU32Dirty(true, false)
	return x
}
func (x *TestEditor) GetU64() uint64 {
	return x.u64
}
func (x *TestEditor) SetU64(v uint64) *TestEditor {
	if x.u64 == v {
		return x
	}
	x.u64 = v
	x.setU64Dirty(true, false)
	return x
}
func (x *TestEditor) GetB() bool {
	return x.b
}
func (x *TestEditor) SetB(v bool) *TestEditor {
	if x.b == v {
		return x
	}
	x.b = v
	x.setBDirty(true, false)
	return x
}
func (x *TestEditor) GetStr() string {
	return x.str
}
func (x *TestEditor) SetStr(v string) *TestEditor {
	if x.str == v {
		return x
	}
	x.str = v
	x.setStrDirty(true, false)
	return x
}
func (x *TestEditor) GetBys() []byte {
	return x.bys
}
func (x *TestEditor) SetBys(v []byte) *TestEditor {
	x.bys = v
	x.setBysDirty(true, false)
	return x
}
func (x *TestEditor) GetTe() TestEnum {
	return x.te
}
func (x *TestEditor) SetTe(v TestEnum) *TestEditor {
	if x.te == v {
		return x
	}
	x.te = v
	x.setTeDirty(true, false)
	return x
}
func (x *TestEditor) GetF32() float32 {
	return x.f32
}
func (x *TestEditor) SetF32(v float32) *TestEditor {
	if x.f32 == v {
		return x
	}
	x.f32 = v
	x.setF32Dirty(true, false)
	return x
}
func (x *TestEditor) GetF64() float64 {
	return x.f64
}
func (x *TestEditor) SetF64(v float64) *TestEditor {
	if x.f64 == v {
		return x
	}
	x.f64 = v
	x.setF64Dirty(true, false)
	return x
}
func (x *TestEditor) GetMsg() *TestMsgEditor {
	if x.msg == nil {
		x.msg = NewTestMsgEditor()
		x.msg.SetParent(x, 17)
	}
	return x.msg
}
func (x *TestEditor) SetMsg(v *TestMsgEditor) *TestEditor {
	if v != nil {
		v.SetParent(x, 17)
	}
	if x.msg != nil {
		x.msg.SetParent(nil, -1)
	}
	x.msg = v
	x.setMsgDirty(true, false)
	return x
}
func (x *TestEditor) GetI32List() *pbeditor.I32List {
	if x.i32List == nil {
		x.i32List = pbeditor.NewI32List()
		x.i32List.SetParent(x, 18)
	}
	return x.i32List
}
func (x *TestEditor) AddI32List(v int32) *TestEditor {
	if x.i32List == nil {
		x.i32List = pbeditor.NewI32List()
		x.i32List.SetParent(x, 18)
	}
	x.i32List.Add(v)
	return x
}
func (x *TestEditor) GetI64List() *pbeditor.I64List {
	if x.i64List == nil {
		x.i64List = pbeditor.NewI64List()
		x.i64List.SetParent(x, 28)
	}
	return x.i64List
}
func (x *TestEditor) AddI64List(v int64) *TestEditor {
	if x.i64List == nil {
		x.i64List = pbeditor.NewI64List()
		x.i64List.SetParent(x, 28)
	}
	x.i64List.Add(v)
	return x
}
func (x *TestEditor) GetU32List() *pbeditor.U32List {
	if x.u32List == nil {
		x.u32List = pbeditor.NewU32List()
		x.u32List.SetParent(x, 38)
	}
	return x.u32List
}
func (x *TestEditor) AddU32List(v uint32) *TestEditor {
	if x.u32List == nil {
		x.u32List = pbeditor.NewU32List()
		x.u32List.SetParent(x, 38)
	}
	x.u32List.Add(v)
	return x
}
func (x *TestEditor) GetU64List() *pbeditor.U64List {
	if x.u64List == nil {
		x.u64List = pbeditor.NewU64List()
		x.u64List.SetParent(x, 48)
	}
	return x.u64List
}
func (x *TestEditor) AddU64List(v uint64) *TestEditor {
	if x.u64List == nil {
		x.u64List = pbeditor.NewU64List()
		x.u64List.SetParent(x, 48)
	}
	x.u64List.Add(v)
	return x
}
func (x *TestEditor) GetBList() *pbeditor.BoolList {
	if x.bList == nil {
		x.bList = pbeditor.NewBoolList()
		x.bList.SetParent(x, 118)
	}
	return x.bList
}
func (x *TestEditor) AddBList(v bool) *TestEditor {
	if x.bList == nil {
		x.bList = pbeditor.NewBoolList()
		x.bList.SetParent(x, 118)
	}
	x.bList.Add(v)
	return x
}
func (x *TestEditor) GetStrList() *pbeditor.StrList {
	if x.strList == nil {
		x.strList = pbeditor.NewStrList()
		x.strList.SetParent(x, 128)
	}
	return x.strList
}
func (x *TestEditor) AddStrList(v string) *TestEditor {
	if x.strList == nil {
		x.strList = pbeditor.NewStrList()
		x.strList.SetParent(x, 128)
	}
	x.strList.Add(v)
	return x
}
func (x *TestEditor) GetTeList() *TestEnumEditorList {
	if x.teList == nil {
		x.teList = NewTestEnumEditorList()
		x.teList.SetParent(x, 148)
	}
	return x.teList
}
func (x *TestEditor) AddTeList(v TestEnum) *TestEditor {
	if x.teList == nil {
		x.teList = NewTestEnumEditorList()
		x.teList.SetParent(x, 148)
	}
	x.teList.Add(v)
	return x
}
func (x *TestEditor) GetF32List() *pbeditor.F32List {
	if x.f32List == nil {
		x.f32List = pbeditor.NewF32List()
		x.f32List.SetParent(x, 158)
	}
	return x.f32List
}
func (x *TestEditor) AddF32List(v float32) *TestEditor {
	if x.f32List == nil {
		x.f32List = pbeditor.NewF32List()
		x.f32List.SetParent(x, 158)
	}
	x.f32List.Add(v)
	return x
}
func (x *TestEditor) GetF64List() *pbeditor.F64List {
	if x.f64List == nil {
		x.f64List = pbeditor.NewF64List()
		x.f64List.SetParent(x, 168)
	}
	return x.f64List
}
func (x *TestEditor) AddF64List(v float64) *TestEditor {
	if x.f64List == nil {
		x.f64List = pbeditor.NewF64List()
		x.f64List.SetParent(x, 168)
	}
	x.f64List.Add(v)
	return x
}
func (x *TestEditor) GetMsgList() *TestMsgEditorList {
	if x.msgList == nil {
		x.msgList = NewTestMsgEditorList()
		x.msgList.SetParent(x, 178)
	}
	return x.msgList
}
func (x *TestEditor) AddMsgList(v *TestMsgEditor) *TestEditor {
	if x.msgList == nil {
		x.msgList = NewTestMsgEditorList()
		x.msgList.SetParent(x, 178)
	}
	x.msgList.Add(v)
	return x
}
func (x *TestEditor) GetI32StrMap() *pbeditor.I32StrMap {
	if x.i32StrMap == nil {
		x.i32StrMap = pbeditor.NewI32StrMap()
		x.i32StrMap.SetParent(x, 19)
	}
	return x.i32StrMap
}
func (x *TestEditor) GetStrF32Map() *pbeditor.StrF32Map {
	if x.strF32Map == nil {
		x.strF32Map = pbeditor.NewStrF32Map()
		x.strF32Map.SetParent(x, 29)
	}
	return x.strF32Map
}
func (x *TestEditor) GetBoolU64Map() *pbeditor.BoolU64Map {
	if x.boolU64Map == nil {
		x.boolU64Map = pbeditor.NewBoolU64Map()
		x.boolU64Map.SetParent(x, 39)
	}
	return x.boolU64Map
}
func (x *TestEditor) GetU64F64Map() *pbeditor.U64F64Map {
	if x.u64F64Map == nil {
		x.u64F64Map = pbeditor.NewU64F64Map()
		x.u64F64Map.SetParent(x, 49)
	}
	return x.u64F64Map
}
func (x *TestEditor) GetMsgMap() *I32TestMsgEditorMap {
	if x.msgMap == nil {
		x.msgMap = NewI32TestMsgEditorMap()
		x.msgMap.SetParent(x, 59)
	}
	return x.msgMap
}
func (x *TestEditor) GetStrMsgMap() *StrPlainMsgEditorMap {
	if x.strMsgMap == nil {
		x.strMsgMap = NewStrPlainMsgEditorMap()
		x.strMsgMap.SetParent(x, 69)
	}
	return x.strMsgMap
}
func (x *TestEditor) GetEnumMap() *U32TestEnumEditorMap {
	if x.enumMap == nil {
		x.enumMap = NewU32TestEnumEditorMap()
		x.enumMap.SetParent(x, 79)
	}
	return x.enumMap
}
func (xs *Test) SetI32(v int32) {
	xs.I32 = &v
}
func (xs *Test) SetI64(v int64) {
	xs.I64 = &v
}
func (xs *Test) SetU32(v uint32) {
	xs.U32 = &v
}
func (xs *Test) SetU64(v uint64) {
	xs.U64 = &v
}
func (xs *Test) SetB(v bool) {
	xs.B = &v
}
func (xs *Test) SetStr(v string) {
	xs.Str = &v
}
func (xs *Test) SetBys(v []byte) {
	xs.Bys = v
}
func (xs *Test) SetTe(v TestEnum) {
	xs.Te = &v
}
func (xs *Test) SetF32(v float32) {
	xs.F32 = &v
}
func (xs *Test) SetF64(v float64) {
	xs.F64 = &v
}
func (xs *Test) SetMsg(v *TestMsg) {
	xs.Msg = v
}
func (xs *Test) SetI32List(v []int32) {
	xs.I32List = v
}
func (xs *Test) AddI32List(v int32) {
	xs.I32List = append(xs.I32List, v)
}
func (xs *Test) SetI64List(v []int64) {
	xs.I64List = v
}
func (xs *Test) AddI64List(v int64) {
	xs.I64List = append(xs.I64List, v)
}
func (xs *Test) SetU32List(v []uint32) {
	xs.U32List = v
}
func (xs *Test) AddU32List(v uint32) {
	xs.U32List = append(xs.U32List, v)
}
func (xs *Test) SetU64List(v []uint64) {
	xs.U64List = v
}
func (xs *Test) AddU64List(v uint64) {
	xs.U64List = append(xs.U64List, v)
}
func (xs *Test) SetBList(v []bool) {
	xs.BList = v
}
func (xs *Test) AddBList(v bool) {
	xs.BList = append(xs.BList, v)
}
func (xs *Test) SetStrList(v []string) {
	xs.StrList = v
}
func (xs *Test) AddStrList(v string) {
	xs.StrList = append(xs.StrList, v)
}
func (xs *Test) SetTeList(v []TestEnum) {
	xs.TeList = v
}
func (xs *Test) AddTeList(v TestEnum) {
	xs.TeList = append(xs.TeList, v)
}
func (xs *Test) SetF32List(v []float32) {
	xs.F32List = v
}
func (xs *Test) AddF32List(v float32) {
	xs.F32List = append(xs.F32List, v)
}
func (xs *Test) SetF64List(v []float64) {
	xs.F64List = v
}
func (xs *Test) AddF64List(v float64) {
	xs.F64List = append(xs.F64List, v)
}
func (xs *Test) SetMsgList(v []*TestMsg) {
	xs.MsgList = v
}
func (xs *Test) AddMsgList(v *TestMsg) {
	xs.MsgList = append(xs.MsgList, v)
}
func (xs *Test) SetI32StrMap(v map[int32]string) {
	xs.I32StrMap = v
}
func (xs *Test) SetStrF32Map(v map[string]float32) {
	xs.StrF32Map = v
}
func (xs *Test) SetBoolU64Map(v map[bool]uint64) {
	xs.BoolU64Map = v
}
func (xs *Test) SetU64F64Map(v map[uint64]float64) {
	xs.U64F64Map = v
}
func (xs *Test) SetMsgMap(v map[int32]*TestMsg) {
	xs.MsgMap = v
}
func (xs *Test) SetStrMsgMap(v map[string]*PlainMsg) {
	xs.StrMsgMap = v
}
func (xs *Test) SetEnumMap(v map[uint32]TestEnum) {
	xs.EnumMap = v
}

type TestMsgEditorList struct {
	parent   pbeditor.Editor
	idxInPar any
	value    []*TestMsgEditor
	monitor  []*TestMsgEditor
	change   []*protodef.ListOpLog
}

func NewTestMsgEditorList() *TestMsgEditorList {
	return &TestMsgEditorList{value: make([]*TestMsgEditor, 0)}
}

func (ms *TestMsgEditorList) SetParent(editor pbeditor.Editor, key any) {
	ms.parent = editor
	ms.idxInPar = key
}
func (ms *TestMsgEditorList) SetParentDirty() {
	if ms.parent != nil {
		ms.parent.SetDirty(ms.idxInPar, true, ms)
	}
}
func (ms *TestMsgEditorList) SetDirty(key any, dirty bool, sync pbeditor.Editor) {
	var pos = -1
	if key != nil {
		pos = key.(int)
	} else {
		pos = slices.Index(ms.value, sync.(*TestMsgEditor))
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
func (ms *TestMsgEditorList) FlushDirty(dirty bool) {
	if !dirty {
		ms.monitor = nil
		ms.change = nil
	}
	for _, v := range ms.value {
		v.FlushDirty(dirty)
	}
}

func (ms *TestMsgEditorList) Add(t *TestMsgEditor) *TestMsgEditorList {
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

func (ms *TestMsgEditorList) AddAll(t []*TestMsgEditor) *TestMsgEditorList {
	for _, v := range t {
		ms.Add(v)
	}
	return ms
}

func (ms *TestMsgEditorList) Insert(pos int, v *TestMsgEditor) *TestMsgEditorList {
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
func (ms *TestMsgEditorList) Set(pos int, v *TestMsgEditor) *TestMsgEditorList {
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
func (ms *TestMsgEditorList) Iter() iter.Seq2[int, *TestMsgEditor] {
	return func(yield func(int, *TestMsgEditor) bool) {
		for i := range ms.value {
			if !yield(i, ms.value[i]) {
				return
			}
		}
	}
}

func (ms *TestMsgEditorList) Clear() {
	for _, v := range ms.value {
		v.SetParent(nil, nil)
	}
	ms.value = make([]*TestMsgEditor, 0)
	ms.monitor = make([]*TestMsgEditor, 0)
	ms.change = nil
	ms.change = append(ms.change, &protodef.ListOpLog{
		OpType: protodef.ListOpType_Clear,
		Pos:    0,
		RefIdx: 0,
	})
	ms.SetParentDirty()
}
func (ms *TestMsgEditorList) Len() int {
	return len(ms.value)
}

func (ms *TestMsgEditorList) Remove(v *TestMsgEditor) {
	ms.RemoveByIdx(slices.Index(ms.value, v))
}

func (ms *TestMsgEditorList) RemoveByIdx(pos int) {
	v := ms.value[pos]
	v.SetParent(nil, nil)
	ms.value = slices.Delete(ms.value, pos, pos+1)
	ms.change = append(ms.change, &protodef.ListOpLog{
		OpType: protodef.ListOpType_Delete,
		Pos:    int32(pos),
	})
	ms.SetParentDirty()
}

func (ms *TestMsgEditorList) ChangeLog() []*protodef.ListOpLog {
	return ms.change
}
func (ms *TestMsgEditorList) MonitorDirtyList() []*TestMsgEditor {
	return ms.monitor
}
func (ms *TestMsgEditorList) Values() []*TestMsgEditor {
	vv := make([]*TestMsgEditor, ms.Len())
	copy(vv, ms.value)
	return vv
}

func (ms *TestMsgEditorList) MergeDirtyFromPb(logs []*protodef.ListOpLog, dirties []*TestMsg) {
	for _, log := range logs {
		switch log.OpType {
		case protodef.ListOpType_Insert:
			v := dirties[log.RefIdx]
			tmp := NewTestMsgEditor().CopyFrom(v)
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

func (ms *TestMsgEditorList) Get(pos int32) *TestMsgEditor {
	return ms.value[pos]
}

type I32TestMsgEditorMap struct {
	parent      pbeditor.Editor
	keyInParent any
	value       map[int32]*TestMsgEditor
	dirtied     map[int32]struct{}
	deleted     map[int32]struct{}
}

func NewI32TestMsgEditorMap() *I32TestMsgEditorMap {
	return &I32TestMsgEditorMap{value: make(map[int32]*TestMsgEditor), dirtied: make(map[int32]struct{}), deleted: make(map[int32]struct{})}
}

func (ms *I32TestMsgEditorMap) SetParent(editor pbeditor.Editor, key any) {
	ms.parent = editor
	ms.keyInParent = key
}

func (ms *I32TestMsgEditorMap) SetParentDirty() {
	if ms.parent != nil {
		ms.parent.SetDirty(ms.keyInParent, true, ms)
	}
}

func (ms *I32TestMsgEditorMap) Entries() map[int32]*TestMsgEditor {
	return ms.value
}

func (ms *I32TestMsgEditorMap) SetDirty(key any, dirty bool, editor pbeditor.Editor) {
	if dirty {
		ms.dirtied[key.(int32)] = struct{}{}
		ms.SetParentDirty()
	} else {
		delete(ms.dirtied, key.(int32))
	}

}

func (ms *I32TestMsgEditorMap) FlushDirty(dirty bool) {
	if !dirty {
		ms.dirtied = map[int32]struct{}{}
		ms.deleted = map[int32]struct{}{}
	}
	for _, v := range ms.value {
		v.FlushDirty(dirty)
	}
}

func (ms *I32TestMsgEditorMap) Put(k int32, v *TestMsgEditor) *I32TestMsgEditorMap {
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

func (ms *I32TestMsgEditorMap) PutAll(kv map[int32]*TestMsgEditor) *I32TestMsgEditorMap {
	for k, v := range kv {
		ms.Put(k, v)
	}
	return ms
}

func (ms *I32TestMsgEditorMap) Len() int {
	return len(ms.value)
}

func (ms *I32TestMsgEditorMap) Clear() *I32TestMsgEditorMap {
	if ms.Len() <= 0 {
		return ms
	}
	for k, v := range ms.value {
		v.SetParent(nil, nil)
		ms.deleted[k] = struct{}{}
	}
	ms.value = map[int32]*TestMsgEditor{}
	ms.dirtied = map[int32]struct{}{}
	ms.SetParentDirty()
	return ms
}

func (ms *I32TestMsgEditorMap) Get(k int32) *TestMsgEditor {
	v := ms.value[k]
	return v
}

func (ms *I32TestMsgEditorMap) Remove(k int32) *TestMsgEditor {
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

func (ms *I32TestMsgEditorMap) RemoveAll(k []int32) {
	for _, kk := range k {
		ms.Remove(kk)
	}
}

func (ms *I32TestMsgEditorMap) Iter() iter.Seq2[int32, *TestMsgEditor] {
	return func(yield func(int32, *TestMsgEditor) bool) {
		for k, v := range ms.value {
			if !yield(k, v) {
				return
			}
		}
	}
}

func (ms *I32TestMsgEditorMap) Dirtied() map[int32]struct{} {
	return ms.dirtied
}

func (ms *I32TestMsgEditorMap) Deleted() map[int32]struct{} {
	return ms.deleted
}

func (ms *I32TestMsgEditorMap) ContainDirtied(kk int32) bool {
	for k := range ms.dirtied {
		if k == kk {
			return true
		}
	}
	return false
}

func (ms *I32TestMsgEditorMap) ContainDeleted(kk int32) bool {
	for k := range ms.deleted {
		if k == kk {
			return true
		}
	}
	return false
}

func (ms *I32TestMsgEditorMap) MergeDirtyFromPb(delList []int32, dirties map[int32]*TestMsg) {
	ms.RemoveAll(delList)
	for k, v := range dirties {
		var tmp = ms.Get(k)
		if tmp == nil {
			tmp = NewTestMsgEditor()
			tmp.MergeDirtyFromPb(v)
			ms.Put(k, tmp)
		} else {
			tmp.MergeDirtyFromPb(v)
		}
	}
}

type TestMsgEditor struct {
	field1       int32
	field2       string
	innerMsg     *PlainMsgEditor
	_dirties     map[int]struct{}
	_parent      pbeditor.Editor
	_keyInParent any
}

func NewTestMsgEditor() *TestMsgEditor {
	return &TestMsgEditor{
		_dirties: make(map[int]struct{}),
	}
}
func (x *TestMsgEditor) SetDirty(key any, dirty bool, sync pbeditor.Editor) {

	if dirty {
		x._dirties[key.(int)] = struct{}{}
		x.SetParentDirty()
	} else {
		delete(x._dirties, key.(int))
	}
}
func (x *TestMsgEditor) SetParentDirty() {
	if x._parent != nil {
		x._parent.SetDirty(x._keyInParent, true, x)
	}
}
func (x *TestMsgEditor) SetParent(editor pbeditor.Editor, key any) {
	x._parent = editor
	x._keyInParent = key
}
func (x *TestMsgEditor) FlushDirty(dirty bool) {
	if dirty || x.isField1Dirty() {
		x.setField1Dirty(dirty, true)
	}
	if dirty || x.isField2Dirty() {
		x.setField2Dirty(dirty, true)
	}
	if dirty || x.isInnerMsgDirty() {
		x.setInnerMsgDirty(dirty, true)
	}
}
func (x *TestMsgEditor) setField1Dirty(dirty bool, recur bool) {
	x.SetDirty(1, dirty, x)
}
func (x *TestMsgEditor) isField1Dirty() bool {
	_, ok := x._dirties[1]
	return ok
}
func (x *TestMsgEditor) setField2Dirty(dirty bool, recur bool) {
	x.SetDirty(2, dirty, x)
}
func (x *TestMsgEditor) isField2Dirty() bool {
	_, ok := x._dirties[2]
	return ok
}
func (x *TestMsgEditor) setInnerMsgDirty(dirty bool, recur bool) {
	x.SetDirty(3, dirty, x)
	if recur && x.innerMsg != nil {
		x.innerMsg.FlushDirty(dirty)
	}
}
func (x *TestMsgEditor) isInnerMsgDirty() bool {
	_, ok := x._dirties[3]
	return ok
}
func (x *TestMsgEditor) Key() any {
	return x._keyInParent
}
func (x *TestMsgEditor) SetKey(v any) {
	x._keyInParent = v
}
func (x *TestMsgEditor) Clear() *TestMsgEditor {
	x.SetField1(0)
	x.SetField2("")
	if x.innerMsg != nil {
		x.innerMsg.Clear()
	}
	return x
}
func (x *TestMsgEditor) CopyTo(r *TestMsg) *TestMsgEditor {
	r.SetField1(x.field1)
	r.SetField2(x.field2)
	if x.innerMsg != nil {
		tmp := &PlainMsg{}
		x.innerMsg.CopyTo(tmp)
		r.SetInnerMsg(tmp)
	}
	return x
}
func (x *TestMsgEditor) CopyFrom(r *TestMsg) *TestMsgEditor {
	if r.Field1 != nil {
		x.SetField1(*r.Field1)
	}
	if r.Field2 != nil {
		x.SetField2(*r.Field2)
	}
	if r.InnerMsg != nil {
		x.GetInnerMsg().CopyFrom(r.InnerMsg)
	}
	return x
}
func (x *TestMsgEditor) MergeDirtyFromPb(r *TestMsg) {
	if r.Field1 != nil {
		x.SetField1(*r.Field1)
	}
	if r.Field2 != nil {
		x.SetField2(*r.Field2)
	}
	if r.InnerMsg != nil {
		x.GetInnerMsg().MergeDirtyFromPb(r.InnerMsg)
	}
}
func (x *TestMsgEditor) MergeDirtyFromBytes(buf []byte) *TestMsgEditor {
	fds := pbeditor.ToRawMessage(buf)
	for _, rawFields := range fds.RawFields {
		if len(rawFields) <= 0 {
			continue
		}
		rawF := rawFields[0]
		switch rawF.Number {
		case 1:
			x.SetField1(int32(rawF.Variant))
		case 2:
			x.SetField2(pbeditor.Bys2Str(rawF.Bytes))
		case 3:
			x.GetInnerMsg().MergeDirtyFromBytes(rawF.Bytes)
		}
	}
	return x
}
func (x *TestMsgEditor) GatherDirtyToPb(r *TestMsg) {
	if x.isField1Dirty() {
		r.SetField1(x.field1)
	}
	if x.isField2Dirty() {
		r.SetField2(x.field2)
	}
	if x.isInnerMsgDirty() && x.innerMsg != nil {
		if r.InnerMsg == nil {
			r.InnerMsg = &PlainMsg{}
		}
		x.innerMsg.GatherDirtyToPb(r.InnerMsg)
	}
}
func (x *TestMsgEditor) GatherDirtyToBytes() []byte {
	var buf []byte
	if x.isField1Dirty() {
		buf = protowire.AppendTag(buf, 1, 0)
		buf = protowire.AppendVarint(buf, uint64(x.field1))
	}
	if x.isField2Dirty() {
		buf = protowire.AppendTag(buf, 2, 2)
		buf = protowire.AppendString(buf, x.field2)
	}
	if x.isInnerMsgDirty() && x.innerMsg != nil {
		buf = protowire.AppendTag(buf, 3, protowire.BytesType)
		buf = protowire.AppendBytes(buf, x.innerMsg.GatherDirtyToBytes())
	}
	return buf
}
func (x *TestMsgEditor) GetField1() int32 {
	return x.field1
}
func (x *TestMsgEditor) SetField1(v int32) *TestMsgEditor {
	if x.field1 == v {
		return x
	}
	x.field1 = v
	x.setField1Dirty(true, false)
	return x
}
func (x *TestMsgEditor) GetField2() string {
	return x.field2
}
func (x *TestMsgEditor) SetField2(v string) *TestMsgEditor {
	if x.field2 == v {
		return x
	}
	x.field2 = v
	x.setField2Dirty(true, false)
	return x
}
func (x *TestMsgEditor) GetInnerMsg() *PlainMsgEditor {
	if x.innerMsg == nil {
		x.innerMsg = NewPlainMsgEditor()
		x.innerMsg.SetParent(x, 3)
	}
	return x.innerMsg
}
func (x *TestMsgEditor) SetInnerMsg(v *PlainMsgEditor) *TestMsgEditor {
	if v != nil {
		v.SetParent(x, 3)
	}
	if x.innerMsg != nil {
		x.innerMsg.SetParent(nil, -1)
	}
	x.innerMsg = v
	x.setInnerMsgDirty(true, false)
	return x
}
func (xs *TestMsg) SetField1(v int32) {
	xs.Field1 = &v
}
func (xs *TestMsg) SetField2(v string) {
	xs.Field2 = &v
}
func (xs *TestMsg) SetInnerMsg(v *PlainMsg) {
	xs.InnerMsg = v
}

type StrPlainMsgEditorMap struct {
	parent      pbeditor.Editor
	keyInParent any
	value       map[string]*PlainMsgEditor
	dirtied     map[string]struct{}
	deleted     map[string]struct{}
}

func NewStrPlainMsgEditorMap() *StrPlainMsgEditorMap {
	return &StrPlainMsgEditorMap{value: make(map[string]*PlainMsgEditor), dirtied: make(map[string]struct{}), deleted: make(map[string]struct{})}
}

func (ms *StrPlainMsgEditorMap) SetParent(editor pbeditor.Editor, key any) {
	ms.parent = editor
	ms.keyInParent = key
}

func (ms *StrPlainMsgEditorMap) SetParentDirty() {
	if ms.parent != nil {
		ms.parent.SetDirty(ms.keyInParent, true, ms)
	}
}

func (ms *StrPlainMsgEditorMap) Entries() map[string]*PlainMsgEditor {
	return ms.value
}

func (ms *StrPlainMsgEditorMap) SetDirty(key any, dirty bool, editor pbeditor.Editor) {
	if dirty {
		ms.dirtied[key.(string)] = struct{}{}
		ms.SetParentDirty()
	} else {
		delete(ms.dirtied, key.(string))
	}

}

func (ms *StrPlainMsgEditorMap) FlushDirty(dirty bool) {
	if !dirty {
		ms.dirtied = map[string]struct{}{}
		ms.deleted = map[string]struct{}{}
	}
	for _, v := range ms.value {
		v.FlushDirty(dirty)
	}
}

func (ms *StrPlainMsgEditorMap) Put(k string, v *PlainMsgEditor) *StrPlainMsgEditorMap {
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

func (ms *StrPlainMsgEditorMap) PutAll(kv map[string]*PlainMsgEditor) *StrPlainMsgEditorMap {
	for k, v := range kv {
		ms.Put(k, v)
	}
	return ms
}

func (ms *StrPlainMsgEditorMap) Len() int {
	return len(ms.value)
}

func (ms *StrPlainMsgEditorMap) Clear() *StrPlainMsgEditorMap {
	if ms.Len() <= 0 {
		return ms
	}
	for k, v := range ms.value {
		v.SetParent(nil, nil)
		ms.deleted[k] = struct{}{}
	}
	ms.value = map[string]*PlainMsgEditor{}
	ms.dirtied = map[string]struct{}{}
	ms.SetParentDirty()
	return ms
}

func (ms *StrPlainMsgEditorMap) Get(k string) *PlainMsgEditor {
	v := ms.value[k]
	return v
}

func (ms *StrPlainMsgEditorMap) Remove(k string) *PlainMsgEditor {
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

func (ms *StrPlainMsgEditorMap) RemoveAll(k []string) {
	for _, kk := range k {
		ms.Remove(kk)
	}
}

func (ms *StrPlainMsgEditorMap) Iter() iter.Seq2[string, *PlainMsgEditor] {
	return func(yield func(string, *PlainMsgEditor) bool) {
		for k, v := range ms.value {
			if !yield(k, v) {
				return
			}
		}
	}
}

func (ms *StrPlainMsgEditorMap) Dirtied() map[string]struct{} {
	return ms.dirtied
}

func (ms *StrPlainMsgEditorMap) Deleted() map[string]struct{} {
	return ms.deleted
}

func (ms *StrPlainMsgEditorMap) ContainDirtied(kk string) bool {
	for k := range ms.dirtied {
		if k == kk {
			return true
		}
	}
	return false
}

func (ms *StrPlainMsgEditorMap) ContainDeleted(kk string) bool {
	for k := range ms.deleted {
		if k == kk {
			return true
		}
	}
	return false
}

func (ms *StrPlainMsgEditorMap) MergeDirtyFromPb(delList []string, dirties map[string]*PlainMsg) {
	ms.RemoveAll(delList)
	for k, v := range dirties {
		var tmp = ms.Get(k)
		if tmp == nil {
			tmp = NewPlainMsgEditor()
			tmp.MergeDirtyFromPb(v)
			ms.Put(k, tmp)
		} else {
			tmp.MergeDirtyFromPb(v)
		}
	}
}

type PlainMsgEditor struct {
	i32          int32
	str          string
	b            bool
	_dirties     map[int]struct{}
	_parent      pbeditor.Editor
	_keyInParent any
}

func NewPlainMsgEditor() *PlainMsgEditor {
	return &PlainMsgEditor{
		_dirties: make(map[int]struct{}),
	}
}
func (x *PlainMsgEditor) SetDirty(key any, dirty bool, sync pbeditor.Editor) {

	if dirty {
		x._dirties[key.(int)] = struct{}{}
		x.SetParentDirty()
	} else {
		delete(x._dirties, key.(int))
	}
}
func (x *PlainMsgEditor) SetParentDirty() {
	if x._parent != nil {
		x._parent.SetDirty(x._keyInParent, true, x)
	}
}
func (x *PlainMsgEditor) SetParent(editor pbeditor.Editor, key any) {
	x._parent = editor
	x._keyInParent = key
}
func (x *PlainMsgEditor) FlushDirty(dirty bool) {
	if dirty || x.isI32Dirty() {
		x.setI32Dirty(dirty, true)
	}
	if dirty || x.isStrDirty() {
		x.setStrDirty(dirty, true)
	}
	if dirty || x.isBDirty() {
		x.setBDirty(dirty, true)
	}
}
func (x *PlainMsgEditor) setI32Dirty(dirty bool, recur bool) {
	x.SetDirty(1, dirty, x)
}
func (x *PlainMsgEditor) isI32Dirty() bool {
	_, ok := x._dirties[1]
	return ok
}
func (x *PlainMsgEditor) setStrDirty(dirty bool, recur bool) {
	x.SetDirty(2, dirty, x)
}
func (x *PlainMsgEditor) isStrDirty() bool {
	_, ok := x._dirties[2]
	return ok
}
func (x *PlainMsgEditor) setBDirty(dirty bool, recur bool) {
	x.SetDirty(3, dirty, x)
}
func (x *PlainMsgEditor) isBDirty() bool {
	_, ok := x._dirties[3]
	return ok
}
func (x *PlainMsgEditor) Key() any {
	return x._keyInParent
}
func (x *PlainMsgEditor) SetKey(v any) {
	x._keyInParent = v
}
func (x *PlainMsgEditor) Clear() *PlainMsgEditor {
	x.SetI32(0)
	x.SetStr("")
	x.SetB(false)
	return x
}
func (x *PlainMsgEditor) CopyTo(r *PlainMsg) *PlainMsgEditor {
	r.SetI32(x.i32)
	r.SetStr(x.str)
	r.SetB(x.b)
	return x
}
func (x *PlainMsgEditor) CopyFrom(r *PlainMsg) *PlainMsgEditor {
	if r.I32 != nil {
		x.SetI32(*r.I32)
	}
	if r.Str != nil {
		x.SetStr(*r.Str)
	}
	if r.B != nil {
		x.SetB(*r.B)
	}
	return x
}
func (x *PlainMsgEditor) MergeDirtyFromPb(r *PlainMsg) {
	if r.I32 != nil {
		x.SetI32(*r.I32)
	}
	if r.Str != nil {
		x.SetStr(*r.Str)
	}
	if r.B != nil {
		x.SetB(*r.B)
	}
}
func (x *PlainMsgEditor) MergeDirtyFromBytes(buf []byte) *PlainMsgEditor {
	fds := pbeditor.ToRawMessage(buf)
	for _, rawFields := range fds.RawFields {
		if len(rawFields) <= 0 {
			continue
		}
		rawF := rawFields[0]
		switch rawF.Number {
		case 1:
			x.SetI32(int32(rawF.Variant))
		case 2:
			x.SetStr(pbeditor.Bys2Str(rawF.Bytes))
		case 3:
			x.SetB(protowire.DecodeBool(rawF.Variant))
		}
	}
	return x
}
func (x *PlainMsgEditor) GatherDirtyToPb(r *PlainMsg) {
	if x.isI32Dirty() {
		r.SetI32(x.i32)
	}
	if x.isStrDirty() {
		r.SetStr(x.str)
	}
	if x.isBDirty() {
		r.SetB(x.b)
	}
}
func (x *PlainMsgEditor) GatherDirtyToBytes() []byte {
	var buf []byte
	if x.isI32Dirty() {
		buf = protowire.AppendTag(buf, 1, 0)
		buf = protowire.AppendVarint(buf, uint64(x.i32))
	}
	if x.isStrDirty() {
		buf = protowire.AppendTag(buf, 2, 2)
		buf = protowire.AppendString(buf, x.str)
	}
	if x.isBDirty() {
		buf = protowire.AppendTag(buf, 3, 0)
		var t uint64 = 0
		if x.b {
			t = 1
		}
		buf = protowire.AppendVarint(buf, t)
	}
	return buf
}
func (x *PlainMsgEditor) GetI32() int32 {
	return x.i32
}
func (x *PlainMsgEditor) SetI32(v int32) *PlainMsgEditor {
	if x.i32 == v {
		return x
	}
	x.i32 = v
	x.setI32Dirty(true, false)
	return x
}
func (x *PlainMsgEditor) GetStr() string {
	return x.str
}
func (x *PlainMsgEditor) SetStr(v string) *PlainMsgEditor {
	if x.str == v {
		return x
	}
	x.str = v
	x.setStrDirty(true, false)
	return x
}
func (x *PlainMsgEditor) GetB() bool {
	return x.b
}
func (x *PlainMsgEditor) SetB(v bool) *PlainMsgEditor {
	if x.b == v {
		return x
	}
	x.b = v
	x.setBDirty(true, false)
	return x
}
func (xs *PlainMsg) SetI32(v int32) {
	xs.I32 = &v
}
func (xs *PlainMsg) SetStr(v string) {
	xs.Str = &v
}
func (xs *PlainMsg) SetB(v bool) {
	xs.B = &v
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
