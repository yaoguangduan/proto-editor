package pbeditor

import "iter"

type StrU32Map struct {
	parent      Editor
	keyInParent any
	value       map[string]uint32
	dirtied     map[string]struct{}
	deleted     map[string]struct{}
}

func NewStrU32Map() *StrU32Map {
	return &StrU32Map{value: make(map[string]uint32), dirtied: make(map[string]struct{}), deleted: make(map[string]struct{})}
}

func (ms *StrU32Map) SetParent(sync Editor, key any) {
	ms.parent = sync
	ms.keyInParent = key
}
func (ms *StrU32Map) SetParentDirty() {
	if ms.parent != nil {
		ms.parent.SetDirty(ms.keyInParent, true, ms)
	}
}
func (ms *StrU32Map) SetDirty(key any, dirty bool, editor Editor) {
	if dirty {
		ms.dirtied[key.(string)] = struct{}{}
		ms.SetParentDirty()
	} else {
		delete(ms.dirtied, key.(string))
	}

}

func (ms *StrU32Map) FlushDirty(dirty bool) {
	if !dirty {
		ms.dirtied = map[string]struct{}{}
		ms.deleted = map[string]struct{}{}
	}
}

func (ms *StrU32Map) Key() any {
	return nil
}
func (ms *StrU32Map) SetKey(i any) {

}

func (ms *StrU32Map) Put(k string, v uint32) *StrU32Map {
	_, exist := ms.value[k]
	if exist {
		ms.deleted[k] = struct{}{}
	}
	ms.value[k] = v
	ms.dirtied[k] = struct{}{}
	ms.SetParentDirty()
	return ms
}

func (ms *StrU32Map) PutAll(kv map[string]uint32) *StrU32Map {
	for k, v := range kv {
		ms.Put(k, v)
	}
	return ms
}

func (ms *StrU32Map) Len() int {
	return len(ms.value)
}

func (ms *StrU32Map) Clear() *StrU32Map {
	if ms.Len() <= 0 {
		return ms
	}
	for k, _ := range ms.value {
		ms.deleted[k] = struct{}{}
	}
	ms.value = map[string]uint32{}
	ms.dirtied = map[string]struct{}{}
	ms.SetParentDirty()
	return ms
}
func (ms *StrU32Map) Get(k string) (uint32, bool) {
	v, exist := ms.value[k]
	return v, exist
}

func (ms *StrU32Map) Remove(k string) (uint32, bool) {
	v, exist := ms.value[k]
	if !exist {
		return v, false
	}
	delete(ms.value, k)
	delete(ms.dirtied, k)
	ms.deleted[k] = struct{}{}
	ms.SetParentDirty()
	return v, true
}

func (ms *StrU32Map) RemoveAll(k []string) {
	for _, kk := range k {
		ms.Remove(kk)
	}
}
func (ms *StrU32Map) Iter() iter.Seq2[string, uint32] {
	return func(yield func(string, uint32) bool) {
		for k, v := range ms.value {
			if !yield(k, v) {
				return
			}
		}
	}
}

func (ms *StrU32Map) Dirtied() map[string]struct{} {
	return ms.dirtied
}

func (ms *StrU32Map) Deleted() map[string]struct{} {
	return ms.deleted
}

func (ms *StrU32Map) ContainDirtied(kk string) bool {
	for k := range ms.dirtied {
		if k == kk {
			return true
		}
	}
	return false
}

func (ms *StrU32Map) ContainDeleted(kk string) bool {
	for k := range ms.deleted {
		if k == kk {
			return true
		}
	}
	return false
}

func (ms *StrU32Map) MergeDirtyFromPb(delList []string, dirties map[string]uint32) {
	ms.RemoveAll(delList)
	for k, v := range dirties {
		ms.Put(k, v)
	}
}

func (ms *StrU32Map) Entries() map[string]uint32 {
	return ms.value
}

type StrI64Map struct {
	parent      Editor
	keyInParent any
	value       map[string]int64
	dirtied     map[string]struct{}
	deleted     map[string]struct{}
}

func NewStrI64Map() *StrI64Map {
	return &StrI64Map{value: make(map[string]int64), dirtied: make(map[string]struct{}), deleted: make(map[string]struct{})}
}

func (ms *StrI64Map) SetParent(sync Editor, key any) {
	ms.parent = sync
	ms.keyInParent = key
}
func (ms *StrI64Map) SetParentDirty() {
	if ms.parent != nil {
		ms.parent.SetDirty(ms.keyInParent, true, ms)
	}
}
func (ms *StrI64Map) SetDirty(key any, dirty bool, editor Editor) {
	if dirty {
		ms.dirtied[key.(string)] = struct{}{}
		ms.SetParentDirty()
	} else {
		delete(ms.dirtied, key.(string))
	}

}

func (ms *StrI64Map) FlushDirty(dirty bool) {
	if !dirty {
		ms.dirtied = map[string]struct{}{}
		ms.deleted = map[string]struct{}{}
	}
}

func (ms *StrI64Map) Key() any {
	return nil
}
func (ms *StrI64Map) SetKey(i any) {

}

func (ms *StrI64Map) Put(k string, v int64) *StrI64Map {
	_, exist := ms.value[k]
	if exist {
		ms.deleted[k] = struct{}{}
	}
	ms.value[k] = v
	ms.dirtied[k] = struct{}{}
	ms.SetParentDirty()
	return ms
}

func (ms *StrI64Map) PutAll(kv map[string]int64) *StrI64Map {
	for k, v := range kv {
		ms.Put(k, v)
	}
	return ms
}

func (ms *StrI64Map) Len() int {
	return len(ms.value)
}

func (ms *StrI64Map) Clear() *StrI64Map {
	if ms.Len() <= 0 {
		return ms
	}
	for k, _ := range ms.value {
		ms.deleted[k] = struct{}{}
	}
	ms.value = map[string]int64{}
	ms.dirtied = map[string]struct{}{}
	ms.SetParentDirty()
	return ms
}
func (ms *StrI64Map) Get(k string) (int64, bool) {
	v, exist := ms.value[k]
	return v, exist
}

func (ms *StrI64Map) Remove(k string) (int64, bool) {
	v, exist := ms.value[k]
	if !exist {
		return v, false
	}
	delete(ms.value, k)
	delete(ms.dirtied, k)
	ms.deleted[k] = struct{}{}
	ms.SetParentDirty()
	return v, true
}

func (ms *StrI64Map) RemoveAll(k []string) {
	for _, kk := range k {
		ms.Remove(kk)
	}
}
func (ms *StrI64Map) Iter() iter.Seq2[string, int64] {
	return func(yield func(string, int64) bool) {
		for k, v := range ms.value {
			if !yield(k, v) {
				return
			}
		}
	}
}

func (ms *StrI64Map) Dirtied() map[string]struct{} {
	return ms.dirtied
}

func (ms *StrI64Map) Deleted() map[string]struct{} {
	return ms.deleted
}

func (ms *StrI64Map) ContainDirtied(kk string) bool {
	for k := range ms.dirtied {
		if k == kk {
			return true
		}
	}
	return false
}

func (ms *StrI64Map) ContainDeleted(kk string) bool {
	for k := range ms.deleted {
		if k == kk {
			return true
		}
	}
	return false
}

func (ms *StrI64Map) MergeDirtyFromPb(delList []string, dirties map[string]int64) {
	ms.RemoveAll(delList)
	for k, v := range dirties {
		ms.Put(k, v)
	}
}

func (ms *StrI64Map) Entries() map[string]int64 {
	return ms.value
}

type StrStrMap struct {
	parent      Editor
	keyInParent any
	value       map[string]string
	dirtied     map[string]struct{}
	deleted     map[string]struct{}
}

func NewStrStrMap() *StrStrMap {
	return &StrStrMap{value: make(map[string]string), dirtied: make(map[string]struct{}), deleted: make(map[string]struct{})}
}

func (ms *StrStrMap) SetParent(sync Editor, key any) {
	ms.parent = sync
	ms.keyInParent = key
}
func (ms *StrStrMap) SetParentDirty() {
	if ms.parent != nil {
		ms.parent.SetDirty(ms.keyInParent, true, ms)
	}
}
func (ms *StrStrMap) SetDirty(key any, dirty bool, editor Editor) {
	if dirty {
		ms.dirtied[key.(string)] = struct{}{}
		ms.SetParentDirty()
	} else {
		delete(ms.dirtied, key.(string))
	}

}

func (ms *StrStrMap) FlushDirty(dirty bool) {
	if !dirty {
		ms.dirtied = map[string]struct{}{}
		ms.deleted = map[string]struct{}{}
	}
}

func (ms *StrStrMap) Key() any {
	return nil
}
func (ms *StrStrMap) SetKey(i any) {

}

func (ms *StrStrMap) Put(k string, v string) *StrStrMap {
	_, exist := ms.value[k]
	if exist {
		ms.deleted[k] = struct{}{}
	}
	ms.value[k] = v
	ms.dirtied[k] = struct{}{}
	ms.SetParentDirty()
	return ms
}

func (ms *StrStrMap) PutAll(kv map[string]string) *StrStrMap {
	for k, v := range kv {
		ms.Put(k, v)
	}
	return ms
}

func (ms *StrStrMap) Len() int {
	return len(ms.value)
}

func (ms *StrStrMap) Clear() *StrStrMap {
	if ms.Len() <= 0 {
		return ms
	}
	for k, _ := range ms.value {
		ms.deleted[k] = struct{}{}
	}
	ms.value = map[string]string{}
	ms.dirtied = map[string]struct{}{}
	ms.SetParentDirty()
	return ms
}
func (ms *StrStrMap) Get(k string) (string, bool) {
	v, exist := ms.value[k]
	return v, exist
}

func (ms *StrStrMap) Remove(k string) (string, bool) {
	v, exist := ms.value[k]
	if !exist {
		return v, false
	}
	delete(ms.value, k)
	delete(ms.dirtied, k)
	ms.deleted[k] = struct{}{}
	ms.SetParentDirty()
	return v, true
}

func (ms *StrStrMap) RemoveAll(k []string) {
	for _, kk := range k {
		ms.Remove(kk)
	}
}
func (ms *StrStrMap) Iter() iter.Seq2[string, string] {
	return func(yield func(string, string) bool) {
		for k, v := range ms.value {
			if !yield(k, v) {
				return
			}
		}
	}
}

func (ms *StrStrMap) Dirtied() map[string]struct{} {
	return ms.dirtied
}

func (ms *StrStrMap) Deleted() map[string]struct{} {
	return ms.deleted
}

func (ms *StrStrMap) ContainDirtied(kk string) bool {
	for k := range ms.dirtied {
		if k == kk {
			return true
		}
	}
	return false
}

func (ms *StrStrMap) ContainDeleted(kk string) bool {
	for k := range ms.deleted {
		if k == kk {
			return true
		}
	}
	return false
}

func (ms *StrStrMap) MergeDirtyFromPb(delList []string, dirties map[string]string) {
	ms.RemoveAll(delList)
	for k, v := range dirties {
		ms.Put(k, v)
	}
}

func (ms *StrStrMap) Entries() map[string]string {
	return ms.value
}

type StrF32Map struct {
	parent      Editor
	keyInParent any
	value       map[string]float32
	dirtied     map[string]struct{}
	deleted     map[string]struct{}
}

func NewStrF32Map() *StrF32Map {
	return &StrF32Map{value: make(map[string]float32), dirtied: make(map[string]struct{}), deleted: make(map[string]struct{})}
}

func (ms *StrF32Map) SetParent(sync Editor, key any) {
	ms.parent = sync
	ms.keyInParent = key
}
func (ms *StrF32Map) SetParentDirty() {
	if ms.parent != nil {
		ms.parent.SetDirty(ms.keyInParent, true, ms)
	}
}
func (ms *StrF32Map) SetDirty(key any, dirty bool, editor Editor) {
	if dirty {
		ms.dirtied[key.(string)] = struct{}{}
		ms.SetParentDirty()
	} else {
		delete(ms.dirtied, key.(string))
	}

}

func (ms *StrF32Map) FlushDirty(dirty bool) {
	if !dirty {
		ms.dirtied = map[string]struct{}{}
		ms.deleted = map[string]struct{}{}
	}
}

func (ms *StrF32Map) Key() any {
	return nil
}
func (ms *StrF32Map) SetKey(i any) {

}

func (ms *StrF32Map) Put(k string, v float32) *StrF32Map {
	_, exist := ms.value[k]
	if exist {
		ms.deleted[k] = struct{}{}
	}
	ms.value[k] = v
	ms.dirtied[k] = struct{}{}
	ms.SetParentDirty()
	return ms
}

func (ms *StrF32Map) PutAll(kv map[string]float32) *StrF32Map {
	for k, v := range kv {
		ms.Put(k, v)
	}
	return ms
}

func (ms *StrF32Map) Len() int {
	return len(ms.value)
}

func (ms *StrF32Map) Clear() *StrF32Map {
	if ms.Len() <= 0 {
		return ms
	}
	for k, _ := range ms.value {
		ms.deleted[k] = struct{}{}
	}
	ms.value = map[string]float32{}
	ms.dirtied = map[string]struct{}{}
	ms.SetParentDirty()
	return ms
}
func (ms *StrF32Map) Get(k string) (float32, bool) {
	v, exist := ms.value[k]
	return v, exist
}

func (ms *StrF32Map) Remove(k string) (float32, bool) {
	v, exist := ms.value[k]
	if !exist {
		return v, false
	}
	delete(ms.value, k)
	delete(ms.dirtied, k)
	ms.deleted[k] = struct{}{}
	ms.SetParentDirty()
	return v, true
}

func (ms *StrF32Map) RemoveAll(k []string) {
	for _, kk := range k {
		ms.Remove(kk)
	}
}
func (ms *StrF32Map) Iter() iter.Seq2[string, float32] {
	return func(yield func(string, float32) bool) {
		for k, v := range ms.value {
			if !yield(k, v) {
				return
			}
		}
	}
}

func (ms *StrF32Map) Dirtied() map[string]struct{} {
	return ms.dirtied
}

func (ms *StrF32Map) Deleted() map[string]struct{} {
	return ms.deleted
}

func (ms *StrF32Map) ContainDirtied(kk string) bool {
	for k := range ms.dirtied {
		if k == kk {
			return true
		}
	}
	return false
}

func (ms *StrF32Map) ContainDeleted(kk string) bool {
	for k := range ms.deleted {
		if k == kk {
			return true
		}
	}
	return false
}

func (ms *StrF32Map) MergeDirtyFromPb(delList []string, dirties map[string]float32) {
	ms.RemoveAll(delList)
	for k, v := range dirties {
		ms.Put(k, v)
	}
}

func (ms *StrF32Map) Entries() map[string]float32 {
	return ms.value
}

type StrU64Map struct {
	parent      Editor
	keyInParent any
	value       map[string]uint64
	dirtied     map[string]struct{}
	deleted     map[string]struct{}
}

func NewStrU64Map() *StrU64Map {
	return &StrU64Map{value: make(map[string]uint64), dirtied: make(map[string]struct{}), deleted: make(map[string]struct{})}
}

func (ms *StrU64Map) SetParent(sync Editor, key any) {
	ms.parent = sync
	ms.keyInParent = key
}
func (ms *StrU64Map) SetParentDirty() {
	if ms.parent != nil {
		ms.parent.SetDirty(ms.keyInParent, true, ms)
	}
}
func (ms *StrU64Map) SetDirty(key any, dirty bool, editor Editor) {
	if dirty {
		ms.dirtied[key.(string)] = struct{}{}
		ms.SetParentDirty()
	} else {
		delete(ms.dirtied, key.(string))
	}

}

func (ms *StrU64Map) FlushDirty(dirty bool) {
	if !dirty {
		ms.dirtied = map[string]struct{}{}
		ms.deleted = map[string]struct{}{}
	}
}

func (ms *StrU64Map) Key() any {
	return nil
}
func (ms *StrU64Map) SetKey(i any) {

}

func (ms *StrU64Map) Put(k string, v uint64) *StrU64Map {
	_, exist := ms.value[k]
	if exist {
		ms.deleted[k] = struct{}{}
	}
	ms.value[k] = v
	ms.dirtied[k] = struct{}{}
	ms.SetParentDirty()
	return ms
}

func (ms *StrU64Map) PutAll(kv map[string]uint64) *StrU64Map {
	for k, v := range kv {
		ms.Put(k, v)
	}
	return ms
}

func (ms *StrU64Map) Len() int {
	return len(ms.value)
}

func (ms *StrU64Map) Clear() *StrU64Map {
	if ms.Len() <= 0 {
		return ms
	}
	for k, _ := range ms.value {
		ms.deleted[k] = struct{}{}
	}
	ms.value = map[string]uint64{}
	ms.dirtied = map[string]struct{}{}
	ms.SetParentDirty()
	return ms
}
func (ms *StrU64Map) Get(k string) (uint64, bool) {
	v, exist := ms.value[k]
	return v, exist
}

func (ms *StrU64Map) Remove(k string) (uint64, bool) {
	v, exist := ms.value[k]
	if !exist {
		return v, false
	}
	delete(ms.value, k)
	delete(ms.dirtied, k)
	ms.deleted[k] = struct{}{}
	ms.SetParentDirty()
	return v, true
}

func (ms *StrU64Map) RemoveAll(k []string) {
	for _, kk := range k {
		ms.Remove(kk)
	}
}
func (ms *StrU64Map) Iter() iter.Seq2[string, uint64] {
	return func(yield func(string, uint64) bool) {
		for k, v := range ms.value {
			if !yield(k, v) {
				return
			}
		}
	}
}

func (ms *StrU64Map) Dirtied() map[string]struct{} {
	return ms.dirtied
}

func (ms *StrU64Map) Deleted() map[string]struct{} {
	return ms.deleted
}

func (ms *StrU64Map) ContainDirtied(kk string) bool {
	for k := range ms.dirtied {
		if k == kk {
			return true
		}
	}
	return false
}

func (ms *StrU64Map) ContainDeleted(kk string) bool {
	for k := range ms.deleted {
		if k == kk {
			return true
		}
	}
	return false
}

func (ms *StrU64Map) MergeDirtyFromPb(delList []string, dirties map[string]uint64) {
	ms.RemoveAll(delList)
	for k, v := range dirties {
		ms.Put(k, v)
	}
}

func (ms *StrU64Map) Entries() map[string]uint64 {
	return ms.value
}

type StrI32Map struct {
	parent      Editor
	keyInParent any
	value       map[string]int32
	dirtied     map[string]struct{}
	deleted     map[string]struct{}
}

func NewStrI32Map() *StrI32Map {
	return &StrI32Map{value: make(map[string]int32), dirtied: make(map[string]struct{}), deleted: make(map[string]struct{})}
}

func (ms *StrI32Map) SetParent(sync Editor, key any) {
	ms.parent = sync
	ms.keyInParent = key
}
func (ms *StrI32Map) SetParentDirty() {
	if ms.parent != nil {
		ms.parent.SetDirty(ms.keyInParent, true, ms)
	}
}
func (ms *StrI32Map) SetDirty(key any, dirty bool, editor Editor) {
	if dirty {
		ms.dirtied[key.(string)] = struct{}{}
		ms.SetParentDirty()
	} else {
		delete(ms.dirtied, key.(string))
	}

}

func (ms *StrI32Map) FlushDirty(dirty bool) {
	if !dirty {
		ms.dirtied = map[string]struct{}{}
		ms.deleted = map[string]struct{}{}
	}
}

func (ms *StrI32Map) Key() any {
	return nil
}
func (ms *StrI32Map) SetKey(i any) {

}

func (ms *StrI32Map) Put(k string, v int32) *StrI32Map {
	_, exist := ms.value[k]
	if exist {
		ms.deleted[k] = struct{}{}
	}
	ms.value[k] = v
	ms.dirtied[k] = struct{}{}
	ms.SetParentDirty()
	return ms
}

func (ms *StrI32Map) PutAll(kv map[string]int32) *StrI32Map {
	for k, v := range kv {
		ms.Put(k, v)
	}
	return ms
}

func (ms *StrI32Map) Len() int {
	return len(ms.value)
}

func (ms *StrI32Map) Clear() *StrI32Map {
	if ms.Len() <= 0 {
		return ms
	}
	for k, _ := range ms.value {
		ms.deleted[k] = struct{}{}
	}
	ms.value = map[string]int32{}
	ms.dirtied = map[string]struct{}{}
	ms.SetParentDirty()
	return ms
}
func (ms *StrI32Map) Get(k string) (int32, bool) {
	v, exist := ms.value[k]
	return v, exist
}

func (ms *StrI32Map) Remove(k string) (int32, bool) {
	v, exist := ms.value[k]
	if !exist {
		return v, false
	}
	delete(ms.value, k)
	delete(ms.dirtied, k)
	ms.deleted[k] = struct{}{}
	ms.SetParentDirty()
	return v, true
}

func (ms *StrI32Map) RemoveAll(k []string) {
	for _, kk := range k {
		ms.Remove(kk)
	}
}
func (ms *StrI32Map) Iter() iter.Seq2[string, int32] {
	return func(yield func(string, int32) bool) {
		for k, v := range ms.value {
			if !yield(k, v) {
				return
			}
		}
	}
}

func (ms *StrI32Map) Dirtied() map[string]struct{} {
	return ms.dirtied
}

func (ms *StrI32Map) Deleted() map[string]struct{} {
	return ms.deleted
}

func (ms *StrI32Map) ContainDirtied(kk string) bool {
	for k := range ms.dirtied {
		if k == kk {
			return true
		}
	}
	return false
}

func (ms *StrI32Map) ContainDeleted(kk string) bool {
	for k := range ms.deleted {
		if k == kk {
			return true
		}
	}
	return false
}

func (ms *StrI32Map) MergeDirtyFromPb(delList []string, dirties map[string]int32) {
	ms.RemoveAll(delList)
	for k, v := range dirties {
		ms.Put(k, v)
	}
}

func (ms *StrI32Map) Entries() map[string]int32 {
	return ms.value
}

type StrBoolMap struct {
	parent      Editor
	keyInParent any
	value       map[string]bool
	dirtied     map[string]struct{}
	deleted     map[string]struct{}
}

func NewStrBoolMap() *StrBoolMap {
	return &StrBoolMap{value: make(map[string]bool), dirtied: make(map[string]struct{}), deleted: make(map[string]struct{})}
}

func (ms *StrBoolMap) SetParent(sync Editor, key any) {
	ms.parent = sync
	ms.keyInParent = key
}
func (ms *StrBoolMap) SetParentDirty() {
	if ms.parent != nil {
		ms.parent.SetDirty(ms.keyInParent, true, ms)
	}
}
func (ms *StrBoolMap) SetDirty(key any, dirty bool, editor Editor) {
	if dirty {
		ms.dirtied[key.(string)] = struct{}{}
		ms.SetParentDirty()
	} else {
		delete(ms.dirtied, key.(string))
	}

}

func (ms *StrBoolMap) FlushDirty(dirty bool) {
	if !dirty {
		ms.dirtied = map[string]struct{}{}
		ms.deleted = map[string]struct{}{}
	}
}

func (ms *StrBoolMap) Key() any {
	return nil
}
func (ms *StrBoolMap) SetKey(i any) {

}

func (ms *StrBoolMap) Put(k string, v bool) *StrBoolMap {
	_, exist := ms.value[k]
	if exist {
		ms.deleted[k] = struct{}{}
	}
	ms.value[k] = v
	ms.dirtied[k] = struct{}{}
	ms.SetParentDirty()
	return ms
}

func (ms *StrBoolMap) PutAll(kv map[string]bool) *StrBoolMap {
	for k, v := range kv {
		ms.Put(k, v)
	}
	return ms
}

func (ms *StrBoolMap) Len() int {
	return len(ms.value)
}

func (ms *StrBoolMap) Clear() *StrBoolMap {
	if ms.Len() <= 0 {
		return ms
	}
	for k, _ := range ms.value {
		ms.deleted[k] = struct{}{}
	}
	ms.value = map[string]bool{}
	ms.dirtied = map[string]struct{}{}
	ms.SetParentDirty()
	return ms
}
func (ms *StrBoolMap) Get(k string) (bool, bool) {
	v, exist := ms.value[k]
	return v, exist
}

func (ms *StrBoolMap) Remove(k string) (bool, bool) {
	v, exist := ms.value[k]
	if !exist {
		return v, false
	}
	delete(ms.value, k)
	delete(ms.dirtied, k)
	ms.deleted[k] = struct{}{}
	ms.SetParentDirty()
	return v, true
}

func (ms *StrBoolMap) RemoveAll(k []string) {
	for _, kk := range k {
		ms.Remove(kk)
	}
}
func (ms *StrBoolMap) Iter() iter.Seq2[string, bool] {
	return func(yield func(string, bool) bool) {
		for k, v := range ms.value {
			if !yield(k, v) {
				return
			}
		}
	}
}

func (ms *StrBoolMap) Dirtied() map[string]struct{} {
	return ms.dirtied
}

func (ms *StrBoolMap) Deleted() map[string]struct{} {
	return ms.deleted
}

func (ms *StrBoolMap) ContainDirtied(kk string) bool {
	for k := range ms.dirtied {
		if k == kk {
			return true
		}
	}
	return false
}

func (ms *StrBoolMap) ContainDeleted(kk string) bool {
	for k := range ms.deleted {
		if k == kk {
			return true
		}
	}
	return false
}

func (ms *StrBoolMap) MergeDirtyFromPb(delList []string, dirties map[string]bool) {
	ms.RemoveAll(delList)
	for k, v := range dirties {
		ms.Put(k, v)
	}
}

func (ms *StrBoolMap) Entries() map[string]bool {
	return ms.value
}

type StrF64Map struct {
	parent      Editor
	keyInParent any
	value       map[string]float64
	dirtied     map[string]struct{}
	deleted     map[string]struct{}
}

func NewStrF64Map() *StrF64Map {
	return &StrF64Map{value: make(map[string]float64), dirtied: make(map[string]struct{}), deleted: make(map[string]struct{})}
}

func (ms *StrF64Map) SetParent(sync Editor, key any) {
	ms.parent = sync
	ms.keyInParent = key
}
func (ms *StrF64Map) SetParentDirty() {
	if ms.parent != nil {
		ms.parent.SetDirty(ms.keyInParent, true, ms)
	}
}
func (ms *StrF64Map) SetDirty(key any, dirty bool, editor Editor) {
	if dirty {
		ms.dirtied[key.(string)] = struct{}{}
		ms.SetParentDirty()
	} else {
		delete(ms.dirtied, key.(string))
	}

}

func (ms *StrF64Map) FlushDirty(dirty bool) {
	if !dirty {
		ms.dirtied = map[string]struct{}{}
		ms.deleted = map[string]struct{}{}
	}
}

func (ms *StrF64Map) Key() any {
	return nil
}
func (ms *StrF64Map) SetKey(i any) {

}

func (ms *StrF64Map) Put(k string, v float64) *StrF64Map {
	_, exist := ms.value[k]
	if exist {
		ms.deleted[k] = struct{}{}
	}
	ms.value[k] = v
	ms.dirtied[k] = struct{}{}
	ms.SetParentDirty()
	return ms
}

func (ms *StrF64Map) PutAll(kv map[string]float64) *StrF64Map {
	for k, v := range kv {
		ms.Put(k, v)
	}
	return ms
}

func (ms *StrF64Map) Len() int {
	return len(ms.value)
}

func (ms *StrF64Map) Clear() *StrF64Map {
	if ms.Len() <= 0 {
		return ms
	}
	for k, _ := range ms.value {
		ms.deleted[k] = struct{}{}
	}
	ms.value = map[string]float64{}
	ms.dirtied = map[string]struct{}{}
	ms.SetParentDirty()
	return ms
}
func (ms *StrF64Map) Get(k string) (float64, bool) {
	v, exist := ms.value[k]
	return v, exist
}

func (ms *StrF64Map) Remove(k string) (float64, bool) {
	v, exist := ms.value[k]
	if !exist {
		return v, false
	}
	delete(ms.value, k)
	delete(ms.dirtied, k)
	ms.deleted[k] = struct{}{}
	ms.SetParentDirty()
	return v, true
}

func (ms *StrF64Map) RemoveAll(k []string) {
	for _, kk := range k {
		ms.Remove(kk)
	}
}
func (ms *StrF64Map) Iter() iter.Seq2[string, float64] {
	return func(yield func(string, float64) bool) {
		for k, v := range ms.value {
			if !yield(k, v) {
				return
			}
		}
	}
}

func (ms *StrF64Map) Dirtied() map[string]struct{} {
	return ms.dirtied
}

func (ms *StrF64Map) Deleted() map[string]struct{} {
	return ms.deleted
}

func (ms *StrF64Map) ContainDirtied(kk string) bool {
	for k := range ms.dirtied {
		if k == kk {
			return true
		}
	}
	return false
}

func (ms *StrF64Map) ContainDeleted(kk string) bool {
	for k := range ms.deleted {
		if k == kk {
			return true
		}
	}
	return false
}

func (ms *StrF64Map) MergeDirtyFromPb(delList []string, dirties map[string]float64) {
	ms.RemoveAll(delList)
	for k, v := range dirties {
		ms.Put(k, v)
	}
}

func (ms *StrF64Map) Entries() map[string]float64 {
	return ms.value
}

type StrBytesMap struct {
	parent      Editor
	keyInParent any
	value       map[string][]byte
	dirtied     map[string]struct{}
	deleted     map[string]struct{}
}

func NewStrBytesMap() *StrBytesMap {
	return &StrBytesMap{value: make(map[string][]byte), dirtied: make(map[string]struct{}), deleted: make(map[string]struct{})}
}

func (ms *StrBytesMap) SetParent(sync Editor, key any) {
	ms.parent = sync
	ms.keyInParent = key
}
func (ms *StrBytesMap) SetParentDirty() {
	if ms.parent != nil {
		ms.parent.SetDirty(ms.keyInParent, true, ms)
	}
}
func (ms *StrBytesMap) SetDirty(key any, dirty bool, editor Editor) {
	if dirty {
		ms.dirtied[key.(string)] = struct{}{}
		ms.SetParentDirty()
	} else {
		delete(ms.dirtied, key.(string))
	}

}

func (ms *StrBytesMap) FlushDirty(dirty bool) {
	if !dirty {
		ms.dirtied = map[string]struct{}{}
		ms.deleted = map[string]struct{}{}
	}
}

func (ms *StrBytesMap) Key() any {
	return nil
}
func (ms *StrBytesMap) SetKey(i any) {

}

func (ms *StrBytesMap) Put(k string, v []byte) *StrBytesMap {
	_, exist := ms.value[k]
	if exist {
		ms.deleted[k] = struct{}{}
	}
	ms.value[k] = v
	ms.dirtied[k] = struct{}{}
	ms.SetParentDirty()
	return ms
}

func (ms *StrBytesMap) PutAll(kv map[string][]byte) *StrBytesMap {
	for k, v := range kv {
		ms.Put(k, v)
	}
	return ms
}

func (ms *StrBytesMap) Len() int {
	return len(ms.value)
}

func (ms *StrBytesMap) Clear() *StrBytesMap {
	if ms.Len() <= 0 {
		return ms
	}
	for k, _ := range ms.value {
		ms.deleted[k] = struct{}{}
	}
	ms.value = map[string][]byte{}
	ms.dirtied = map[string]struct{}{}
	ms.SetParentDirty()
	return ms
}
func (ms *StrBytesMap) Get(k string) ([]byte, bool) {
	v, exist := ms.value[k]
	return v, exist
}

func (ms *StrBytesMap) Remove(k string) ([]byte, bool) {
	v, exist := ms.value[k]
	if !exist {
		return v, false
	}
	delete(ms.value, k)
	delete(ms.dirtied, k)
	ms.deleted[k] = struct{}{}
	ms.SetParentDirty()
	return v, true
}

func (ms *StrBytesMap) RemoveAll(k []string) {
	for _, kk := range k {
		ms.Remove(kk)
	}
}
func (ms *StrBytesMap) Iter() iter.Seq2[string, []byte] {
	return func(yield func(string, []byte) bool) {
		for k, v := range ms.value {
			if !yield(k, v) {
				return
			}
		}
	}
}

func (ms *StrBytesMap) Dirtied() map[string]struct{} {
	return ms.dirtied
}

func (ms *StrBytesMap) Deleted() map[string]struct{} {
	return ms.deleted
}

func (ms *StrBytesMap) ContainDirtied(kk string) bool {
	for k := range ms.dirtied {
		if k == kk {
			return true
		}
	}
	return false
}

func (ms *StrBytesMap) ContainDeleted(kk string) bool {
	for k := range ms.deleted {
		if k == kk {
			return true
		}
	}
	return false
}

func (ms *StrBytesMap) MergeDirtyFromPb(delList []string, dirties map[string][]byte) {
	ms.RemoveAll(delList)
	for k, v := range dirties {
		ms.Put(k, v)
	}
}

func (ms *StrBytesMap) Entries() map[string][]byte {
	return ms.value
}

type BoolU32Map struct {
	parent      Editor
	keyInParent any
	value       map[bool]uint32
	dirtied     map[bool]struct{}
	deleted     map[bool]struct{}
}

func NewBoolU32Map() *BoolU32Map {
	return &BoolU32Map{value: make(map[bool]uint32), dirtied: make(map[bool]struct{}), deleted: make(map[bool]struct{})}
}

func (ms *BoolU32Map) SetParent(sync Editor, key any) {
	ms.parent = sync
	ms.keyInParent = key
}
func (ms *BoolU32Map) SetParentDirty() {
	if ms.parent != nil {
		ms.parent.SetDirty(ms.keyInParent, true, ms)
	}
}
func (ms *BoolU32Map) SetDirty(key any, dirty bool, editor Editor) {
	if dirty {
		ms.dirtied[key.(bool)] = struct{}{}
		ms.SetParentDirty()
	} else {
		delete(ms.dirtied, key.(bool))
	}

}

func (ms *BoolU32Map) FlushDirty(dirty bool) {
	if !dirty {
		ms.dirtied = map[bool]struct{}{}
		ms.deleted = map[bool]struct{}{}
	}
}

func (ms *BoolU32Map) Key() any {
	return nil
}
func (ms *BoolU32Map) SetKey(i any) {

}

func (ms *BoolU32Map) Put(k bool, v uint32) *BoolU32Map {
	_, exist := ms.value[k]
	if exist {
		ms.deleted[k] = struct{}{}
	}
	ms.value[k] = v
	ms.dirtied[k] = struct{}{}
	ms.SetParentDirty()
	return ms
}

func (ms *BoolU32Map) PutAll(kv map[bool]uint32) *BoolU32Map {
	for k, v := range kv {
		ms.Put(k, v)
	}
	return ms
}

func (ms *BoolU32Map) Len() int {
	return len(ms.value)
}

func (ms *BoolU32Map) Clear() *BoolU32Map {
	if ms.Len() <= 0 {
		return ms
	}
	for k, _ := range ms.value {
		ms.deleted[k] = struct{}{}
	}
	ms.value = map[bool]uint32{}
	ms.dirtied = map[bool]struct{}{}
	ms.SetParentDirty()
	return ms
}
func (ms *BoolU32Map) Get(k bool) (uint32, bool) {
	v, exist := ms.value[k]
	return v, exist
}

func (ms *BoolU32Map) Remove(k bool) (uint32, bool) {
	v, exist := ms.value[k]
	if !exist {
		return v, false
	}
	delete(ms.value, k)
	delete(ms.dirtied, k)
	ms.deleted[k] = struct{}{}
	ms.SetParentDirty()
	return v, true
}

func (ms *BoolU32Map) RemoveAll(k []bool) {
	for _, kk := range k {
		ms.Remove(kk)
	}
}
func (ms *BoolU32Map) Iter() iter.Seq2[bool, uint32] {
	return func(yield func(bool, uint32) bool) {
		for k, v := range ms.value {
			if !yield(k, v) {
				return
			}
		}
	}
}

func (ms *BoolU32Map) Dirtied() map[bool]struct{} {
	return ms.dirtied
}

func (ms *BoolU32Map) Deleted() map[bool]struct{} {
	return ms.deleted
}

func (ms *BoolU32Map) ContainDirtied(kk bool) bool {
	for k := range ms.dirtied {
		if k == kk {
			return true
		}
	}
	return false
}

func (ms *BoolU32Map) ContainDeleted(kk bool) bool {
	for k := range ms.deleted {
		if k == kk {
			return true
		}
	}
	return false
}

func (ms *BoolU32Map) MergeDirtyFromPb(delList []bool, dirties map[bool]uint32) {
	ms.RemoveAll(delList)
	for k, v := range dirties {
		ms.Put(k, v)
	}
}

func (ms *BoolU32Map) Entries() map[bool]uint32 {
	return ms.value
}

type BoolI64Map struct {
	parent      Editor
	keyInParent any
	value       map[bool]int64
	dirtied     map[bool]struct{}
	deleted     map[bool]struct{}
}

func NewBoolI64Map() *BoolI64Map {
	return &BoolI64Map{value: make(map[bool]int64), dirtied: make(map[bool]struct{}), deleted: make(map[bool]struct{})}
}

func (ms *BoolI64Map) SetParent(sync Editor, key any) {
	ms.parent = sync
	ms.keyInParent = key
}
func (ms *BoolI64Map) SetParentDirty() {
	if ms.parent != nil {
		ms.parent.SetDirty(ms.keyInParent, true, ms)
	}
}
func (ms *BoolI64Map) SetDirty(key any, dirty bool, editor Editor) {
	if dirty {
		ms.dirtied[key.(bool)] = struct{}{}
		ms.SetParentDirty()
	} else {
		delete(ms.dirtied, key.(bool))
	}

}

func (ms *BoolI64Map) FlushDirty(dirty bool) {
	if !dirty {
		ms.dirtied = map[bool]struct{}{}
		ms.deleted = map[bool]struct{}{}
	}
}

func (ms *BoolI64Map) Key() any {
	return nil
}
func (ms *BoolI64Map) SetKey(i any) {

}

func (ms *BoolI64Map) Put(k bool, v int64) *BoolI64Map {
	_, exist := ms.value[k]
	if exist {
		ms.deleted[k] = struct{}{}
	}
	ms.value[k] = v
	ms.dirtied[k] = struct{}{}
	ms.SetParentDirty()
	return ms
}

func (ms *BoolI64Map) PutAll(kv map[bool]int64) *BoolI64Map {
	for k, v := range kv {
		ms.Put(k, v)
	}
	return ms
}

func (ms *BoolI64Map) Len() int {
	return len(ms.value)
}

func (ms *BoolI64Map) Clear() *BoolI64Map {
	if ms.Len() <= 0 {
		return ms
	}
	for k, _ := range ms.value {
		ms.deleted[k] = struct{}{}
	}
	ms.value = map[bool]int64{}
	ms.dirtied = map[bool]struct{}{}
	ms.SetParentDirty()
	return ms
}
func (ms *BoolI64Map) Get(k bool) (int64, bool) {
	v, exist := ms.value[k]
	return v, exist
}

func (ms *BoolI64Map) Remove(k bool) (int64, bool) {
	v, exist := ms.value[k]
	if !exist {
		return v, false
	}
	delete(ms.value, k)
	delete(ms.dirtied, k)
	ms.deleted[k] = struct{}{}
	ms.SetParentDirty()
	return v, true
}

func (ms *BoolI64Map) RemoveAll(k []bool) {
	for _, kk := range k {
		ms.Remove(kk)
	}
}
func (ms *BoolI64Map) Iter() iter.Seq2[bool, int64] {
	return func(yield func(bool, int64) bool) {
		for k, v := range ms.value {
			if !yield(k, v) {
				return
			}
		}
	}
}

func (ms *BoolI64Map) Dirtied() map[bool]struct{} {
	return ms.dirtied
}

func (ms *BoolI64Map) Deleted() map[bool]struct{} {
	return ms.deleted
}

func (ms *BoolI64Map) ContainDirtied(kk bool) bool {
	for k := range ms.dirtied {
		if k == kk {
			return true
		}
	}
	return false
}

func (ms *BoolI64Map) ContainDeleted(kk bool) bool {
	for k := range ms.deleted {
		if k == kk {
			return true
		}
	}
	return false
}

func (ms *BoolI64Map) MergeDirtyFromPb(delList []bool, dirties map[bool]int64) {
	ms.RemoveAll(delList)
	for k, v := range dirties {
		ms.Put(k, v)
	}
}

func (ms *BoolI64Map) Entries() map[bool]int64 {
	return ms.value
}

type BoolStrMap struct {
	parent      Editor
	keyInParent any
	value       map[bool]string
	dirtied     map[bool]struct{}
	deleted     map[bool]struct{}
}

func NewBoolStrMap() *BoolStrMap {
	return &BoolStrMap{value: make(map[bool]string), dirtied: make(map[bool]struct{}), deleted: make(map[bool]struct{})}
}

func (ms *BoolStrMap) SetParent(sync Editor, key any) {
	ms.parent = sync
	ms.keyInParent = key
}
func (ms *BoolStrMap) SetParentDirty() {
	if ms.parent != nil {
		ms.parent.SetDirty(ms.keyInParent, true, ms)
	}
}
func (ms *BoolStrMap) SetDirty(key any, dirty bool, editor Editor) {
	if dirty {
		ms.dirtied[key.(bool)] = struct{}{}
		ms.SetParentDirty()
	} else {
		delete(ms.dirtied, key.(bool))
	}

}

func (ms *BoolStrMap) FlushDirty(dirty bool) {
	if !dirty {
		ms.dirtied = map[bool]struct{}{}
		ms.deleted = map[bool]struct{}{}
	}
}

func (ms *BoolStrMap) Key() any {
	return nil
}
func (ms *BoolStrMap) SetKey(i any) {

}

func (ms *BoolStrMap) Put(k bool, v string) *BoolStrMap {
	_, exist := ms.value[k]
	if exist {
		ms.deleted[k] = struct{}{}
	}
	ms.value[k] = v
	ms.dirtied[k] = struct{}{}
	ms.SetParentDirty()
	return ms
}

func (ms *BoolStrMap) PutAll(kv map[bool]string) *BoolStrMap {
	for k, v := range kv {
		ms.Put(k, v)
	}
	return ms
}

func (ms *BoolStrMap) Len() int {
	return len(ms.value)
}

func (ms *BoolStrMap) Clear() *BoolStrMap {
	if ms.Len() <= 0 {
		return ms
	}
	for k, _ := range ms.value {
		ms.deleted[k] = struct{}{}
	}
	ms.value = map[bool]string{}
	ms.dirtied = map[bool]struct{}{}
	ms.SetParentDirty()
	return ms
}
func (ms *BoolStrMap) Get(k bool) (string, bool) {
	v, exist := ms.value[k]
	return v, exist
}

func (ms *BoolStrMap) Remove(k bool) (string, bool) {
	v, exist := ms.value[k]
	if !exist {
		return v, false
	}
	delete(ms.value, k)
	delete(ms.dirtied, k)
	ms.deleted[k] = struct{}{}
	ms.SetParentDirty()
	return v, true
}

func (ms *BoolStrMap) RemoveAll(k []bool) {
	for _, kk := range k {
		ms.Remove(kk)
	}
}
func (ms *BoolStrMap) Iter() iter.Seq2[bool, string] {
	return func(yield func(bool, string) bool) {
		for k, v := range ms.value {
			if !yield(k, v) {
				return
			}
		}
	}
}

func (ms *BoolStrMap) Dirtied() map[bool]struct{} {
	return ms.dirtied
}

func (ms *BoolStrMap) Deleted() map[bool]struct{} {
	return ms.deleted
}

func (ms *BoolStrMap) ContainDirtied(kk bool) bool {
	for k := range ms.dirtied {
		if k == kk {
			return true
		}
	}
	return false
}

func (ms *BoolStrMap) ContainDeleted(kk bool) bool {
	for k := range ms.deleted {
		if k == kk {
			return true
		}
	}
	return false
}

func (ms *BoolStrMap) MergeDirtyFromPb(delList []bool, dirties map[bool]string) {
	ms.RemoveAll(delList)
	for k, v := range dirties {
		ms.Put(k, v)
	}
}

func (ms *BoolStrMap) Entries() map[bool]string {
	return ms.value
}

type BoolF32Map struct {
	parent      Editor
	keyInParent any
	value       map[bool]float32
	dirtied     map[bool]struct{}
	deleted     map[bool]struct{}
}

func NewBoolF32Map() *BoolF32Map {
	return &BoolF32Map{value: make(map[bool]float32), dirtied: make(map[bool]struct{}), deleted: make(map[bool]struct{})}
}

func (ms *BoolF32Map) SetParent(sync Editor, key any) {
	ms.parent = sync
	ms.keyInParent = key
}
func (ms *BoolF32Map) SetParentDirty() {
	if ms.parent != nil {
		ms.parent.SetDirty(ms.keyInParent, true, ms)
	}
}
func (ms *BoolF32Map) SetDirty(key any, dirty bool, editor Editor) {
	if dirty {
		ms.dirtied[key.(bool)] = struct{}{}
		ms.SetParentDirty()
	} else {
		delete(ms.dirtied, key.(bool))
	}

}

func (ms *BoolF32Map) FlushDirty(dirty bool) {
	if !dirty {
		ms.dirtied = map[bool]struct{}{}
		ms.deleted = map[bool]struct{}{}
	}
}

func (ms *BoolF32Map) Key() any {
	return nil
}
func (ms *BoolF32Map) SetKey(i any) {

}

func (ms *BoolF32Map) Put(k bool, v float32) *BoolF32Map {
	_, exist := ms.value[k]
	if exist {
		ms.deleted[k] = struct{}{}
	}
	ms.value[k] = v
	ms.dirtied[k] = struct{}{}
	ms.SetParentDirty()
	return ms
}

func (ms *BoolF32Map) PutAll(kv map[bool]float32) *BoolF32Map {
	for k, v := range kv {
		ms.Put(k, v)
	}
	return ms
}

func (ms *BoolF32Map) Len() int {
	return len(ms.value)
}

func (ms *BoolF32Map) Clear() *BoolF32Map {
	if ms.Len() <= 0 {
		return ms
	}
	for k, _ := range ms.value {
		ms.deleted[k] = struct{}{}
	}
	ms.value = map[bool]float32{}
	ms.dirtied = map[bool]struct{}{}
	ms.SetParentDirty()
	return ms
}
func (ms *BoolF32Map) Get(k bool) (float32, bool) {
	v, exist := ms.value[k]
	return v, exist
}

func (ms *BoolF32Map) Remove(k bool) (float32, bool) {
	v, exist := ms.value[k]
	if !exist {
		return v, false
	}
	delete(ms.value, k)
	delete(ms.dirtied, k)
	ms.deleted[k] = struct{}{}
	ms.SetParentDirty()
	return v, true
}

func (ms *BoolF32Map) RemoveAll(k []bool) {
	for _, kk := range k {
		ms.Remove(kk)
	}
}
func (ms *BoolF32Map) Iter() iter.Seq2[bool, float32] {
	return func(yield func(bool, float32) bool) {
		for k, v := range ms.value {
			if !yield(k, v) {
				return
			}
		}
	}
}

func (ms *BoolF32Map) Dirtied() map[bool]struct{} {
	return ms.dirtied
}

func (ms *BoolF32Map) Deleted() map[bool]struct{} {
	return ms.deleted
}

func (ms *BoolF32Map) ContainDirtied(kk bool) bool {
	for k := range ms.dirtied {
		if k == kk {
			return true
		}
	}
	return false
}

func (ms *BoolF32Map) ContainDeleted(kk bool) bool {
	for k := range ms.deleted {
		if k == kk {
			return true
		}
	}
	return false
}

func (ms *BoolF32Map) MergeDirtyFromPb(delList []bool, dirties map[bool]float32) {
	ms.RemoveAll(delList)
	for k, v := range dirties {
		ms.Put(k, v)
	}
}

func (ms *BoolF32Map) Entries() map[bool]float32 {
	return ms.value
}

type BoolU64Map struct {
	parent      Editor
	keyInParent any
	value       map[bool]uint64
	dirtied     map[bool]struct{}
	deleted     map[bool]struct{}
}

func NewBoolU64Map() *BoolU64Map {
	return &BoolU64Map{value: make(map[bool]uint64), dirtied: make(map[bool]struct{}), deleted: make(map[bool]struct{})}
}

func (ms *BoolU64Map) SetParent(sync Editor, key any) {
	ms.parent = sync
	ms.keyInParent = key
}
func (ms *BoolU64Map) SetParentDirty() {
	if ms.parent != nil {
		ms.parent.SetDirty(ms.keyInParent, true, ms)
	}
}
func (ms *BoolU64Map) SetDirty(key any, dirty bool, editor Editor) {
	if dirty {
		ms.dirtied[key.(bool)] = struct{}{}
		ms.SetParentDirty()
	} else {
		delete(ms.dirtied, key.(bool))
	}

}

func (ms *BoolU64Map) FlushDirty(dirty bool) {
	if !dirty {
		ms.dirtied = map[bool]struct{}{}
		ms.deleted = map[bool]struct{}{}
	}
}

func (ms *BoolU64Map) Key() any {
	return nil
}
func (ms *BoolU64Map) SetKey(i any) {

}

func (ms *BoolU64Map) Put(k bool, v uint64) *BoolU64Map {
	_, exist := ms.value[k]
	if exist {
		ms.deleted[k] = struct{}{}
	}
	ms.value[k] = v
	ms.dirtied[k] = struct{}{}
	ms.SetParentDirty()
	return ms
}

func (ms *BoolU64Map) PutAll(kv map[bool]uint64) *BoolU64Map {
	for k, v := range kv {
		ms.Put(k, v)
	}
	return ms
}

func (ms *BoolU64Map) Len() int {
	return len(ms.value)
}

func (ms *BoolU64Map) Clear() *BoolU64Map {
	if ms.Len() <= 0 {
		return ms
	}
	for k, _ := range ms.value {
		ms.deleted[k] = struct{}{}
	}
	ms.value = map[bool]uint64{}
	ms.dirtied = map[bool]struct{}{}
	ms.SetParentDirty()
	return ms
}
func (ms *BoolU64Map) Get(k bool) (uint64, bool) {
	v, exist := ms.value[k]
	return v, exist
}

func (ms *BoolU64Map) Remove(k bool) (uint64, bool) {
	v, exist := ms.value[k]
	if !exist {
		return v, false
	}
	delete(ms.value, k)
	delete(ms.dirtied, k)
	ms.deleted[k] = struct{}{}
	ms.SetParentDirty()
	return v, true
}

func (ms *BoolU64Map) RemoveAll(k []bool) {
	for _, kk := range k {
		ms.Remove(kk)
	}
}
func (ms *BoolU64Map) Iter() iter.Seq2[bool, uint64] {
	return func(yield func(bool, uint64) bool) {
		for k, v := range ms.value {
			if !yield(k, v) {
				return
			}
		}
	}
}

func (ms *BoolU64Map) Dirtied() map[bool]struct{} {
	return ms.dirtied
}

func (ms *BoolU64Map) Deleted() map[bool]struct{} {
	return ms.deleted
}

func (ms *BoolU64Map) ContainDirtied(kk bool) bool {
	for k := range ms.dirtied {
		if k == kk {
			return true
		}
	}
	return false
}

func (ms *BoolU64Map) ContainDeleted(kk bool) bool {
	for k := range ms.deleted {
		if k == kk {
			return true
		}
	}
	return false
}

func (ms *BoolU64Map) MergeDirtyFromPb(delList []bool, dirties map[bool]uint64) {
	ms.RemoveAll(delList)
	for k, v := range dirties {
		ms.Put(k, v)
	}
}

func (ms *BoolU64Map) Entries() map[bool]uint64 {
	return ms.value
}

type BoolI32Map struct {
	parent      Editor
	keyInParent any
	value       map[bool]int32
	dirtied     map[bool]struct{}
	deleted     map[bool]struct{}
}

func NewBoolI32Map() *BoolI32Map {
	return &BoolI32Map{value: make(map[bool]int32), dirtied: make(map[bool]struct{}), deleted: make(map[bool]struct{})}
}

func (ms *BoolI32Map) SetParent(sync Editor, key any) {
	ms.parent = sync
	ms.keyInParent = key
}
func (ms *BoolI32Map) SetParentDirty() {
	if ms.parent != nil {
		ms.parent.SetDirty(ms.keyInParent, true, ms)
	}
}
func (ms *BoolI32Map) SetDirty(key any, dirty bool, editor Editor) {
	if dirty {
		ms.dirtied[key.(bool)] = struct{}{}
		ms.SetParentDirty()
	} else {
		delete(ms.dirtied, key.(bool))
	}

}

func (ms *BoolI32Map) FlushDirty(dirty bool) {
	if !dirty {
		ms.dirtied = map[bool]struct{}{}
		ms.deleted = map[bool]struct{}{}
	}
}

func (ms *BoolI32Map) Key() any {
	return nil
}
func (ms *BoolI32Map) SetKey(i any) {

}

func (ms *BoolI32Map) Put(k bool, v int32) *BoolI32Map {
	_, exist := ms.value[k]
	if exist {
		ms.deleted[k] = struct{}{}
	}
	ms.value[k] = v
	ms.dirtied[k] = struct{}{}
	ms.SetParentDirty()
	return ms
}

func (ms *BoolI32Map) PutAll(kv map[bool]int32) *BoolI32Map {
	for k, v := range kv {
		ms.Put(k, v)
	}
	return ms
}

func (ms *BoolI32Map) Len() int {
	return len(ms.value)
}

func (ms *BoolI32Map) Clear() *BoolI32Map {
	if ms.Len() <= 0 {
		return ms
	}
	for k, _ := range ms.value {
		ms.deleted[k] = struct{}{}
	}
	ms.value = map[bool]int32{}
	ms.dirtied = map[bool]struct{}{}
	ms.SetParentDirty()
	return ms
}
func (ms *BoolI32Map) Get(k bool) (int32, bool) {
	v, exist := ms.value[k]
	return v, exist
}

func (ms *BoolI32Map) Remove(k bool) (int32, bool) {
	v, exist := ms.value[k]
	if !exist {
		return v, false
	}
	delete(ms.value, k)
	delete(ms.dirtied, k)
	ms.deleted[k] = struct{}{}
	ms.SetParentDirty()
	return v, true
}

func (ms *BoolI32Map) RemoveAll(k []bool) {
	for _, kk := range k {
		ms.Remove(kk)
	}
}
func (ms *BoolI32Map) Iter() iter.Seq2[bool, int32] {
	return func(yield func(bool, int32) bool) {
		for k, v := range ms.value {
			if !yield(k, v) {
				return
			}
		}
	}
}

func (ms *BoolI32Map) Dirtied() map[bool]struct{} {
	return ms.dirtied
}

func (ms *BoolI32Map) Deleted() map[bool]struct{} {
	return ms.deleted
}

func (ms *BoolI32Map) ContainDirtied(kk bool) bool {
	for k := range ms.dirtied {
		if k == kk {
			return true
		}
	}
	return false
}

func (ms *BoolI32Map) ContainDeleted(kk bool) bool {
	for k := range ms.deleted {
		if k == kk {
			return true
		}
	}
	return false
}

func (ms *BoolI32Map) MergeDirtyFromPb(delList []bool, dirties map[bool]int32) {
	ms.RemoveAll(delList)
	for k, v := range dirties {
		ms.Put(k, v)
	}
}

func (ms *BoolI32Map) Entries() map[bool]int32 {
	return ms.value
}

type BoolBoolMap struct {
	parent      Editor
	keyInParent any
	value       map[bool]bool
	dirtied     map[bool]struct{}
	deleted     map[bool]struct{}
}

func NewBoolBoolMap() *BoolBoolMap {
	return &BoolBoolMap{value: make(map[bool]bool), dirtied: make(map[bool]struct{}), deleted: make(map[bool]struct{})}
}

func (ms *BoolBoolMap) SetParent(sync Editor, key any) {
	ms.parent = sync
	ms.keyInParent = key
}
func (ms *BoolBoolMap) SetParentDirty() {
	if ms.parent != nil {
		ms.parent.SetDirty(ms.keyInParent, true, ms)
	}
}
func (ms *BoolBoolMap) SetDirty(key any, dirty bool, editor Editor) {
	if dirty {
		ms.dirtied[key.(bool)] = struct{}{}
		ms.SetParentDirty()
	} else {
		delete(ms.dirtied, key.(bool))
	}

}

func (ms *BoolBoolMap) FlushDirty(dirty bool) {
	if !dirty {
		ms.dirtied = map[bool]struct{}{}
		ms.deleted = map[bool]struct{}{}
	}
}

func (ms *BoolBoolMap) Key() any {
	return nil
}
func (ms *BoolBoolMap) SetKey(i any) {

}

func (ms *BoolBoolMap) Put(k bool, v bool) *BoolBoolMap {
	_, exist := ms.value[k]
	if exist {
		ms.deleted[k] = struct{}{}
	}
	ms.value[k] = v
	ms.dirtied[k] = struct{}{}
	ms.SetParentDirty()
	return ms
}

func (ms *BoolBoolMap) PutAll(kv map[bool]bool) *BoolBoolMap {
	for k, v := range kv {
		ms.Put(k, v)
	}
	return ms
}

func (ms *BoolBoolMap) Len() int {
	return len(ms.value)
}

func (ms *BoolBoolMap) Clear() *BoolBoolMap {
	if ms.Len() <= 0 {
		return ms
	}
	for k, _ := range ms.value {
		ms.deleted[k] = struct{}{}
	}
	ms.value = map[bool]bool{}
	ms.dirtied = map[bool]struct{}{}
	ms.SetParentDirty()
	return ms
}
func (ms *BoolBoolMap) Get(k bool) (bool, bool) {
	v, exist := ms.value[k]
	return v, exist
}

func (ms *BoolBoolMap) Remove(k bool) (bool, bool) {
	v, exist := ms.value[k]
	if !exist {
		return v, false
	}
	delete(ms.value, k)
	delete(ms.dirtied, k)
	ms.deleted[k] = struct{}{}
	ms.SetParentDirty()
	return v, true
}

func (ms *BoolBoolMap) RemoveAll(k []bool) {
	for _, kk := range k {
		ms.Remove(kk)
	}
}
func (ms *BoolBoolMap) Iter() iter.Seq2[bool, bool] {
	return func(yield func(bool, bool) bool) {
		for k, v := range ms.value {
			if !yield(k, v) {
				return
			}
		}
	}
}

func (ms *BoolBoolMap) Dirtied() map[bool]struct{} {
	return ms.dirtied
}

func (ms *BoolBoolMap) Deleted() map[bool]struct{} {
	return ms.deleted
}

func (ms *BoolBoolMap) ContainDirtied(kk bool) bool {
	for k := range ms.dirtied {
		if k == kk {
			return true
		}
	}
	return false
}

func (ms *BoolBoolMap) ContainDeleted(kk bool) bool {
	for k := range ms.deleted {
		if k == kk {
			return true
		}
	}
	return false
}

func (ms *BoolBoolMap) MergeDirtyFromPb(delList []bool, dirties map[bool]bool) {
	ms.RemoveAll(delList)
	for k, v := range dirties {
		ms.Put(k, v)
	}
}

func (ms *BoolBoolMap) Entries() map[bool]bool {
	return ms.value
}

type BoolF64Map struct {
	parent      Editor
	keyInParent any
	value       map[bool]float64
	dirtied     map[bool]struct{}
	deleted     map[bool]struct{}
}

func NewBoolF64Map() *BoolF64Map {
	return &BoolF64Map{value: make(map[bool]float64), dirtied: make(map[bool]struct{}), deleted: make(map[bool]struct{})}
}

func (ms *BoolF64Map) SetParent(sync Editor, key any) {
	ms.parent = sync
	ms.keyInParent = key
}
func (ms *BoolF64Map) SetParentDirty() {
	if ms.parent != nil {
		ms.parent.SetDirty(ms.keyInParent, true, ms)
	}
}
func (ms *BoolF64Map) SetDirty(key any, dirty bool, editor Editor) {
	if dirty {
		ms.dirtied[key.(bool)] = struct{}{}
		ms.SetParentDirty()
	} else {
		delete(ms.dirtied, key.(bool))
	}

}

func (ms *BoolF64Map) FlushDirty(dirty bool) {
	if !dirty {
		ms.dirtied = map[bool]struct{}{}
		ms.deleted = map[bool]struct{}{}
	}
}

func (ms *BoolF64Map) Key() any {
	return nil
}
func (ms *BoolF64Map) SetKey(i any) {

}

func (ms *BoolF64Map) Put(k bool, v float64) *BoolF64Map {
	_, exist := ms.value[k]
	if exist {
		ms.deleted[k] = struct{}{}
	}
	ms.value[k] = v
	ms.dirtied[k] = struct{}{}
	ms.SetParentDirty()
	return ms
}

func (ms *BoolF64Map) PutAll(kv map[bool]float64) *BoolF64Map {
	for k, v := range kv {
		ms.Put(k, v)
	}
	return ms
}

func (ms *BoolF64Map) Len() int {
	return len(ms.value)
}

func (ms *BoolF64Map) Clear() *BoolF64Map {
	if ms.Len() <= 0 {
		return ms
	}
	for k, _ := range ms.value {
		ms.deleted[k] = struct{}{}
	}
	ms.value = map[bool]float64{}
	ms.dirtied = map[bool]struct{}{}
	ms.SetParentDirty()
	return ms
}
func (ms *BoolF64Map) Get(k bool) (float64, bool) {
	v, exist := ms.value[k]
	return v, exist
}

func (ms *BoolF64Map) Remove(k bool) (float64, bool) {
	v, exist := ms.value[k]
	if !exist {
		return v, false
	}
	delete(ms.value, k)
	delete(ms.dirtied, k)
	ms.deleted[k] = struct{}{}
	ms.SetParentDirty()
	return v, true
}

func (ms *BoolF64Map) RemoveAll(k []bool) {
	for _, kk := range k {
		ms.Remove(kk)
	}
}
func (ms *BoolF64Map) Iter() iter.Seq2[bool, float64] {
	return func(yield func(bool, float64) bool) {
		for k, v := range ms.value {
			if !yield(k, v) {
				return
			}
		}
	}
}

func (ms *BoolF64Map) Dirtied() map[bool]struct{} {
	return ms.dirtied
}

func (ms *BoolF64Map) Deleted() map[bool]struct{} {
	return ms.deleted
}

func (ms *BoolF64Map) ContainDirtied(kk bool) bool {
	for k := range ms.dirtied {
		if k == kk {
			return true
		}
	}
	return false
}

func (ms *BoolF64Map) ContainDeleted(kk bool) bool {
	for k := range ms.deleted {
		if k == kk {
			return true
		}
	}
	return false
}

func (ms *BoolF64Map) MergeDirtyFromPb(delList []bool, dirties map[bool]float64) {
	ms.RemoveAll(delList)
	for k, v := range dirties {
		ms.Put(k, v)
	}
}

func (ms *BoolF64Map) Entries() map[bool]float64 {
	return ms.value
}

type BoolBytesMap struct {
	parent      Editor
	keyInParent any
	value       map[bool][]byte
	dirtied     map[bool]struct{}
	deleted     map[bool]struct{}
}

func NewBoolBytesMap() *BoolBytesMap {
	return &BoolBytesMap{value: make(map[bool][]byte), dirtied: make(map[bool]struct{}), deleted: make(map[bool]struct{})}
}

func (ms *BoolBytesMap) SetParent(sync Editor, key any) {
	ms.parent = sync
	ms.keyInParent = key
}
func (ms *BoolBytesMap) SetParentDirty() {
	if ms.parent != nil {
		ms.parent.SetDirty(ms.keyInParent, true, ms)
	}
}
func (ms *BoolBytesMap) SetDirty(key any, dirty bool, editor Editor) {
	if dirty {
		ms.dirtied[key.(bool)] = struct{}{}
		ms.SetParentDirty()
	} else {
		delete(ms.dirtied, key.(bool))
	}

}

func (ms *BoolBytesMap) FlushDirty(dirty bool) {
	if !dirty {
		ms.dirtied = map[bool]struct{}{}
		ms.deleted = map[bool]struct{}{}
	}
}

func (ms *BoolBytesMap) Key() any {
	return nil
}
func (ms *BoolBytesMap) SetKey(i any) {

}

func (ms *BoolBytesMap) Put(k bool, v []byte) *BoolBytesMap {
	_, exist := ms.value[k]
	if exist {
		ms.deleted[k] = struct{}{}
	}
	ms.value[k] = v
	ms.dirtied[k] = struct{}{}
	ms.SetParentDirty()
	return ms
}

func (ms *BoolBytesMap) PutAll(kv map[bool][]byte) *BoolBytesMap {
	for k, v := range kv {
		ms.Put(k, v)
	}
	return ms
}

func (ms *BoolBytesMap) Len() int {
	return len(ms.value)
}

func (ms *BoolBytesMap) Clear() *BoolBytesMap {
	if ms.Len() <= 0 {
		return ms
	}
	for k, _ := range ms.value {
		ms.deleted[k] = struct{}{}
	}
	ms.value = map[bool][]byte{}
	ms.dirtied = map[bool]struct{}{}
	ms.SetParentDirty()
	return ms
}
func (ms *BoolBytesMap) Get(k bool) ([]byte, bool) {
	v, exist := ms.value[k]
	return v, exist
}

func (ms *BoolBytesMap) Remove(k bool) ([]byte, bool) {
	v, exist := ms.value[k]
	if !exist {
		return v, false
	}
	delete(ms.value, k)
	delete(ms.dirtied, k)
	ms.deleted[k] = struct{}{}
	ms.SetParentDirty()
	return v, true
}

func (ms *BoolBytesMap) RemoveAll(k []bool) {
	for _, kk := range k {
		ms.Remove(kk)
	}
}
func (ms *BoolBytesMap) Iter() iter.Seq2[bool, []byte] {
	return func(yield func(bool, []byte) bool) {
		for k, v := range ms.value {
			if !yield(k, v) {
				return
			}
		}
	}
}

func (ms *BoolBytesMap) Dirtied() map[bool]struct{} {
	return ms.dirtied
}

func (ms *BoolBytesMap) Deleted() map[bool]struct{} {
	return ms.deleted
}

func (ms *BoolBytesMap) ContainDirtied(kk bool) bool {
	for k := range ms.dirtied {
		if k == kk {
			return true
		}
	}
	return false
}

func (ms *BoolBytesMap) ContainDeleted(kk bool) bool {
	for k := range ms.deleted {
		if k == kk {
			return true
		}
	}
	return false
}

func (ms *BoolBytesMap) MergeDirtyFromPb(delList []bool, dirties map[bool][]byte) {
	ms.RemoveAll(delList)
	for k, v := range dirties {
		ms.Put(k, v)
	}
}

func (ms *BoolBytesMap) Entries() map[bool][]byte {
	return ms.value
}

type U32U32Map struct {
	parent      Editor
	keyInParent any
	value       map[uint32]uint32
	dirtied     map[uint32]struct{}
	deleted     map[uint32]struct{}
}

func NewU32U32Map() *U32U32Map {
	return &U32U32Map{value: make(map[uint32]uint32), dirtied: make(map[uint32]struct{}), deleted: make(map[uint32]struct{})}
}

func (ms *U32U32Map) SetParent(sync Editor, key any) {
	ms.parent = sync
	ms.keyInParent = key
}
func (ms *U32U32Map) SetParentDirty() {
	if ms.parent != nil {
		ms.parent.SetDirty(ms.keyInParent, true, ms)
	}
}
func (ms *U32U32Map) SetDirty(key any, dirty bool, editor Editor) {
	if dirty {
		ms.dirtied[key.(uint32)] = struct{}{}
		ms.SetParentDirty()
	} else {
		delete(ms.dirtied, key.(uint32))
	}

}

func (ms *U32U32Map) FlushDirty(dirty bool) {
	if !dirty {
		ms.dirtied = map[uint32]struct{}{}
		ms.deleted = map[uint32]struct{}{}
	}
}

func (ms *U32U32Map) Key() any {
	return nil
}
func (ms *U32U32Map) SetKey(i any) {

}

func (ms *U32U32Map) Put(k uint32, v uint32) *U32U32Map {
	_, exist := ms.value[k]
	if exist {
		ms.deleted[k] = struct{}{}
	}
	ms.value[k] = v
	ms.dirtied[k] = struct{}{}
	ms.SetParentDirty()
	return ms
}

func (ms *U32U32Map) PutAll(kv map[uint32]uint32) *U32U32Map {
	for k, v := range kv {
		ms.Put(k, v)
	}
	return ms
}

func (ms *U32U32Map) Len() int {
	return len(ms.value)
}

func (ms *U32U32Map) Clear() *U32U32Map {
	if ms.Len() <= 0 {
		return ms
	}
	for k, _ := range ms.value {
		ms.deleted[k] = struct{}{}
	}
	ms.value = map[uint32]uint32{}
	ms.dirtied = map[uint32]struct{}{}
	ms.SetParentDirty()
	return ms
}
func (ms *U32U32Map) Get(k uint32) (uint32, bool) {
	v, exist := ms.value[k]
	return v, exist
}

func (ms *U32U32Map) Remove(k uint32) (uint32, bool) {
	v, exist := ms.value[k]
	if !exist {
		return v, false
	}
	delete(ms.value, k)
	delete(ms.dirtied, k)
	ms.deleted[k] = struct{}{}
	ms.SetParentDirty()
	return v, true
}

func (ms *U32U32Map) RemoveAll(k []uint32) {
	for _, kk := range k {
		ms.Remove(kk)
	}
}
func (ms *U32U32Map) Iter() iter.Seq2[uint32, uint32] {
	return func(yield func(uint32, uint32) bool) {
		for k, v := range ms.value {
			if !yield(k, v) {
				return
			}
		}
	}
}

func (ms *U32U32Map) Dirtied() map[uint32]struct{} {
	return ms.dirtied
}

func (ms *U32U32Map) Deleted() map[uint32]struct{} {
	return ms.deleted
}

func (ms *U32U32Map) ContainDirtied(kk uint32) bool {
	for k := range ms.dirtied {
		if k == kk {
			return true
		}
	}
	return false
}

func (ms *U32U32Map) ContainDeleted(kk uint32) bool {
	for k := range ms.deleted {
		if k == kk {
			return true
		}
	}
	return false
}

func (ms *U32U32Map) MergeDirtyFromPb(delList []uint32, dirties map[uint32]uint32) {
	ms.RemoveAll(delList)
	for k, v := range dirties {
		ms.Put(k, v)
	}
}

func (ms *U32U32Map) Entries() map[uint32]uint32 {
	return ms.value
}

type U32I64Map struct {
	parent      Editor
	keyInParent any
	value       map[uint32]int64
	dirtied     map[uint32]struct{}
	deleted     map[uint32]struct{}
}

func NewU32I64Map() *U32I64Map {
	return &U32I64Map{value: make(map[uint32]int64), dirtied: make(map[uint32]struct{}), deleted: make(map[uint32]struct{})}
}

func (ms *U32I64Map) SetParent(sync Editor, key any) {
	ms.parent = sync
	ms.keyInParent = key
}
func (ms *U32I64Map) SetParentDirty() {
	if ms.parent != nil {
		ms.parent.SetDirty(ms.keyInParent, true, ms)
	}
}
func (ms *U32I64Map) SetDirty(key any, dirty bool, editor Editor) {
	if dirty {
		ms.dirtied[key.(uint32)] = struct{}{}
		ms.SetParentDirty()
	} else {
		delete(ms.dirtied, key.(uint32))
	}

}

func (ms *U32I64Map) FlushDirty(dirty bool) {
	if !dirty {
		ms.dirtied = map[uint32]struct{}{}
		ms.deleted = map[uint32]struct{}{}
	}
}

func (ms *U32I64Map) Key() any {
	return nil
}
func (ms *U32I64Map) SetKey(i any) {

}

func (ms *U32I64Map) Put(k uint32, v int64) *U32I64Map {
	_, exist := ms.value[k]
	if exist {
		ms.deleted[k] = struct{}{}
	}
	ms.value[k] = v
	ms.dirtied[k] = struct{}{}
	ms.SetParentDirty()
	return ms
}

func (ms *U32I64Map) PutAll(kv map[uint32]int64) *U32I64Map {
	for k, v := range kv {
		ms.Put(k, v)
	}
	return ms
}

func (ms *U32I64Map) Len() int {
	return len(ms.value)
}

func (ms *U32I64Map) Clear() *U32I64Map {
	if ms.Len() <= 0 {
		return ms
	}
	for k, _ := range ms.value {
		ms.deleted[k] = struct{}{}
	}
	ms.value = map[uint32]int64{}
	ms.dirtied = map[uint32]struct{}{}
	ms.SetParentDirty()
	return ms
}
func (ms *U32I64Map) Get(k uint32) (int64, bool) {
	v, exist := ms.value[k]
	return v, exist
}

func (ms *U32I64Map) Remove(k uint32) (int64, bool) {
	v, exist := ms.value[k]
	if !exist {
		return v, false
	}
	delete(ms.value, k)
	delete(ms.dirtied, k)
	ms.deleted[k] = struct{}{}
	ms.SetParentDirty()
	return v, true
}

func (ms *U32I64Map) RemoveAll(k []uint32) {
	for _, kk := range k {
		ms.Remove(kk)
	}
}
func (ms *U32I64Map) Iter() iter.Seq2[uint32, int64] {
	return func(yield func(uint32, int64) bool) {
		for k, v := range ms.value {
			if !yield(k, v) {
				return
			}
		}
	}
}

func (ms *U32I64Map) Dirtied() map[uint32]struct{} {
	return ms.dirtied
}

func (ms *U32I64Map) Deleted() map[uint32]struct{} {
	return ms.deleted
}

func (ms *U32I64Map) ContainDirtied(kk uint32) bool {
	for k := range ms.dirtied {
		if k == kk {
			return true
		}
	}
	return false
}

func (ms *U32I64Map) ContainDeleted(kk uint32) bool {
	for k := range ms.deleted {
		if k == kk {
			return true
		}
	}
	return false
}

func (ms *U32I64Map) MergeDirtyFromPb(delList []uint32, dirties map[uint32]int64) {
	ms.RemoveAll(delList)
	for k, v := range dirties {
		ms.Put(k, v)
	}
}

func (ms *U32I64Map) Entries() map[uint32]int64 {
	return ms.value
}

type U32StrMap struct {
	parent      Editor
	keyInParent any
	value       map[uint32]string
	dirtied     map[uint32]struct{}
	deleted     map[uint32]struct{}
}

func NewU32StrMap() *U32StrMap {
	return &U32StrMap{value: make(map[uint32]string), dirtied: make(map[uint32]struct{}), deleted: make(map[uint32]struct{})}
}

func (ms *U32StrMap) SetParent(sync Editor, key any) {
	ms.parent = sync
	ms.keyInParent = key
}
func (ms *U32StrMap) SetParentDirty() {
	if ms.parent != nil {
		ms.parent.SetDirty(ms.keyInParent, true, ms)
	}
}
func (ms *U32StrMap) SetDirty(key any, dirty bool, editor Editor) {
	if dirty {
		ms.dirtied[key.(uint32)] = struct{}{}
		ms.SetParentDirty()
	} else {
		delete(ms.dirtied, key.(uint32))
	}

}

func (ms *U32StrMap) FlushDirty(dirty bool) {
	if !dirty {
		ms.dirtied = map[uint32]struct{}{}
		ms.deleted = map[uint32]struct{}{}
	}
}

func (ms *U32StrMap) Key() any {
	return nil
}
func (ms *U32StrMap) SetKey(i any) {

}

func (ms *U32StrMap) Put(k uint32, v string) *U32StrMap {
	_, exist := ms.value[k]
	if exist {
		ms.deleted[k] = struct{}{}
	}
	ms.value[k] = v
	ms.dirtied[k] = struct{}{}
	ms.SetParentDirty()
	return ms
}

func (ms *U32StrMap) PutAll(kv map[uint32]string) *U32StrMap {
	for k, v := range kv {
		ms.Put(k, v)
	}
	return ms
}

func (ms *U32StrMap) Len() int {
	return len(ms.value)
}

func (ms *U32StrMap) Clear() *U32StrMap {
	if ms.Len() <= 0 {
		return ms
	}
	for k, _ := range ms.value {
		ms.deleted[k] = struct{}{}
	}
	ms.value = map[uint32]string{}
	ms.dirtied = map[uint32]struct{}{}
	ms.SetParentDirty()
	return ms
}
func (ms *U32StrMap) Get(k uint32) (string, bool) {
	v, exist := ms.value[k]
	return v, exist
}

func (ms *U32StrMap) Remove(k uint32) (string, bool) {
	v, exist := ms.value[k]
	if !exist {
		return v, false
	}
	delete(ms.value, k)
	delete(ms.dirtied, k)
	ms.deleted[k] = struct{}{}
	ms.SetParentDirty()
	return v, true
}

func (ms *U32StrMap) RemoveAll(k []uint32) {
	for _, kk := range k {
		ms.Remove(kk)
	}
}
func (ms *U32StrMap) Iter() iter.Seq2[uint32, string] {
	return func(yield func(uint32, string) bool) {
		for k, v := range ms.value {
			if !yield(k, v) {
				return
			}
		}
	}
}

func (ms *U32StrMap) Dirtied() map[uint32]struct{} {
	return ms.dirtied
}

func (ms *U32StrMap) Deleted() map[uint32]struct{} {
	return ms.deleted
}

func (ms *U32StrMap) ContainDirtied(kk uint32) bool {
	for k := range ms.dirtied {
		if k == kk {
			return true
		}
	}
	return false
}

func (ms *U32StrMap) ContainDeleted(kk uint32) bool {
	for k := range ms.deleted {
		if k == kk {
			return true
		}
	}
	return false
}

func (ms *U32StrMap) MergeDirtyFromPb(delList []uint32, dirties map[uint32]string) {
	ms.RemoveAll(delList)
	for k, v := range dirties {
		ms.Put(k, v)
	}
}

func (ms *U32StrMap) Entries() map[uint32]string {
	return ms.value
}

type U32F32Map struct {
	parent      Editor
	keyInParent any
	value       map[uint32]float32
	dirtied     map[uint32]struct{}
	deleted     map[uint32]struct{}
}

func NewU32F32Map() *U32F32Map {
	return &U32F32Map{value: make(map[uint32]float32), dirtied: make(map[uint32]struct{}), deleted: make(map[uint32]struct{})}
}

func (ms *U32F32Map) SetParent(sync Editor, key any) {
	ms.parent = sync
	ms.keyInParent = key
}
func (ms *U32F32Map) SetParentDirty() {
	if ms.parent != nil {
		ms.parent.SetDirty(ms.keyInParent, true, ms)
	}
}
func (ms *U32F32Map) SetDirty(key any, dirty bool, editor Editor) {
	if dirty {
		ms.dirtied[key.(uint32)] = struct{}{}
		ms.SetParentDirty()
	} else {
		delete(ms.dirtied, key.(uint32))
	}

}

func (ms *U32F32Map) FlushDirty(dirty bool) {
	if !dirty {
		ms.dirtied = map[uint32]struct{}{}
		ms.deleted = map[uint32]struct{}{}
	}
}

func (ms *U32F32Map) Key() any {
	return nil
}
func (ms *U32F32Map) SetKey(i any) {

}

func (ms *U32F32Map) Put(k uint32, v float32) *U32F32Map {
	_, exist := ms.value[k]
	if exist {
		ms.deleted[k] = struct{}{}
	}
	ms.value[k] = v
	ms.dirtied[k] = struct{}{}
	ms.SetParentDirty()
	return ms
}

func (ms *U32F32Map) PutAll(kv map[uint32]float32) *U32F32Map {
	for k, v := range kv {
		ms.Put(k, v)
	}
	return ms
}

func (ms *U32F32Map) Len() int {
	return len(ms.value)
}

func (ms *U32F32Map) Clear() *U32F32Map {
	if ms.Len() <= 0 {
		return ms
	}
	for k, _ := range ms.value {
		ms.deleted[k] = struct{}{}
	}
	ms.value = map[uint32]float32{}
	ms.dirtied = map[uint32]struct{}{}
	ms.SetParentDirty()
	return ms
}
func (ms *U32F32Map) Get(k uint32) (float32, bool) {
	v, exist := ms.value[k]
	return v, exist
}

func (ms *U32F32Map) Remove(k uint32) (float32, bool) {
	v, exist := ms.value[k]
	if !exist {
		return v, false
	}
	delete(ms.value, k)
	delete(ms.dirtied, k)
	ms.deleted[k] = struct{}{}
	ms.SetParentDirty()
	return v, true
}

func (ms *U32F32Map) RemoveAll(k []uint32) {
	for _, kk := range k {
		ms.Remove(kk)
	}
}
func (ms *U32F32Map) Iter() iter.Seq2[uint32, float32] {
	return func(yield func(uint32, float32) bool) {
		for k, v := range ms.value {
			if !yield(k, v) {
				return
			}
		}
	}
}

func (ms *U32F32Map) Dirtied() map[uint32]struct{} {
	return ms.dirtied
}

func (ms *U32F32Map) Deleted() map[uint32]struct{} {
	return ms.deleted
}

func (ms *U32F32Map) ContainDirtied(kk uint32) bool {
	for k := range ms.dirtied {
		if k == kk {
			return true
		}
	}
	return false
}

func (ms *U32F32Map) ContainDeleted(kk uint32) bool {
	for k := range ms.deleted {
		if k == kk {
			return true
		}
	}
	return false
}

func (ms *U32F32Map) MergeDirtyFromPb(delList []uint32, dirties map[uint32]float32) {
	ms.RemoveAll(delList)
	for k, v := range dirties {
		ms.Put(k, v)
	}
}

func (ms *U32F32Map) Entries() map[uint32]float32 {
	return ms.value
}

type U32U64Map struct {
	parent      Editor
	keyInParent any
	value       map[uint32]uint64
	dirtied     map[uint32]struct{}
	deleted     map[uint32]struct{}
}

func NewU32U64Map() *U32U64Map {
	return &U32U64Map{value: make(map[uint32]uint64), dirtied: make(map[uint32]struct{}), deleted: make(map[uint32]struct{})}
}

func (ms *U32U64Map) SetParent(sync Editor, key any) {
	ms.parent = sync
	ms.keyInParent = key
}
func (ms *U32U64Map) SetParentDirty() {
	if ms.parent != nil {
		ms.parent.SetDirty(ms.keyInParent, true, ms)
	}
}
func (ms *U32U64Map) SetDirty(key any, dirty bool, editor Editor) {
	if dirty {
		ms.dirtied[key.(uint32)] = struct{}{}
		ms.SetParentDirty()
	} else {
		delete(ms.dirtied, key.(uint32))
	}

}

func (ms *U32U64Map) FlushDirty(dirty bool) {
	if !dirty {
		ms.dirtied = map[uint32]struct{}{}
		ms.deleted = map[uint32]struct{}{}
	}
}

func (ms *U32U64Map) Key() any {
	return nil
}
func (ms *U32U64Map) SetKey(i any) {

}

func (ms *U32U64Map) Put(k uint32, v uint64) *U32U64Map {
	_, exist := ms.value[k]
	if exist {
		ms.deleted[k] = struct{}{}
	}
	ms.value[k] = v
	ms.dirtied[k] = struct{}{}
	ms.SetParentDirty()
	return ms
}

func (ms *U32U64Map) PutAll(kv map[uint32]uint64) *U32U64Map {
	for k, v := range kv {
		ms.Put(k, v)
	}
	return ms
}

func (ms *U32U64Map) Len() int {
	return len(ms.value)
}

func (ms *U32U64Map) Clear() *U32U64Map {
	if ms.Len() <= 0 {
		return ms
	}
	for k, _ := range ms.value {
		ms.deleted[k] = struct{}{}
	}
	ms.value = map[uint32]uint64{}
	ms.dirtied = map[uint32]struct{}{}
	ms.SetParentDirty()
	return ms
}
func (ms *U32U64Map) Get(k uint32) (uint64, bool) {
	v, exist := ms.value[k]
	return v, exist
}

func (ms *U32U64Map) Remove(k uint32) (uint64, bool) {
	v, exist := ms.value[k]
	if !exist {
		return v, false
	}
	delete(ms.value, k)
	delete(ms.dirtied, k)
	ms.deleted[k] = struct{}{}
	ms.SetParentDirty()
	return v, true
}

func (ms *U32U64Map) RemoveAll(k []uint32) {
	for _, kk := range k {
		ms.Remove(kk)
	}
}
func (ms *U32U64Map) Iter() iter.Seq2[uint32, uint64] {
	return func(yield func(uint32, uint64) bool) {
		for k, v := range ms.value {
			if !yield(k, v) {
				return
			}
		}
	}
}

func (ms *U32U64Map) Dirtied() map[uint32]struct{} {
	return ms.dirtied
}

func (ms *U32U64Map) Deleted() map[uint32]struct{} {
	return ms.deleted
}

func (ms *U32U64Map) ContainDirtied(kk uint32) bool {
	for k := range ms.dirtied {
		if k == kk {
			return true
		}
	}
	return false
}

func (ms *U32U64Map) ContainDeleted(kk uint32) bool {
	for k := range ms.deleted {
		if k == kk {
			return true
		}
	}
	return false
}

func (ms *U32U64Map) MergeDirtyFromPb(delList []uint32, dirties map[uint32]uint64) {
	ms.RemoveAll(delList)
	for k, v := range dirties {
		ms.Put(k, v)
	}
}

func (ms *U32U64Map) Entries() map[uint32]uint64 {
	return ms.value
}

type U32I32Map struct {
	parent      Editor
	keyInParent any
	value       map[uint32]int32
	dirtied     map[uint32]struct{}
	deleted     map[uint32]struct{}
}

func NewU32I32Map() *U32I32Map {
	return &U32I32Map{value: make(map[uint32]int32), dirtied: make(map[uint32]struct{}), deleted: make(map[uint32]struct{})}
}

func (ms *U32I32Map) SetParent(sync Editor, key any) {
	ms.parent = sync
	ms.keyInParent = key
}
func (ms *U32I32Map) SetParentDirty() {
	if ms.parent != nil {
		ms.parent.SetDirty(ms.keyInParent, true, ms)
	}
}
func (ms *U32I32Map) SetDirty(key any, dirty bool, editor Editor) {
	if dirty {
		ms.dirtied[key.(uint32)] = struct{}{}
		ms.SetParentDirty()
	} else {
		delete(ms.dirtied, key.(uint32))
	}

}

func (ms *U32I32Map) FlushDirty(dirty bool) {
	if !dirty {
		ms.dirtied = map[uint32]struct{}{}
		ms.deleted = map[uint32]struct{}{}
	}
}

func (ms *U32I32Map) Key() any {
	return nil
}
func (ms *U32I32Map) SetKey(i any) {

}

func (ms *U32I32Map) Put(k uint32, v int32) *U32I32Map {
	_, exist := ms.value[k]
	if exist {
		ms.deleted[k] = struct{}{}
	}
	ms.value[k] = v
	ms.dirtied[k] = struct{}{}
	ms.SetParentDirty()
	return ms
}

func (ms *U32I32Map) PutAll(kv map[uint32]int32) *U32I32Map {
	for k, v := range kv {
		ms.Put(k, v)
	}
	return ms
}

func (ms *U32I32Map) Len() int {
	return len(ms.value)
}

func (ms *U32I32Map) Clear() *U32I32Map {
	if ms.Len() <= 0 {
		return ms
	}
	for k, _ := range ms.value {
		ms.deleted[k] = struct{}{}
	}
	ms.value = map[uint32]int32{}
	ms.dirtied = map[uint32]struct{}{}
	ms.SetParentDirty()
	return ms
}
func (ms *U32I32Map) Get(k uint32) (int32, bool) {
	v, exist := ms.value[k]
	return v, exist
}

func (ms *U32I32Map) Remove(k uint32) (int32, bool) {
	v, exist := ms.value[k]
	if !exist {
		return v, false
	}
	delete(ms.value, k)
	delete(ms.dirtied, k)
	ms.deleted[k] = struct{}{}
	ms.SetParentDirty()
	return v, true
}

func (ms *U32I32Map) RemoveAll(k []uint32) {
	for _, kk := range k {
		ms.Remove(kk)
	}
}
func (ms *U32I32Map) Iter() iter.Seq2[uint32, int32] {
	return func(yield func(uint32, int32) bool) {
		for k, v := range ms.value {
			if !yield(k, v) {
				return
			}
		}
	}
}

func (ms *U32I32Map) Dirtied() map[uint32]struct{} {
	return ms.dirtied
}

func (ms *U32I32Map) Deleted() map[uint32]struct{} {
	return ms.deleted
}

func (ms *U32I32Map) ContainDirtied(kk uint32) bool {
	for k := range ms.dirtied {
		if k == kk {
			return true
		}
	}
	return false
}

func (ms *U32I32Map) ContainDeleted(kk uint32) bool {
	for k := range ms.deleted {
		if k == kk {
			return true
		}
	}
	return false
}

func (ms *U32I32Map) MergeDirtyFromPb(delList []uint32, dirties map[uint32]int32) {
	ms.RemoveAll(delList)
	for k, v := range dirties {
		ms.Put(k, v)
	}
}

func (ms *U32I32Map) Entries() map[uint32]int32 {
	return ms.value
}

type U32BoolMap struct {
	parent      Editor
	keyInParent any
	value       map[uint32]bool
	dirtied     map[uint32]struct{}
	deleted     map[uint32]struct{}
}

func NewU32BoolMap() *U32BoolMap {
	return &U32BoolMap{value: make(map[uint32]bool), dirtied: make(map[uint32]struct{}), deleted: make(map[uint32]struct{})}
}

func (ms *U32BoolMap) SetParent(sync Editor, key any) {
	ms.parent = sync
	ms.keyInParent = key
}
func (ms *U32BoolMap) SetParentDirty() {
	if ms.parent != nil {
		ms.parent.SetDirty(ms.keyInParent, true, ms)
	}
}
func (ms *U32BoolMap) SetDirty(key any, dirty bool, editor Editor) {
	if dirty {
		ms.dirtied[key.(uint32)] = struct{}{}
		ms.SetParentDirty()
	} else {
		delete(ms.dirtied, key.(uint32))
	}

}

func (ms *U32BoolMap) FlushDirty(dirty bool) {
	if !dirty {
		ms.dirtied = map[uint32]struct{}{}
		ms.deleted = map[uint32]struct{}{}
	}
}

func (ms *U32BoolMap) Key() any {
	return nil
}
func (ms *U32BoolMap) SetKey(i any) {

}

func (ms *U32BoolMap) Put(k uint32, v bool) *U32BoolMap {
	_, exist := ms.value[k]
	if exist {
		ms.deleted[k] = struct{}{}
	}
	ms.value[k] = v
	ms.dirtied[k] = struct{}{}
	ms.SetParentDirty()
	return ms
}

func (ms *U32BoolMap) PutAll(kv map[uint32]bool) *U32BoolMap {
	for k, v := range kv {
		ms.Put(k, v)
	}
	return ms
}

func (ms *U32BoolMap) Len() int {
	return len(ms.value)
}

func (ms *U32BoolMap) Clear() *U32BoolMap {
	if ms.Len() <= 0 {
		return ms
	}
	for k, _ := range ms.value {
		ms.deleted[k] = struct{}{}
	}
	ms.value = map[uint32]bool{}
	ms.dirtied = map[uint32]struct{}{}
	ms.SetParentDirty()
	return ms
}
func (ms *U32BoolMap) Get(k uint32) (bool, bool) {
	v, exist := ms.value[k]
	return v, exist
}

func (ms *U32BoolMap) Remove(k uint32) (bool, bool) {
	v, exist := ms.value[k]
	if !exist {
		return v, false
	}
	delete(ms.value, k)
	delete(ms.dirtied, k)
	ms.deleted[k] = struct{}{}
	ms.SetParentDirty()
	return v, true
}

func (ms *U32BoolMap) RemoveAll(k []uint32) {
	for _, kk := range k {
		ms.Remove(kk)
	}
}
func (ms *U32BoolMap) Iter() iter.Seq2[uint32, bool] {
	return func(yield func(uint32, bool) bool) {
		for k, v := range ms.value {
			if !yield(k, v) {
				return
			}
		}
	}
}

func (ms *U32BoolMap) Dirtied() map[uint32]struct{} {
	return ms.dirtied
}

func (ms *U32BoolMap) Deleted() map[uint32]struct{} {
	return ms.deleted
}

func (ms *U32BoolMap) ContainDirtied(kk uint32) bool {
	for k := range ms.dirtied {
		if k == kk {
			return true
		}
	}
	return false
}

func (ms *U32BoolMap) ContainDeleted(kk uint32) bool {
	for k := range ms.deleted {
		if k == kk {
			return true
		}
	}
	return false
}

func (ms *U32BoolMap) MergeDirtyFromPb(delList []uint32, dirties map[uint32]bool) {
	ms.RemoveAll(delList)
	for k, v := range dirties {
		ms.Put(k, v)
	}
}

func (ms *U32BoolMap) Entries() map[uint32]bool {
	return ms.value
}

type U32F64Map struct {
	parent      Editor
	keyInParent any
	value       map[uint32]float64
	dirtied     map[uint32]struct{}
	deleted     map[uint32]struct{}
}

func NewU32F64Map() *U32F64Map {
	return &U32F64Map{value: make(map[uint32]float64), dirtied: make(map[uint32]struct{}), deleted: make(map[uint32]struct{})}
}

func (ms *U32F64Map) SetParent(sync Editor, key any) {
	ms.parent = sync
	ms.keyInParent = key
}
func (ms *U32F64Map) SetParentDirty() {
	if ms.parent != nil {
		ms.parent.SetDirty(ms.keyInParent, true, ms)
	}
}
func (ms *U32F64Map) SetDirty(key any, dirty bool, editor Editor) {
	if dirty {
		ms.dirtied[key.(uint32)] = struct{}{}
		ms.SetParentDirty()
	} else {
		delete(ms.dirtied, key.(uint32))
	}

}

func (ms *U32F64Map) FlushDirty(dirty bool) {
	if !dirty {
		ms.dirtied = map[uint32]struct{}{}
		ms.deleted = map[uint32]struct{}{}
	}
}

func (ms *U32F64Map) Key() any {
	return nil
}
func (ms *U32F64Map) SetKey(i any) {

}

func (ms *U32F64Map) Put(k uint32, v float64) *U32F64Map {
	_, exist := ms.value[k]
	if exist {
		ms.deleted[k] = struct{}{}
	}
	ms.value[k] = v
	ms.dirtied[k] = struct{}{}
	ms.SetParentDirty()
	return ms
}

func (ms *U32F64Map) PutAll(kv map[uint32]float64) *U32F64Map {
	for k, v := range kv {
		ms.Put(k, v)
	}
	return ms
}

func (ms *U32F64Map) Len() int {
	return len(ms.value)
}

func (ms *U32F64Map) Clear() *U32F64Map {
	if ms.Len() <= 0 {
		return ms
	}
	for k, _ := range ms.value {
		ms.deleted[k] = struct{}{}
	}
	ms.value = map[uint32]float64{}
	ms.dirtied = map[uint32]struct{}{}
	ms.SetParentDirty()
	return ms
}
func (ms *U32F64Map) Get(k uint32) (float64, bool) {
	v, exist := ms.value[k]
	return v, exist
}

func (ms *U32F64Map) Remove(k uint32) (float64, bool) {
	v, exist := ms.value[k]
	if !exist {
		return v, false
	}
	delete(ms.value, k)
	delete(ms.dirtied, k)
	ms.deleted[k] = struct{}{}
	ms.SetParentDirty()
	return v, true
}

func (ms *U32F64Map) RemoveAll(k []uint32) {
	for _, kk := range k {
		ms.Remove(kk)
	}
}
func (ms *U32F64Map) Iter() iter.Seq2[uint32, float64] {
	return func(yield func(uint32, float64) bool) {
		for k, v := range ms.value {
			if !yield(k, v) {
				return
			}
		}
	}
}

func (ms *U32F64Map) Dirtied() map[uint32]struct{} {
	return ms.dirtied
}

func (ms *U32F64Map) Deleted() map[uint32]struct{} {
	return ms.deleted
}

func (ms *U32F64Map) ContainDirtied(kk uint32) bool {
	for k := range ms.dirtied {
		if k == kk {
			return true
		}
	}
	return false
}

func (ms *U32F64Map) ContainDeleted(kk uint32) bool {
	for k := range ms.deleted {
		if k == kk {
			return true
		}
	}
	return false
}

func (ms *U32F64Map) MergeDirtyFromPb(delList []uint32, dirties map[uint32]float64) {
	ms.RemoveAll(delList)
	for k, v := range dirties {
		ms.Put(k, v)
	}
}

func (ms *U32F64Map) Entries() map[uint32]float64 {
	return ms.value
}

type U32BytesMap struct {
	parent      Editor
	keyInParent any
	value       map[uint32][]byte
	dirtied     map[uint32]struct{}
	deleted     map[uint32]struct{}
}

func NewU32BytesMap() *U32BytesMap {
	return &U32BytesMap{value: make(map[uint32][]byte), dirtied: make(map[uint32]struct{}), deleted: make(map[uint32]struct{})}
}

func (ms *U32BytesMap) SetParent(sync Editor, key any) {
	ms.parent = sync
	ms.keyInParent = key
}
func (ms *U32BytesMap) SetParentDirty() {
	if ms.parent != nil {
		ms.parent.SetDirty(ms.keyInParent, true, ms)
	}
}
func (ms *U32BytesMap) SetDirty(key any, dirty bool, editor Editor) {
	if dirty {
		ms.dirtied[key.(uint32)] = struct{}{}
		ms.SetParentDirty()
	} else {
		delete(ms.dirtied, key.(uint32))
	}

}

func (ms *U32BytesMap) FlushDirty(dirty bool) {
	if !dirty {
		ms.dirtied = map[uint32]struct{}{}
		ms.deleted = map[uint32]struct{}{}
	}
}

func (ms *U32BytesMap) Key() any {
	return nil
}
func (ms *U32BytesMap) SetKey(i any) {

}

func (ms *U32BytesMap) Put(k uint32, v []byte) *U32BytesMap {
	_, exist := ms.value[k]
	if exist {
		ms.deleted[k] = struct{}{}
	}
	ms.value[k] = v
	ms.dirtied[k] = struct{}{}
	ms.SetParentDirty()
	return ms
}

func (ms *U32BytesMap) PutAll(kv map[uint32][]byte) *U32BytesMap {
	for k, v := range kv {
		ms.Put(k, v)
	}
	return ms
}

func (ms *U32BytesMap) Len() int {
	return len(ms.value)
}

func (ms *U32BytesMap) Clear() *U32BytesMap {
	if ms.Len() <= 0 {
		return ms
	}
	for k, _ := range ms.value {
		ms.deleted[k] = struct{}{}
	}
	ms.value = map[uint32][]byte{}
	ms.dirtied = map[uint32]struct{}{}
	ms.SetParentDirty()
	return ms
}
func (ms *U32BytesMap) Get(k uint32) ([]byte, bool) {
	v, exist := ms.value[k]
	return v, exist
}

func (ms *U32BytesMap) Remove(k uint32) ([]byte, bool) {
	v, exist := ms.value[k]
	if !exist {
		return v, false
	}
	delete(ms.value, k)
	delete(ms.dirtied, k)
	ms.deleted[k] = struct{}{}
	ms.SetParentDirty()
	return v, true
}

func (ms *U32BytesMap) RemoveAll(k []uint32) {
	for _, kk := range k {
		ms.Remove(kk)
	}
}
func (ms *U32BytesMap) Iter() iter.Seq2[uint32, []byte] {
	return func(yield func(uint32, []byte) bool) {
		for k, v := range ms.value {
			if !yield(k, v) {
				return
			}
		}
	}
}

func (ms *U32BytesMap) Dirtied() map[uint32]struct{} {
	return ms.dirtied
}

func (ms *U32BytesMap) Deleted() map[uint32]struct{} {
	return ms.deleted
}

func (ms *U32BytesMap) ContainDirtied(kk uint32) bool {
	for k := range ms.dirtied {
		if k == kk {
			return true
		}
	}
	return false
}

func (ms *U32BytesMap) ContainDeleted(kk uint32) bool {
	for k := range ms.deleted {
		if k == kk {
			return true
		}
	}
	return false
}

func (ms *U32BytesMap) MergeDirtyFromPb(delList []uint32, dirties map[uint32][]byte) {
	ms.RemoveAll(delList)
	for k, v := range dirties {
		ms.Put(k, v)
	}
}

func (ms *U32BytesMap) Entries() map[uint32][]byte {
	return ms.value
}

type U64U32Map struct {
	parent      Editor
	keyInParent any
	value       map[uint64]uint32
	dirtied     map[uint64]struct{}
	deleted     map[uint64]struct{}
}

func NewU64U32Map() *U64U32Map {
	return &U64U32Map{value: make(map[uint64]uint32), dirtied: make(map[uint64]struct{}), deleted: make(map[uint64]struct{})}
}

func (ms *U64U32Map) SetParent(sync Editor, key any) {
	ms.parent = sync
	ms.keyInParent = key
}
func (ms *U64U32Map) SetParentDirty() {
	if ms.parent != nil {
		ms.parent.SetDirty(ms.keyInParent, true, ms)
	}
}
func (ms *U64U32Map) SetDirty(key any, dirty bool, editor Editor) {
	if dirty {
		ms.dirtied[key.(uint64)] = struct{}{}
		ms.SetParentDirty()
	} else {
		delete(ms.dirtied, key.(uint64))
	}

}

func (ms *U64U32Map) FlushDirty(dirty bool) {
	if !dirty {
		ms.dirtied = map[uint64]struct{}{}
		ms.deleted = map[uint64]struct{}{}
	}
}

func (ms *U64U32Map) Key() any {
	return nil
}
func (ms *U64U32Map) SetKey(i any) {

}

func (ms *U64U32Map) Put(k uint64, v uint32) *U64U32Map {
	_, exist := ms.value[k]
	if exist {
		ms.deleted[k] = struct{}{}
	}
	ms.value[k] = v
	ms.dirtied[k] = struct{}{}
	ms.SetParentDirty()
	return ms
}

func (ms *U64U32Map) PutAll(kv map[uint64]uint32) *U64U32Map {
	for k, v := range kv {
		ms.Put(k, v)
	}
	return ms
}

func (ms *U64U32Map) Len() int {
	return len(ms.value)
}

func (ms *U64U32Map) Clear() *U64U32Map {
	if ms.Len() <= 0 {
		return ms
	}
	for k, _ := range ms.value {
		ms.deleted[k] = struct{}{}
	}
	ms.value = map[uint64]uint32{}
	ms.dirtied = map[uint64]struct{}{}
	ms.SetParentDirty()
	return ms
}
func (ms *U64U32Map) Get(k uint64) (uint32, bool) {
	v, exist := ms.value[k]
	return v, exist
}

func (ms *U64U32Map) Remove(k uint64) (uint32, bool) {
	v, exist := ms.value[k]
	if !exist {
		return v, false
	}
	delete(ms.value, k)
	delete(ms.dirtied, k)
	ms.deleted[k] = struct{}{}
	ms.SetParentDirty()
	return v, true
}

func (ms *U64U32Map) RemoveAll(k []uint64) {
	for _, kk := range k {
		ms.Remove(kk)
	}
}
func (ms *U64U32Map) Iter() iter.Seq2[uint64, uint32] {
	return func(yield func(uint64, uint32) bool) {
		for k, v := range ms.value {
			if !yield(k, v) {
				return
			}
		}
	}
}

func (ms *U64U32Map) Dirtied() map[uint64]struct{} {
	return ms.dirtied
}

func (ms *U64U32Map) Deleted() map[uint64]struct{} {
	return ms.deleted
}

func (ms *U64U32Map) ContainDirtied(kk uint64) bool {
	for k := range ms.dirtied {
		if k == kk {
			return true
		}
	}
	return false
}

func (ms *U64U32Map) ContainDeleted(kk uint64) bool {
	for k := range ms.deleted {
		if k == kk {
			return true
		}
	}
	return false
}

func (ms *U64U32Map) MergeDirtyFromPb(delList []uint64, dirties map[uint64]uint32) {
	ms.RemoveAll(delList)
	for k, v := range dirties {
		ms.Put(k, v)
	}
}

func (ms *U64U32Map) Entries() map[uint64]uint32 {
	return ms.value
}

type U64I64Map struct {
	parent      Editor
	keyInParent any
	value       map[uint64]int64
	dirtied     map[uint64]struct{}
	deleted     map[uint64]struct{}
}

func NewU64I64Map() *U64I64Map {
	return &U64I64Map{value: make(map[uint64]int64), dirtied: make(map[uint64]struct{}), deleted: make(map[uint64]struct{})}
}

func (ms *U64I64Map) SetParent(sync Editor, key any) {
	ms.parent = sync
	ms.keyInParent = key
}
func (ms *U64I64Map) SetParentDirty() {
	if ms.parent != nil {
		ms.parent.SetDirty(ms.keyInParent, true, ms)
	}
}
func (ms *U64I64Map) SetDirty(key any, dirty bool, editor Editor) {
	if dirty {
		ms.dirtied[key.(uint64)] = struct{}{}
		ms.SetParentDirty()
	} else {
		delete(ms.dirtied, key.(uint64))
	}

}

func (ms *U64I64Map) FlushDirty(dirty bool) {
	if !dirty {
		ms.dirtied = map[uint64]struct{}{}
		ms.deleted = map[uint64]struct{}{}
	}
}

func (ms *U64I64Map) Key() any {
	return nil
}
func (ms *U64I64Map) SetKey(i any) {

}

func (ms *U64I64Map) Put(k uint64, v int64) *U64I64Map {
	_, exist := ms.value[k]
	if exist {
		ms.deleted[k] = struct{}{}
	}
	ms.value[k] = v
	ms.dirtied[k] = struct{}{}
	ms.SetParentDirty()
	return ms
}

func (ms *U64I64Map) PutAll(kv map[uint64]int64) *U64I64Map {
	for k, v := range kv {
		ms.Put(k, v)
	}
	return ms
}

func (ms *U64I64Map) Len() int {
	return len(ms.value)
}

func (ms *U64I64Map) Clear() *U64I64Map {
	if ms.Len() <= 0 {
		return ms
	}
	for k, _ := range ms.value {
		ms.deleted[k] = struct{}{}
	}
	ms.value = map[uint64]int64{}
	ms.dirtied = map[uint64]struct{}{}
	ms.SetParentDirty()
	return ms
}
func (ms *U64I64Map) Get(k uint64) (int64, bool) {
	v, exist := ms.value[k]
	return v, exist
}

func (ms *U64I64Map) Remove(k uint64) (int64, bool) {
	v, exist := ms.value[k]
	if !exist {
		return v, false
	}
	delete(ms.value, k)
	delete(ms.dirtied, k)
	ms.deleted[k] = struct{}{}
	ms.SetParentDirty()
	return v, true
}

func (ms *U64I64Map) RemoveAll(k []uint64) {
	for _, kk := range k {
		ms.Remove(kk)
	}
}
func (ms *U64I64Map) Iter() iter.Seq2[uint64, int64] {
	return func(yield func(uint64, int64) bool) {
		for k, v := range ms.value {
			if !yield(k, v) {
				return
			}
		}
	}
}

func (ms *U64I64Map) Dirtied() map[uint64]struct{} {
	return ms.dirtied
}

func (ms *U64I64Map) Deleted() map[uint64]struct{} {
	return ms.deleted
}

func (ms *U64I64Map) ContainDirtied(kk uint64) bool {
	for k := range ms.dirtied {
		if k == kk {
			return true
		}
	}
	return false
}

func (ms *U64I64Map) ContainDeleted(kk uint64) bool {
	for k := range ms.deleted {
		if k == kk {
			return true
		}
	}
	return false
}

func (ms *U64I64Map) MergeDirtyFromPb(delList []uint64, dirties map[uint64]int64) {
	ms.RemoveAll(delList)
	for k, v := range dirties {
		ms.Put(k, v)
	}
}

func (ms *U64I64Map) Entries() map[uint64]int64 {
	return ms.value
}

type U64StrMap struct {
	parent      Editor
	keyInParent any
	value       map[uint64]string
	dirtied     map[uint64]struct{}
	deleted     map[uint64]struct{}
}

func NewU64StrMap() *U64StrMap {
	return &U64StrMap{value: make(map[uint64]string), dirtied: make(map[uint64]struct{}), deleted: make(map[uint64]struct{})}
}

func (ms *U64StrMap) SetParent(sync Editor, key any) {
	ms.parent = sync
	ms.keyInParent = key
}
func (ms *U64StrMap) SetParentDirty() {
	if ms.parent != nil {
		ms.parent.SetDirty(ms.keyInParent, true, ms)
	}
}
func (ms *U64StrMap) SetDirty(key any, dirty bool, editor Editor) {
	if dirty {
		ms.dirtied[key.(uint64)] = struct{}{}
		ms.SetParentDirty()
	} else {
		delete(ms.dirtied, key.(uint64))
	}

}

func (ms *U64StrMap) FlushDirty(dirty bool) {
	if !dirty {
		ms.dirtied = map[uint64]struct{}{}
		ms.deleted = map[uint64]struct{}{}
	}
}

func (ms *U64StrMap) Key() any {
	return nil
}
func (ms *U64StrMap) SetKey(i any) {

}

func (ms *U64StrMap) Put(k uint64, v string) *U64StrMap {
	_, exist := ms.value[k]
	if exist {
		ms.deleted[k] = struct{}{}
	}
	ms.value[k] = v
	ms.dirtied[k] = struct{}{}
	ms.SetParentDirty()
	return ms
}

func (ms *U64StrMap) PutAll(kv map[uint64]string) *U64StrMap {
	for k, v := range kv {
		ms.Put(k, v)
	}
	return ms
}

func (ms *U64StrMap) Len() int {
	return len(ms.value)
}

func (ms *U64StrMap) Clear() *U64StrMap {
	if ms.Len() <= 0 {
		return ms
	}
	for k, _ := range ms.value {
		ms.deleted[k] = struct{}{}
	}
	ms.value = map[uint64]string{}
	ms.dirtied = map[uint64]struct{}{}
	ms.SetParentDirty()
	return ms
}
func (ms *U64StrMap) Get(k uint64) (string, bool) {
	v, exist := ms.value[k]
	return v, exist
}

func (ms *U64StrMap) Remove(k uint64) (string, bool) {
	v, exist := ms.value[k]
	if !exist {
		return v, false
	}
	delete(ms.value, k)
	delete(ms.dirtied, k)
	ms.deleted[k] = struct{}{}
	ms.SetParentDirty()
	return v, true
}

func (ms *U64StrMap) RemoveAll(k []uint64) {
	for _, kk := range k {
		ms.Remove(kk)
	}
}
func (ms *U64StrMap) Iter() iter.Seq2[uint64, string] {
	return func(yield func(uint64, string) bool) {
		for k, v := range ms.value {
			if !yield(k, v) {
				return
			}
		}
	}
}

func (ms *U64StrMap) Dirtied() map[uint64]struct{} {
	return ms.dirtied
}

func (ms *U64StrMap) Deleted() map[uint64]struct{} {
	return ms.deleted
}

func (ms *U64StrMap) ContainDirtied(kk uint64) bool {
	for k := range ms.dirtied {
		if k == kk {
			return true
		}
	}
	return false
}

func (ms *U64StrMap) ContainDeleted(kk uint64) bool {
	for k := range ms.deleted {
		if k == kk {
			return true
		}
	}
	return false
}

func (ms *U64StrMap) MergeDirtyFromPb(delList []uint64, dirties map[uint64]string) {
	ms.RemoveAll(delList)
	for k, v := range dirties {
		ms.Put(k, v)
	}
}

func (ms *U64StrMap) Entries() map[uint64]string {
	return ms.value
}

type U64F32Map struct {
	parent      Editor
	keyInParent any
	value       map[uint64]float32
	dirtied     map[uint64]struct{}
	deleted     map[uint64]struct{}
}

func NewU64F32Map() *U64F32Map {
	return &U64F32Map{value: make(map[uint64]float32), dirtied: make(map[uint64]struct{}), deleted: make(map[uint64]struct{})}
}

func (ms *U64F32Map) SetParent(sync Editor, key any) {
	ms.parent = sync
	ms.keyInParent = key
}
func (ms *U64F32Map) SetParentDirty() {
	if ms.parent != nil {
		ms.parent.SetDirty(ms.keyInParent, true, ms)
	}
}
func (ms *U64F32Map) SetDirty(key any, dirty bool, editor Editor) {
	if dirty {
		ms.dirtied[key.(uint64)] = struct{}{}
		ms.SetParentDirty()
	} else {
		delete(ms.dirtied, key.(uint64))
	}

}

func (ms *U64F32Map) FlushDirty(dirty bool) {
	if !dirty {
		ms.dirtied = map[uint64]struct{}{}
		ms.deleted = map[uint64]struct{}{}
	}
}

func (ms *U64F32Map) Key() any {
	return nil
}
func (ms *U64F32Map) SetKey(i any) {

}

func (ms *U64F32Map) Put(k uint64, v float32) *U64F32Map {
	_, exist := ms.value[k]
	if exist {
		ms.deleted[k] = struct{}{}
	}
	ms.value[k] = v
	ms.dirtied[k] = struct{}{}
	ms.SetParentDirty()
	return ms
}

func (ms *U64F32Map) PutAll(kv map[uint64]float32) *U64F32Map {
	for k, v := range kv {
		ms.Put(k, v)
	}
	return ms
}

func (ms *U64F32Map) Len() int {
	return len(ms.value)
}

func (ms *U64F32Map) Clear() *U64F32Map {
	if ms.Len() <= 0 {
		return ms
	}
	for k, _ := range ms.value {
		ms.deleted[k] = struct{}{}
	}
	ms.value = map[uint64]float32{}
	ms.dirtied = map[uint64]struct{}{}
	ms.SetParentDirty()
	return ms
}
func (ms *U64F32Map) Get(k uint64) (float32, bool) {
	v, exist := ms.value[k]
	return v, exist
}

func (ms *U64F32Map) Remove(k uint64) (float32, bool) {
	v, exist := ms.value[k]
	if !exist {
		return v, false
	}
	delete(ms.value, k)
	delete(ms.dirtied, k)
	ms.deleted[k] = struct{}{}
	ms.SetParentDirty()
	return v, true
}

func (ms *U64F32Map) RemoveAll(k []uint64) {
	for _, kk := range k {
		ms.Remove(kk)
	}
}
func (ms *U64F32Map) Iter() iter.Seq2[uint64, float32] {
	return func(yield func(uint64, float32) bool) {
		for k, v := range ms.value {
			if !yield(k, v) {
				return
			}
		}
	}
}

func (ms *U64F32Map) Dirtied() map[uint64]struct{} {
	return ms.dirtied
}

func (ms *U64F32Map) Deleted() map[uint64]struct{} {
	return ms.deleted
}

func (ms *U64F32Map) ContainDirtied(kk uint64) bool {
	for k := range ms.dirtied {
		if k == kk {
			return true
		}
	}
	return false
}

func (ms *U64F32Map) ContainDeleted(kk uint64) bool {
	for k := range ms.deleted {
		if k == kk {
			return true
		}
	}
	return false
}

func (ms *U64F32Map) MergeDirtyFromPb(delList []uint64, dirties map[uint64]float32) {
	ms.RemoveAll(delList)
	for k, v := range dirties {
		ms.Put(k, v)
	}
}

func (ms *U64F32Map) Entries() map[uint64]float32 {
	return ms.value
}

type U64U64Map struct {
	parent      Editor
	keyInParent any
	value       map[uint64]uint64
	dirtied     map[uint64]struct{}
	deleted     map[uint64]struct{}
}

func NewU64U64Map() *U64U64Map {
	return &U64U64Map{value: make(map[uint64]uint64), dirtied: make(map[uint64]struct{}), deleted: make(map[uint64]struct{})}
}

func (ms *U64U64Map) SetParent(sync Editor, key any) {
	ms.parent = sync
	ms.keyInParent = key
}
func (ms *U64U64Map) SetParentDirty() {
	if ms.parent != nil {
		ms.parent.SetDirty(ms.keyInParent, true, ms)
	}
}
func (ms *U64U64Map) SetDirty(key any, dirty bool, editor Editor) {
	if dirty {
		ms.dirtied[key.(uint64)] = struct{}{}
		ms.SetParentDirty()
	} else {
		delete(ms.dirtied, key.(uint64))
	}

}

func (ms *U64U64Map) FlushDirty(dirty bool) {
	if !dirty {
		ms.dirtied = map[uint64]struct{}{}
		ms.deleted = map[uint64]struct{}{}
	}
}

func (ms *U64U64Map) Key() any {
	return nil
}
func (ms *U64U64Map) SetKey(i any) {

}

func (ms *U64U64Map) Put(k uint64, v uint64) *U64U64Map {
	_, exist := ms.value[k]
	if exist {
		ms.deleted[k] = struct{}{}
	}
	ms.value[k] = v
	ms.dirtied[k] = struct{}{}
	ms.SetParentDirty()
	return ms
}

func (ms *U64U64Map) PutAll(kv map[uint64]uint64) *U64U64Map {
	for k, v := range kv {
		ms.Put(k, v)
	}
	return ms
}

func (ms *U64U64Map) Len() int {
	return len(ms.value)
}

func (ms *U64U64Map) Clear() *U64U64Map {
	if ms.Len() <= 0 {
		return ms
	}
	for k, _ := range ms.value {
		ms.deleted[k] = struct{}{}
	}
	ms.value = map[uint64]uint64{}
	ms.dirtied = map[uint64]struct{}{}
	ms.SetParentDirty()
	return ms
}
func (ms *U64U64Map) Get(k uint64) (uint64, bool) {
	v, exist := ms.value[k]
	return v, exist
}

func (ms *U64U64Map) Remove(k uint64) (uint64, bool) {
	v, exist := ms.value[k]
	if !exist {
		return v, false
	}
	delete(ms.value, k)
	delete(ms.dirtied, k)
	ms.deleted[k] = struct{}{}
	ms.SetParentDirty()
	return v, true
}

func (ms *U64U64Map) RemoveAll(k []uint64) {
	for _, kk := range k {
		ms.Remove(kk)
	}
}
func (ms *U64U64Map) Iter() iter.Seq2[uint64, uint64] {
	return func(yield func(uint64, uint64) bool) {
		for k, v := range ms.value {
			if !yield(k, v) {
				return
			}
		}
	}
}

func (ms *U64U64Map) Dirtied() map[uint64]struct{} {
	return ms.dirtied
}

func (ms *U64U64Map) Deleted() map[uint64]struct{} {
	return ms.deleted
}

func (ms *U64U64Map) ContainDirtied(kk uint64) bool {
	for k := range ms.dirtied {
		if k == kk {
			return true
		}
	}
	return false
}

func (ms *U64U64Map) ContainDeleted(kk uint64) bool {
	for k := range ms.deleted {
		if k == kk {
			return true
		}
	}
	return false
}

func (ms *U64U64Map) MergeDirtyFromPb(delList []uint64, dirties map[uint64]uint64) {
	ms.RemoveAll(delList)
	for k, v := range dirties {
		ms.Put(k, v)
	}
}

func (ms *U64U64Map) Entries() map[uint64]uint64 {
	return ms.value
}

type U64I32Map struct {
	parent      Editor
	keyInParent any
	value       map[uint64]int32
	dirtied     map[uint64]struct{}
	deleted     map[uint64]struct{}
}

func NewU64I32Map() *U64I32Map {
	return &U64I32Map{value: make(map[uint64]int32), dirtied: make(map[uint64]struct{}), deleted: make(map[uint64]struct{})}
}

func (ms *U64I32Map) SetParent(sync Editor, key any) {
	ms.parent = sync
	ms.keyInParent = key
}
func (ms *U64I32Map) SetParentDirty() {
	if ms.parent != nil {
		ms.parent.SetDirty(ms.keyInParent, true, ms)
	}
}
func (ms *U64I32Map) SetDirty(key any, dirty bool, editor Editor) {
	if dirty {
		ms.dirtied[key.(uint64)] = struct{}{}
		ms.SetParentDirty()
	} else {
		delete(ms.dirtied, key.(uint64))
	}

}

func (ms *U64I32Map) FlushDirty(dirty bool) {
	if !dirty {
		ms.dirtied = map[uint64]struct{}{}
		ms.deleted = map[uint64]struct{}{}
	}
}

func (ms *U64I32Map) Key() any {
	return nil
}
func (ms *U64I32Map) SetKey(i any) {

}

func (ms *U64I32Map) Put(k uint64, v int32) *U64I32Map {
	_, exist := ms.value[k]
	if exist {
		ms.deleted[k] = struct{}{}
	}
	ms.value[k] = v
	ms.dirtied[k] = struct{}{}
	ms.SetParentDirty()
	return ms
}

func (ms *U64I32Map) PutAll(kv map[uint64]int32) *U64I32Map {
	for k, v := range kv {
		ms.Put(k, v)
	}
	return ms
}

func (ms *U64I32Map) Len() int {
	return len(ms.value)
}

func (ms *U64I32Map) Clear() *U64I32Map {
	if ms.Len() <= 0 {
		return ms
	}
	for k, _ := range ms.value {
		ms.deleted[k] = struct{}{}
	}
	ms.value = map[uint64]int32{}
	ms.dirtied = map[uint64]struct{}{}
	ms.SetParentDirty()
	return ms
}
func (ms *U64I32Map) Get(k uint64) (int32, bool) {
	v, exist := ms.value[k]
	return v, exist
}

func (ms *U64I32Map) Remove(k uint64) (int32, bool) {
	v, exist := ms.value[k]
	if !exist {
		return v, false
	}
	delete(ms.value, k)
	delete(ms.dirtied, k)
	ms.deleted[k] = struct{}{}
	ms.SetParentDirty()
	return v, true
}

func (ms *U64I32Map) RemoveAll(k []uint64) {
	for _, kk := range k {
		ms.Remove(kk)
	}
}
func (ms *U64I32Map) Iter() iter.Seq2[uint64, int32] {
	return func(yield func(uint64, int32) bool) {
		for k, v := range ms.value {
			if !yield(k, v) {
				return
			}
		}
	}
}

func (ms *U64I32Map) Dirtied() map[uint64]struct{} {
	return ms.dirtied
}

func (ms *U64I32Map) Deleted() map[uint64]struct{} {
	return ms.deleted
}

func (ms *U64I32Map) ContainDirtied(kk uint64) bool {
	for k := range ms.dirtied {
		if k == kk {
			return true
		}
	}
	return false
}

func (ms *U64I32Map) ContainDeleted(kk uint64) bool {
	for k := range ms.deleted {
		if k == kk {
			return true
		}
	}
	return false
}

func (ms *U64I32Map) MergeDirtyFromPb(delList []uint64, dirties map[uint64]int32) {
	ms.RemoveAll(delList)
	for k, v := range dirties {
		ms.Put(k, v)
	}
}

func (ms *U64I32Map) Entries() map[uint64]int32 {
	return ms.value
}

type U64BoolMap struct {
	parent      Editor
	keyInParent any
	value       map[uint64]bool
	dirtied     map[uint64]struct{}
	deleted     map[uint64]struct{}
}

func NewU64BoolMap() *U64BoolMap {
	return &U64BoolMap{value: make(map[uint64]bool), dirtied: make(map[uint64]struct{}), deleted: make(map[uint64]struct{})}
}

func (ms *U64BoolMap) SetParent(sync Editor, key any) {
	ms.parent = sync
	ms.keyInParent = key
}
func (ms *U64BoolMap) SetParentDirty() {
	if ms.parent != nil {
		ms.parent.SetDirty(ms.keyInParent, true, ms)
	}
}
func (ms *U64BoolMap) SetDirty(key any, dirty bool, editor Editor) {
	if dirty {
		ms.dirtied[key.(uint64)] = struct{}{}
		ms.SetParentDirty()
	} else {
		delete(ms.dirtied, key.(uint64))
	}

}

func (ms *U64BoolMap) FlushDirty(dirty bool) {
	if !dirty {
		ms.dirtied = map[uint64]struct{}{}
		ms.deleted = map[uint64]struct{}{}
	}
}

func (ms *U64BoolMap) Key() any {
	return nil
}
func (ms *U64BoolMap) SetKey(i any) {

}

func (ms *U64BoolMap) Put(k uint64, v bool) *U64BoolMap {
	_, exist := ms.value[k]
	if exist {
		ms.deleted[k] = struct{}{}
	}
	ms.value[k] = v
	ms.dirtied[k] = struct{}{}
	ms.SetParentDirty()
	return ms
}

func (ms *U64BoolMap) PutAll(kv map[uint64]bool) *U64BoolMap {
	for k, v := range kv {
		ms.Put(k, v)
	}
	return ms
}

func (ms *U64BoolMap) Len() int {
	return len(ms.value)
}

func (ms *U64BoolMap) Clear() *U64BoolMap {
	if ms.Len() <= 0 {
		return ms
	}
	for k, _ := range ms.value {
		ms.deleted[k] = struct{}{}
	}
	ms.value = map[uint64]bool{}
	ms.dirtied = map[uint64]struct{}{}
	ms.SetParentDirty()
	return ms
}
func (ms *U64BoolMap) Get(k uint64) (bool, bool) {
	v, exist := ms.value[k]
	return v, exist
}

func (ms *U64BoolMap) Remove(k uint64) (bool, bool) {
	v, exist := ms.value[k]
	if !exist {
		return v, false
	}
	delete(ms.value, k)
	delete(ms.dirtied, k)
	ms.deleted[k] = struct{}{}
	ms.SetParentDirty()
	return v, true
}

func (ms *U64BoolMap) RemoveAll(k []uint64) {
	for _, kk := range k {
		ms.Remove(kk)
	}
}
func (ms *U64BoolMap) Iter() iter.Seq2[uint64, bool] {
	return func(yield func(uint64, bool) bool) {
		for k, v := range ms.value {
			if !yield(k, v) {
				return
			}
		}
	}
}

func (ms *U64BoolMap) Dirtied() map[uint64]struct{} {
	return ms.dirtied
}

func (ms *U64BoolMap) Deleted() map[uint64]struct{} {
	return ms.deleted
}

func (ms *U64BoolMap) ContainDirtied(kk uint64) bool {
	for k := range ms.dirtied {
		if k == kk {
			return true
		}
	}
	return false
}

func (ms *U64BoolMap) ContainDeleted(kk uint64) bool {
	for k := range ms.deleted {
		if k == kk {
			return true
		}
	}
	return false
}

func (ms *U64BoolMap) MergeDirtyFromPb(delList []uint64, dirties map[uint64]bool) {
	ms.RemoveAll(delList)
	for k, v := range dirties {
		ms.Put(k, v)
	}
}

func (ms *U64BoolMap) Entries() map[uint64]bool {
	return ms.value
}

type U64F64Map struct {
	parent      Editor
	keyInParent any
	value       map[uint64]float64
	dirtied     map[uint64]struct{}
	deleted     map[uint64]struct{}
}

func NewU64F64Map() *U64F64Map {
	return &U64F64Map{value: make(map[uint64]float64), dirtied: make(map[uint64]struct{}), deleted: make(map[uint64]struct{})}
}

func (ms *U64F64Map) SetParent(sync Editor, key any) {
	ms.parent = sync
	ms.keyInParent = key
}
func (ms *U64F64Map) SetParentDirty() {
	if ms.parent != nil {
		ms.parent.SetDirty(ms.keyInParent, true, ms)
	}
}
func (ms *U64F64Map) SetDirty(key any, dirty bool, editor Editor) {
	if dirty {
		ms.dirtied[key.(uint64)] = struct{}{}
		ms.SetParentDirty()
	} else {
		delete(ms.dirtied, key.(uint64))
	}

}

func (ms *U64F64Map) FlushDirty(dirty bool) {
	if !dirty {
		ms.dirtied = map[uint64]struct{}{}
		ms.deleted = map[uint64]struct{}{}
	}
}

func (ms *U64F64Map) Key() any {
	return nil
}
func (ms *U64F64Map) SetKey(i any) {

}

func (ms *U64F64Map) Put(k uint64, v float64) *U64F64Map {
	_, exist := ms.value[k]
	if exist {
		ms.deleted[k] = struct{}{}
	}
	ms.value[k] = v
	ms.dirtied[k] = struct{}{}
	ms.SetParentDirty()
	return ms
}

func (ms *U64F64Map) PutAll(kv map[uint64]float64) *U64F64Map {
	for k, v := range kv {
		ms.Put(k, v)
	}
	return ms
}

func (ms *U64F64Map) Len() int {
	return len(ms.value)
}

func (ms *U64F64Map) Clear() *U64F64Map {
	if ms.Len() <= 0 {
		return ms
	}
	for k, _ := range ms.value {
		ms.deleted[k] = struct{}{}
	}
	ms.value = map[uint64]float64{}
	ms.dirtied = map[uint64]struct{}{}
	ms.SetParentDirty()
	return ms
}
func (ms *U64F64Map) Get(k uint64) (float64, bool) {
	v, exist := ms.value[k]
	return v, exist
}

func (ms *U64F64Map) Remove(k uint64) (float64, bool) {
	v, exist := ms.value[k]
	if !exist {
		return v, false
	}
	delete(ms.value, k)
	delete(ms.dirtied, k)
	ms.deleted[k] = struct{}{}
	ms.SetParentDirty()
	return v, true
}

func (ms *U64F64Map) RemoveAll(k []uint64) {
	for _, kk := range k {
		ms.Remove(kk)
	}
}
func (ms *U64F64Map) Iter() iter.Seq2[uint64, float64] {
	return func(yield func(uint64, float64) bool) {
		for k, v := range ms.value {
			if !yield(k, v) {
				return
			}
		}
	}
}

func (ms *U64F64Map) Dirtied() map[uint64]struct{} {
	return ms.dirtied
}

func (ms *U64F64Map) Deleted() map[uint64]struct{} {
	return ms.deleted
}

func (ms *U64F64Map) ContainDirtied(kk uint64) bool {
	for k := range ms.dirtied {
		if k == kk {
			return true
		}
	}
	return false
}

func (ms *U64F64Map) ContainDeleted(kk uint64) bool {
	for k := range ms.deleted {
		if k == kk {
			return true
		}
	}
	return false
}

func (ms *U64F64Map) MergeDirtyFromPb(delList []uint64, dirties map[uint64]float64) {
	ms.RemoveAll(delList)
	for k, v := range dirties {
		ms.Put(k, v)
	}
}

func (ms *U64F64Map) Entries() map[uint64]float64 {
	return ms.value
}

type U64BytesMap struct {
	parent      Editor
	keyInParent any
	value       map[uint64][]byte
	dirtied     map[uint64]struct{}
	deleted     map[uint64]struct{}
}

func NewU64BytesMap() *U64BytesMap {
	return &U64BytesMap{value: make(map[uint64][]byte), dirtied: make(map[uint64]struct{}), deleted: make(map[uint64]struct{})}
}

func (ms *U64BytesMap) SetParent(sync Editor, key any) {
	ms.parent = sync
	ms.keyInParent = key
}
func (ms *U64BytesMap) SetParentDirty() {
	if ms.parent != nil {
		ms.parent.SetDirty(ms.keyInParent, true, ms)
	}
}
func (ms *U64BytesMap) SetDirty(key any, dirty bool, editor Editor) {
	if dirty {
		ms.dirtied[key.(uint64)] = struct{}{}
		ms.SetParentDirty()
	} else {
		delete(ms.dirtied, key.(uint64))
	}

}

func (ms *U64BytesMap) FlushDirty(dirty bool) {
	if !dirty {
		ms.dirtied = map[uint64]struct{}{}
		ms.deleted = map[uint64]struct{}{}
	}
}

func (ms *U64BytesMap) Key() any {
	return nil
}
func (ms *U64BytesMap) SetKey(i any) {

}

func (ms *U64BytesMap) Put(k uint64, v []byte) *U64BytesMap {
	_, exist := ms.value[k]
	if exist {
		ms.deleted[k] = struct{}{}
	}
	ms.value[k] = v
	ms.dirtied[k] = struct{}{}
	ms.SetParentDirty()
	return ms
}

func (ms *U64BytesMap) PutAll(kv map[uint64][]byte) *U64BytesMap {
	for k, v := range kv {
		ms.Put(k, v)
	}
	return ms
}

func (ms *U64BytesMap) Len() int {
	return len(ms.value)
}

func (ms *U64BytesMap) Clear() *U64BytesMap {
	if ms.Len() <= 0 {
		return ms
	}
	for k, _ := range ms.value {
		ms.deleted[k] = struct{}{}
	}
	ms.value = map[uint64][]byte{}
	ms.dirtied = map[uint64]struct{}{}
	ms.SetParentDirty()
	return ms
}
func (ms *U64BytesMap) Get(k uint64) ([]byte, bool) {
	v, exist := ms.value[k]
	return v, exist
}

func (ms *U64BytesMap) Remove(k uint64) ([]byte, bool) {
	v, exist := ms.value[k]
	if !exist {
		return v, false
	}
	delete(ms.value, k)
	delete(ms.dirtied, k)
	ms.deleted[k] = struct{}{}
	ms.SetParentDirty()
	return v, true
}

func (ms *U64BytesMap) RemoveAll(k []uint64) {
	for _, kk := range k {
		ms.Remove(kk)
	}
}
func (ms *U64BytesMap) Iter() iter.Seq2[uint64, []byte] {
	return func(yield func(uint64, []byte) bool) {
		for k, v := range ms.value {
			if !yield(k, v) {
				return
			}
		}
	}
}

func (ms *U64BytesMap) Dirtied() map[uint64]struct{} {
	return ms.dirtied
}

func (ms *U64BytesMap) Deleted() map[uint64]struct{} {
	return ms.deleted
}

func (ms *U64BytesMap) ContainDirtied(kk uint64) bool {
	for k := range ms.dirtied {
		if k == kk {
			return true
		}
	}
	return false
}

func (ms *U64BytesMap) ContainDeleted(kk uint64) bool {
	for k := range ms.deleted {
		if k == kk {
			return true
		}
	}
	return false
}

func (ms *U64BytesMap) MergeDirtyFromPb(delList []uint64, dirties map[uint64][]byte) {
	ms.RemoveAll(delList)
	for k, v := range dirties {
		ms.Put(k, v)
	}
}

func (ms *U64BytesMap) Entries() map[uint64][]byte {
	return ms.value
}

type I32U32Map struct {
	parent      Editor
	keyInParent any
	value       map[int32]uint32
	dirtied     map[int32]struct{}
	deleted     map[int32]struct{}
}

func NewI32U32Map() *I32U32Map {
	return &I32U32Map{value: make(map[int32]uint32), dirtied: make(map[int32]struct{}), deleted: make(map[int32]struct{})}
}

func (ms *I32U32Map) SetParent(sync Editor, key any) {
	ms.parent = sync
	ms.keyInParent = key
}
func (ms *I32U32Map) SetParentDirty() {
	if ms.parent != nil {
		ms.parent.SetDirty(ms.keyInParent, true, ms)
	}
}
func (ms *I32U32Map) SetDirty(key any, dirty bool, editor Editor) {
	if dirty {
		ms.dirtied[key.(int32)] = struct{}{}
		ms.SetParentDirty()
	} else {
		delete(ms.dirtied, key.(int32))
	}

}

func (ms *I32U32Map) FlushDirty(dirty bool) {
	if !dirty {
		ms.dirtied = map[int32]struct{}{}
		ms.deleted = map[int32]struct{}{}
	}
}

func (ms *I32U32Map) Key() any {
	return nil
}
func (ms *I32U32Map) SetKey(i any) {

}

func (ms *I32U32Map) Put(k int32, v uint32) *I32U32Map {
	_, exist := ms.value[k]
	if exist {
		ms.deleted[k] = struct{}{}
	}
	ms.value[k] = v
	ms.dirtied[k] = struct{}{}
	ms.SetParentDirty()
	return ms
}

func (ms *I32U32Map) PutAll(kv map[int32]uint32) *I32U32Map {
	for k, v := range kv {
		ms.Put(k, v)
	}
	return ms
}

func (ms *I32U32Map) Len() int {
	return len(ms.value)
}

func (ms *I32U32Map) Clear() *I32U32Map {
	if ms.Len() <= 0 {
		return ms
	}
	for k, _ := range ms.value {
		ms.deleted[k] = struct{}{}
	}
	ms.value = map[int32]uint32{}
	ms.dirtied = map[int32]struct{}{}
	ms.SetParentDirty()
	return ms
}
func (ms *I32U32Map) Get(k int32) (uint32, bool) {
	v, exist := ms.value[k]
	return v, exist
}

func (ms *I32U32Map) Remove(k int32) (uint32, bool) {
	v, exist := ms.value[k]
	if !exist {
		return v, false
	}
	delete(ms.value, k)
	delete(ms.dirtied, k)
	ms.deleted[k] = struct{}{}
	ms.SetParentDirty()
	return v, true
}

func (ms *I32U32Map) RemoveAll(k []int32) {
	for _, kk := range k {
		ms.Remove(kk)
	}
}
func (ms *I32U32Map) Iter() iter.Seq2[int32, uint32] {
	return func(yield func(int32, uint32) bool) {
		for k, v := range ms.value {
			if !yield(k, v) {
				return
			}
		}
	}
}

func (ms *I32U32Map) Dirtied() map[int32]struct{} {
	return ms.dirtied
}

func (ms *I32U32Map) Deleted() map[int32]struct{} {
	return ms.deleted
}

func (ms *I32U32Map) ContainDirtied(kk int32) bool {
	for k := range ms.dirtied {
		if k == kk {
			return true
		}
	}
	return false
}

func (ms *I32U32Map) ContainDeleted(kk int32) bool {
	for k := range ms.deleted {
		if k == kk {
			return true
		}
	}
	return false
}

func (ms *I32U32Map) MergeDirtyFromPb(delList []int32, dirties map[int32]uint32) {
	ms.RemoveAll(delList)
	for k, v := range dirties {
		ms.Put(k, v)
	}
}

func (ms *I32U32Map) Entries() map[int32]uint32 {
	return ms.value
}

type I32I64Map struct {
	parent      Editor
	keyInParent any
	value       map[int32]int64
	dirtied     map[int32]struct{}
	deleted     map[int32]struct{}
}

func NewI32I64Map() *I32I64Map {
	return &I32I64Map{value: make(map[int32]int64), dirtied: make(map[int32]struct{}), deleted: make(map[int32]struct{})}
}

func (ms *I32I64Map) SetParent(sync Editor, key any) {
	ms.parent = sync
	ms.keyInParent = key
}
func (ms *I32I64Map) SetParentDirty() {
	if ms.parent != nil {
		ms.parent.SetDirty(ms.keyInParent, true, ms)
	}
}
func (ms *I32I64Map) SetDirty(key any, dirty bool, editor Editor) {
	if dirty {
		ms.dirtied[key.(int32)] = struct{}{}
		ms.SetParentDirty()
	} else {
		delete(ms.dirtied, key.(int32))
	}

}

func (ms *I32I64Map) FlushDirty(dirty bool) {
	if !dirty {
		ms.dirtied = map[int32]struct{}{}
		ms.deleted = map[int32]struct{}{}
	}
}

func (ms *I32I64Map) Key() any {
	return nil
}
func (ms *I32I64Map) SetKey(i any) {

}

func (ms *I32I64Map) Put(k int32, v int64) *I32I64Map {
	_, exist := ms.value[k]
	if exist {
		ms.deleted[k] = struct{}{}
	}
	ms.value[k] = v
	ms.dirtied[k] = struct{}{}
	ms.SetParentDirty()
	return ms
}

func (ms *I32I64Map) PutAll(kv map[int32]int64) *I32I64Map {
	for k, v := range kv {
		ms.Put(k, v)
	}
	return ms
}

func (ms *I32I64Map) Len() int {
	return len(ms.value)
}

func (ms *I32I64Map) Clear() *I32I64Map {
	if ms.Len() <= 0 {
		return ms
	}
	for k, _ := range ms.value {
		ms.deleted[k] = struct{}{}
	}
	ms.value = map[int32]int64{}
	ms.dirtied = map[int32]struct{}{}
	ms.SetParentDirty()
	return ms
}
func (ms *I32I64Map) Get(k int32) (int64, bool) {
	v, exist := ms.value[k]
	return v, exist
}

func (ms *I32I64Map) Remove(k int32) (int64, bool) {
	v, exist := ms.value[k]
	if !exist {
		return v, false
	}
	delete(ms.value, k)
	delete(ms.dirtied, k)
	ms.deleted[k] = struct{}{}
	ms.SetParentDirty()
	return v, true
}

func (ms *I32I64Map) RemoveAll(k []int32) {
	for _, kk := range k {
		ms.Remove(kk)
	}
}
func (ms *I32I64Map) Iter() iter.Seq2[int32, int64] {
	return func(yield func(int32, int64) bool) {
		for k, v := range ms.value {
			if !yield(k, v) {
				return
			}
		}
	}
}

func (ms *I32I64Map) Dirtied() map[int32]struct{} {
	return ms.dirtied
}

func (ms *I32I64Map) Deleted() map[int32]struct{} {
	return ms.deleted
}

func (ms *I32I64Map) ContainDirtied(kk int32) bool {
	for k := range ms.dirtied {
		if k == kk {
			return true
		}
	}
	return false
}

func (ms *I32I64Map) ContainDeleted(kk int32) bool {
	for k := range ms.deleted {
		if k == kk {
			return true
		}
	}
	return false
}

func (ms *I32I64Map) MergeDirtyFromPb(delList []int32, dirties map[int32]int64) {
	ms.RemoveAll(delList)
	for k, v := range dirties {
		ms.Put(k, v)
	}
}

func (ms *I32I64Map) Entries() map[int32]int64 {
	return ms.value
}

type I32StrMap struct {
	parent      Editor
	keyInParent any
	value       map[int32]string
	dirtied     map[int32]struct{}
	deleted     map[int32]struct{}
}

func NewI32StrMap() *I32StrMap {
	return &I32StrMap{value: make(map[int32]string), dirtied: make(map[int32]struct{}), deleted: make(map[int32]struct{})}
}

func (ms *I32StrMap) SetParent(sync Editor, key any) {
	ms.parent = sync
	ms.keyInParent = key
}
func (ms *I32StrMap) SetParentDirty() {
	if ms.parent != nil {
		ms.parent.SetDirty(ms.keyInParent, true, ms)
	}
}
func (ms *I32StrMap) SetDirty(key any, dirty bool, editor Editor) {
	if dirty {
		ms.dirtied[key.(int32)] = struct{}{}
		ms.SetParentDirty()
	} else {
		delete(ms.dirtied, key.(int32))
	}

}

func (ms *I32StrMap) FlushDirty(dirty bool) {
	if !dirty {
		ms.dirtied = map[int32]struct{}{}
		ms.deleted = map[int32]struct{}{}
	}
}

func (ms *I32StrMap) Key() any {
	return nil
}
func (ms *I32StrMap) SetKey(i any) {

}

func (ms *I32StrMap) Put(k int32, v string) *I32StrMap {
	_, exist := ms.value[k]
	if exist {
		ms.deleted[k] = struct{}{}
	}
	ms.value[k] = v
	ms.dirtied[k] = struct{}{}
	ms.SetParentDirty()
	return ms
}

func (ms *I32StrMap) PutAll(kv map[int32]string) *I32StrMap {
	for k, v := range kv {
		ms.Put(k, v)
	}
	return ms
}

func (ms *I32StrMap) Len() int {
	return len(ms.value)
}

func (ms *I32StrMap) Clear() *I32StrMap {
	if ms.Len() <= 0 {
		return ms
	}
	for k, _ := range ms.value {
		ms.deleted[k] = struct{}{}
	}
	ms.value = map[int32]string{}
	ms.dirtied = map[int32]struct{}{}
	ms.SetParentDirty()
	return ms
}
func (ms *I32StrMap) Get(k int32) (string, bool) {
	v, exist := ms.value[k]
	return v, exist
}

func (ms *I32StrMap) Remove(k int32) (string, bool) {
	v, exist := ms.value[k]
	if !exist {
		return v, false
	}
	delete(ms.value, k)
	delete(ms.dirtied, k)
	ms.deleted[k] = struct{}{}
	ms.SetParentDirty()
	return v, true
}

func (ms *I32StrMap) RemoveAll(k []int32) {
	for _, kk := range k {
		ms.Remove(kk)
	}
}
func (ms *I32StrMap) Iter() iter.Seq2[int32, string] {
	return func(yield func(int32, string) bool) {
		for k, v := range ms.value {
			if !yield(k, v) {
				return
			}
		}
	}
}

func (ms *I32StrMap) Dirtied() map[int32]struct{} {
	return ms.dirtied
}

func (ms *I32StrMap) Deleted() map[int32]struct{} {
	return ms.deleted
}

func (ms *I32StrMap) ContainDirtied(kk int32) bool {
	for k := range ms.dirtied {
		if k == kk {
			return true
		}
	}
	return false
}

func (ms *I32StrMap) ContainDeleted(kk int32) bool {
	for k := range ms.deleted {
		if k == kk {
			return true
		}
	}
	return false
}

func (ms *I32StrMap) MergeDirtyFromPb(delList []int32, dirties map[int32]string) {
	ms.RemoveAll(delList)
	for k, v := range dirties {
		ms.Put(k, v)
	}
}

func (ms *I32StrMap) Entries() map[int32]string {
	return ms.value
}

type I32F32Map struct {
	parent      Editor
	keyInParent any
	value       map[int32]float32
	dirtied     map[int32]struct{}
	deleted     map[int32]struct{}
}

func NewI32F32Map() *I32F32Map {
	return &I32F32Map{value: make(map[int32]float32), dirtied: make(map[int32]struct{}), deleted: make(map[int32]struct{})}
}

func (ms *I32F32Map) SetParent(sync Editor, key any) {
	ms.parent = sync
	ms.keyInParent = key
}
func (ms *I32F32Map) SetParentDirty() {
	if ms.parent != nil {
		ms.parent.SetDirty(ms.keyInParent, true, ms)
	}
}
func (ms *I32F32Map) SetDirty(key any, dirty bool, editor Editor) {
	if dirty {
		ms.dirtied[key.(int32)] = struct{}{}
		ms.SetParentDirty()
	} else {
		delete(ms.dirtied, key.(int32))
	}

}

func (ms *I32F32Map) FlushDirty(dirty bool) {
	if !dirty {
		ms.dirtied = map[int32]struct{}{}
		ms.deleted = map[int32]struct{}{}
	}
}

func (ms *I32F32Map) Key() any {
	return nil
}
func (ms *I32F32Map) SetKey(i any) {

}

func (ms *I32F32Map) Put(k int32, v float32) *I32F32Map {
	_, exist := ms.value[k]
	if exist {
		ms.deleted[k] = struct{}{}
	}
	ms.value[k] = v
	ms.dirtied[k] = struct{}{}
	ms.SetParentDirty()
	return ms
}

func (ms *I32F32Map) PutAll(kv map[int32]float32) *I32F32Map {
	for k, v := range kv {
		ms.Put(k, v)
	}
	return ms
}

func (ms *I32F32Map) Len() int {
	return len(ms.value)
}

func (ms *I32F32Map) Clear() *I32F32Map {
	if ms.Len() <= 0 {
		return ms
	}
	for k, _ := range ms.value {
		ms.deleted[k] = struct{}{}
	}
	ms.value = map[int32]float32{}
	ms.dirtied = map[int32]struct{}{}
	ms.SetParentDirty()
	return ms
}
func (ms *I32F32Map) Get(k int32) (float32, bool) {
	v, exist := ms.value[k]
	return v, exist
}

func (ms *I32F32Map) Remove(k int32) (float32, bool) {
	v, exist := ms.value[k]
	if !exist {
		return v, false
	}
	delete(ms.value, k)
	delete(ms.dirtied, k)
	ms.deleted[k] = struct{}{}
	ms.SetParentDirty()
	return v, true
}

func (ms *I32F32Map) RemoveAll(k []int32) {
	for _, kk := range k {
		ms.Remove(kk)
	}
}
func (ms *I32F32Map) Iter() iter.Seq2[int32, float32] {
	return func(yield func(int32, float32) bool) {
		for k, v := range ms.value {
			if !yield(k, v) {
				return
			}
		}
	}
}

func (ms *I32F32Map) Dirtied() map[int32]struct{} {
	return ms.dirtied
}

func (ms *I32F32Map) Deleted() map[int32]struct{} {
	return ms.deleted
}

func (ms *I32F32Map) ContainDirtied(kk int32) bool {
	for k := range ms.dirtied {
		if k == kk {
			return true
		}
	}
	return false
}

func (ms *I32F32Map) ContainDeleted(kk int32) bool {
	for k := range ms.deleted {
		if k == kk {
			return true
		}
	}
	return false
}

func (ms *I32F32Map) MergeDirtyFromPb(delList []int32, dirties map[int32]float32) {
	ms.RemoveAll(delList)
	for k, v := range dirties {
		ms.Put(k, v)
	}
}

func (ms *I32F32Map) Entries() map[int32]float32 {
	return ms.value
}

type I32U64Map struct {
	parent      Editor
	keyInParent any
	value       map[int32]uint64
	dirtied     map[int32]struct{}
	deleted     map[int32]struct{}
}

func NewI32U64Map() *I32U64Map {
	return &I32U64Map{value: make(map[int32]uint64), dirtied: make(map[int32]struct{}), deleted: make(map[int32]struct{})}
}

func (ms *I32U64Map) SetParent(sync Editor, key any) {
	ms.parent = sync
	ms.keyInParent = key
}
func (ms *I32U64Map) SetParentDirty() {
	if ms.parent != nil {
		ms.parent.SetDirty(ms.keyInParent, true, ms)
	}
}
func (ms *I32U64Map) SetDirty(key any, dirty bool, editor Editor) {
	if dirty {
		ms.dirtied[key.(int32)] = struct{}{}
		ms.SetParentDirty()
	} else {
		delete(ms.dirtied, key.(int32))
	}

}

func (ms *I32U64Map) FlushDirty(dirty bool) {
	if !dirty {
		ms.dirtied = map[int32]struct{}{}
		ms.deleted = map[int32]struct{}{}
	}
}

func (ms *I32U64Map) Key() any {
	return nil
}
func (ms *I32U64Map) SetKey(i any) {

}

func (ms *I32U64Map) Put(k int32, v uint64) *I32U64Map {
	_, exist := ms.value[k]
	if exist {
		ms.deleted[k] = struct{}{}
	}
	ms.value[k] = v
	ms.dirtied[k] = struct{}{}
	ms.SetParentDirty()
	return ms
}

func (ms *I32U64Map) PutAll(kv map[int32]uint64) *I32U64Map {
	for k, v := range kv {
		ms.Put(k, v)
	}
	return ms
}

func (ms *I32U64Map) Len() int {
	return len(ms.value)
}

func (ms *I32U64Map) Clear() *I32U64Map {
	if ms.Len() <= 0 {
		return ms
	}
	for k, _ := range ms.value {
		ms.deleted[k] = struct{}{}
	}
	ms.value = map[int32]uint64{}
	ms.dirtied = map[int32]struct{}{}
	ms.SetParentDirty()
	return ms
}
func (ms *I32U64Map) Get(k int32) (uint64, bool) {
	v, exist := ms.value[k]
	return v, exist
}

func (ms *I32U64Map) Remove(k int32) (uint64, bool) {
	v, exist := ms.value[k]
	if !exist {
		return v, false
	}
	delete(ms.value, k)
	delete(ms.dirtied, k)
	ms.deleted[k] = struct{}{}
	ms.SetParentDirty()
	return v, true
}

func (ms *I32U64Map) RemoveAll(k []int32) {
	for _, kk := range k {
		ms.Remove(kk)
	}
}
func (ms *I32U64Map) Iter() iter.Seq2[int32, uint64] {
	return func(yield func(int32, uint64) bool) {
		for k, v := range ms.value {
			if !yield(k, v) {
				return
			}
		}
	}
}

func (ms *I32U64Map) Dirtied() map[int32]struct{} {
	return ms.dirtied
}

func (ms *I32U64Map) Deleted() map[int32]struct{} {
	return ms.deleted
}

func (ms *I32U64Map) ContainDirtied(kk int32) bool {
	for k := range ms.dirtied {
		if k == kk {
			return true
		}
	}
	return false
}

func (ms *I32U64Map) ContainDeleted(kk int32) bool {
	for k := range ms.deleted {
		if k == kk {
			return true
		}
	}
	return false
}

func (ms *I32U64Map) MergeDirtyFromPb(delList []int32, dirties map[int32]uint64) {
	ms.RemoveAll(delList)
	for k, v := range dirties {
		ms.Put(k, v)
	}
}

func (ms *I32U64Map) Entries() map[int32]uint64 {
	return ms.value
}

type I32I32Map struct {
	parent      Editor
	keyInParent any
	value       map[int32]int32
	dirtied     map[int32]struct{}
	deleted     map[int32]struct{}
}

func NewI32I32Map() *I32I32Map {
	return &I32I32Map{value: make(map[int32]int32), dirtied: make(map[int32]struct{}), deleted: make(map[int32]struct{})}
}

func (ms *I32I32Map) SetParent(sync Editor, key any) {
	ms.parent = sync
	ms.keyInParent = key
}
func (ms *I32I32Map) SetParentDirty() {
	if ms.parent != nil {
		ms.parent.SetDirty(ms.keyInParent, true, ms)
	}
}
func (ms *I32I32Map) SetDirty(key any, dirty bool, editor Editor) {
	if dirty {
		ms.dirtied[key.(int32)] = struct{}{}
		ms.SetParentDirty()
	} else {
		delete(ms.dirtied, key.(int32))
	}

}

func (ms *I32I32Map) FlushDirty(dirty bool) {
	if !dirty {
		ms.dirtied = map[int32]struct{}{}
		ms.deleted = map[int32]struct{}{}
	}
}

func (ms *I32I32Map) Key() any {
	return nil
}
func (ms *I32I32Map) SetKey(i any) {

}

func (ms *I32I32Map) Put(k int32, v int32) *I32I32Map {
	_, exist := ms.value[k]
	if exist {
		ms.deleted[k] = struct{}{}
	}
	ms.value[k] = v
	ms.dirtied[k] = struct{}{}
	ms.SetParentDirty()
	return ms
}

func (ms *I32I32Map) PutAll(kv map[int32]int32) *I32I32Map {
	for k, v := range kv {
		ms.Put(k, v)
	}
	return ms
}

func (ms *I32I32Map) Len() int {
	return len(ms.value)
}

func (ms *I32I32Map) Clear() *I32I32Map {
	if ms.Len() <= 0 {
		return ms
	}
	for k, _ := range ms.value {
		ms.deleted[k] = struct{}{}
	}
	ms.value = map[int32]int32{}
	ms.dirtied = map[int32]struct{}{}
	ms.SetParentDirty()
	return ms
}
func (ms *I32I32Map) Get(k int32) (int32, bool) {
	v, exist := ms.value[k]
	return v, exist
}

func (ms *I32I32Map) Remove(k int32) (int32, bool) {
	v, exist := ms.value[k]
	if !exist {
		return v, false
	}
	delete(ms.value, k)
	delete(ms.dirtied, k)
	ms.deleted[k] = struct{}{}
	ms.SetParentDirty()
	return v, true
}

func (ms *I32I32Map) RemoveAll(k []int32) {
	for _, kk := range k {
		ms.Remove(kk)
	}
}
func (ms *I32I32Map) Iter() iter.Seq2[int32, int32] {
	return func(yield func(int32, int32) bool) {
		for k, v := range ms.value {
			if !yield(k, v) {
				return
			}
		}
	}
}

func (ms *I32I32Map) Dirtied() map[int32]struct{} {
	return ms.dirtied
}

func (ms *I32I32Map) Deleted() map[int32]struct{} {
	return ms.deleted
}

func (ms *I32I32Map) ContainDirtied(kk int32) bool {
	for k := range ms.dirtied {
		if k == kk {
			return true
		}
	}
	return false
}

func (ms *I32I32Map) ContainDeleted(kk int32) bool {
	for k := range ms.deleted {
		if k == kk {
			return true
		}
	}
	return false
}

func (ms *I32I32Map) MergeDirtyFromPb(delList []int32, dirties map[int32]int32) {
	ms.RemoveAll(delList)
	for k, v := range dirties {
		ms.Put(k, v)
	}
}

func (ms *I32I32Map) Entries() map[int32]int32 {
	return ms.value
}

type I32BoolMap struct {
	parent      Editor
	keyInParent any
	value       map[int32]bool
	dirtied     map[int32]struct{}
	deleted     map[int32]struct{}
}

func NewI32BoolMap() *I32BoolMap {
	return &I32BoolMap{value: make(map[int32]bool), dirtied: make(map[int32]struct{}), deleted: make(map[int32]struct{})}
}

func (ms *I32BoolMap) SetParent(sync Editor, key any) {
	ms.parent = sync
	ms.keyInParent = key
}
func (ms *I32BoolMap) SetParentDirty() {
	if ms.parent != nil {
		ms.parent.SetDirty(ms.keyInParent, true, ms)
	}
}
func (ms *I32BoolMap) SetDirty(key any, dirty bool, editor Editor) {
	if dirty {
		ms.dirtied[key.(int32)] = struct{}{}
		ms.SetParentDirty()
	} else {
		delete(ms.dirtied, key.(int32))
	}

}

func (ms *I32BoolMap) FlushDirty(dirty bool) {
	if !dirty {
		ms.dirtied = map[int32]struct{}{}
		ms.deleted = map[int32]struct{}{}
	}
}

func (ms *I32BoolMap) Key() any {
	return nil
}
func (ms *I32BoolMap) SetKey(i any) {

}

func (ms *I32BoolMap) Put(k int32, v bool) *I32BoolMap {
	_, exist := ms.value[k]
	if exist {
		ms.deleted[k] = struct{}{}
	}
	ms.value[k] = v
	ms.dirtied[k] = struct{}{}
	ms.SetParentDirty()
	return ms
}

func (ms *I32BoolMap) PutAll(kv map[int32]bool) *I32BoolMap {
	for k, v := range kv {
		ms.Put(k, v)
	}
	return ms
}

func (ms *I32BoolMap) Len() int {
	return len(ms.value)
}

func (ms *I32BoolMap) Clear() *I32BoolMap {
	if ms.Len() <= 0 {
		return ms
	}
	for k, _ := range ms.value {
		ms.deleted[k] = struct{}{}
	}
	ms.value = map[int32]bool{}
	ms.dirtied = map[int32]struct{}{}
	ms.SetParentDirty()
	return ms
}
func (ms *I32BoolMap) Get(k int32) (bool, bool) {
	v, exist := ms.value[k]
	return v, exist
}

func (ms *I32BoolMap) Remove(k int32) (bool, bool) {
	v, exist := ms.value[k]
	if !exist {
		return v, false
	}
	delete(ms.value, k)
	delete(ms.dirtied, k)
	ms.deleted[k] = struct{}{}
	ms.SetParentDirty()
	return v, true
}

func (ms *I32BoolMap) RemoveAll(k []int32) {
	for _, kk := range k {
		ms.Remove(kk)
	}
}
func (ms *I32BoolMap) Iter() iter.Seq2[int32, bool] {
	return func(yield func(int32, bool) bool) {
		for k, v := range ms.value {
			if !yield(k, v) {
				return
			}
		}
	}
}

func (ms *I32BoolMap) Dirtied() map[int32]struct{} {
	return ms.dirtied
}

func (ms *I32BoolMap) Deleted() map[int32]struct{} {
	return ms.deleted
}

func (ms *I32BoolMap) ContainDirtied(kk int32) bool {
	for k := range ms.dirtied {
		if k == kk {
			return true
		}
	}
	return false
}

func (ms *I32BoolMap) ContainDeleted(kk int32) bool {
	for k := range ms.deleted {
		if k == kk {
			return true
		}
	}
	return false
}

func (ms *I32BoolMap) MergeDirtyFromPb(delList []int32, dirties map[int32]bool) {
	ms.RemoveAll(delList)
	for k, v := range dirties {
		ms.Put(k, v)
	}
}

func (ms *I32BoolMap) Entries() map[int32]bool {
	return ms.value
}

type I32F64Map struct {
	parent      Editor
	keyInParent any
	value       map[int32]float64
	dirtied     map[int32]struct{}
	deleted     map[int32]struct{}
}

func NewI32F64Map() *I32F64Map {
	return &I32F64Map{value: make(map[int32]float64), dirtied: make(map[int32]struct{}), deleted: make(map[int32]struct{})}
}

func (ms *I32F64Map) SetParent(sync Editor, key any) {
	ms.parent = sync
	ms.keyInParent = key
}
func (ms *I32F64Map) SetParentDirty() {
	if ms.parent != nil {
		ms.parent.SetDirty(ms.keyInParent, true, ms)
	}
}
func (ms *I32F64Map) SetDirty(key any, dirty bool, editor Editor) {
	if dirty {
		ms.dirtied[key.(int32)] = struct{}{}
		ms.SetParentDirty()
	} else {
		delete(ms.dirtied, key.(int32))
	}

}

func (ms *I32F64Map) FlushDirty(dirty bool) {
	if !dirty {
		ms.dirtied = map[int32]struct{}{}
		ms.deleted = map[int32]struct{}{}
	}
}

func (ms *I32F64Map) Key() any {
	return nil
}
func (ms *I32F64Map) SetKey(i any) {

}

func (ms *I32F64Map) Put(k int32, v float64) *I32F64Map {
	_, exist := ms.value[k]
	if exist {
		ms.deleted[k] = struct{}{}
	}
	ms.value[k] = v
	ms.dirtied[k] = struct{}{}
	ms.SetParentDirty()
	return ms
}

func (ms *I32F64Map) PutAll(kv map[int32]float64) *I32F64Map {
	for k, v := range kv {
		ms.Put(k, v)
	}
	return ms
}

func (ms *I32F64Map) Len() int {
	return len(ms.value)
}

func (ms *I32F64Map) Clear() *I32F64Map {
	if ms.Len() <= 0 {
		return ms
	}
	for k, _ := range ms.value {
		ms.deleted[k] = struct{}{}
	}
	ms.value = map[int32]float64{}
	ms.dirtied = map[int32]struct{}{}
	ms.SetParentDirty()
	return ms
}
func (ms *I32F64Map) Get(k int32) (float64, bool) {
	v, exist := ms.value[k]
	return v, exist
}

func (ms *I32F64Map) Remove(k int32) (float64, bool) {
	v, exist := ms.value[k]
	if !exist {
		return v, false
	}
	delete(ms.value, k)
	delete(ms.dirtied, k)
	ms.deleted[k] = struct{}{}
	ms.SetParentDirty()
	return v, true
}

func (ms *I32F64Map) RemoveAll(k []int32) {
	for _, kk := range k {
		ms.Remove(kk)
	}
}
func (ms *I32F64Map) Iter() iter.Seq2[int32, float64] {
	return func(yield func(int32, float64) bool) {
		for k, v := range ms.value {
			if !yield(k, v) {
				return
			}
		}
	}
}

func (ms *I32F64Map) Dirtied() map[int32]struct{} {
	return ms.dirtied
}

func (ms *I32F64Map) Deleted() map[int32]struct{} {
	return ms.deleted
}

func (ms *I32F64Map) ContainDirtied(kk int32) bool {
	for k := range ms.dirtied {
		if k == kk {
			return true
		}
	}
	return false
}

func (ms *I32F64Map) ContainDeleted(kk int32) bool {
	for k := range ms.deleted {
		if k == kk {
			return true
		}
	}
	return false
}

func (ms *I32F64Map) MergeDirtyFromPb(delList []int32, dirties map[int32]float64) {
	ms.RemoveAll(delList)
	for k, v := range dirties {
		ms.Put(k, v)
	}
}

func (ms *I32F64Map) Entries() map[int32]float64 {
	return ms.value
}

type I32BytesMap struct {
	parent      Editor
	keyInParent any
	value       map[int32][]byte
	dirtied     map[int32]struct{}
	deleted     map[int32]struct{}
}

func NewI32BytesMap() *I32BytesMap {
	return &I32BytesMap{value: make(map[int32][]byte), dirtied: make(map[int32]struct{}), deleted: make(map[int32]struct{})}
}

func (ms *I32BytesMap) SetParent(sync Editor, key any) {
	ms.parent = sync
	ms.keyInParent = key
}
func (ms *I32BytesMap) SetParentDirty() {
	if ms.parent != nil {
		ms.parent.SetDirty(ms.keyInParent, true, ms)
	}
}
func (ms *I32BytesMap) SetDirty(key any, dirty bool, editor Editor) {
	if dirty {
		ms.dirtied[key.(int32)] = struct{}{}
		ms.SetParentDirty()
	} else {
		delete(ms.dirtied, key.(int32))
	}

}

func (ms *I32BytesMap) FlushDirty(dirty bool) {
	if !dirty {
		ms.dirtied = map[int32]struct{}{}
		ms.deleted = map[int32]struct{}{}
	}
}

func (ms *I32BytesMap) Key() any {
	return nil
}
func (ms *I32BytesMap) SetKey(i any) {

}

func (ms *I32BytesMap) Put(k int32, v []byte) *I32BytesMap {
	_, exist := ms.value[k]
	if exist {
		ms.deleted[k] = struct{}{}
	}
	ms.value[k] = v
	ms.dirtied[k] = struct{}{}
	ms.SetParentDirty()
	return ms
}

func (ms *I32BytesMap) PutAll(kv map[int32][]byte) *I32BytesMap {
	for k, v := range kv {
		ms.Put(k, v)
	}
	return ms
}

func (ms *I32BytesMap) Len() int {
	return len(ms.value)
}

func (ms *I32BytesMap) Clear() *I32BytesMap {
	if ms.Len() <= 0 {
		return ms
	}
	for k, _ := range ms.value {
		ms.deleted[k] = struct{}{}
	}
	ms.value = map[int32][]byte{}
	ms.dirtied = map[int32]struct{}{}
	ms.SetParentDirty()
	return ms
}
func (ms *I32BytesMap) Get(k int32) ([]byte, bool) {
	v, exist := ms.value[k]
	return v, exist
}

func (ms *I32BytesMap) Remove(k int32) ([]byte, bool) {
	v, exist := ms.value[k]
	if !exist {
		return v, false
	}
	delete(ms.value, k)
	delete(ms.dirtied, k)
	ms.deleted[k] = struct{}{}
	ms.SetParentDirty()
	return v, true
}

func (ms *I32BytesMap) RemoveAll(k []int32) {
	for _, kk := range k {
		ms.Remove(kk)
	}
}
func (ms *I32BytesMap) Iter() iter.Seq2[int32, []byte] {
	return func(yield func(int32, []byte) bool) {
		for k, v := range ms.value {
			if !yield(k, v) {
				return
			}
		}
	}
}

func (ms *I32BytesMap) Dirtied() map[int32]struct{} {
	return ms.dirtied
}

func (ms *I32BytesMap) Deleted() map[int32]struct{} {
	return ms.deleted
}

func (ms *I32BytesMap) ContainDirtied(kk int32) bool {
	for k := range ms.dirtied {
		if k == kk {
			return true
		}
	}
	return false
}

func (ms *I32BytesMap) ContainDeleted(kk int32) bool {
	for k := range ms.deleted {
		if k == kk {
			return true
		}
	}
	return false
}

func (ms *I32BytesMap) MergeDirtyFromPb(delList []int32, dirties map[int32][]byte) {
	ms.RemoveAll(delList)
	for k, v := range dirties {
		ms.Put(k, v)
	}
}

func (ms *I32BytesMap) Entries() map[int32][]byte {
	return ms.value
}

type I64BoolMap struct {
	parent      Editor
	keyInParent any
	value       map[int64]bool
	dirtied     map[int64]struct{}
	deleted     map[int64]struct{}
}

func NewI64BoolMap() *I64BoolMap {
	return &I64BoolMap{value: make(map[int64]bool), dirtied: make(map[int64]struct{}), deleted: make(map[int64]struct{})}
}

func (ms *I64BoolMap) SetParent(sync Editor, key any) {
	ms.parent = sync
	ms.keyInParent = key
}
func (ms *I64BoolMap) SetParentDirty() {
	if ms.parent != nil {
		ms.parent.SetDirty(ms.keyInParent, true, ms)
	}
}
func (ms *I64BoolMap) SetDirty(key any, dirty bool, editor Editor) {
	if dirty {
		ms.dirtied[key.(int64)] = struct{}{}
		ms.SetParentDirty()
	} else {
		delete(ms.dirtied, key.(int64))
	}

}

func (ms *I64BoolMap) FlushDirty(dirty bool) {
	if !dirty {
		ms.dirtied = map[int64]struct{}{}
		ms.deleted = map[int64]struct{}{}
	}
}

func (ms *I64BoolMap) Key() any {
	return nil
}
func (ms *I64BoolMap) SetKey(i any) {

}

func (ms *I64BoolMap) Put(k int64, v bool) *I64BoolMap {
	_, exist := ms.value[k]
	if exist {
		ms.deleted[k] = struct{}{}
	}
	ms.value[k] = v
	ms.dirtied[k] = struct{}{}
	ms.SetParentDirty()
	return ms
}

func (ms *I64BoolMap) PutAll(kv map[int64]bool) *I64BoolMap {
	for k, v := range kv {
		ms.Put(k, v)
	}
	return ms
}

func (ms *I64BoolMap) Len() int {
	return len(ms.value)
}

func (ms *I64BoolMap) Clear() *I64BoolMap {
	if ms.Len() <= 0 {
		return ms
	}
	for k, _ := range ms.value {
		ms.deleted[k] = struct{}{}
	}
	ms.value = map[int64]bool{}
	ms.dirtied = map[int64]struct{}{}
	ms.SetParentDirty()
	return ms
}
func (ms *I64BoolMap) Get(k int64) (bool, bool) {
	v, exist := ms.value[k]
	return v, exist
}

func (ms *I64BoolMap) Remove(k int64) (bool, bool) {
	v, exist := ms.value[k]
	if !exist {
		return v, false
	}
	delete(ms.value, k)
	delete(ms.dirtied, k)
	ms.deleted[k] = struct{}{}
	ms.SetParentDirty()
	return v, true
}

func (ms *I64BoolMap) RemoveAll(k []int64) {
	for _, kk := range k {
		ms.Remove(kk)
	}
}
func (ms *I64BoolMap) Iter() iter.Seq2[int64, bool] {
	return func(yield func(int64, bool) bool) {
		for k, v := range ms.value {
			if !yield(k, v) {
				return
			}
		}
	}
}

func (ms *I64BoolMap) Dirtied() map[int64]struct{} {
	return ms.dirtied
}

func (ms *I64BoolMap) Deleted() map[int64]struct{} {
	return ms.deleted
}

func (ms *I64BoolMap) ContainDirtied(kk int64) bool {
	for k := range ms.dirtied {
		if k == kk {
			return true
		}
	}
	return false
}

func (ms *I64BoolMap) ContainDeleted(kk int64) bool {
	for k := range ms.deleted {
		if k == kk {
			return true
		}
	}
	return false
}

func (ms *I64BoolMap) MergeDirtyFromPb(delList []int64, dirties map[int64]bool) {
	ms.RemoveAll(delList)
	for k, v := range dirties {
		ms.Put(k, v)
	}
}

func (ms *I64BoolMap) Entries() map[int64]bool {
	return ms.value
}

type I64F64Map struct {
	parent      Editor
	keyInParent any
	value       map[int64]float64
	dirtied     map[int64]struct{}
	deleted     map[int64]struct{}
}

func NewI64F64Map() *I64F64Map {
	return &I64F64Map{value: make(map[int64]float64), dirtied: make(map[int64]struct{}), deleted: make(map[int64]struct{})}
}

func (ms *I64F64Map) SetParent(sync Editor, key any) {
	ms.parent = sync
	ms.keyInParent = key
}
func (ms *I64F64Map) SetParentDirty() {
	if ms.parent != nil {
		ms.parent.SetDirty(ms.keyInParent, true, ms)
	}
}
func (ms *I64F64Map) SetDirty(key any, dirty bool, editor Editor) {
	if dirty {
		ms.dirtied[key.(int64)] = struct{}{}
		ms.SetParentDirty()
	} else {
		delete(ms.dirtied, key.(int64))
	}

}

func (ms *I64F64Map) FlushDirty(dirty bool) {
	if !dirty {
		ms.dirtied = map[int64]struct{}{}
		ms.deleted = map[int64]struct{}{}
	}
}

func (ms *I64F64Map) Key() any {
	return nil
}
func (ms *I64F64Map) SetKey(i any) {

}

func (ms *I64F64Map) Put(k int64, v float64) *I64F64Map {
	_, exist := ms.value[k]
	if exist {
		ms.deleted[k] = struct{}{}
	}
	ms.value[k] = v
	ms.dirtied[k] = struct{}{}
	ms.SetParentDirty()
	return ms
}

func (ms *I64F64Map) PutAll(kv map[int64]float64) *I64F64Map {
	for k, v := range kv {
		ms.Put(k, v)
	}
	return ms
}

func (ms *I64F64Map) Len() int {
	return len(ms.value)
}

func (ms *I64F64Map) Clear() *I64F64Map {
	if ms.Len() <= 0 {
		return ms
	}
	for k, _ := range ms.value {
		ms.deleted[k] = struct{}{}
	}
	ms.value = map[int64]float64{}
	ms.dirtied = map[int64]struct{}{}
	ms.SetParentDirty()
	return ms
}
func (ms *I64F64Map) Get(k int64) (float64, bool) {
	v, exist := ms.value[k]
	return v, exist
}

func (ms *I64F64Map) Remove(k int64) (float64, bool) {
	v, exist := ms.value[k]
	if !exist {
		return v, false
	}
	delete(ms.value, k)
	delete(ms.dirtied, k)
	ms.deleted[k] = struct{}{}
	ms.SetParentDirty()
	return v, true
}

func (ms *I64F64Map) RemoveAll(k []int64) {
	for _, kk := range k {
		ms.Remove(kk)
	}
}
func (ms *I64F64Map) Iter() iter.Seq2[int64, float64] {
	return func(yield func(int64, float64) bool) {
		for k, v := range ms.value {
			if !yield(k, v) {
				return
			}
		}
	}
}

func (ms *I64F64Map) Dirtied() map[int64]struct{} {
	return ms.dirtied
}

func (ms *I64F64Map) Deleted() map[int64]struct{} {
	return ms.deleted
}

func (ms *I64F64Map) ContainDirtied(kk int64) bool {
	for k := range ms.dirtied {
		if k == kk {
			return true
		}
	}
	return false
}

func (ms *I64F64Map) ContainDeleted(kk int64) bool {
	for k := range ms.deleted {
		if k == kk {
			return true
		}
	}
	return false
}

func (ms *I64F64Map) MergeDirtyFromPb(delList []int64, dirties map[int64]float64) {
	ms.RemoveAll(delList)
	for k, v := range dirties {
		ms.Put(k, v)
	}
}

func (ms *I64F64Map) Entries() map[int64]float64 {
	return ms.value
}

type I64BytesMap struct {
	parent      Editor
	keyInParent any
	value       map[int64][]byte
	dirtied     map[int64]struct{}
	deleted     map[int64]struct{}
}

func NewI64BytesMap() *I64BytesMap {
	return &I64BytesMap{value: make(map[int64][]byte), dirtied: make(map[int64]struct{}), deleted: make(map[int64]struct{})}
}

func (ms *I64BytesMap) SetParent(sync Editor, key any) {
	ms.parent = sync
	ms.keyInParent = key
}
func (ms *I64BytesMap) SetParentDirty() {
	if ms.parent != nil {
		ms.parent.SetDirty(ms.keyInParent, true, ms)
	}
}
func (ms *I64BytesMap) SetDirty(key any, dirty bool, editor Editor) {
	if dirty {
		ms.dirtied[key.(int64)] = struct{}{}
		ms.SetParentDirty()
	} else {
		delete(ms.dirtied, key.(int64))
	}

}

func (ms *I64BytesMap) FlushDirty(dirty bool) {
	if !dirty {
		ms.dirtied = map[int64]struct{}{}
		ms.deleted = map[int64]struct{}{}
	}
}

func (ms *I64BytesMap) Key() any {
	return nil
}
func (ms *I64BytesMap) SetKey(i any) {

}

func (ms *I64BytesMap) Put(k int64, v []byte) *I64BytesMap {
	_, exist := ms.value[k]
	if exist {
		ms.deleted[k] = struct{}{}
	}
	ms.value[k] = v
	ms.dirtied[k] = struct{}{}
	ms.SetParentDirty()
	return ms
}

func (ms *I64BytesMap) PutAll(kv map[int64][]byte) *I64BytesMap {
	for k, v := range kv {
		ms.Put(k, v)
	}
	return ms
}

func (ms *I64BytesMap) Len() int {
	return len(ms.value)
}

func (ms *I64BytesMap) Clear() *I64BytesMap {
	if ms.Len() <= 0 {
		return ms
	}
	for k, _ := range ms.value {
		ms.deleted[k] = struct{}{}
	}
	ms.value = map[int64][]byte{}
	ms.dirtied = map[int64]struct{}{}
	ms.SetParentDirty()
	return ms
}
func (ms *I64BytesMap) Get(k int64) ([]byte, bool) {
	v, exist := ms.value[k]
	return v, exist
}

func (ms *I64BytesMap) Remove(k int64) ([]byte, bool) {
	v, exist := ms.value[k]
	if !exist {
		return v, false
	}
	delete(ms.value, k)
	delete(ms.dirtied, k)
	ms.deleted[k] = struct{}{}
	ms.SetParentDirty()
	return v, true
}

func (ms *I64BytesMap) RemoveAll(k []int64) {
	for _, kk := range k {
		ms.Remove(kk)
	}
}
func (ms *I64BytesMap) Iter() iter.Seq2[int64, []byte] {
	return func(yield func(int64, []byte) bool) {
		for k, v := range ms.value {
			if !yield(k, v) {
				return
			}
		}
	}
}

func (ms *I64BytesMap) Dirtied() map[int64]struct{} {
	return ms.dirtied
}

func (ms *I64BytesMap) Deleted() map[int64]struct{} {
	return ms.deleted
}

func (ms *I64BytesMap) ContainDirtied(kk int64) bool {
	for k := range ms.dirtied {
		if k == kk {
			return true
		}
	}
	return false
}

func (ms *I64BytesMap) ContainDeleted(kk int64) bool {
	for k := range ms.deleted {
		if k == kk {
			return true
		}
	}
	return false
}

func (ms *I64BytesMap) MergeDirtyFromPb(delList []int64, dirties map[int64][]byte) {
	ms.RemoveAll(delList)
	for k, v := range dirties {
		ms.Put(k, v)
	}
}

func (ms *I64BytesMap) Entries() map[int64][]byte {
	return ms.value
}

type I64U64Map struct {
	parent      Editor
	keyInParent any
	value       map[int64]uint64
	dirtied     map[int64]struct{}
	deleted     map[int64]struct{}
}

func NewI64U64Map() *I64U64Map {
	return &I64U64Map{value: make(map[int64]uint64), dirtied: make(map[int64]struct{}), deleted: make(map[int64]struct{})}
}

func (ms *I64U64Map) SetParent(sync Editor, key any) {
	ms.parent = sync
	ms.keyInParent = key
}
func (ms *I64U64Map) SetParentDirty() {
	if ms.parent != nil {
		ms.parent.SetDirty(ms.keyInParent, true, ms)
	}
}
func (ms *I64U64Map) SetDirty(key any, dirty bool, editor Editor) {
	if dirty {
		ms.dirtied[key.(int64)] = struct{}{}
		ms.SetParentDirty()
	} else {
		delete(ms.dirtied, key.(int64))
	}

}

func (ms *I64U64Map) FlushDirty(dirty bool) {
	if !dirty {
		ms.dirtied = map[int64]struct{}{}
		ms.deleted = map[int64]struct{}{}
	}
}

func (ms *I64U64Map) Key() any {
	return nil
}
func (ms *I64U64Map) SetKey(i any) {

}

func (ms *I64U64Map) Put(k int64, v uint64) *I64U64Map {
	_, exist := ms.value[k]
	if exist {
		ms.deleted[k] = struct{}{}
	}
	ms.value[k] = v
	ms.dirtied[k] = struct{}{}
	ms.SetParentDirty()
	return ms
}

func (ms *I64U64Map) PutAll(kv map[int64]uint64) *I64U64Map {
	for k, v := range kv {
		ms.Put(k, v)
	}
	return ms
}

func (ms *I64U64Map) Len() int {
	return len(ms.value)
}

func (ms *I64U64Map) Clear() *I64U64Map {
	if ms.Len() <= 0 {
		return ms
	}
	for k, _ := range ms.value {
		ms.deleted[k] = struct{}{}
	}
	ms.value = map[int64]uint64{}
	ms.dirtied = map[int64]struct{}{}
	ms.SetParentDirty()
	return ms
}
func (ms *I64U64Map) Get(k int64) (uint64, bool) {
	v, exist := ms.value[k]
	return v, exist
}

func (ms *I64U64Map) Remove(k int64) (uint64, bool) {
	v, exist := ms.value[k]
	if !exist {
		return v, false
	}
	delete(ms.value, k)
	delete(ms.dirtied, k)
	ms.deleted[k] = struct{}{}
	ms.SetParentDirty()
	return v, true
}

func (ms *I64U64Map) RemoveAll(k []int64) {
	for _, kk := range k {
		ms.Remove(kk)
	}
}
func (ms *I64U64Map) Iter() iter.Seq2[int64, uint64] {
	return func(yield func(int64, uint64) bool) {
		for k, v := range ms.value {
			if !yield(k, v) {
				return
			}
		}
	}
}

func (ms *I64U64Map) Dirtied() map[int64]struct{} {
	return ms.dirtied
}

func (ms *I64U64Map) Deleted() map[int64]struct{} {
	return ms.deleted
}

func (ms *I64U64Map) ContainDirtied(kk int64) bool {
	for k := range ms.dirtied {
		if k == kk {
			return true
		}
	}
	return false
}

func (ms *I64U64Map) ContainDeleted(kk int64) bool {
	for k := range ms.deleted {
		if k == kk {
			return true
		}
	}
	return false
}

func (ms *I64U64Map) MergeDirtyFromPb(delList []int64, dirties map[int64]uint64) {
	ms.RemoveAll(delList)
	for k, v := range dirties {
		ms.Put(k, v)
	}
}

func (ms *I64U64Map) Entries() map[int64]uint64 {
	return ms.value
}

type I64I32Map struct {
	parent      Editor
	keyInParent any
	value       map[int64]int32
	dirtied     map[int64]struct{}
	deleted     map[int64]struct{}
}

func NewI64I32Map() *I64I32Map {
	return &I64I32Map{value: make(map[int64]int32), dirtied: make(map[int64]struct{}), deleted: make(map[int64]struct{})}
}

func (ms *I64I32Map) SetParent(sync Editor, key any) {
	ms.parent = sync
	ms.keyInParent = key
}
func (ms *I64I32Map) SetParentDirty() {
	if ms.parent != nil {
		ms.parent.SetDirty(ms.keyInParent, true, ms)
	}
}
func (ms *I64I32Map) SetDirty(key any, dirty bool, editor Editor) {
	if dirty {
		ms.dirtied[key.(int64)] = struct{}{}
		ms.SetParentDirty()
	} else {
		delete(ms.dirtied, key.(int64))
	}

}

func (ms *I64I32Map) FlushDirty(dirty bool) {
	if !dirty {
		ms.dirtied = map[int64]struct{}{}
		ms.deleted = map[int64]struct{}{}
	}
}

func (ms *I64I32Map) Key() any {
	return nil
}
func (ms *I64I32Map) SetKey(i any) {

}

func (ms *I64I32Map) Put(k int64, v int32) *I64I32Map {
	_, exist := ms.value[k]
	if exist {
		ms.deleted[k] = struct{}{}
	}
	ms.value[k] = v
	ms.dirtied[k] = struct{}{}
	ms.SetParentDirty()
	return ms
}

func (ms *I64I32Map) PutAll(kv map[int64]int32) *I64I32Map {
	for k, v := range kv {
		ms.Put(k, v)
	}
	return ms
}

func (ms *I64I32Map) Len() int {
	return len(ms.value)
}

func (ms *I64I32Map) Clear() *I64I32Map {
	if ms.Len() <= 0 {
		return ms
	}
	for k, _ := range ms.value {
		ms.deleted[k] = struct{}{}
	}
	ms.value = map[int64]int32{}
	ms.dirtied = map[int64]struct{}{}
	ms.SetParentDirty()
	return ms
}
func (ms *I64I32Map) Get(k int64) (int32, bool) {
	v, exist := ms.value[k]
	return v, exist
}

func (ms *I64I32Map) Remove(k int64) (int32, bool) {
	v, exist := ms.value[k]
	if !exist {
		return v, false
	}
	delete(ms.value, k)
	delete(ms.dirtied, k)
	ms.deleted[k] = struct{}{}
	ms.SetParentDirty()
	return v, true
}

func (ms *I64I32Map) RemoveAll(k []int64) {
	for _, kk := range k {
		ms.Remove(kk)
	}
}
func (ms *I64I32Map) Iter() iter.Seq2[int64, int32] {
	return func(yield func(int64, int32) bool) {
		for k, v := range ms.value {
			if !yield(k, v) {
				return
			}
		}
	}
}

func (ms *I64I32Map) Dirtied() map[int64]struct{} {
	return ms.dirtied
}

func (ms *I64I32Map) Deleted() map[int64]struct{} {
	return ms.deleted
}

func (ms *I64I32Map) ContainDirtied(kk int64) bool {
	for k := range ms.dirtied {
		if k == kk {
			return true
		}
	}
	return false
}

func (ms *I64I32Map) ContainDeleted(kk int64) bool {
	for k := range ms.deleted {
		if k == kk {
			return true
		}
	}
	return false
}

func (ms *I64I32Map) MergeDirtyFromPb(delList []int64, dirties map[int64]int32) {
	ms.RemoveAll(delList)
	for k, v := range dirties {
		ms.Put(k, v)
	}
}

func (ms *I64I32Map) Entries() map[int64]int32 {
	return ms.value
}

type I64StrMap struct {
	parent      Editor
	keyInParent any
	value       map[int64]string
	dirtied     map[int64]struct{}
	deleted     map[int64]struct{}
}

func NewI64StrMap() *I64StrMap {
	return &I64StrMap{value: make(map[int64]string), dirtied: make(map[int64]struct{}), deleted: make(map[int64]struct{})}
}

func (ms *I64StrMap) SetParent(sync Editor, key any) {
	ms.parent = sync
	ms.keyInParent = key
}
func (ms *I64StrMap) SetParentDirty() {
	if ms.parent != nil {
		ms.parent.SetDirty(ms.keyInParent, true, ms)
	}
}
func (ms *I64StrMap) SetDirty(key any, dirty bool, editor Editor) {
	if dirty {
		ms.dirtied[key.(int64)] = struct{}{}
		ms.SetParentDirty()
	} else {
		delete(ms.dirtied, key.(int64))
	}

}

func (ms *I64StrMap) FlushDirty(dirty bool) {
	if !dirty {
		ms.dirtied = map[int64]struct{}{}
		ms.deleted = map[int64]struct{}{}
	}
}

func (ms *I64StrMap) Key() any {
	return nil
}
func (ms *I64StrMap) SetKey(i any) {

}

func (ms *I64StrMap) Put(k int64, v string) *I64StrMap {
	_, exist := ms.value[k]
	if exist {
		ms.deleted[k] = struct{}{}
	}
	ms.value[k] = v
	ms.dirtied[k] = struct{}{}
	ms.SetParentDirty()
	return ms
}

func (ms *I64StrMap) PutAll(kv map[int64]string) *I64StrMap {
	for k, v := range kv {
		ms.Put(k, v)
	}
	return ms
}

func (ms *I64StrMap) Len() int {
	return len(ms.value)
}

func (ms *I64StrMap) Clear() *I64StrMap {
	if ms.Len() <= 0 {
		return ms
	}
	for k, _ := range ms.value {
		ms.deleted[k] = struct{}{}
	}
	ms.value = map[int64]string{}
	ms.dirtied = map[int64]struct{}{}
	ms.SetParentDirty()
	return ms
}
func (ms *I64StrMap) Get(k int64) (string, bool) {
	v, exist := ms.value[k]
	return v, exist
}

func (ms *I64StrMap) Remove(k int64) (string, bool) {
	v, exist := ms.value[k]
	if !exist {
		return v, false
	}
	delete(ms.value, k)
	delete(ms.dirtied, k)
	ms.deleted[k] = struct{}{}
	ms.SetParentDirty()
	return v, true
}

func (ms *I64StrMap) RemoveAll(k []int64) {
	for _, kk := range k {
		ms.Remove(kk)
	}
}
func (ms *I64StrMap) Iter() iter.Seq2[int64, string] {
	return func(yield func(int64, string) bool) {
		for k, v := range ms.value {
			if !yield(k, v) {
				return
			}
		}
	}
}

func (ms *I64StrMap) Dirtied() map[int64]struct{} {
	return ms.dirtied
}

func (ms *I64StrMap) Deleted() map[int64]struct{} {
	return ms.deleted
}

func (ms *I64StrMap) ContainDirtied(kk int64) bool {
	for k := range ms.dirtied {
		if k == kk {
			return true
		}
	}
	return false
}

func (ms *I64StrMap) ContainDeleted(kk int64) bool {
	for k := range ms.deleted {
		if k == kk {
			return true
		}
	}
	return false
}

func (ms *I64StrMap) MergeDirtyFromPb(delList []int64, dirties map[int64]string) {
	ms.RemoveAll(delList)
	for k, v := range dirties {
		ms.Put(k, v)
	}
}

func (ms *I64StrMap) Entries() map[int64]string {
	return ms.value
}

type I64F32Map struct {
	parent      Editor
	keyInParent any
	value       map[int64]float32
	dirtied     map[int64]struct{}
	deleted     map[int64]struct{}
}

func NewI64F32Map() *I64F32Map {
	return &I64F32Map{value: make(map[int64]float32), dirtied: make(map[int64]struct{}), deleted: make(map[int64]struct{})}
}

func (ms *I64F32Map) SetParent(sync Editor, key any) {
	ms.parent = sync
	ms.keyInParent = key
}
func (ms *I64F32Map) SetParentDirty() {
	if ms.parent != nil {
		ms.parent.SetDirty(ms.keyInParent, true, ms)
	}
}
func (ms *I64F32Map) SetDirty(key any, dirty bool, editor Editor) {
	if dirty {
		ms.dirtied[key.(int64)] = struct{}{}
		ms.SetParentDirty()
	} else {
		delete(ms.dirtied, key.(int64))
	}

}

func (ms *I64F32Map) FlushDirty(dirty bool) {
	if !dirty {
		ms.dirtied = map[int64]struct{}{}
		ms.deleted = map[int64]struct{}{}
	}
}

func (ms *I64F32Map) Key() any {
	return nil
}
func (ms *I64F32Map) SetKey(i any) {

}

func (ms *I64F32Map) Put(k int64, v float32) *I64F32Map {
	_, exist := ms.value[k]
	if exist {
		ms.deleted[k] = struct{}{}
	}
	ms.value[k] = v
	ms.dirtied[k] = struct{}{}
	ms.SetParentDirty()
	return ms
}

func (ms *I64F32Map) PutAll(kv map[int64]float32) *I64F32Map {
	for k, v := range kv {
		ms.Put(k, v)
	}
	return ms
}

func (ms *I64F32Map) Len() int {
	return len(ms.value)
}

func (ms *I64F32Map) Clear() *I64F32Map {
	if ms.Len() <= 0 {
		return ms
	}
	for k, _ := range ms.value {
		ms.deleted[k] = struct{}{}
	}
	ms.value = map[int64]float32{}
	ms.dirtied = map[int64]struct{}{}
	ms.SetParentDirty()
	return ms
}
func (ms *I64F32Map) Get(k int64) (float32, bool) {
	v, exist := ms.value[k]
	return v, exist
}

func (ms *I64F32Map) Remove(k int64) (float32, bool) {
	v, exist := ms.value[k]
	if !exist {
		return v, false
	}
	delete(ms.value, k)
	delete(ms.dirtied, k)
	ms.deleted[k] = struct{}{}
	ms.SetParentDirty()
	return v, true
}

func (ms *I64F32Map) RemoveAll(k []int64) {
	for _, kk := range k {
		ms.Remove(kk)
	}
}
func (ms *I64F32Map) Iter() iter.Seq2[int64, float32] {
	return func(yield func(int64, float32) bool) {
		for k, v := range ms.value {
			if !yield(k, v) {
				return
			}
		}
	}
}

func (ms *I64F32Map) Dirtied() map[int64]struct{} {
	return ms.dirtied
}

func (ms *I64F32Map) Deleted() map[int64]struct{} {
	return ms.deleted
}

func (ms *I64F32Map) ContainDirtied(kk int64) bool {
	for k := range ms.dirtied {
		if k == kk {
			return true
		}
	}
	return false
}

func (ms *I64F32Map) ContainDeleted(kk int64) bool {
	for k := range ms.deleted {
		if k == kk {
			return true
		}
	}
	return false
}

func (ms *I64F32Map) MergeDirtyFromPb(delList []int64, dirties map[int64]float32) {
	ms.RemoveAll(delList)
	for k, v := range dirties {
		ms.Put(k, v)
	}
}

func (ms *I64F32Map) Entries() map[int64]float32 {
	return ms.value
}

type I64U32Map struct {
	parent      Editor
	keyInParent any
	value       map[int64]uint32
	dirtied     map[int64]struct{}
	deleted     map[int64]struct{}
}

func NewI64U32Map() *I64U32Map {
	return &I64U32Map{value: make(map[int64]uint32), dirtied: make(map[int64]struct{}), deleted: make(map[int64]struct{})}
}

func (ms *I64U32Map) SetParent(sync Editor, key any) {
	ms.parent = sync
	ms.keyInParent = key
}
func (ms *I64U32Map) SetParentDirty() {
	if ms.parent != nil {
		ms.parent.SetDirty(ms.keyInParent, true, ms)
	}
}
func (ms *I64U32Map) SetDirty(key any, dirty bool, editor Editor) {
	if dirty {
		ms.dirtied[key.(int64)] = struct{}{}
		ms.SetParentDirty()
	} else {
		delete(ms.dirtied, key.(int64))
	}

}

func (ms *I64U32Map) FlushDirty(dirty bool) {
	if !dirty {
		ms.dirtied = map[int64]struct{}{}
		ms.deleted = map[int64]struct{}{}
	}
}

func (ms *I64U32Map) Key() any {
	return nil
}
func (ms *I64U32Map) SetKey(i any) {

}

func (ms *I64U32Map) Put(k int64, v uint32) *I64U32Map {
	_, exist := ms.value[k]
	if exist {
		ms.deleted[k] = struct{}{}
	}
	ms.value[k] = v
	ms.dirtied[k] = struct{}{}
	ms.SetParentDirty()
	return ms
}

func (ms *I64U32Map) PutAll(kv map[int64]uint32) *I64U32Map {
	for k, v := range kv {
		ms.Put(k, v)
	}
	return ms
}

func (ms *I64U32Map) Len() int {
	return len(ms.value)
}

func (ms *I64U32Map) Clear() *I64U32Map {
	if ms.Len() <= 0 {
		return ms
	}
	for k, _ := range ms.value {
		ms.deleted[k] = struct{}{}
	}
	ms.value = map[int64]uint32{}
	ms.dirtied = map[int64]struct{}{}
	ms.SetParentDirty()
	return ms
}
func (ms *I64U32Map) Get(k int64) (uint32, bool) {
	v, exist := ms.value[k]
	return v, exist
}

func (ms *I64U32Map) Remove(k int64) (uint32, bool) {
	v, exist := ms.value[k]
	if !exist {
		return v, false
	}
	delete(ms.value, k)
	delete(ms.dirtied, k)
	ms.deleted[k] = struct{}{}
	ms.SetParentDirty()
	return v, true
}

func (ms *I64U32Map) RemoveAll(k []int64) {
	for _, kk := range k {
		ms.Remove(kk)
	}
}
func (ms *I64U32Map) Iter() iter.Seq2[int64, uint32] {
	return func(yield func(int64, uint32) bool) {
		for k, v := range ms.value {
			if !yield(k, v) {
				return
			}
		}
	}
}

func (ms *I64U32Map) Dirtied() map[int64]struct{} {
	return ms.dirtied
}

func (ms *I64U32Map) Deleted() map[int64]struct{} {
	return ms.deleted
}

func (ms *I64U32Map) ContainDirtied(kk int64) bool {
	for k := range ms.dirtied {
		if k == kk {
			return true
		}
	}
	return false
}

func (ms *I64U32Map) ContainDeleted(kk int64) bool {
	for k := range ms.deleted {
		if k == kk {
			return true
		}
	}
	return false
}

func (ms *I64U32Map) MergeDirtyFromPb(delList []int64, dirties map[int64]uint32) {
	ms.RemoveAll(delList)
	for k, v := range dirties {
		ms.Put(k, v)
	}
}

func (ms *I64U32Map) Entries() map[int64]uint32 {
	return ms.value
}

type I64I64Map struct {
	parent      Editor
	keyInParent any
	value       map[int64]int64
	dirtied     map[int64]struct{}
	deleted     map[int64]struct{}
}

func NewI64I64Map() *I64I64Map {
	return &I64I64Map{value: make(map[int64]int64), dirtied: make(map[int64]struct{}), deleted: make(map[int64]struct{})}
}

func (ms *I64I64Map) SetParent(sync Editor, key any) {
	ms.parent = sync
	ms.keyInParent = key
}
func (ms *I64I64Map) SetParentDirty() {
	if ms.parent != nil {
		ms.parent.SetDirty(ms.keyInParent, true, ms)
	}
}
func (ms *I64I64Map) SetDirty(key any, dirty bool, editor Editor) {
	if dirty {
		ms.dirtied[key.(int64)] = struct{}{}
		ms.SetParentDirty()
	} else {
		delete(ms.dirtied, key.(int64))
	}

}

func (ms *I64I64Map) FlushDirty(dirty bool) {
	if !dirty {
		ms.dirtied = map[int64]struct{}{}
		ms.deleted = map[int64]struct{}{}
	}
}

func (ms *I64I64Map) Key() any {
	return nil
}
func (ms *I64I64Map) SetKey(i any) {

}

func (ms *I64I64Map) Put(k int64, v int64) *I64I64Map {
	_, exist := ms.value[k]
	if exist {
		ms.deleted[k] = struct{}{}
	}
	ms.value[k] = v
	ms.dirtied[k] = struct{}{}
	ms.SetParentDirty()
	return ms
}

func (ms *I64I64Map) PutAll(kv map[int64]int64) *I64I64Map {
	for k, v := range kv {
		ms.Put(k, v)
	}
	return ms
}

func (ms *I64I64Map) Len() int {
	return len(ms.value)
}

func (ms *I64I64Map) Clear() *I64I64Map {
	if ms.Len() <= 0 {
		return ms
	}
	for k, _ := range ms.value {
		ms.deleted[k] = struct{}{}
	}
	ms.value = map[int64]int64{}
	ms.dirtied = map[int64]struct{}{}
	ms.SetParentDirty()
	return ms
}
func (ms *I64I64Map) Get(k int64) (int64, bool) {
	v, exist := ms.value[k]
	return v, exist
}

func (ms *I64I64Map) Remove(k int64) (int64, bool) {
	v, exist := ms.value[k]
	if !exist {
		return v, false
	}
	delete(ms.value, k)
	delete(ms.dirtied, k)
	ms.deleted[k] = struct{}{}
	ms.SetParentDirty()
	return v, true
}

func (ms *I64I64Map) RemoveAll(k []int64) {
	for _, kk := range k {
		ms.Remove(kk)
	}
}
func (ms *I64I64Map) Iter() iter.Seq2[int64, int64] {
	return func(yield func(int64, int64) bool) {
		for k, v := range ms.value {
			if !yield(k, v) {
				return
			}
		}
	}
}

func (ms *I64I64Map) Dirtied() map[int64]struct{} {
	return ms.dirtied
}

func (ms *I64I64Map) Deleted() map[int64]struct{} {
	return ms.deleted
}

func (ms *I64I64Map) ContainDirtied(kk int64) bool {
	for k := range ms.dirtied {
		if k == kk {
			return true
		}
	}
	return false
}

func (ms *I64I64Map) ContainDeleted(kk int64) bool {
	for k := range ms.deleted {
		if k == kk {
			return true
		}
	}
	return false
}

func (ms *I64I64Map) MergeDirtyFromPb(delList []int64, dirties map[int64]int64) {
	ms.RemoveAll(delList)
	for k, v := range dirties {
		ms.Put(k, v)
	}
}

func (ms *I64I64Map) Entries() map[int64]int64 {
	return ms.value
}
