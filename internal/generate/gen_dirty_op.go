package generate

import (
	"fmt"
	"github.com/samber/lo"
	"github.com/yaoguangduan/proto-editor/internal/generate/util"
	"github.com/yaoguangduan/proto-editor/pbeditor"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/reflect/protoreflect"
)

func GenDirtyFunc(gen *protogen.Plugin, g *protogen.GeneratedFile, f *util.FileInfo, msg *protogen.Message) {
	editorMsgName := util.EditorMsgName(g.QualifiedGoIdent(msg.GoIdent))
	var hasArrOrMap = lo.CountBy(msg.Fields, func(item *protogen.Field) bool {
		return item.Desc.IsMap() || item.Desc.IsList()
	}) > 0

	g.P(fmt.Sprintf("func (x *%s) MergeDirtyFromPb(r *%s) {", editorMsgName, g.QualifiedGoIdent(msg.GoIdent)))
	if hasArrOrMap {
		g.P("var mutState = pbeditor.UnmarshalMutState(r.ProtoReflect().GetUnknown())")
	}
	for _, field := range msg.Fields {
		if (util.IsBuildInType(field.Desc.Kind()) || field.Desc.Kind() == protoreflect.EnumKind) && !field.Desc.IsList() && field.Desc.Kind() != protoreflect.BytesKind {
			g.P(fmt.Sprintf("if r.%s != nil {", field.GoName))
			g.P(fmt.Sprintf("x.Set%s(*r.%s)", field.GoName, field.GoName))
			g.P(fmt.Sprintf("}"))
		}
		if field.Desc.Kind() == protoreflect.BytesKind {
			g.P(fmt.Sprintf("if len(r.%s) > 0 {", field.GoName))
			g.P(fmt.Sprintf("x.Set%s(slices.Clone(r.%s))", field.GoName, field.GoName))
			g.P(fmt.Sprintf("}"))
		}
		if field.Desc.IsList() {
			g.P(fmt.Sprintf("x.Get%s().MergeDirtyFromPb(mutState.ListOpLog(%d),r.%s)", field.GoName, field.Desc.Number(), field.GoName))
		}
		if util.IsMsg(field) {
			g.P(fmt.Sprintf("if r.%s != nil {", field.GoName))
			g.P(fmt.Sprintf("x.Get%s().MergeDirtyFromPb(r.%s)", field.GoName, field.GoName))
			g.P(fmt.Sprintf("}"))
		}
		if field.Desc.IsMap() {
			mapFieldInfo := util.MapFieldInfo(g, f, field)
			g.P(fmt.Sprintf("x.Get%s().MergeDirtyFromPb(mutState.%sList(%d),r.%s)", field.GoName, util.Capital(mapFieldInfo.KeyAlias), field.Desc.Number(), field.GoName))
		}
	}
	g.P(fmt.Sprintf("}"))

	g.P(fmt.Sprintf("func (x *%s) MergeDirtyFromBytes(buf []byte) *%s{", editorMsgName, editorMsgName))
	g.P(fmt.Sprintf("fds := pbeditor.ToRawMessage(buf)"))

	if hasArrOrMap {
		for _, field := range msg.Fields {
			fieldName := util.EditorFieldName(field.GoName)
			if field.Desc.IsMap() {
				g.P(fmt.Sprintf("if x.%s != nil {", fieldName))
				mapFieldInfo := util.MapFieldInfo(g, f, field)
				g.P(fmt.Sprintf("x.%s.RemoveAll(fds.MutState.%sList(%d))", fieldName, mapFieldInfo.KeyAlias, field.Desc.Number()))
				g.P(fmt.Sprintf("}"))
			}
		}
	}

	g.P(fmt.Sprintf("for _,rawFields := range fds.RawFields {"))
	g.P(fmt.Sprintf("if len(rawFields) <= 0 {"))
	g.P("continue")
	g.P("}")
	g.P("rawF := rawFields[0]")
	g.P(fmt.Sprintf("switch rawF.Number {"))
	for _, field := range msg.Fields {
		fieldCapitalName := util.Capital(field.GoName)
		g.P(fmt.Sprintf("case %d:", field.Desc.Number()))
		if (util.IsBuildInType(field.Desc.Kind()) && !field.Desc.IsList() || field.Desc.Kind() == protoreflect.EnumKind || field.Desc.Kind() == protoreflect.BytesKind) && !field.Desc.IsList() {
			goKind := util.ToGoKind(field.Desc.Kind())
			goType, _ := util.FieldGoType(g, f, field)
			if goKind == "float" {
				g.P(fmt.Sprintf("x.Set%s(rawF.F32)", fieldCapitalName))
			} else if goKind == "double" {
				g.P(fmt.Sprintf("x.Set%s(rawF.F64)", fieldCapitalName))
			} else if goKind == "bool" {
				g.P(fmt.Sprintf("x.Set%s(protowire.DecodeBool(rawF.Variant))", fieldCapitalName))
			} else if goKind == "enum" {
				g.P(fmt.Sprintf("x.Set%s(%s(int32(rawF.Variant)))", fieldCapitalName, goType))
			} else if goKind == "string" {
				g.P(fmt.Sprintf("x.Set%s(pbeditor.Bys2Str(rawF.Bytes))", fieldCapitalName))
			} else if goKind == "bytes" {
				g.P(fmt.Sprintf("x.Set%s(rawF.Bytes)", fieldCapitalName))
			} else {
				g.P(fmt.Sprintf("x.Set%s(%s(rawF.Variant))", fieldCapitalName, goType))
			}
		}
		if util.IsMsg(field) {
			g.P(fmt.Sprintf("x.Get%s().MergeDirtyFromBytes(rawF.Bytes)", fieldCapitalName))
		}
		if field.Desc.IsList() {
			goType := util.ToGoKind(field.Desc.Kind())
			listInfo := util.ListFieldInfo(g, f, field)
			if goType == "string" {
				g.P("sl := make([]string,0)")
				g.P("for _,field := range rawFields {")
				g.P("sl = append(sl,pbeditor.Bys2Str(field.Bytes))")
				g.P("}")
				g.P(fmt.Sprintf("x.Get%s().MergeDirtyFromPb(fds.MutState.ListOpLog(rawF.Number),sl)", fieldCapitalName))
			} else if listInfo.IsRef {
				g.P(fmt.Sprintf("sl := make([]%s,0)", listInfo.TypeWithPtr))
				g.P("for _,field := range rawFields {")
				g.P(fmt.Sprintf("tmp := &%s{}", listInfo.Type))
				g.P(fmt.Sprintf("proto.Unmarshal(field.Bytes,tmp)"))
				g.P("sl = append(sl,tmp)")
				g.P("}")
				g.P(fmt.Sprintf("x.Get%s().MergeDirtyFromPb(fds.MutState.ListOpLog(rawF.Number),sl)", fieldCapitalName))
			} else {
				g.P(fmt.Sprintf("tmp := rawF.Bytes"))
				g.P(fmt.Sprintf("dl := make([]%s,0)", listInfo.Type))
				g.P(fmt.Sprintf("for len(tmp) > 0 {"))
				if goType == "float" {
					g.P(fmt.Sprintf("val, n := protowire.ConsumeFixed32(tmp)"))
					g.P(fmt.Sprintf("if n <0 {"))
					g.P(fmt.Sprintf("panic(n)"))
					g.P(fmt.Sprintf("}"))
					g.P(fmt.Sprintf("tmp = tmp[n:]"))
					g.P(fmt.Sprintf("dl = append(dl,math.Float32frombits(val))"))
				} else if goType == "double" {
					g.P(fmt.Sprintf("val, n := protowire.ConsumeFixed64(tmp)"))
					g.P(fmt.Sprintf("if n <0 {"))
					g.P(fmt.Sprintf("panic(n)"))
					g.P(fmt.Sprintf("}"))
					g.P(fmt.Sprintf("tmp = tmp[n:]"))
					g.P(fmt.Sprintf("dl = append(dl,math.Float64frombits(val))"))

				} else if goType == "bool" {
					g.P(fmt.Sprintf("val, n := protowire.ConsumeVarint(tmp)"))
					g.P(fmt.Sprintf("if n <0 {"))
					g.P(fmt.Sprintf("panic(n)"))
					g.P(fmt.Sprintf("}"))
					g.P(fmt.Sprintf("tmp = tmp[n:]"))
					g.P(fmt.Sprintf("dl = append(dl,protowire.DecodeBool(val))"))
				} else {
					g.P(fmt.Sprintf("val, n := protowire.ConsumeVarint(tmp)"))
					g.P(fmt.Sprintf("if n <0 {"))
					g.P(fmt.Sprintf("panic(n)"))
					g.P(fmt.Sprintf("}"))
					g.P(fmt.Sprintf("tmp = tmp[n:]"))
					g.P(fmt.Sprintf("dl = append(dl,%s(val))", listInfo.Type))

				}
				g.P(fmt.Sprintf("}"))
				g.P(fmt.Sprintf("x.Get%s().MergeDirtyFromPb(fds.MutState.ListOpLog(rawF.Number),dl)", fieldCapitalName))
			}
		}
		if field.Desc.IsMap() {
			mapInfo := util.MapFieldInfo(g, f, field)
			g.P(fmt.Sprintf("r := make(map[%s]%s)", mapInfo.KeyType, mapInfo.ValuePointerType))
			g.P("for _,field := range rawFields {")
			g.P("tmp := pbeditor.ParseMap(field.Bytes)")
			g.P(fmt.Sprintf("k := tmp.Key%s()", mapInfo.KeyAlias))
			if mapInfo.IsRef {
				g.P(fmt.Sprintf("v := &%s{}", g.QualifiedGoIdent(mapInfo.ValueMsg.GoIdent)))
				g.P(fmt.Sprintf("proto.Unmarshal(field.Bytes,v)"))
				g.P(fmt.Sprintf("r[k] = v"))
			} else {
				if field.Message.Fields[1].Enum != nil {
					g.P(fmt.Sprintf("v := tmp.ValI32()"))
					g.P(fmt.Sprintf("r[k] = %s(v)", mapInfo.ValueType))
				} else {
					g.P(fmt.Sprintf("v := tmp.Val%s()", mapInfo.ValueAlias))
					g.P(fmt.Sprintf("r[k] = v"))
				}
			}
			g.P("}")
			g.P(fmt.Sprintf("x.Get%s().MergeDirtyFromPb(fds.MutState.%sList(rawF.Number),r)", field.GoName, mapInfo.KeyAlias))
		}
	}
	g.P(fmt.Sprintf("}"))
	g.P(fmt.Sprintf("}"))
	g.P(fmt.Sprintf("return x"))
	g.P(fmt.Sprintf("}"))

	g.P(fmt.Sprintf("func (x *%s) GatherDirtyToPb(r *%s) {", editorMsgName, g.QualifiedGoIdent(msg.GoIdent)))
	if hasArrOrMap {
		g.P("mutState := pbeditor.UnmarshalMutState(r.unknownFields)")
	}
	for _, field := range msg.Fields {
		fieldName := util.EditorFieldName(field.GoName)
		if (util.IsBuildInType(field.Desc.Kind()) || field.Desc.Kind() == protoreflect.EnumKind) && !field.Desc.IsList() {
			g.P(fmt.Sprintf("if x.is%sDirty() {", field.GoName))
			g.P(fmt.Sprintf("r.Set%s(x.%s)", field.GoName, fieldName))
			g.P("}")
		}
		if field.Desc.Kind() == protoreflect.BytesKind {
			g.P(fmt.Sprintf("if x.is%sDirty() {", field.GoName))
			g.P(fmt.Sprintf("r.Set%s(slices.Clone(x.%s))", field.GoName, fieldName))
			g.P("}")
		}
		if util.IsMsg(field) {
			g.P(fmt.Sprintf("if x.is%sDirty()  && x.%s != nil{", field.GoName, fieldName))
			g.P(fmt.Sprintf("if r.%s == nil {", field.GoName))
			g.P(fmt.Sprintf("r.%s = &%s{}", field.GoName, g.QualifiedGoIdent(field.Message.GoIdent)))
			g.P(fmt.Sprintf("}"))
			g.P(fmt.Sprintf("x.%s.GatherDirtyToPb(r.%s)", fieldName, field.GoName))
			g.P(fmt.Sprintf("}"))
		}
		if field.Desc.IsList() {
			listInfo := util.ListFieldInfo(g, f, field)
			g.P(fmt.Sprintf("if x.is%sDirty() {", field.GoName))
			g.P(fmt.Sprintf("ori := r.%s", field.GoName))
			g.P(fmt.Sprintf("lenOf%s := len(ori)", field.GoName))
			g.P(fmt.Sprintf("logsOf%s := x.%s.ChangeLog()", field.GoName, fieldName))
			g.P(fmt.Sprintf("dirtiesOf%s := x.%s.MonitorDirtyList()", field.GoName, fieldName))
			if listInfo.IsRef {
				g.P(fmt.Sprintf("for _,val := range dirtiesOf%s {", field.GoName))
				g.P(fmt.Sprintf("o := &%s{}", listInfo.Type))
				g.P(fmt.Sprintf("val.GatherDirtyToPb(o)"))
				g.P(fmt.Sprintf("r.%s = append(r.%s, o)", field.GoName, field.GoName))
				g.P("}")
				g.P(fmt.Sprintf("for _,log := range logsOf%s {", field.GoName))
				g.P(fmt.Sprintf("log.Pos = int32(lenOf%s) + log.Pos", field.GoName))
				g.P(fmt.Sprintf("mutState.ListAppendLog(%d,log)", field.Desc.Number()))
				g.P("}")
			} else {
				g.P(fmt.Sprintf("for _,val := range dirtiesOf%s {", field.GoName))
				g.P(fmt.Sprintf("r.%s = append(r.%s, val)", field.GoName, field.GoName))
				g.P("}")
				g.P(fmt.Sprintf("for _,log := range logsOf%s {", field.GoName))
				g.P(fmt.Sprintf("log.Pos = int32(lenOf%s) + log.Pos", field.GoName))
				g.P(fmt.Sprintf("mutState.ListAppendLog(%d,log)", field.Desc.Number()))
				g.P("}")
			}
			g.P(fmt.Sprintf("}"))
		}
		if field.Desc.IsMap() {
			mapField := util.MapFieldInfo(g, f, field)
			g.P(fmt.Sprintf("if x.is%sDirty() {", field.GoName))
			g.P(fmt.Sprintf("updated := make([]%s,0)", mapField.KeyType))
			g.P(fmt.Sprintf("if r.%s != nil {", field.GoName))
			g.P(fmt.Sprintf("for k,v := range r.%s {", field.GoName))
			g.P(fmt.Sprintf("if x.%s.ContainDeleted(k) {", fieldName))
			g.P(fmt.Sprintf("delete(r.%s,k)", field.GoName))
			g.P(fmt.Sprintf("continue"))
			g.P(fmt.Sprintf("}"))
			g.P(fmt.Sprintf("if x.%s.ContainDirtied(k) {", fieldName))
			g.P(fmt.Sprintf("updated = append(updated,k)"))
			if mapField.IsRef {
				g.P(fmt.Sprintf("tmp := x.%s.Get(k)", fieldName))
				g.P("if tmp == nil {")
				g.P(fmt.Sprintf("r.%s[k] = nil", field.GoName))
				g.P("continue")
				g.P("}")
				g.P(fmt.Sprintf("if v == nil{"))
				g.P(fmt.Sprintf("v = &%s{}", mapField.ValueType))
				g.P(fmt.Sprintf("tmp.GatherDirtyToPb(v)"))
				g.P(fmt.Sprintf("r.%s[k] = v", field.GoName))
				g.P(fmt.Sprintf(""))
				g.P(fmt.Sprintf("} else {"))
				g.P(fmt.Sprintf("tmp.GatherDirtyToPb(v)"))
				g.P("}")
			} else {
				g.P(fmt.Sprintf("v,_ = x.%s.Get(k)", fieldName))
				g.P(fmt.Sprintf("r.%s[k] = v", field.GoName))
			}

			g.P(fmt.Sprintf("}"))
			g.P(fmt.Sprintf("}"))
			g.P(fmt.Sprintf("} else {"))
			g.P(fmt.Sprintf("r.%s = make(map[%s]%s)", field.GoName, mapField.KeyType, mapField.ValuePointerType))
			g.P(fmt.Sprintf("}"))

			g.P(fmt.Sprintf("for k := range x.%s.Dirtied() {", fieldName))
			g.P(fmt.Sprintf("if !slices.Contains(updated,k) {"))
			if mapField.IsRef {
				g.P(fmt.Sprintf("tmp := x.%s.Get(k)", fieldName))
				g.P("if tmp == nil {")
				g.P(fmt.Sprintf("r.%s[k] = nil", field.GoName))
				g.P("} else {")
				g.P(fmt.Sprintf("v := &%s{}", mapField.ValueType))
				g.P(fmt.Sprintf("tmp.GatherDirtyToPb(v)"))
				g.P(fmt.Sprintf("r.%s[k] = v", field.GoName))
				g.P(fmt.Sprintf("}"))
			} else {
				g.P(fmt.Sprintf("tmp,_ := x.%s.Get(k)", fieldName))
				g.P(fmt.Sprintf("r.%s[k] = tmp", field.GoName))
			}
			g.P(fmt.Sprintf("}"))
			g.P(fmt.Sprintf("}"))

			g.P(fmt.Sprintf("deleted := mutState.%sList(%d)", util.Capital(mapField.KeyAlias), field.Desc.Number()))

			g.P(fmt.Sprintf("for k := range x.%s.Deleted() {", fieldName))
			g.P(fmt.Sprintf("if !slices.Contains(deleted,k) {"))
			g.P(fmt.Sprintf("mutState.Add%s(%d,k)", util.Capital(mapField.KeyAlias), field.Desc.Number()))
			g.P(fmt.Sprintf("}"))
			g.P(fmt.Sprintf("}"))

			g.P(fmt.Sprintf("}"))
		}
	}
	if hasArrOrMap {
		g.P(fmt.Sprintf("r.ProtoReflect().SetUnknown(mutState.MarshalToFieldBytes())"))
	}
	g.P(fmt.Sprintf("}"))

	g.P(fmt.Sprintf("func (x *%s) GatherDirtyToBytes() []byte {", editorMsgName))
	if hasArrOrMap {
		g.P("mutState := pbeditor.CreateMutStateInfo()")
	}
	g.P("var buf []byte")
	for _, field := range msg.Fields {
		fieldName := util.EditorFieldName(field.GoName)
		if (util.IsBuildInType(field.Desc.Kind()) || field.Desc.Kind() == protoreflect.EnumKind) && !field.Desc.IsList() {
			g.P(fmt.Sprintf("if x.is%sDirty() {", field.GoName))
			g.P(fmt.Sprintf("buf = protowire.AppendTag(buf, %d, %+v)", field.Desc.Number(), pbeditor.FieldTypeToWireType(field)))
			switch util.ToGoKind(field.Desc.Kind()) {
			case "int32", "uint32", "int64", "uint64", "enum":
				g.P(fmt.Sprintf("buf = protowire.AppendVarint(buf,uint64(x.%s))", fieldName))
			case "string":
				g.P(fmt.Sprintf("buf = protowire.AppendString(buf,x.%s)", fieldName))
			case "bytes":
				g.P(fmt.Sprintf("buf = protowire.AppendBytes(buf,x.%s)", fieldName))
			case "float":
				g.P(fmt.Sprintf("buf = protowire.AppendFixed32(buf,math.Float32bits(x.%s))", fieldName))
			case "double":
				g.P(fmt.Sprintf("buf = protowire.AppendFixed64(buf,math.Float64bits(x.%s))", fieldName))
			case "bool":
				g.P("var t uint64 = 0")
				g.P(fmt.Sprintf("if x.%s {", fieldName))
				g.P("t = 1")
				g.P("}")
				g.P(fmt.Sprintf("buf = protowire.AppendVarint(buf,t)"))
			}
			g.P("}")
		}
		if util.IsMsg(field) {
			g.P(fmt.Sprintf("if x.is%sDirty() && x.%s != nil {", field.GoName, fieldName))
			g.P(fmt.Sprintf("buf = protowire.AppendTag(buf, %d, protowire.BytesType)", field.Desc.Number()))
			g.P(fmt.Sprintf("buf = protowire.AppendBytes(buf,x.%s.GatherDirtyToBytes())", fieldName))
			g.P("}")
		}
		if field.Desc.IsList() {
			g.P(fmt.Sprintf("if x.is%sDirty() {", field.GoName))
			listInfo := util.ListFieldInfo(g, f, field)
			g.P(fmt.Sprintf("for _,log := range x.%s.ChangeLog() {", fieldName))
			g.P(fmt.Sprintf("mutState.ListAppendLog(%d,log)", field.Desc.Number()))
			g.P("}")
			if listInfo.IsRef {
				g.P(fmt.Sprintf("for _, e := range x.%s.MonitorDirtyList() {", fieldName))
				g.P(fmt.Sprintf("buf = protowire.AppendTag(buf,%d,protowire.BytesType)", field.Desc.Number()))
				g.P(fmt.Sprintf("buf = protowire.AppendBytes(buf,e.GatherDirtyToBytes())"))
				g.P("}")
			} else {
				kind := util.ToGoKind(field.Desc.Kind())
				if kind == "string" {
					g.P(fmt.Sprintf("for _, s := range x.%s.MonitorDirtyList() {", fieldName))
					g.P(fmt.Sprintf("buf = protowire.AppendTag(buf,%d,protowire.BytesType)", field.Desc.Number()))
					g.P(fmt.Sprintf("buf = protowire.AppendString(buf,s)"))
					g.P("}")
				} else if kind == "bool" {
					g.P("tmp := make([]byte, 0)")
					g.P(fmt.Sprintf("for _, b := range x.%s.MonitorDirtyList() {", fieldName))
					g.P("var t uint64 = 0")
					g.P(fmt.Sprintf("if b {"))
					g.P("t = 1")
					g.P("}")
					g.P(fmt.Sprintf("tmp = protowire.AppendVarint(tmp,t)"))
					g.P("}")
					g.P(fmt.Sprintf("buf = protowire.AppendTag(buf,%d,protowire.BytesType)", field.Desc.Number()))
					g.P("buf = protowire.AppendBytes(buf,tmp)")
				} else if kind == "float" {
					g.P("tmp := make([]byte, 0)")
					g.P(fmt.Sprintf("for _, f := range x.%s.MonitorDirtyList() {", fieldName))
					g.P(fmt.Sprintf("tmp = protowire.AppendFixed32(tmp,math.Float32bits(f))"))
					g.P("}")
					g.P(fmt.Sprintf("buf = protowire.AppendTag(buf,%d,protowire.BytesType)", field.Desc.Number()))
					g.P("buf = protowire.AppendBytes(buf,tmp)")

				} else if kind == "double" {
					g.P("tmp := make([]byte, 0)")
					g.P(fmt.Sprintf("for _, f := range x.%s.MonitorDirtyList() {", fieldName))
					g.P(fmt.Sprintf("tmp = protowire.AppendFixed64(tmp,math.Float64bits(f))"))
					g.P("}")
					g.P(fmt.Sprintf("buf = protowire.AppendTag(buf,%d,protowire.BytesType)", field.Desc.Number()))
					g.P("buf = protowire.AppendBytes(buf,tmp)")
				} else {
					g.P("tmp := make([]byte, 0)")
					g.P(fmt.Sprintf("for _, v := range x.%s.MonitorDirtyList() {", fieldName))
					g.P(fmt.Sprintf("tmp = protowire.AppendVarint(tmp,uint64(v))"))
					g.P("}")
					g.P(fmt.Sprintf("buf = protowire.AppendTag(buf,%d,protowire.BytesType)", field.Desc.Number()))
					g.P("buf = protowire.AppendBytes(buf,tmp)")
				}
			}
			g.P("}")
		}

		if field.Desc.IsMap() {
			g.P(fmt.Sprintf("if x.is%sDirty() {", field.GoName))
			mapInfo := util.MapFieldInfo(g, f, field)
			g.P(fmt.Sprintf("if len(x.%s.Deleted()) > 0 {", fieldName))
			g.P(fmt.Sprintf("deleted := x.%s.Deleted()", fieldName))
			g.P(fmt.Sprintf("for v := range deleted {"))
			g.P(fmt.Sprintf("mutState.Add%s(%d,v)", mapInfo.KeyAlias, field.Desc.Number()))
			g.P(fmt.Sprintf("}"))

			g.P(fmt.Sprintf("}"))
			g.P(fmt.Sprintf("if x.%s.Len() > 0 {", fieldName))
			g.P(fmt.Sprintf("for k ,v := range x.%s.Iter() {", fieldName))
			g.P(fmt.Sprintf("if !x.%s.ContainDirtied(k) {", fieldName))
			g.P(fmt.Sprintf("continue"))
			g.P(fmt.Sprintf("}"))

			g.P("var tmp []byte")
			keyTp := util.ToGoKind(field.Message.Fields[0].Desc.Kind())
			if keyTp == "string" {
				g.P(fmt.Sprintf("tmp = protowire.AppendTag(tmp, 1, protowire.BytesType)"))
				g.P(fmt.Sprintf("tmp = protowire.AppendString(tmp, k)"))
			} else if keyTp == "bool" {
				g.P("var t uint64 = 0")
				g.P(fmt.Sprintf("if k {"))
				g.P("t = 1")
				g.P("}")
				g.P(fmt.Sprintf("tmp = protowire.AppendTag(tmp, 1, protowire.VarintType)"))
				g.P(fmt.Sprintf("tmp = protowire.AppendVarint(tmp, t)"))
			} else {
				g.P(fmt.Sprintf("tmp = protowire.AppendTag(tmp, 1, protowire.VarintType)"))
				g.P(fmt.Sprintf("tmp = protowire.AppendVarint(tmp, uint64(k))"))
			}

			if mapInfo.IsRef {
				g.P(fmt.Sprintf("vb := v.GatherDirtyToBytes()"))
				g.P(fmt.Sprintf("tmp = protowire.AppendTag(tmp, 2, protowire.BytesType)"))
				g.P(fmt.Sprintf("tmp = protowire.AppendBytes(tmp, vb)"))
			} else {
				valType := util.ToGoKind(field.Message.Fields[1].Desc.Kind())
				if valType == "string" {
					g.P(fmt.Sprintf("tmp = protowire.AppendTag(tmp, 2, protowire.BytesType)"))
					g.P(fmt.Sprintf("tmp = protowire.AppendString(tmp, v)"))
				} else if valType == "bool" {
					g.P("var t uint64 = 0")
					g.P(fmt.Sprintf("if v {"))
					g.P("t = 1")
					g.P("}")
					g.P(fmt.Sprintf("tmp = protowire.AppendTag(tmp, 2, protowire.VarintType)"))
					g.P(fmt.Sprintf("tmp = protowire.AppendVarint(tmp, t)"))
				} else if valType == "float" {
					g.P(fmt.Sprintf("tmp = protowire.AppendTag(tmp, 2, protowire.Fixed32Type)"))
					g.P(fmt.Sprintf("tmp = protowire.AppendFixed32(tmp,  math.Float32bits(v))"))
				} else if valType == "double" {
					g.P(fmt.Sprintf("tmp = protowire.AppendTag(tmp, 2, protowire.Fixed64Type)"))
					g.P(fmt.Sprintf("tmp = protowire.AppendFixed64(tmp, math.Float64bits(v))"))
				} else {
					g.P(fmt.Sprintf("tmp = protowire.AppendTag(tmp, 2, protowire.VarintType)"))
					g.P(fmt.Sprintf("tmp = protowire.AppendVarint(tmp, uint64(v))"))
				}
			}
			g.P(fmt.Sprintf("buf = protowire.AppendTag(buf,%d,protowire.BytesType)", field.Desc.Number()))
			g.P("buf = protowire.AppendBytes(buf,tmp)")
			g.P(fmt.Sprintf("}"))
			g.P("}")

			g.P("}")
		}
	}
	if hasArrOrMap {
		g.P("buf = append(buf,mutState.MarshalToFieldBytes()...)")
	}
	g.P("return buf")
	g.P("}")
}
