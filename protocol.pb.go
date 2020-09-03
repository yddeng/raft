// Code generated by protoc-gen-go. DO NOT EDIT.
// source: protocol.proto

package raft

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
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

type LogEntry struct {
	Index                uint64   `protobuf:"varint,1,opt,name=Index,proto3" json:"Index,omitempty"`
	Term                 uint64   `protobuf:"varint,2,opt,name=Term,proto3" json:"Term,omitempty"`
	CommandName          string   `protobuf:"bytes,3,opt,name=CommandName,proto3" json:"CommandName,omitempty"`
	Command              []byte   `protobuf:"bytes,4,opt,name=Command,proto3" json:"Command,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *LogEntry) Reset()         { *m = LogEntry{} }
func (m *LogEntry) String() string { return proto.CompactTextString(m) }
func (*LogEntry) ProtoMessage()    {}
func (*LogEntry) Descriptor() ([]byte, []int) {
	return fileDescriptor_2bc2336598a3f7e0, []int{0}
}

func (m *LogEntry) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_LogEntry.Unmarshal(m, b)
}
func (m *LogEntry) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_LogEntry.Marshal(b, m, deterministic)
}
func (m *LogEntry) XXX_Merge(src proto.Message) {
	xxx_messageInfo_LogEntry.Merge(m, src)
}
func (m *LogEntry) XXX_Size() int {
	return xxx_messageInfo_LogEntry.Size(m)
}
func (m *LogEntry) XXX_DiscardUnknown() {
	xxx_messageInfo_LogEntry.DiscardUnknown(m)
}

var xxx_messageInfo_LogEntry proto.InternalMessageInfo

func (m *LogEntry) GetIndex() uint64 {
	if m != nil {
		return m.Index
	}
	return 0
}

func (m *LogEntry) GetTerm() uint64 {
	if m != nil {
		return m.Term
	}
	return 0
}

func (m *LogEntry) GetCommandName() string {
	if m != nil {
		return m.CommandName
	}
	return ""
}

func (m *LogEntry) GetCommand() []byte {
	if m != nil {
		return m.Command
	}
	return nil
}

type AppendEntriesReq struct {
	Term                 uint64      `protobuf:"varint,1,opt,name=Term,proto3" json:"Term,omitempty"`
	LeaderName           string      `protobuf:"bytes,2,opt,name=LeaderName,proto3" json:"LeaderName,omitempty"`
	PrevLogIndex         uint64      `protobuf:"varint,3,opt,name=PrevLogIndex,proto3" json:"PrevLogIndex,omitempty"`
	PrevLogTerm          uint64      `protobuf:"varint,4,opt,name=PrevLogTerm,proto3" json:"PrevLogTerm,omitempty"`
	LeaderCommit         uint64      `protobuf:"varint,5,opt,name=LeaderCommit,proto3" json:"LeaderCommit,omitempty"`
	Entries              []*LogEntry `protobuf:"bytes,6,rep,name=Entries,proto3" json:"Entries,omitempty"`
	XXX_NoUnkeyedLiteral struct{}    `json:"-"`
	XXX_unrecognized     []byte      `json:"-"`
	XXX_sizecache        int32       `json:"-"`
}

func (m *AppendEntriesReq) Reset()         { *m = AppendEntriesReq{} }
func (m *AppendEntriesReq) String() string { return proto.CompactTextString(m) }
func (*AppendEntriesReq) ProtoMessage()    {}
func (*AppendEntriesReq) Descriptor() ([]byte, []int) {
	return fileDescriptor_2bc2336598a3f7e0, []int{1}
}

func (m *AppendEntriesReq) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_AppendEntriesReq.Unmarshal(m, b)
}
func (m *AppendEntriesReq) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_AppendEntriesReq.Marshal(b, m, deterministic)
}
func (m *AppendEntriesReq) XXX_Merge(src proto.Message) {
	xxx_messageInfo_AppendEntriesReq.Merge(m, src)
}
func (m *AppendEntriesReq) XXX_Size() int {
	return xxx_messageInfo_AppendEntriesReq.Size(m)
}
func (m *AppendEntriesReq) XXX_DiscardUnknown() {
	xxx_messageInfo_AppendEntriesReq.DiscardUnknown(m)
}

var xxx_messageInfo_AppendEntriesReq proto.InternalMessageInfo

func (m *AppendEntriesReq) GetTerm() uint64 {
	if m != nil {
		return m.Term
	}
	return 0
}

func (m *AppendEntriesReq) GetLeaderName() string {
	if m != nil {
		return m.LeaderName
	}
	return ""
}

func (m *AppendEntriesReq) GetPrevLogIndex() uint64 {
	if m != nil {
		return m.PrevLogIndex
	}
	return 0
}

func (m *AppendEntriesReq) GetPrevLogTerm() uint64 {
	if m != nil {
		return m.PrevLogTerm
	}
	return 0
}

func (m *AppendEntriesReq) GetLeaderCommit() uint64 {
	if m != nil {
		return m.LeaderCommit
	}
	return 0
}

func (m *AppendEntriesReq) GetEntries() []*LogEntry {
	if m != nil {
		return m.Entries
	}
	return nil
}

type AppendEntriesResp struct {
	Term                 uint64   `protobuf:"varint,1,opt,name=Term,proto3" json:"Term,omitempty"`
	Index                uint64   `protobuf:"varint,2,opt,name=Index,proto3" json:"Index,omitempty"`
	CommitIndex          uint64   `protobuf:"varint,3,opt,name=CommitIndex,proto3" json:"CommitIndex,omitempty"`
	Success              bool     `protobuf:"varint,4,opt,name=Success,proto3" json:"Success,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *AppendEntriesResp) Reset()         { *m = AppendEntriesResp{} }
func (m *AppendEntriesResp) String() string { return proto.CompactTextString(m) }
func (*AppendEntriesResp) ProtoMessage()    {}
func (*AppendEntriesResp) Descriptor() ([]byte, []int) {
	return fileDescriptor_2bc2336598a3f7e0, []int{2}
}

func (m *AppendEntriesResp) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_AppendEntriesResp.Unmarshal(m, b)
}
func (m *AppendEntriesResp) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_AppendEntriesResp.Marshal(b, m, deterministic)
}
func (m *AppendEntriesResp) XXX_Merge(src proto.Message) {
	xxx_messageInfo_AppendEntriesResp.Merge(m, src)
}
func (m *AppendEntriesResp) XXX_Size() int {
	return xxx_messageInfo_AppendEntriesResp.Size(m)
}
func (m *AppendEntriesResp) XXX_DiscardUnknown() {
	xxx_messageInfo_AppendEntriesResp.DiscardUnknown(m)
}

var xxx_messageInfo_AppendEntriesResp proto.InternalMessageInfo

func (m *AppendEntriesResp) GetTerm() uint64 {
	if m != nil {
		return m.Term
	}
	return 0
}

func (m *AppendEntriesResp) GetIndex() uint64 {
	if m != nil {
		return m.Index
	}
	return 0
}

func (m *AppendEntriesResp) GetCommitIndex() uint64 {
	if m != nil {
		return m.CommitIndex
	}
	return 0
}

func (m *AppendEntriesResp) GetSuccess() bool {
	if m != nil {
		return m.Success
	}
	return false
}

type RequestVoteReq struct {
	Term                 uint64   `protobuf:"varint,1,opt,name=Term,proto3" json:"Term,omitempty"`
	LastLogIndex         uint64   `protobuf:"varint,2,opt,name=LastLogIndex,proto3" json:"LastLogIndex,omitempty"`
	LastLogTerm          uint64   `protobuf:"varint,3,opt,name=LastLogTerm,proto3" json:"LastLogTerm,omitempty"`
	CandidateName        string   `protobuf:"bytes,4,opt,name=CandidateName,proto3" json:"CandidateName,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *RequestVoteReq) Reset()         { *m = RequestVoteReq{} }
func (m *RequestVoteReq) String() string { return proto.CompactTextString(m) }
func (*RequestVoteReq) ProtoMessage()    {}
func (*RequestVoteReq) Descriptor() ([]byte, []int) {
	return fileDescriptor_2bc2336598a3f7e0, []int{3}
}

func (m *RequestVoteReq) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_RequestVoteReq.Unmarshal(m, b)
}
func (m *RequestVoteReq) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_RequestVoteReq.Marshal(b, m, deterministic)
}
func (m *RequestVoteReq) XXX_Merge(src proto.Message) {
	xxx_messageInfo_RequestVoteReq.Merge(m, src)
}
func (m *RequestVoteReq) XXX_Size() int {
	return xxx_messageInfo_RequestVoteReq.Size(m)
}
func (m *RequestVoteReq) XXX_DiscardUnknown() {
	xxx_messageInfo_RequestVoteReq.DiscardUnknown(m)
}

var xxx_messageInfo_RequestVoteReq proto.InternalMessageInfo

func (m *RequestVoteReq) GetTerm() uint64 {
	if m != nil {
		return m.Term
	}
	return 0
}

func (m *RequestVoteReq) GetLastLogIndex() uint64 {
	if m != nil {
		return m.LastLogIndex
	}
	return 0
}

func (m *RequestVoteReq) GetLastLogTerm() uint64 {
	if m != nil {
		return m.LastLogTerm
	}
	return 0
}

func (m *RequestVoteReq) GetCandidateName() string {
	if m != nil {
		return m.CandidateName
	}
	return ""
}

type RequestVoteResp struct {
	Term                 uint64   `protobuf:"varint,1,opt,name=Term,proto3" json:"Term,omitempty"`
	VoteGranted          bool     `protobuf:"varint,2,opt,name=VoteGranted,proto3" json:"VoteGranted,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *RequestVoteResp) Reset()         { *m = RequestVoteResp{} }
func (m *RequestVoteResp) String() string { return proto.CompactTextString(m) }
func (*RequestVoteResp) ProtoMessage()    {}
func (*RequestVoteResp) Descriptor() ([]byte, []int) {
	return fileDescriptor_2bc2336598a3f7e0, []int{4}
}

func (m *RequestVoteResp) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_RequestVoteResp.Unmarshal(m, b)
}
func (m *RequestVoteResp) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_RequestVoteResp.Marshal(b, m, deterministic)
}
func (m *RequestVoteResp) XXX_Merge(src proto.Message) {
	xxx_messageInfo_RequestVoteResp.Merge(m, src)
}
func (m *RequestVoteResp) XXX_Size() int {
	return xxx_messageInfo_RequestVoteResp.Size(m)
}
func (m *RequestVoteResp) XXX_DiscardUnknown() {
	xxx_messageInfo_RequestVoteResp.DiscardUnknown(m)
}

var xxx_messageInfo_RequestVoteResp proto.InternalMessageInfo

func (m *RequestVoteResp) GetTerm() uint64 {
	if m != nil {
		return m.Term
	}
	return 0
}

func (m *RequestVoteResp) GetVoteGranted() bool {
	if m != nil {
		return m.VoteGranted
	}
	return false
}

type SnapshotReq struct {
	LeaderName           string   `protobuf:"bytes,1,opt,name=LeaderName,proto3" json:"LeaderName,omitempty"`
	LastIndex            uint64   `protobuf:"varint,2,opt,name=LastIndex,proto3" json:"LastIndex,omitempty"`
	LastTerm             uint64   `protobuf:"varint,3,opt,name=LastTerm,proto3" json:"LastTerm,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *SnapshotReq) Reset()         { *m = SnapshotReq{} }
func (m *SnapshotReq) String() string { return proto.CompactTextString(m) }
func (*SnapshotReq) ProtoMessage()    {}
func (*SnapshotReq) Descriptor() ([]byte, []int) {
	return fileDescriptor_2bc2336598a3f7e0, []int{5}
}

func (m *SnapshotReq) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SnapshotReq.Unmarshal(m, b)
}
func (m *SnapshotReq) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SnapshotReq.Marshal(b, m, deterministic)
}
func (m *SnapshotReq) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SnapshotReq.Merge(m, src)
}
func (m *SnapshotReq) XXX_Size() int {
	return xxx_messageInfo_SnapshotReq.Size(m)
}
func (m *SnapshotReq) XXX_DiscardUnknown() {
	xxx_messageInfo_SnapshotReq.DiscardUnknown(m)
}

var xxx_messageInfo_SnapshotReq proto.InternalMessageInfo

func (m *SnapshotReq) GetLeaderName() string {
	if m != nil {
		return m.LeaderName
	}
	return ""
}

func (m *SnapshotReq) GetLastIndex() uint64 {
	if m != nil {
		return m.LastIndex
	}
	return 0
}

func (m *SnapshotReq) GetLastTerm() uint64 {
	if m != nil {
		return m.LastTerm
	}
	return 0
}

type SnapshotResp struct {
	Success              bool     `protobuf:"varint,1,opt,name=Success,proto3" json:"Success,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *SnapshotResp) Reset()         { *m = SnapshotResp{} }
func (m *SnapshotResp) String() string { return proto.CompactTextString(m) }
func (*SnapshotResp) ProtoMessage()    {}
func (*SnapshotResp) Descriptor() ([]byte, []int) {
	return fileDescriptor_2bc2336598a3f7e0, []int{6}
}

func (m *SnapshotResp) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SnapshotResp.Unmarshal(m, b)
}
func (m *SnapshotResp) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SnapshotResp.Marshal(b, m, deterministic)
}
func (m *SnapshotResp) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SnapshotResp.Merge(m, src)
}
func (m *SnapshotResp) XXX_Size() int {
	return xxx_messageInfo_SnapshotResp.Size(m)
}
func (m *SnapshotResp) XXX_DiscardUnknown() {
	xxx_messageInfo_SnapshotResp.DiscardUnknown(m)
}

var xxx_messageInfo_SnapshotResp proto.InternalMessageInfo

func (m *SnapshotResp) GetSuccess() bool {
	if m != nil {
		return m.Success
	}
	return false
}

func init() {
	proto.RegisterType((*LogEntry)(nil), "LogEntry")
	proto.RegisterType((*AppendEntriesReq)(nil), "AppendEntriesReq")
	proto.RegisterType((*AppendEntriesResp)(nil), "AppendEntriesResp")
	proto.RegisterType((*RequestVoteReq)(nil), "RequestVoteReq")
	proto.RegisterType((*RequestVoteResp)(nil), "RequestVoteResp")
	proto.RegisterType((*SnapshotReq)(nil), "SnapshotReq")
	proto.RegisterType((*SnapshotResp)(nil), "SnapshotResp")
}

func init() { proto.RegisterFile("protocol.proto", fileDescriptor_2bc2336598a3f7e0) }

var fileDescriptor_2bc2336598a3f7e0 = []byte{
	// 394 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x74, 0x53, 0x4d, 0x6b, 0xdb, 0x30,
	0x18, 0x46, 0x8e, 0x93, 0x38, 0xaf, 0xb3, 0x6c, 0x13, 0x3b, 0x98, 0x31, 0x86, 0xd1, 0x76, 0xf0,
	0x29, 0x87, 0xed, 0x17, 0x6c, 0xa1, 0x84, 0x42, 0x28, 0x45, 0x29, 0x3d, 0xf4, 0xa6, 0xc6, 0x6f,
	0x53, 0x43, 0x6d, 0x39, 0x96, 0x52, 0x9a, 0x7f, 0xd1, 0x3f, 0xd7, 0xff, 0x53, 0x24, 0x7f, 0x44,
	0x6e, 0x9b, 0x9b, 0x9e, 0x47, 0xd6, 0xfb, 0x7c, 0x48, 0x86, 0x59, 0x59, 0x49, 0x2d, 0x37, 0xf2,
	0x61, 0x6e, 0x17, 0xac, 0x84, 0x60, 0x25, 0xb7, 0x67, 0x85, 0xae, 0x0e, 0xf4, 0x1b, 0x0c, 0xcf,
	0x8b, 0x14, 0x9f, 0x22, 0x12, 0x93, 0xc4, 0xe7, 0x35, 0xa0, 0x14, 0xfc, 0x2b, 0xac, 0xf2, 0xc8,
	0xb3, 0xa4, 0x5d, 0xd3, 0x18, 0xc2, 0x85, 0xcc, 0x73, 0x51, 0xa4, 0x17, 0x22, 0xc7, 0x68, 0x10,
	0x93, 0x64, 0xc2, 0x5d, 0x8a, 0x46, 0x30, 0x6e, 0x60, 0xe4, 0xc7, 0x24, 0x99, 0xf2, 0x16, 0xb2,
	0x17, 0x02, 0x5f, 0xfe, 0x95, 0x25, 0x16, 0xa9, 0x51, 0xcd, 0x50, 0x71, 0xdc, 0x75, 0x22, 0xc4,
	0x11, 0xf9, 0x09, 0xb0, 0x42, 0x91, 0x62, 0x65, 0x35, 0x3c, 0xab, 0xe1, 0x30, 0x94, 0xc1, 0xf4,
	0xb2, 0xc2, 0xc7, 0x95, 0xdc, 0xd6, 0xae, 0x07, 0xf6, 0x6c, 0x8f, 0x33, 0x46, 0x1b, 0x6c, 0xc7,
	0xfb, 0xf6, 0x13, 0x97, 0x32, 0x53, 0xea, 0x99, 0xc6, 0x5f, 0xa6, 0xa3, 0x61, 0x3d, 0xc5, 0xe5,
	0xe8, 0x2f, 0x18, 0x37, 0x5e, 0xa3, 0x51, 0x3c, 0x48, 0xc2, 0x3f, 0x93, 0x79, 0x5b, 0x1a, 0x6f,
	0x77, 0xd8, 0x01, 0xbe, 0xbe, 0x89, 0xa5, 0xca, 0x0f, 0x73, 0x75, 0x35, 0x7b, 0x6e, 0xcd, 0x4d,
	0xa5, 0x99, 0x76, 0xc3, 0xb8, 0x94, 0xa9, 0x74, 0xbd, 0xdf, 0x6c, 0x50, 0x29, 0x9b, 0x23, 0xe0,
	0x2d, 0x64, 0xcf, 0x04, 0x66, 0x1c, 0x77, 0x7b, 0x54, 0xfa, 0x5a, 0x6a, 0x3c, 0x55, 0xa8, 0x89,
	0x2a, 0x94, 0xee, 0x0a, 0xf3, 0x9a, 0xa8, 0x0e, 0x67, 0x6c, 0x34, 0xd8, 0x1e, 0x6f, 0x6c, 0x38,
	0x14, 0xfd, 0x0d, 0x9f, 0x16, 0xa2, 0x48, 0xb3, 0x54, 0x68, 0xb4, 0x37, 0xe3, 0xdb, 0x9b, 0xe9,
	0x93, 0x6c, 0x09, 0x9f, 0x7b, 0x8e, 0x4e, 0x74, 0x11, 0x43, 0x68, 0xf6, 0x97, 0x95, 0x28, 0x34,
	0xa6, 0xd6, 0x51, 0xc0, 0x5d, 0x8a, 0x6d, 0x21, 0x5c, 0x17, 0xa2, 0x54, 0xf7, 0x52, 0x9b, 0x5c,
	0xfd, 0x47, 0x41, 0xde, 0x3d, 0x8a, 0x1f, 0x30, 0x31, 0x66, 0xdd, 0x80, 0x47, 0x82, 0x7e, 0x87,
	0xc0, 0x00, 0x27, 0x5a, 0x87, 0x59, 0x02, 0xd3, 0xa3, 0x90, 0x2a, 0xdd, 0xba, 0x49, 0xaf, 0xee,
	0xff, 0xa3, 0x1b, 0xbf, 0x12, 0x77, 0xfa, 0x76, 0x64, 0x7f, 0xa1, 0xbf, 0xaf, 0x01, 0x00, 0x00,
	0xff, 0xff, 0xd4, 0x81, 0x8d, 0x15, 0x54, 0x03, 0x00, 0x00,
}
