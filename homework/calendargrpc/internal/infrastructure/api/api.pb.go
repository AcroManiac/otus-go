// Code generated by protoc-gen-go. DO NOT EDIT.
// source: api.proto

package api

import (
	context "context"
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	duration "github.com/golang/protobuf/ptypes/duration"
	timestamp "github.com/golang/protobuf/ptypes/timestamp"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

type TimePeriod int32

const (
	TimePeriod_TIME_UNKNOWN TimePeriod = 0
	TimePeriod_TIME_DAY     TimePeriod = 1
	TimePeriod_TIME_WEEK    TimePeriod = 2
	TimePeriod_TIME_MONTH   TimePeriod = 3
)

var TimePeriod_name = map[int32]string{
	0: "TIME_UNKNOWN",
	1: "TIME_DAY",
	2: "TIME_WEEK",
	3: "TIME_MONTH",
}

var TimePeriod_value = map[string]int32{
	"TIME_UNKNOWN": 0,
	"TIME_DAY":     1,
	"TIME_WEEK":    2,
	"TIME_MONTH":   3,
}

func (x TimePeriod) String() string {
	return proto.EnumName(TimePeriod_name, int32(x))
}

func (TimePeriod) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_00212fb1f9d3bf1c, []int{0}
}

type CreateEventRequest struct {
	Title                string               `protobuf:"bytes,1,opt,name=title,proto3" json:"title,omitempty"`
	Description          string               `protobuf:"bytes,2,opt,name=description,proto3" json:"description,omitempty"`
	StartTime            *timestamp.Timestamp `protobuf:"bytes,4,opt,name=start_time,json=startTime,proto3" json:"start_time,omitempty"`
	Duration             *duration.Duration   `protobuf:"bytes,5,opt,name=duration,proto3" json:"duration,omitempty"`
	XXX_NoUnkeyedLiteral struct{}             `json:"-"`
	XXX_unrecognized     []byte               `json:"-"`
	XXX_sizecache        int32                `json:"-"`
}

func (m *CreateEventRequest) Reset()         { *m = CreateEventRequest{} }
func (m *CreateEventRequest) String() string { return proto.CompactTextString(m) }
func (*CreateEventRequest) ProtoMessage()    {}
func (*CreateEventRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_00212fb1f9d3bf1c, []int{0}
}

func (m *CreateEventRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_CreateEventRequest.Unmarshal(m, b)
}
func (m *CreateEventRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_CreateEventRequest.Marshal(b, m, deterministic)
}
func (m *CreateEventRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CreateEventRequest.Merge(m, src)
}
func (m *CreateEventRequest) XXX_Size() int {
	return xxx_messageInfo_CreateEventRequest.Size(m)
}
func (m *CreateEventRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_CreateEventRequest.DiscardUnknown(m)
}

var xxx_messageInfo_CreateEventRequest proto.InternalMessageInfo

func (m *CreateEventRequest) GetTitle() string {
	if m != nil {
		return m.Title
	}
	return ""
}

func (m *CreateEventRequest) GetDescription() string {
	if m != nil {
		return m.Description
	}
	return ""
}

func (m *CreateEventRequest) GetStartTime() *timestamp.Timestamp {
	if m != nil {
		return m.StartTime
	}
	return nil
}

func (m *CreateEventRequest) GetDuration() *duration.Duration {
	if m != nil {
		return m.Duration
	}
	return nil
}

type CreateEventResponse struct {
	// Types that are valid to be assigned to Result:
	//	*CreateEventResponse_Event
	//	*CreateEventResponse_Error
	Result               isCreateEventResponse_Result `protobuf_oneof:"result"`
	XXX_NoUnkeyedLiteral struct{}                     `json:"-"`
	XXX_unrecognized     []byte                       `json:"-"`
	XXX_sizecache        int32                        `json:"-"`
}

func (m *CreateEventResponse) Reset()         { *m = CreateEventResponse{} }
func (m *CreateEventResponse) String() string { return proto.CompactTextString(m) }
func (*CreateEventResponse) ProtoMessage()    {}
func (*CreateEventResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_00212fb1f9d3bf1c, []int{1}
}

func (m *CreateEventResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_CreateEventResponse.Unmarshal(m, b)
}
func (m *CreateEventResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_CreateEventResponse.Marshal(b, m, deterministic)
}
func (m *CreateEventResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CreateEventResponse.Merge(m, src)
}
func (m *CreateEventResponse) XXX_Size() int {
	return xxx_messageInfo_CreateEventResponse.Size(m)
}
func (m *CreateEventResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_CreateEventResponse.DiscardUnknown(m)
}

var xxx_messageInfo_CreateEventResponse proto.InternalMessageInfo

type isCreateEventResponse_Result interface {
	isCreateEventResponse_Result()
}

type CreateEventResponse_Event struct {
	Event *Event `protobuf:"bytes,1,opt,name=event,proto3,oneof"`
}

type CreateEventResponse_Error struct {
	Error string `protobuf:"bytes,2,opt,name=error,proto3,oneof"`
}

func (*CreateEventResponse_Event) isCreateEventResponse_Result() {}

func (*CreateEventResponse_Error) isCreateEventResponse_Result() {}

func (m *CreateEventResponse) GetResult() isCreateEventResponse_Result {
	if m != nil {
		return m.Result
	}
	return nil
}

func (m *CreateEventResponse) GetEvent() *Event {
	if x, ok := m.GetResult().(*CreateEventResponse_Event); ok {
		return x.Event
	}
	return nil
}

func (m *CreateEventResponse) GetError() string {
	if x, ok := m.GetResult().(*CreateEventResponse_Error); ok {
		return x.Error
	}
	return ""
}

// XXX_OneofWrappers is for the internal use of the proto package.
func (*CreateEventResponse) XXX_OneofWrappers() []interface{} {
	return []interface{}{
		(*CreateEventResponse_Event)(nil),
		(*CreateEventResponse_Error)(nil),
	}
}

type EditEventRequest struct {
	Id                   string   `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Event                *Event   `protobuf:"bytes,2,opt,name=event,proto3" json:"event,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *EditEventRequest) Reset()         { *m = EditEventRequest{} }
func (m *EditEventRequest) String() string { return proto.CompactTextString(m) }
func (*EditEventRequest) ProtoMessage()    {}
func (*EditEventRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_00212fb1f9d3bf1c, []int{2}
}

func (m *EditEventRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_EditEventRequest.Unmarshal(m, b)
}
func (m *EditEventRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_EditEventRequest.Marshal(b, m, deterministic)
}
func (m *EditEventRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_EditEventRequest.Merge(m, src)
}
func (m *EditEventRequest) XXX_Size() int {
	return xxx_messageInfo_EditEventRequest.Size(m)
}
func (m *EditEventRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_EditEventRequest.DiscardUnknown(m)
}

var xxx_messageInfo_EditEventRequest proto.InternalMessageInfo

func (m *EditEventRequest) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *EditEventRequest) GetEvent() *Event {
	if m != nil {
		return m.Event
	}
	return nil
}

type EditEventResponse struct {
	Error                string   `protobuf:"bytes,1,opt,name=error,proto3" json:"error,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *EditEventResponse) Reset()         { *m = EditEventResponse{} }
func (m *EditEventResponse) String() string { return proto.CompactTextString(m) }
func (*EditEventResponse) ProtoMessage()    {}
func (*EditEventResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_00212fb1f9d3bf1c, []int{3}
}

func (m *EditEventResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_EditEventResponse.Unmarshal(m, b)
}
func (m *EditEventResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_EditEventResponse.Marshal(b, m, deterministic)
}
func (m *EditEventResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_EditEventResponse.Merge(m, src)
}
func (m *EditEventResponse) XXX_Size() int {
	return xxx_messageInfo_EditEventResponse.Size(m)
}
func (m *EditEventResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_EditEventResponse.DiscardUnknown(m)
}

var xxx_messageInfo_EditEventResponse proto.InternalMessageInfo

func (m *EditEventResponse) GetError() string {
	if m != nil {
		return m.Error
	}
	return ""
}

type DeleteEventRequest struct {
	Id                   string   `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *DeleteEventRequest) Reset()         { *m = DeleteEventRequest{} }
func (m *DeleteEventRequest) String() string { return proto.CompactTextString(m) }
func (*DeleteEventRequest) ProtoMessage()    {}
func (*DeleteEventRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_00212fb1f9d3bf1c, []int{4}
}

func (m *DeleteEventRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_DeleteEventRequest.Unmarshal(m, b)
}
func (m *DeleteEventRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_DeleteEventRequest.Marshal(b, m, deterministic)
}
func (m *DeleteEventRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_DeleteEventRequest.Merge(m, src)
}
func (m *DeleteEventRequest) XXX_Size() int {
	return xxx_messageInfo_DeleteEventRequest.Size(m)
}
func (m *DeleteEventRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_DeleteEventRequest.DiscardUnknown(m)
}

var xxx_messageInfo_DeleteEventRequest proto.InternalMessageInfo

func (m *DeleteEventRequest) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

type DeleteEventResponse struct {
	Error                string   `protobuf:"bytes,1,opt,name=error,proto3" json:"error,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *DeleteEventResponse) Reset()         { *m = DeleteEventResponse{} }
func (m *DeleteEventResponse) String() string { return proto.CompactTextString(m) }
func (*DeleteEventResponse) ProtoMessage()    {}
func (*DeleteEventResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_00212fb1f9d3bf1c, []int{5}
}

func (m *DeleteEventResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_DeleteEventResponse.Unmarshal(m, b)
}
func (m *DeleteEventResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_DeleteEventResponse.Marshal(b, m, deterministic)
}
func (m *DeleteEventResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_DeleteEventResponse.Merge(m, src)
}
func (m *DeleteEventResponse) XXX_Size() int {
	return xxx_messageInfo_DeleteEventResponse.Size(m)
}
func (m *DeleteEventResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_DeleteEventResponse.DiscardUnknown(m)
}

var xxx_messageInfo_DeleteEventResponse proto.InternalMessageInfo

func (m *DeleteEventResponse) GetError() string {
	if m != nil {
		return m.Error
	}
	return ""
}

type GetEventsRequest struct {
	Period               TimePeriod           `protobuf:"varint,1,opt,name=period,proto3,enum=TimePeriod" json:"period,omitempty"`
	StartTime            *timestamp.Timestamp `protobuf:"bytes,2,opt,name=start_time,json=startTime,proto3" json:"start_time,omitempty"`
	XXX_NoUnkeyedLiteral struct{}             `json:"-"`
	XXX_unrecognized     []byte               `json:"-"`
	XXX_sizecache        int32                `json:"-"`
}

func (m *GetEventsRequest) Reset()         { *m = GetEventsRequest{} }
func (m *GetEventsRequest) String() string { return proto.CompactTextString(m) }
func (*GetEventsRequest) ProtoMessage()    {}
func (*GetEventsRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_00212fb1f9d3bf1c, []int{6}
}

func (m *GetEventsRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GetEventsRequest.Unmarshal(m, b)
}
func (m *GetEventsRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GetEventsRequest.Marshal(b, m, deterministic)
}
func (m *GetEventsRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GetEventsRequest.Merge(m, src)
}
func (m *GetEventsRequest) XXX_Size() int {
	return xxx_messageInfo_GetEventsRequest.Size(m)
}
func (m *GetEventsRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_GetEventsRequest.DiscardUnknown(m)
}

var xxx_messageInfo_GetEventsRequest proto.InternalMessageInfo

func (m *GetEventsRequest) GetPeriod() TimePeriod {
	if m != nil {
		return m.Period
	}
	return TimePeriod_TIME_UNKNOWN
}

func (m *GetEventsRequest) GetStartTime() *timestamp.Timestamp {
	if m != nil {
		return m.StartTime
	}
	return nil
}

type GetEventsResponse struct {
	Events               []*Event `protobuf:"bytes,1,rep,name=events,proto3" json:"events,omitempty"`
	Error                string   `protobuf:"bytes,2,opt,name=error,proto3" json:"error,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GetEventsResponse) Reset()         { *m = GetEventsResponse{} }
func (m *GetEventsResponse) String() string { return proto.CompactTextString(m) }
func (*GetEventsResponse) ProtoMessage()    {}
func (*GetEventsResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_00212fb1f9d3bf1c, []int{7}
}

func (m *GetEventsResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GetEventsResponse.Unmarshal(m, b)
}
func (m *GetEventsResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GetEventsResponse.Marshal(b, m, deterministic)
}
func (m *GetEventsResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GetEventsResponse.Merge(m, src)
}
func (m *GetEventsResponse) XXX_Size() int {
	return xxx_messageInfo_GetEventsResponse.Size(m)
}
func (m *GetEventsResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_GetEventsResponse.DiscardUnknown(m)
}

var xxx_messageInfo_GetEventsResponse proto.InternalMessageInfo

func (m *GetEventsResponse) GetEvents() []*Event {
	if m != nil {
		return m.Events
	}
	return nil
}

func (m *GetEventsResponse) GetError() string {
	if m != nil {
		return m.Error
	}
	return ""
}

func init() {
	proto.RegisterEnum("TimePeriod", TimePeriod_name, TimePeriod_value)
	proto.RegisterType((*CreateEventRequest)(nil), "CreateEventRequest")
	proto.RegisterType((*CreateEventResponse)(nil), "CreateEventResponse")
	proto.RegisterType((*EditEventRequest)(nil), "EditEventRequest")
	proto.RegisterType((*EditEventResponse)(nil), "EditEventResponse")
	proto.RegisterType((*DeleteEventRequest)(nil), "DeleteEventRequest")
	proto.RegisterType((*DeleteEventResponse)(nil), "DeleteEventResponse")
	proto.RegisterType((*GetEventsRequest)(nil), "GetEventsRequest")
	proto.RegisterType((*GetEventsResponse)(nil), "GetEventsResponse")
}

func init() { proto.RegisterFile("api.proto", fileDescriptor_00212fb1f9d3bf1c) }

var fileDescriptor_00212fb1f9d3bf1c = []byte{
	// 528 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x94, 0x53, 0xc1, 0x6e, 0xd3, 0x40,
	0x10, 0x8d, 0x9d, 0xd8, 0x24, 0xe3, 0x52, 0x39, 0x93, 0x08, 0x05, 0x0b, 0x95, 0x68, 0xe1, 0x50,
	0x40, 0xda, 0x4a, 0x01, 0x0e, 0x70, 0xa2, 0x6d, 0x22, 0xda, 0x46, 0x4d, 0x91, 0x15, 0x14, 0xc1,
	0xa5, 0x72, 0xf1, 0xb6, 0x5a, 0xc9, 0xf5, 0x9a, 0xf5, 0x06, 0xc4, 0xa7, 0xf2, 0x19, 0xfc, 0x01,
	0xf2, 0xda, 0x4e, 0x1d, 0xbb, 0x42, 0xe2, 0xf8, 0x66, 0xc6, 0xf3, 0xde, 0xbe, 0x37, 0x86, 0x5e,
	0x90, 0x70, 0x9a, 0x48, 0xa1, 0x84, 0xe7, 0xb0, 0x1f, 0x2c, 0x56, 0x05, 0x78, 0x7a, 0x23, 0xc4,
	0x4d, 0xc4, 0x0e, 0x34, 0xba, 0x5a, 0x5f, 0x1f, 0x28, 0x7e, 0xcb, 0x52, 0x15, 0xdc, 0x26, 0xc5,
	0xc0, 0x5e, 0x7d, 0x20, 0x5c, 0xcb, 0x40, 0x71, 0x11, 0xe7, 0x7d, 0xf2, 0xdb, 0x00, 0x3c, 0x96,
	0x2c, 0x50, 0x6c, 0x96, 0xad, 0xf5, 0xd9, 0xf7, 0x35, 0x4b, 0x15, 0x0e, 0xc1, 0x52, 0x5c, 0x45,
	0x6c, 0x64, 0x8c, 0x8d, 0xfd, 0x9e, 0x9f, 0x03, 0x1c, 0x83, 0x13, 0xb2, 0xf4, 0x9b, 0xe4, 0x49,
	0xb6, 0x61, 0x64, 0xea, 0x5e, 0xb5, 0x84, 0xef, 0x00, 0x52, 0x15, 0x48, 0x75, 0x99, 0xe9, 0x18,
	0x75, 0xc6, 0xc6, 0xbe, 0x33, 0xf1, 0x68, 0xae, 0x81, 0x96, 0x1a, 0xe8, 0xb2, 0x14, 0xe9, 0xf7,
	0xf4, 0x74, 0x86, 0xf1, 0x2d, 0x74, 0x4b, 0x6d, 0x23, 0x4b, 0x7f, 0xf8, 0xb8, 0xf1, 0xe1, 0xb4,
	0x18, 0xf0, 0x37, 0xa3, 0x67, 0x9d, 0x6e, 0xdb, 0xed, 0x9c, 0x75, 0xba, 0xb6, 0xfb, 0xc0, 0xb7,
	0xc4, 0xcf, 0x98, 0x49, 0xdf, 0x8e, 0x85, 0xe2, 0xd7, 0xbf, 0xc8, 0x0a, 0x06, 0x5b, 0x4f, 0x4b,
	0x13, 0x11, 0xa7, 0x0c, 0xf7, 0xc0, 0xd2, 0x16, 0xea, 0xb7, 0x39, 0x13, 0x9b, 0xea, 0xf6, 0x49,
	0xcb, 0xcf, 0xcb, 0xf8, 0x08, 0x2c, 0x26, 0xa5, 0x90, 0xf9, 0xfb, 0x74, 0x3d, 0x83, 0x47, 0x5d,
	0xb0, 0x25, 0x4b, 0xd7, 0x91, 0x22, 0x1f, 0xc0, 0x9d, 0x85, 0x5c, 0x6d, 0x39, 0xb6, 0x0b, 0x26,
	0x0f, 0x0b, 0xbb, 0x4c, 0x1e, 0xe2, 0x93, 0x92, 0xc5, 0xac, 0xb2, 0x14, 0x1c, 0xe4, 0x05, 0xf4,
	0x2b, 0x1b, 0x0a, 0x61, 0xc3, 0x92, 0xb8, 0x30, 0x5d, 0x03, 0xf2, 0x1c, 0x70, 0xca, 0x22, 0x56,
	0x0b, 0xa8, 0x46, 0x47, 0x5e, 0xc1, 0x60, 0x6b, 0xea, 0x9f, 0x2b, 0x25, 0xb8, 0x1f, 0x59, 0x4e,
	0x9e, 0x96, 0x0b, 0x9f, 0x81, 0x9d, 0x30, 0xc9, 0x45, 0xbe, 0x74, 0x77, 0xe2, 0xe8, 0x94, 0x3e,
	0xe9, 0x92, 0x5f, 0xb4, 0x6a, 0xf1, 0x9a, 0xff, 0x11, 0x2f, 0x39, 0x85, 0x7e, 0x85, 0x73, 0x13,
	0x85, 0xad, 0xfd, 0x48, 0x47, 0xc6, 0xb8, 0x5d, 0x71, 0xa9, 0xa8, 0xde, 0xc9, 0x37, 0x2b, 0xf2,
	0x5f, 0xce, 0x01, 0xee, 0xb4, 0xa1, 0x0b, 0x3b, 0xcb, 0xd3, 0xf3, 0xd9, 0xe5, 0xe7, 0xc5, 0x7c,
	0x71, 0xb1, 0x5a, 0xb8, 0x2d, 0xdc, 0x81, 0xae, 0xae, 0x4c, 0x0f, 0xbf, 0xb8, 0x06, 0x3e, 0x84,
	0x9e, 0x46, 0xab, 0xd9, 0x6c, 0xee, 0x9a, 0xb8, 0x0b, 0xa0, 0xe1, 0xf9, 0xc5, 0x62, 0x79, 0xe2,
	0xb6, 0x27, 0x7f, 0x0c, 0x70, 0x8e, 0x83, 0x88, 0xc5, 0x61, 0x20, 0x0f, 0x13, 0x8e, 0xef, 0xc1,
	0xa9, 0x1c, 0x0d, 0x0e, 0x68, 0xf3, 0xef, 0xf0, 0x86, 0xf4, 0x9e, 0xbb, 0x22, 0x2d, 0x7c, 0x03,
	0xbd, 0x4d, 0xaa, 0xd8, 0xa7, 0xf5, 0x1b, 0xf1, 0x90, 0x36, 0x42, 0x27, 0xad, 0x8c, 0xb1, 0x12,
	0x1d, 0x0e, 0x68, 0x33, 0x6e, 0x6f, 0x48, 0xef, 0x49, 0x37, 0x67, 0xdc, 0xb8, 0x8a, 0x7d, 0x5a,
	0x4f, 0xd5, 0x43, 0xda, 0x30, 0x9d, 0xb4, 0x8e, 0xac, 0xaf, 0xed, 0x20, 0xe1, 0x57, 0xb6, 0x4e,
	0xec, 0xf5, 0xdf, 0x00, 0x00, 0x00, 0xff, 0xff, 0x1a, 0xc1, 0xcb, 0x5c, 0x5d, 0x04, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConnInterface

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion6

// CalendarApiClient is the client API for CalendarApi service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type CalendarApiClient interface {
	CreateEvent(ctx context.Context, in *CreateEventRequest, opts ...grpc.CallOption) (*CreateEventResponse, error)
	EditEvent(ctx context.Context, in *EditEventRequest, opts ...grpc.CallOption) (*EditEventResponse, error)
	DeleteEvent(ctx context.Context, in *DeleteEventRequest, opts ...grpc.CallOption) (*DeleteEventResponse, error)
	GetEvents(ctx context.Context, in *GetEventsRequest, opts ...grpc.CallOption) (*GetEventsResponse, error)
}

type calendarApiClient struct {
	cc grpc.ClientConnInterface
}

func NewCalendarApiClient(cc grpc.ClientConnInterface) CalendarApiClient {
	return &calendarApiClient{cc}
}

func (c *calendarApiClient) CreateEvent(ctx context.Context, in *CreateEventRequest, opts ...grpc.CallOption) (*CreateEventResponse, error) {
	out := new(CreateEventResponse)
	err := c.cc.Invoke(ctx, "/CalendarApi/CreateEvent", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *calendarApiClient) EditEvent(ctx context.Context, in *EditEventRequest, opts ...grpc.CallOption) (*EditEventResponse, error) {
	out := new(EditEventResponse)
	err := c.cc.Invoke(ctx, "/CalendarApi/EditEvent", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *calendarApiClient) DeleteEvent(ctx context.Context, in *DeleteEventRequest, opts ...grpc.CallOption) (*DeleteEventResponse, error) {
	out := new(DeleteEventResponse)
	err := c.cc.Invoke(ctx, "/CalendarApi/DeleteEvent", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *calendarApiClient) GetEvents(ctx context.Context, in *GetEventsRequest, opts ...grpc.CallOption) (*GetEventsResponse, error) {
	out := new(GetEventsResponse)
	err := c.cc.Invoke(ctx, "/CalendarApi/GetEvents", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// CalendarApiServer is the server API for CalendarApi service.
type CalendarApiServer interface {
	CreateEvent(context.Context, *CreateEventRequest) (*CreateEventResponse, error)
	EditEvent(context.Context, *EditEventRequest) (*EditEventResponse, error)
	DeleteEvent(context.Context, *DeleteEventRequest) (*DeleteEventResponse, error)
	GetEvents(context.Context, *GetEventsRequest) (*GetEventsResponse, error)
}

// UnimplementedCalendarApiServer can be embedded to have forward compatible implementations.
type UnimplementedCalendarApiServer struct {
}

func (*UnimplementedCalendarApiServer) CreateEvent(ctx context.Context, req *CreateEventRequest) (*CreateEventResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateEvent not implemented")
}
func (*UnimplementedCalendarApiServer) EditEvent(ctx context.Context, req *EditEventRequest) (*EditEventResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method EditEvent not implemented")
}
func (*UnimplementedCalendarApiServer) DeleteEvent(ctx context.Context, req *DeleteEventRequest) (*DeleteEventResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteEvent not implemented")
}
func (*UnimplementedCalendarApiServer) GetEvents(ctx context.Context, req *GetEventsRequest) (*GetEventsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetEvents not implemented")
}

func RegisterCalendarApiServer(s *grpc.Server, srv CalendarApiServer) {
	s.RegisterService(&_CalendarApi_serviceDesc, srv)
}

func _CalendarApi_CreateEvent_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateEventRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CalendarApiServer).CreateEvent(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/CalendarApi/CreateEvent",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CalendarApiServer).CreateEvent(ctx, req.(*CreateEventRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _CalendarApi_EditEvent_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(EditEventRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CalendarApiServer).EditEvent(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/CalendarApi/EditEvent",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CalendarApiServer).EditEvent(ctx, req.(*EditEventRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _CalendarApi_DeleteEvent_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteEventRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CalendarApiServer).DeleteEvent(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/CalendarApi/DeleteEvent",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CalendarApiServer).DeleteEvent(ctx, req.(*DeleteEventRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _CalendarApi_GetEvents_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetEventsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CalendarApiServer).GetEvents(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/CalendarApi/GetEvents",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CalendarApiServer).GetEvents(ctx, req.(*GetEventsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _CalendarApi_serviceDesc = grpc.ServiceDesc{
	ServiceName: "CalendarApi",
	HandlerType: (*CalendarApiServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreateEvent",
			Handler:    _CalendarApi_CreateEvent_Handler,
		},
		{
			MethodName: "EditEvent",
			Handler:    _CalendarApi_EditEvent_Handler,
		},
		{
			MethodName: "DeleteEvent",
			Handler:    _CalendarApi_DeleteEvent_Handler,
		},
		{
			MethodName: "GetEvents",
			Handler:    _CalendarApi_GetEvents_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "api.proto",
}