// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.31.0
// 	protoc        v3.15.8
// source: service.proto

package pb

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

var File_service_proto protoreflect.FileDescriptor

var file_service_proto_rawDesc = []byte{
	0x0a, 0x0d, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a,
	0x0e, 0x70, 0x69, 0x6e, 0x67, 0x5f, 0x6d, 0x73, 0x67, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a,
	0x0f, 0x63, 0x68, 0x65, 0x63, 0x6b, 0x5f, 0x6d, 0x73, 0x67, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x1a, 0x0e, 0x62, 0x65, 0x61, 0x74, 0x5f, 0x6d, 0x73, 0x67, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x1a, 0x0d, 0x61, 0x70, 0x70, 0x5f, 0x6d, 0x73, 0x67, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a,
	0x11, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x5f, 0x6d, 0x73, 0x67, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x32, 0xcd, 0x01, 0x0a, 0x07, 0x50, 0x72, 0x6f, 0x6d, 0x6f, 0x74, 0x65, 0x12, 0x1c,
	0x0a, 0x04, 0x50, 0x69, 0x6e, 0x67, 0x12, 0x08, 0x2e, 0x52, 0x65, 0x71, 0x50, 0x69, 0x6e, 0x67,
	0x1a, 0x08, 0x2e, 0x52, 0x65, 0x73, 0x50, 0x69, 0x6e, 0x67, 0x22, 0x00, 0x12, 0x24, 0x0a, 0x0c,
	0x41, 0x72, 0x65, 0x59, 0x6f, 0x75, 0x4c, 0x65, 0x61, 0x64, 0x65, 0x72, 0x12, 0x08, 0x2e, 0x52,
	0x65, 0x71, 0x4c, 0x65, 0x61, 0x64, 0x1a, 0x08, 0x2e, 0x52, 0x65, 0x73, 0x4c, 0x65, 0x61, 0x64,
	0x22, 0x00, 0x12, 0x21, 0x0a, 0x09, 0x48, 0x65, 0x61, 0x72, 0x74, 0x62, 0x65, 0x61, 0x74, 0x12,
	0x08, 0x2e, 0x52, 0x65, 0x71, 0x42, 0x65, 0x61, 0x74, 0x1a, 0x08, 0x2e, 0x52, 0x65, 0x73, 0x42,
	0x65, 0x61, 0x74, 0x22, 0x00, 0x12, 0x2c, 0x0a, 0x0a, 0x41, 0x70, 0x70, 0x73, 0x53, 0x74, 0x61,
	0x74, 0x75, 0x73, 0x12, 0x0d, 0x2e, 0x52, 0x65, 0x71, 0x41, 0x70, 0x70, 0x53, 0x74, 0x61, 0x74,
	0x75, 0x73, 0x1a, 0x0d, 0x2e, 0x52, 0x65, 0x73, 0x41, 0x70, 0x70, 0x53, 0x74, 0x61, 0x74, 0x75,
	0x73, 0x22, 0x00, 0x12, 0x2d, 0x0a, 0x09, 0x41, 0x70, 0x70, 0x41, 0x63, 0x74, 0x69, 0x6f, 0x6e,
	0x12, 0x0e, 0x2e, 0x52, 0x65, 0x71, 0x41, 0x70, 0x70, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65,
	0x1a, 0x0e, 0x2e, 0x52, 0x65, 0x73, 0x41, 0x70, 0x70, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65,
	0x22, 0x00, 0x42, 0x0d, 0x5a, 0x0b, 0x2e, 0x2f, 0x6d, 0x65, 0x64, 0x69, 0x75, 0x6d, 0x2f, 0x70,
	0x62, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var file_service_proto_goTypes = []interface{}{
	(*ReqPing)(nil),       // 0: ReqPing
	(*ReqLead)(nil),       // 1: ReqLead
	(*ReqBeat)(nil),       // 2: ReqBeat
	(*ReqAppStatus)(nil),  // 3: ReqAppStatus
	(*ReqAppService)(nil), // 4: ReqAppService
	(*ResPing)(nil),       // 5: ResPing
	(*ResLead)(nil),       // 6: ResLead
	(*ResBeat)(nil),       // 7: ResBeat
	(*ResAppStatus)(nil),  // 8: ResAppStatus
	(*ResAppService)(nil), // 9: ResAppService
}
var file_service_proto_depIdxs = []int32{
	0, // 0: Promote.Ping:input_type -> ReqPing
	1, // 1: Promote.AreYouLeader:input_type -> ReqLead
	2, // 2: Promote.Heartbeat:input_type -> ReqBeat
	3, // 3: Promote.AppsStatus:input_type -> ReqAppStatus
	4, // 4: Promote.AppAction:input_type -> ReqAppService
	5, // 5: Promote.Ping:output_type -> ResPing
	6, // 6: Promote.AreYouLeader:output_type -> ResLead
	7, // 7: Promote.Heartbeat:output_type -> ResBeat
	8, // 8: Promote.AppsStatus:output_type -> ResAppStatus
	9, // 9: Promote.AppAction:output_type -> ResAppService
	5, // [5:10] is the sub-list for method output_type
	0, // [0:5] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_service_proto_init() }
func file_service_proto_init() {
	if File_service_proto != nil {
		return
	}
	file_ping_msg_proto_init()
	file_check_msg_proto_init()
	file_beat_msg_proto_init()
	file_app_msg_proto_init()
	file_service_msg_proto_init()
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_service_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   0,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_service_proto_goTypes,
		DependencyIndexes: file_service_proto_depIdxs,
	}.Build()
	File_service_proto = out.File
	file_service_proto_rawDesc = nil
	file_service_proto_goTypes = nil
	file_service_proto_depIdxs = nil
}
