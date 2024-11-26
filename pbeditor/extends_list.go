package pbeditor

import "iter"
import "slices"
import "github.com/yaoguangduan/proto-editor/protodef"

type I64List struct {
	parent   Editor
	idxInPar any
	value    []int64
	monitor  []int64
	change   []*protodef.ListOpLog
}

func NewI64List() *I64List {
	return &I64List{value: make([]int64, 0)}
}

func (ms *I64List) SetParent(editor Editor, key any) {
	ms.parent = editor
	ms.idxInPar = key
}
func (ms *I64List) SetParentDirty() {
	if ms.parent != nil {
		ms.parent.SetDirty(ms.idxInPar, true, ms)
	}
}
func (ms *I64List) SetDirty(key any, dirty bool, sync Editor) {
	pos := key.(int)
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
func (ms *I64List) FlushDirty(dirty bool) {
	if !dirty {
		ms.monitor = nil
		ms.change = nil
	}
}

func (ms *I64List) Add(t int64) *I64List {
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

func (ms *I64List) AddAll(t []int64) *I64List {
	for _, v := range t {
		ms.Add(v)
	}
	return ms
}

func (ms *I64List) Insert(pos int, v int64) *I64List {
	if pos < 0 {
		pos = 0
	}
	if pos >= len(ms.value) {
		ms.Add(v)
	}
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
func (ms *I64List) Set(pos int, v int64) *I64List {
	ms.value[pos] = v
	ms.SetDirty(pos, true, ms)
	return ms
}
func (ms *I64List) Iter() iter.Seq2[int, int64] {
	return func(yield func(int, int64) bool) {
		for i := range ms.value {
			if !yield(i, ms.value[i]) {
				return
			}
		}
	}
}

func (ms *I64List) Clear() {
	ms.value = make([]int64, 0)
	ms.monitor = make([]int64, 0)
	ms.change = nil
	ms.change = append(ms.change, &protodef.ListOpLog{
		OpType: protodef.ListOpType_Clear,
		Pos:    0,
		RefIdx: 0,
	})
	ms.SetParentDirty()
}
func (ms *I64List) Len() int {
	return len(ms.value)
}

func (ms *I64List) Remove(v int64) {
	ms.RemoveByIdx(slices.Index(ms.value, v))
}

func (ms *I64List) RemoveByIdx(pos int) {
	ms.value = slices.Delete(ms.value, pos, pos+1)
	ms.change = append(ms.change, &protodef.ListOpLog{
		OpType: protodef.ListOpType_Delete,
		Pos:    int32(pos),
	})
	ms.SetParentDirty()
}
func (ms *I64List) ChangeLog() []*protodef.ListOpLog {
	return ms.change
}
func (ms *I64List) MonitorDirtyList() []int64 {
	return ms.monitor
}
func (ms *I64List) Values() []int64 {
	vv := make([]int64, ms.Len())
	copy(vv, ms.value)
	return vv
}

func (ms *I64List) MergeDirtyFromPb(logs []*protodef.ListOpLog, dirties []int64) {
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

func (ms *I64List) Get(pos int32) int64 {
	return ms.value[pos]
}

type StrList struct {
	parent   Editor
	idxInPar any
	value    []string
	monitor  []string
	change   []*protodef.ListOpLog
}

func NewStrList() *StrList {
	return &StrList{value: make([]string, 0)}
}

func (ms *StrList) SetParent(editor Editor, key any) {
	ms.parent = editor
	ms.idxInPar = key
}
func (ms *StrList) SetParentDirty() {
	if ms.parent != nil {
		ms.parent.SetDirty(ms.idxInPar, true, ms)
	}
}
func (ms *StrList) SetDirty(key any, dirty bool, sync Editor) {
	pos := key.(int)
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
func (ms *StrList) FlushDirty(dirty bool) {
	if !dirty {
		ms.monitor = nil
		ms.change = nil
	}
}

func (ms *StrList) Add(t string) *StrList {
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

func (ms *StrList) AddAll(t []string) *StrList {
	for _, v := range t {
		ms.Add(v)
	}
	return ms
}

func (ms *StrList) Insert(pos int, v string) *StrList {
	if pos < 0 {
		pos = 0
	}
	if pos >= len(ms.value) {
		ms.Add(v)
	}
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
func (ms *StrList) Set(pos int, v string) *StrList {
	ms.value[pos] = v
	ms.SetDirty(pos, true, ms)
	return ms
}
func (ms *StrList) Iter() iter.Seq2[int, string] {
	return func(yield func(int, string) bool) {
		for i := range ms.value {
			if !yield(i, ms.value[i]) {
				return
			}
		}
	}
}

func (ms *StrList) Clear() {
	ms.value = make([]string, 0)
	ms.monitor = make([]string, 0)
	ms.change = nil
	ms.change = append(ms.change, &protodef.ListOpLog{
		OpType: protodef.ListOpType_Clear,
		Pos:    0,
		RefIdx: 0,
	})
	ms.SetParentDirty()
}
func (ms *StrList) Len() int {
	return len(ms.value)
}

func (ms *StrList) Remove(v string) {
	ms.RemoveByIdx(slices.Index(ms.value, v))
}

func (ms *StrList) RemoveByIdx(pos int) {
	ms.value = slices.Delete(ms.value, pos, pos+1)
	ms.change = append(ms.change, &protodef.ListOpLog{
		OpType: protodef.ListOpType_Delete,
		Pos:    int32(pos),
	})
	ms.SetParentDirty()
}
func (ms *StrList) ChangeLog() []*protodef.ListOpLog {
	return ms.change
}
func (ms *StrList) MonitorDirtyList() []string {
	return ms.monitor
}
func (ms *StrList) Values() []string {
	vv := make([]string, ms.Len())
	copy(vv, ms.value)
	return vv
}

func (ms *StrList) MergeDirtyFromPb(logs []*protodef.ListOpLog, dirties []string) {
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

func (ms *StrList) Get(pos int32) string {
	return ms.value[pos]
}

type BoolList struct {
	parent   Editor
	idxInPar any
	value    []bool
	monitor  []bool
	change   []*protodef.ListOpLog
}

func NewBoolList() *BoolList {
	return &BoolList{value: make([]bool, 0)}
}

func (ms *BoolList) SetParent(editor Editor, key any) {
	ms.parent = editor
	ms.idxInPar = key
}
func (ms *BoolList) SetParentDirty() {
	if ms.parent != nil {
		ms.parent.SetDirty(ms.idxInPar, true, ms)
	}
}
func (ms *BoolList) SetDirty(key any, dirty bool, sync Editor) {
	pos := key.(int)
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
func (ms *BoolList) FlushDirty(dirty bool) {
	if !dirty {
		ms.monitor = nil
		ms.change = nil
	}
}

func (ms *BoolList) Add(t bool) *BoolList {
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

func (ms *BoolList) AddAll(t []bool) *BoolList {
	for _, v := range t {
		ms.Add(v)
	}
	return ms
}

func (ms *BoolList) Insert(pos int, v bool) *BoolList {
	if pos < 0 {
		pos = 0
	}
	if pos >= len(ms.value) {
		ms.Add(v)
	}
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
func (ms *BoolList) Set(pos int, v bool) *BoolList {
	ms.value[pos] = v
	ms.SetDirty(pos, true, ms)
	return ms
}
func (ms *BoolList) Iter() iter.Seq2[int, bool] {
	return func(yield func(int, bool) bool) {
		for i := range ms.value {
			if !yield(i, ms.value[i]) {
				return
			}
		}
	}
}

func (ms *BoolList) Clear() {
	ms.value = make([]bool, 0)
	ms.monitor = make([]bool, 0)
	ms.change = nil
	ms.change = append(ms.change, &protodef.ListOpLog{
		OpType: protodef.ListOpType_Clear,
		Pos:    0,
		RefIdx: 0,
	})
	ms.SetParentDirty()
}
func (ms *BoolList) Len() int {
	return len(ms.value)
}

func (ms *BoolList) Remove(v bool) {
	ms.RemoveByIdx(slices.Index(ms.value, v))
}

func (ms *BoolList) RemoveByIdx(pos int) {
	ms.value = slices.Delete(ms.value, pos, pos+1)
	ms.change = append(ms.change, &protodef.ListOpLog{
		OpType: protodef.ListOpType_Delete,
		Pos:    int32(pos),
	})
	ms.SetParentDirty()
}
func (ms *BoolList) ChangeLog() []*protodef.ListOpLog {
	return ms.change
}
func (ms *BoolList) MonitorDirtyList() []bool {
	return ms.monitor
}
func (ms *BoolList) Values() []bool {
	vv := make([]bool, ms.Len())
	copy(vv, ms.value)
	return vv
}

func (ms *BoolList) MergeDirtyFromPb(logs []*protodef.ListOpLog, dirties []bool) {
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

func (ms *BoolList) Get(pos int32) bool {
	return ms.value[pos]
}

type F32List struct {
	parent   Editor
	idxInPar any
	value    []float32
	monitor  []float32
	change   []*protodef.ListOpLog
}

func NewF32List() *F32List {
	return &F32List{value: make([]float32, 0)}
}

func (ms *F32List) SetParent(editor Editor, key any) {
	ms.parent = editor
	ms.idxInPar = key
}
func (ms *F32List) SetParentDirty() {
	if ms.parent != nil {
		ms.parent.SetDirty(ms.idxInPar, true, ms)
	}
}
func (ms *F32List) SetDirty(key any, dirty bool, sync Editor) {
	pos := key.(int)
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
func (ms *F32List) FlushDirty(dirty bool) {
	if !dirty {
		ms.monitor = nil
		ms.change = nil
	}
}

func (ms *F32List) Add(t float32) *F32List {
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

func (ms *F32List) AddAll(t []float32) *F32List {
	for _, v := range t {
		ms.Add(v)
	}
	return ms
}

func (ms *F32List) Insert(pos int, v float32) *F32List {
	if pos < 0 {
		pos = 0
	}
	if pos >= len(ms.value) {
		ms.Add(v)
	}
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
func (ms *F32List) Set(pos int, v float32) *F32List {
	ms.value[pos] = v
	ms.SetDirty(pos, true, ms)
	return ms
}
func (ms *F32List) Iter() iter.Seq2[int, float32] {
	return func(yield func(int, float32) bool) {
		for i := range ms.value {
			if !yield(i, ms.value[i]) {
				return
			}
		}
	}
}

func (ms *F32List) Clear() {
	ms.value = make([]float32, 0)
	ms.monitor = make([]float32, 0)
	ms.change = nil
	ms.change = append(ms.change, &protodef.ListOpLog{
		OpType: protodef.ListOpType_Clear,
		Pos:    0,
		RefIdx: 0,
	})
	ms.SetParentDirty()
}
func (ms *F32List) Len() int {
	return len(ms.value)
}

func (ms *F32List) Remove(v float32) {
	ms.RemoveByIdx(slices.Index(ms.value, v))
}

func (ms *F32List) RemoveByIdx(pos int) {
	ms.value = slices.Delete(ms.value, pos, pos+1)
	ms.change = append(ms.change, &protodef.ListOpLog{
		OpType: protodef.ListOpType_Delete,
		Pos:    int32(pos),
	})
	ms.SetParentDirty()
}
func (ms *F32List) ChangeLog() []*protodef.ListOpLog {
	return ms.change
}
func (ms *F32List) MonitorDirtyList() []float32 {
	return ms.monitor
}
func (ms *F32List) Values() []float32 {
	vv := make([]float32, ms.Len())
	copy(vv, ms.value)
	return vv
}

func (ms *F32List) MergeDirtyFromPb(logs []*protodef.ListOpLog, dirties []float32) {
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

func (ms *F32List) Get(pos int32) float32 {
	return ms.value[pos]
}

type F64List struct {
	parent   Editor
	idxInPar any
	value    []float64
	monitor  []float64
	change   []*protodef.ListOpLog
}

func NewF64List() *F64List {
	return &F64List{value: make([]float64, 0)}
}

func (ms *F64List) SetParent(editor Editor, key any) {
	ms.parent = editor
	ms.idxInPar = key
}
func (ms *F64List) SetParentDirty() {
	if ms.parent != nil {
		ms.parent.SetDirty(ms.idxInPar, true, ms)
	}
}
func (ms *F64List) SetDirty(key any, dirty bool, sync Editor) {
	pos := key.(int)
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
func (ms *F64List) FlushDirty(dirty bool) {
	if !dirty {
		ms.monitor = nil
		ms.change = nil
	}
}

func (ms *F64List) Add(t float64) *F64List {
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

func (ms *F64List) AddAll(t []float64) *F64List {
	for _, v := range t {
		ms.Add(v)
	}
	return ms
}

func (ms *F64List) Insert(pos int, v float64) *F64List {
	if pos < 0 {
		pos = 0
	}
	if pos >= len(ms.value) {
		ms.Add(v)
	}
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
func (ms *F64List) Set(pos int, v float64) *F64List {
	ms.value[pos] = v
	ms.SetDirty(pos, true, ms)
	return ms
}
func (ms *F64List) Iter() iter.Seq2[int, float64] {
	return func(yield func(int, float64) bool) {
		for i := range ms.value {
			if !yield(i, ms.value[i]) {
				return
			}
		}
	}
}

func (ms *F64List) Clear() {
	ms.value = make([]float64, 0)
	ms.monitor = make([]float64, 0)
	ms.change = nil
	ms.change = append(ms.change, &protodef.ListOpLog{
		OpType: protodef.ListOpType_Clear,
		Pos:    0,
		RefIdx: 0,
	})
	ms.SetParentDirty()
}
func (ms *F64List) Len() int {
	return len(ms.value)
}

func (ms *F64List) Remove(v float64) {
	ms.RemoveByIdx(slices.Index(ms.value, v))
}

func (ms *F64List) RemoveByIdx(pos int) {
	ms.value = slices.Delete(ms.value, pos, pos+1)
	ms.change = append(ms.change, &protodef.ListOpLog{
		OpType: protodef.ListOpType_Delete,
		Pos:    int32(pos),
	})
	ms.SetParentDirty()
}
func (ms *F64List) ChangeLog() []*protodef.ListOpLog {
	return ms.change
}
func (ms *F64List) MonitorDirtyList() []float64 {
	return ms.monitor
}
func (ms *F64List) Values() []float64 {
	vv := make([]float64, ms.Len())
	copy(vv, ms.value)
	return vv
}

func (ms *F64List) MergeDirtyFromPb(logs []*protodef.ListOpLog, dirties []float64) {
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

func (ms *F64List) Get(pos int32) float64 {
	return ms.value[pos]
}

type U32List struct {
	parent   Editor
	idxInPar any
	value    []uint32
	monitor  []uint32
	change   []*protodef.ListOpLog
}

func NewU32List() *U32List {
	return &U32List{value: make([]uint32, 0)}
}

func (ms *U32List) SetParent(editor Editor, key any) {
	ms.parent = editor
	ms.idxInPar = key
}
func (ms *U32List) SetParentDirty() {
	if ms.parent != nil {
		ms.parent.SetDirty(ms.idxInPar, true, ms)
	}
}
func (ms *U32List) SetDirty(key any, dirty bool, sync Editor) {
	pos := key.(int)
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
func (ms *U32List) FlushDirty(dirty bool) {
	if !dirty {
		ms.monitor = nil
		ms.change = nil
	}
}

func (ms *U32List) Add(t uint32) *U32List {
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

func (ms *U32List) AddAll(t []uint32) *U32List {
	for _, v := range t {
		ms.Add(v)
	}
	return ms
}

func (ms *U32List) Insert(pos int, v uint32) *U32List {
	if pos < 0 {
		pos = 0
	}
	if pos >= len(ms.value) {
		ms.Add(v)
	}
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
func (ms *U32List) Set(pos int, v uint32) *U32List {
	ms.value[pos] = v
	ms.SetDirty(pos, true, ms)
	return ms
}
func (ms *U32List) Iter() iter.Seq2[int, uint32] {
	return func(yield func(int, uint32) bool) {
		for i := range ms.value {
			if !yield(i, ms.value[i]) {
				return
			}
		}
	}
}

func (ms *U32List) Clear() {
	ms.value = make([]uint32, 0)
	ms.monitor = make([]uint32, 0)
	ms.change = nil
	ms.change = append(ms.change, &protodef.ListOpLog{
		OpType: protodef.ListOpType_Clear,
		Pos:    0,
		RefIdx: 0,
	})
	ms.SetParentDirty()
}
func (ms *U32List) Len() int {
	return len(ms.value)
}

func (ms *U32List) Remove(v uint32) {
	ms.RemoveByIdx(slices.Index(ms.value, v))
}

func (ms *U32List) RemoveByIdx(pos int) {
	ms.value = slices.Delete(ms.value, pos, pos+1)
	ms.change = append(ms.change, &protodef.ListOpLog{
		OpType: protodef.ListOpType_Delete,
		Pos:    int32(pos),
	})
	ms.SetParentDirty()
}
func (ms *U32List) ChangeLog() []*protodef.ListOpLog {
	return ms.change
}
func (ms *U32List) MonitorDirtyList() []uint32 {
	return ms.monitor
}
func (ms *U32List) Values() []uint32 {
	vv := make([]uint32, ms.Len())
	copy(vv, ms.value)
	return vv
}

func (ms *U32List) MergeDirtyFromPb(logs []*protodef.ListOpLog, dirties []uint32) {
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

func (ms *U32List) Get(pos int32) uint32 {
	return ms.value[pos]
}

type U64List struct {
	parent   Editor
	idxInPar any
	value    []uint64
	monitor  []uint64
	change   []*protodef.ListOpLog
}

func NewU64List() *U64List {
	return &U64List{value: make([]uint64, 0)}
}

func (ms *U64List) SetParent(editor Editor, key any) {
	ms.parent = editor
	ms.idxInPar = key
}
func (ms *U64List) SetParentDirty() {
	if ms.parent != nil {
		ms.parent.SetDirty(ms.idxInPar, true, ms)
	}
}
func (ms *U64List) SetDirty(key any, dirty bool, sync Editor) {
	pos := key.(int)
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
func (ms *U64List) FlushDirty(dirty bool) {
	if !dirty {
		ms.monitor = nil
		ms.change = nil
	}
}

func (ms *U64List) Add(t uint64) *U64List {
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

func (ms *U64List) AddAll(t []uint64) *U64List {
	for _, v := range t {
		ms.Add(v)
	}
	return ms
}

func (ms *U64List) Insert(pos int, v uint64) *U64List {
	if pos < 0 {
		pos = 0
	}
	if pos >= len(ms.value) {
		ms.Add(v)
	}
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
func (ms *U64List) Set(pos int, v uint64) *U64List {
	ms.value[pos] = v
	ms.SetDirty(pos, true, ms)
	return ms
}
func (ms *U64List) Iter() iter.Seq2[int, uint64] {
	return func(yield func(int, uint64) bool) {
		for i := range ms.value {
			if !yield(i, ms.value[i]) {
				return
			}
		}
	}
}

func (ms *U64List) Clear() {
	ms.value = make([]uint64, 0)
	ms.monitor = make([]uint64, 0)
	ms.change = nil
	ms.change = append(ms.change, &protodef.ListOpLog{
		OpType: protodef.ListOpType_Clear,
		Pos:    0,
		RefIdx: 0,
	})
	ms.SetParentDirty()
}
func (ms *U64List) Len() int {
	return len(ms.value)
}

func (ms *U64List) Remove(v uint64) {
	ms.RemoveByIdx(slices.Index(ms.value, v))
}

func (ms *U64List) RemoveByIdx(pos int) {
	ms.value = slices.Delete(ms.value, pos, pos+1)
	ms.change = append(ms.change, &protodef.ListOpLog{
		OpType: protodef.ListOpType_Delete,
		Pos:    int32(pos),
	})
	ms.SetParentDirty()
}
func (ms *U64List) ChangeLog() []*protodef.ListOpLog {
	return ms.change
}
func (ms *U64List) MonitorDirtyList() []uint64 {
	return ms.monitor
}
func (ms *U64List) Values() []uint64 {
	vv := make([]uint64, ms.Len())
	copy(vv, ms.value)
	return vv
}

func (ms *U64List) MergeDirtyFromPb(logs []*protodef.ListOpLog, dirties []uint64) {
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

func (ms *U64List) Get(pos int32) uint64 {
	return ms.value[pos]
}

type I32List struct {
	parent   Editor
	idxInPar any
	value    []int32
	monitor  []int32
	change   []*protodef.ListOpLog
}

func NewI32List() *I32List {
	return &I32List{value: make([]int32, 0)}
}

func (ms *I32List) SetParent(editor Editor, key any) {
	ms.parent = editor
	ms.idxInPar = key
}
func (ms *I32List) SetParentDirty() {
	if ms.parent != nil {
		ms.parent.SetDirty(ms.idxInPar, true, ms)
	}
}
func (ms *I32List) SetDirty(key any, dirty bool, sync Editor) {
	pos := key.(int)
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
func (ms *I32List) FlushDirty(dirty bool) {
	if !dirty {
		ms.monitor = nil
		ms.change = nil
	}
}

func (ms *I32List) Add(t int32) *I32List {
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

func (ms *I32List) AddAll(t []int32) *I32List {
	for _, v := range t {
		ms.Add(v)
	}
	return ms
}

func (ms *I32List) Insert(pos int, v int32) *I32List {
	if pos < 0 {
		pos = 0
	}
	if pos >= len(ms.value) {
		ms.Add(v)
	}
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
func (ms *I32List) Set(pos int, v int32) *I32List {
	ms.value[pos] = v
	ms.SetDirty(pos, true, ms)
	return ms
}
func (ms *I32List) Iter() iter.Seq2[int, int32] {
	return func(yield func(int, int32) bool) {
		for i := range ms.value {
			if !yield(i, ms.value[i]) {
				return
			}
		}
	}
}

func (ms *I32List) Clear() {
	ms.value = make([]int32, 0)
	ms.monitor = make([]int32, 0)
	ms.change = nil
	ms.change = append(ms.change, &protodef.ListOpLog{
		OpType: protodef.ListOpType_Clear,
		Pos:    0,
		RefIdx: 0,
	})
	ms.SetParentDirty()
}
func (ms *I32List) Len() int {
	return len(ms.value)
}

func (ms *I32List) Remove(v int32) {
	ms.RemoveByIdx(slices.Index(ms.value, v))
}

func (ms *I32List) RemoveByIdx(pos int) {
	ms.value = slices.Delete(ms.value, pos, pos+1)
	ms.change = append(ms.change, &protodef.ListOpLog{
		OpType: protodef.ListOpType_Delete,
		Pos:    int32(pos),
	})
	ms.SetParentDirty()
}
func (ms *I32List) ChangeLog() []*protodef.ListOpLog {
	return ms.change
}
func (ms *I32List) MonitorDirtyList() []int32 {
	return ms.monitor
}
func (ms *I32List) Values() []int32 {
	vv := make([]int32, ms.Len())
	copy(vv, ms.value)
	return vv
}

func (ms *I32List) MergeDirtyFromPb(logs []*protodef.ListOpLog, dirties []int32) {
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

func (ms *I32List) Get(pos int32) int32 {
	return ms.value[pos]
}
