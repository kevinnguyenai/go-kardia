// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: kardiachain/dualnode/types.proto

package kardiachain_dualnode

import (
	fmt "fmt"
	_ "github.com/gogo/protobuf/gogoproto"
	proto "github.com/gogo/protobuf/proto"
	io "io"
	math "math"
	math_bits "math/bits"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.GoGoProtoPackageIsVersion3 // please upgrade the proto package

type ProposalType int32

const (
	ProposalUnLock   ProposalType = 0
	ProposalTransfer ProposalType = 1
	ProposalMintLock ProposalType = 2
)

var ProposalType_name = map[int32]string{
	0: "PROPOSAL_UNLOCK",
	1: "PROPOSAL_TRANSFER",
	2: "PROPOSAL_MIN_LOCK",
}

var ProposalType_value = map[string]int32{
	"PROPOSAL_UNLOCK":   0,
	"PROPOSAL_TRANSFER": 1,
	"PROPOSAL_MIN_LOCK": 2,
}

func (x ProposalType) String() string {
	return proto.EnumName(ProposalType_name, int32(x))
}

func (ProposalType) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_fdeddabbf1b9645a, []int{0}
}

type Vote struct {
	Hash             string `protobuf:"bytes,1,opt,name=hash,proto3" json:"hash,omitempty"`
	Destination      string `protobuf:"bytes,2,opt,name=destination,proto3" json:"destination,omitempty"`
	ValidatorAddress []byte `protobuf:"bytes,6,opt,name=validator_address,json=validatorAddress,proto3" json:"validator_address,omitempty"`
	ValidatorIndex   uint32 `protobuf:"varint,7,opt,name=validator_index,json=validatorIndex,proto3" json:"validator_index,omitempty"`
	Signature        []byte `protobuf:"bytes,8,opt,name=signature,proto3" json:"signature,omitempty"`
}

func (m *Vote) Reset()         { *m = Vote{} }
func (m *Vote) String() string { return proto.CompactTextString(m) }
func (*Vote) ProtoMessage()    {}
func (*Vote) Descriptor() ([]byte, []int) {
	return fileDescriptor_fdeddabbf1b9645a, []int{0}
}
func (m *Vote) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Vote) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Vote.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *Vote) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Vote.Merge(m, src)
}
func (m *Vote) XXX_Size() int {
	return m.Size()
}
func (m *Vote) XXX_DiscardUnknown() {
	xxx_messageInfo_Vote.DiscardUnknown(m)
}

var xxx_messageInfo_Vote proto.InternalMessageInfo

func (m *Vote) GetHash() string {
	if m != nil {
		return m.Hash
	}
	return ""
}

func (m *Vote) GetDestination() string {
	if m != nil {
		return m.Destination
	}
	return ""
}

func (m *Vote) GetValidatorAddress() []byte {
	if m != nil {
		return m.ValidatorAddress
	}
	return nil
}

func (m *Vote) GetValidatorIndex() uint32 {
	if m != nil {
		return m.ValidatorIndex
	}
	return 0
}

func (m *Vote) GetSignature() []byte {
	if m != nil {
		return m.Signature
	}
	return nil
}

type Proposal struct {
	Source      string       `protobuf:"bytes,1,opt,name=source,proto3" json:"source,omitempty"`
	Destination string       `protobuf:"bytes,2,opt,name=destination,proto3" json:"destination,omitempty"`
	Type        ProposalType `protobuf:"varint,3,opt,name=type,proto3,enum=kardiachain.dualnode.ProposalType" json:"type,omitempty"`
	Args        [][]byte     `protobuf:"bytes,4,rep,name=args,proto3" json:"args,omitempty"`
}

func (m *Proposal) Reset()         { *m = Proposal{} }
func (m *Proposal) String() string { return proto.CompactTextString(m) }
func (*Proposal) ProtoMessage()    {}
func (*Proposal) Descriptor() ([]byte, []int) {
	return fileDescriptor_fdeddabbf1b9645a, []int{1}
}
func (m *Proposal) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Proposal) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Proposal.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *Proposal) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Proposal.Merge(m, src)
}
func (m *Proposal) XXX_Size() int {
	return m.Size()
}
func (m *Proposal) XXX_DiscardUnknown() {
	xxx_messageInfo_Proposal.DiscardUnknown(m)
}

var xxx_messageInfo_Proposal proto.InternalMessageInfo

func (m *Proposal) GetSource() string {
	if m != nil {
		return m.Source
	}
	return ""
}

func (m *Proposal) GetDestination() string {
	if m != nil {
		return m.Destination
	}
	return ""
}

func (m *Proposal) GetType() ProposalType {
	if m != nil {
		return m.Type
	}
	return ProposalUnLock
}

func (m *Proposal) GetArgs() [][]byte {
	if m != nil {
		return m.Args
	}
	return nil
}

type Message struct {
	// Types that are valid to be assigned to Sum:
	//	*Message_Vote
	//	*Message_Proposal
	Sum isMessage_Sum `protobuf_oneof:"sum"`
}

func (m *Message) Reset()         { *m = Message{} }
func (m *Message) String() string { return proto.CompactTextString(m) }
func (*Message) ProtoMessage()    {}
func (*Message) Descriptor() ([]byte, []int) {
	return fileDescriptor_fdeddabbf1b9645a, []int{2}
}
func (m *Message) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Message) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Message.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *Message) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Message.Merge(m, src)
}
func (m *Message) XXX_Size() int {
	return m.Size()
}
func (m *Message) XXX_DiscardUnknown() {
	xxx_messageInfo_Message.DiscardUnknown(m)
}

var xxx_messageInfo_Message proto.InternalMessageInfo

type isMessage_Sum interface {
	isMessage_Sum()
	MarshalTo([]byte) (int, error)
	Size() int
}

type Message_Vote struct {
	Vote *Vote `protobuf:"bytes,1,opt,name=vote,proto3,oneof" json:"vote,omitempty"`
}
type Message_Proposal struct {
	Proposal *Proposal `protobuf:"bytes,2,opt,name=proposal,proto3,oneof" json:"proposal,omitempty"`
}

func (*Message_Vote) isMessage_Sum()     {}
func (*Message_Proposal) isMessage_Sum() {}

func (m *Message) GetSum() isMessage_Sum {
	if m != nil {
		return m.Sum
	}
	return nil
}

func (m *Message) GetVote() *Vote {
	if x, ok := m.GetSum().(*Message_Vote); ok {
		return x.Vote
	}
	return nil
}

func (m *Message) GetProposal() *Proposal {
	if x, ok := m.GetSum().(*Message_Proposal); ok {
		return x.Proposal
	}
	return nil
}

// XXX_OneofWrappers is for the internal use of the proto package.
func (*Message) XXX_OneofWrappers() []interface{} {
	return []interface{}{
		(*Message_Vote)(nil),
		(*Message_Proposal)(nil),
	}
}

func init() {
	proto.RegisterEnum("kardiachain.dualnode.ProposalType", ProposalType_name, ProposalType_value)
	proto.RegisterType((*Vote)(nil), "kardiachain.dualnode.Vote")
	proto.RegisterType((*Proposal)(nil), "kardiachain.dualnode.Proposal")
	proto.RegisterType((*Message)(nil), "kardiachain.dualnode.Message")
}

func init() { proto.RegisterFile("kardiachain/dualnode/types.proto", fileDescriptor_fdeddabbf1b9645a) }

var fileDescriptor_fdeddabbf1b9645a = []byte{
	// 453 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x8c, 0x92, 0x4f, 0x6b, 0x13, 0x41,
	0x18, 0xc6, 0x77, 0x9a, 0x35, 0x4d, 0x27, 0x31, 0xc6, 0x21, 0xc8, 0xb2, 0xc8, 0x30, 0xe4, 0xe2,
	0xa2, 0x90, 0x48, 0x04, 0x4f, 0x5e, 0x52, 0x51, 0x5a, 0xcc, 0x3f, 0xa6, 0xad, 0xd7, 0x30, 0x66,
	0xc7, 0xcd, 0xd2, 0x38, 0xb3, 0xcc, 0x4c, 0x8a, 0xbd, 0x79, 0x94, 0x9e, 0x7a, 0x97, 0x9e, 0xf4,
	0xe0, 0xd1, 0x8f, 0xe1, 0xb1, 0x47, 0x8f, 0x92, 0x7c, 0x11, 0x99, 0x49, 0xb2, 0x4d, 0x21, 0x48,
	0x6f, 0xef, 0x3e, 0xef, 0x6f, 0x9f, 0xf7, 0xe1, 0xd9, 0x85, 0xe4, 0x94, 0xa9, 0x38, 0x65, 0xe3,
	0x09, 0x4b, 0x45, 0x2b, 0x9e, 0xb1, 0xa9, 0x90, 0x31, 0x6f, 0x99, 0xf3, 0x8c, 0xeb, 0x66, 0xa6,
	0xa4, 0x91, 0xa8, 0xbe, 0x41, 0x34, 0xd7, 0x44, 0x58, 0x4f, 0x64, 0x22, 0x1d, 0xd0, 0xb2, 0xd3,
	0x92, 0x6d, 0xfc, 0x02, 0xd0, 0x7f, 0x2f, 0x0d, 0x47, 0x08, 0xfa, 0x13, 0xa6, 0x27, 0x01, 0x20,
	0x20, 0xda, 0xa3, 0x6e, 0x46, 0x04, 0x96, 0x63, 0xae, 0x4d, 0x2a, 0x98, 0x49, 0xa5, 0x08, 0x76,
	0xdc, 0x6a, 0x53, 0x42, 0xcf, 0xe0, 0xc3, 0x33, 0x36, 0x4d, 0x63, 0x66, 0xa4, 0x1a, 0xb1, 0x38,
	0x56, 0x5c, 0xeb, 0xa0, 0x48, 0x40, 0x54, 0xa1, 0xb5, 0x7c, 0xd1, 0x59, 0xea, 0xe8, 0x09, 0x7c,
	0x70, 0x03, 0xa7, 0x22, 0xe6, 0x9f, 0x83, 0x5d, 0x02, 0xa2, 0xfb, 0xb4, 0x9a, 0xcb, 0x87, 0x56,
	0x45, 0x8f, 0xe1, 0x9e, 0x4e, 0x13, 0xc1, 0xcc, 0x4c, 0xf1, 0xa0, 0xe4, 0xdc, 0x6e, 0x84, 0xc6,
	0x25, 0x80, 0xa5, 0xa1, 0x92, 0x99, 0xd4, 0x6c, 0x8a, 0x1e, 0xc1, 0xa2, 0x96, 0x33, 0x35, 0xe6,
	0xab, 0xe0, 0xab, 0xa7, 0x3b, 0x44, 0x7f, 0x09, 0x7d, 0x5b, 0x5a, 0x50, 0x20, 0x20, 0xaa, 0xb6,
	0x1b, 0xcd, 0x6d, 0xa5, 0x35, 0xd7, 0x77, 0x8e, 0xcf, 0x33, 0x4e, 0x1d, 0x6f, 0x8b, 0x62, 0x2a,
	0xd1, 0x81, 0x4f, 0x0a, 0x51, 0x85, 0xba, 0xb9, 0xf1, 0x05, 0xc0, 0xdd, 0x1e, 0xd7, 0x9a, 0x25,
	0x1c, 0x3d, 0x87, 0xfe, 0x99, 0x34, 0xcb, 0x3c, 0xe5, 0x76, 0xb8, 0xdd, 0xd7, 0x56, 0x7e, 0xe0,
	0x51, 0x47, 0xa2, 0x57, 0xb0, 0x94, 0xad, 0xee, 0xb8, 0xa0, 0xe5, 0x36, 0xfe, 0x7f, 0x9a, 0x03,
	0x8f, 0xe6, 0x6f, 0xec, 0xdf, 0x83, 0x05, 0x3d, 0xfb, 0xf4, 0xf4, 0x1b, 0x80, 0x95, 0xcd, 0xb4,
	0xb6, 0xed, 0x21, 0x1d, 0x0c, 0x07, 0x47, 0x9d, 0xee, 0xe8, 0xa4, 0xdf, 0x1d, 0xbc, 0x7e, 0x57,
	0xf3, 0x42, 0x74, 0x71, 0x45, 0xaa, 0x6b, 0xec, 0x44, 0x74, 0xe5, 0xf8, 0xd4, 0x7e, 0xc3, 0x1c,
	0x3c, 0xa6, 0x9d, 0xfe, 0xd1, 0xdb, 0x37, 0xb4, 0x06, 0xc2, 0xfa, 0xc5, 0x15, 0xa9, 0xe5, 0x8e,
	0x8a, 0x09, 0xfd, 0x91, 0xab, 0x5b, 0x70, 0xef, 0xb0, 0x3f, 0x72, 0xbe, 0x3b, 0xb7, 0xe1, 0x5e,
	0x2a, 0x8c, 0x75, 0x0e, 0x4b, 0x5f, 0xbf, 0x63, 0xef, 0xe7, 0x0f, 0x0c, 0xf6, 0x83, 0xdf, 0x73,
	0x0c, 0xae, 0xe7, 0x18, 0xfc, 0x9d, 0x63, 0x70, 0xb9, 0xc0, 0xde, 0xf5, 0x02, 0x7b, 0x7f, 0x16,
	0xd8, 0xfb, 0x50, 0x74, 0xff, 0xe1, 0x8b, 0x7f, 0x01, 0x00, 0x00, 0xff, 0xff, 0x3c, 0x7d, 0x0e,
	0xd0, 0xd7, 0x02, 0x00, 0x00,
}

func (m *Vote) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Vote) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Vote) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.Signature) > 0 {
		i -= len(m.Signature)
		copy(dAtA[i:], m.Signature)
		i = encodeVarintTypes(dAtA, i, uint64(len(m.Signature)))
		i--
		dAtA[i] = 0x42
	}
	if m.ValidatorIndex != 0 {
		i = encodeVarintTypes(dAtA, i, uint64(m.ValidatorIndex))
		i--
		dAtA[i] = 0x38
	}
	if len(m.ValidatorAddress) > 0 {
		i -= len(m.ValidatorAddress)
		copy(dAtA[i:], m.ValidatorAddress)
		i = encodeVarintTypes(dAtA, i, uint64(len(m.ValidatorAddress)))
		i--
		dAtA[i] = 0x32
	}
	if len(m.Destination) > 0 {
		i -= len(m.Destination)
		copy(dAtA[i:], m.Destination)
		i = encodeVarintTypes(dAtA, i, uint64(len(m.Destination)))
		i--
		dAtA[i] = 0x12
	}
	if len(m.Hash) > 0 {
		i -= len(m.Hash)
		copy(dAtA[i:], m.Hash)
		i = encodeVarintTypes(dAtA, i, uint64(len(m.Hash)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *Proposal) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Proposal) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Proposal) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.Args) > 0 {
		for iNdEx := len(m.Args) - 1; iNdEx >= 0; iNdEx-- {
			i -= len(m.Args[iNdEx])
			copy(dAtA[i:], m.Args[iNdEx])
			i = encodeVarintTypes(dAtA, i, uint64(len(m.Args[iNdEx])))
			i--
			dAtA[i] = 0x22
		}
	}
	if m.Type != 0 {
		i = encodeVarintTypes(dAtA, i, uint64(m.Type))
		i--
		dAtA[i] = 0x18
	}
	if len(m.Destination) > 0 {
		i -= len(m.Destination)
		copy(dAtA[i:], m.Destination)
		i = encodeVarintTypes(dAtA, i, uint64(len(m.Destination)))
		i--
		dAtA[i] = 0x12
	}
	if len(m.Source) > 0 {
		i -= len(m.Source)
		copy(dAtA[i:], m.Source)
		i = encodeVarintTypes(dAtA, i, uint64(len(m.Source)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *Message) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Message) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Message) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.Sum != nil {
		{
			size := m.Sum.Size()
			i -= size
			if _, err := m.Sum.MarshalTo(dAtA[i:]); err != nil {
				return 0, err
			}
		}
	}
	return len(dAtA) - i, nil
}

func (m *Message_Vote) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Message_Vote) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	if m.Vote != nil {
		{
			size, err := m.Vote.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintTypes(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}
func (m *Message_Proposal) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Message_Proposal) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	if m.Proposal != nil {
		{
			size, err := m.Proposal.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintTypes(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x12
	}
	return len(dAtA) - i, nil
}
func encodeVarintTypes(dAtA []byte, offset int, v uint64) int {
	offset -= sovTypes(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *Vote) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Hash)
	if l > 0 {
		n += 1 + l + sovTypes(uint64(l))
	}
	l = len(m.Destination)
	if l > 0 {
		n += 1 + l + sovTypes(uint64(l))
	}
	l = len(m.ValidatorAddress)
	if l > 0 {
		n += 1 + l + sovTypes(uint64(l))
	}
	if m.ValidatorIndex != 0 {
		n += 1 + sovTypes(uint64(m.ValidatorIndex))
	}
	l = len(m.Signature)
	if l > 0 {
		n += 1 + l + sovTypes(uint64(l))
	}
	return n
}

func (m *Proposal) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Source)
	if l > 0 {
		n += 1 + l + sovTypes(uint64(l))
	}
	l = len(m.Destination)
	if l > 0 {
		n += 1 + l + sovTypes(uint64(l))
	}
	if m.Type != 0 {
		n += 1 + sovTypes(uint64(m.Type))
	}
	if len(m.Args) > 0 {
		for _, b := range m.Args {
			l = len(b)
			n += 1 + l + sovTypes(uint64(l))
		}
	}
	return n
}

func (m *Message) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.Sum != nil {
		n += m.Sum.Size()
	}
	return n
}

func (m *Message_Vote) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.Vote != nil {
		l = m.Vote.Size()
		n += 1 + l + sovTypes(uint64(l))
	}
	return n
}
func (m *Message_Proposal) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.Proposal != nil {
		l = m.Proposal.Size()
		n += 1 + l + sovTypes(uint64(l))
	}
	return n
}

func sovTypes(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozTypes(x uint64) (n int) {
	return sovTypes(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *Vote) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowTypes
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: Vote: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Vote: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Hash", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTypes
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthTypes
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthTypes
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Hash = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Destination", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTypes
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthTypes
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthTypes
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Destination = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 6:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ValidatorAddress", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTypes
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				byteLen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if byteLen < 0 {
				return ErrInvalidLengthTypes
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthTypes
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.ValidatorAddress = append(m.ValidatorAddress[:0], dAtA[iNdEx:postIndex]...)
			if m.ValidatorAddress == nil {
				m.ValidatorAddress = []byte{}
			}
			iNdEx = postIndex
		case 7:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field ValidatorIndex", wireType)
			}
			m.ValidatorIndex = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTypes
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.ValidatorIndex |= uint32(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 8:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Signature", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTypes
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				byteLen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if byteLen < 0 {
				return ErrInvalidLengthTypes
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthTypes
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Signature = append(m.Signature[:0], dAtA[iNdEx:postIndex]...)
			if m.Signature == nil {
				m.Signature = []byte{}
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipTypes(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthTypes
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *Proposal) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowTypes
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: Proposal: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Proposal: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Source", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTypes
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthTypes
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthTypes
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Source = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Destination", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTypes
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthTypes
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthTypes
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Destination = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 3:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Type", wireType)
			}
			m.Type = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTypes
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Type |= ProposalType(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Args", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTypes
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				byteLen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if byteLen < 0 {
				return ErrInvalidLengthTypes
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthTypes
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Args = append(m.Args, make([]byte, postIndex-iNdEx))
			copy(m.Args[len(m.Args)-1], dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipTypes(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthTypes
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *Message) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowTypes
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: Message: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Message: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Vote", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTypes
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthTypes
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthTypes
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			v := &Vote{}
			if err := v.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			m.Sum = &Message_Vote{v}
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Proposal", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTypes
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthTypes
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthTypes
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			v := &Proposal{}
			if err := v.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			m.Sum = &Message_Proposal{v}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipTypes(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthTypes
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func skipTypes(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowTypes
			}
			if iNdEx >= l {
				return 0, io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		wireType := int(wire & 0x7)
		switch wireType {
		case 0:
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowTypes
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				iNdEx++
				if dAtA[iNdEx-1] < 0x80 {
					break
				}
			}
		case 1:
			iNdEx += 8
		case 2:
			var length int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowTypes
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				length |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if length < 0 {
				return 0, ErrInvalidLengthTypes
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupTypes
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthTypes
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthTypes        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowTypes          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupTypes = fmt.Errorf("proto: unexpected end of group")
)
