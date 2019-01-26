// Code generated by protoc-gen-go. DO NOT EDIT.
// source: transaction_receipts.proto

package template_pb

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

type Receipt struct {
	Date                 int64     `protobuf:"varint,1,opt,name=date,proto3" json:"date,omitempty"`
	StateAddress         string    `protobuf:"bytes,2,opt,name=state_address,json=stateAddress,proto3" json:"state_address,omitempty"`
	RpcMethod            Method    `protobuf:"varint,3,opt,name=rpc_method,json=rpcMethod,proto3,enum=template_pb.Method" json:"rpc_method,omitempty"`
	Template             *Template `protobuf:"bytes,4,opt,name=template,proto3" json:"template,omitempty"`
	XXX_NoUnkeyedLiteral struct{}  `json:"-"`
	XXX_unrecognized     []byte    `json:"-"`
	XXX_sizecache        int32     `json:"-"`
}

func (m *Receipt) Reset()         { *m = Receipt{} }
func (m *Receipt) String() string { return proto.CompactTextString(m) }
func (*Receipt) ProtoMessage()    {}
func (*Receipt) Descriptor() ([]byte, []int) {
	return fileDescriptor_transaction_receipts_89072fd7395da3d7, []int{0}
}
func (m *Receipt) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Receipt.Unmarshal(m, b)
}
func (m *Receipt) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Receipt.Marshal(b, m, deterministic)
}
func (dst *Receipt) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Receipt.Merge(dst, src)
}
func (m *Receipt) XXX_Size() int {
	return xxx_messageInfo_Receipt.Size(m)
}
func (m *Receipt) XXX_DiscardUnknown() {
	xxx_messageInfo_Receipt.DiscardUnknown(m)
}

var xxx_messageInfo_Receipt proto.InternalMessageInfo

func (m *Receipt) GetDate() int64 {
	if m != nil {
		return m.Date
	}
	return 0
}

func (m *Receipt) GetStateAddress() string {
	if m != nil {
		return m.StateAddress
	}
	return ""
}

func (m *Receipt) GetRpcMethod() Method {
	if m != nil {
		return m.RpcMethod
	}
	return Method_CREATE
}

func (m *Receipt) GetTemplate() *Template {
	if m != nil {
		return m.Template
	}
	return nil
}

func init() {
	proto.RegisterType((*Receipt)(nil), "template_pb.Receipt")
}

func init() {
	proto.RegisterFile("transaction_receipts.proto", fileDescriptor_transaction_receipts_89072fd7395da3d7)
}

var fileDescriptor_transaction_receipts_89072fd7395da3d7 = []byte{
	// 200 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x54, 0x4f, 0x3d, 0x6a, 0xc3, 0x30,
	0x14, 0x46, 0xb5, 0x69, 0x6b, 0xb9, 0xf6, 0xa0, 0x52, 0x10, 0x9e, 0x44, 0xbb, 0x68, 0x32, 0xd4,
	0x3d, 0x41, 0xb3, 0x67, 0x11, 0x99, 0xb2, 0x18, 0x59, 0x12, 0xc4, 0x60, 0x5b, 0x42, 0x7a, 0x4b,
	0x8e, 0x94, 0x5b, 0x86, 0x48, 0x71, 0x70, 0xb6, 0xef, 0xf7, 0xf1, 0x3e, 0xdc, 0x80, 0x97, 0x4b,
	0x90, 0x0a, 0x46, 0xbb, 0xf4, 0xde, 0x28, 0x33, 0x3a, 0x08, 0xad, 0xf3, 0x16, 0x2c, 0x29, 0xc1,
	0xcc, 0x6e, 0x92, 0x60, 0x7a, 0x37, 0x34, 0xf5, 0x4a, 0x92, 0xd9, 0x54, 0x4e, 0x9e, 0x27, 0x2b,
	0x75, 0xa2, 0xdf, 0x17, 0x84, 0xdf, 0x44, 0xaa, 0x13, 0x82, 0x73, 0x2d, 0xc1, 0x50, 0xc4, 0x10,
	0xcf, 0x44, 0xc4, 0xe4, 0x07, 0x57, 0x01, 0x6e, 0xa7, 0xa4, 0xd6, 0xde, 0x84, 0x40, 0x5f, 0x18,
	0xe2, 0x85, 0xf8, 0x88, 0xe2, 0x7f, 0xd2, 0x48, 0x87, 0xb1, 0x77, 0xaa, 0x9f, 0x0d, 0x9c, 0xac,
	0xa6, 0x19, 0x43, 0xbc, 0xee, 0x3e, 0xdb, 0xcd, 0x17, 0xed, 0x3e, 0x5a, 0xa2, 0xf0, 0x4e, 0x25,
	0x48, 0x7e, 0xf1, 0xfb, 0x1a, 0xa0, 0x39, 0x43, 0xbc, 0xec, 0xbe, 0x9e, 0x1a, 0x87, 0x3b, 0x16,
	0x8f, 0xd8, 0xae, 0x3a, 0x6e, 0x97, 0x0d, 0xaf, 0x71, 0xc1, 0xdf, 0x35, 0x00, 0x00, 0xff, 0xff,
	0x4e, 0xd2, 0x9e, 0x9d, 0x0b, 0x01, 0x00, 0x00,
}