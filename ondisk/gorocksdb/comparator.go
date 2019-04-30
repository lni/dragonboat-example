package gorocksdb

// #include "rocksdb/c.h"
import "C"

// A Comparator object provides a total order across slices that are
// used as keys in an sstable or a database.
type Comparator interface {
	// Three-way comparison. Returns value:
	//   < 0 iff "a" < "b",
	//   == 0 iff "a" == "b",
	//   > 0 iff "a" > "b"
	Compare(a, b []byte) int

	// The name of the comparator.
	Name() string
}

// NewNativeComparator creates a Comparator object.
func NewNativeComparator(c *C.rocksdb_comparator_t) Comparator {
	return nativeComparator{c}
}

type nativeComparator struct {
	c *C.rocksdb_comparator_t
}

func (c nativeComparator) Compare(a, b []byte) int { return 0 }
func (c nativeComparator) Name() string            { return "" }

// Hold references to comperators.
var comperators []Comparator

func registerComperator(cmp Comparator) int {
	comperators = append(comperators, cmp)
	return len(comperators) - 1
}
