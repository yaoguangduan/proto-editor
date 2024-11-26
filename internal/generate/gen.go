package generate

import (
	"github.com/samber/lo"
	"github.com/yaoguangduan/proto-editor/internal/generate/util"
	"github.com/yaoguangduan/proto-editor/protodef"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/proto"
	"log"
)

type fileMessageInfo struct {
	file     *protogen.File
	messages []*protogen.Message
}

func Generate(gen *protogen.Plugin) {
	allNeedGenMsg := make(map[string]*protogen.Message)
	for _, file := range gen.Files {
		for _, msg := range file.Messages {
			if proto.HasExtension(msg.Desc.Options(), protodef.E_EditorGen) {
				genIt := proto.GetExtension(msg.Desc.Options(), protodef.E_EditorGen).(bool)
				if genIt {
					findMessageDep(allNeedGenMsg, msg, gen)
				}
			}
		}
	}
	log.Println(allNeedGenMsg)
	for _, files := range gen.Files {
		messages := make([]*protogen.Message, 0)
		for _, msg := range files.Messages {
			if lo.HasKey(allNeedGenMsg, string(msg.Desc.FullName())) {
				messages = append(messages, msg)
			}
		}
		if len(messages) > 0 {
			generateOneFile(gen, files, messages, allNeedGenMsg)
		}
		if len(files.Messages) <= 0 && len(files.Enums) > 0 {
			generateOneFileOnlyEnum(gen, files, messages, allNeedGenMsg)

		}
	}

}
func generateOneFileOnlyEnum(gen *protogen.Plugin, file *protogen.File, messages []*protogen.Message, allNeedGen map[string]*protogen.Message) *protogen.GeneratedFile {
	filename := file.GeneratedFilenamePrefix + ".editor.go"
	g := gen.NewGeneratedFile(filename, file.GoImportPath)
	f := util.NewFileInfo(file)

	g.P("package ", f.GoPackageName)
	g.P()

	g.P("import \"github.com/yaoguangduan/proto-editor/pbeditor\"")
	g.P("import \"iter\"")
	g.P("import \"math\"")
	g.P("import \"slices\"")
	g.P("import \"github.com/yaoguangduan/proto-editor/protodef\"")
	g.P("import \"google.golang.org/protobuf/encoding/protowire\"")
	g.P("import \"google.golang.org/protobuf/proto\"")

	for _, enum := range f.Enums {
		maps, list := allArrayAndMapDependencyOfEnum(g, f, allNeedGen, enum)
		if list != nil {
			GenListEditorForEnum(gen, g, f, enum, list)
		}
		for _, field := range maps {
			GenMapEditorForEnum(gen, g, f, enum, field)
		}

	}

	g.P(`
func unused() {
	_ = math.MinInt
	_ = slices.Min([]int32{})
	_ = iter.Seq[int32](func(yield func(int32) bool) {
		
	})
	_ = protodef.ListOpType_Delete
	_ = pbeditor.U32List{}
	_ = protowire.Number(1)
	_ = proto.Error
}`)
	return g
}
func generateOneFile(gen *protogen.Plugin, file *protogen.File, messages []*protogen.Message, allNeedGen map[string]*protogen.Message) *protogen.GeneratedFile {
	filename := file.GeneratedFilenamePrefix + ".editor.go"
	g := gen.NewGeneratedFile(filename, file.GoImportPath)
	f := util.NewFileInfo(file)

	g.P("package ", f.GoPackageName)
	g.P()

	g.P("import \"github.com/yaoguangduan/proto-editor/pbeditor\"")
	g.P("import \"iter\"")
	g.P("import \"math\"")
	g.P("import \"slices\"")
	g.P("import \"github.com/yaoguangduan/proto-editor/protodef\"")
	g.P("import \"google.golang.org/protobuf/encoding/protowire\"")
	g.P("import \"google.golang.org/protobuf/proto\"")
	for i, imps := 0, f.Desc.Imports(); i < imps.Len(); i++ {
		genImport(gen, g, f, imps.Get(i))
	}
	for _, enum := range f.Enums {
		maps, list := allArrayAndMapDependencyOfEnum(g, f, allNeedGen, enum)
		if list != nil {
			GenListEditorForEnum(gen, g, f, enum, list)
		}
		for _, field := range maps {
			GenMapEditorForEnum(gen, g, f, enum, field)
		}

	}
	for _, msg := range messages {
		maps, list := allArrayAndMapDependency(g, f, allNeedGen, msg)
		if list != nil {
			GenListEditor(gen, g, f, msg, list)
		}
		for _, field := range maps {
			GenMapEditor(gen, g, f, msg, field)
		}
		GenStruct(gen, g, f, msg)
		GenNewMethod(gen, g, f, msg)
		GenInterfaceMethod(gen, g, f, msg)
		GenClearFunc(gen, g, f, msg)
		GenCopyFunc(gen, g, f, msg)
		GenDirtyFunc(gen, g, f, msg)
		GenGetSet(gen, g, f, msg)
		GenOriginalSetFunc(gen, g, f, msg)
	}
	g.P(`
func unused() {
	_ = math.MinInt
	_ = slices.Min([]int32{})
	_ = iter.Seq[int32](func(yield func(int32) bool) {
		
	})
	_ = protodef.ListOpType_Delete
	_ = pbeditor.U32List{}
	_ = protowire.Number(1)
	_ = proto.Error
}`)
	return g
}

func allArrayAndMapDependencyOfEnum(g *protogen.GeneratedFile, f *util.FileInfo, allNeed map[string]*protogen.Message, enum *protogen.Enum) (map[string]*protogen.Field, *protogen.Field) {
	var arrField *protogen.Field
	maps := make(map[string]*protogen.Field)
	for _, msg := range allNeed {
		for _, field := range msg.Fields {
			if field.Desc.IsMap() {
				if field.Message.Fields[1].Enum != nil && field.Message.Fields[1].Enum == enum {
					gt, _ := util.OriFieldGoType(g, f, field.Message.Fields[0])
					maps[gt] = field
				}
			}
			if field.Desc.IsList() && field.Enum != nil && field.Enum == enum {
				arrField = field
			}
		}
	}
	return maps, arrField
}

func allArrayAndMapDependency(g *protogen.GeneratedFile, f *util.FileInfo, allNeed map[string]*protogen.Message, message *protogen.Message) (map[string]*protogen.Field, *protogen.Field) {
	var arrField *protogen.Field
	maps := make(map[string]*protogen.Field)
	for _, msg := range allNeed {
		for _, field := range msg.Fields {
			if field.Desc.IsMap() {
				if field.Message.Fields[1].Message != nil && field.Message.Fields[1].Message == message {
					gt, _ := util.OriFieldGoType(g, f, field.Message.Fields[0])
					maps[gt] = field
				}
			}
			if field.Desc.IsList() && field.Message != nil && field.Message == message {
				arrField = field
			}
		}
	}
	return maps, arrField
}

func findMessageDep(allNeedGenMsg map[string]*protogen.Message, msg *protogen.Message, msgMap *protogen.Plugin) {
	allNeedGenMsg[string(msg.Desc.FullName())] = msg
	for _, field := range msg.Fields {
		if field.Desc.IsMap() {
			mapVal := field.Message.Fields[1]
			if mapVal.Message != nil {
				findMessageDep(allNeedGenMsg, mapVal.Message, msgMap)
			}
		} else {
			if field.Message != nil {
				findMessageDep(allNeedGenMsg, field.Message, msgMap)
			}
		}

	}
}
