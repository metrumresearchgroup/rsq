// Code generated by protoc-gen-go. DO NOT EDIT.
// source: internal.proto

package internal

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"
import timestamp "github.com/golang/protobuf/ptypes/timestamp"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type Job_StatusType int32

const (
	Job_QUEUED    Job_StatusType = 0
	Job_RUNNING   Job_StatusType = 1
	Job_COMPLETED Job_StatusType = 2
	Job_ERROR     Job_StatusType = 3
)

var Job_StatusType_name = map[int32]string{
	0: "QUEUED",
	1: "RUNNING",
	2: "COMPLETED",
	3: "ERROR",
}
var Job_StatusType_value = map[string]int32{
	"QUEUED":    0,
	"RUNNING":   1,
	"COMPLETED": 2,
	"ERROR":     3,
}

func (x Job_StatusType) String() string {
	return proto.EnumName(Job_StatusType_name, int32(x))
}
func (Job_StatusType) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_internal_fad1fda03c34fe21, []int{3, 0}
}

type RunDetails struct {
	QueueTime            *timestamp.Timestamp `protobuf:"bytes,1,opt,name=queue_time,json=queueTime,proto3" json:"queue_time,omitempty"`
	StartTime            *timestamp.Timestamp `protobuf:"bytes,2,opt,name=start_time,json=startTime,proto3" json:"start_time,omitempty"`
	EndTime              *timestamp.Timestamp `protobuf:"bytes,3,opt,name=end_time,json=endTime,proto3" json:"end_time,omitempty"`
	Error                string               `protobuf:"bytes,4,opt,name=error,proto3" json:"error,omitempty"`
	XXX_NoUnkeyedLiteral struct{}             `json:"-"`
	XXX_unrecognized     []byte               `json:"-"`
	XXX_sizecache        int32                `json:"-"`
}

func (m *RunDetails) Reset()         { *m = RunDetails{} }
func (m *RunDetails) String() string { return proto.CompactTextString(m) }
func (*RunDetails) ProtoMessage()    {}
func (*RunDetails) Descriptor() ([]byte, []int) {
	return fileDescriptor_internal_fad1fda03c34fe21, []int{0}
}
func (m *RunDetails) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_RunDetails.Unmarshal(m, b)
}
func (m *RunDetails) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_RunDetails.Marshal(b, m, deterministic)
}
func (dst *RunDetails) XXX_Merge(src proto.Message) {
	xxx_messageInfo_RunDetails.Merge(dst, src)
}
func (m *RunDetails) XXX_Size() int {
	return xxx_messageInfo_RunDetails.Size(m)
}
func (m *RunDetails) XXX_DiscardUnknown() {
	xxx_messageInfo_RunDetails.DiscardUnknown(m)
}

var xxx_messageInfo_RunDetails proto.InternalMessageInfo

func (m *RunDetails) GetQueueTime() *timestamp.Timestamp {
	if m != nil {
		return m.QueueTime
	}
	return nil
}

func (m *RunDetails) GetStartTime() *timestamp.Timestamp {
	if m != nil {
		return m.StartTime
	}
	return nil
}

func (m *RunDetails) GetEndTime() *timestamp.Timestamp {
	if m != nil {
		return m.EndTime
	}
	return nil
}

func (m *RunDetails) GetError() string {
	if m != nil {
		return m.Error
	}
	return ""
}

type RScript struct {
	RscriptPath          string   `protobuf:"bytes,1,opt,name=rscript_path,json=rscriptPath,proto3" json:"rscript_path,omitempty"`
	Renv                 []string `protobuf:"bytes,3,rep,name=renv,proto3" json:"renv,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *RScript) Reset()         { *m = RScript{} }
func (m *RScript) String() string { return proto.CompactTextString(m) }
func (*RScript) ProtoMessage()    {}
func (*RScript) Descriptor() ([]byte, []int) {
	return fileDescriptor_internal_fad1fda03c34fe21, []int{1}
}
func (m *RScript) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_RScript.Unmarshal(m, b)
}
func (m *RScript) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_RScript.Marshal(b, m, deterministic)
}
func (dst *RScript) XXX_Merge(src proto.Message) {
	xxx_messageInfo_RScript.Merge(dst, src)
}
func (m *RScript) XXX_Size() int {
	return xxx_messageInfo_RScript.Size(m)
}
func (m *RScript) XXX_DiscardUnknown() {
	xxx_messageInfo_RScript.DiscardUnknown(m)
}

var xxx_messageInfo_RScript proto.InternalMessageInfo

func (m *RScript) GetRscriptPath() string {
	if m != nil {
		return m.RscriptPath
	}
	return ""
}

func (m *RScript) GetRenv() []string {
	if m != nil {
		return m.Renv
	}
	return nil
}

type RScriptResult struct {
	Output               string   `protobuf:"bytes,1,opt,name=output,proto3" json:"output,omitempty"`
	ExitCode             int32    `protobuf:"varint,2,opt,name=exit_code,json=exitCode,proto3" json:"exit_code,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *RScriptResult) Reset()         { *m = RScriptResult{} }
func (m *RScriptResult) String() string { return proto.CompactTextString(m) }
func (*RScriptResult) ProtoMessage()    {}
func (*RScriptResult) Descriptor() ([]byte, []int) {
	return fileDescriptor_internal_fad1fda03c34fe21, []int{2}
}
func (m *RScriptResult) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_RScriptResult.Unmarshal(m, b)
}
func (m *RScriptResult) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_RScriptResult.Marshal(b, m, deterministic)
}
func (dst *RScriptResult) XXX_Merge(src proto.Message) {
	xxx_messageInfo_RScriptResult.Merge(dst, src)
}
func (m *RScriptResult) XXX_Size() int {
	return xxx_messageInfo_RScriptResult.Size(m)
}
func (m *RScriptResult) XXX_DiscardUnknown() {
	xxx_messageInfo_RScriptResult.DiscardUnknown(m)
}

var xxx_messageInfo_RScriptResult proto.InternalMessageInfo

func (m *RScriptResult) GetOutput() string {
	if m != nil {
		return m.Output
	}
	return ""
}

func (m *RScriptResult) GetExitCode() int32 {
	if m != nil {
		return m.ExitCode
	}
	return 0
}

type Job struct {
	Id                   uint64         `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	Status               Job_StatusType `protobuf:"varint,2,opt,name=status,proto3,enum=internal.Job_StatusType" json:"status,omitempty"`
	RunDetails           *RunDetails    `protobuf:"bytes,3,opt,name=run_details,json=runDetails,proto3" json:"run_details,omitempty"`
	Context              string         `protobuf:"bytes,4,opt,name=context,proto3" json:"context,omitempty"`
	Rscript              *RScript       `protobuf:"bytes,5,opt,name=rscript,proto3" json:"rscript,omitempty"`
	Result               *RScriptResult `protobuf:"bytes,6,opt,name=result,proto3" json:"result,omitempty"`
	XXX_NoUnkeyedLiteral struct{}       `json:"-"`
	XXX_unrecognized     []byte         `json:"-"`
	XXX_sizecache        int32          `json:"-"`
}

func (m *Job) Reset()         { *m = Job{} }
func (m *Job) String() string { return proto.CompactTextString(m) }
func (*Job) ProtoMessage()    {}
func (*Job) Descriptor() ([]byte, []int) {
	return fileDescriptor_internal_fad1fda03c34fe21, []int{3}
}
func (m *Job) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Job.Unmarshal(m, b)
}
func (m *Job) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Job.Marshal(b, m, deterministic)
}
func (dst *Job) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Job.Merge(dst, src)
}
func (m *Job) XXX_Size() int {
	return xxx_messageInfo_Job.Size(m)
}
func (m *Job) XXX_DiscardUnknown() {
	xxx_messageInfo_Job.DiscardUnknown(m)
}

var xxx_messageInfo_Job proto.InternalMessageInfo

func (m *Job) GetId() uint64 {
	if m != nil {
		return m.Id
	}
	return 0
}

func (m *Job) GetStatus() Job_StatusType {
	if m != nil {
		return m.Status
	}
	return Job_QUEUED
}

func (m *Job) GetRunDetails() *RunDetails {
	if m != nil {
		return m.RunDetails
	}
	return nil
}

func (m *Job) GetContext() string {
	if m != nil {
		return m.Context
	}
	return ""
}

func (m *Job) GetRscript() *RScript {
	if m != nil {
		return m.Rscript
	}
	return nil
}

func (m *Job) GetResult() *RScriptResult {
	if m != nil {
		return m.Result
	}
	return nil
}

func init() {
	proto.RegisterType((*RunDetails)(nil), "internal.RunDetails")
	proto.RegisterType((*RScript)(nil), "internal.RScript")
	proto.RegisterType((*RScriptResult)(nil), "internal.RScriptResult")
	proto.RegisterType((*Job)(nil), "internal.Job")
	proto.RegisterEnum("internal.Job_StatusType", Job_StatusType_name, Job_StatusType_value)
}

func init() { proto.RegisterFile("internal.proto", fileDescriptor_internal_fad1fda03c34fe21) }

var fileDescriptor_internal_fad1fda03c34fe21 = []byte{
	// 424 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x84, 0x91, 0x41, 0x6f, 0xd3, 0x40,
	0x10, 0x85, 0xb1, 0x9d, 0xd8, 0xf1, 0x84, 0x46, 0x61, 0x55, 0x81, 0x55, 0x0e, 0x04, 0x9f, 0x22,
	0x21, 0x39, 0xa8, 0xa8, 0x07, 0x4e, 0x20, 0x35, 0x16, 0xa2, 0x82, 0xb4, 0x4c, 0x93, 0x73, 0xe4,
	0xc4, 0x4b, 0x6b, 0x29, 0xdd, 0x35, 0xeb, 0x59, 0x54, 0x7e, 0x19, 0x3f, 0x83, 0xbf, 0x84, 0x32,
	0x5e, 0x13, 0xa4, 0x1e, 0x7a, 0xdb, 0x37, 0xfb, 0xbe, 0x27, 0xcd, 0x3c, 0x18, 0x55, 0x8a, 0xa4,
	0x51, 0xc5, 0x2e, 0xab, 0x8d, 0x26, 0x2d, 0x06, 0x9d, 0x3e, 0x79, 0x75, 0xa3, 0xf5, 0xcd, 0x4e,
	0xce, 0x78, 0xbe, 0xb1, 0xdf, 0x67, 0x54, 0xdd, 0xc9, 0x86, 0x8a, 0xbb, 0xba, 0xb5, 0xa6, 0x7f,
	0x3c, 0x00, 0xb4, 0x6a, 0x2e, 0xa9, 0xa8, 0x76, 0x8d, 0x78, 0x0f, 0xf0, 0xc3, 0x4a, 0x2b, 0xd7,
	0x7b, 0x5f, 0xe2, 0x4d, 0xbc, 0xe9, 0xf0, 0xf4, 0x24, 0x6b, 0x43, 0xb2, 0x2e, 0x24, 0x5b, 0x76,
	0x21, 0x18, 0xb3, 0x7b, 0xaf, 0xf7, 0x68, 0x43, 0x85, 0xa1, 0x16, 0xf5, 0x1f, 0x47, 0xd9, 0xcd,
	0xe8, 0x19, 0x0c, 0xa4, 0x2a, 0x5b, 0x30, 0x78, 0x14, 0x8c, 0xa4, 0x2a, 0x19, 0x3b, 0x86, 0xbe,
	0x34, 0x46, 0x9b, 0xa4, 0x37, 0xf1, 0xa6, 0x31, 0xb6, 0x22, 0xfd, 0x08, 0x11, 0x5e, 0x6f, 0x4d,
	0x55, 0x93, 0x78, 0x0d, 0x4f, 0x4d, 0xc3, 0xcf, 0x75, 0x5d, 0xd0, 0x2d, 0xef, 0x13, 0xe3, 0xd0,
	0xcd, 0xae, 0x0a, 0xba, 0x15, 0x02, 0x7a, 0x46, 0xaa, 0x9f, 0x49, 0x30, 0x09, 0xa6, 0x31, 0xf2,
	0x3b, 0x9d, 0xc3, 0x91, 0x4b, 0x40, 0xd9, 0xd8, 0x1d, 0x89, 0xe7, 0x10, 0x6a, 0x4b, 0xb5, 0x25,
	0x97, 0xe0, 0x94, 0x78, 0x09, 0xb1, 0xbc, 0xaf, 0x68, 0xbd, 0xd5, 0x65, 0xbb, 0x71, 0x1f, 0x07,
	0xfb, 0xc1, 0xb9, 0x2e, 0x65, 0xfa, 0xdb, 0x87, 0xe0, 0x42, 0x6f, 0xc4, 0x08, 0xfc, 0xaa, 0x64,
	0xb0, 0x87, 0x7e, 0x55, 0x8a, 0xb7, 0x10, 0x36, 0x54, 0x90, 0x6d, 0x98, 0x18, 0x9d, 0x26, 0xd9,
	0xbf, 0xf6, 0x2e, 0xf4, 0x26, 0xbb, 0xe6, 0xbf, 0xe5, 0xaf, 0x5a, 0xa2, 0xf3, 0x89, 0x33, 0x18,
	0x1a, 0xab, 0xd6, 0x65, 0xdb, 0x91, 0xbb, 0xd0, 0xf1, 0x01, 0x3b, 0xf4, 0x87, 0x60, 0x0e, 0x5d,
	0x26, 0x10, 0x6d, 0xb5, 0x22, 0x79, 0x4f, 0xee, 0x40, 0x9d, 0x14, 0x6f, 0x20, 0x72, 0x37, 0x48,
	0xfa, 0x1c, 0xf6, 0xec, 0xbf, 0x30, 0xb7, 0x79, 0xe7, 0x10, 0x33, 0x08, 0x0d, 0x9f, 0x21, 0x09,
	0xd9, 0xfb, 0xe2, 0xa1, 0x97, 0xbf, 0xd1, 0xd9, 0xd2, 0x0f, 0x00, 0x87, 0x25, 0x04, 0x40, 0xf8,
	0x6d, 0x95, 0xaf, 0xf2, 0xf9, 0xf8, 0x89, 0x18, 0x42, 0x84, 0xab, 0xc5, 0xe2, 0xf3, 0xe2, 0xd3,
	0xd8, 0x13, 0x47, 0x10, 0x9f, 0x5f, 0x7e, 0xbd, 0xfa, 0x92, 0x2f, 0xf3, 0xf9, 0xd8, 0x17, 0x31,
	0xf4, 0x73, 0xc4, 0x4b, 0x1c, 0x07, 0x9b, 0x90, 0x4b, 0x7f, 0xf7, 0x37, 0x00, 0x00, 0xff, 0xff,
	0xfe, 0xe3, 0xe9, 0xc4, 0xd7, 0x02, 0x00, 0x00,
}
