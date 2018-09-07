// Code generated by protoc-gen-go. DO NOT EDIT.
// source: crdexample.proto

package crdexample

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

// +k8s:openapi-gen=true
type CrdExample struct {
	Name                 string                        `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Uuid                 string                        `protobuf:"bytes,2,opt,name=uuid,proto3" json:"uuid,omitempty"`
	Repeats              []*CrdExample_CrdExampleEmbed `protobuf:"bytes,4,rep,name=repeats,proto3" json:"repeats,omitempty"`
	XXX_NoUnkeyedLiteral struct{}                      `json:"-"`
	XXX_unrecognized     []byte                        `json:"-"`
	XXX_sizecache        int32                         `json:"-"`
}

func (m *CrdExample) Reset()         { *m = CrdExample{} }
func (m *CrdExample) String() string { return proto.CompactTextString(m) }
func (*CrdExample) ProtoMessage()    {}
func (*CrdExample) Descriptor() ([]byte, []int) {
	return fileDescriptor_crdexample_74da81f9cdfe3757, []int{0}
}
func (m *CrdExample) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_CrdExample.Unmarshal(m, b)
}
func (m *CrdExample) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_CrdExample.Marshal(b, m, deterministic)
}
func (dst *CrdExample) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CrdExample.Merge(dst, src)
}
func (m *CrdExample) XXX_Size() int {
	return xxx_messageInfo_CrdExample.Size(m)
}
func (m *CrdExample) XXX_DiscardUnknown() {
	xxx_messageInfo_CrdExample.DiscardUnknown(m)
}

var xxx_messageInfo_CrdExample proto.InternalMessageInfo

func (m *CrdExample) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *CrdExample) GetUuid() string {
	if m != nil {
		return m.Uuid
	}
	return ""
}

func (m *CrdExample) GetRepeats() []*CrdExample_CrdExampleEmbed {
	if m != nil {
		return m.Repeats
	}
	return nil
}

type CrdExample_CrdExampleEmbed struct {
	Name                 string   `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Uuid                 string   `protobuf:"bytes,2,opt,name=uuid,proto3" json:"uuid,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *CrdExample_CrdExampleEmbed) Reset()         { *m = CrdExample_CrdExampleEmbed{} }
func (m *CrdExample_CrdExampleEmbed) String() string { return proto.CompactTextString(m) }
func (*CrdExample_CrdExampleEmbed) ProtoMessage()    {}
func (*CrdExample_CrdExampleEmbed) Descriptor() ([]byte, []int) {
	return fileDescriptor_crdexample_74da81f9cdfe3757, []int{0, 0}
}
func (m *CrdExample_CrdExampleEmbed) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_CrdExample_CrdExampleEmbed.Unmarshal(m, b)
}
func (m *CrdExample_CrdExampleEmbed) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_CrdExample_CrdExampleEmbed.Marshal(b, m, deterministic)
}
func (dst *CrdExample_CrdExampleEmbed) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CrdExample_CrdExampleEmbed.Merge(dst, src)
}
func (m *CrdExample_CrdExampleEmbed) XXX_Size() int {
	return xxx_messageInfo_CrdExample_CrdExampleEmbed.Size(m)
}
func (m *CrdExample_CrdExampleEmbed) XXX_DiscardUnknown() {
	xxx_messageInfo_CrdExample_CrdExampleEmbed.DiscardUnknown(m)
}

var xxx_messageInfo_CrdExample_CrdExampleEmbed proto.InternalMessageInfo

func (m *CrdExample_CrdExampleEmbed) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *CrdExample_CrdExampleEmbed) GetUuid() string {
	if m != nil {
		return m.Uuid
	}
	return ""
}

func init() {
	proto.RegisterType((*CrdExample)(nil), "crdexample.CrdExample")
	proto.RegisterType((*CrdExample_CrdExampleEmbed)(nil), "crdexample.CrdExample.CrdExampleEmbed")
}

func init() { proto.RegisterFile("crdexample.proto", fileDescriptor_crdexample_74da81f9cdfe3757) }

var fileDescriptor_crdexample_74da81f9cdfe3757 = []byte{
	// 133 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x12, 0x48, 0x2e, 0x4a, 0x49,
	0xad, 0x48, 0xcc, 0x2d, 0xc8, 0x49, 0xd5, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0xe2, 0x42, 0x88,
	0x28, 0x6d, 0x64, 0xe4, 0xe2, 0x72, 0x2e, 0x4a, 0x71, 0x85, 0x70, 0x85, 0x84, 0xb8, 0x58, 0xf2,
	0x12, 0x73, 0x53, 0x25, 0x18, 0x15, 0x18, 0x35, 0x38, 0x83, 0xc0, 0x6c, 0x90, 0x58, 0x69, 0x69,
	0x66, 0x8a, 0x04, 0x13, 0x44, 0x0c, 0xc4, 0x16, 0x72, 0xe0, 0x62, 0x2f, 0x4a, 0x2d, 0x48, 0x4d,
	0x2c, 0x29, 0x96, 0x60, 0x51, 0x60, 0xd6, 0xe0, 0x36, 0x52, 0xd3, 0x43, 0xb2, 0x06, 0x61, 0x20,
	0x12, 0xd3, 0x35, 0x37, 0x29, 0x35, 0x25, 0x08, 0xa6, 0x4d, 0xca, 0x92, 0x8b, 0x1f, 0x4d, 0x8e,
	0x58, 0xcb, 0x93, 0xd8, 0xc0, 0xde, 0x30, 0x06, 0x04, 0x00, 0x00, 0xff, 0xff, 0xba, 0xcd, 0xfc,
	0x51, 0xda, 0x00, 0x00, 0x00,
}
