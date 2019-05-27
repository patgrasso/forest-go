package fields

import (
	"bytes"
	"encoding"
)

const sizeofDescriptor = sizeofgenericType + sizeofContentLength

func marshalTextDescriptor(descriptorType encoding.TextMarshaler, length encoding.TextMarshaler) ([]byte, error) {
	buf := new(bytes.Buffer)
	b, err := descriptorType.MarshalText()
	if err != nil {
		return nil, err
	}
	_, _ = buf.Write(b)
	_, _ = buf.Write([]byte("_"))
	b, err = length.MarshalText()
	if err != nil {
		return nil, err
	}
	_, _ = buf.Write(b)
	return buf.Bytes(), nil
}

// concrete descriptors
type HashDescriptor struct {
	Type   HashType
	Length ContentLength
}

const sizeofHashDescriptor = sizeofDescriptor

func NewHashDescriptor(t HashType, length int) (*HashDescriptor, error) {
	cLength, err := NewContentLength(length)
	if err != nil {
		return nil, err
	}
	return &HashDescriptor{t, *cLength}, nil
}

func (d *HashDescriptor) SerializationOrder() []BidirectionalBinaryMarshaler {
	return []BidirectionalBinaryMarshaler{&d.Type, &d.Length}
}

func (d *HashDescriptor) Equals(other *HashDescriptor) bool {
	return d.Type.Equals(&other.Type) && d.Length.Equals(&other.Length)
}

func (d *HashDescriptor) MarshalText() ([]byte, error) {
	return marshalTextDescriptor(d.Type, d.Length)
}

type ContentDescriptor struct {
	Type   ContentType
	Length ContentLength
}

const sizeofContentDescriptor = sizeofDescriptor

func NewContentDescriptor(t ContentType, length int) (*ContentDescriptor, error) {
	cLength, err := NewContentLength(length)
	if err != nil {
		return nil, err
	}
	return &ContentDescriptor{t, *cLength}, nil
}

func (d *ContentDescriptor) SerializationOrder() []BidirectionalBinaryMarshaler {
	return []BidirectionalBinaryMarshaler{&d.Type, &d.Length}
}

func (d *ContentDescriptor) Equals(other *ContentDescriptor) bool {
	return d.Type.Equals(&other.Type) && d.Length.Equals(&other.Length)
}

func (d *ContentDescriptor) MarshalText() ([]byte, error) {
	return marshalTextDescriptor(d.Type, d.Length)
}

type SignatureDescriptor struct {
	Type   SignatureType
	Length ContentLength
}

const sizeofSignatureDescriptor = sizeofDescriptor

func NewSignatureDescriptor(t SignatureType, length int) (*SignatureDescriptor, error) {
	cLength, err := NewContentLength(length)
	if err != nil {
		return nil, err
	}
	return &SignatureDescriptor{t, *cLength}, nil
}

func (d *SignatureDescriptor) SerializationOrder() []BidirectionalBinaryMarshaler {
	return []BidirectionalBinaryMarshaler{&d.Type, &d.Length}
}

func (d *SignatureDescriptor) Equals(other *SignatureDescriptor) bool {
	return d.Type.Equals(&other.Type) && d.Length.Equals(&other.Length)
}

func (d *SignatureDescriptor) MarshalText() ([]byte, error) {
	return marshalTextDescriptor(d.Type, d.Length)
}

type KeyDescriptor struct {
	Type   KeyType
	Length ContentLength
}

const sizeofKeyDescriptor = sizeofDescriptor

func NewKeyDescriptor(t KeyType, length int) (*KeyDescriptor, error) {
	cLength, err := NewContentLength(length)
	if err != nil {
		return nil, err
	}
	return &KeyDescriptor{t, *cLength}, nil
}

func (d *KeyDescriptor) SerializationOrder() []BidirectionalBinaryMarshaler {
	return []BidirectionalBinaryMarshaler{&d.Type, &d.Length}
}

func (d *KeyDescriptor) Equals(other *KeyDescriptor) bool {
	return d.Type.Equals(&other.Type) && d.Length.Equals(&other.Length)
}

func (d *KeyDescriptor) MarshalText() ([]byte, error) {
	return marshalTextDescriptor(d.Type, d.Length)
}
