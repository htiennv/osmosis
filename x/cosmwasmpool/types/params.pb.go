// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: osmosis/cosmwasmpool/v1beta1/params.proto

package types

import (
	fmt "fmt"
	_ "github.com/cosmos/cosmos-proto"
	_ "github.com/gogo/protobuf/gogoproto"
	proto "github.com/gogo/protobuf/proto"
	_ "github.com/gogo/protobuf/types"
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

type Params struct {
	// code_ide_whitelist contains the list of code ids that are allowed to be
	// instantiated.
	CodeIdWhitelist []uint64 `protobuf:"varint,1,rep,packed,name=code_id_whitelist,json=codeIdWhitelist,proto3" json:"code_id_whitelist,omitempty" yaml:"code_id_whitelist"`
	// pool_migration_limit is the maximum number of pools that can be migrated
	// at once via governance proposal. This is to have a constant bound on the
	// number of pools that can be migrated at once and remove the possibility
	// of an unlikely scenario of causing a chain halt due to a large migration.
	PoolMigrationLimit uint64 `protobuf:"varint,2,opt,name=pool_migration_limit,json=poolMigrationLimit,proto3" json:"pool_migration_limit,omitempty" yaml:"pool_migration_limit"`
}

func (m *Params) Reset()         { *m = Params{} }
func (m *Params) String() string { return proto.CompactTextString(m) }
func (*Params) ProtoMessage()    {}
func (*Params) Descriptor() ([]byte, []int) {
	return fileDescriptor_6cf69242a2b5e68e, []int{0}
}
func (m *Params) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Params) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Params.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *Params) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Params.Merge(m, src)
}
func (m *Params) XXX_Size() int {
	return m.Size()
}
func (m *Params) XXX_DiscardUnknown() {
	xxx_messageInfo_Params.DiscardUnknown(m)
}

var xxx_messageInfo_Params proto.InternalMessageInfo

func (m *Params) GetCodeIdWhitelist() []uint64 {
	if m != nil {
		return m.CodeIdWhitelist
	}
	return nil
}

func (m *Params) GetPoolMigrationLimit() uint64 {
	if m != nil {
		return m.PoolMigrationLimit
	}
	return 0
}

func init() {
	proto.RegisterType((*Params)(nil), "osmosis.cosmwasmpool.v1beta1.Params")
}

func init() {
	proto.RegisterFile("osmosis/cosmwasmpool/v1beta1/params.proto", fileDescriptor_6cf69242a2b5e68e)
}

var fileDescriptor_6cf69242a2b5e68e = []byte{
	// 300 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xd2, 0xcc, 0x2f, 0xce, 0xcd,
	0x2f, 0xce, 0x2c, 0xd6, 0x4f, 0xce, 0x2f, 0xce, 0x2d, 0x4f, 0x2c, 0xce, 0x2d, 0xc8, 0xcf, 0xcf,
	0xd1, 0x2f, 0x33, 0x4c, 0x4a, 0x2d, 0x49, 0x34, 0xd4, 0x2f, 0x48, 0x2c, 0x4a, 0xcc, 0x2d, 0xd6,
	0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0x92, 0x81, 0x2a, 0xd5, 0x43, 0x56, 0xaa, 0x07, 0x55, 0x2a,
	0x25, 0x92, 0x9e, 0x9f, 0x9e, 0x0f, 0x56, 0xa8, 0x0f, 0x62, 0x41, 0xf4, 0x48, 0x49, 0x26, 0x83,
	0x35, 0xc5, 0x43, 0x24, 0x20, 0x1c, 0xa8, 0x94, 0x5c, 0x7a, 0x7e, 0x7e, 0x7a, 0x4e, 0xaa, 0x3e,
	0x98, 0x97, 0x54, 0x9a, 0xa6, 0x9f, 0x52, 0x5a, 0x94, 0x58, 0x92, 0x99, 0x9f, 0x07, 0x91, 0x57,
	0x5a, 0xca, 0xc8, 0xc5, 0x16, 0x00, 0xb6, 0x5f, 0xc8, 0x83, 0x4b, 0x30, 0x39, 0x3f, 0x25, 0x35,
	0x3e, 0x33, 0x25, 0xbe, 0x3c, 0x23, 0xb3, 0x24, 0x35, 0x27, 0xb3, 0xb8, 0x44, 0x82, 0x51, 0x81,
	0x59, 0x83, 0xc5, 0x49, 0xe6, 0xd3, 0x3d, 0x79, 0x89, 0xca, 0xc4, 0xdc, 0x1c, 0x2b, 0x25, 0x0c,
	0x25, 0x4a, 0x41, 0xfc, 0x20, 0x31, 0xcf, 0x94, 0x70, 0x98, 0x88, 0x50, 0x20, 0x97, 0x08, 0xc8,
	0xd5, 0xf1, 0xb9, 0x99, 0xe9, 0x10, 0xcb, 0xe2, 0x73, 0x32, 0x73, 0x33, 0x4b, 0x24, 0x98, 0x14,
	0x18, 0x35, 0x58, 0x9c, 0xe4, 0x3f, 0xdd, 0x93, 0x97, 0x86, 0x18, 0x86, 0x4d, 0x95, 0x52, 0x90,
	0x10, 0x48, 0xd8, 0x17, 0x26, 0xea, 0x03, 0x12, 0x74, 0x0a, 0x3e, 0xf1, 0x48, 0x8e, 0xf1, 0xc2,
	0x23, 0x39, 0xc6, 0x07, 0x8f, 0xe4, 0x18, 0x27, 0x3c, 0x96, 0x63, 0xb8, 0xf0, 0x58, 0x8e, 0xe1,
	0xc6, 0x63, 0x39, 0x86, 0x28, 0xcb, 0xf4, 0xcc, 0x92, 0x8c, 0xd2, 0x24, 0xbd, 0xe4, 0xfc, 0x5c,
	0x7d, 0x68, 0xd8, 0xe9, 0xe6, 0x24, 0x26, 0x15, 0xc3, 0x38, 0xfa, 0x65, 0x65, 0x46, 0x86, 0xfa,
	0x15, 0xa8, 0x41, 0x5f, 0x52, 0x59, 0x90, 0x5a, 0x9c, 0xc4, 0x06, 0x0e, 0x03, 0x63, 0x40, 0x00,
	0x00, 0x00, 0xff, 0xff, 0x25, 0x57, 0x3e, 0xdb, 0x9f, 0x01, 0x00, 0x00,
}

func (m *Params) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Params) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Params) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.PoolMigrationLimit != 0 {
		i = encodeVarintParams(dAtA, i, uint64(m.PoolMigrationLimit))
		i--
		dAtA[i] = 0x10
	}
	if len(m.CodeIdWhitelist) > 0 {
		dAtA2 := make([]byte, len(m.CodeIdWhitelist)*10)
		var j1 int
		for _, num := range m.CodeIdWhitelist {
			for num >= 1<<7 {
				dAtA2[j1] = uint8(uint64(num)&0x7f | 0x80)
				num >>= 7
				j1++
			}
			dAtA2[j1] = uint8(num)
			j1++
		}
		i -= j1
		copy(dAtA[i:], dAtA2[:j1])
		i = encodeVarintParams(dAtA, i, uint64(j1))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func encodeVarintParams(dAtA []byte, offset int, v uint64) int {
	offset -= sovParams(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *Params) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if len(m.CodeIdWhitelist) > 0 {
		l = 0
		for _, e := range m.CodeIdWhitelist {
			l += sovParams(uint64(e))
		}
		n += 1 + sovParams(uint64(l)) + l
	}
	if m.PoolMigrationLimit != 0 {
		n += 1 + sovParams(uint64(m.PoolMigrationLimit))
	}
	return n
}

func sovParams(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozParams(x uint64) (n int) {
	return sovParams(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *Params) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowParams
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
			return fmt.Errorf("proto: Params: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Params: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType == 0 {
				var v uint64
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return ErrIntOverflowParams
					}
					if iNdEx >= l {
						return io.ErrUnexpectedEOF
					}
					b := dAtA[iNdEx]
					iNdEx++
					v |= uint64(b&0x7F) << shift
					if b < 0x80 {
						break
					}
				}
				m.CodeIdWhitelist = append(m.CodeIdWhitelist, v)
			} else if wireType == 2 {
				var packedLen int
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return ErrIntOverflowParams
					}
					if iNdEx >= l {
						return io.ErrUnexpectedEOF
					}
					b := dAtA[iNdEx]
					iNdEx++
					packedLen |= int(b&0x7F) << shift
					if b < 0x80 {
						break
					}
				}
				if packedLen < 0 {
					return ErrInvalidLengthParams
				}
				postIndex := iNdEx + packedLen
				if postIndex < 0 {
					return ErrInvalidLengthParams
				}
				if postIndex > l {
					return io.ErrUnexpectedEOF
				}
				var elementCount int
				var count int
				for _, integer := range dAtA[iNdEx:postIndex] {
					if integer < 128 {
						count++
					}
				}
				elementCount = count
				if elementCount != 0 && len(m.CodeIdWhitelist) == 0 {
					m.CodeIdWhitelist = make([]uint64, 0, elementCount)
				}
				for iNdEx < postIndex {
					var v uint64
					for shift := uint(0); ; shift += 7 {
						if shift >= 64 {
							return ErrIntOverflowParams
						}
						if iNdEx >= l {
							return io.ErrUnexpectedEOF
						}
						b := dAtA[iNdEx]
						iNdEx++
						v |= uint64(b&0x7F) << shift
						if b < 0x80 {
							break
						}
					}
					m.CodeIdWhitelist = append(m.CodeIdWhitelist, v)
				}
			} else {
				return fmt.Errorf("proto: wrong wireType = %d for field CodeIdWhitelist", wireType)
			}
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field PoolMigrationLimit", wireType)
			}
			m.PoolMigrationLimit = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowParams
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.PoolMigrationLimit |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		default:
			iNdEx = preIndex
			skippy, err := skipParams(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthParams
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
func skipParams(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowParams
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
					return 0, ErrIntOverflowParams
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
					return 0, ErrIntOverflowParams
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
				return 0, ErrInvalidLengthParams
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupParams
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthParams
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthParams        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowParams          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupParams = fmt.Errorf("proto: unexpected end of group")
)
