// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.34.2
// 	protoc        v5.27.2
// source: protos_user/user.proto

package protos_user

import (
	protos_animal "github.com/yaoguangduan/proto-editor/pbgen/protos_animal"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type User struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name       *string                         `protobuf:"bytes,1,opt,name=name,proto3,oneof" json:"name,omitempty"`
	Age        *int32                          `protobuf:"varint,2,opt,name=age,proto3,oneof" json:"age,omitempty"`
	Pet        map[int64]*protos_animal.Animal `protobuf:"bytes,3,rep,name=pet,proto3" json:"pet,omitempty" protobuf_key:"varint,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	Friends    []*Friend                       `protobuf:"bytes,4,rep,name=friends,proto3" json:"friends,omitempty"`
	TempChange []float32                       `protobuf:"fixed32,5,rep,packed,name=tempChange,proto3" json:"tempChange,omitempty"`
}

func (x *User) Reset() {
	*x = User{}
	if protoimpl.UnsafeEnabled {
		mi := &file_protos_user_user_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *User) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*User) ProtoMessage() {}

func (x *User) ProtoReflect() protoreflect.Message {
	mi := &file_protos_user_user_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use User.ProtoReflect.Descriptor instead.
func (*User) Descriptor() ([]byte, []int) {
	return file_protos_user_user_proto_rawDescGZIP(), []int{0}
}

func (x *User) GetName() string {
	if x != nil && x.Name != nil {
		return *x.Name
	}
	return ""
}

func (x *User) GetAge() int32 {
	if x != nil && x.Age != nil {
		return *x.Age
	}
	return 0
}

func (x *User) GetPet() map[int64]*protos_animal.Animal {
	if x != nil {
		return x.Pet
	}
	return nil
}

func (x *User) GetFriends() []*Friend {
	if x != nil {
		return x.Friends
	}
	return nil
}

func (x *User) GetTempChange() []float32 {
	if x != nil {
		return x.TempChange
	}
	return nil
}

type Friend struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name     *string  `protobuf:"bytes,1,opt,name=name,proto3,oneof" json:"name,omitempty"`
	Relation *float32 `protobuf:"fixed32,2,opt,name=relation,proto3,oneof" json:"relation,omitempty"`
}

func (x *Friend) Reset() {
	*x = Friend{}
	if protoimpl.UnsafeEnabled {
		mi := &file_protos_user_user_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Friend) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Friend) ProtoMessage() {}

func (x *Friend) ProtoReflect() protoreflect.Message {
	mi := &file_protos_user_user_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Friend.ProtoReflect.Descriptor instead.
func (*Friend) Descriptor() ([]byte, []int) {
	return file_protos_user_user_proto_rawDescGZIP(), []int{1}
}

func (x *Friend) GetName() string {
	if x != nil && x.Name != nil {
		return *x.Name
	}
	return ""
}

func (x *Friend) GetRelation() float32 {
	if x != nil && x.Relation != nil {
		return *x.Relation
	}
	return 0
}

var File_protos_user_user_proto protoreflect.FileDescriptor

var file_protos_user_user_proto_rawDesc = []byte{
	0x0a, 0x16, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x73, 0x5f, 0x75, 0x73, 0x65, 0x72, 0x2f, 0x75, 0x73,
	0x65, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1a, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x73,
	0x5f, 0x61, 0x6e, 0x69, 0x6d, 0x61, 0x6c, 0x2f, 0x61, 0x6e, 0x69, 0x6d, 0x61, 0x6c, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x22, 0xed, 0x01, 0x0a, 0x04, 0x55, 0x73, 0x65, 0x72, 0x12, 0x17, 0x0a,
	0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x48, 0x00, 0x52, 0x04, 0x6e,
	0x61, 0x6d, 0x65, 0x88, 0x01, 0x01, 0x12, 0x15, 0x0a, 0x03, 0x61, 0x67, 0x65, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x05, 0x48, 0x01, 0x52, 0x03, 0x61, 0x67, 0x65, 0x88, 0x01, 0x01, 0x12, 0x20, 0x0a,
	0x03, 0x70, 0x65, 0x74, 0x18, 0x03, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x0e, 0x2e, 0x55, 0x73, 0x65,
	0x72, 0x2e, 0x50, 0x65, 0x74, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x03, 0x70, 0x65, 0x74, 0x12,
	0x21, 0x0a, 0x07, 0x66, 0x72, 0x69, 0x65, 0x6e, 0x64, 0x73, 0x18, 0x04, 0x20, 0x03, 0x28, 0x0b,
	0x32, 0x07, 0x2e, 0x46, 0x72, 0x69, 0x65, 0x6e, 0x64, 0x52, 0x07, 0x66, 0x72, 0x69, 0x65, 0x6e,
	0x64, 0x73, 0x12, 0x1e, 0x0a, 0x0a, 0x74, 0x65, 0x6d, 0x70, 0x43, 0x68, 0x61, 0x6e, 0x67, 0x65,
	0x18, 0x05, 0x20, 0x03, 0x28, 0x02, 0x52, 0x0a, 0x74, 0x65, 0x6d, 0x70, 0x43, 0x68, 0x61, 0x6e,
	0x67, 0x65, 0x1a, 0x3f, 0x0a, 0x08, 0x50, 0x65, 0x74, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x10,
	0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x03, 0x6b, 0x65, 0x79,
	0x12, 0x1d, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x07, 0x2e, 0x41, 0x6e, 0x69, 0x6d, 0x61, 0x6c, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a,
	0x02, 0x38, 0x01, 0x42, 0x07, 0x0a, 0x05, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x42, 0x06, 0x0a, 0x04,
	0x5f, 0x61, 0x67, 0x65, 0x22, 0x58, 0x0a, 0x06, 0x46, 0x72, 0x69, 0x65, 0x6e, 0x64, 0x12, 0x17,
	0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x48, 0x00, 0x52, 0x04,
	0x6e, 0x61, 0x6d, 0x65, 0x88, 0x01, 0x01, 0x12, 0x1f, 0x0a, 0x08, 0x72, 0x65, 0x6c, 0x61, 0x74,
	0x69, 0x6f, 0x6e, 0x18, 0x02, 0x20, 0x01, 0x28, 0x02, 0x48, 0x01, 0x52, 0x08, 0x72, 0x65, 0x6c,
	0x61, 0x74, 0x69, 0x6f, 0x6e, 0x88, 0x01, 0x01, 0x42, 0x07, 0x0a, 0x05, 0x5f, 0x6e, 0x61, 0x6d,
	0x65, 0x42, 0x0b, 0x0a, 0x09, 0x5f, 0x72, 0x65, 0x6c, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x42, 0x38,
	0x5a, 0x36, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x79, 0x61, 0x6f,
	0x67, 0x75, 0x61, 0x6e, 0x67, 0x64, 0x75, 0x61, 0x6e, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2d,
	0x65, 0x64, 0x69, 0x74, 0x6f, 0x72, 0x2f, 0x70, 0x62, 0x67, 0x65, 0x6e, 0x2f, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x73, 0x5f, 0x75, 0x73, 0x65, 0x72, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_protos_user_user_proto_rawDescOnce sync.Once
	file_protos_user_user_proto_rawDescData = file_protos_user_user_proto_rawDesc
)

func file_protos_user_user_proto_rawDescGZIP() []byte {
	file_protos_user_user_proto_rawDescOnce.Do(func() {
		file_protos_user_user_proto_rawDescData = protoimpl.X.CompressGZIP(file_protos_user_user_proto_rawDescData)
	})
	return file_protos_user_user_proto_rawDescData
}

var file_protos_user_user_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_protos_user_user_proto_goTypes = []any{
	(*User)(nil),                 // 0: User
	(*Friend)(nil),               // 1: Friend
	nil,                          // 2: User.PetEntry
	(*protos_animal.Animal)(nil), // 3: Animal
}
var file_protos_user_user_proto_depIdxs = []int32{
	2, // 0: User.pet:type_name -> User.PetEntry
	1, // 1: User.friends:type_name -> Friend
	3, // 2: User.PetEntry.value:type_name -> Animal
	3, // [3:3] is the sub-list for method output_type
	3, // [3:3] is the sub-list for method input_type
	3, // [3:3] is the sub-list for extension type_name
	3, // [3:3] is the sub-list for extension extendee
	0, // [0:3] is the sub-list for field type_name
}

func init() { file_protos_user_user_proto_init() }
func file_protos_user_user_proto_init() {
	if File_protos_user_user_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_protos_user_user_proto_msgTypes[0].Exporter = func(v any, i int) any {
			switch v := v.(*User); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_protos_user_user_proto_msgTypes[1].Exporter = func(v any, i int) any {
			switch v := v.(*Friend); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	file_protos_user_user_proto_msgTypes[0].OneofWrappers = []any{}
	file_protos_user_user_proto_msgTypes[1].OneofWrappers = []any{}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_protos_user_user_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_protos_user_user_proto_goTypes,
		DependencyIndexes: file_protos_user_user_proto_depIdxs,
		MessageInfos:      file_protos_user_user_proto_msgTypes,
	}.Build()
	File_protos_user_user_proto = out.File
	file_protos_user_user_proto_rawDesc = nil
	file_protos_user_user_proto_goTypes = nil
	file_protos_user_user_proto_depIdxs = nil
}
