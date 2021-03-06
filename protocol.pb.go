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
	Success              bool     `protobuf:"varint,2,opt,name=Success,proto3" json:"Success,omitempty"`
	ConflictTerm         int64    `protobuf:"varint,3,opt,name=ConflictTerm,proto3" json:"ConflictTerm,omitempty"`
	ConflictIndex        int64    `protobuf:"varint,4,opt,name=ConflictIndex,proto3" json:"ConflictIndex,omitempty"`
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

func (m *AppendEntriesResp) GetSuccess() bool {
	if m != nil {
		return m.Success
	}
	return false
}

func (m *AppendEntriesResp) GetConflictTerm() int64 {
	if m != nil {
		return m.ConflictTerm
	}
	return 0
}

func (m *AppendEntriesResp) GetConflictIndex() int64 {
	if m != nil {
		return m.ConflictIndex
	}
	return 0
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

type InstallSnapshotReq struct {
	Term                 uint64   `protobuf:"varint,1,opt,name=Term,proto3" json:"Term,omitempty"`
	LeaderName           string   `protobuf:"bytes,2,opt,name=LeaderName,proto3" json:"LeaderName,omitempty"`
	LastIncludedIndex    uint64   `protobuf:"varint,3,opt,name=LastIncludedIndex,proto3" json:"LastIncludedIndex,omitempty"`
	LastIncludedTerm     uint64   `protobuf:"varint,4,opt,name=LastIncludedTerm,proto3" json:"LastIncludedTerm,omitempty"`
	Offset               uint64   `protobuf:"varint,5,opt,name=Offset,proto3" json:"Offset,omitempty"`
	Data                 []byte   `protobuf:"bytes,6,opt,name=Data,proto3" json:"Data,omitempty"`
	Done                 bool     `protobuf:"varint,7,opt,name=Done,proto3" json:"Done,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *InstallSnapshotReq) Reset()         { *m = InstallSnapshotReq{} }
func (m *InstallSnapshotReq) String() string { return proto.CompactTextString(m) }
func (*InstallSnapshotReq) ProtoMessage()    {}
func (*InstallSnapshotReq) Descriptor() ([]byte, []int) {
	return fileDescriptor_2bc2336598a3f7e0, []int{5}
}

func (m *InstallSnapshotReq) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_InstallSnapshotReq.Unmarshal(m, b)
}
func (m *InstallSnapshotReq) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_InstallSnapshotReq.Marshal(b, m, deterministic)
}
func (m *InstallSnapshotReq) XXX_Merge(src proto.Message) {
	xxx_messageInfo_InstallSnapshotReq.Merge(m, src)
}
func (m *InstallSnapshotReq) XXX_Size() int {
	return xxx_messageInfo_InstallSnapshotReq.Size(m)
}
func (m *InstallSnapshotReq) XXX_DiscardUnknown() {
	xxx_messageInfo_InstallSnapshotReq.DiscardUnknown(m)
}

var xxx_messageInfo_InstallSnapshotReq proto.InternalMessageInfo

func (m *InstallSnapshotReq) GetTerm() uint64 {
	if m != nil {
		return m.Term
	}
	return 0
}

func (m *InstallSnapshotReq) GetLeaderName() string {
	if m != nil {
		return m.LeaderName
	}
	return ""
}

func (m *InstallSnapshotReq) GetLastIncludedIndex() uint64 {
	if m != nil {
		return m.LastIncludedIndex
	}
	return 0
}

func (m *InstallSnapshotReq) GetLastIncludedTerm() uint64 {
	if m != nil {
		return m.LastIncludedTerm
	}
	return 0
}

func (m *InstallSnapshotReq) GetOffset() uint64 {
	if m != nil {
		return m.Offset
	}
	return 0
}

func (m *InstallSnapshotReq) GetData() []byte {
	if m != nil {
		return m.Data
	}
	return nil
}

func (m *InstallSnapshotReq) GetDone() bool {
	if m != nil {
		return m.Done
	}
	return false
}

type InstallSnapshotResp struct {
	Term                 uint64   `protobuf:"varint,1,opt,name=Term,proto3" json:"Term,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *InstallSnapshotResp) Reset()         { *m = InstallSnapshotResp{} }
func (m *InstallSnapshotResp) String() string { return proto.CompactTextString(m) }
func (*InstallSnapshotResp) ProtoMessage()    {}
func (*InstallSnapshotResp) Descriptor() ([]byte, []int) {
	return fileDescriptor_2bc2336598a3f7e0, []int{6}
}

func (m *InstallSnapshotResp) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_InstallSnapshotResp.Unmarshal(m, b)
}
func (m *InstallSnapshotResp) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_InstallSnapshotResp.Marshal(b, m, deterministic)
}
func (m *InstallSnapshotResp) XXX_Merge(src proto.Message) {
	xxx_messageInfo_InstallSnapshotResp.Merge(m, src)
}
func (m *InstallSnapshotResp) XXX_Size() int {
	return xxx_messageInfo_InstallSnapshotResp.Size(m)
}
func (m *InstallSnapshotResp) XXX_DiscardUnknown() {
	xxx_messageInfo_InstallSnapshotResp.DiscardUnknown(m)
}

var xxx_messageInfo_InstallSnapshotResp proto.InternalMessageInfo

func (m *InstallSnapshotResp) GetTerm() uint64 {
	if m != nil {
		return m.Term
	}
	return 0
}

func init() {
	proto.RegisterType((*LogEntry)(nil), "LogEntry")
	proto.RegisterType((*AppendEntriesReq)(nil), "AppendEntriesReq")
	proto.RegisterType((*AppendEntriesResp)(nil), "AppendEntriesResp")
	proto.RegisterType((*RequestVoteReq)(nil), "RequestVoteReq")
	proto.RegisterType((*RequestVoteResp)(nil), "RequestVoteResp")
	proto.RegisterType((*InstallSnapshotReq)(nil), "InstallSnapshotReq")
	proto.RegisterType((*InstallSnapshotResp)(nil), "InstallSnapshotResp")
}

func init() { proto.RegisterFile("protocol.proto", fileDescriptor_2bc2336598a3f7e0) }

var fileDescriptor_2bc2336598a3f7e0 = []byte{
	// 445 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x9c, 0x93, 0xc1, 0x6e, 0xd4, 0x30,
	0x10, 0x86, 0xe5, 0xdd, 0x34, 0xdb, 0xce, 0x96, 0xd2, 0x1a, 0x84, 0x7c, 0x42, 0x91, 0xe1, 0xb0,
	0x20, 0xd4, 0x03, 0x3c, 0x01, 0x14, 0x54, 0xad, 0xb4, 0x02, 0xe4, 0x22, 0x0e, 0xdc, 0x4c, 0x3c,
	0x5b, 0x56, 0xca, 0xda, 0x69, 0xec, 0x45, 0xf0, 0x0c, 0x5c, 0x78, 0x39, 0x5e, 0x80, 0x27, 0x41,
	0x9e, 0x24, 0x2b, 0x87, 0xed, 0x5e, 0x7a, 0xf3, 0xff, 0x39, 0x99, 0xf9, 0xff, 0xc9, 0x04, 0x4e,
	0xea, 0xc6, 0x05, 0x57, 0xba, 0xea, 0x9c, 0x0e, 0xb2, 0x86, 0xc3, 0x85, 0xbb, 0x7e, 0x67, 0x43,
	0xf3, 0x93, 0x3f, 0x84, 0x83, 0xb9, 0x35, 0xf8, 0x43, 0xb0, 0x82, 0xcd, 0x32, 0xd5, 0x0a, 0xce,
	0x21, 0xfb, 0x84, 0xcd, 0x5a, 0x8c, 0x08, 0xd2, 0x99, 0x17, 0x30, 0xbd, 0x70, 0xeb, 0xb5, 0xb6,
	0xe6, 0xbd, 0x5e, 0xa3, 0x18, 0x17, 0x6c, 0x76, 0xa4, 0x52, 0xc4, 0x05, 0x4c, 0x3a, 0x29, 0xb2,
	0x82, 0xcd, 0x8e, 0x55, 0x2f, 0xe5, 0x1f, 0x06, 0xa7, 0xaf, 0xeb, 0x1a, 0xad, 0x89, 0x5d, 0x57,
	0xe8, 0x15, 0xde, 0x6c, 0x9b, 0xb0, 0xa4, 0xc9, 0x63, 0x80, 0x05, 0x6a, 0x83, 0x0d, 0xf5, 0x18,
	0x51, 0x8f, 0x84, 0x70, 0x09, 0xc7, 0x1f, 0x1b, 0xfc, 0xbe, 0x70, 0xd7, 0xad, 0xeb, 0x31, 0xbd,
	0x3b, 0x60, 0xd1, 0x68, 0xa7, 0xa9, 0x7c, 0x46, 0x8f, 0xa4, 0x28, 0x56, 0x69, 0x6b, 0x46, 0x7f,
	0xab, 0x20, 0x0e, 0xda, 0x2a, 0x29, 0xe3, 0x4f, 0x60, 0xd2, 0x79, 0x15, 0x79, 0x31, 0x9e, 0x4d,
	0x5f, 0x1e, 0x9d, 0xf7, 0x43, 0x53, 0xfd, 0x8d, 0xfc, 0xc5, 0xe0, 0xec, 0xbf, 0x5c, 0xbe, 0xbe,
	0x35, 0x98, 0x80, 0xc9, 0xd5, 0xa6, 0x2c, 0xd1, 0x7b, 0x4a, 0x75, 0xa8, 0x7a, 0x19, 0xcd, 0x5c,
	0x38, 0xbb, 0xac, 0x56, 0x65, 0xa0, 0xb7, 0x62, 0xa4, 0xb1, 0x1a, 0x30, 0xfe, 0x14, 0xee, 0xf5,
	0xba, 0xcd, 0x9d, 0xd1, 0x43, 0x43, 0x28, 0x7f, 0x33, 0x38, 0x51, 0x78, 0xb3, 0x41, 0x1f, 0x3e,
	0xbb, 0x80, 0xfb, 0x66, 0x1c, 0xd3, 0x6b, 0x1f, 0xb6, 0x33, 0x1c, 0x75, 0xe9, 0x13, 0x16, 0x67,
	0xd8, 0xe9, 0xad, 0xa7, 0x4c, 0xa5, 0x88, 0x2c, 0x69, 0x6b, 0x56, 0x46, 0x07, 0xa4, 0x8f, 0x95,
	0xd1, 0xc7, 0x1a, 0x42, 0x79, 0x09, 0xf7, 0x07, 0x8e, 0xf6, 0x4c, 0xa7, 0x80, 0x69, 0xbc, 0xbf,
	0x6c, 0xb4, 0x0d, 0x68, 0xba, 0x09, 0xa5, 0x48, 0xfe, 0x65, 0xc0, 0xe7, 0xd6, 0x07, 0x5d, 0x55,
	0x57, 0x56, 0xd7, 0xfe, 0x9b, 0x0b, 0x77, 0xdd, 0xa1, 0x17, 0x70, 0x16, 0x83, 0xcc, 0x6d, 0x59,
	0x6d, 0x0c, 0x9a, 0x74, 0x91, 0x76, 0x2f, 0xf8, 0x73, 0x38, 0x4d, 0x61, 0xb2, 0x52, 0x3b, 0x9c,
	0x3f, 0x82, 0xfc, 0xc3, 0x72, 0xe9, 0xb1, 0xdf, 0xa8, 0x4e, 0x45, 0x97, 0x6f, 0x75, 0xd0, 0x22,
	0xa7, 0xbf, 0x82, 0xce, 0xc4, 0x9c, 0x45, 0x31, 0xa1, 0xac, 0x74, 0x96, 0xcf, 0xe0, 0xc1, 0x4e,
	0xc6, 0xdb, 0x27, 0xf6, 0x26, 0xff, 0x92, 0x35, 0x7a, 0x19, 0xbe, 0xe6, 0xf4, 0x4b, 0xbf, 0xfa,
	0x17, 0x00, 0x00, 0xff, 0xff, 0xfe, 0x43, 0xca, 0xc3, 0xe4, 0x03, 0x00, 0x00,
}
