// Code generated by protoc-gen-go. DO NOT EDIT.
// source: block.proto

package poset

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
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type BlockBody struct {
	Index                int64    `protobuf:"varint,1,opt,name=Index,proto3" json:"Index,omitempty"`
	RoundReceived        int64    `protobuf:"varint,2,opt,name=RoundReceived,proto3" json:"RoundReceived,omitempty"`
	Transactions         [][]byte `protobuf:"bytes,5,rep,name=Transactions,proto3" json:"Transactions,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *BlockBody) Reset()         { *m = BlockBody{} }
func (m *BlockBody) String() string { return proto.CompactTextString(m) }
func (*BlockBody) ProtoMessage()    {}
func (*BlockBody) Descriptor() ([]byte, []int) {
	return fileDescriptor_8e550b1f5926e92d, []int{0}
}

func (m *BlockBody) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_BlockBody.Unmarshal(m, b)
}
func (m *BlockBody) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_BlockBody.Marshal(b, m, deterministic)
}
func (m *BlockBody) XXX_Merge(src proto.Message) {
	xxx_messageInfo_BlockBody.Merge(m, src)
}
func (m *BlockBody) XXX_Size() int {
	return xxx_messageInfo_BlockBody.Size(m)
}
func (m *BlockBody) XXX_DiscardUnknown() {
	xxx_messageInfo_BlockBody.DiscardUnknown(m)
}

var xxx_messageInfo_BlockBody proto.InternalMessageInfo

func (m *BlockBody) GetIndex() int64 {
	if m != nil {
		return m.Index
	}
	return 0
}

func (m *BlockBody) GetRoundReceived() int64 {
	if m != nil {
		return m.RoundReceived
	}
	return 0
}

func (m *BlockBody) GetTransactions() [][]byte {
	if m != nil {
		return m.Transactions
	}
	return nil
}

type WireBlockSignature struct {
	Index                int64    `protobuf:"varint,1,opt,name=Index,proto3" json:"Index,omitempty"`
	Signature            string   `protobuf:"bytes,2,opt,name=Signature,proto3" json:"Signature,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *WireBlockSignature) Reset()         { *m = WireBlockSignature{} }
func (m *WireBlockSignature) String() string { return proto.CompactTextString(m) }
func (*WireBlockSignature) ProtoMessage()    {}
func (*WireBlockSignature) Descriptor() ([]byte, []int) {
	return fileDescriptor_8e550b1f5926e92d, []int{1}
}

func (m *WireBlockSignature) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_WireBlockSignature.Unmarshal(m, b)
}
func (m *WireBlockSignature) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_WireBlockSignature.Marshal(b, m, deterministic)
}
func (m *WireBlockSignature) XXX_Merge(src proto.Message) {
	xxx_messageInfo_WireBlockSignature.Merge(m, src)
}
func (m *WireBlockSignature) XXX_Size() int {
	return xxx_messageInfo_WireBlockSignature.Size(m)
}
func (m *WireBlockSignature) XXX_DiscardUnknown() {
	xxx_messageInfo_WireBlockSignature.DiscardUnknown(m)
}

var xxx_messageInfo_WireBlockSignature proto.InternalMessageInfo

func (m *WireBlockSignature) GetIndex() int64 {
	if m != nil {
		return m.Index
	}
	return 0
}

func (m *WireBlockSignature) GetSignature() string {
	if m != nil {
		return m.Signature
	}
	return ""
}

type Block struct {
	Body                 *BlockBody        `protobuf:"bytes,1,opt,name=Body,proto3" json:"Body,omitempty"`
	Signatures           map[string]string `protobuf:"bytes,2,rep,name=Signatures,proto3" json:"Signatures,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	Hash                 []byte            `protobuf:"bytes,3,opt,name=hash,proto3" json:"hash,omitempty"`
	Hex                  string            `protobuf:"bytes,4,opt,name=hex,proto3" json:"hex,omitempty"`
	StateHash            []byte            `protobuf:"bytes,5,opt,name=StateHash,proto3" json:"StateHash,omitempty"`
	FrameHash            []byte            `protobuf:"bytes,6,opt,name=FrameHash,proto3" json:"FrameHash,omitempty"`
	CreatedTime          int64             `protobuf:"varint,7,opt,name=CreatedTime,proto3" json:"CreatedTime,omitempty"`
	XXX_NoUnkeyedLiteral struct{}          `json:"-"`
	XXX_unrecognized     []byte            `json:"-"`
	XXX_sizecache        int32             `json:"-"`
}

func (m *Block) Reset()         { *m = Block{} }
func (m *Block) String() string { return proto.CompactTextString(m) }
func (*Block) ProtoMessage()    {}
func (*Block) Descriptor() ([]byte, []int) {
	return fileDescriptor_8e550b1f5926e92d, []int{2}
}

func (m *Block) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Block.Unmarshal(m, b)
}
func (m *Block) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Block.Marshal(b, m, deterministic)
}
func (m *Block) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Block.Merge(m, src)
}
func (m *Block) XXX_Size() int {
	return xxx_messageInfo_Block.Size(m)
}
func (m *Block) XXX_DiscardUnknown() {
	xxx_messageInfo_Block.DiscardUnknown(m)
}

var xxx_messageInfo_Block proto.InternalMessageInfo

func (m *Block) GetBody() *BlockBody {
	if m != nil {
		return m.Body
	}
	return nil
}

func (m *Block) GetSignatures() map[string]string {
	if m != nil {
		return m.Signatures
	}
	return nil
}

func (m *Block) GetHash() []byte {
	if m != nil {
		return m.Hash
	}
	return nil
}

func (m *Block) GetHex() string {
	if m != nil {
		return m.Hex
	}
	return ""
}

func (m *Block) GetStateHash() []byte {
	if m != nil {
		return m.StateHash
	}
	return nil
}

func (m *Block) GetFrameHash() []byte {
	if m != nil {
		return m.FrameHash
	}
	return nil
}

func (m *Block) GetCreatedTime() int64 {
	if m != nil {
		return m.CreatedTime
	}
	return 0
}

func init() {
	proto.RegisterType((*BlockBody)(nil), "poset.BlockBody")
	proto.RegisterType((*WireBlockSignature)(nil), "poset.WireBlockSignature")
	proto.RegisterType((*Block)(nil), "poset.Block")
	proto.RegisterMapType((map[string]string)(nil), "poset.Block.SignaturesEntry")
}

func init() { proto.RegisterFile("block.proto", fileDescriptor_8e550b1f5926e92d) }

var fileDescriptor_8e550b1f5926e92d = []byte{
	// 306 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x74, 0x91, 0x41, 0x4b, 0xc3, 0x40,
	0x10, 0x85, 0x49, 0xd2, 0x54, 0x32, 0xa9, 0x58, 0x16, 0x0f, 0x8b, 0xf4, 0x10, 0x42, 0x0f, 0x39,
	0xe5, 0x50, 0x2f, 0x22, 0x7a, 0xa9, 0x28, 0xf5, 0xba, 0x16, 0x3c, 0x6f, 0x9b, 0xc1, 0x86, 0xb6,
	0xbb, 0x65, 0xb3, 0x29, 0xed, 0x9f, 0xf2, 0x37, 0xca, 0x4e, 0x30, 0x49, 0x05, 0x6f, 0xb3, 0xdf,
	0x7b, 0x93, 0x79, 0x33, 0x81, 0x78, 0xb5, 0xd3, 0xeb, 0x6d, 0x7e, 0x30, 0xda, 0x6a, 0x16, 0x1e,
	0x74, 0x85, 0x36, 0xdd, 0x42, 0x34, 0x77, 0x74, 0xae, 0x8b, 0x33, 0xbb, 0x85, 0xf0, 0x5d, 0x15,
	0x78, 0xe2, 0x5e, 0xe2, 0x65, 0x81, 0x68, 0x1e, 0x6c, 0x0a, 0xd7, 0x42, 0xd7, 0xaa, 0x10, 0xb8,
	0xc6, 0xf2, 0x88, 0x05, 0xf7, 0x49, 0xbd, 0x84, 0x2c, 0x85, 0xd1, 0xd2, 0x48, 0x55, 0xc9, 0xb5,
	0x2d, 0xb5, 0xaa, 0x78, 0x98, 0x04, 0xd9, 0x48, 0x5c, 0xb0, 0x74, 0x01, 0xec, 0xb3, 0x34, 0x48,
	0x03, 0x3f, 0xca, 0x2f, 0x25, 0x6d, 0x6d, 0xf0, 0x9f, 0xa9, 0x13, 0x88, 0x5a, 0x0b, 0x4d, 0x8c,
	0x44, 0x07, 0xd2, 0x6f, 0x1f, 0x42, 0xfa, 0x0c, 0x9b, 0xc2, 0xc0, 0x65, 0xa7, 0xe6, 0x78, 0x36,
	0xce, 0x69, 0xad, 0xbc, 0xdd, 0x49, 0x90, 0xca, 0x9e, 0x00, 0xda, 0xe6, 0x8a, 0xfb, 0x49, 0x90,
	0xc5, 0xb3, 0x49, 0xdf, 0x9b, 0x77, 0xf2, 0xab, 0xb2, 0xe6, 0x2c, 0x7a, 0x7e, 0xc6, 0x60, 0xb0,
	0x91, 0xd5, 0x86, 0x07, 0x89, 0x97, 0x8d, 0x04, 0xd5, 0x6c, 0x0c, 0xc1, 0x06, 0x4f, 0x7c, 0x40,
	0xc9, 0x5c, 0x49, 0x89, 0xad, 0xb4, 0xb8, 0x70, 0xd6, 0x90, 0xac, 0x1d, 0x70, 0xea, 0x9b, 0x91,
	0xfb, 0x46, 0x1d, 0x36, 0x6a, 0x0b, 0x58, 0x02, 0xf1, 0x8b, 0x41, 0x69, 0xb1, 0x58, 0x96, 0x7b,
	0xe4, 0x57, 0x74, 0x89, 0x3e, 0xba, 0x7b, 0x86, 0x9b, 0x3f, 0x11, 0x5d, 0x84, 0x2d, 0x36, 0x9b,
	0x47, 0xc2, 0x95, 0xee, 0x94, 0x47, 0xb9, 0xab, 0x7f, 0x0f, 0xd6, 0x3c, 0x1e, 0xfd, 0x07, 0x6f,
	0x35, 0xa4, 0xbf, 0x7e, 0xff, 0x13, 0x00, 0x00, 0xff, 0xff, 0x10, 0xb8, 0x9f, 0xf4, 0x04, 0x02,
	0x00, 0x00,
}
