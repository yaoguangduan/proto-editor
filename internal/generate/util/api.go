package util

import (
	"fmt"
	"github.com/samber/lo"
	"github.com/yaoguangduan/proto-editor/protodef"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"log"
	"slices"
	"strings"
)

const DirtyFieldName = "_dirties"
const IdxInParentFieldName = "_IdxInParent"
const ParentEditorFieldName = "_parent"
const KeyInParentFieldName = "_keyInParent"

var BuildInTypes = []string{"int32", "uint32", "int64", "uint64", "float", "double", "string", "bool", "bytes"}

func EditorMsgName(name string) string {
	return fmt.Sprintf("%sEditor", name)
}
func EditorFieldIndexName(name string) string {
	return fmt.Sprintf("_%sIdx", strings.ToLower(name[0:1])+name[1:])
}
func Capital(name string) string {
	return lo.Capitalize(name[0:1]) + name[1:]
}

var keywords = map[string]bool{
	"break":       true,
	"default":     true,
	"func":        true,
	"interface":   true,
	"select":      true,
	"case":        true,
	"defer":       true,
	"go":          true,
	"map":         true,
	"struct":      true,
	"chan":        true,
	"else":        true,
	"goto":        true,
	"package":     true,
	"switch":      true,
	"const":       true,
	"fallthrough": true,
	"if":          true,
	"range":       true,
	"type":        true,
	"continue":    true,
	"for":         true,
	"import":      true,
	"return":      true,
	"var":         true,
}

func isKeyword(word string) bool {
	_, ok := keywords[word]
	return ok
}
func EditorFieldName(name string) string {
	nn := strings.ToLower(name[0:1]) + name[1:]
	if isKeyword(nn) {
		nn = nn + "_"
	}
	return nn
}

func IsBuildInType(s protoreflect.Kind) bool {
	return slices.Contains(BuildInTypes, ToGoKind(s))
}

func ToGoKind(desc protoreflect.Kind) string {
	switch desc {
	case protoreflect.BoolKind:
		return "bool"
	case protoreflect.EnumKind:
		return "enum"
	case protoreflect.Int32Kind, protoreflect.Sint32Kind, protoreflect.Sfixed32Kind:
		return "int32"
	case protoreflect.Uint32Kind, protoreflect.Fixed32Kind:
		return "uint32"
	case protoreflect.Int64Kind, protoreflect.Sint64Kind, protoreflect.Sfixed64Kind:
		return "int64"
	case protoreflect.Uint64Kind, protoreflect.Fixed64Kind:
		return "uint64"
	case protoreflect.FloatKind:
		return "float"
	case protoreflect.DoubleKind:
		return "double"
	case protoreflect.StringKind:
		return "string"
	case protoreflect.BytesKind:
		return "bytes"
	case protoreflect.MessageKind:
		return "message"
	}
	panic(fmt.Sprintf("unknown enum kind %v", desc))
}

func IsMsg(field *protogen.Field) bool {
	return !field.Desc.IsList() && !field.Desc.IsMap() && field.Desc.Message() != nil
}

func GetMsgKey(msg *protogen.Message) *protogen.Field {
	opts := msg.Desc.Options()
	if proto.HasExtension(opts, protodef.E_EditorKey) {
		key := proto.GetExtension(opts, protodef.E_EditorKey).(int32)
		fieldDesc := msg.Desc.Fields().ByNumber(protoreflect.FieldNumber(key))
		if fieldDesc == nil {
			return nil
		}
		for _, field := range msg.Fields {
			if field.Desc.FullName() == fieldDesc.FullName() {
				return field
			}
		}
		return nil
	} else {
		return nil
	}
}

func MapFieldValueMsg(field *protogen.Field) *protogen.Message {
	return field.Message.Fields[1].Message
}

func FieldDefaultVal(g *protogen.GeneratedFile, field *protogen.Field) interface{} {
	switch ToGoKind(field.Desc.Kind()) {
	case "int32", "int64", "uint32", "uint64", "float", "double":
		return 0
	case "string":
		return "\"\""
	case "bool":
		return false
	case "bytes":
		return "make([]byte,0)"
	case "enum":
		log.Println(field.Message)
		return g.QualifiedGoIdent(field.Enum.Values[0].GoIdent)
	}
	panic(fmt.Sprintf("invalid kind:%s", field.Desc.Kind()))
}

func ConstructString(g *protogen.GeneratedFile, f *FileInfo, m *protogen.Message) string {
	name := g.QualifiedGoIdent(m.GoIdent)
	idx := strings.Index(name, ".")
	if idx == -1 {
		return fmt.Sprintf("New%s()", EditorMsgName(name))
	} else {
		return name[0:idx] + ".New()" + EditorMsgName(name[idx+1:])
	}
}

type ListFieldData struct {
	Msg               *protogen.Message
	Type              string
	TypeAlias         string
	TypeWithPtr       string
	TypeName          string
	TypeConsName      string
	IsRef             bool
	ConsName          string
	EditorName        string
	EditorNameWithPtr string
}

func ListFieldInfo(g *protogen.GeneratedFile, f *FileInfo, field *protogen.Field) *ListFieldData {
	info := &ListFieldData{Msg: field.Message}
	tp, _ := OriFieldGoType(g, f, field)
	if strings.HasPrefix(tp, "[]") {
		tp = tp[2:]
	}
	info.Type = tp
	info.TypeWithPtr = tp
	if strings.HasPrefix(tp, "*") {
		info.Type = tp[1:]
	}
	info.TypeName = info.Type
	idx := strings.Index(tp, ".")
	if idx != -1 {
		info.TypeName = tp[idx+1:]
	}
	alias, ok := mapValName[tp]
	info.TypeAlias = tp
	if ok {
		info.TypeAlias = alias
		info.ConsName = fmt.Sprintf("pbeditor.New%sList()", alias)
		info.EditorName = info.Type
		info.EditorNameWithPtr = info.Type
	} else {
		info.IsRef = field.Enum == nil
		if idx == -1 {
			info.ConsName = fmt.Sprintf("New%sEditorList()", info.TypeName)
			info.TypeConsName = fmt.Sprintf("New%sEditor()", info.TypeName)
		} else {
			info.ConsName = fmt.Sprintf("%s.New%sEditorList()", tp[0:idx], info.TypeName)
			info.TypeConsName = fmt.Sprintf("%s.New%sEditor()", tp[0:idx], info.TypeName)
		}
		info.EditorName = info.Type + "Editor"
		info.EditorNameWithPtr = "*" + info.Type + "Editor"
		if !info.IsRef {
			info.EditorName = info.Type
			info.EditorNameWithPtr = info.Type
		}

	}

	return info
}

type MapFieldData struct {
	KeyMsg               *protogen.Message
	KeyType              string
	KeyAlias             string
	ValueMsg             *protogen.Message
	ValueType            string
	ValueAlias           string
	ValuePointerType     string
	ValueEditorType      string
	ConstructName        string
	IsRef                bool
	ValueEditorTypeNoPtr string
	ValueConstructName   string
}

func MapFieldInfo(g *protogen.GeneratedFile, f *FileInfo, field *protogen.Field) *MapFieldData {
	if !field.Desc.IsMap() {
		return nil
	}
	kf := field.Message.Fields[0]
	vf := field.Message.Fields[1]
	mfd := &MapFieldData{}
	mfd.KeyMsg = kf.Message
	mfd.KeyType, _ = OriFieldGoType(g, f, kf)
	mfd.KeyAlias = mapKeyName[mfd.KeyType]
	mfd.ValueMsg = vf.Message
	vt, _ := OriFieldGoType(g, f, vf)
	mfd.ValuePointerType = vt
	if strings.HasPrefix(mfd.ValuePointerType, "*") {
		mfd.ValueType = mfd.ValuePointerType[1:]
	} else {
		mfd.ValueType = mfd.ValuePointerType
	}
	if vf.Message != nil {
		mfd.ValueEditorType = "*" + EditorMsgName(g.QualifiedGoIdent(vf.Message.GoIdent))
		mfd.ValueEditorTypeNoPtr = EditorMsgName(g.QualifiedGoIdent(vf.Message.GoIdent))
		mfd.IsRef = true
	} else {
		mfd.ValueEditorType = mfd.ValueType
	}
	valName, ok := mapValName[mfd.ValueEditorType]
	if ok {
		mfd.ConstructName = fmt.Sprintf("pbeditor.New%s%sMap()", mapKeyName[mfd.KeyType], valName)
		mfd.ValueAlias = valName
	} else {

		fullName := ""
		if field.Message.Fields[1].Message != nil {
			fullName = g.QualifiedGoIdent(field.Message.Fields[1].Message.GoIdent)
		} else {
			fullName = g.QualifiedGoIdent(field.Message.Fields[1].Enum.GoIdent)
		}
		pkgName := ""
		msgName := fullName
		if strings.Contains(fullName, ".") {
			pkgName = fullName[0:strings.Index(fullName, ".")]
			msgName = fullName[strings.Index(fullName, ".")+1:]
		}
		if pkgName == "" {
			mfd.ConstructName = fmt.Sprintf("New%s%sEditorMap()", mapKeyName[mfd.KeyType], msgName)
			mfd.ValueConstructName = fmt.Sprintf("New%sEditor()", msgName)
		} else {
			mfd.ConstructName = fmt.Sprintf("%s.New%s%sEditorMap()", pkgName, mapKeyName[mfd.KeyType], msgName)
			mfd.ValueConstructName = fmt.Sprintf("%s.New%sEditor()", pkgName, msgName)
		}
	}
	return mfd
}

func OriFieldGoType(g *protogen.GeneratedFile, f *FileInfo, field *protogen.Field) (goType string, pointer bool) {
	if field.Desc.IsWeak() {
		return "struct{}", false
	}

	pointer = field.Desc.HasPresence()
	switch field.Desc.Kind() {
	case protoreflect.BoolKind:
		goType = "bool"
	case protoreflect.EnumKind:
		goType = g.QualifiedGoIdent(field.Enum.GoIdent)
	case protoreflect.Int32Kind, protoreflect.Sint32Kind, protoreflect.Sfixed32Kind:
		goType = "int32"
	case protoreflect.Uint32Kind, protoreflect.Fixed32Kind:
		goType = "uint32"
	case protoreflect.Int64Kind, protoreflect.Sint64Kind, protoreflect.Sfixed64Kind:
		goType = "int64"
	case protoreflect.Uint64Kind, protoreflect.Fixed64Kind:
		goType = "uint64"
	case protoreflect.FloatKind:
		goType = "float32"
	case protoreflect.DoubleKind:
		goType = "float64"
	case protoreflect.StringKind:
		goType = "string"
	case protoreflect.BytesKind:
		goType = "[]byte"
		pointer = false // rely on nullability of slices for presence
	case protoreflect.MessageKind, protoreflect.GroupKind:
		goType = "*" + g.QualifiedGoIdent(field.Message.GoIdent)
		pointer = false // pointer captured as part of the type
	}
	switch {
	case field.Desc.IsList():
		return "[]" + goType, false
	case field.Desc.IsMap():
		keyType, _ := OriFieldGoType(g, f, field.Message.Fields[0])
		valType, _ := OriFieldGoType(g, f, field.Message.Fields[1])
		return fmt.Sprintf("map[%v]%v", keyType, valType), false
	}
	return goType, pointer
}

// FieldGoType returns the Go type used for a field.
//
// If it returns pointer=true, the struct field is a pointer to the type.
func FieldGoType(g *protogen.GeneratedFile, f *FileInfo, field *protogen.Field) (goType string, pointer bool) {
	if field.Desc.IsWeak() {
		return "struct{}", false
	}
	pointer = field.Desc.HasPresence()
	switch field.Desc.Kind() {
	case protoreflect.BoolKind:
		goType = "bool"
	case protoreflect.EnumKind:
		goType = g.QualifiedGoIdent(field.Enum.GoIdent)
	case protoreflect.Int32Kind, protoreflect.Sint32Kind, protoreflect.Sfixed32Kind:
		goType = "int32"
	case protoreflect.Uint32Kind, protoreflect.Fixed32Kind:
		goType = "uint32"
	case protoreflect.Int64Kind, protoreflect.Sint64Kind, protoreflect.Sfixed64Kind:
		goType = "int64"
	case protoreflect.Uint64Kind, protoreflect.Fixed64Kind:
		goType = "uint64"
	case protoreflect.FloatKind:
		goType = "float32"
	case protoreflect.DoubleKind:
		goType = "float64"
	case protoreflect.StringKind:
		goType = "string"
	case protoreflect.BytesKind:
		goType = "[]byte"
		pointer = false // rely on nullability of slices for presence
	case protoreflect.MessageKind, protoreflect.GroupKind:
		goType = "*" + EditorMsgName(g.QualifiedGoIdent(field.Message.GoIdent))
		pointer = false // pointer captured as part of the type
	}
	switch {
	case field.Desc.IsList():
		keyName, ok := mapValName[goType]
		if ok {
			return fmt.Sprintf("*pbeditor.%sList", keyName), false
		} else {
			log.Println(goType, field.Desc)
			fullName := ""
			if field.Message != nil {
				fullName = g.QualifiedGoIdent(field.Message.GoIdent)
			} else {
				fullName = g.QualifiedGoIdent(field.Enum.GoIdent)
			}
			pkgName := ""
			msgName := fullName
			if strings.Contains(fullName, ".") {
				pkgName = fullName[0:strings.Index(fullName, ".")]
				msgName = fullName[strings.Index(fullName, ".")+1:]
			}
			if pkgName == "" {
				return fmt.Sprintf("*%sEditorList", msgName), false
			} else {
				return fmt.Sprintf("*%s.%sEditorList", pkgName, msgName), false
			}
		}
	case field.Desc.IsMap():
		keyType, _ := FieldGoType(g, f, field.Message.Fields[0])
		valType, _ := FieldGoType(g, f, field.Message.Fields[1])
		valName, ok := mapValName[valType]
		if ok {
			return fmt.Sprintf("*pbeditor.%s%sMap", mapKeyName[keyType], valName), false
		} else {
			var fullName string
			if field.Message.Fields[1].Message != nil {
				fullName = g.QualifiedGoIdent(field.Message.Fields[1].Message.GoIdent)
			} else {
				fullName = g.QualifiedGoIdent(field.Message.Fields[1].Enum.GoIdent)
			}
			pkgName := ""
			msgName := fullName
			if strings.Contains(fullName, ".") {
				pkgName = fullName[0:strings.Index(fullName, ".")]
				msgName = fullName[strings.Index(fullName, ".")+1:]
			}
			if pkgName == "" {
				return fmt.Sprintf("*%s%sEditorMap", mapKeyName[keyType], msgName), false
			} else {
				return fmt.Sprintf("*%s.%s%sEditorMap", pkgName, mapKeyName[keyType], msgName), false
			}
		}
	}
	return goType, pointer
}

var mapKeyName = map[string]string{
	"uint32": "U32",
	"uint64": "U64",
	"int32":  "I32",
	"int64":  "I64",
	"string": "Str",
	"bool":   "Bool",
}
var mapValName = map[string]string{
	"uint32":  "U32",
	"uint64":  "U64",
	"int32":   "I32",
	"int64":   "I64",
	"string":  "Str",
	"bool":    "Bool",
	"float32": "F32",
	"float64": "F64",
	"[]byte":  "Bytes",
}
