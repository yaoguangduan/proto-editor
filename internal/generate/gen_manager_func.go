package generate

import (
	"fmt"
	"github.com/yaoguangduan/proto-editor/internal/generate/util"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/reflect/protoreflect"
)

func GenClearFunc(gen *protogen.Plugin, g *protogen.GeneratedFile, f *util.FileInfo, msg *protogen.Message) {
	editorMsgName := util.EditorMsgName(g.QualifiedGoIdent(msg.GoIdent))
	g.P(fmt.Sprintf("func (x *%s) Clear() *%s {", editorMsgName, editorMsgName))
	for _, field := range msg.Fields {
		fieldName := util.EditorFieldName(field.GoName)
		if (util.IsBuildInType(field.Desc.Kind()) || field.Desc.Kind() == protoreflect.EnumKind) && !field.Desc.IsList() {
			g.P(fmt.Sprintf("x.Set%s(%v)", field.GoName, util.FieldDefaultVal(g, field)))
		}
		if util.IsMsg(field) {
			g.P(fmt.Sprintf("if x.%s != nil {", fieldName))
			g.P(fmt.Sprintf("x.%s.Clear()", fieldName))
			g.P(fmt.Sprintf("}"))
		}
		if field.Desc.IsList() || field.Desc.IsMap() {
			g.P(fmt.Sprintf("if x.%s != nil {", fieldName))
			g.P(fmt.Sprintf("x.%s.Clear()", fieldName))
			g.P(fmt.Sprintf("}"))
		}
	}
	g.P(fmt.Sprintf("return x"))
	g.P(fmt.Sprintf("}"))
}
