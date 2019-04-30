package gorocksdb

// #include "rocksdb/c.h"
import "C"

// A SliceTransform can be used as a prefix extractor.
type SliceTransform interface {
	// Transform a src in domain to a dst in the range.
	Transform(src []byte) []byte

	// Determine whether this is a valid src upon the function applies.
	InDomain(src []byte) bool

	// Determine whether dst=Transform(src) for some src.
	InRange(src []byte) bool

	// Return the name of this transformation.
	Name() string
}

// NewFixedPrefixTransform creates a new fixed prefix transform.
func NewFixedPrefixTransform(prefixLen int) SliceTransform {
	return NewNativeSliceTransform(C.rocksdb_slicetransform_create_fixed_prefix(C.size_t(prefixLen)))
}

// NewNativeSliceTransform creates a SliceTransform object.
func NewNativeSliceTransform(c *C.rocksdb_slicetransform_t) SliceTransform {
	return nativeSliceTransform{c}
}

type nativeSliceTransform struct {
	c *C.rocksdb_slicetransform_t
}

func (st nativeSliceTransform) Transform(src []byte) []byte { return nil }
func (st nativeSliceTransform) InDomain(src []byte) bool    { return false }
func (st nativeSliceTransform) InRange(src []byte) bool     { return false }
func (st nativeSliceTransform) Name() string                { return "" }

// Hold references to slice transforms.
var sliceTransforms []SliceTransform

func registerSliceTransform(st SliceTransform) int {
	sliceTransforms = append(sliceTransforms, st)
	return len(sliceTransforms) - 1
}
