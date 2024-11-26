package main

import (
	"fmt"
	"go/format"
	"os"
	"strings"
)

func main() {
	fw := strings.Builder{}
	fw.WriteString("package pbeditor\n")
	fw.WriteString("import \"iter\"\n")
	kts := map[string]string{"U32": "uint32", "U64": "uint64", "I32": "int32", "I64": "int64", "Str": "string", "Bool": "bool"}
	vts := map[string]string{"U32": "uint32", "U64": "uint64", "I32": "int32", "I64": "int64", "Str": "string", "Bool": "bool", "F32": "float32", "F64": "float64", "Bytes": "[]byte"}
	//vts := map[string]string{"Editor": "Editor"}
	for kt, kv := range kts {
		for vt, vv := range vts {
			fw.WriteString(fmt.Sprintf(`type %s%sMap struct {
	parent      Editor
	keyInParent any
	value       map[%s]%s
	dirtied     map[%s]struct{}
	deleted     map[%s]struct{}
}
`, kt, vt, kv, vv, kv, kv))
			fw.WriteString(fmt.Sprintf(`
func New%s%sMap() *%s%sMap {
	return &%s%sMap{value: make(map[%s]%s), dirtied: make(map[%s]struct{}), deleted: make(map[%s]struct{})}
}
`, kt, vt, kt, vt, kt, vt, kv, vv, kv, kv))
			fw.WriteString(fmt.Sprintf(`
func (ms *%s%sMap) SetParent(sync Editor, key any) {
	ms.parent = sync
	ms.keyInParent = key
}
func (ms *%s%sMap) SetParentDirty() {
	if ms.parent != nil {
		ms.parent.SetDirty(ms.keyInParent, true, ms)
	}
}
func (ms *%s%sMap) SetDirty(key any, dirty bool, editor Editor) {
	if dirty {
		ms.dirtied[key.(%s)] = struct{}{}
		ms.SetParentDirty()
	} else {
		delete(ms.dirtied, key.(%s))
	}

}
`, kt, vt, kt, vt, kt, vt, kv, kv))

			fw.WriteString(fmt.Sprintf(`
func (ms *%s%sMap) FlushDirty(dirty bool) {
	if !dirty {
		ms.dirtied = map[%s]struct{}{}
		ms.deleted = map[%s]struct{}{}
	}
}

func (ms *%s%sMap) Key() any {
	return nil
}
func (ms *%s%sMap) SetKey(i any) {

}
`, kt, vt, kv, kv, kt, vt, kt, vt))

			fw.WriteString(fmt.Sprintf(`
func (ms *%s%sMap) Put(k %s, v %s) *%s%sMap {
	_, exist := ms.value[k]
	if exist {
		ms.deleted[k] = struct{}{}
	}
	ms.value[k] = v
	ms.dirtied[k] = struct{}{}
	ms.SetParentDirty()
	return ms
}

func (ms *%s%sMap) PutAll(kv map[%s]%s) *%s%sMap {
	for k, v := range kv {
		ms.Put(k,v)
	}
	return ms
}

func (ms *%s%sMap) Len() int {
	return len(ms.value)
}
`, kt, vt, kv, vv, kt, vt, kt, vt, kv, vv, kt, vt, kt, vt))

			fw.WriteString(fmt.Sprintf(`
func (ms *%s%sMap) Clear() *%s%sMap {
	if ms.Len() <= 0 {
		return ms
	}
	for k, _ := range ms.value {
		ms.deleted[k] = struct{}{}
	}
	ms.value = map[%s]%s{}
	ms.dirtied = map[%s]struct{}{}
	ms.SetParentDirty()
	return ms
}
func (ms *%s%sMap) Get(k %s) (%s,bool) {
	v ,exist:= ms.value[k]
	return v,exist
}

func (ms *%s%sMap) Remove(k %s) (%s,bool) {
	v, exist := ms.value[k]
	if !exist {
		return v,false
	}
	delete(ms.value, k)
	delete(ms.dirtied, k)
	ms.deleted[k] = struct{}{}
	ms.SetParentDirty()
	return v,true
}
`, kt, vt, kt, vt, kv, vv, kv, kt, vt, kv, vv, kt, vt, kv, vv))

			fw.WriteString(fmt.Sprintf(`
				func (ms *%s%sMap) RemoveAll(k []%s) {
					for _, kk := range k {
						ms.Remove(kk)
					}
				}
				func (ms *%s%sMap)  Iter() iter.Seq2[%s, %s] {
					return func(yield func(%s, %s) bool) {
						for k, v := range ms.value {
							if !yield(k, v) {
								return
							}
						}
					}
				}
				`, kt, vt, kv, kt, vt, kv, vv, kv, vv))

			fw.WriteString(fmt.Sprintf(`
func (ms *%s%sMap) Dirtied() map[%s]struct{} {
	return ms.dirtied
}

func (ms *%s%sMap) Deleted() map[%s]struct{} {
	return ms.deleted
}

func (ms *%s%sMap) ContainDirtied(kk %s) bool {
	for k := range ms.dirtied {
		if k == kk {
			return true
		}
	}
	return false
}

func (ms *%s%sMap) ContainDeleted(kk %s) bool {
	for k := range ms.deleted {
		if k == kk {
			return true
		}
	}
	return false
}
`, kt, vt, kv, kt, vt, kv, kt, vt, kv, kt, vt, kv))

			fw.WriteString(fmt.Sprintf(`
		func (ms *%s%sMap) MergeDirtyFromPb(delList []%s, dirties map[%s]%s) {
			ms.RemoveAll(delList)
			for k, v := range dirties {
				ms.Put(k, v)
			}
		}
		`, kt, vt, kv, kv, vv))
			fw.WriteString(fmt.Sprintf(`
		func (ms *%s%sMap) Entries() map[%s]%s {
			return ms.value
		}
		`, kt, vt, kv, vv))
		}
	}
	fmt.Println(fw.String())

	source, err := format.Source([]byte(fw.String()))
	err = os.WriteFile("D:\\code\\go\\gogame\\proto-editor\\pbeditor\\extends_map.go", source, os.ModePerm)
	if err != nil {
		panic(err)
	}
}
