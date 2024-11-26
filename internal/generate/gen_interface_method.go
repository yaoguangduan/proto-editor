package generate

import (
	"fmt"
	"github.com/yaoguangduan/proto-editor/internal/generate/util"
	"google.golang.org/protobuf/compiler/protogen"
)

func GenInterfaceMethod(gen *protogen.Plugin, g *protogen.GeneratedFile, f *util.FileInfo, msg *protogen.Message) {

	editorMsgName := util.EditorMsgName(g.QualifiedGoIdent(msg.GoIdent))
	g.P(fmt.Sprintf("func (x *%s) SetDirty(key any, dirty bool,sync pbeditor.Editor) {", editorMsgName))
	g.P(fmt.Sprintf(""))

	g.P("if dirty {")
	g.P(fmt.Sprintf("x.%s[key.(int)] = struct{}{}", util.DirtyFieldName))
	g.P("x.SetParentDirty()")
	g.P("} else {")
	g.P(fmt.Sprintf("delete(x.%s,key.(int))", util.DirtyFieldName))
	g.P("}")
	g.P("}")

	g.P(fmt.Sprintf("func (x *%s) SetParentDirty() {", util.EditorMsgName(g.QualifiedGoIdent(msg.GoIdent))))
	g.P(fmt.Sprintf("if x.%s != nil {", util.ParentEditorFieldName))
	g.P(fmt.Sprintf("x.%s.SetDirty(x.%s,true,x)", util.ParentEditorFieldName, util.KeyInParentFieldName))
	//g.P(fmt.Sprintf("x.parent.SetParentDirty()"))
	g.P("}")
	g.P("}")

	g.P(fmt.Sprintf("func (x *%s) SetParent(editor pbeditor.Editor, key any) {", editorMsgName))
	g.P(fmt.Sprintf("x.%s = editor", util.ParentEditorFieldName))
	g.P(fmt.Sprintf("x.%s = key", util.KeyInParentFieldName))
	g.P("}")

	g.P(fmt.Sprintf("func (x *%s) FlushDirty(dirty bool) {", editorMsgName))
	for _, field := range msg.Fields {
		g.P(fmt.Sprintf("if dirty || x.is%sDirty() {", util.Capital(field.GoName)))
		g.P(fmt.Sprintf("x.set%sDirty(dirty,true)", util.Capital(field.GoName)))
		g.P("}")
	}
	g.P("}")

	for _, field := range msg.Fields {
		fieldName := util.EditorFieldName(field.GoName)
		g.P(fmt.Sprintf("func (x *%s) set%sDirty(dirty bool,recur bool) {", editorMsgName, util.Capital(field.GoName)))
		g.P(fmt.Sprintf("x.SetDirty(%d,dirty,x)", field.Desc.Number()))
		if util.IsMsg(field) || field.Desc.IsMap() || field.Desc.IsList() {
			g.P(fmt.Sprintf("if recur && x.%s != nil {", fieldName))
			g.P(fmt.Sprintf("x.%s.FlushDirty(dirty)", fieldName))
			g.P("}")
		}

		g.P("}")

		g.P(fmt.Sprintf("func (x *%s) is%sDirty() bool{", editorMsgName, field.GoName))
		g.P(fmt.Sprintf("_,ok := x.%s[%d]", util.DirtyFieldName, field.Desc.Number()))
		g.P(fmt.Sprintf("return ok"))
		g.P("}")
	}

	//Keys
	g.P(fmt.Sprintf("func (x *%s) Key() any {", editorMsgName))
	g.P(fmt.Sprintf("return x.%s", util.KeyInParentFieldName))
	g.P("}")

	g.P(fmt.Sprintf("func (x *%s) SetKey(v any) {", editorMsgName))
	g.P(fmt.Sprintf("x.%s = v", util.KeyInParentFieldName))
	g.P(fmt.Sprintf("}"))
}
