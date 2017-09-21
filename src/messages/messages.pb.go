// Code generated by protoc-gen-go. DO NOT EDIT.
// source: messages.proto

/*
Package messages is a generated protocol buffer package.

It is generated from these files:
	messages.proto

It has these top-level messages:
	Message
	Request
	Response
	FindClosest
	Ping
	FindData
	Store
*/
package messages

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

type Message_Type int32

const (
	Message_REQUEST  Message_Type = 0
	Message_RESPONSE Message_Type = 1
)

var Message_Type_name = map[int32]string{
	0: "REQUEST",
	1: "RESPONSE",
}
var Message_Type_value = map[string]int32{
	"REQUEST":  0,
	"RESPONSE": 1,
}

func (x Message_Type) String() string {
	return proto.EnumName(Message_Type_name, int32(x))
}
func (Message_Type) EnumDescriptor() ([]byte, []int) { return fileDescriptor0, []int{0, 0} }

type Request_Type int32

const (
	Request_PING     Request_Type = 0
	Request_FINDNODE Request_Type = 1
	Request_FINDDATA Request_Type = 2
	Request_STORE    Request_Type = 3
)

var Request_Type_name = map[int32]string{
	0: "PING",
	1: "FINDNODE",
	2: "FINDDATA",
	3: "STORE",
}
var Request_Type_value = map[string]int32{
	"PING":     0,
	"FINDNODE": 1,
	"FINDDATA": 2,
	"STORE":    3,
}

func (x Request_Type) String() string {
	return proto.EnumName(Request_Type_name, int32(x))
}
func (Request_Type) EnumDescriptor() ([]byte, []int) { return fileDescriptor0, []int{1, 0} }

type Response_Type int32

const (
	Response_PING Response_Type = 0
)

var Response_Type_name = map[int32]string{
	0: "PING",
}
var Response_Type_value = map[string]int32{
	"PING": 0,
}

func (x Response_Type) String() string {
	return proto.EnumName(Response_Type_name, int32(x))
}
func (Response_Type) EnumDescriptor() ([]byte, []int) { return fileDescriptor0, []int{2, 0} }

type Message struct {
	Type          Message_Type `protobuf:"varint,1,opt,name=type,enum=messages.Message_Type" json:"type,omitempty"`
	SenderID      string       `protobuf:"bytes,2,opt,name=senderID" json:"senderID,omitempty"`
	SenderAddress string       `protobuf:"bytes,3,opt,name=senderAddress" json:"senderAddress,omitempty"`
	Request       *Request     `protobuf:"bytes,4,opt,name=Request" json:"Request,omitempty"`
	Response      *Response    `protobuf:"bytes,5,opt,name=Response" json:"Response,omitempty"`
}

func (m *Message) Reset()                    { *m = Message{} }
func (m *Message) String() string            { return proto.CompactTextString(m) }
func (*Message) ProtoMessage()               {}
func (*Message) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *Message) GetType() Message_Type {
	if m != nil {
		return m.Type
	}
	return Message_REQUEST
}

func (m *Message) GetSenderID() string {
	if m != nil {
		return m.SenderID
	}
	return ""
}

func (m *Message) GetSenderAddress() string {
	if m != nil {
		return m.SenderAddress
	}
	return ""
}

func (m *Message) GetRequest() *Request {
	if m != nil {
		return m.Request
	}
	return nil
}

func (m *Message) GetResponse() *Response {
	if m != nil {
		return m.Response
	}
	return nil
}

type Request struct {
	Type Request_Type `protobuf:"varint,1,opt,name=type,enum=messages.Request_Type" json:"type,omitempty"`
	ID   string       `protobuf:"bytes,2,opt,name=ID" json:"ID,omitempty"`
}

func (m *Request) Reset()                    { *m = Request{} }
func (m *Request) String() string            { return proto.CompactTextString(m) }
func (*Request) ProtoMessage()               {}
func (*Request) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *Request) GetType() Request_Type {
	if m != nil {
		return m.Type
	}
	return Request_PING
}

func (m *Request) GetID() string {
	if m != nil {
		return m.ID
	}
	return ""
}

type Response struct {
	Type Response_Type `protobuf:"varint,1,opt,name=type,enum=messages.Response_Type" json:"type,omitempty"`
}

func (m *Response) Reset()                    { *m = Response{} }
func (m *Response) String() string            { return proto.CompactTextString(m) }
func (*Response) ProtoMessage()               {}
func (*Response) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

func (m *Response) GetType() Response_Type {
	if m != nil {
		return m.Type
	}
	return Response_PING
}

type FindClosest struct {
	Find string `protobuf:"bytes,1,opt,name=find" json:"find,omitempty"`
}

func (m *FindClosest) Reset()                    { *m = FindClosest{} }
func (m *FindClosest) String() string            { return proto.CompactTextString(m) }
func (*FindClosest) ProtoMessage()               {}
func (*FindClosest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{3} }

func (m *FindClosest) GetFind() string {
	if m != nil {
		return m.Find
	}
	return ""
}

type Ping struct {
}

func (m *Ping) Reset()                    { *m = Ping{} }
func (m *Ping) String() string            { return proto.CompactTextString(m) }
func (*Ping) ProtoMessage()               {}
func (*Ping) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{4} }

type FindData struct {
	Find string `protobuf:"bytes,1,opt,name=find" json:"find,omitempty"`
}

func (m *FindData) Reset()                    { *m = FindData{} }
func (m *FindData) String() string            { return proto.CompactTextString(m) }
func (*FindData) ProtoMessage()               {}
func (*FindData) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{5} }

func (m *FindData) GetFind() string {
	if m != nil {
		return m.Find
	}
	return ""
}

type Store struct {
	Store string `protobuf:"bytes,1,opt,name=store" json:"store,omitempty"`
}

func (m *Store) Reset()                    { *m = Store{} }
func (m *Store) String() string            { return proto.CompactTextString(m) }
func (*Store) ProtoMessage()               {}
func (*Store) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{6} }

func (m *Store) GetStore() string {
	if m != nil {
		return m.Store
	}
	return ""
}

func init() {
	proto.RegisterType((*Message)(nil), "messages.message")
	proto.RegisterType((*Request)(nil), "messages.request")
	proto.RegisterType((*Response)(nil), "messages.response")
	proto.RegisterType((*FindClosest)(nil), "messages.findClosest")
	proto.RegisterType((*Ping)(nil), "messages.ping")
	proto.RegisterType((*FindData)(nil), "messages.findData")
	proto.RegisterType((*Store)(nil), "messages.store")
	proto.RegisterEnum("messages.Message_Type", Message_Type_name, Message_Type_value)
	proto.RegisterEnum("messages.Request_Type", Request_Type_name, Request_Type_value)
	proto.RegisterEnum("messages.Response_Type", Response_Type_name, Response_Type_value)
}

func init() { proto.RegisterFile("messages.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 334 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0x74, 0x92, 0x4f, 0x4f, 0xb3, 0x40,
	0x18, 0xc4, 0xdf, 0xa5, 0xb4, 0x6c, 0x9f, 0xbe, 0x36, 0xf8, 0xc4, 0x28, 0x31, 0xd1, 0xb4, 0x1b,
	0x0f, 0x8d, 0x4d, 0x38, 0xd4, 0x83, 0xe7, 0x46, 0xd0, 0x70, 0xa1, 0x75, 0xc1, 0x0f, 0x50, 0xc3,
	0xda, 0x34, 0x51, 0x40, 0x16, 0x0f, 0xbd, 0xf8, 0xb5, 0xbd, 0x1a, 0x96, 0x05, 0x5b, 0xff, 0xdc,
	0x66, 0x78, 0x7e, 0x30, 0x33, 0x09, 0x30, 0x7c, 0x11, 0x52, 0xae, 0xd6, 0x42, 0xba, 0x79, 0x91,
	0x95, 0x19, 0xd2, 0xc6, 0xb3, 0x0f, 0x02, 0x96, 0x36, 0x78, 0x09, 0x66, 0xb9, 0xcd, 0x85, 0x43,
	0x46, 0x64, 0x32, 0x9c, 0x1d, 0xbb, 0xed, 0x4b, 0x5a, 0xb8, 0xf1, 0x36, 0x17, 0x5c, 0x31, 0x78,
	0x0a, 0x54, 0x8a, 0x34, 0x11, 0x45, 0xe0, 0x39, 0xc6, 0x88, 0x4c, 0xfa, 0xbc, 0xf5, 0x78, 0x01,
	0x07, 0xb5, 0x9e, 0x27, 0x49, 0x21, 0xa4, 0x74, 0x3a, 0x0a, 0xd8, 0x7f, 0x88, 0x53, 0xb0, 0xb8,
	0x78, 0x7d, 0x13, 0xb2, 0x74, 0xcc, 0x11, 0x99, 0x0c, 0x66, 0x87, 0x5f, 0x81, 0x45, 0x7d, 0xe0,
	0x0d, 0x81, 0x2e, 0x50, 0x2e, 0x64, 0x9e, 0xa5, 0x52, 0x38, 0x5d, 0x45, 0xe3, 0x2e, 0x5d, 0x5f,
	0x78, 0xcb, 0xb0, 0x31, 0x98, 0x55, 0x59, 0x1c, 0x80, 0xc5, 0xfd, 0xfb, 0x07, 0x3f, 0x8a, 0xed,
	0x7f, 0xf8, 0x1f, 0x28, 0xf7, 0xa3, 0xe5, 0x22, 0x8c, 0x7c, 0x9b, 0xb0, 0x77, 0xb0, 0x74, 0xcc,
	0xdf, 0xc3, 0x35, 0xb0, 0x3b, 0x7c, 0x08, 0x46, 0x3b, 0xd9, 0x08, 0x3c, 0x76, 0xad, 0x93, 0x28,
	0x98, 0xcb, 0x20, 0xbc, 0xab, 0x63, 0x6e, 0x83, 0xd0, 0x0b, 0x17, 0x9e, 0x6f, 0x93, 0xc6, 0x79,
	0xf3, 0x78, 0x6e, 0x1b, 0xd8, 0x87, 0x6e, 0x14, 0x2f, 0xb8, 0x6f, 0x77, 0x58, 0x00, 0xb4, 0x29,
	0x8e, 0xd3, 0xbd, 0x02, 0x27, 0x3f, 0xa7, 0xed, 0x34, 0x60, 0xf6, 0xf7, 0x44, 0x36, 0x86, 0xc1,
	0xd3, 0x26, 0x4d, 0x6e, 0x9e, 0x33, 0x59, 0xcd, 0x41, 0x30, 0x2b, 0xab, 0xbe, 0xd6, 0xe7, 0x4a,
	0xb3, 0x1e, 0x98, 0xf9, 0x26, 0x5d, 0xb3, 0x73, 0xa0, 0x95, 0xf7, 0x56, 0xe5, 0xea, 0x57, 0xee,
	0x0c, 0xba, 0xb2, 0xcc, 0x0a, 0x81, 0x47, 0x5a, 0xe8, 0x6b, 0x6d, 0x1e, 0x7b, 0xea, 0xff, 0xb9,
	0xfa, 0x0c, 0x00, 0x00, 0xff, 0xff, 0x33, 0x4f, 0x7c, 0x1e, 0x51, 0x02, 0x00, 0x00,
}
