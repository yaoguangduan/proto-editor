package generate

import (
	"fmt"
	"github.com/yaoguangduan/proto-editor/internal/generate/util"
	"google.golang.org/protobuf/compiler/protogen"
	"strings"
)

func GenListEditorForEnum(gen *protogen.Plugin, g *protogen.GeneratedFile, f *util.FileInfo, msg *protogen.Enum, field *protogen.Field) {
	fullName := g.QualifiedGoIdent(field.Enum.GoIdent)
	v := fullName
	k := fullName
	idx := strings.Index(fullName, ".")
	if idx != -1 {
		k = fullName[idx+1:]
	}
	g.P(fmt.Sprintf(`
			type %sEditorList struct {
				parent   pbeditor.Editor
				idxInPar any
				value    []%s
				monitor  []%s
				change   []*protodef.ListOpLog
			}
			
			func New%sEditorList() *%sEditorList {
				return &%sEditorList{value: make([]%s, 0)}
			}
			`, k, v, v, k, k, k, v))

	g.P(fmt.Sprintf(`
			func (ms *%sEditorList) SetParent(editor pbeditor.Editor, key any) {
				ms.parent = editor
				ms.idxInPar = key
			}
			func (ms *%sEditorList) SetParentDirty() {
				if ms.parent != nil {
					ms.parent.SetDirty(ms.idxInPar, true, ms)
				}
			}
			func (ms *%sEditorList) SetDirty(key any, dirty bool, sync pbeditor.Editor) {
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
			func (ms *%sEditorList) FlushDirty(dirty bool) {
				if !dirty {
					ms.monitor = nil
					ms.change = nil
				}
			}
			`, k, k, k, k))
	g.P(fmt.Sprintf(`
			func (ms *%sEditorList) Add(t %s) *%sEditorList {
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
			
			func (ms *%sEditorList) AddAll(t []%s) *%sEditorList {
				for _, v := range t {
					ms.Add(v)
				}
				return ms
			}
			`, k, v, k,
		k, v, k))

	g.P(fmt.Sprintf(`
			func (ms *%sEditorList) Insert(pos int, v %s) *%sEditorList {
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
			func (ms *%sEditorList) Set(pos int, v %s) *%sEditorList {
				ms.value[pos] = v
				ms.SetDirty(pos, true, ms)
				return ms
			}
			func (ms *%sEditorList) Iter() iter.Seq2[int, %s] {
				return func(yield func(int, %s) bool) {
					for i := range ms.value {
						if !yield(i, ms.value[i]) {
							return
						}
					}
				}
			}
			`, k, v, k,
		k, v, k,
		k, v,
		v))
	g.P(fmt.Sprintf(`
			func (ms *%sEditorList) Clear() {
				ms.value = make([]%s, 0)
				ms.monitor = make([]%s, 0)
				ms.change = nil
				ms.change = append(ms.change, &protodef.ListOpLog{
					OpType: protodef.ListOpType_Clear,
					Pos:    0,
					RefIdx: 0,
				})
				ms.SetParentDirty()
			}
			func (ms *%sEditorList) Len() int {
				return len(ms.value)
			}
			
			func (ms *%sEditorList) Remove(v %s) {
				ms.RemoveByIdx(slices.Index(ms.value, v))
			}
			
			func (ms *%sEditorList) RemoveByIdx(pos int) {
				ms.value = slices.Delete(ms.value, pos, pos+1)
				ms.change = append(ms.change, &protodef.ListOpLog{
					OpType: protodef.ListOpType_Delete,
					Pos:    int32(pos),
				})
				ms.SetParentDirty()
			}`, k,
		v,
		v,
		k,
		k, v,
		k))
	g.P(fmt.Sprintf(`
			func (ms *%sEditorList) ChangeLog() []*protodef.ListOpLog {
				return ms.change
			}
			func (ms *%sEditorList) MonitorDirtyList() []%s {
				return ms.monitor
			}
			func (ms *%sEditorList) Values() []%s {
				vv := make([]%s, ms.Len())
				copy(vv, ms.value)
				return vv
			}
			`, k,
		k, v,
		k, v,
		v))

	g.P(fmt.Sprintf(`
		func (ms *%sEditorList) MergeDirtyFromPb(logs []*protodef.ListOpLog, dirties []%s) {
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
					ms.Set(int(log.Pos),dirties[log.RefIdx])
				case protodef.ListOpType_Delete:
					ms.RemoveByIdx(int(log.Pos))
				case protodef.ListOpType_Clear:
					ms.Clear()
				}
			}
		}
		`, k, k))
	g.P(fmt.Sprintf(`
		func (ms *%sEditorList) Get(pos int32) %s {
			return ms.value[pos]
		}
		`, k, k))

}

func GenListEditor(gen *protogen.Plugin, g *protogen.GeneratedFile, f *util.FileInfo, msg *protogen.Message, field *protogen.Field) {
	fullName := g.QualifiedGoIdent(field.Message.GoIdent)
	v := "*" + fullName + "Editor"
	k := fullName
	idx := strings.Index(fullName, ".")
	if idx != -1 {
		k = fullName[idx+1:]
	}
	g.P(fmt.Sprintf(`
			type %sEditorList struct {
				parent   pbeditor.Editor
				idxInPar any
				value    []%s
				monitor  []%s
				change   []*protodef.ListOpLog
			}
			
			func New%sEditorList() *%sEditorList {
				return &%sEditorList{value: make([]%s, 0)}
			}
			`, k, v, v, k, k, k, v))

	g.P(fmt.Sprintf(`
			func (ms *%sEditorList) SetParent(editor pbeditor.Editor, key any) {
				ms.parent = editor
				ms.idxInPar = key
			}
			func (ms *%sEditorList) SetParentDirty() {
				if ms.parent != nil {
					ms.parent.SetDirty(ms.idxInPar, true, ms)
				}
			}
			func (ms *%sEditorList) SetDirty(key any, dirty bool, sync pbeditor.Editor) {
				var pos = -1
				if key != nil {
					pos = key.(int)
				} else {
					pos = slices.Index(ms.value,sync.(%s))
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
			func (ms *%sEditorList) FlushDirty(dirty bool) {
				if !dirty {
					ms.monitor = nil
					ms.change = nil
				}
				for _,v := range ms.value {
					v.FlushDirty(dirty)
				}
			}
			`, k, k, k, v, k))
	g.P(fmt.Sprintf(`
			func (ms *%sEditorList) Add(t %s) *%sEditorList {
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
			
			func (ms *%sEditorList) AddAll(t []%s) *%sEditorList {
				for _, v := range t {
					ms.Add(v)
				}
				return ms
			}
			`, k, v, k,
		k, v, k))

	g.P(fmt.Sprintf(`
			func (ms *%sEditorList) Insert(pos int, v %s) *%sEditorList {
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
			func (ms *%sEditorList) Set(pos int, v %s) *%sEditorList {
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
			func (ms *%sEditorList) Iter() iter.Seq2[int, %s] {
				return func(yield func(int, %s) bool) {
					for i := range ms.value {
						if !yield(i, ms.value[i]) {
							return
						}
					}
				}
			}
			`, k, v, k,
		k, v, k,
		k, v,
		v))
	g.P(fmt.Sprintf(`
			func (ms *%sEditorList) Clear() {
				for _,v := range ms.value {
					v.SetParent(nil, nil)
				}
				ms.value = make([]%s, 0)
				ms.monitor = make([]%s, 0)
				ms.change = nil
				ms.change = append(ms.change, &protodef.ListOpLog{
					OpType: protodef.ListOpType_Clear,
					Pos:    0,
					RefIdx: 0,
				})
				ms.SetParentDirty()
			}
			func (ms *%sEditorList) Len() int {
				return len(ms.value)
			}
			
			func (ms *%sEditorList) Remove(v %s) {
				ms.RemoveByIdx(slices.Index(ms.value, v))
			}
			
			func (ms *%sEditorList) RemoveByIdx(pos int) {
				v := ms.value[pos]
				v.SetParent(nil, nil)
				ms.value = slices.Delete(ms.value, pos, pos+1)
				ms.change = append(ms.change, &protodef.ListOpLog{
					OpType: protodef.ListOpType_Delete,
					Pos:    int32(pos),
				})
				ms.SetParentDirty()
			}`, k,
		v,
		v,
		k,
		k, v,
		k))
	g.P(fmt.Sprintf(`
			func (ms *%sEditorList) ChangeLog() []*protodef.ListOpLog {
				return ms.change
			}
			func (ms *%sEditorList) MonitorDirtyList() []%s {
				return ms.monitor
			}
			func (ms *%sEditorList) Values() []%s {
				vv := make([]%s, ms.Len())
				copy(vv, ms.value)
				return vv
			}
			`, k,
		k, v,
		k, v,
		v))

	g.P(fmt.Sprintf(`
		func (ms *%sEditorList) MergeDirtyFromPb(logs []*protodef.ListOpLog, dirties []*%s) {
			for _, log := range logs {
				switch log.OpType {
				case protodef.ListOpType_Insert:
					v := dirties[log.RefIdx]
					tmp := New%sEditor().CopyFrom(v)
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
		`, k, k, k))
	g.P(fmt.Sprintf(`
		func (ms *%sEditorList) Get(pos int32) *%sEditor {
			return ms.value[pos]
		}
		`, k, k))

}
