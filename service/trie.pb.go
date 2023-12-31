// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.31.0
// 	protoc        v3.19.6
// source: trie.proto

package main

import (
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

type TrieNode struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	IsEndOfWord bool                 `protobuf:"varint,2,opt,name=IsEndOfWord,proto3" json:"IsEndOfWord,omitempty"`
	Children    map[string]*TrieNode `protobuf:"bytes,3,rep,name=Children,proto3" json:"Children,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
}

func (x *TrieNode) Reset() {
	*x = TrieNode{}
	if protoimpl.UnsafeEnabled {
		mi := &file_trie_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TrieNode) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TrieNode) ProtoMessage() {}

func (x *TrieNode) ProtoReflect() protoreflect.Message {
	mi := &file_trie_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TrieNode.ProtoReflect.Descriptor instead.
func (*TrieNode) Descriptor() ([]byte, []int) {
	return file_trie_proto_rawDescGZIP(), []int{0}
}

func (x *TrieNode) GetIsEndOfWord() bool {
	if x != nil {
		return x.IsEndOfWord
	}
	return false
}

func (x *TrieNode) GetChildren() map[string]*TrieNode {
	if x != nil {
		return x.Children
	}
	return nil
}

type HybridTrie struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Root TrieNode `protobuf:"bytes,1,opt,name=Root,proto3" json:"Root,omitempty"`
}

func (x *HybridTrie) Reset() {
	*x = HybridTrie{}
	if protoimpl.UnsafeEnabled {
		mi := &file_trie_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *HybridTrie) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*HybridTrie) ProtoMessage() {}

func (x *HybridTrie) ProtoReflect() protoreflect.Message {
	mi := &file_trie_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use HybridTrie.ProtoReflect.Descriptor instead.
func (*HybridTrie) Descriptor() ([]byte, []int) {
	return file_trie_proto_rawDescGZIP(), []int{1}
}

func (x *HybridTrie) GetRoot() *TrieNode {
	if x != nil {
		return &x.Root
	}
	return nil
}

var File_trie_proto protoreflect.FileDescriptor

var file_trie_proto_rawDesc = []byte{
	0x0a, 0x0a, 0x74, 0x72, 0x69, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xa9, 0x01, 0x0a,
	0x08, 0x54, 0x72, 0x69, 0x65, 0x4e, 0x6f, 0x64, 0x65, 0x12, 0x20, 0x0a, 0x0b, 0x49, 0x73, 0x45,
	0x6e, 0x64, 0x4f, 0x66, 0x57, 0x6f, 0x72, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x08, 0x52, 0x0b,
	0x49, 0x73, 0x45, 0x6e, 0x64, 0x4f, 0x66, 0x57, 0x6f, 0x72, 0x64, 0x12, 0x33, 0x0a, 0x08, 0x43,
	0x68, 0x69, 0x6c, 0x64, 0x72, 0x65, 0x6e, 0x18, 0x03, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x17, 0x2e,
	0x54, 0x72, 0x69, 0x65, 0x4e, 0x6f, 0x64, 0x65, 0x2e, 0x43, 0x68, 0x69, 0x6c, 0x64, 0x72, 0x65,
	0x6e, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x08, 0x43, 0x68, 0x69, 0x6c, 0x64, 0x72, 0x65, 0x6e,
	0x1a, 0x46, 0x0a, 0x0d, 0x43, 0x68, 0x69, 0x6c, 0x64, 0x72, 0x65, 0x6e, 0x45, 0x6e, 0x74, 0x72,
	0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03,
	0x6b, 0x65, 0x79, 0x12, 0x1f, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x09, 0x2e, 0x54, 0x72, 0x69, 0x65, 0x4e, 0x6f, 0x64, 0x65, 0x52, 0x05, 0x76,
	0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38, 0x01, 0x22, 0x2b, 0x0a, 0x0a, 0x48, 0x79, 0x62, 0x72,
	0x69, 0x64, 0x54, 0x72, 0x69, 0x65, 0x12, 0x1d, 0x0a, 0x04, 0x52, 0x6f, 0x6f, 0x74, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x09, 0x2e, 0x54, 0x72, 0x69, 0x65, 0x4e, 0x6f, 0x64, 0x65, 0x52,
	0x04, 0x52, 0x6f, 0x6f, 0x74, 0x42, 0x03, 0x5a, 0x01, 0x2e, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x33,
}

var (
	file_trie_proto_rawDescOnce sync.Once
	file_trie_proto_rawDescData = file_trie_proto_rawDesc
)

func file_trie_proto_rawDescGZIP() []byte {
	file_trie_proto_rawDescOnce.Do(func() {
		file_trie_proto_rawDescData = protoimpl.X.CompressGZIP(file_trie_proto_rawDescData)
	})
	return file_trie_proto_rawDescData
}

var file_trie_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_trie_proto_goTypes = []interface{}{
	(*TrieNode)(nil),   // 0: TrieNode
	(*HybridTrie)(nil), // 1: HybridTrie
	nil,                // 2: TrieNode.ChildrenEntry
}
var file_trie_proto_depIdxs = []int32{
	2, // 0: TrieNode.Children:type_name -> TrieNode.ChildrenEntry
	0, // 1: HybridTrie.Root:type_name -> TrieNode
	0, // 2: TrieNode.ChildrenEntry.value:type_name -> TrieNode
	3, // [3:3] is the sub-list for method output_type
	3, // [3:3] is the sub-list for method input_type
	3, // [3:3] is the sub-list for extension type_name
	3, // [3:3] is the sub-list for extension extendee
	0, // [0:3] is the sub-list for field type_name
}

func init() { file_trie_proto_init() }
func file_trie_proto_init() {
	if File_trie_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_trie_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*TrieNode); i {
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
		file_trie_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*HybridTrie); i {
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
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_trie_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_trie_proto_goTypes,
		DependencyIndexes: file_trie_proto_depIdxs,
		MessageInfos:      file_trie_proto_msgTypes,
	}.Build()
	File_trie_proto = out.File
	file_trie_proto_rawDesc = nil
	file_trie_proto_goTypes = nil
	file_trie_proto_depIdxs = nil
}
