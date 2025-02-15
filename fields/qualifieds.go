package fields

import (
	"bytes"
	"encoding"
	"fmt"
)

const minSizeofQualified = sizeofDescriptor

// concrete qualified data types
type QualifiedHash struct {
	Descriptor HashDescriptor
	Blob      Blob
}

const minSizeofQualifiedHash = sizeofHashDescriptor

func marshalTextQualified(first, second encoding.TextMarshaler) ([]byte, error) {
	buf := new(bytes.Buffer)
	b, err := first.MarshalText()
	if err != nil {
		return nil, err
	}
	_, _ = buf.Write(b)
	_, _ = buf.Write([]byte("__"))
	b, err = second.MarshalText()
	if err != nil {
		return nil, err
	}
	_, _ = buf.Write(b)
	return buf.Bytes(), nil
}

// NewQualifiedHash returns a valid QualifiedHash from the given data
func NewQualifiedHash(t HashType, content []byte) (*QualifiedHash, error) {
	hd, err := NewHashDescriptor(t, len(content))
	if err != nil {
		return nil, err
	}
	return &QualifiedHash{*hd, Blob(content)}, nil
}

func NullHash() *QualifiedHash {
	return &QualifiedHash{
		Descriptor: HashDescriptor{
			Type:   HashTypeNullHash,
			Length: 0,
		},
		Blob: []byte{},
	}
}

func (q *QualifiedHash) UnmarshalBinary(b []byte) error {
	unused, err := UnmarshalAll(b, AsUnmarshaler(q.Descriptor.SerializationOrder())...)
	if err != nil {
		return err
	}
	if err := q.Blob.UnmarshalBinary(unused[:q.Descriptor.Length]); err != nil {
		return err
	}
	return nil
}

func (q *QualifiedHash) MarshalBinary() ([]byte, error) {
	buf := new(bytes.Buffer)
	if err := MarshalAllInto(buf, AsMarshaler(q.SerializationOrder())...); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (q *QualifiedHash) BytesConsumed() int {
	return TotalBytesConsumed(q.SerializationOrder()...)
}

func (q *QualifiedHash) SerializationOrder() []BidirectionalBinaryMarshaler {
	return append(q.Descriptor.SerializationOrder(), &q.Blob)
}

func (q *QualifiedHash) Equals(other *QualifiedHash) bool {
	return q.Descriptor.Equals(&other.Descriptor) && q.Blob.Equals(&other.Blob)
}

func (q *QualifiedHash) MarshalText() ([]byte, error) {
	return marshalTextQualified(&q.Descriptor, q.Blob)
}

func (q *QualifiedHash) MarshalString() (string, error) {
	s, e := q.MarshalText()
	return string(s), e
}

func (q *QualifiedHash) Validate() error {
	if err := q.Descriptor.Validate(); err != nil {
		return err
	}
	if int(q.Descriptor.Length) != len(q.Blob) {
		return fmt.Errorf("Descriptor length %d does not match value length %d", q.Descriptor.Length, len(q.Blob))
	}
	return nil
}

type QualifiedContent struct {
	Descriptor ContentDescriptor
	Blob      Blob
}

const minSizeofQualifiedContent = sizeofContentDescriptor

// NewQualifiedContent returns a valid QualifiedContent from the given data
func NewQualifiedContent(t ContentType, content []byte) (*QualifiedContent, error) {
	hd, err := NewContentDescriptor(t, len(content))
	if err != nil {
		return nil, err
	}
	return &QualifiedContent{*hd, Blob(content)}, nil
}

func (q *QualifiedContent) SerializationOrder() []BidirectionalBinaryMarshaler {
	return append(q.Descriptor.SerializationOrder(), &q.Blob)
}

func (q *QualifiedContent) Equals(other *QualifiedContent) bool {
	return q.Descriptor.Equals(&other.Descriptor) && q.Blob.Equals(&other.Blob)
}

func (q *QualifiedContent) UnmarshalBinary(b []byte) error {
	unused, err := UnmarshalAll(b, AsUnmarshaler(q.Descriptor.SerializationOrder())...)
	if err != nil {
		return err
	}
	if err := q.Blob.UnmarshalBinary(unused[:q.Descriptor.Length]); err != nil {
		return err
	}
	return nil
}

func (q *QualifiedContent) MarshalBinary() ([]byte, error) {
	buf := new(bytes.Buffer)
	if err := MarshalAllInto(buf, AsMarshaler(q.SerializationOrder())...); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (q *QualifiedContent) BytesConsumed() int {
	return TotalBytesConsumed(q.SerializationOrder()...)
}

func (q *QualifiedContent) MarshalText() ([]byte, error) {
	switch q.Descriptor.Type {
	case ContentTypeUTF8String:
		fallthrough
	case ContentTypeJSON:
		descText, err := (&q.Descriptor).MarshalText()
		if err != nil {
			return nil, err
		}
		buf := bytes.NewBuffer(descText)
		_, err = buf.Write(q.Blob)
		if err != nil {
			return nil, err
		}
		return buf.Bytes(), nil
	default:
		return marshalTextQualified(&q.Descriptor, q.Blob)
	}
}

func (q *QualifiedContent) Validate() error {
	if err := q.Descriptor.Validate(); err != nil {
		return err
	}
	if int(q.Descriptor.Length) != len(q.Blob) {
		return fmt.Errorf("Descriptor length %d does not match value length %d", q.Descriptor.Length, len(q.Blob))
	}
	return nil
}

type QualifiedKey struct {
	Descriptor KeyDescriptor
	Blob      Blob
}

const minSizeofQualifiedKey = sizeofKeyDescriptor

// NewQualifiedKey returns a valid QualifiedKey from the given data
func NewQualifiedKey(t KeyType, content []byte) (*QualifiedKey, error) {
	hd, err := NewKeyDescriptor(t, len(content))
	if err != nil {
		return nil, err
	}
	return &QualifiedKey{*hd, Blob(content)}, nil
}

func (q *QualifiedKey) SerializationOrder() []BidirectionalBinaryMarshaler {
	return append(q.Descriptor.SerializationOrder(), &q.Blob)
}

func (q *QualifiedKey) Equals(other *QualifiedKey) bool {
	return q.Descriptor.Equals(&other.Descriptor) && q.Blob.Equals(&other.Blob)
}

func (q *QualifiedKey) UnmarshalBinary(b []byte) error {
	unused, err := UnmarshalAll(b, AsUnmarshaler(q.Descriptor.SerializationOrder())...)
	if err != nil {
		return err
	}
	if err := q.Blob.UnmarshalBinary(unused[:q.Descriptor.Length]); err != nil {
		return err
	}
	return nil
}

func (q *QualifiedKey) MarshalBinary() ([]byte, error) {
	buf := new(bytes.Buffer)
	if err := MarshalAllInto(buf, AsMarshaler(q.SerializationOrder())...); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (q *QualifiedKey) BytesConsumed() int {
	return TotalBytesConsumed(q.SerializationOrder()...)
}

func (q *QualifiedKey) MarshalText() ([]byte, error) {
	return marshalTextQualified(&q.Descriptor, q.Blob)
}

func (q *QualifiedKey) Validate() error {
	if err := q.Descriptor.Validate(); err != nil {
		return err
	}
	if int(q.Descriptor.Length) != len(q.Blob) {
		return fmt.Errorf("Descriptor length %d does not match value length %d", q.Descriptor.Length, len(q.Blob))
	}
	return nil
}

type QualifiedSignature struct {
	Descriptor SignatureDescriptor
	Blob      Blob
}

const minSizeofQualifiedSignature = sizeofSignatureDescriptor

// NewQualifiedSignature returns a valid QualifiedSignature from the given data
func NewQualifiedSignature(t SignatureType, content []byte) (*QualifiedSignature, error) {
	hd, err := NewSignatureDescriptor(t, len(content))
	if err != nil {
		return nil, err
	}
	return &QualifiedSignature{*hd, Blob(content)}, nil
}

func (q *QualifiedSignature) SerializationOrder() []BidirectionalBinaryMarshaler {
	return append(q.Descriptor.SerializationOrder(), &q.Blob)
}

func (q *QualifiedSignature) Equals(other *QualifiedSignature) bool {
	return q.Descriptor.Equals(&other.Descriptor) && q.Blob.Equals(&other.Blob)
}

func (q *QualifiedSignature) UnmarshalBinary(b []byte) error {
	unused, err := UnmarshalAll(b, AsUnmarshaler(q.Descriptor.SerializationOrder())...)
	if err != nil {
		return err
	}
	if err := q.Blob.UnmarshalBinary(unused[:q.Descriptor.Length]); err != nil {
		return err
	}
	return nil
}

func (q *QualifiedSignature) MarshalBinary() ([]byte, error) {
	buf := new(bytes.Buffer)
	if err := MarshalAllInto(buf, AsMarshaler(q.SerializationOrder())...); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (q *QualifiedSignature) BytesConsumed() int {
	return TotalBytesConsumed(q.SerializationOrder()...)
}

func (q *QualifiedSignature) MarshalText() ([]byte, error) {
	return marshalTextQualified(&q.Descriptor, q.Blob)
}

func (q *QualifiedSignature) Validate() error {
	if err := q.Descriptor.Validate(); err != nil {
		return err
	}
	if int(q.Descriptor.Length) != len(q.Blob) {
		return fmt.Errorf("Descriptor length %d does not match value length %d", q.Descriptor.Length, len(q.Blob))
	}
	return nil
}
