package generate

import (
	"fmt"
	"github.com/yaoguangduan/proto-editor/internal/generate/util"
	"google.golang.org/protobuf/compiler/protogen"
)

func GenMapEditorForEnum(gen *protogen.Plugin, g *protogen.GeneratedFile, f *util.FileInfo, msg *protogen.Enum, field *protogen.Field) {

	mapInfo := util.MapFieldInfo(g, f, field)
	g.P(fmt.Sprintf(`
		type %s%sEditorMap struct {
			parent      pbeditor.Editor
			keyInParent any
			value       map[%s]%s
			dirtied     map[%s]struct{}
			deleted     map[%s]struct{}
		}
		`, mapInfo.KeyAlias, mapInfo.ValueType, mapInfo.KeyType, mapInfo.ValueEditorType, mapInfo.KeyType, mapInfo.KeyType))
	g.P(fmt.Sprintf(`
		func New%s%sEditorMap() *%s%sEditorMap {
			return &%s%sEditorMap{value: make(map[%s]%s), dirtied: make(map[%s]struct{}), deleted: make(map[%s]struct{})}
		}`, mapInfo.KeyAlias, mapInfo.ValueType, mapInfo.KeyAlias, mapInfo.ValueType, mapInfo.KeyAlias, mapInfo.ValueType, mapInfo.KeyType, mapInfo.ValueEditorType, mapInfo.KeyType, mapInfo.KeyType))

	g.P(fmt.Sprintf(`
		func (ms *%s%sEditorMap) SetParent(editor pbeditor.Editor, key any) {
			ms.parent = editor
			ms.keyInParent = key
		}`, mapInfo.KeyAlias, mapInfo.ValueType))

	g.P(fmt.Sprintf(`
		func (ms *%s%sEditorMap) SetParentDirty() {
			if ms.parent != nil {
				ms.parent.SetDirty(ms.keyInParent, true, ms)
			}
		}`, mapInfo.KeyAlias, mapInfo.ValueType))

	g.P(fmt.Sprintf(`
		func (ms *%s%sEditorMap) Entries() map[%s]%s {
			return ms.value
		}`, mapInfo.KeyAlias, mapInfo.ValueType, mapInfo.KeyType, mapInfo.ValueType))
	g.P(fmt.Sprintf(`
		func (ms *%s%sEditorMap) SetDirty(key any, dirty bool, editor pbeditor.Editor) {
			if dirty {
				ms.dirtied[key.(%s)] = struct{}{}
				ms.SetParentDirty()
			} else {
				delete(ms.dirtied, key.(%s))
			}
		
		}
		`, mapInfo.KeyAlias, mapInfo.ValueType, mapInfo.KeyType, mapInfo.KeyType))

	g.P(fmt.Sprintf(`
		func (ms *%s%sEditorMap) FlushDirty(dirty bool) {
			if !dirty {
				ms.dirtied = map[%s]struct{}{}
				ms.deleted = map[%s]struct{}{}
			}
		}
		`, mapInfo.KeyAlias, mapInfo.ValueType, mapInfo.KeyType, mapInfo.KeyType))

	g.P(fmt.Sprintf(`
		func (ms *%s%sEditorMap) Put(k %s, v %s) *%s%sEditorMap {
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
		
		func (ms *%s%sEditorMap) PutAll(kv map[%s]%s) *%s%sEditorMap {
			for k, v := range kv {
				ms.Put(k, v)
			}
			return ms
		}`, mapInfo.KeyAlias, mapInfo.ValueType, mapInfo.KeyType, mapInfo.ValueEditorType, mapInfo.KeyAlias, mapInfo.ValueType, mapInfo.KeyAlias, mapInfo.ValueType, mapInfo.KeyType, mapInfo.ValueEditorType, mapInfo.KeyAlias, mapInfo.ValueType))

	g.P(fmt.Sprintf(`
		func (ms *%s%sEditorMap) Len() int {
			return len(ms.value)
		}
		
		func (ms *%s%sEditorMap) Clear() *%s%sEditorMap {
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
		}`, mapInfo.KeyAlias, mapInfo.ValueType, mapInfo.KeyAlias, mapInfo.ValueType, mapInfo.KeyAlias, mapInfo.ValueType, mapInfo.KeyType, mapInfo.ValueEditorType, mapInfo.KeyType))

	g.P(fmt.Sprintf(`
		func (ms *%s%sEditorMap) Get(k %s) (%s,bool) {
			v ,ok := ms.value[k]
			return v,ok
		}
		
		func (ms *%s%sEditorMap) Remove(k %s) (%s,bool) {
			v, exist := ms.value[k]
			if !exist {
				return v,exist
			}
			delete(ms.value, k)
			delete(ms.dirtied, k)
			ms.deleted[k] = struct{}{}
			ms.SetParentDirty()
			return v,exist
		}
		`, mapInfo.KeyAlias, mapInfo.ValueType, mapInfo.KeyType, mapInfo.ValueEditorType, mapInfo.KeyAlias, mapInfo.ValueType, mapInfo.KeyType, mapInfo.ValueEditorType))
	g.P(fmt.Sprintf(`
		func (ms *%s%sEditorMap) RemoveAll(k []%s) {
			for _, kk := range k {
				ms.Remove(kk)
			}
		}
		
		func (ms *%s%sEditorMap) Iter() iter.Seq2[%s, %s] {
			return func(yield func(%s, %s) bool) {
				for k, v := range ms.value {
					if !yield(k, v) {
						return
					}
				}
			}
		}
		`, mapInfo.KeyAlias, mapInfo.ValueType, mapInfo.KeyType, mapInfo.KeyAlias, mapInfo.ValueType, mapInfo.KeyType, mapInfo.ValueEditorType, mapInfo.KeyType, mapInfo.ValueEditorType))
	g.P(fmt.Sprintf(`
		func (ms *%s%sEditorMap) Dirtied() map[%s]struct{} {
			return ms.dirtied
		}
		
		func (ms *%s%sEditorMap) Deleted() map[%s]struct{} {
			return ms.deleted
		}
		
		func (ms *%s%sEditorMap) ContainDirtied(kk %s) bool {
			for k := range ms.dirtied {
				if k == kk {
					return true
				}
			}
			return false
		}
		
		func (ms *%s%sEditorMap) ContainDeleted(kk %s) bool {
			for k := range ms.deleted {
				if k == kk {
					return true
				}
			}
			return false
		}
		`, mapInfo.KeyAlias, mapInfo.ValueType, mapInfo.KeyType,
		mapInfo.KeyAlias, mapInfo.ValueType, mapInfo.KeyType,
		mapInfo.KeyAlias, mapInfo.ValueType, mapInfo.KeyType,
		mapInfo.KeyAlias, mapInfo.ValueType, mapInfo.KeyType,
	))

	g.P(fmt.Sprintf(`
		func (ms *%s%sEditorMap) MergeDirtyFromPb(delList []%s, dirties map[%s]%s) {
			ms.RemoveAll(delList)
			for k, v := range dirties {
				ms.Put(k, v)
			}
		}
		`, mapInfo.KeyAlias, mapInfo.ValueType, mapInfo.KeyType, mapInfo.KeyType, mapInfo.ValueType))
}
func GenMapEditor(gen *protogen.Plugin, g *protogen.GeneratedFile, f *util.FileInfo, msg *protogen.Message, field *protogen.Field) {
	mapInfo := util.MapFieldInfo(g, f, field)
	g.P(fmt.Sprintf(`
		type %s%sEditorMap struct {
			parent      pbeditor.Editor
			keyInParent any
			value       map[%s]%s
			dirtied     map[%s]struct{}
			deleted     map[%s]struct{}
		}
		`, mapInfo.KeyAlias, mapInfo.ValueType, mapInfo.KeyType, mapInfo.ValueEditorType, mapInfo.KeyType, mapInfo.KeyType))
	g.P(fmt.Sprintf(`
		func New%s%sEditorMap() *%s%sEditorMap {
			return &%s%sEditorMap{value: make(map[%s]%s), dirtied: make(map[%s]struct{}), deleted: make(map[%s]struct{})}
		}`, mapInfo.KeyAlias, mapInfo.ValueType, mapInfo.KeyAlias, mapInfo.ValueType, mapInfo.KeyAlias, mapInfo.ValueType, mapInfo.KeyType, mapInfo.ValueEditorType, mapInfo.KeyType, mapInfo.KeyType))

	g.P(fmt.Sprintf(`
		func (ms *%s%sEditorMap) SetParent(editor pbeditor.Editor, key any) {
			ms.parent = editor
			ms.keyInParent = key
		}`, mapInfo.KeyAlias, mapInfo.ValueType))

	g.P(fmt.Sprintf(`
		func (ms *%s%sEditorMap) SetParentDirty() {
			if ms.parent != nil {
				ms.parent.SetDirty(ms.keyInParent, true, ms)
			}
		}`, mapInfo.KeyAlias, mapInfo.ValueType))

	g.P(fmt.Sprintf(`
		func (ms *%s%sEditorMap) Entries() map[%s]%s{
			return ms.value
		}`, mapInfo.KeyAlias, mapInfo.ValueType, mapInfo.KeyType, mapInfo.ValueEditorType))
	g.P(fmt.Sprintf(`
		func (ms *%s%sEditorMap) SetDirty(key any, dirty bool, editor pbeditor.Editor) {
			if dirty {
				ms.dirtied[key.(%s)] = struct{}{}
				ms.SetParentDirty()
			} else {
				delete(ms.dirtied, key.(%s))
			}
		
		}
		`, mapInfo.KeyAlias, mapInfo.ValueType, mapInfo.KeyType, mapInfo.KeyType))

	g.P(fmt.Sprintf(`
		func (ms *%s%sEditorMap) FlushDirty(dirty bool) {
			if !dirty {
				ms.dirtied = map[%s]struct{}{}
				ms.deleted = map[%s]struct{}{}
			}
			for _,v := range ms.value {
				v.FlushDirty(dirty)
			}
		}
		`, mapInfo.KeyAlias, mapInfo.ValueType, mapInfo.KeyType, mapInfo.KeyType))

	g.P(fmt.Sprintf(`
		func (ms *%s%sEditorMap) Put(k %s, v %s) *%s%sEditorMap {
			v.SetParent(ms, k)
			old, exist := ms.value[k]
			if exist {
				if old != v {
					old.SetParent(nil,nil)
					ms.deleted[k] = struct{}{}
				}
			}
			v.FlushDirty(true)
			ms.value[k] = v
			ms.dirtied[k] = struct{}{}
			ms.SetParentDirty()
			return ms
		}
		
		func (ms *%s%sEditorMap) PutAll(kv map[%s]%s) *%s%sEditorMap {
			for k, v := range kv {
				ms.Put(k, v)
			}
			return ms
		}`, mapInfo.KeyAlias, mapInfo.ValueType, mapInfo.KeyType, mapInfo.ValueEditorType, mapInfo.KeyAlias, mapInfo.ValueType, mapInfo.KeyAlias, mapInfo.ValueType, mapInfo.KeyType, mapInfo.ValueEditorType, mapInfo.KeyAlias, mapInfo.ValueType))

	g.P(fmt.Sprintf(`
		func (ms *%s%sEditorMap) Len() int {
			return len(ms.value)
		}
		
		func (ms *%s%sEditorMap) Clear() *%s%sEditorMap {
			if ms.Len() <= 0 {
				return ms
			}
			for k, v := range ms.value {
				v.SetParent(nil,nil)
				ms.deleted[k] = struct{}{}
			}
			ms.value = map[%s]%s{}
			ms.dirtied = map[%s]struct{}{}
			ms.SetParentDirty()
			return ms
		}`, mapInfo.KeyAlias, mapInfo.ValueType, mapInfo.KeyAlias, mapInfo.ValueType, mapInfo.KeyAlias, mapInfo.ValueType, mapInfo.KeyType, mapInfo.ValueEditorType, mapInfo.KeyType))

	g.P(fmt.Sprintf(`
		func (ms *%s%sEditorMap) Get(k %s) %s {
			v := ms.value[k]
			return v
		}
		
		func (ms *%s%sEditorMap) Remove(k %s) %s {
			v, exist := ms.value[k]
			if !exist {
				return v
			}
			v.SetParent(nil,nil)
			delete(ms.value, k)
			delete(ms.dirtied, k)
			ms.deleted[k] = struct{}{}
			ms.SetParentDirty()
			return v
		}
		`, mapInfo.KeyAlias, mapInfo.ValueType, mapInfo.KeyType, mapInfo.ValueEditorType, mapInfo.KeyAlias, mapInfo.ValueType, mapInfo.KeyType, mapInfo.ValueEditorType))
	g.P(fmt.Sprintf(`
		func (ms *%s%sEditorMap) RemoveAll(k []%s) {
			for _, kk := range k {
				ms.Remove(kk)
			}
		}
		
		func (ms *%s%sEditorMap) Iter() iter.Seq2[%s, %s] {
			return func(yield func(%s, %s) bool) {
				for k, v := range ms.value {
					if !yield(k, v) {
						return
					}
				}
			}
		}
		`, mapInfo.KeyAlias, mapInfo.ValueType, mapInfo.KeyType, mapInfo.KeyAlias, mapInfo.ValueType, mapInfo.KeyType, mapInfo.ValueEditorType, mapInfo.KeyType, mapInfo.ValueEditorType))
	g.P(fmt.Sprintf(`
		func (ms *%s%sEditorMap) Dirtied() map[%s]struct{} {
			return ms.dirtied
		}
		
		func (ms *%s%sEditorMap) Deleted() map[%s]struct{} {
			return ms.deleted
		}
		
		func (ms *%s%sEditorMap) ContainDirtied(kk %s) bool {
			for k := range ms.dirtied {
				if k == kk {
					return true
				}
			}
			return false
		}
		
		func (ms *%s%sEditorMap) ContainDeleted(kk %s) bool {
			for k := range ms.deleted {
				if k == kk {
					return true
				}
			}
			return false
		}
		`, mapInfo.KeyAlias, mapInfo.ValueType, mapInfo.KeyType,
		mapInfo.KeyAlias, mapInfo.ValueType, mapInfo.KeyType,
		mapInfo.KeyAlias, mapInfo.ValueType, mapInfo.KeyType,
		mapInfo.KeyAlias, mapInfo.ValueType, mapInfo.KeyType,
	))

	g.P(fmt.Sprintf(`
		func (ms *%s%sEditorMap) MergeDirtyFromPb(delList []%s, dirties map[%s]*%s) {
			ms.RemoveAll(delList)
			for k, v := range dirties {
				var tmp = ms.Get(k)
				if tmp == nil {
					tmp = New%sEditor()
					tmp.MergeDirtyFromPb(v)
					ms.Put(k, tmp)
				} else {
					tmp.MergeDirtyFromPb(v)
				}
			}
		}
		`, mapInfo.KeyAlias, mapInfo.ValueType, mapInfo.KeyType, mapInfo.KeyType, mapInfo.ValueType,
		mapInfo.ValueType))
}
