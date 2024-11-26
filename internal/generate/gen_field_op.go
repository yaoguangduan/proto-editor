package generate

import (
	"fmt"
	"github.com/yaoguangduan/proto-editor/internal/generate/util"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/reflect/protoreflect"
	"log"
	"strings"
)

func GenGetSet(gen *protogen.Plugin, g *protogen.GeneratedFile, f *util.FileInfo, msg *protogen.Message) {
	editorMsgName := util.EditorMsgName(g.QualifiedGoIdent(msg.GoIdent))
	for _, field := range msg.Fields {
		//GetXXX
		fieldType, _ := util.FieldGoType(g, f, field)
		fieldName := util.EditorFieldName(field.GoName)
		g.P(fmt.Sprintf("func (x *%s) Get%s() %s {", editorMsgName, field.GoName, fieldType))
		if util.IsMsg(field) {
			g.P(fmt.Sprintf("if x.%s == nil {", fieldName))
			fullName := util.EditorMsgName(g.QualifiedGoIdent(field.Message.GoIdent))
			idx := strings.LastIndex(fullName, ".")
			if idx == -1 {
				g.P(fmt.Sprintf("x.%s = New%s()", fieldName, fullName))
			} else {
				g.P(fmt.Sprintf("x.%s = %s.New%s()", fieldName, fullName[0:idx], fullName[idx+1:]))
			}
			g.P(fmt.Sprintf("x.%s.SetParent(x,%d)", fieldName, field.Desc.Number()))
			g.P("}")
		}
		if field.Desc.IsMap() {
			mapFieldInfo := util.MapFieldInfo(g, f, field)
			g.P(fmt.Sprintf("if x.%s == nil {", fieldName))
			g.P(fmt.Sprintf("x.%s = %s", fieldName, mapFieldInfo.ConstructName))
			g.P(fmt.Sprintf("x.%s.SetParent(x,%d)", fieldName, field.Desc.Number()))
			g.P("}")
		}
		if field.Desc.IsList() {
			fieldInfo := util.ListFieldInfo(g, f, field)
			g.P(fmt.Sprintf("if x.%s == nil {", fieldName))
			g.P(fmt.Sprintf("x.%s = %s", fieldName, fieldInfo.ConsName))
			g.P(fmt.Sprintf("x.%s.SetParent(x,%d)", fieldName, field.Desc.Number()))
			g.P("}")
		}
		g.P(fmt.Sprintf("return x.%s", fieldName))
		g.P("}")

		//SetXXX
		if !field.Desc.IsList() && !field.Desc.IsMap() {
			if (util.IsBuildInType(field.Desc.Kind()) || field.Desc.Kind() == protoreflect.EnumKind) && field.Desc.Kind() != protoreflect.BytesKind {
				g.P(fmt.Sprintf("func (x *%s) Set%s(v %s) *%s{", editorMsgName, field.GoName, fieldType, editorMsgName))
				g.P(fmt.Sprintf("if x.%s == v {", fieldName))
				g.P(fmt.Sprintf("return x"))
				g.P(fmt.Sprintf("}"))
			} else if field.Desc.Kind() == protoreflect.BytesKind {
				g.P(fmt.Sprintf("func (x *%s) Set%s(v %s) *%s{", editorMsgName, field.GoName, fieldType, editorMsgName))
			} else {
				fullName := util.EditorMsgName(g.QualifiedGoIdent(field.Message.GoIdent))
				if util.IsMsg(field) {
					fullName = "*" + fullName
				}
				g.P(fmt.Sprintf("func (x *%s) Set%s(v %s) *%s{", editorMsgName, field.GoName, fullName, editorMsgName))
			}

			if util.IsMsg(field) {
				g.P(fmt.Sprintf("if v != nil {"))
				g.P(fmt.Sprintf("v.SetParent(x,%d)", field.Desc.Number()))
				g.P("}")
				g.P(fmt.Sprintf("if x.%s != nil {", fieldName))
				g.P(fmt.Sprintf("x.%s.SetParent(nil,-1)", fieldName))
				g.P("}")
			}
			g.P(fmt.Sprintf("x.%s = v", fieldName))
			g.P(fmt.Sprintf("x.set%sDirty(true,false)", field.GoName))
			g.P(fmt.Sprintf("return x"))

			g.P("}")
		}

		if field.Desc.IsList() {
			fieldInfo := util.ListFieldInfo(g, f, field)
			log.Println(fieldInfo)
			g.P(fmt.Sprintf("func (x *%s) Add%s(v %s) *%s{", editorMsgName, field.GoName, fieldInfo.EditorNameWithPtr, editorMsgName))
			g.P(fmt.Sprintf("if x.%s == nil {", fieldName))
			g.P(fmt.Sprintf("x.%s = %s", fieldName, fieldInfo.ConsName))
			g.P(fmt.Sprintf("x.%s.SetParent(x,%d)", fieldName, field.Desc.Number()))
			g.P("}")
			g.P(fmt.Sprintf("x.%s.Add(v)", fieldName))
			g.P(fmt.Sprintf("return x"))
			g.P("}")
		}
	}
}
