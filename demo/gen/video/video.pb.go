// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: demo/api/video.proto

package model

import (
	fmt "fmt"
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

type Video struct {
	ID       string   `protobuf:"bytes,1,opt,name=ID,proto3" json:"ID,omitempty"`
	Title    string   `protobuf:"bytes,2,opt,name=title,proto3" json:"title,omitempty"`
	Author   string   `protobuf:"bytes,3,opt,name=Author,proto3" json:"Author,omitempty"`
	PostTime int64    `protobuf:"varint,4,opt,name=PostTime,proto3" json:"PostTime,omitempty"`
	View     int32    `protobuf:"varint,5,opt,name=View,proto3" json:"View,omitempty"`
	Like     int32    `protobuf:"varint,6,opt,name=Like,proto3" json:"Like,omitempty"`
	Coin     int32    `protobuf:"varint,7,opt,name=Coin,proto3" json:"Coin,omitempty"`
	Favorite int32    `protobuf:"varint,8,opt,name=Favorite,proto3" json:"Favorite,omitempty"`
	Share    int32    `protobuf:"varint,9,opt,name=Share,proto3" json:"Share,omitempty"`
	KeyWords []string `protobuf:"bytes,10,rep,name=KeyWords,proto3" json:"KeyWords,omitempty"`
}

func (m *Video) Reset()         { *m = Video{} }
func (m *Video) String() string { return proto.CompactTextString(m) }
func (*Video) ProtoMessage()    {}
func (*Video) Descriptor() ([]byte, []int) {
	return fileDescriptor_b8ab23b2436d5792, []int{0}
}
func (m *Video) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Video) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Video.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *Video) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Video.Merge(m, src)
}
func (m *Video) XXX_Size() int {
	return m.Size()
}
func (m *Video) XXX_DiscardUnknown() {
	xxx_messageInfo_Video.DiscardUnknown(m)
}

var xxx_messageInfo_Video proto.InternalMessageInfo

func (m *Video) GetID() string {
	if m != nil {
		return m.ID
	}
	return ""
}

func (m *Video) GetTitle() string {
	if m != nil {
		return m.Title
	}
	return ""
}

func (m *Video) GetAuthor() string {
	if m != nil {
		return m.Author
	}
	return ""
}

func (m *Video) GetPostTime() int64 {
	if m != nil {
		return m.PostTime
	}
	return 0
}

func (m *Video) GetView() int32 {
	if m != nil {
		return m.View
	}
	return 0
}

func (m *Video) GetLike() int32 {
	if m != nil {
		return m.Like
	}
	return 0
}

func (m *Video) GetCoin() int32 {
	if m != nil {
		return m.Coin
	}
	return 0
}

func (m *Video) GetFavorite() int32 {
	if m != nil {
		return m.Favorite
	}
	return 0
}

func (m *Video) GetShare() int32 {
	if m != nil {
		return m.Share
	}
	return 0
}

func (m *Video) GetKeyWords() []string {
	if m != nil {
		return m.KeyWords
	}
	return nil
}

func init() {
	proto.RegisterType((*Video)(nil), "model.Video")
}

func init() { proto.RegisterFile("demo/api/video.proto", fileDescriptor_b8ab23b2436d5792) }

var fileDescriptor_b8ab23b2436d5792 = []byte{
	// 261 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x3c, 0x90, 0x41, 0x4b, 0xfb, 0x30,
	0x18, 0x87, 0x9b, 0x76, 0xe9, 0x7f, 0xcd, 0xe1, 0x7f, 0x08, 0x63, 0xbc, 0x78, 0x08, 0xc5, 0x53,
	0x4f, 0xab, 0xe0, 0xd1, 0x93, 0x3a, 0x84, 0xa1, 0x07, 0xa9, 0x32, 0xc1, 0x5b, 0xa5, 0x2f, 0x2e,
	0xb8, 0xee, 0x1d, 0x59, 0x9c, 0xf8, 0x2d, 0xfc, 0x58, 0x1e, 0x77, 0xf4, 0x28, 0xed, 0xd1, 0x2f,
	0x21, 0x49, 0xa4, 0xb7, 0xdf, 0xf3, 0x84, 0x84, 0xf0, 0x88, 0x49, 0x83, 0x2d, 0x95, 0xf5, 0x56,
	0x97, 0x7b, 0xdd, 0x20, 0xcd, 0xb6, 0x86, 0x2c, 0x49, 0xde, 0x52, 0x83, 0xeb, 0xe3, 0x1f, 0x26,
	0xf8, 0xd2, 0x69, 0xf9, 0x5f, 0xc4, 0x8b, 0x39, 0xb0, 0x9c, 0x15, 0x59, 0x15, 0x2f, 0xe6, 0x72,
	0x22, 0xb8, 0xd5, 0x76, 0x8d, 0x10, 0x7b, 0x15, 0x40, 0x4e, 0x45, 0x7a, 0xfe, 0x6a, 0x57, 0x64,
	0x20, 0xf1, 0xfa, 0x8f, 0xe4, 0x91, 0x18, 0xdf, 0xd2, 0xce, 0xde, 0xeb, 0x16, 0x61, 0x94, 0xb3,
	0x22, 0xa9, 0x06, 0x96, 0x52, 0x8c, 0x96, 0x1a, 0xdf, 0x80, 0xe7, 0xac, 0xe0, 0x95, 0xdf, 0xce,
	0xdd, 0xe8, 0x17, 0x84, 0x34, 0x38, 0xb7, 0x9d, 0xbb, 0x24, 0xbd, 0x81, 0x7f, 0xc1, 0xb9, 0xed,
	0xde, 0xbd, 0xaa, 0xf7, 0x64, 0xb4, 0x45, 0x18, 0x7b, 0x3f, 0xb0, 0xfb, 0xe1, 0xdd, 0xaa, 0x36,
	0x08, 0x99, 0x3f, 0x08, 0xe0, 0x6e, 0x5c, 0xe3, 0xfb, 0x03, 0x99, 0x66, 0x07, 0x22, 0x4f, 0x8a,
	0xac, 0x1a, 0xf8, 0xe2, 0xe4, 0xb3, 0x53, 0xec, 0xd0, 0x29, 0xf6, 0xdd, 0x29, 0xf6, 0xd1, 0xab,
	0xe8, 0xd0, 0xab, 0xe8, 0xab, 0x57, 0xd1, 0xe3, 0x74, 0x56, 0xfa, 0x4c, 0xcf, 0xb8, 0x09, 0x99,
	0xce, 0x7c, 0x9f, 0xa7, 0xd4, 0xd7, 0x3a, 0xfd, 0x0d, 0x00, 0x00, 0xff, 0xff, 0xdf, 0xfb, 0x9c,
	0xe8, 0x45, 0x01, 0x00, 0x00,
}

func (m *Video) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Video) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Video) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.KeyWords) > 0 {
		for iNdEx := len(m.KeyWords) - 1; iNdEx >= 0; iNdEx-- {
			i -= len(m.KeyWords[iNdEx])
			copy(dAtA[i:], m.KeyWords[iNdEx])
			i = encodeVarintVideo(dAtA, i, uint64(len(m.KeyWords[iNdEx])))
			i--
			dAtA[i] = 0x52
		}
	}
	if m.Share != 0 {
		i = encodeVarintVideo(dAtA, i, uint64(m.Share))
		i--
		dAtA[i] = 0x48
	}
	if m.Favorite != 0 {
		i = encodeVarintVideo(dAtA, i, uint64(m.Favorite))
		i--
		dAtA[i] = 0x40
	}
	if m.Coin != 0 {
		i = encodeVarintVideo(dAtA, i, uint64(m.Coin))
		i--
		dAtA[i] = 0x38
	}
	if m.Like != 0 {
		i = encodeVarintVideo(dAtA, i, uint64(m.Like))
		i--
		dAtA[i] = 0x30
	}
	if m.View != 0 {
		i = encodeVarintVideo(dAtA, i, uint64(m.View))
		i--
		dAtA[i] = 0x28
	}
	if m.PostTime != 0 {
		i = encodeVarintVideo(dAtA, i, uint64(m.PostTime))
		i--
		dAtA[i] = 0x20
	}
	if len(m.Author) > 0 {
		i -= len(m.Author)
		copy(dAtA[i:], m.Author)
		i = encodeVarintVideo(dAtA, i, uint64(len(m.Author)))
		i--
		dAtA[i] = 0x1a
	}
	if len(m.Title) > 0 {
		i -= len(m.Title)
		copy(dAtA[i:], m.Title)
		i = encodeVarintVideo(dAtA, i, uint64(len(m.Title)))
		i--
		dAtA[i] = 0x12
	}
	if len(m.ID) > 0 {
		i -= len(m.ID)
		copy(dAtA[i:], m.ID)
		i = encodeVarintVideo(dAtA, i, uint64(len(m.ID)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func encodeVarintVideo(dAtA []byte, offset int, v uint64) int {
	offset -= sovVideo(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *Video) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.ID)
	if l > 0 {
		n += 1 + l + sovVideo(uint64(l))
	}
	l = len(m.Title)
	if l > 0 {
		n += 1 + l + sovVideo(uint64(l))
	}
	l = len(m.Author)
	if l > 0 {
		n += 1 + l + sovVideo(uint64(l))
	}
	if m.PostTime != 0 {
		n += 1 + sovVideo(uint64(m.PostTime))
	}
	if m.View != 0 {
		n += 1 + sovVideo(uint64(m.View))
	}
	if m.Like != 0 {
		n += 1 + sovVideo(uint64(m.Like))
	}
	if m.Coin != 0 {
		n += 1 + sovVideo(uint64(m.Coin))
	}
	if m.Favorite != 0 {
		n += 1 + sovVideo(uint64(m.Favorite))
	}
	if m.Share != 0 {
		n += 1 + sovVideo(uint64(m.Share))
	}
	if len(m.KeyWords) > 0 {
		for _, s := range m.KeyWords {
			l = len(s)
			n += 1 + l + sovVideo(uint64(l))
		}
	}
	return n
}

func sovVideo(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozVideo(x uint64) (n int) {
	return sovVideo(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *Video) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowVideo
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
			return fmt.Errorf("proto: Video: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Video: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ID", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowVideo
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
				return ErrInvalidLengthVideo
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthVideo
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.ID = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Title", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowVideo
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
				return ErrInvalidLengthVideo
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthVideo
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Title = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Author", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowVideo
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
				return ErrInvalidLengthVideo
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthVideo
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Author = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 4:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field PostTime", wireType)
			}
			m.PostTime = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowVideo
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.PostTime |= int64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 5:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field View", wireType)
			}
			m.View = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowVideo
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.View |= int32(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 6:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Like", wireType)
			}
			m.Like = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowVideo
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Like |= int32(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 7:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Coin", wireType)
			}
			m.Coin = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowVideo
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Coin |= int32(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 8:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Favorite", wireType)
			}
			m.Favorite = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowVideo
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Favorite |= int32(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 9:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Share", wireType)
			}
			m.Share = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowVideo
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Share |= int32(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 10:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field KeyWords", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowVideo
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
				return ErrInvalidLengthVideo
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthVideo
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.KeyWords = append(m.KeyWords, string(dAtA[iNdEx:postIndex]))
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipVideo(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthVideo
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
func skipVideo(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowVideo
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
					return 0, ErrIntOverflowVideo
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
					return 0, ErrIntOverflowVideo
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
				return 0, ErrInvalidLengthVideo
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupVideo
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthVideo
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthVideo        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowVideo          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupVideo = fmt.Errorf("proto: unexpected end of group")
)
