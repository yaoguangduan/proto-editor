package generate

import (
	"fmt"
	"github.com/yaoguangduan/proto-editor/internal/generate/util"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/reflect/protoreflect"
)

func GenCopyFunc(gen *protogen.Plugin, g *protogen.GeneratedFile, f *util.FileInfo, msg *protogen.Message) {
	editorMsgName := util.EditorMsgName(g.QualifiedGoIdent(msg.GoIdent))
	g.P(fmt.Sprintf("func (x *%s) CopyTo(r *%s) *%s {", editorMsgName, g.QualifiedGoIdent(msg.GoIdent), editorMsgName))
	for _, field := range msg.Fields {
		fieldName := util.EditorFieldName(field.GoName)
		if (util.IsBuildInType(field.Desc.Kind()) || field.Desc.Kind() == protoreflect.EnumKind) && !field.Desc.IsList() {
			g.P(fmt.Sprintf("r.Set%s(x.%s)", field.GoName, fieldName))
		}
		if field.Desc.Kind() == protoreflect.BytesKind {
			g.P(fmt.Sprintf("r.Set%s(slices.Clone(x.%s))", field.GoName, fieldName))
		}
		if util.IsMsg(field) {
			g.P(fmt.Sprintf("if x.%s != nil {", fieldName))
			g.P(fmt.Sprintf("tmp := &%s{}", g.QualifiedGoIdent(field.Message.GoIdent)))
			g.P(fmt.Sprintf("x.%s.CopyTo(tmp)", fieldName))
			g.P(fmt.Sprintf("r.Set%s(tmp)", field.GoName))
			g.P(fmt.Sprintf("}"))
		}
		if field.Desc.IsList() {
			g.P(fmt.Sprintf("if x.%s != nil && x.%s.Len() > 0 {", fieldName, fieldName))
			listInfo := util.ListFieldInfo(g, f, field)
			if listInfo.IsRef {
				g.P(fmt.Sprintf(`
					for _, v := range x.%s.Iter() {
						if v == nil {
							r.Add%s(nil)
						} else {
							tmp := &%s{}
							v.CopyTo(tmp)
							r.Add%s(tmp)
						}
					}
					`, fieldName, field.GoName, listInfo.TypeName, field.GoName))
			} else {
				g.P(fmt.Sprintf("r.Set%s(x.%s.Values())", field.GoName, fieldName))
			}
			g.P(fmt.Sprintf("}"))
		}
		if field.Desc.IsMap() {
			mapFieldInfo := util.MapFieldInfo(g, f, field)
			g.P(fmt.Sprintf("if x.%s != nil && x.%s.Len() > 0 {", fieldName, fieldName))
			g.P(fmt.Sprintf("tmp := make(map[%s]%s)", mapFieldInfo.KeyType, mapFieldInfo.ValuePointerType))
			g.P(fmt.Sprintf("for k,v := range x.%s.Iter() {", fieldName))
			if mapFieldInfo.IsRef {
				g.P(fmt.Sprintf("tmpV := &%s{}", mapFieldInfo.ValueType))
				g.P(fmt.Sprintf("v.CopyTo(tmpV)"))
				g.P(fmt.Sprintf("tmp[k] = tmpV"))
			} else {
				g.P(fmt.Sprintf("tmp[k] = v"))
			}
			g.P(fmt.Sprintf("}"))
			g.P(fmt.Sprintf("r.Set%s(tmp)", field.GoName))
			g.P(fmt.Sprintf("}"))
		}
	}
	g.P(fmt.Sprintf("return x"))
	g.P(fmt.Sprintf("}"))

	g.P(fmt.Sprintf("func (x *%s) CopyFrom(r *%s) *%s{", editorMsgName, g.QualifiedGoIdent(msg.GoIdent), editorMsgName))
	for _, field := range msg.Fields {
		fieldName := util.Capital(field.GoName)
		if (util.IsBuildInType(field.Desc.Kind()) || field.Desc.Kind() == protoreflect.EnumKind) && !field.Desc.IsList() && field.Desc.Kind() != protoreflect.BytesKind {
			g.P(fmt.Sprintf("if r.%s != nil {", fieldName))
			g.P(fmt.Sprintf("x.Set%s(*r.%s)", field.GoName, fieldName))
			g.P(fmt.Sprintf("}"))
		}
		if field.Desc.Kind() == protoreflect.BytesKind {
			g.P(fmt.Sprintf("if len(r.%s) > 0 {", field.GoName))
			g.P(fmt.Sprintf("x.Set%s(slices.Clone(r.%s))", field.GoName, field.GoName))
			g.P(fmt.Sprintf("}"))
		}
		if util.IsMsg(field) {
			g.P(fmt.Sprintf("if r.%s != nil {", field.GoName))
			g.P(fmt.Sprintf("x.Get%s().CopyFrom(r.%s)", field.GoName, field.GoName))
			g.P(fmt.Sprintf("}"))
		}
		if field.Desc.IsList() {
			g.P(fmt.Sprintf("if r.%s != nil {", field.GoName))
			listInfo := util.ListFieldInfo(g, f, field)
			if listInfo.IsRef {
				g.P(fmt.Sprintf(`
					for _, v := range r.%s {
						if v == nil {
							x.Get%s().Add(nil)
						} else {
							tmp := %s
							tmp.CopyFrom(v)
							x.Get%s().Add(tmp)
						}
					}
					`, field.GoName, field.GoName, listInfo.TypeConsName, field.GoName))
			} else {
				g.P(fmt.Sprintf("x.Get%s().AddAll(r.%s)", fieldName, fieldName))
			}
			g.P(fmt.Sprintf("}"))
		}
		if field.Desc.IsMap() {
			mapFieldInfo := util.MapFieldInfo(g, f, field)
			g.P(fmt.Sprintf("for k,v := range r.%s {", fieldName))
			if mapFieldInfo.IsRef {
				g.P(fmt.Sprintf("if v != nil {"))
				g.P(fmt.Sprintf("vv := %s", mapFieldInfo.ValueConstructName))
				g.P(fmt.Sprintf("vv.CopyFrom(v)"))
				g.P(fmt.Sprintf("x.Get%s().Put(k,vv)", fieldName))
				g.P(fmt.Sprintf("}"))
			} else {
				g.P(fmt.Sprintf("x.Get%s().Put(k,v)", fieldName))
			}
			g.P(fmt.Sprintf("}"))
		}
	}
	g.P(fmt.Sprintf("return x"))
	g.P(fmt.Sprintf("}"))

}
