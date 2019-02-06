// Code generated by protoc-gen-go. DO NOT EDIT.
// source: peer.proto

package peers

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type Peer struct {
	ID                   uint64   `protobuf:"varint,1,opt,name=ID,proto3" json:"ID,omitempty"`
	NetAddr              string   `protobuf:"bytes,2,opt,name=NetAddr,proto3" json:"NetAddr,omitempty"`
	PubKeyHex            string   `protobuf:"bytes,3,opt,name=PubKeyHex,proto3" json:"PubKeyHex,omitempty"`
	Used                 int64    `protobuf:"varint,4,opt,name=used,proto3" json:"used,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Peer) Reset()         { *m = Peer{} }
func (m *Peer) String() string { return proto.CompactTextString(m) }
func (*Peer) ProtoMessage()    {}
func (*Peer) Descriptor() ([]byte, []int) {
	return fileDescriptor_peer_3c6a3f501ad1abf2, []int{0}
}
func (m *Peer) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Peer.Unmarshal(m, b)
}
func (m *Peer) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Peer.Marshal(b, m, deterministic)
}
func (dst *Peer) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Peer.Merge(dst, src)
}
func (m *Peer) XXX_Size() int {
	return xxx_messageInfo_Peer.Size(m)
}
func (m *Peer) XXX_DiscardUnknown() {
	xxx_messageInfo_Peer.DiscardUnknown(m)
}

var xxx_messageInfo_Peer proto.InternalMessageInfo

func (m *Peer) GetID() uint64 {
	if m != nil {
		return m.ID
	}
	return 0
}

func (m *Peer) GetNetAddr() string {
	if m != nil {
		return m.NetAddr
	}
	return ""
}

func (m *Peer) GetPubKeyHex() string {
	if m != nil {
		return m.PubKeyHex
	}
	return ""
}

func (m *Peer) GetUsed() int64 {
	if m != nil {
		return m.Used
	}
	return 0
}

func init() {
	proto.RegisterType((*Peer)(nil), "peers.Peer")
}

func init() { proto.RegisterFile("peer.proto", fileDescriptor_peer_3c6a3f501ad1abf2) }

var fileDescriptor_peer_3c6a3f501ad1abf2 = []byte{
	// 125 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xe2, 0x2a, 0x48, 0x4d, 0x2d,
	0xd2, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0x62, 0x05, 0xb1, 0x8b, 0x95, 0x92, 0xb8, 0x58, 0x02,
	0x52, 0x53, 0x8b, 0x84, 0xf8, 0xb8, 0x98, 0x3c, 0x5d, 0x24, 0x18, 0x15, 0x18, 0x35, 0x58, 0x82,
	0x98, 0x3c, 0x5d, 0x84, 0x24, 0xb8, 0xd8, 0xfd, 0x52, 0x4b, 0x1c, 0x53, 0x52, 0x8a, 0x24, 0x98,
	0x14, 0x18, 0x35, 0x38, 0x83, 0x60, 0x5c, 0x21, 0x19, 0x2e, 0xce, 0x80, 0xd2, 0x24, 0xef, 0xd4,
	0x4a, 0x8f, 0xd4, 0x0a, 0x09, 0x66, 0xb0, 0x1c, 0x42, 0x40, 0x48, 0x88, 0x8b, 0xa5, 0xb4, 0x38,
	0x35, 0x45, 0x82, 0x45, 0x81, 0x51, 0x83, 0x39, 0x08, 0xcc, 0x4e, 0x62, 0x03, 0xdb, 0x68, 0x0c,
	0x08, 0x00, 0x00, 0xff, 0xff, 0x70, 0x30, 0x8e, 0xd3, 0x7f, 0x00, 0x00, 0x00,
}