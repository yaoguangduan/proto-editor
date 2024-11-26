package generate

import (
	"fmt"
	"github.com/yaoguangduan/proto-editor/internal/generate/util"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/reflect/protoreflect"
)

func GenOriginalSetFunc(gen *protogen.Plugin, g *protogen.GeneratedFile, f *util.FileInfo, msg *protogen.Message) {
	msgName := g.QualifiedGoIdent(msg.GoIdent)
	for _, field := range msg.Fields {
		gt, _ := util.OriFieldGoType(g, f, field)
		if (util.IsBuildInType(field.Desc.Kind()) || field.Desc.Kind() == protoreflect.EnumKind) && !field.Desc.IsList() && field.Desc.Kind() != protoreflect.BytesKind {
			g.P(fmt.Sprintf("func (xs *%s) Set%s(v %s) {", msgName, field.GoName, gt))
			g.P(fmt.Sprintf("xs.%s = &v", field.GoName))
			g.P(fmt.Sprintf("}"))
		} else {
			g.P(fmt.Sprintf("func (xs *%s) Set%s(v %s) {", msgName, field.GoName, gt))
			g.P(fmt.Sprintf("xs.%s = v", field.GoName))
			g.P(fmt.Sprintf("}"))
		}
		if field.Desc.IsList() {
			fieldInfo := util.ListFieldInfo(g, f, field)
			g.P(fmt.Sprintf("func (xs *%s) Add%s(v %s) {", msgName, field.GoName, fieldInfo.TypeWithPtr))
			g.P(fmt.Sprintf("xs.%s = append(xs.%s,v)", field.GoName, field.GoName))
			g.P(fmt.Sprintf("}"))
		}
	}
}
