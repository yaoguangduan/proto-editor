package generate

import (
	"fmt"
	"github.com/yaoguangduan/proto-editor/internal/generate/util"
	"google.golang.org/protobuf/compiler/protogen"
)

func GenStruct(gen *protogen.Plugin, g *protogen.GeneratedFile, f *util.FileInfo, msg *protogen.Message) {
	g.P("type ", util.EditorMsgName(g.QualifiedGoIdent(msg.GoIdent)), " struct {")

	for _, field := range msg.Fields {
		goType, _ := util.FieldGoType(g, f, field)
		g.P(util.EditorFieldName(field.GoName), " ", goType)
	}
	g.P(util.DirtyFieldName, " map[int]struct{}")
	g.P(util.ParentEditorFieldName, " pbeditor.Editor")
	g.P(util.KeyInParentFieldName, " any")
	g.P("}")
}
func GenNewMethod(gen *protogen.Plugin, g *protogen.GeneratedFile, f *util.FileInfo, msg *protogen.Message) {
	g.P("func New", util.EditorMsgName(g.QualifiedGoIdent(msg.GoIdent)), " () *", util.EditorMsgName(g.QualifiedGoIdent(msg.GoIdent)), " {")

	g.P("return &", util.EditorMsgName(g.QualifiedGoIdent(msg.GoIdent)), "{")
	g.P(fmt.Sprintf("%s : make(map[int]struct{}),", util.DirtyFieldName))
	g.P("}")

	g.P("}")
}
