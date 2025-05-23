// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.6
// 	protoc        v6.30.2
// source: rpc_refresh_token.proto

package pb

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
	reflect "reflect"
	sync "sync"
	unsafe "unsafe"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type RefreshTokenRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	RefreshToken  string                 `protobuf:"bytes,1,opt,name=refresh_token,json=refreshToken,proto3" json:"refresh_token,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *RefreshTokenRequest) Reset() {
	*x = RefreshTokenRequest{}
	mi := &file_rpc_refresh_token_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *RefreshTokenRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RefreshTokenRequest) ProtoMessage() {}

func (x *RefreshTokenRequest) ProtoReflect() protoreflect.Message {
	mi := &file_rpc_refresh_token_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RefreshTokenRequest.ProtoReflect.Descriptor instead.
func (*RefreshTokenRequest) Descriptor() ([]byte, []int) {
	return file_rpc_refresh_token_proto_rawDescGZIP(), []int{0}
}

func (x *RefreshTokenRequest) GetRefreshToken() string {
	if x != nil {
		return x.RefreshToken
	}
	return ""
}

type RefreshTokenResponse struct {
	state                protoimpl.MessageState `protogen:"open.v1"`
	AccessToken          string                 `protobuf:"bytes,1,opt,name=access_token,json=accessToken,proto3" json:"access_token,omitempty"`
	AccessTokenExpiresAt *timestamppb.Timestamp `protobuf:"bytes,2,opt,name=access_token_expires_at,json=accessTokenExpiresAt,proto3" json:"access_token_expires_at,omitempty"`
	User                 *User                  `protobuf:"bytes,3,opt,name=user,proto3" json:"user,omitempty"`
	Session              *Session               `protobuf:"bytes,4,opt,name=session,proto3" json:"session,omitempty"`
	unknownFields        protoimpl.UnknownFields
	sizeCache            protoimpl.SizeCache
}

func (x *RefreshTokenResponse) Reset() {
	*x = RefreshTokenResponse{}
	mi := &file_rpc_refresh_token_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *RefreshTokenResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RefreshTokenResponse) ProtoMessage() {}

func (x *RefreshTokenResponse) ProtoReflect() protoreflect.Message {
	mi := &file_rpc_refresh_token_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RefreshTokenResponse.ProtoReflect.Descriptor instead.
func (*RefreshTokenResponse) Descriptor() ([]byte, []int) {
	return file_rpc_refresh_token_proto_rawDescGZIP(), []int{1}
}

func (x *RefreshTokenResponse) GetAccessToken() string {
	if x != nil {
		return x.AccessToken
	}
	return ""
}

func (x *RefreshTokenResponse) GetAccessTokenExpiresAt() *timestamppb.Timestamp {
	if x != nil {
		return x.AccessTokenExpiresAt
	}
	return nil
}

func (x *RefreshTokenResponse) GetUser() *User {
	if x != nil {
		return x.User
	}
	return nil
}

func (x *RefreshTokenResponse) GetSession() *Session {
	if x != nil {
		return x.Session
	}
	return nil
}

var File_rpc_refresh_token_proto protoreflect.FileDescriptor

const file_rpc_refresh_token_proto_rawDesc = "" +
	"\n" +
	"\x17rpc_refresh_token.proto\x12\x02pb\x1a\n" +
	"user.proto\x1a\rsession.proto\x1a\x1fgoogle/protobuf/timestamp.proto\":\n" +
	"\x13RefreshTokenRequest\x12#\n" +
	"\rrefresh_token\x18\x01 \x01(\tR\frefreshToken\"\xd1\x01\n" +
	"\x14RefreshTokenResponse\x12!\n" +
	"\faccess_token\x18\x01 \x01(\tR\vaccessToken\x12Q\n" +
	"\x17access_token_expires_at\x18\x02 \x01(\v2\x1a.google.protobuf.TimestampR\x14accessTokenExpiresAt\x12\x1c\n" +
	"\x04user\x18\x03 \x01(\v2\b.pb.UserR\x04user\x12%\n" +
	"\asession\x18\x04 \x01(\v2\v.pb.SessionR\asessionB.Z,github.com/noueii/gonuxt-starter/internal/pbb\x06proto3"

var (
	file_rpc_refresh_token_proto_rawDescOnce sync.Once
	file_rpc_refresh_token_proto_rawDescData []byte
)

func file_rpc_refresh_token_proto_rawDescGZIP() []byte {
	file_rpc_refresh_token_proto_rawDescOnce.Do(func() {
		file_rpc_refresh_token_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_rpc_refresh_token_proto_rawDesc), len(file_rpc_refresh_token_proto_rawDesc)))
	})
	return file_rpc_refresh_token_proto_rawDescData
}

var file_rpc_refresh_token_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_rpc_refresh_token_proto_goTypes = []any{
	(*RefreshTokenRequest)(nil),   // 0: pb.RefreshTokenRequest
	(*RefreshTokenResponse)(nil),  // 1: pb.RefreshTokenResponse
	(*timestamppb.Timestamp)(nil), // 2: google.protobuf.Timestamp
	(*User)(nil),                  // 3: pb.User
	(*Session)(nil),               // 4: pb.Session
}
var file_rpc_refresh_token_proto_depIdxs = []int32{
	2, // 0: pb.RefreshTokenResponse.access_token_expires_at:type_name -> google.protobuf.Timestamp
	3, // 1: pb.RefreshTokenResponse.user:type_name -> pb.User
	4, // 2: pb.RefreshTokenResponse.session:type_name -> pb.Session
	3, // [3:3] is the sub-list for method output_type
	3, // [3:3] is the sub-list for method input_type
	3, // [3:3] is the sub-list for extension type_name
	3, // [3:3] is the sub-list for extension extendee
	0, // [0:3] is the sub-list for field type_name
}

func init() { file_rpc_refresh_token_proto_init() }
func file_rpc_refresh_token_proto_init() {
	if File_rpc_refresh_token_proto != nil {
		return
	}
	file_user_proto_init()
	file_session_proto_init()
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_rpc_refresh_token_proto_rawDesc), len(file_rpc_refresh_token_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_rpc_refresh_token_proto_goTypes,
		DependencyIndexes: file_rpc_refresh_token_proto_depIdxs,
		MessageInfos:      file_rpc_refresh_token_proto_msgTypes,
	}.Build()
	File_rpc_refresh_token_proto = out.File
	file_rpc_refresh_token_proto_goTypes = nil
	file_rpc_refresh_token_proto_depIdxs = nil
}
