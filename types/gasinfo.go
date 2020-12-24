package types

import (
	"fmt"
	"io"
	"strings"

	math_bits "math/bits"

	"github.com/gogo/protobuf/proto"
)

// GasInfo defines tx execution gas context.
type GasInfo struct {
	// GasWanted is the maximum units of work we allow this tx to perform.
	GasWanted uint64 `protobuf:"varint,1,opt,name=gas_wanted,json=gasWanted,proto3" json:"gas_wanted,omitempty" yaml:"gas_wanted"`
	// GasUsed is the amount of gas actually consumed.
	GasUsed uint64 `protobuf:"varint,2,opt,name=gas_used,json=gasUsed,proto3" json:"gas_used,omitempty" yaml:"gas_used"`
}

func NewGasInfo() GasInfo {
	return GasInfo{}
}

var xxx_messageInfo_GasInfo proto.InternalMessageInfo

func (m *GasInfo) Reset()      { *m = GasInfo{} }
func (*GasInfo) ProtoMessage() {}
func (*GasInfo) Descriptor() ([]byte, []int) {
	return fileDescriptor_2c0f90c600ad7e2e, []int{5}
}
func (m *GasInfo) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *GasInfo) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_GasInfo.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *GasInfo) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GasInfo.Merge(m, src)
}
func (m *GasInfo) XXX_Size() int {
	return m.Size()
}
func (m *GasInfo) XXX_DiscardUnknown() {
	xxx_messageInfo_GasInfo.DiscardUnknown(m)
}

func (m *GasInfo) GetGasWanted() uint64 {
	if m != nil {
		return m.GasWanted
	}
	return 0
}

func (m *GasInfo) GetGasUsed() uint64 {
	if m != nil {
		return m.GasUsed
	}
	return 0
}

func (m *GasInfo) String() string {
	if m != nil {
		return "nil"
	}

	s := strings.Join([]string{`&GasInfo{`,
		`GasWanted:` + fmt.Sprintf("%d", m.GasWanted) + `,`,
		`GasUsed:` + fmt.Sprintf("%d", m.GasUsed) + `,`,
		`}`,
	}, "")
	return s
}

var fileDescriptor_2c0f90c600ad7e2e = []byte{
	// 534 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x9c, 0x52, 0x4f, 0x6f, 0xd3, 0x4e,
	0x10, 0xb5, 0x7f, 0xf6, 0x2f, 0x7f, 0x36, 0xe1, 0x4f, 0x17, 0x8a, 0xa2, 0x0a, 0xec, 0xc8, 0x48,
	0x28, 0x48, 0xd4, 0x16, 0x29, 0xa7, 0x70, 0xc2, 0x04, 0x55, 0xe1, 0x84, 0x16, 0x01, 0x12, 0x97,
	0x68, 0xe3, 0xdd, 0xba, 0x56, 0xe3, 0xdd, 0xc8, 0xbb, 0x29, 0xca, 0x2d, 0x47, 0x8e, 0x7c, 0x84,
	0x7e, 0x9c, 0x1e, 0x73, 0xac, 0x10, 0xb2, 0x20, 0xb9, 0x70, 0xee, 0x91, 0x13, 0xda, 0xb5, 0xc1,
	0x6a, 0x7a, 0xe3, 0x92, 0xcc, 0xce, 0xbc, 0x79, 0x33, 0xf3, 0xfc, 0xc0, 0x8e, 0x5c, 0xcc, 0xa8,
	0x08, 0xf4, 0xaf, 0x3f, 0xcb, 0xb8, 0xe4, 0xf0, 0x46, 0xc4, 0x45, 0xca, 0xc5, 0x58, 0x90, 0x13,
	0xff, 0xf4, 0xe9, 0xde, 0x23, 0x79, 0x9c, 0x64, 0x64, 0x3c, 0xc3, 0x99, 0x5c, 0x04, 0x1a, 0x11,
	0xc4, 0x3c, 0xe6, 0x55, 0x54, 0xb4, 0xed, 0x1d, 0x5c, 0xc7, 0x49, 0xca, 0x08, 0xcd, 0xd2, 0x84,
	0xc9, 0x00, 0x4f, 0xa2, 0x24, 0xb8, 0x36, 0xcb, 0x3b, 0x04, 0xf6, 0x4b, 0x9e, 0x30, 0x78, 0x17,
	0xfc, 0x4f, 0x28, 0xe3, 0x69, 0xc7, 0xec, 0x9a, 0xbd, 0x26, 0x2a, 0x1e, 0xf0, 0x21, 0xa8, 0xe1,
	0x94, 0xcf, 0x99, 0xec, 0xfc, 0xa7, 0xd2, 0x61, 0xeb, 0x3c, 0x77, 0x8d, 0xaf, 0xb9, 0x6b, 0x8d,
	0x98, 0x44, 0x65, 0x69, 0x60, 0xff, 0x3c, 0x73, 0x4d, 0xef, 0x35, 0xa8, 0x0f, 0x69, 0xf4, 0x2f,
	0x5c, 0x43, 0x1a, 0x6d, 0x71, 0x3d, 0x06, 0x8d, 0x11, 0x93, 0x6f, 0xb4, 0x18, 0x0f, 0x80, 0x95,
	0x30, 0x59, 0x50, 0x5d, 0x9d, 0xaf, 0xf2, 0x0a, 0x3a, 0xa4, 0xd1, 0x5f, 0x28, 0xa1, 0xd1, 0x36,
	0x54, 0xd1, 0xab, 0xbc, 0x17, 0x82, 0xf6, 0x7b, 0x3c, 0x7d, 0x41, 0x48, 0x46, 0x85, 0xa0, 0x02,
	0x3e, 0x01, 0x4d, 0xfc, 0xe7, 0xd1, 0x31, 0xbb, 0x56, 0xaf, 0x1d, 0xde, 0xfc, 0x95, 0xbb, 0xa0,
	0x02, 0xa1, 0x0a, 0x30, 0xb0, 0x97, 0xdf, 0xba, 0xa6, 0xc7, 0x41, 0xfd, 0x10, 0x8b, 0x11, 0x3b,
	0xe2, 0xf0, 0x19, 0x00, 0x31, 0x16, 0xe3, 0x4f, 0x98, 0x49, 0x4a, 0xf4, 0x50, 0x3b, 0xdc, 0xbd,
	0xcc, 0xdd, 0x9d, 0x05, 0x4e, 0xa7, 0x03, 0xaf, 0xaa, 0x79, 0xa8, 0x19, 0x63, 0xf1, 0x41, 0xc7,
	0xd0, 0x07, 0x0d, 0x55, 0x99, 0x0b, 0x4a, 0xb4, 0x0e, 0x76, 0x78, 0xe7, 0x32, 0x77, 0x6f, 0x55,
	0x3d, 0xaa, 0xe2, 0xa1, 0x7a, 0x8c, 0xc5, 0x3b, 0x15, 0xcd, 0x40, 0x0d, 0x51, 0x31, 0x9f, 0x4a,
	0x08, 0x81, 0x4d, 0xb0, 0xc4, 0x7a, 0x52, 0x1b, 0xe9, 0x18, 0xde, 0x06, 0xd6, 0x94, 0xc7, 0x85,
	0xa0, 0x48, 0x85, 0x70, 0x00, 0x6a, 0xf4, 0x94, 0x32, 0x29, 0x3a, 0x56, 0xd7, 0xea, 0xb5, 0xfa,
	0xf7, 0xfd, 0xca, 0x03, 0xbe, 0xf2, 0x80, 0x5f, 0x7c, 0xfd, 0x57, 0x0a, 0x14, 0xda, 0x4a, 0x24,
	0x54, 0x76, 0x0c, 0xec, 0xcf, 0x67, 0xae, 0xe1, 0x2d, 0x4d, 0x00, 0xdf, 0x26, 0xe9, 0x7c, 0x8a,
	0x65, 0xc2, 0x19, 0xa2, 0x62, 0xc6, 0x99, 0xa0, 0xf0, 0x79, 0xb1, 0x78, 0xc2, 0x8e, 0xb8, 0x5e,
	0xa1, 0xd5, 0xbf, 0xe7, 0x5f, 0xf1, 0xa9, 0x5f, 0x0a, 0x13, 0x36, 0x14, 0xe9, 0x2a, 0x77, 0x4d,
	0x7d, 0x85, 0xd6, 0x6a, 0x1f, 0xd4, 0x32, 0x7d, 0x85, 0x5e, 0xb5, 0xd5, 0xdf, 0xdd, 0x6a, 0x2d,
	0x4e, 0x44, 0x25, 0x28, 0x1c, 0x5e, 0xfc, 0x70, 0x8c, 0xe5, 0xda, 0x31, 0xce, 0xd7, 0x8e, 0xb9,
	0x5a, 0x3b, 0xe6, 0xf7, 0xb5, 0x63, 0x7e, 0xd9, 0x38, 0xc6, 0x6a, 0xe3, 0x18, 0x17, 0x1b, 0xc7,
	0xf8, 0xe8, 0xc5, 0x89, 0x3c, 0x9e, 0x4f, 0xfc, 0x88, 0xa7, 0x41, 0x41, 0x55, 0xfe, 0xed, 0x0b,
	0x72, 0x52, 0x18, 0x7c, 0x52, 0xd3, 0x0e, 0x3f, 0xf8, 0x1d, 0x00, 0x00, 0xff, 0xff, 0x0e, 0xe8,
	0xc3, 0x0c, 0x62, 0x03, 0x00, 0x00,
}

var (
	ErrInvalidLengthTypes        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowTypes          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupTypes = fmt.Errorf("proto: unexpected end of group")
)

func (m *GasInfo) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *GasInfo) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *GasInfo) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.GasUsed != 0 {
		i = encodeVarintTypes(dAtA, i, uint64(m.GasUsed))
		i--
		dAtA[i] = 0x10
	}
	if m.GasWanted != 0 {
		i = encodeVarintTypes(dAtA, i, uint64(m.GasWanted))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func (m *GasInfo) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.GasWanted != 0 {
		n += 1 + sovTypes(uint64(m.GasWanted))
	}
	if m.GasUsed != 0 {
		n += 1 + sovTypes(uint64(m.GasUsed))
	}
	return n
}

func (m *GasInfo) Unmarshal(dAtA []byte) error {
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
			return fmt.Errorf("proto: GasInfo: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: GasInfo: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field GasWanted", wireType)
			}
			m.GasWanted = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTypes
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.GasWanted |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field GasUsed", wireType)
			}
			m.GasUsed = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTypes
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.GasUsed |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		default:
			iNdEx = preIndex
			skippy, err := skipTypes(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthTypes
			}
			if (iNdEx + skippy) < 0 {
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

func sovTypes(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozTypes(x uint64) (n int) {
	return sovTypes(uint64((x << 1) ^ uint64((int64(x) >> 63))))
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
