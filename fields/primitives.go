package fields

import (
	"bytes"
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"math"
)

const (
	// CurrentVersion is the Forest version that this library writes
	CurrentVersion Version = 1

	// HashDigestLengthSHA512_256 is the length of the digest produced by the SHA512/256 hash algorithm
	HashDigestLengthSHA512_256 ContentLength = 32
)

// multiByteSerializationOrder defines the order in which multi-byte
// integers are serialized into binary
var multiByteSerializationOrder binary.ByteOrder = binary.BigEndian

// fundamental types
type genericType uint8

const sizeofgenericType = 1

func (g genericType) MarshalBinary() ([]byte, error) {
	b := new(bytes.Buffer)
	err := binary.Write(b, multiByteSerializationOrder, g)
	return b.Bytes(), err
}

func (g *genericType) UnmarshalBinary(b []byte) error {
	buf := bytes.NewBuffer(b)
	return binary.Read(buf, multiByteSerializationOrder, g)
}

func (g *genericType) SizeConstraints() (int, bool) {
	return sizeofgenericType, false
}

func (g *genericType) BytesConsumed() int {
	return sizeofgenericType
}

func (g *genericType) Equals(g2 *genericType) bool {
	return *g == *g2
}

// ContentLength represents the length of a piece of data in the Forest
type ContentLength uint16

const sizeofContentLength = 2

const (
	// MaxContentLength is the maximum representable content length in this
	// version of the Forest
	MaxContentLength = math.MaxUint16
)

func NewContentLength(size int) (*ContentLength, error) {
	if size > MaxContentLength {
		return nil, fmt.Errorf("Cannot represent content of size %d, max is %d", size, MaxContentLength)
	}
	c := ContentLength(size)
	return &c, nil
}

// MarshalBinary converts the ContentLength into its binary representation
func (c ContentLength) MarshalBinary() ([]byte, error) {
	b := new(bytes.Buffer)
	err := binary.Write(b, multiByteSerializationOrder, c)
	return b.Bytes(), err
}

func (c ContentLength) MarshalText() ([]byte, error) {
	return []byte(fmt.Sprintf("B%d", c)), nil
}

// UnmarshalBinary converts from the binary representation of a ContentLength
// back to its structured form
func (c *ContentLength) UnmarshalBinary(b []byte) error {
	buf := bytes.NewBuffer(b)
	return binary.Read(buf, multiByteSerializationOrder, c)
}

func (c *ContentLength) BytesConsumed() int {
	return sizeofContentLength
}

func (c *ContentLength) Equals(c2 *ContentLength) bool {
	return *c == *c2
}

// TreeDepth represents the depth of a node within a tree
type TreeDepth uint32

const sizeofTreeDepth = 4

// MarshalBinary converts the TreeDepth into its binary representation
func (t TreeDepth) MarshalBinary() ([]byte, error) {
	b := new(bytes.Buffer)
	err := binary.Write(b, multiByteSerializationOrder, t)
	return b.Bytes(), err
}

func (t TreeDepth) MarshalText() ([]byte, error) {
	return []byte(fmt.Sprintf("L%d", t)), nil
}

// UnmarshalBinary converts from the binary representation of a TreeDepth
// back to its structured form
func (t *TreeDepth) UnmarshalBinary(b []byte) error {
	buf := bytes.NewBuffer(b)
	return binary.Read(buf, multiByteSerializationOrder, t)
}

func (t *TreeDepth) BytesConsumed() int {
	return sizeofTreeDepth
}

func (t *TreeDepth) Equals(t2 *TreeDepth) bool {
	return *t == *t2
}

// Blob represents a quantity of arbitrary binary data in the Forest
type Blob []byte

// MarshalBinary converts the Blob into its binary representation
func (v Blob) MarshalBinary() ([]byte, error) {
	return v, nil
}

func (v Blob) MarshalText() ([]byte, error) {
	based := base64.RawURLEncoding.EncodeToString([]byte(v))
	return []byte(based), nil
}

// UnmarshalBinary converts from the binary representation of a Blob
// back to its structured form
func (v *Blob) UnmarshalBinary(b []byte) error {
	*v = b
	return nil
}

func (v *Blob) BytesConsumed() int {
	return len([]byte(*v))
}

func (v *Blob) Equals(v2 *Blob) bool {
	return bytes.Equal([]byte(*v), []byte(*v2))
}

// Version represents the version of the Arbor Forest Schema used to construct
// a particular node
type Version uint64

const sizeofVersion = 8

// MarshalBinary converts the Version into its binary representation
func (v Version) MarshalBinary() ([]byte, error) {
	b := new(bytes.Buffer)
	err := binary.Write(b, multiByteSerializationOrder, v)
	return b.Bytes(), err
}

func (v Version) MarshalText() ([]byte, error) {
	return []byte(fmt.Sprintf("V%d", v)), nil
}

// UnmarshalBinary converts from the binary representation of a Version
// back to its structured form
func (v *Version) UnmarshalBinary(b []byte) error {
	buf := bytes.NewBuffer(b)
	return binary.Read(buf, multiByteSerializationOrder, v)
}

func (v *Version) BytesConsumed() int {
	return sizeofVersion
}

func (v *Version) Equals(v2 *Version) bool {
	return *v == *v2
}

// specialized types
type NodeType genericType

const (
	sizeofNodeType            = sizeofgenericType
	NodeTypeIdentity NodeType = iota
	NodeTypeCommunity
	NodeTypeReply
)

var ValidNodeTypes = map[NodeType]struct{}{
	NodeTypeIdentity:  struct{}{},
	NodeTypeCommunity: struct{}{},
	NodeTypeReply:     struct{}{},
}

var nodeTypeNames = map[NodeType]string{
	NodeTypeIdentity:  "identity",
	NodeTypeCommunity: "community",
	NodeTypeReply:     "reply",
}

func (t NodeType) MarshalBinary() ([]byte, error) {
	return genericType(t).MarshalBinary()
}

func (t NodeType) MarshalText() ([]byte, error) {
	return []byte(nodeTypeNames[t]), nil
}

func (t *NodeType) UnmarshalBinary(b []byte) error {
	if err := (*genericType)(t).UnmarshalBinary(b); err != nil {
		return err
	}
	if _, valid := ValidNodeTypes[*t]; !valid {
		return fmt.Errorf("%d is not a valid node type", *t)
	}
	return nil
}

func (t *NodeType) BytesConsumed() int {
	return sizeofNodeType
}

func (t *NodeType) Equals(t2 *NodeType) bool {
	return ((*genericType)(t)).Equals((*genericType)(t2))
}

type HashType genericType

const (
	sizeofHashType            = sizeofgenericType
	HashTypeNullHash HashType = iota
	HashTypeSHA512
)

// map to valid lengths
var ValidHashTypes = map[HashType][]ContentLength{
	HashTypeNullHash: []ContentLength{0},
	HashTypeSHA512:   []ContentLength{HashDigestLengthSHA512_256},
}

var hashNames = map[HashType]string{
	HashTypeNullHash: "NullHash",
	HashTypeSHA512:   "SHA512",
}

func (t HashType) MarshalBinary() ([]byte, error) {
	return genericType(t).MarshalBinary()
}

func (t HashType) MarshalText() ([]byte, error) {
	return []byte(hashNames[t]), nil
}

func (t *HashType) UnmarshalBinary(b []byte) error {
	if err := (*genericType)(t).UnmarshalBinary(b); err != nil {
		return err
	}
	if _, valid := ValidHashTypes[*t]; !valid {
		return fmt.Errorf("%d is not a valid hash type", *t)
	}
	return nil
}

func (t *HashType) BytesConsumed() int {
	return sizeofHashType
}

func (t *HashType) Equals(t2 *HashType) bool {
	return ((*genericType)(t)).Equals((*genericType)(t2))
}

type ContentType genericType

const (
	sizeofContentType                 = sizeofgenericType
	ContentTypeUTF8String ContentType = 1
	ContentTypeJSON       ContentType = 2
)

var ValidContentTypes = map[ContentType]struct{}{
	ContentTypeUTF8String: struct{}{},
	ContentTypeJSON:       struct{}{},
}

var contentNames = map[ContentType]string{
	ContentTypeUTF8String: "UTF-8",
	ContentTypeJSON:       "JSON",
}

func (t ContentType) MarshalBinary() ([]byte, error) {
	return genericType(t).MarshalBinary()
}

func (t ContentType) MarshalText() ([]byte, error) {
	return []byte(contentNames[t]), nil
}

func (t *ContentType) UnmarshalBinary(b []byte) error {
	if err := (*genericType)(t).UnmarshalBinary(b); err != nil {
		return err
	}
	if _, valid := ValidContentTypes[*t]; !valid {
		return fmt.Errorf("%d is not a valid content type", *t)
	}
	return nil
}

func (t *ContentType) BytesConsumed() int {
	return sizeofContentType
}

func (t *ContentType) Equals(t2 *ContentType) bool {
	return ((*genericType)(t)).Equals((*genericType)(t2))
}

type KeyType genericType

const (
	sizeofKeyType          = sizeofgenericType
	KeyTypeNoKey   KeyType = 0
	KeyTypeOpenPGP KeyType = 1
)

var ValidKeyTypes = map[KeyType]struct{}{
	KeyTypeNoKey:   struct{}{},
	KeyTypeOpenPGP: struct{}{},
}

var keyNames = map[KeyType]string{
	KeyTypeNoKey:   "None",
	KeyTypeOpenPGP: "OpenPGP",
}

func (t KeyType) MarshalBinary() ([]byte, error) {
	return genericType(t).MarshalBinary()
}

func (t KeyType) MarshalText() ([]byte, error) {
	return []byte(keyNames[t]), nil
}

func (t *KeyType) UnmarshalBinary(b []byte) error {
	if err := (*genericType)(t).UnmarshalBinary(b); err != nil {
		return err
	}
	if _, valid := ValidKeyTypes[*t]; !valid {
		return fmt.Errorf("%d is not a valid key type", *t)
	}
	return nil
}

func (t *KeyType) BytesConsumed() int {
	return sizeofKeyType
}

func (t *KeyType) Equals(t2 *KeyType) bool {
	return ((*genericType)(t)).Equals((*genericType)(t2))
}

type SignatureType genericType

const (
	sizeofSignatureType                = sizeofgenericType
	SignatureTypeOpenPGP SignatureType = 1
)

var ValidSignatureTypes = map[SignatureType]struct{}{
	SignatureTypeOpenPGP: struct{}{},
}

var signatureNames = map[SignatureType]string{
	SignatureTypeOpenPGP: "OpenPGP",
}

func (t SignatureType) MarshalBinary() ([]byte, error) {
	return genericType(t).MarshalBinary()
}

func (t SignatureType) MarshalText() ([]byte, error) {
	return []byte(signatureNames[t]), nil
}

func (t *SignatureType) UnmarshalBinary(b []byte) error {
	if err := (*genericType)(t).UnmarshalBinary(b); err != nil {
		return err
	}
	if _, valid := ValidSignatureTypes[*t]; !valid {
		return fmt.Errorf("%d is not a valid signature type", *t)
	}
	return nil
}

func (t *SignatureType) BytesConsumed() int {
	return sizeofSignatureType
}

func (t *SignatureType) Equals(t2 *SignatureType) bool {
	return ((*genericType)(t)).Equals((*genericType)(t2))
}
