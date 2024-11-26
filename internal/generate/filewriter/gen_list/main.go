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
	fw.WriteString("import \"slices\"\n")
	fw.WriteString("import \"github.com/yaoguangduan/proto-editor/protodef\"\n")
	kts := map[string]string{"U32": "uint32", "U64": "uint64", "I32": "int32", "I64": "int64", "Str": "string", "Bool": "bool", "F32": "float32", "F64": "float64"}

	for k, v := range kts {
		fw.WriteString(fmt.Sprintf(`
			type %sList struct {
				parent   Editor
				idxInPar any
				value    []%s
				monitor  []%s
				change   []*protodef.ListOpLog
			}
			
			func New%sList() *%sList {
				return &%sList{value: make([]%s, 0)}
			}
			`, k, v, v, k, k, k, v))
		fw.WriteString(fmt.Sprintf(`
			func (ms *%sList) SetParent(editor Editor, key any) {
				ms.parent = editor
				ms.idxInPar = key
			}
			func (ms *%sList) SetParentDirty() {
				if ms.parent != nil {
					ms.parent.SetDirty(ms.idxInPar, true, ms)
				}
			}
			func (ms *%sList) SetDirty(key any, dirty bool, sync Editor) {
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
			func (ms *%sList) FlushDirty(dirty bool) {
				if !dirty {
					ms.monitor = nil
					ms.change = nil
				}
			}
			`, k, k, k, k))
		fw.WriteString(fmt.Sprintf(`
			func (ms *%sList) Add(t %s) *%sList {
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
			
			func (ms *%sList) AddAll(t []%s) *%sList {
				for _, v := range t {
					ms.Add(v)
				}
				return ms
			}
			`, k, v, k,
			k, v, k))

		fw.WriteString(fmt.Sprintf(`
			func (ms *%sList) Insert(pos int, v %s) *%sList {
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
			func (ms *%sList) Set(pos int, v %s) *%sList {
				ms.value[pos] = v
				ms.SetDirty(pos, true, ms)
				return ms
			}
			func (ms *%sList) Iter() iter.Seq2[int, %s] {
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
		fw.WriteString(fmt.Sprintf(`
			func (ms *%sList) Clear() {
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
			func (ms *%sList) Len() int {
				return len(ms.value)
			}
			
			func (ms *%sList) Remove(v %s) {
				ms.RemoveByIdx(slices.Index(ms.value, v))
			}
			
			func (ms *%sList) RemoveByIdx(pos int) {
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
		fw.WriteString(fmt.Sprintf(`
			func (ms *%sList) ChangeLog() []*protodef.ListOpLog {
				return ms.change
			}
			func (ms *%sList) MonitorDirtyList() []%s {
				return ms.monitor
			}
			func (ms *%sList) Values() []%s {
				vv := make([]%s, ms.Len())
				copy(vv, ms.value)
				return vv
			}
			`, k,
			k, v,
			k, v,
			v))

		fw.WriteString(fmt.Sprintf(`
		func (ms *%sList) MergeDirtyFromPb(logs []*protodef.ListOpLog, dirties []%s) {
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
		`, k, v))
		fw.WriteString(fmt.Sprintf(`
		func (ms *%sList) Get(pos int32) %s {
			return ms.value[pos]
		}
		`, k, v))
	}

	source, err := format.Source([]byte(fw.String()))
	err = os.WriteFile("D:\\code\\go\\gogame\\proto-editor\\pbeditor\\extends_list.go", source, os.ModePerm)
	if err != nil {
		panic(err)
	}
}
