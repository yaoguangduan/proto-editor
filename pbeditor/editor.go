package pbeditor

type Editor interface {
	SetParent(editor Editor, key any)
	SetDirty(key any, dirty bool, editor Editor)
	SetParentDirty()
	FlushDirty(dirty bool)
}

func Keys[K comparable, V any](m map[K]V) []K {
	keys := make([]K, 0, len(m))
	for key := range m {
		keys = append(keys, key)
	}
	return keys
}

type ArrayType interface {
	~int32 | int64 | uint32 | uint64 | string | bool | float32 | float64
}

type MsgKey interface {
	int32 | int64 | uint32 | uint64 | string | bool
}

type MsgValue interface {
	~int32 | int64 | uint32 | uint64 | string | bool | float32 | float64
}
type MapEditor[K MsgKey, V Editor] struct {
	parent      Editor
	keyInParent any
	value       map[K]V
	dirtied     map[K]struct{}
	deleted     map[K]struct{}
}

func NewMapEditor[K MsgKey, V Editor]() *MapEditor[K, V] {
	return &MapEditor[K, V]{value: make(map[K]V), dirtied: make(map[K]struct{}), deleted: make(map[K]struct{})}
}

func (ms *MapEditor[K, V]) SetParent(sync Editor, key any) {
	ms.parent = sync
	ms.keyInParent = key
}
func (ms *MapEditor[K, V]) SetParentDirty() {
	if ms.parent != nil {
		ms.parent.SetDirty(ms.keyInParent, true, ms)
	}
}
func (ms *MapEditor[K, V]) SetDirty(key any, dirty bool, editor Editor) {
	if dirty {
		//ms.dirtied[key] = struct{}{}
		ms.SetParentDirty()
	} else {
		//delete(ms.dirtied, key)
	}

}
func (ms *MapEditor[K, V]) FlushDirty(dirty bool) {
	if !dirty {
		ms.dirtied = map[K]struct{}{}
		ms.deleted = map[K]struct{}{}
	}
	for k := range ms.value {
		any(ms.value[k]).(Editor).FlushDirty(dirty)
	}
}

func (ms *MapEditor[K, V]) Key() interface{} {
	return nil
}
func (ms *MapEditor[K, V]) SetKey(i interface{}) {

}

func (ms *MapEditor[K, V]) Put(k K, s V) *MapEditor[K, V] {
	old, exist := ms.value[k]
	if exist {
		old.SetParent(nil, -1)
		ms.deleted[k] = struct{}{}
	}
	s.SetParent(ms, -1)
	ms.value[k] = s
	ms.dirtied[k] = struct{}{}
	ms.SetParentDirty()
	return ms
}

func (ms *MapEditor[K, V]) PutAll(kv map[K]V) *MapEditor[K, V] {
	for k, v := range kv {
		ms.Put(k, v)
	}
	return ms
}

func (ms *MapEditor[K, V]) Len() int {
	return len(ms.value)
}

func (ms *MapEditor[K, V]) Clear() *MapEditor[K, V] {
	if ms.Len() <= 0 {
		return ms
	}
	for k, v := range ms.value {
		v.SetParent(nil, -1)
		ms.deleted[k] = struct{}{}
	}
	ms.value = map[K]V{}
	ms.dirtied = map[K]struct{}{}
	ms.SetParentDirty()
	return ms
}
func (ms *MapEditor[K, V]) Get(k K) V {
	v := ms.value[k]
	return v
}

func (ms *MapEditor[K, V]) Remove(k K) V {
	v, exist := ms.value[k]
	if !exist {
		return v
	}
	delete(ms.value, k)
	delete(ms.dirtied, k)
	if vPtr, ok := any(v).(interface {
		SetParent(sync Editor, idx int)
	}); ok && vPtr != nil {
		vPtr.SetParent(nil, -1)
	}
	ms.deleted[k] = struct{}{}
	ms.SetParentDirty()
	return v
}
func (ms *MapEditor[K, V]) RemoveAll(k []K) {
	for _, kk := range k {
		ms.Remove(kk)
	}
}
func (ms *MapEditor[K, V]) Each(f func(k K, v V) bool) {
	for k, v := range ms.value {
		if !f(k, v) {
			break
		}
	}
}

func (ms *MapEditor[K, V]) Dirtied() map[K]struct{} {
	return ms.dirtied
}

func (ms *MapEditor[K, V]) Deleted() map[K]struct{} {
	return ms.deleted
}

func (ms *MapEditor[K, V]) ContainDirtied(kk K) bool {
	for k := range ms.dirtied {
		if k == kk {
			return true
		}
	}
	return false
}

func (ms *MapEditor[K, V]) ContainDeleted(kk K) bool {
	for k := range ms.deleted {
		if k == kk {
			return true
		}
	}
	return false
}

type ArrayEditor[V ArrayType] struct {
	parent   Editor
	idxInPar any
	value    []V
}

func NewArrayEditor[V ArrayType]() *ArrayEditor[V] {
	return &ArrayEditor[V]{value: make([]V, 0)}
}

func (ms *ArrayEditor[V]) SetParent(sync Editor, key any) {
	ms.parent = sync
	ms.idxInPar = key
}
func (ms *ArrayEditor[V]) SetParentDirty() {
	if ms.parent != nil {
		ms.parent.SetDirty(ms.idxInPar, true, ms)
	}
}
func (ms *ArrayEditor[V]) SetDirty(key any, dirty bool, sync Editor) {
	ms.SetParentDirty()
}
func (ms *ArrayEditor[V]) FlushDirty(dirty bool) {
}

func (ms *ArrayEditor[V]) Key() interface{} {
	return nil
}
func (ms *ArrayEditor[V]) SetKey(i interface{}) {

}
func (ms *ArrayEditor[V]) Add(t V) *ArrayEditor[V] {
	ms.value = append(ms.value, t)
	ms.SetParentDirty()
	return ms
}

func (ms *ArrayEditor[V]) AddAll(t []V) *ArrayEditor[V] {
	if len(t) <= 0 {
		return ms
	}
	ms.value = append(ms.value, t...)
	ms.SetParentDirty()
	return ms
}

func (ms *ArrayEditor[V]) Each(f func(i int, v V) bool) {
	for i := range ms.value {
		if !f(i, ms.value[i]) {
			break
		}
	}
}

func (ms *ArrayEditor[V]) Clear() {
	ms.value = make([]V, 0)
	ms.SetParentDirty()
}
func (ms *ArrayEditor[V]) Len() int {
	return len(ms.value)
}

func (ms *ArrayEditor[V]) Remove(v V) {
	var idx = -1
	for i := range ms.value {
		if ms.value[i] == v {
			idx = i
			break
		}
	}
	ms.RemoveByIdx(idx)
}

func (ms *ArrayEditor[V]) RemoveByIdx(idx int) {
	if idx != -1 {
		ms.value = append(ms.value[0:idx], ms.value[idx+1:]...)
		ms.SetParentDirty()
	}
}

func (ms *ArrayEditor[V]) Values() []V {
	vv := make([]V, ms.Len())
	copy(vv, ms.value)
	return vv
}
