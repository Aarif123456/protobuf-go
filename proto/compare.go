// Copyright 2019 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package proto

import (
	"reflect"

	"google.golang.org/protobuf/reflect/protoreflect"
)

// Compare recursively compares the two message. Compare will only return 0
// if and only if proto.Equal will return true.
//
// This function is meant to be used to sort arbitrary messages.
// For example, if you want to ensure two arrays contained the same elements
// and didn't care about the ordering.
func Compare(x, y Message) int {
	if x == nil || y == nil {
		return boolCompare(x != nil, y != nil)
	}
	if reflect.TypeOf(x).Kind() == reflect.Ptr && x == y {
		// Avoid an expensive comparison if both inputs are identical pointers.
		return 0
	}

	mx := x.ProtoReflect()
	my := y.ProtoReflect()

	if validCmp := boolCompare(mx.IsValid(), my.IsValid()); validCmp != 0 {
		return validCmp
	}

	vx := protoreflect.ValueOfMessage(mx)
	vy := protoreflect.ValueOfMessage(my)

	return vx.Compare(vy)
}

// LessThan returns true is x is less than y.
// It has implements the 3 properties required to be used as a Less function
// by Go's "sort" package
//
//   - Deterministic: less(x, y) == less(x, y)
//   - Irreflexive: !less(x, x)
//   - Transitive: if !less(x, y) and !less(y, z), then !less(x, z)
func LessThan(x, y Message) bool {
	return Compare(x, y) == -1
}

func boolCompare[T ~bool](x, y T) int {
	if x == y {
		return 0
	}

	if !x {
		return -1
	}

	return 1
}
