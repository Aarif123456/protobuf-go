// Copyright 2019 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package proto_test

import (
	"math"
	"testing"

	"google.golang.org/protobuf/encoding/prototext"
	"google.golang.org/protobuf/internal/pragma"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/testing/protopack"

	testpb "google.golang.org/protobuf/internal/testprotos/test"
	test3pb "google.golang.org/protobuf/internal/testprotos/test3"
	testeditionspb "google.golang.org/protobuf/internal/testprotos/testeditions"
)

func TestCompare(t *testing.T) {
	identicalPtrPb := &testpb.TestAllTypes{MapStringString: map[string]string{"a": "b", "c": "d"}}

	type incomparableMessage struct {
		*testpb.TestAllTypes
		pragma.DoNotCompare
	}

	tests := []struct {
		x, y proto.Message
		want int
	}{
		{
			x:    nil,
			y:    nil,
			want: 0,
		},
		{
			x:    nil,
			y:    (*testpb.TestAllTypes)(nil),
			want: -1,
		},
		{
			x:    (*testpb.TestAllTypes)(nil),
			y:    nil,
			want: 1,
		},
		{
			x:    (*testpb.TestAllTypes)(nil),
			y:    (*testpb.TestAllTypes)(nil),
			want: 0,
		},
		{
			x:    (*testpb.TestAllTypes)(nil),
			y:    new(testpb.TestAllTypes),
			want: -1,
		},
		{
			x:    new(testpb.TestAllTypes),
			y:    (*testpb.TestAllTypes)(nil),
			want: 1,
		},
		{
			x:    new(testpb.TestAllTypes),
			y:    new(testpb.TestAllTypes),
			want: 0,
		},
		{
			x:    (*testpb.TestAllTypes)(nil),
			y:    (*testpb.TestAllExtensions)(nil),
			want: 1,
		},
		{
			x:    (*testpb.TestAllExtensions)(nil),
			y:    (*testpb.TestAllTypes)(nil),
			want: -1,
		},
		{
			x:    (*testpb.TestAllTypes)(nil),
			y:    new(testpb.TestAllExtensions),
			want: -1,
		},
		{
			x:    new(testpb.TestAllExtensions),
			y:    (*testpb.TestAllTypes)(nil),
			want: 1,
		},
		{
			x:    (*testpb.TestAllExtensions)(nil),
			y:    new(testpb.TestAllTypes),
			want: -1,
		},
		{
			x:    new(testpb.TestAllTypes),
			y:    (*testpb.TestAllExtensions)(nil),
			want: 1,
		},
		{
			x:    new(testpb.TestAllExtensions),
			y:    new(testpb.TestAllTypes),
			want: -1,
		},
		{
			x:    new(testpb.TestAllTypes),
			y:    new(testpb.TestAllExtensions),
			want: 1,
		},

		// Identical input pointers
		{
			x:    identicalPtrPb,
			y:    identicalPtrPb,
			want: 0,
		},

		// Incomparable types. The top-level types are not actually directly
		// compared (which would panic), but rather the comparison happens on the
		// objects returned by ProtoReflect(). These tests are here just to ensure
		// that any short-circuit checks do not accidentally try to compare
		// incomparable top-level types.
		{
			x:    incomparableMessage{TestAllTypes: identicalPtrPb},
			y:    incomparableMessage{TestAllTypes: identicalPtrPb},
			want: 0,
		},
		{
			x:    identicalPtrPb,
			y:    incomparableMessage{TestAllTypes: identicalPtrPb},
			want: 0,
		},
		{
			x:    identicalPtrPb,
			y:    &incomparableMessage{TestAllTypes: identicalPtrPb},
			want: 0,
		},

		// Proto2 scalars.
		{
			x: &testpb.TestAllTypes{
				OptionalInt32: proto.Int32(-1),
				OptionalInt64: proto.Int64(1),
			},
			y: &testpb.TestAllTypes{
				OptionalInt32: proto.Int32(1),
				OptionalInt64: proto.Int64(-1),
			},
			want: -1,
		},
		{
			x: &testpb.TestAllTypes{
				OptionalInt32: proto.Int32(1),
				OptionalInt64: proto.Int64(-1),
			},
			y: &testpb.TestAllTypes{
				OptionalInt32: proto.Int32(-1),
				OptionalInt64: proto.Int64(1),
			},
			want: 1,
		},
		{
			x:    &testpb.TestAllTypes{OptionalInt32: proto.Int32(1)},
			y:    &testpb.TestAllTypes{OptionalInt64: proto.Int64(1)},
			want: -1,
		},
		{
			x:    &testpb.TestAllTypes{OptionalInt64: proto.Int64(1)},
			y:    &testpb.TestAllTypes{OptionalInt32: proto.Int32(1)},
			want: 1,
		},
		{
			x:    &testpb.TestAllTypes{OptionalInt32: proto.Int32(1)},
			y:    &testpb.TestAllTypes{OptionalInt32: proto.Int32(2)},
			want: -1,
		},
		{
			x:    &testpb.TestAllTypes{OptionalInt32: proto.Int32(2)},
			y:    &testpb.TestAllTypes{OptionalInt32: proto.Int32(1)},
			want: 1,
		},
		{
			x:    &testpb.TestAllTypes{OptionalInt64: proto.Int64(1)},
			y:    &testpb.TestAllTypes{OptionalInt64: proto.Int64(2)},
			want: -1,
		},
		{
			x:    &testpb.TestAllTypes{OptionalInt64: proto.Int64(2)},
			y:    &testpb.TestAllTypes{OptionalInt64: proto.Int64(1)},
			want: 1,
		},
		{
			x:    &testpb.TestAllTypes{OptionalUint32: proto.Uint32(1)},
			y:    &testpb.TestAllTypes{OptionalUint32: proto.Uint32(2)},
			want: -1,
		},
		{
			x:    &testpb.TestAllTypes{OptionalUint32: proto.Uint32(2)},
			y:    &testpb.TestAllTypes{OptionalUint32: proto.Uint32(1)},
			want: 1,
		},
		{
			x:    &testpb.TestAllTypes{OptionalUint64: proto.Uint64(1)},
			y:    &testpb.TestAllTypes{OptionalUint64: proto.Uint64(2)},
			want: -1,
		},
		{
			x:    &testpb.TestAllTypes{OptionalUint64: proto.Uint64(2)},
			y:    &testpb.TestAllTypes{OptionalUint64: proto.Uint64(1)},
			want: 1,
		},
		{
			x:    &testpb.TestAllTypes{OptionalSint32: proto.Int32(1)},
			y:    &testpb.TestAllTypes{OptionalSint32: proto.Int32(2)},
			want: -1,
		},
		{
			x:    &testpb.TestAllTypes{OptionalSint32: proto.Int32(2)},
			y:    &testpb.TestAllTypes{OptionalSint32: proto.Int32(1)},
			want: 1,
		},
		{
			x:    &testpb.TestAllTypes{OptionalSint64: proto.Int64(1)},
			y:    &testpb.TestAllTypes{OptionalSint64: proto.Int64(2)},
			want: -1,
		},
		{
			x:    &testpb.TestAllTypes{OptionalSint64: proto.Int64(2)},
			y:    &testpb.TestAllTypes{OptionalSint64: proto.Int64(1)},
			want: 1,
		},
		{
			x:    &testpb.TestAllTypes{OptionalFixed32: proto.Uint32(1)},
			y:    &testpb.TestAllTypes{OptionalFixed32: proto.Uint32(2)},
			want: -1,
		},
		{
			x:    &testpb.TestAllTypes{OptionalFixed32: proto.Uint32(2)},
			y:    &testpb.TestAllTypes{OptionalFixed32: proto.Uint32(1)},
			want: 1,
		},
		{
			x:    &testpb.TestAllTypes{OptionalFixed64: proto.Uint64(1)},
			y:    &testpb.TestAllTypes{OptionalFixed64: proto.Uint64(2)},
			want: -1,
		},
		{
			x:    &testpb.TestAllTypes{OptionalFixed64: proto.Uint64(2)},
			y:    &testpb.TestAllTypes{OptionalFixed64: proto.Uint64(1)},
			want: 1,
		},
		{
			x:    &testpb.TestAllTypes{OptionalSfixed32: proto.Int32(1)},
			y:    &testpb.TestAllTypes{OptionalSfixed32: proto.Int32(2)},
			want: -1,
		},
		{
			x:    &testpb.TestAllTypes{OptionalSfixed32: proto.Int32(2)},
			y:    &testpb.TestAllTypes{OptionalSfixed32: proto.Int32(1)},
			want: 1,
		},
		{
			x:    &testpb.TestAllTypes{OptionalSfixed64: proto.Int64(1)},
			y:    &testpb.TestAllTypes{OptionalSfixed64: proto.Int64(2)},
			want: -1,
		},
		{
			x:    &testpb.TestAllTypes{OptionalSfixed64: proto.Int64(2)},
			y:    &testpb.TestAllTypes{OptionalSfixed64: proto.Int64(1)},
			want: 1,
		},
		{
			x:    &testpb.TestAllTypes{OptionalFloat: proto.Float32(1)},
			y:    &testpb.TestAllTypes{OptionalFloat: proto.Float32(2)},
			want: -1,
		},
		{
			x:    &testpb.TestAllTypes{OptionalFloat: proto.Float32(2)},
			y:    &testpb.TestAllTypes{OptionalFloat: proto.Float32(1)},
			want: 1,
		},
		{
			x:    &testpb.TestAllTypes{OptionalDouble: proto.Float64(1)},
			y:    &testpb.TestAllTypes{OptionalDouble: proto.Float64(2)},
			want: -1,
		},
		{
			x:    &testpb.TestAllTypes{OptionalDouble: proto.Float64(2)},
			y:    &testpb.TestAllTypes{OptionalDouble: proto.Float64(1)},
			want: 1,
		},
		{
			x:    &testpb.TestAllTypes{OptionalFloat: proto.Float32(float32(math.NaN()))},
			y:    &testpb.TestAllTypes{OptionalFloat: proto.Float32(0)},
			want: -1,
		},
		{
			x:    &testpb.TestAllTypes{OptionalFloat: proto.Float32(0)},
			y:    &testpb.TestAllTypes{OptionalFloat: proto.Float32(float32(math.NaN()))},
			want: 1,
		},
		{
			x:    &testpb.TestAllTypes{OptionalDouble: proto.Float64(float64(math.NaN()))},
			y:    &testpb.TestAllTypes{OptionalDouble: proto.Float64(0)},
			want: -1,
		},
		{
			x:    &testpb.TestAllTypes{OptionalDouble: proto.Float64(0)},
			y:    &testpb.TestAllTypes{OptionalDouble: proto.Float64(float64(math.NaN()))},
			want: 1,
		},
		{
			x:    &testpb.TestAllTypes{OptionalBool: proto.Bool(false)},
			y:    &testpb.TestAllTypes{OptionalBool: proto.Bool(true)},
			want: -1,
		},
		{
			x:    &testpb.TestAllTypes{OptionalBool: proto.Bool(true)},
			y:    &testpb.TestAllTypes{OptionalBool: proto.Bool(false)},
			want: 1,
		},
		{
			x:    &testpb.TestAllTypes{OptionalString: proto.String("a")},
			y:    &testpb.TestAllTypes{OptionalString: proto.String("b")},
			want: -1,
		},
		{
			x:    &testpb.TestAllTypes{OptionalString: proto.String("b")},
			y:    &testpb.TestAllTypes{OptionalString: proto.String("a")},
			want: 1,
		},
		{
			x:    &testpb.TestAllTypes{OptionalBytes: []byte("a")},
			y:    &testpb.TestAllTypes{OptionalBytes: []byte("b")},
			want: -1,
		},
		{
			x:    &testpb.TestAllTypes{OptionalBytes: []byte("b")},
			y:    &testpb.TestAllTypes{OptionalBytes: []byte("a")},
			want: 1,
		},
		{
			x:    &testpb.TestAllTypes{OptionalNestedEnum: testpb.TestAllTypes_FOO.Enum()},
			y:    &testpb.TestAllTypes{OptionalNestedEnum: testpb.TestAllTypes_BAR.Enum()},
			want: -1,
		},
		{
			x:    &testpb.TestAllTypes{OptionalNestedEnum: testpb.TestAllTypes_BAR.Enum()},
			y:    &testpb.TestAllTypes{OptionalNestedEnum: testpb.TestAllTypes_FOO.Enum()},
			want: 1,
		},
		{
			x:    &testpb.TestAllTypes{OptionalInt32: proto.Int32(2)},
			y:    &testpb.TestAllTypes{OptionalInt32: proto.Int32(2)},
			want: 0,
		},
		{
			x:    &testpb.TestAllTypes{OptionalInt64: proto.Int64(2)},
			y:    &testpb.TestAllTypes{OptionalInt64: proto.Int64(2)},
			want: 0,
		},
		{
			x:    &testpb.TestAllTypes{OptionalUint32: proto.Uint32(2)},
			y:    &testpb.TestAllTypes{OptionalUint32: proto.Uint32(2)},
			want: 0,
		},
		{
			x:    &testpb.TestAllTypes{OptionalUint64: proto.Uint64(2)},
			y:    &testpb.TestAllTypes{OptionalUint64: proto.Uint64(2)},
			want: 0,
		},
		{
			x:    &testpb.TestAllTypes{OptionalSint32: proto.Int32(2)},
			y:    &testpb.TestAllTypes{OptionalSint32: proto.Int32(2)},
			want: 0,
		},
		{
			x:    &testpb.TestAllTypes{OptionalSint64: proto.Int64(2)},
			y:    &testpb.TestAllTypes{OptionalSint64: proto.Int64(2)},
			want: 0,
		},
		{
			x:    &testpb.TestAllTypes{OptionalFixed32: proto.Uint32(2)},
			y:    &testpb.TestAllTypes{OptionalFixed32: proto.Uint32(2)},
			want: 0,
		},
		{
			x:    &testpb.TestAllTypes{OptionalFixed64: proto.Uint64(2)},
			y:    &testpb.TestAllTypes{OptionalFixed64: proto.Uint64(2)},
			want: 0,
		},
		{
			x:    &testpb.TestAllTypes{OptionalSfixed32: proto.Int32(2)},
			y:    &testpb.TestAllTypes{OptionalSfixed32: proto.Int32(2)},
			want: 0,
		},
		{
			x:    &testpb.TestAllTypes{OptionalSfixed64: proto.Int64(2)},
			y:    &testpb.TestAllTypes{OptionalSfixed64: proto.Int64(2)},
			want: 0,
		},
		{
			x:    &testpb.TestAllTypes{OptionalFloat: proto.Float32(2)},
			y:    &testpb.TestAllTypes{OptionalFloat: proto.Float32(2)},
			want: 0,
		},
		{
			x:    &testpb.TestAllTypes{OptionalDouble: proto.Float64(2)},
			y:    &testpb.TestAllTypes{OptionalDouble: proto.Float64(2)},
			want: 0,
		},
		{
			x:    &testpb.TestAllTypes{OptionalFloat: proto.Float32(float32(math.NaN()))},
			y:    &testpb.TestAllTypes{OptionalFloat: proto.Float32(float32(math.NaN()))},
			want: 0,
		},
		{
			x:    &testpb.TestAllTypes{OptionalDouble: proto.Float64(float64(math.NaN()))},
			y:    &testpb.TestAllTypes{OptionalDouble: proto.Float64(float64(math.NaN()))},
			want: 0,
		},
		{
			x:    &testpb.TestAllTypes{OptionalBool: proto.Bool(true)},
			y:    &testpb.TestAllTypes{OptionalBool: proto.Bool(true)},
			want: 0,
		},
		{
			x:    &testpb.TestAllTypes{OptionalString: proto.String("abc")},
			y:    &testpb.TestAllTypes{OptionalString: proto.String("abc")},
			want: 0,
		},
		{
			x:    &testpb.TestAllTypes{OptionalBytes: []byte("abc")},
			y:    &testpb.TestAllTypes{OptionalBytes: []byte("abc")},
			want: 0,
		},
		{
			x:    &testpb.TestAllTypes{OptionalNestedEnum: testpb.TestAllTypes_FOO.Enum()},
			y:    &testpb.TestAllTypes{OptionalNestedEnum: testpb.TestAllTypes_FOO.Enum()},
			want: 0,
		},

		// Editions scalars.
		{
			x:    &testeditionspb.TestAllTypes{OptionalInt32: proto.Int32(1)},
			y:    &testeditionspb.TestAllTypes{OptionalInt64: proto.Int64(2)},
			want: -1,
		},
		{
			x:    &testeditionspb.TestAllTypes{OptionalInt64: proto.Int64(2)},
			y:    &testeditionspb.TestAllTypes{OptionalInt32: proto.Int32(1)},
			want: 1,
		},
		{
			x:    &testeditionspb.TestAllTypes{OptionalInt32: proto.Int32(1)},
			y:    &testeditionspb.TestAllTypes{OptionalInt32: proto.Int32(2)},
			want: -1,
		},
		{
			x:    &testeditionspb.TestAllTypes{OptionalInt32: proto.Int32(2)},
			y:    &testeditionspb.TestAllTypes{OptionalInt32: proto.Int32(1)},
			want: 1,
		},
		{
			x:    &testeditionspb.TestAllTypes{OptionalInt64: proto.Int64(1)},
			y:    &testeditionspb.TestAllTypes{OptionalInt64: proto.Int64(2)},
			want: -1,
		},
		{
			x:    &testeditionspb.TestAllTypes{OptionalInt64: proto.Int64(2)},
			y:    &testeditionspb.TestAllTypes{OptionalInt64: proto.Int64(1)},
			want: 1,
		},
		{
			x:    &testeditionspb.TestAllTypes{OptionalUint32: proto.Uint32(1)},
			y:    &testeditionspb.TestAllTypes{OptionalUint32: proto.Uint32(2)},
			want: -1,
		},
		{
			x:    &testeditionspb.TestAllTypes{OptionalUint32: proto.Uint32(2)},
			y:    &testeditionspb.TestAllTypes{OptionalUint32: proto.Uint32(1)},
			want: 1,
		},
		{
			x:    &testeditionspb.TestAllTypes{OptionalUint64: proto.Uint64(1)},
			y:    &testeditionspb.TestAllTypes{OptionalUint64: proto.Uint64(2)},
			want: -1,
		},
		{
			x:    &testeditionspb.TestAllTypes{OptionalUint64: proto.Uint64(2)},
			y:    &testeditionspb.TestAllTypes{OptionalUint64: proto.Uint64(1)},
			want: 1,
		},
		{
			x:    &testeditionspb.TestAllTypes{OptionalSint32: proto.Int32(1)},
			y:    &testeditionspb.TestAllTypes{OptionalSint32: proto.Int32(2)},
			want: -1,
		},
		{
			x:    &testeditionspb.TestAllTypes{OptionalSint32: proto.Int32(2)},
			y:    &testeditionspb.TestAllTypes{OptionalSint32: proto.Int32(1)},
			want: 1,
		},
		{
			x:    &testeditionspb.TestAllTypes{OptionalSint64: proto.Int64(1)},
			y:    &testeditionspb.TestAllTypes{OptionalSint64: proto.Int64(2)},
			want: -1,
		},
		{
			x:    &testeditionspb.TestAllTypes{OptionalSint64: proto.Int64(2)},
			y:    &testeditionspb.TestAllTypes{OptionalSint64: proto.Int64(1)},
			want: 1,
		},
		{
			x:    &testeditionspb.TestAllTypes{OptionalFixed32: proto.Uint32(1)},
			y:    &testeditionspb.TestAllTypes{OptionalFixed32: proto.Uint32(2)},
			want: -1,
		},
		{
			x:    &testeditionspb.TestAllTypes{OptionalFixed32: proto.Uint32(2)},
			y:    &testeditionspb.TestAllTypes{OptionalFixed32: proto.Uint32(1)},
			want: 1,
		},
		{
			x:    &testeditionspb.TestAllTypes{OptionalFixed64: proto.Uint64(1)},
			y:    &testeditionspb.TestAllTypes{OptionalFixed64: proto.Uint64(2)},
			want: -1,
		},
		{
			x:    &testeditionspb.TestAllTypes{OptionalFixed64: proto.Uint64(2)},
			y:    &testeditionspb.TestAllTypes{OptionalFixed64: proto.Uint64(1)},
			want: 1,
		},
		{
			x:    &testeditionspb.TestAllTypes{OptionalSfixed32: proto.Int32(1)},
			y:    &testeditionspb.TestAllTypes{OptionalSfixed32: proto.Int32(2)},
			want: -1,
		},
		{
			x:    &testeditionspb.TestAllTypes{OptionalSfixed32: proto.Int32(2)},
			y:    &testeditionspb.TestAllTypes{OptionalSfixed32: proto.Int32(1)},
			want: 1,
		},
		{
			x:    &testeditionspb.TestAllTypes{OptionalSfixed64: proto.Int64(1)},
			y:    &testeditionspb.TestAllTypes{OptionalSfixed64: proto.Int64(2)},
			want: -1,
		},
		{
			x:    &testeditionspb.TestAllTypes{OptionalSfixed64: proto.Int64(2)},
			y:    &testeditionspb.TestAllTypes{OptionalSfixed64: proto.Int64(1)},
			want: 1,
		},
		{
			x:    &testeditionspb.TestAllTypes{OptionalFloat: proto.Float32(1)},
			y:    &testeditionspb.TestAllTypes{OptionalFloat: proto.Float32(2)},
			want: -1,
		},
		{
			x:    &testeditionspb.TestAllTypes{OptionalFloat: proto.Float32(2)},
			y:    &testeditionspb.TestAllTypes{OptionalFloat: proto.Float32(1)},
			want: 1,
		},
		{
			x:    &testeditionspb.TestAllTypes{OptionalDouble: proto.Float64(1)},
			y:    &testeditionspb.TestAllTypes{OptionalDouble: proto.Float64(2)},
			want: -1,
		},
		{
			x:    &testeditionspb.TestAllTypes{OptionalDouble: proto.Float64(2)},
			y:    &testeditionspb.TestAllTypes{OptionalDouble: proto.Float64(1)},
			want: 1,
		},
		{
			x:    &testeditionspb.TestAllTypes{OptionalFloat: proto.Float32(float32(math.NaN()))},
			y:    &testeditionspb.TestAllTypes{OptionalFloat: proto.Float32(0)},
			want: -1,
		},
		{
			x:    &testeditionspb.TestAllTypes{OptionalFloat: proto.Float32(0)},
			y:    &testeditionspb.TestAllTypes{OptionalFloat: proto.Float32(float32(math.NaN()))},
			want: 1,
		},
		{
			x:    &testeditionspb.TestAllTypes{OptionalDouble: proto.Float64(float64(math.NaN()))},
			y:    &testeditionspb.TestAllTypes{OptionalDouble: proto.Float64(0)},
			want: -1,
		},
		{
			x:    &testeditionspb.TestAllTypes{OptionalDouble: proto.Float64(0)},
			y:    &testeditionspb.TestAllTypes{OptionalDouble: proto.Float64(float64(math.NaN()))},
			want: 1,
		},
		{
			x:    &testeditionspb.TestAllTypes{OptionalBool: proto.Bool(false)},
			y:    &testeditionspb.TestAllTypes{OptionalBool: proto.Bool(true)},
			want: -1,
		},
		{
			x:    &testeditionspb.TestAllTypes{OptionalBool: proto.Bool(true)},
			y:    &testeditionspb.TestAllTypes{OptionalBool: proto.Bool(false)},
			want: 1,
		},
		{
			x:    &testeditionspb.TestAllTypes{OptionalString: proto.String("a")},
			y:    &testeditionspb.TestAllTypes{OptionalString: proto.String("b")},
			want: -1,
		},
		{
			x:    &testeditionspb.TestAllTypes{OptionalString: proto.String("b")},
			y:    &testeditionspb.TestAllTypes{OptionalString: proto.String("a")},
			want: 1,
		},
		{
			x:    &testeditionspb.TestAllTypes{OptionalBytes: []byte("a")},
			y:    &testeditionspb.TestAllTypes{OptionalBytes: []byte("b")},
			want: -1,
		},
		{
			x:    &testeditionspb.TestAllTypes{OptionalBytes: []byte("b")},
			y:    &testeditionspb.TestAllTypes{OptionalBytes: []byte("a")},
			want: 1,
		},
		{
			x:    &testeditionspb.TestAllTypes{OptionalNestedEnum: testeditionspb.TestAllTypes_FOO.Enum()},
			y:    &testeditionspb.TestAllTypes{OptionalNestedEnum: testeditionspb.TestAllTypes_BAR.Enum()},
			want: -1,
		},
		{
			x:    &testeditionspb.TestAllTypes{OptionalNestedEnum: testeditionspb.TestAllTypes_BAR.Enum()},
			y:    &testeditionspb.TestAllTypes{OptionalNestedEnum: testeditionspb.TestAllTypes_FOO.Enum()},
			want: 1,
		},
		{
			x:    &testeditionspb.TestAllTypes{OptionalInt32: proto.Int32(2)},
			y:    &testeditionspb.TestAllTypes{OptionalInt32: proto.Int32(2)},
			want: 0,
		},
		{
			x:    &testeditionspb.TestAllTypes{OptionalInt64: proto.Int64(2)},
			y:    &testeditionspb.TestAllTypes{OptionalInt64: proto.Int64(2)},
			want: 0,
		},
		{
			x:    &testeditionspb.TestAllTypes{OptionalUint32: proto.Uint32(2)},
			y:    &testeditionspb.TestAllTypes{OptionalUint32: proto.Uint32(2)},
			want: 0,
		},
		{
			x:    &testeditionspb.TestAllTypes{OptionalUint64: proto.Uint64(2)},
			y:    &testeditionspb.TestAllTypes{OptionalUint64: proto.Uint64(2)},
			want: 0,
		},
		{
			x:    &testeditionspb.TestAllTypes{OptionalSint32: proto.Int32(2)},
			y:    &testeditionspb.TestAllTypes{OptionalSint32: proto.Int32(2)},
			want: 0,
		},
		{
			x:    &testeditionspb.TestAllTypes{OptionalSint64: proto.Int64(2)},
			y:    &testeditionspb.TestAllTypes{OptionalSint64: proto.Int64(2)},
			want: 0,
		},
		{
			x:    &testeditionspb.TestAllTypes{OptionalFixed32: proto.Uint32(2)},
			y:    &testeditionspb.TestAllTypes{OptionalFixed32: proto.Uint32(2)},
			want: 0,
		},
		{
			x:    &testeditionspb.TestAllTypes{OptionalFixed64: proto.Uint64(2)},
			y:    &testeditionspb.TestAllTypes{OptionalFixed64: proto.Uint64(2)},
			want: 0,
		},
		{
			x:    &testeditionspb.TestAllTypes{OptionalSfixed32: proto.Int32(2)},
			y:    &testeditionspb.TestAllTypes{OptionalSfixed32: proto.Int32(2)},
			want: 0,
		},
		{
			x:    &testeditionspb.TestAllTypes{OptionalSfixed64: proto.Int64(2)},
			y:    &testeditionspb.TestAllTypes{OptionalSfixed64: proto.Int64(2)},
			want: 0,
		},
		{
			x:    &testeditionspb.TestAllTypes{OptionalFloat: proto.Float32(2)},
			y:    &testeditionspb.TestAllTypes{OptionalFloat: proto.Float32(2)},
			want: 0,
		},
		{
			x:    &testeditionspb.TestAllTypes{OptionalDouble: proto.Float64(2)},
			y:    &testeditionspb.TestAllTypes{OptionalDouble: proto.Float64(2)},
			want: 0,
		},
		{
			x:    &testeditionspb.TestAllTypes{OptionalFloat: proto.Float32(float32(math.NaN()))},
			y:    &testeditionspb.TestAllTypes{OptionalFloat: proto.Float32(float32(math.NaN()))},
			want: 0,
		},
		{
			x:    &testeditionspb.TestAllTypes{OptionalDouble: proto.Float64(float64(math.NaN()))},
			y:    &testeditionspb.TestAllTypes{OptionalDouble: proto.Float64(float64(math.NaN()))},
			want: 0,
		},
		{
			x:    &testeditionspb.TestAllTypes{OptionalBool: proto.Bool(true)},
			y:    &testeditionspb.TestAllTypes{OptionalBool: proto.Bool(true)},
			want: 0,
		},
		{
			x:    &testeditionspb.TestAllTypes{OptionalString: proto.String("abc")},
			y:    &testeditionspb.TestAllTypes{OptionalString: proto.String("abc")},
			want: 0,
		},
		{
			x:    &testeditionspb.TestAllTypes{OptionalBytes: []byte("abc")},
			y:    &testeditionspb.TestAllTypes{OptionalBytes: []byte("abc")},
			want: 0,
		},
		{
			x:    &testeditionspb.TestAllTypes{OptionalNestedEnum: testeditionspb.TestAllTypes_FOO.Enum()},
			y:    &testeditionspb.TestAllTypes{OptionalNestedEnum: testeditionspb.TestAllTypes_FOO.Enum()},
			want: 0,
		},

		// Proto2 presence.
		{
			x:    &testpb.TestAllTypes{},
			y:    &testpb.TestAllTypes{OptionalInt32: proto.Int32(0)},
			want: -1,
		},
		{
			x:    &testpb.TestAllTypes{OptionalInt32: proto.Int32(0)},
			y:    &testpb.TestAllTypes{},
			want: 1,
		},
		{
			x:    &testpb.TestAllTypes{},
			y:    &testpb.TestAllTypes{OptionalInt64: proto.Int64(0)},
			want: -1,
		},
		{
			x:    &testpb.TestAllTypes{OptionalInt64: proto.Int64(0)},
			y:    &testpb.TestAllTypes{},
			want: 1,
		},
		{
			x:    &testpb.TestAllTypes{},
			y:    &testpb.TestAllTypes{OptionalUint32: proto.Uint32(0)},
			want: -1,
		},
		{
			x:    &testpb.TestAllTypes{OptionalUint32: proto.Uint32(0)},
			y:    &testpb.TestAllTypes{},
			want: 1,
		},
		{
			x:    &testpb.TestAllTypes{},
			y:    &testpb.TestAllTypes{OptionalUint64: proto.Uint64(0)},
			want: -1,
		},
		{
			x:    &testpb.TestAllTypes{OptionalUint64: proto.Uint64(0)},
			y:    &testpb.TestAllTypes{},
			want: 1,
		},
		{
			x:    &testpb.TestAllTypes{},
			y:    &testpb.TestAllTypes{OptionalSint32: proto.Int32(0)},
			want: -1,
		},
		{
			x:    &testpb.TestAllTypes{OptionalSint32: proto.Int32(0)},
			y:    &testpb.TestAllTypes{},
			want: 1,
		},
		{
			x:    &testpb.TestAllTypes{},
			y:    &testpb.TestAllTypes{OptionalSint64: proto.Int64(0)},
			want: -1,
		},
		{
			x:    &testpb.TestAllTypes{OptionalSint64: proto.Int64(0)},
			y:    &testpb.TestAllTypes{},
			want: 1,
		},
		{
			x:    &testpb.TestAllTypes{},
			y:    &testpb.TestAllTypes{OptionalFixed32: proto.Uint32(0)},
			want: -1,
		},
		{
			x:    &testpb.TestAllTypes{OptionalFixed32: proto.Uint32(0)},
			y:    &testpb.TestAllTypes{},
			want: 1,
		},
		{
			x:    &testpb.TestAllTypes{},
			y:    &testpb.TestAllTypes{OptionalFixed64: proto.Uint64(0)},
			want: -1,
		},
		{
			x:    &testpb.TestAllTypes{OptionalFixed64: proto.Uint64(0)},
			y:    &testpb.TestAllTypes{},
			want: 1,
		},
		{
			x:    &testpb.TestAllTypes{},
			y:    &testpb.TestAllTypes{OptionalSfixed32: proto.Int32(0)},
			want: -1,
		},
		{
			x:    &testpb.TestAllTypes{OptionalSfixed32: proto.Int32(0)},
			y:    &testpb.TestAllTypes{},
			want: 1,
		},
		{
			x:    &testpb.TestAllTypes{},
			y:    &testpb.TestAllTypes{OptionalSfixed64: proto.Int64(0)},
			want: -1,
		},
		{
			x:    &testpb.TestAllTypes{OptionalSfixed64: proto.Int64(0)},
			y:    &testpb.TestAllTypes{},
			want: 1,
		},
		{
			x:    &testpb.TestAllTypes{},
			y:    &testpb.TestAllTypes{OptionalFloat: proto.Float32(0)},
			want: -1,
		},
		{
			x:    &testpb.TestAllTypes{OptionalFloat: proto.Float32(0)},
			y:    &testpb.TestAllTypes{},
			want: 1,
		},
		{
			x:    &testpb.TestAllTypes{},
			y:    &testpb.TestAllTypes{OptionalDouble: proto.Float64(0)},
			want: -1,
		},
		{
			x:    &testpb.TestAllTypes{OptionalDouble: proto.Float64(0)},
			y:    &testpb.TestAllTypes{},
			want: 1,
		},
		{
			x:    &testpb.TestAllTypes{},
			y:    &testpb.TestAllTypes{OptionalBool: proto.Bool(false)},
			want: -1,
		},
		{
			x:    &testpb.TestAllTypes{OptionalBool: proto.Bool(false)},
			y:    &testpb.TestAllTypes{},
			want: 1,
		},
		{
			x:    &testpb.TestAllTypes{},
			y:    &testpb.TestAllTypes{OptionalString: proto.String("")},
			want: -1,
		},
		{
			x:    &testpb.TestAllTypes{OptionalString: proto.String("")},
			y:    &testpb.TestAllTypes{},
			want: 1,
		},
		{
			x:    &testpb.TestAllTypes{},
			y:    &testpb.TestAllTypes{OptionalBytes: []byte{}},
			want: -1,
		},
		{
			x:    &testpb.TestAllTypes{OptionalBytes: []byte{}},
			y:    &testpb.TestAllTypes{},
			want: 1,
		},
		{
			x:    &testpb.TestAllTypes{},
			y:    &testpb.TestAllTypes{OptionalNestedEnum: testpb.TestAllTypes_FOO.Enum()},
			want: -1,
		},
		{
			x:    &testpb.TestAllTypes{OptionalNestedEnum: testpb.TestAllTypes_FOO.Enum()},
			y:    &testpb.TestAllTypes{},
			want: 1,
		},

		// Editions presence.
		{
			x:    &testeditionspb.TestAllTypes{},
			y:    &testeditionspb.TestAllTypes{OptionalInt32: proto.Int32(0)},
			want: -1,
		},
		{
			x:    &testeditionspb.TestAllTypes{OptionalInt32: proto.Int32(0)},
			y:    &testeditionspb.TestAllTypes{},
			want: 1,
		},
		{
			x:    &testeditionspb.TestAllTypes{},
			y:    &testeditionspb.TestAllTypes{OptionalInt64: proto.Int64(0)},
			want: -1,
		},
		{
			x:    &testeditionspb.TestAllTypes{OptionalInt64: proto.Int64(0)},
			y:    &testeditionspb.TestAllTypes{},
			want: 1,
		},
		{
			x:    &testeditionspb.TestAllTypes{},
			y:    &testeditionspb.TestAllTypes{OptionalUint32: proto.Uint32(0)},
			want: -1,
		},
		{
			x:    &testeditionspb.TestAllTypes{OptionalUint32: proto.Uint32(0)},
			y:    &testeditionspb.TestAllTypes{},
			want: 1,
		},
		{
			x:    &testeditionspb.TestAllTypes{},
			y:    &testeditionspb.TestAllTypes{OptionalUint64: proto.Uint64(0)},
			want: -1,
		},
		{
			x:    &testeditionspb.TestAllTypes{OptionalUint64: proto.Uint64(0)},
			y:    &testeditionspb.TestAllTypes{},
			want: 1,
		},
		{
			x:    &testeditionspb.TestAllTypes{},
			y:    &testeditionspb.TestAllTypes{OptionalSint32: proto.Int32(0)},
			want: -1,
		},
		{
			x:    &testeditionspb.TestAllTypes{OptionalSint32: proto.Int32(0)},
			y:    &testeditionspb.TestAllTypes{},
			want: 1,
		},
		{
			x:    &testeditionspb.TestAllTypes{},
			y:    &testeditionspb.TestAllTypes{OptionalSint64: proto.Int64(0)},
			want: -1,
		},
		{
			x:    &testeditionspb.TestAllTypes{OptionalSint64: proto.Int64(0)},
			y:    &testeditionspb.TestAllTypes{},
			want: 1,
		},
		{
			x:    &testeditionspb.TestAllTypes{},
			y:    &testeditionspb.TestAllTypes{OptionalFixed32: proto.Uint32(0)},
			want: -1,
		},
		{
			x:    &testeditionspb.TestAllTypes{OptionalFixed32: proto.Uint32(0)},
			y:    &testeditionspb.TestAllTypes{},
			want: 1,
		},
		{
			x:    &testeditionspb.TestAllTypes{},
			y:    &testeditionspb.TestAllTypes{OptionalFixed64: proto.Uint64(0)},
			want: -1,
		},
		{
			x:    &testeditionspb.TestAllTypes{OptionalFixed64: proto.Uint64(0)},
			y:    &testeditionspb.TestAllTypes{},
			want: 1,
		},
		{
			x:    &testeditionspb.TestAllTypes{},
			y:    &testeditionspb.TestAllTypes{OptionalSfixed32: proto.Int32(0)},
			want: -1,
		},
		{
			x:    &testeditionspb.TestAllTypes{OptionalSfixed32: proto.Int32(0)},
			y:    &testeditionspb.TestAllTypes{},
			want: 1,
		},
		{
			x:    &testeditionspb.TestAllTypes{},
			y:    &testeditionspb.TestAllTypes{OptionalSfixed64: proto.Int64(0)},
			want: -1,
		},
		{
			x:    &testeditionspb.TestAllTypes{OptionalSfixed64: proto.Int64(0)},
			y:    &testeditionspb.TestAllTypes{},
			want: 1,
		},
		{
			x:    &testeditionspb.TestAllTypes{},
			y:    &testeditionspb.TestAllTypes{OptionalFloat: proto.Float32(0)},
			want: -1,
		},
		{
			x:    &testeditionspb.TestAllTypes{OptionalFloat: proto.Float32(0)},
			y:    &testeditionspb.TestAllTypes{},
			want: 1,
		},
		{
			x:    &testeditionspb.TestAllTypes{},
			y:    &testeditionspb.TestAllTypes{OptionalDouble: proto.Float64(0)},
			want: -1,
		},
		{
			x:    &testeditionspb.TestAllTypes{OptionalDouble: proto.Float64(0)},
			y:    &testeditionspb.TestAllTypes{},
			want: 1,
		},
		{
			x:    &testeditionspb.TestAllTypes{},
			y:    &testeditionspb.TestAllTypes{OptionalBool: proto.Bool(false)},
			want: -1,
		},
		{
			x:    &testeditionspb.TestAllTypes{OptionalBool: proto.Bool(false)},
			y:    &testeditionspb.TestAllTypes{},
			want: 1,
		},
		{
			x:    &testeditionspb.TestAllTypes{},
			y:    &testeditionspb.TestAllTypes{OptionalString: proto.String("")},
			want: -1,
		},
		{
			x:    &testeditionspb.TestAllTypes{OptionalString: proto.String("")},
			y:    &testeditionspb.TestAllTypes{},
			want: 1,
		},
		{
			x:    &testeditionspb.TestAllTypes{},
			y:    &testeditionspb.TestAllTypes{OptionalBytes: []byte{}},
			want: -1,
		},
		{
			x:    &testeditionspb.TestAllTypes{OptionalBytes: []byte{}},
			y:    &testeditionspb.TestAllTypes{},
			want: 1,
		},
		{
			x:    &testeditionspb.TestAllTypes{},
			y:    &testeditionspb.TestAllTypes{OptionalNestedEnum: testeditionspb.TestAllTypes_FOO.Enum()},
			want: -1,
		},
		{
			x:    &testeditionspb.TestAllTypes{OptionalNestedEnum: testeditionspb.TestAllTypes_FOO.Enum()},
			y:    &testeditionspb.TestAllTypes{},
			want: 1,
		},

		// Proto3 presence.
		{
			x:    &test3pb.TestAllTypes{},
			y:    &test3pb.TestAllTypes{OptionalInt32: proto.Int32(0)},
			want: -1,
		},
		{
			x:    &test3pb.TestAllTypes{OptionalInt32: proto.Int32(0)},
			y:    &test3pb.TestAllTypes{},
			want: 1,
		},
		{
			x:    &test3pb.TestAllTypes{},
			y:    &test3pb.TestAllTypes{OptionalInt64: proto.Int64(0)},
			want: -1,
		},
		{
			x:    &test3pb.TestAllTypes{OptionalInt64: proto.Int64(0)},
			y:    &test3pb.TestAllTypes{},
			want: 1,
		},
		{
			x:    &test3pb.TestAllTypes{},
			y:    &test3pb.TestAllTypes{OptionalUint32: proto.Uint32(0)},
			want: -1,
		},
		{
			x:    &test3pb.TestAllTypes{OptionalUint32: proto.Uint32(0)},
			y:    &test3pb.TestAllTypes{},
			want: 1,
		},
		{
			x:    &test3pb.TestAllTypes{},
			y:    &test3pb.TestAllTypes{OptionalUint64: proto.Uint64(0)},
			want: -1,
		},
		{
			x:    &test3pb.TestAllTypes{OptionalUint64: proto.Uint64(0)},
			y:    &test3pb.TestAllTypes{},
			want: 1,
		},
		{
			x:    &test3pb.TestAllTypes{},
			y:    &test3pb.TestAllTypes{OptionalSint32: proto.Int32(0)},
			want: -1,
		},
		{
			x:    &test3pb.TestAllTypes{OptionalSint32: proto.Int32(0)},
			y:    &test3pb.TestAllTypes{},
			want: 1,
		},
		{
			x:    &test3pb.TestAllTypes{},
			y:    &test3pb.TestAllTypes{OptionalSint64: proto.Int64(0)},
			want: -1,
		},
		{
			x:    &test3pb.TestAllTypes{OptionalSint64: proto.Int64(0)},
			y:    &test3pb.TestAllTypes{},
			want: 1,
		},
		{
			x:    &test3pb.TestAllTypes{},
			y:    &test3pb.TestAllTypes{OptionalFixed32: proto.Uint32(0)},
			want: -1,
		},
		{
			x:    &test3pb.TestAllTypes{OptionalFixed32: proto.Uint32(0)},
			y:    &test3pb.TestAllTypes{},
			want: 1,
		},
		{
			x:    &test3pb.TestAllTypes{},
			y:    &test3pb.TestAllTypes{OptionalFixed64: proto.Uint64(0)},
			want: -1,
		},
		{
			x:    &test3pb.TestAllTypes{OptionalFixed64: proto.Uint64(0)},
			y:    &test3pb.TestAllTypes{},
			want: 1,
		},
		{
			x:    &test3pb.TestAllTypes{},
			y:    &test3pb.TestAllTypes{OptionalSfixed32: proto.Int32(0)},
			want: -1,
		},
		{
			x:    &test3pb.TestAllTypes{OptionalSfixed32: proto.Int32(0)},
			y:    &test3pb.TestAllTypes{},
			want: 1,
		},
		{
			x:    &test3pb.TestAllTypes{},
			y:    &test3pb.TestAllTypes{OptionalSfixed64: proto.Int64(0)},
			want: -1,
		},
		{
			x:    &test3pb.TestAllTypes{OptionalSfixed64: proto.Int64(0)},
			y:    &test3pb.TestAllTypes{},
			want: 1,
		},
		{
			x:    &test3pb.TestAllTypes{},
			y:    &test3pb.TestAllTypes{OptionalFloat: proto.Float32(0)},
			want: -1,
		},
		{
			x:    &test3pb.TestAllTypes{OptionalFloat: proto.Float32(0)},
			y:    &test3pb.TestAllTypes{},
			want: 1,
		},
		{
			x:    &test3pb.TestAllTypes{},
			y:    &test3pb.TestAllTypes{OptionalDouble: proto.Float64(0)},
			want: -1,
		},
		{
			x:    &test3pb.TestAllTypes{OptionalDouble: proto.Float64(0)},
			y:    &test3pb.TestAllTypes{},
			want: 1,
		},
		{
			x:    &test3pb.TestAllTypes{},
			y:    &test3pb.TestAllTypes{OptionalBool: proto.Bool(false)},
			want: -1,
		},
		{
			x:    &test3pb.TestAllTypes{OptionalBool: proto.Bool(false)},
			y:    &test3pb.TestAllTypes{},
			want: 1,
		},
		{
			x:    &test3pb.TestAllTypes{},
			y:    &test3pb.TestAllTypes{OptionalString: proto.String("")},
			want: -1,
		},
		{
			x:    &test3pb.TestAllTypes{OptionalString: proto.String("")},
			y:    &test3pb.TestAllTypes{},
			want: 1,
		},
		{
			x:    &test3pb.TestAllTypes{},
			y:    &test3pb.TestAllTypes{OptionalBytes: []byte{}},
			want: -1,
		},
		{
			x:    &test3pb.TestAllTypes{OptionalBytes: []byte{}},
			y:    &test3pb.TestAllTypes{},
			want: 1,
		},
		{
			x:    &test3pb.TestAllTypes{},
			y:    &test3pb.TestAllTypes{OptionalNestedEnum: test3pb.TestAllTypes_FOO.Enum()},
			want: -1,
		},
		{
			x:    &test3pb.TestAllTypes{OptionalNestedEnum: test3pb.TestAllTypes_FOO.Enum()},
			y:    &test3pb.TestAllTypes{},
			want: 1,
		},

		// Proto2 default values are not considered by Equal, so the following are still unequal.
		{
			x:    &testpb.TestAllTypes{},
			y:    &testpb.TestAllTypes{DefaultInt32: proto.Int32(81)},
			want: -1,
		},
		{
			x:    &testpb.TestAllTypes{DefaultInt32: proto.Int32(81)},
			y:    &testpb.TestAllTypes{},
			want: 1,
		},
		{
			x:    &testpb.TestAllTypes{},
			y:    &testpb.TestAllTypes{DefaultInt64: proto.Int64(82)},
			want: -1,
		},
		{
			x:    &testpb.TestAllTypes{DefaultInt64: proto.Int64(82)},
			y:    &testpb.TestAllTypes{},
			want: 1,
		},
		{
			x:    &testpb.TestAllTypes{},
			y:    &testpb.TestAllTypes{DefaultUint32: proto.Uint32(83)},
			want: -1,
		},
		{
			x:    &testpb.TestAllTypes{DefaultUint32: proto.Uint32(83)},
			y:    &testpb.TestAllTypes{},
			want: 1,
		},
		{
			x:    &testpb.TestAllTypes{},
			y:    &testpb.TestAllTypes{DefaultUint64: proto.Uint64(84)},
			want: -1,
		},
		{
			x:    &testpb.TestAllTypes{DefaultUint64: proto.Uint64(84)},
			y:    &testpb.TestAllTypes{},
			want: 1,
		},
		{
			x:    &testpb.TestAllTypes{},
			y:    &testpb.TestAllTypes{DefaultSint32: proto.Int32(-85)},
			want: -1,
		},
		{
			x:    &testpb.TestAllTypes{DefaultSint32: proto.Int32(-85)},
			y:    &testpb.TestAllTypes{},
			want: 1,
		},
		{
			x:    &testpb.TestAllTypes{},
			y:    &testpb.TestAllTypes{DefaultSint64: proto.Int64(86)},
			want: -1,
		},
		{
			x:    &testpb.TestAllTypes{DefaultSint64: proto.Int64(86)},
			y:    &testpb.TestAllTypes{},
			want: 1,
		},
		{
			x:    &testpb.TestAllTypes{},
			y:    &testpb.TestAllTypes{DefaultFixed32: proto.Uint32(87)},
			want: -1,
		},
		{
			x:    &testpb.TestAllTypes{DefaultFixed32: proto.Uint32(87)},
			y:    &testpb.TestAllTypes{},
			want: 1,
		},
		{
			x:    &testpb.TestAllTypes{},
			y:    &testpb.TestAllTypes{DefaultFixed64: proto.Uint64(88)},
			want: -1,
		},
		{
			x:    &testpb.TestAllTypes{DefaultFixed64: proto.Uint64(88)},
			y:    &testpb.TestAllTypes{},
			want: 1,
		},
		{
			x:    &testpb.TestAllTypes{},
			y:    &testpb.TestAllTypes{DefaultSfixed32: proto.Int32(89)},
			want: -1,
		},
		{
			x:    &testpb.TestAllTypes{DefaultSfixed32: proto.Int32(89)},
			y:    &testpb.TestAllTypes{},
			want: 1,
		},
		{
			x:    &testpb.TestAllTypes{},
			y:    &testpb.TestAllTypes{DefaultSfixed64: proto.Int64(-90)},
			want: -1,
		},
		{
			x:    &testpb.TestAllTypes{DefaultSfixed64: proto.Int64(-90)},
			y:    &testpb.TestAllTypes{},
			want: 1,
		},
		{
			x:    &testpb.TestAllTypes{},
			y:    &testpb.TestAllTypes{DefaultFloat: proto.Float32(91.5)},
			want: -1,
		},
		{
			x:    &testpb.TestAllTypes{DefaultFloat: proto.Float32(91.5)},
			y:    &testpb.TestAllTypes{},
			want: 1,
		},
		{
			x:    &testpb.TestAllTypes{},
			y:    &testpb.TestAllTypes{DefaultDouble: proto.Float64(92e3)},
			want: -1,
		},
		{
			x:    &testpb.TestAllTypes{DefaultDouble: proto.Float64(92e3)},
			y:    &testpb.TestAllTypes{},
			want: 1,
		},
		{
			x:    &testpb.TestAllTypes{},
			y:    &testpb.TestAllTypes{DefaultBool: proto.Bool(true)},
			want: -1,
		},
		{
			x:    &testpb.TestAllTypes{DefaultBool: proto.Bool(true)},
			y:    &testpb.TestAllTypes{},
			want: 1,
		},
		{
			x:    &testpb.TestAllTypes{},
			y:    &testpb.TestAllTypes{DefaultString: proto.String("hello")},
			want: -1,
		},
		{
			x:    &testpb.TestAllTypes{DefaultString: proto.String("hello")},
			y:    &testpb.TestAllTypes{},
			want: 1,
		},
		{
			x:    &testpb.TestAllTypes{},
			y:    &testpb.TestAllTypes{DefaultBytes: []byte("world")},
			want: -1,
		},
		{
			x:    &testpb.TestAllTypes{DefaultBytes: []byte("world")},
			y:    &testpb.TestAllTypes{},
			want: 1,
		},
		{
			x:    &testpb.TestAllTypes{},
			y:    &testpb.TestAllTypes{DefaultNestedEnum: testpb.TestAllTypes_BAR.Enum()},
			want: -1,
		},
		{
			x:    &testpb.TestAllTypes{DefaultNestedEnum: testpb.TestAllTypes_BAR.Enum()},
			y:    &testpb.TestAllTypes{},
			want: 1,
		},

		// Edition default values are not considered by Equal, so the following are still unequal.
		{
			x:    &testeditionspb.TestAllTypes{},
			y:    &testeditionspb.TestAllTypes{DefaultInt32: proto.Int32(81)},
			want: -1,
		},
		{
			x:    &testeditionspb.TestAllTypes{DefaultInt32: proto.Int32(81)},
			y:    &testeditionspb.TestAllTypes{},
			want: 1,
		},
		{
			x:    &testeditionspb.TestAllTypes{},
			y:    &testeditionspb.TestAllTypes{DefaultInt64: proto.Int64(82)},
			want: -1,
		},
		{
			x:    &testeditionspb.TestAllTypes{DefaultInt64: proto.Int64(82)},
			y:    &testeditionspb.TestAllTypes{},
			want: 1,
		},
		{
			x:    &testeditionspb.TestAllTypes{},
			y:    &testeditionspb.TestAllTypes{DefaultUint32: proto.Uint32(83)},
			want: -1,
		},
		{
			x:    &testeditionspb.TestAllTypes{DefaultUint32: proto.Uint32(83)},
			y:    &testeditionspb.TestAllTypes{},
			want: 1,
		},
		{
			x:    &testeditionspb.TestAllTypes{},
			y:    &testeditionspb.TestAllTypes{DefaultUint64: proto.Uint64(84)},
			want: -1,
		},
		{
			x:    &testeditionspb.TestAllTypes{DefaultUint64: proto.Uint64(84)},
			y:    &testeditionspb.TestAllTypes{},
			want: 1,
		},
		{
			x:    &testeditionspb.TestAllTypes{},
			y:    &testeditionspb.TestAllTypes{DefaultSint32: proto.Int32(-85)},
			want: -1,
		},
		{
			x:    &testeditionspb.TestAllTypes{DefaultSint32: proto.Int32(-85)},
			y:    &testeditionspb.TestAllTypes{},
			want: 1,
		},
		{
			x:    &testeditionspb.TestAllTypes{},
			y:    &testeditionspb.TestAllTypes{DefaultSint64: proto.Int64(86)},
			want: -1,
		},
		{
			x:    &testeditionspb.TestAllTypes{DefaultSint64: proto.Int64(86)},
			y:    &testeditionspb.TestAllTypes{},
			want: 1,
		},
		{
			x:    &testeditionspb.TestAllTypes{},
			y:    &testeditionspb.TestAllTypes{DefaultFixed32: proto.Uint32(87)},
			want: -1,
		},
		{
			x:    &testeditionspb.TestAllTypes{DefaultFixed32: proto.Uint32(87)},
			y:    &testeditionspb.TestAllTypes{},
			want: 1,
		},
		{
			x:    &testeditionspb.TestAllTypes{},
			y:    &testeditionspb.TestAllTypes{DefaultFixed64: proto.Uint64(88)},
			want: -1,
		},
		{
			x:    &testeditionspb.TestAllTypes{DefaultFixed64: proto.Uint64(88)},
			y:    &testeditionspb.TestAllTypes{},
			want: 1,
		},
		{
			x:    &testeditionspb.TestAllTypes{},
			y:    &testeditionspb.TestAllTypes{DefaultSfixed32: proto.Int32(89)},
			want: -1,
		},
		{
			x:    &testeditionspb.TestAllTypes{DefaultSfixed32: proto.Int32(89)},
			y:    &testeditionspb.TestAllTypes{},
			want: 1,
		},
		{
			x:    &testeditionspb.TestAllTypes{},
			y:    &testeditionspb.TestAllTypes{DefaultSfixed64: proto.Int64(-90)},
			want: -1,
		},
		{
			x:    &testeditionspb.TestAllTypes{DefaultSfixed64: proto.Int64(-90)},
			y:    &testeditionspb.TestAllTypes{},
			want: 1,
		},
		{
			x:    &testeditionspb.TestAllTypes{},
			y:    &testeditionspb.TestAllTypes{DefaultFloat: proto.Float32(91.5)},
			want: -1,
		},
		{
			x:    &testeditionspb.TestAllTypes{DefaultFloat: proto.Float32(91.5)},
			y:    &testeditionspb.TestAllTypes{},
			want: 1,
		},
		{
			x:    &testeditionspb.TestAllTypes{},
			y:    &testeditionspb.TestAllTypes{DefaultDouble: proto.Float64(92e3)},
			want: -1,
		},
		{
			x:    &testeditionspb.TestAllTypes{DefaultDouble: proto.Float64(92e3)},
			y:    &testeditionspb.TestAllTypes{},
			want: 1,
		},
		{
			x:    &testeditionspb.TestAllTypes{},
			y:    &testeditionspb.TestAllTypes{DefaultBool: proto.Bool(true)},
			want: -1,
		},
		{
			x:    &testeditionspb.TestAllTypes{DefaultBool: proto.Bool(true)},
			y:    &testeditionspb.TestAllTypes{},
			want: 1,
		},
		{
			x:    &testeditionspb.TestAllTypes{},
			y:    &testeditionspb.TestAllTypes{DefaultString: proto.String("hello")},
			want: -1,
		},
		{
			x:    &testeditionspb.TestAllTypes{DefaultString: proto.String("hello")},
			y:    &testeditionspb.TestAllTypes{},
			want: 1,
		},
		{
			x:    &testeditionspb.TestAllTypes{},
			y:    &testeditionspb.TestAllTypes{DefaultBytes: []byte("world")},
			want: -1,
		},
		{
			x:    &testeditionspb.TestAllTypes{DefaultBytes: []byte("world")},
			y:    &testeditionspb.TestAllTypes{},
			want: 1,
		},
		{
			x:    &testeditionspb.TestAllTypes{},
			y:    &testeditionspb.TestAllTypes{DefaultNestedEnum: testeditionspb.TestAllTypes_BAR.Enum()},
			want: -1,
		},
		{
			x:    &testeditionspb.TestAllTypes{DefaultNestedEnum: testeditionspb.TestAllTypes_BAR.Enum()},
			y:    &testeditionspb.TestAllTypes{},
			want: 1,
		},

		// Groups.
		{
			x: &testpb.TestAllTypes{Optionalgroup: &testpb.TestAllTypes_OptionalGroup{
				A: proto.Int32(1),
			}},
			y: &testpb.TestAllTypes{Optionalgroup: &testpb.TestAllTypes_OptionalGroup{
				A: proto.Int32(2),
			}},
			want: -1,
		},
		{
			x: &testpb.TestAllTypes{Optionalgroup: &testpb.TestAllTypes_OptionalGroup{
				A: proto.Int32(2),
			}},
			y: &testpb.TestAllTypes{Optionalgroup: &testpb.TestAllTypes_OptionalGroup{
				A: proto.Int32(1),
			}},
			want: 1,
		},
		{
			x:    &testpb.TestAllTypes{},
			y:    &testpb.TestAllTypes{Optionalgroup: &testpb.TestAllTypes_OptionalGroup{}},
			want: -1,
		},
		{
			x:    &testpb.TestAllTypes{Optionalgroup: &testpb.TestAllTypes_OptionalGroup{}},
			y:    &testpb.TestAllTypes{},
			want: 1,
		},
		{
			x: &testeditionspb.TestAllTypes{Optionalgroup: &testeditionspb.TestAllTypes_OptionalGroup{
				A: proto.Int32(1),
			}},
			y: &testeditionspb.TestAllTypes{Optionalgroup: &testeditionspb.TestAllTypes_OptionalGroup{
				A: proto.Int32(2),
			}},
			want: -1,
		},
		{
			x: &testeditionspb.TestAllTypes{Optionalgroup: &testeditionspb.TestAllTypes_OptionalGroup{
				A: proto.Int32(2),
			}},
			y: &testeditionspb.TestAllTypes{Optionalgroup: &testeditionspb.TestAllTypes_OptionalGroup{
				A: proto.Int32(1),
			}},
			want: 1,
		},
		{
			x:    &testeditionspb.TestAllTypes{},
			y:    &testeditionspb.TestAllTypes{Optionalgroup: &testeditionspb.TestAllTypes_OptionalGroup{}},
			want: -1,
		},
		{
			x:    &testeditionspb.TestAllTypes{Optionalgroup: &testeditionspb.TestAllTypes_OptionalGroup{}},
			y:    &testeditionspb.TestAllTypes{},
			want: 1,
		},

		// Messages.
		{
			x: &testpb.TestAllTypes{OptionalNestedMessage: &testpb.TestAllTypes_NestedMessage{
				A: proto.Int32(1),
			}},
			y: &testpb.TestAllTypes{OptionalNestedMessage: &testpb.TestAllTypes_NestedMessage{
				A: proto.Int32(2),
			}},
			want: -1,
		},
		{
			x: &testpb.TestAllTypes{OptionalNestedMessage: &testpb.TestAllTypes_NestedMessage{
				A: proto.Int32(2),
			}},
			y: &testpb.TestAllTypes{OptionalNestedMessage: &testpb.TestAllTypes_NestedMessage{
				A: proto.Int32(1),
			}},
			want: 1,
		},
		{
			x:    &testpb.TestAllTypes{},
			y:    &testpb.TestAllTypes{OptionalNestedMessage: &testpb.TestAllTypes_NestedMessage{}},
			want: -1,
		},
		{
			x:    &testpb.TestAllTypes{OptionalNestedMessage: &testpb.TestAllTypes_NestedMessage{}},
			y:    &testpb.TestAllTypes{},
			want: 1,
		},
		{
			x:    &testeditionspb.TestAllTypes{},
			y:    &testeditionspb.TestAllTypes{OptionalNestedMessage: &testeditionspb.TestAllTypes_NestedMessage{}},
			want: -1,
		},
		{
			x:    &testeditionspb.TestAllTypes{OptionalNestedMessage: &testeditionspb.TestAllTypes_NestedMessage{}},
			y:    &testeditionspb.TestAllTypes{},
			want: 1,
		},
		{
			x:    &test3pb.TestAllTypes{},
			y:    &test3pb.TestAllTypes{OptionalNestedMessage: &test3pb.TestAllTypes_NestedMessage{}},
			want: -1,
		},
		{
			x:    &test3pb.TestAllTypes{OptionalNestedMessage: &test3pb.TestAllTypes_NestedMessage{}},
			y:    &test3pb.TestAllTypes{},
			want: 1,
		},

		// Lists.
		{
			x:    &testpb.TestAllTypes{RepeatedInt32: []int32{1}},
			y:    &testpb.TestAllTypes{RepeatedInt32: []int32{1}},
			want: 0,
		},
		{
			x:    &testpb.TestAllTypes{RepeatedInt32: []int32{1}},
			y:    &testpb.TestAllTypes{RepeatedInt32: []int32{1, 2}},
			want: -1,
		},
		{
			x:    &testpb.TestAllTypes{RepeatedInt32: []int32{1, 2}},
			y:    &testpb.TestAllTypes{RepeatedInt32: []int32{1}},
			want: 1,
		},
		{
			x:    &testpb.TestAllTypes{RepeatedInt32: []int32{7}},
			y:    &testpb.TestAllTypes{RepeatedInt32: []int32{1, 2}},
			want: -1,
		},
		{
			x:    &testpb.TestAllTypes{RepeatedInt32: []int32{1, 2}},
			y:    &testpb.TestAllTypes{RepeatedInt32: []int32{7}},
			want: 1,
		},
		{
			x:    &testpb.TestAllTypes{RepeatedInt32: []int32{1, 2, 3}},
			y:    &testpb.TestAllTypes{RepeatedInt32: []int32{1, 3, 2}},
			want: -1,
		},
		{
			x:    &testpb.TestAllTypes{RepeatedInt32: []int32{1, 3, 2}},
			y:    &testpb.TestAllTypes{RepeatedInt32: []int32{1, 2, 3}},
			want: 1,
		},
		{
			x:    &testpb.TestAllTypes{RepeatedInt64: []int64{1, 2, 3}},
			y:    &testpb.TestAllTypes{RepeatedInt64: []int64{1, 3, 2}},
			want: -1,
		},
		{
			x:    &testpb.TestAllTypes{RepeatedInt64: []int64{1, 3, 2}},
			y:    &testpb.TestAllTypes{RepeatedInt64: []int64{1, 2, 3}},
			want: 1,
		},
		{
			x:    &testpb.TestAllTypes{RepeatedUint32: []uint32{1, 2, 3}},
			y:    &testpb.TestAllTypes{RepeatedUint32: []uint32{1, 3, 2}},
			want: -1,
		},
		{
			x:    &testpb.TestAllTypes{RepeatedUint32: []uint32{1, 3, 2}},
			y:    &testpb.TestAllTypes{RepeatedUint32: []uint32{1, 2, 3}},
			want: 1,
		},
		{
			x:    &testpb.TestAllTypes{RepeatedUint64: []uint64{1, 2, 3}},
			y:    &testpb.TestAllTypes{RepeatedUint64: []uint64{1, 3, 2}},
			want: -1,
		},
		{
			x:    &testpb.TestAllTypes{RepeatedUint64: []uint64{1, 3, 2}},
			y:    &testpb.TestAllTypes{RepeatedUint64: []uint64{1, 2, 3}},
			want: 1,
		},
		{
			x:    &testpb.TestAllTypes{RepeatedSint32: []int32{1, 2, 3}},
			y:    &testpb.TestAllTypes{RepeatedSint32: []int32{1, 3, 2}},
			want: -1,
		},
		{
			x:    &testpb.TestAllTypes{RepeatedSint32: []int32{1, 3, 2}},
			y:    &testpb.TestAllTypes{RepeatedSint32: []int32{1, 2, 3}},
			want: 1,
		},
		{
			x:    &testpb.TestAllTypes{RepeatedSint64: []int64{1, 2, 3}},
			y:    &testpb.TestAllTypes{RepeatedSint64: []int64{1, 3, 2}},
			want: -1,
		},
		{
			x:    &testpb.TestAllTypes{RepeatedSint64: []int64{1, 3, 2}},
			y:    &testpb.TestAllTypes{RepeatedSint64: []int64{1, 2, 3}},
			want: 1,
		},
		{
			x:    &testpb.TestAllTypes{RepeatedFixed32: []uint32{1, 2, 3}},
			y:    &testpb.TestAllTypes{RepeatedFixed32: []uint32{1, 3, 2}},
			want: -1,
		},
		{
			x:    &testpb.TestAllTypes{RepeatedFixed32: []uint32{1, 3, 2}},
			y:    &testpb.TestAllTypes{RepeatedFixed32: []uint32{1, 2, 3}},
			want: 1,
		},
		{
			x:    &testpb.TestAllTypes{RepeatedFixed64: []uint64{1, 2, 3}},
			y:    &testpb.TestAllTypes{RepeatedFixed64: []uint64{1, 3, 2}},
			want: -1,
		},
		{
			x:    &testpb.TestAllTypes{RepeatedFixed64: []uint64{1, 3, 2}},
			y:    &testpb.TestAllTypes{RepeatedFixed64: []uint64{1, 2, 3}},
			want: 1,
		},
		{
			x:    &testpb.TestAllTypes{RepeatedSfixed32: []int32{1, 2, 3}},
			y:    &testpb.TestAllTypes{RepeatedSfixed32: []int32{1, 3, 2}},
			want: -1,
		},
		{
			x:    &testpb.TestAllTypes{RepeatedSfixed32: []int32{1, 3, 2}},
			y:    &testpb.TestAllTypes{RepeatedSfixed32: []int32{1, 2, 3}},
			want: 1,
		},
		{
			x:    &testpb.TestAllTypes{RepeatedSfixed64: []int64{1, 2, 3}},
			y:    &testpb.TestAllTypes{RepeatedSfixed64: []int64{1, 3, 2}},
			want: -1,
		},
		{
			x:    &testpb.TestAllTypes{RepeatedSfixed64: []int64{1, 3, 2}},
			y:    &testpb.TestAllTypes{RepeatedSfixed64: []int64{1, 2, 3}},
			want: 1,
		},
		{
			x:    &testpb.TestAllTypes{RepeatedFloat: []float32{1, 2, 3}},
			y:    &testpb.TestAllTypes{RepeatedFloat: []float32{1, 3, 2}},
			want: -1,
		},
		{
			x:    &testpb.TestAllTypes{RepeatedFloat: []float32{1, 3, 2}},
			y:    &testpb.TestAllTypes{RepeatedFloat: []float32{1, 2, 3}},
			want: 1,
		},
		{
			x:    &testpb.TestAllTypes{RepeatedDouble: []float64{1, 2, 3}},
			y:    &testpb.TestAllTypes{RepeatedDouble: []float64{1, 3, 2}},
			want: -1,
		},
		{
			x:    &testpb.TestAllTypes{RepeatedDouble: []float64{1, 3, 2}},
			y:    &testpb.TestAllTypes{RepeatedDouble: []float64{1, 2, 3}},
			want: 1,
		},
		{
			x:    &testpb.TestAllTypes{RepeatedBool: []bool{true, false}},
			y:    &testpb.TestAllTypes{RepeatedBool: []bool{true, true}},
			want: -1,
		},
		{
			x:    &testpb.TestAllTypes{RepeatedBool: []bool{true, true}},
			y:    &testpb.TestAllTypes{RepeatedBool: []bool{true, false}},
			want: 1,
		},
		{
			x:    &testpb.TestAllTypes{RepeatedString: []string{"a", "b"}},
			y:    &testpb.TestAllTypes{RepeatedString: []string{"a", "c"}},
			want: -1,
		},
		{
			x:    &testpb.TestAllTypes{RepeatedString: []string{"a", "c"}},
			y:    &testpb.TestAllTypes{RepeatedString: []string{"a", "b"}},
			want: 1,
		},
		{
			x:    &testpb.TestAllTypes{RepeatedBytes: [][]byte{[]byte("a"), []byte("b")}},
			y:    &testpb.TestAllTypes{RepeatedBytes: [][]byte{[]byte("a"), []byte("c")}},
			want: -1,
		},
		{
			x:    &testpb.TestAllTypes{RepeatedBytes: [][]byte{[]byte("a"), []byte("c")}},
			y:    &testpb.TestAllTypes{RepeatedBytes: [][]byte{[]byte("a"), []byte("b")}},
			want: 1,
		},
		{
			x:    &testpb.TestAllTypes{RepeatedNestedEnum: []testpb.TestAllTypes_NestedEnum{testpb.TestAllTypes_FOO}},
			y:    &testpb.TestAllTypes{RepeatedNestedEnum: []testpb.TestAllTypes_NestedEnum{testpb.TestAllTypes_BAR}},
			want: -1,
		},
		{
			x:    &testpb.TestAllTypes{RepeatedNestedEnum: []testpb.TestAllTypes_NestedEnum{testpb.TestAllTypes_BAR}},
			y:    &testpb.TestAllTypes{RepeatedNestedEnum: []testpb.TestAllTypes_NestedEnum{testpb.TestAllTypes_FOO}},
			want: 1,
		},
		{
			x: &testpb.TestAllTypes{Repeatedgroup: []*testpb.TestAllTypes_RepeatedGroup{
				{
					A: proto.Int32(1),
				}, {
					A: proto.Int32(2),
				},
			}},
			y: &testpb.TestAllTypes{Repeatedgroup: []*testpb.TestAllTypes_RepeatedGroup{
				{
					A: proto.Int32(1),
				}, {
					A: proto.Int32(3),
				},
			}},
			want: -1,
		},
		{
			x: &testpb.TestAllTypes{Repeatedgroup: []*testpb.TestAllTypes_RepeatedGroup{
				{
					A: proto.Int32(1),
				}, {
					A: proto.Int32(3),
				},
			}},
			y: &testpb.TestAllTypes{Repeatedgroup: []*testpb.TestAllTypes_RepeatedGroup{
				{
					A: proto.Int32(1),
				}, {
					A: proto.Int32(2),
				},
			}},
			want: 1,
		},
		{
			x: &testpb.TestAllTypes{RepeatedNestedMessage: []*testpb.TestAllTypes_NestedMessage{
				{
					A: proto.Int32(1),
				}, {
					A: proto.Int32(2),
				},
			}},
			y: &testpb.TestAllTypes{RepeatedNestedMessage: []*testpb.TestAllTypes_NestedMessage{
				{
					A: proto.Int32(1),
				}, {
					A: proto.Int32(3),
				},
			}},
			want: -1,
		},
		{
			x: &testpb.TestAllTypes{RepeatedNestedMessage: []*testpb.TestAllTypes_NestedMessage{
				{
					A: proto.Int32(1),
				}, {
					A: proto.Int32(3),
				},
			}},
			y: &testpb.TestAllTypes{RepeatedNestedMessage: []*testpb.TestAllTypes_NestedMessage{
				{
					A: proto.Int32(1),
				}, {
					A: proto.Int32(2),
				},
			}},
			want: 1,
		},

		// Editions Lists.
		{
			x:    &testeditionspb.TestAllTypes{RepeatedInt32: []int32{1, 2, 3}},
			y:    &testeditionspb.TestAllTypes{RepeatedInt32: []int32{1, 2, 3}},
			want: 0,
		},
		{
			x:    &testeditionspb.TestAllTypes{RepeatedInt32: []int32{1}},
			y:    &testeditionspb.TestAllTypes{RepeatedInt32: []int32{1, 2, 3}},
			want: -1,
		},
		{
			x:    &testeditionspb.TestAllTypes{RepeatedInt32: []int32{1, 2, 3}},
			y:    &testeditionspb.TestAllTypes{RepeatedInt32: []int32{1}},
			want: 1,
		},
		{
			x:    &testeditionspb.TestAllTypes{RepeatedInt32: []int32{1, 2, 3}},
			y:    &testeditionspb.TestAllTypes{RepeatedInt32: []int32{1, 3, 2}},
			want: -1,
		},
		{
			x:    &testeditionspb.TestAllTypes{RepeatedInt32: []int32{1, 3, 2}},
			y:    &testeditionspb.TestAllTypes{RepeatedInt32: []int32{1, 2, 3}},
			want: 1,
		},
		{
			x:    &testeditionspb.TestAllTypes{RepeatedInt64: []int64{1, 2, 3}},
			y:    &testeditionspb.TestAllTypes{RepeatedInt64: []int64{1, 3, 2}},
			want: -1,
		},
		{
			x:    &testeditionspb.TestAllTypes{RepeatedInt64: []int64{1, 3, 2}},
			y:    &testeditionspb.TestAllTypes{RepeatedInt64: []int64{1, 2, 3}},
			want: 1,
		},
		{
			x:    &testeditionspb.TestAllTypes{RepeatedUint32: []uint32{1, 2, 3}},
			y:    &testeditionspb.TestAllTypes{RepeatedUint32: []uint32{1, 3, 2}},
			want: -1,
		},
		{
			x:    &testeditionspb.TestAllTypes{RepeatedUint32: []uint32{1, 3, 2}},
			y:    &testeditionspb.TestAllTypes{RepeatedUint32: []uint32{1, 2, 3}},
			want: 1,
		},
		{
			x:    &testeditionspb.TestAllTypes{RepeatedUint64: []uint64{1, 2, 3}},
			y:    &testeditionspb.TestAllTypes{RepeatedUint64: []uint64{1, 3, 2}},
			want: -1,
		},
		{
			x:    &testeditionspb.TestAllTypes{RepeatedUint64: []uint64{1, 3, 2}},
			y:    &testeditionspb.TestAllTypes{RepeatedUint64: []uint64{1, 2, 3}},
			want: 1,
		},
		{
			x:    &testeditionspb.TestAllTypes{RepeatedSint32: []int32{1, 2, 3}},
			y:    &testeditionspb.TestAllTypes{RepeatedSint32: []int32{1, 3, 2}},
			want: -1,
		},
		{
			x:    &testeditionspb.TestAllTypes{RepeatedSint32: []int32{1, 3, 2}},
			y:    &testeditionspb.TestAllTypes{RepeatedSint32: []int32{1, 2, 3}},
			want: 1,
		},
		{
			x:    &testeditionspb.TestAllTypes{RepeatedSint64: []int64{1, 2, 3}},
			y:    &testeditionspb.TestAllTypes{RepeatedSint64: []int64{1, 3, 2}},
			want: -1,
		},
		{
			x:    &testeditionspb.TestAllTypes{RepeatedSint64: []int64{1, 3, 2}},
			y:    &testeditionspb.TestAllTypes{RepeatedSint64: []int64{1, 2, 3}},
			want: 1,
		},
		{
			x:    &testeditionspb.TestAllTypes{RepeatedFixed32: []uint32{1, 2, 3}},
			y:    &testeditionspb.TestAllTypes{RepeatedFixed32: []uint32{1, 3, 2}},
			want: -1,
		},
		{
			x:    &testeditionspb.TestAllTypes{RepeatedFixed32: []uint32{1, 3, 2}},
			y:    &testeditionspb.TestAllTypes{RepeatedFixed32: []uint32{1, 2, 3}},
			want: 1,
		},
		{
			x:    &testeditionspb.TestAllTypes{RepeatedFixed64: []uint64{1, 2, 3}},
			y:    &testeditionspb.TestAllTypes{RepeatedFixed64: []uint64{1, 3, 2}},
			want: -1,
		},
		{
			x:    &testeditionspb.TestAllTypes{RepeatedFixed64: []uint64{1, 3, 2}},
			y:    &testeditionspb.TestAllTypes{RepeatedFixed64: []uint64{1, 2, 3}},
			want: 1,
		},
		{
			x:    &testeditionspb.TestAllTypes{RepeatedSfixed32: []int32{1, 2, 3}},
			y:    &testeditionspb.TestAllTypes{RepeatedSfixed32: []int32{1, 3, 2}},
			want: -1,
		},
		{
			x:    &testeditionspb.TestAllTypes{RepeatedSfixed32: []int32{1, 3, 2}},
			y:    &testeditionspb.TestAllTypes{RepeatedSfixed32: []int32{1, 2, 3}},
			want: 1,
		},
		{
			x:    &testeditionspb.TestAllTypes{RepeatedSfixed64: []int64{1, 2, 3}},
			y:    &testeditionspb.TestAllTypes{RepeatedSfixed64: []int64{1, 3, 2}},
			want: -1,
		},
		{
			x:    &testeditionspb.TestAllTypes{RepeatedSfixed64: []int64{1, 3, 2}},
			y:    &testeditionspb.TestAllTypes{RepeatedSfixed64: []int64{1, 2, 3}},
			want: 1,
		},
		{
			x:    &testeditionspb.TestAllTypes{RepeatedFloat: []float32{1, 2, 3}},
			y:    &testeditionspb.TestAllTypes{RepeatedFloat: []float32{1, 3, 2}},
			want: -1,
		},
		{
			x:    &testeditionspb.TestAllTypes{RepeatedFloat: []float32{1, 3, 2}},
			y:    &testeditionspb.TestAllTypes{RepeatedFloat: []float32{1, 2, 3}},
			want: 1,
		},
		{
			x:    &testeditionspb.TestAllTypes{RepeatedDouble: []float64{1, 2, 3}},
			y:    &testeditionspb.TestAllTypes{RepeatedDouble: []float64{1, 3, 2}},
			want: -1,
		},
		{
			x:    &testeditionspb.TestAllTypes{RepeatedDouble: []float64{1, 3, 2}},
			y:    &testeditionspb.TestAllTypes{RepeatedDouble: []float64{1, 2, 3}},
			want: 1,
		},
		{
			x:    &testeditionspb.TestAllTypes{RepeatedBool: []bool{true, false}},
			y:    &testeditionspb.TestAllTypes{RepeatedBool: []bool{true, true}},
			want: -1,
		},
		{
			x:    &testeditionspb.TestAllTypes{RepeatedBool: []bool{true, true}},
			y:    &testeditionspb.TestAllTypes{RepeatedBool: []bool{true, false}},
			want: 1,
		},
		{
			x:    &testeditionspb.TestAllTypes{RepeatedString: []string{"a", "b"}},
			y:    &testeditionspb.TestAllTypes{RepeatedString: []string{"a", "c"}},
			want: -1,
		},
		{
			x:    &testeditionspb.TestAllTypes{RepeatedString: []string{"a", "c"}},
			y:    &testeditionspb.TestAllTypes{RepeatedString: []string{"a", "b"}},
			want: 1,
		},
		{
			x:    &testeditionspb.TestAllTypes{RepeatedBytes: [][]byte{[]byte("a"), []byte("b")}},
			y:    &testeditionspb.TestAllTypes{RepeatedBytes: [][]byte{[]byte("a"), []byte("c")}},
			want: -1,
		},
		{
			x:    &testeditionspb.TestAllTypes{RepeatedBytes: [][]byte{[]byte("a"), []byte("c")}},
			y:    &testeditionspb.TestAllTypes{RepeatedBytes: [][]byte{[]byte("a"), []byte("b")}},
			want: 1,
		},
		{
			x: &testeditionspb.TestAllTypes{
				RepeatedNestedEnum: []testeditionspb.TestAllTypes_NestedEnum{testeditionspb.TestAllTypes_FOO},
			},
			y: &testeditionspb.TestAllTypes{
				RepeatedNestedEnum: []testeditionspb.TestAllTypes_NestedEnum{testeditionspb.TestAllTypes_BAR},
			},
			want: -1,
		},
		{
			x: &testeditionspb.TestAllTypes{
				RepeatedNestedEnum: []testeditionspb.TestAllTypes_NestedEnum{testeditionspb.TestAllTypes_BAR},
			},
			y: &testeditionspb.TestAllTypes{
				RepeatedNestedEnum: []testeditionspb.TestAllTypes_NestedEnum{testeditionspb.TestAllTypes_FOO},
			},
			want: 1,
		},
		{
			x: &testeditionspb.TestAllTypes{Repeatedgroup: []*testeditionspb.TestAllTypes_RepeatedGroup{
				{
					A: proto.Int32(1),
				}, {
					A: proto.Int32(2),
				},
			}},
			y: &testeditionspb.TestAllTypes{Repeatedgroup: []*testeditionspb.TestAllTypes_RepeatedGroup{
				{
					A: proto.Int32(1),
				}, {
					A: proto.Int32(3),
				},
			}},
			want: -1,
		},
		{
			x: &testeditionspb.TestAllTypes{Repeatedgroup: []*testeditionspb.TestAllTypes_RepeatedGroup{
				{
					A: proto.Int32(1),
				}, {
					A: proto.Int32(3),
				},
			}},
			y: &testeditionspb.TestAllTypes{Repeatedgroup: []*testeditionspb.TestAllTypes_RepeatedGroup{
				{
					A: proto.Int32(1),
				}, {
					A: proto.Int32(2),
				},
			}},
			want: 1,
		},
		{
			x: &testeditionspb.TestAllTypes{RepeatedNestedMessage: []*testeditionspb.TestAllTypes_NestedMessage{
				{
					A: proto.Int32(1),
				}, {
					A: proto.Int32(2),
				},
			}},
			y: &testeditionspb.TestAllTypes{RepeatedNestedMessage: []*testeditionspb.TestAllTypes_NestedMessage{
				{
					A: proto.Int32(1),
				}, {
					A: proto.Int32(3),
				},
			}},
			want: -1,
		},
		{
			x: &testeditionspb.TestAllTypes{RepeatedNestedMessage: []*testeditionspb.TestAllTypes_NestedMessage{
				{
					A: proto.Int32(1),
				}, {
					A: proto.Int32(3),
				},
			}},
			y: &testeditionspb.TestAllTypes{RepeatedNestedMessage: []*testeditionspb.TestAllTypes_NestedMessage{
				{
					A: proto.Int32(1),
				}, {
					A: proto.Int32(2),
				},
			}},
			want: 1,
		},

		// Maps: various configurations.
		{
			x:    &testpb.TestAllTypes{MapInt32Int32: map[int32]int32{1: 2}},
			y:    &testpb.TestAllTypes{MapInt32Int32: map[int32]int32{1: 2}},
			want: 0,
		},
		{
			x:    &testpb.TestAllTypes{MapInt32Int32: map[int32]int32{1: 2}},
			y:    &testpb.TestAllTypes{MapInt32Int32: map[int32]int32{3: 4}},
			want: -1,
		},
		{
			x:    &testpb.TestAllTypes{MapInt32Int32: map[int32]int32{3: 4}},
			y:    &testpb.TestAllTypes{MapInt32Int32: map[int32]int32{1: 2}},
			want: 1,
		},
		{
			x:    &testpb.TestAllTypes{MapInt32Int32: map[int32]int32{1: 2}},
			y:    &testpb.TestAllTypes{MapInt32Int32: map[int32]int32{1: 2, 3: 4}},
			want: -1,
		},
		{
			x:    &testpb.TestAllTypes{MapInt32Int32: map[int32]int32{1: 2, 3: 4}},
			y:    &testpb.TestAllTypes{MapInt32Int32: map[int32]int32{1: 2}},
			want: 1,
		},
		{
			x:    &testeditionspb.TestAllTypes{MapInt32Int32: map[int32]int32{1: 2}},
			y:    &testeditionspb.TestAllTypes{MapInt32Int32: map[int32]int32{3: 4}},
			want: -1,
		},
		{
			x:    &testeditionspb.TestAllTypes{MapInt32Int32: map[int32]int32{3: 4}},
			y:    &testeditionspb.TestAllTypes{MapInt32Int32: map[int32]int32{1: 2}},
			want: 1,
		},
		{
			x:    &testeditionspb.TestAllTypes{MapInt32Int32: map[int32]int32{1: 2}},
			y:    &testeditionspb.TestAllTypes{MapInt32Int32: map[int32]int32{1: 2, 3: 4}},
			want: -1,
		},
		{
			x:    &testeditionspb.TestAllTypes{MapInt32Int32: map[int32]int32{1: 2, 3: 4}},
			y:    &testeditionspb.TestAllTypes{MapInt32Int32: map[int32]int32{1: 2}},
			want: 1,
		},

		// Maps: various types.
		{
			x:    &testpb.TestAllTypes{MapInt32Int32: map[int32]int32{1: 2, 3: 4}},
			y:    &testpb.TestAllTypes{MapInt32Int32: map[int32]int32{1: 2, 3: 4}},
			want: 0,
		},
		{
			x:    &testpb.TestAllTypes{MapInt32Int32: map[int32]int32{1: 2, 3: 4}},
			y:    &testpb.TestAllTypes{MapInt32Int32: map[int32]int32{1: 2, 3: 5}},
			want: -1,
		},
		{
			x:    &testpb.TestAllTypes{MapInt32Int32: map[int32]int32{1: 2, 3: 5}},
			y:    &testpb.TestAllTypes{MapInt32Int32: map[int32]int32{1: 2, 3: 4}},
			want: 1,
		},
		{
			x:    &testpb.TestAllTypes{MapInt64Int64: map[int64]int64{1: 2, 3: 4}},
			y:    &testpb.TestAllTypes{MapInt64Int64: map[int64]int64{1: 2, 3: 5}},
			want: -1,
		},
		{
			x:    &testpb.TestAllTypes{MapInt64Int64: map[int64]int64{1: 2, 3: 5}},
			y:    &testpb.TestAllTypes{MapInt64Int64: map[int64]int64{1: 2, 3: 4}},
			want: 1,
		},
		{
			x:    &testpb.TestAllTypes{MapUint32Uint32: map[uint32]uint32{1: 2, 3: 4}},
			y:    &testpb.TestAllTypes{MapUint32Uint32: map[uint32]uint32{1: 2, 3: 5}},
			want: -1,
		},
		{
			x:    &testpb.TestAllTypes{MapUint32Uint32: map[uint32]uint32{1: 2, 3: 5}},
			y:    &testpb.TestAllTypes{MapUint32Uint32: map[uint32]uint32{1: 2, 3: 4}},
			want: 1,
		},
		{
			x:    &testpb.TestAllTypes{MapUint64Uint64: map[uint64]uint64{1: 2, 3: 4}},
			y:    &testpb.TestAllTypes{MapUint64Uint64: map[uint64]uint64{1: 2, 3: 5}},
			want: -1,
		},
		{
			x:    &testpb.TestAllTypes{MapUint64Uint64: map[uint64]uint64{1: 2, 3: 5}},
			y:    &testpb.TestAllTypes{MapUint64Uint64: map[uint64]uint64{1: 2, 3: 4}},
			want: 1,
		},
		{
			x:    &testpb.TestAllTypes{MapSint32Sint32: map[int32]int32{1: 2, 3: 4}},
			y:    &testpb.TestAllTypes{MapSint32Sint32: map[int32]int32{1: 2, 3: 5}},
			want: -1,
		},
		{
			x:    &testpb.TestAllTypes{MapSint32Sint32: map[int32]int32{1: 2, 3: 5}},
			y:    &testpb.TestAllTypes{MapSint32Sint32: map[int32]int32{1: 2, 3: 4}},
			want: 1,
		},
		{
			x:    &testpb.TestAllTypes{MapSint64Sint64: map[int64]int64{1: 2, 3: 4}},
			y:    &testpb.TestAllTypes{MapSint64Sint64: map[int64]int64{1: 2, 3: 5}},
			want: -1,
		},
		{
			x:    &testpb.TestAllTypes{MapSint64Sint64: map[int64]int64{1: 2, 3: 5}},
			y:    &testpb.TestAllTypes{MapSint64Sint64: map[int64]int64{1: 2, 3: 4}},
			want: 1,
		},
		{
			x:    &testpb.TestAllTypes{MapFixed32Fixed32: map[uint32]uint32{1: 2, 3: 4}},
			y:    &testpb.TestAllTypes{MapFixed32Fixed32: map[uint32]uint32{1: 2, 3: 5}},
			want: -1,
		},
		{
			x:    &testpb.TestAllTypes{MapFixed32Fixed32: map[uint32]uint32{1: 2, 3: 5}},
			y:    &testpb.TestAllTypes{MapFixed32Fixed32: map[uint32]uint32{1: 2, 3: 4}},
			want: 1,
		},
		{
			x:    &testpb.TestAllTypes{MapFixed64Fixed64: map[uint64]uint64{1: 2, 3: 4}},
			y:    &testpb.TestAllTypes{MapFixed64Fixed64: map[uint64]uint64{1: 2, 3: 5}},
			want: -1,
		},
		{
			x:    &testpb.TestAllTypes{MapFixed64Fixed64: map[uint64]uint64{1: 2, 3: 5}},
			y:    &testpb.TestAllTypes{MapFixed64Fixed64: map[uint64]uint64{1: 2, 3: 4}},
			want: 1,
		},
		{
			x:    &testpb.TestAllTypes{MapSfixed32Sfixed32: map[int32]int32{1: 2, 3: 4}},
			y:    &testpb.TestAllTypes{MapSfixed32Sfixed32: map[int32]int32{1: 2, 3: 5}},
			want: -1,
		},
		{
			x:    &testpb.TestAllTypes{MapSfixed32Sfixed32: map[int32]int32{1: 2, 3: 5}},
			y:    &testpb.TestAllTypes{MapSfixed32Sfixed32: map[int32]int32{1: 2, 3: 4}},
			want: 1,
		},
		{
			x:    &testpb.TestAllTypes{MapSfixed64Sfixed64: map[int64]int64{1: 2, 3: 4}},
			y:    &testpb.TestAllTypes{MapSfixed64Sfixed64: map[int64]int64{1: 2, 3: 5}},
			want: -1,
		},
		{
			x:    &testpb.TestAllTypes{MapSfixed64Sfixed64: map[int64]int64{1: 2, 3: 5}},
			y:    &testpb.TestAllTypes{MapSfixed64Sfixed64: map[int64]int64{1: 2, 3: 4}},
			want: 1,
		},
		{
			x:    &testpb.TestAllTypes{MapInt32Float: map[int32]float32{1: 2, 3: 4}},
			y:    &testpb.TestAllTypes{MapInt32Float: map[int32]float32{1: 2, 3: 5}},
			want: -1,
		},
		{
			x:    &testpb.TestAllTypes{MapInt32Float: map[int32]float32{1: 2, 3: 5}},
			y:    &testpb.TestAllTypes{MapInt32Float: map[int32]float32{1: 2, 3: 4}},
			want: 1,
		},
		{
			x:    &testpb.TestAllTypes{MapInt32Float: map[int32]float32{0: -1, 1: 1}},
			y:    &testpb.TestAllTypes{MapInt32Float: map[int32]float32{0: 1, 1: -1}},
			want: -1,
		},
		{
			x:    &testpb.TestAllTypes{MapInt32Float: map[int32]float32{0: 1, 1: -1}},
			y:    &testpb.TestAllTypes{MapInt32Float: map[int32]float32{0: -1, 1: 1}},
			want: 1,
		},
		{
			x:    &testpb.TestAllTypes{MapInt32Double: map[int32]float64{1: 2, 3: 4}},
			y:    &testpb.TestAllTypes{MapInt32Double: map[int32]float64{1: 2, 3: 5}},
			want: -1,
		},
		{
			x:    &testpb.TestAllTypes{MapInt32Double: map[int32]float64{1: 2, 3: 5}},
			y:    &testpb.TestAllTypes{MapInt32Double: map[int32]float64{1: 2, 3: 4}},
			want: 1,
		},
		{
			x:    &testpb.TestAllTypes{MapBoolBool: map[bool]bool{true: false, false: false}},
			y:    &testpb.TestAllTypes{MapBoolBool: map[bool]bool{true: false, false: true}},
			want: -1,
		},
		{
			x:    &testpb.TestAllTypes{MapBoolBool: map[bool]bool{true: false, false: true}},
			y:    &testpb.TestAllTypes{MapBoolBool: map[bool]bool{true: false, false: false}},
			want: 1,
		},
		{
			x:    &testpb.TestAllTypes{MapStringString: map[string]string{"a": "b", "c": "d"}},
			y:    &testpb.TestAllTypes{MapStringString: map[string]string{"a": "b", "c": "e"}},
			want: -1,
		},
		{
			x:    &testpb.TestAllTypes{MapStringString: map[string]string{"a": "b", "c": "e"}},
			y:    &testpb.TestAllTypes{MapStringString: map[string]string{"a": "b", "c": "d"}},
			want: 1,
		},
		{
			x:    &testpb.TestAllTypes{MapStringBytes: map[string][]byte{"a": []byte("b"), "c": []byte("d")}},
			y:    &testpb.TestAllTypes{MapStringBytes: map[string][]byte{"a": []byte("b"), "c": []byte("e")}},
			want: -1,
		},
		{
			x:    &testpb.TestAllTypes{MapStringBytes: map[string][]byte{"a": []byte("b"), "c": []byte("e")}},
			y:    &testpb.TestAllTypes{MapStringBytes: map[string][]byte{"a": []byte("b"), "c": []byte("d")}},
			want: 1,
		},
		{
			x: &testpb.TestAllTypes{MapStringNestedMessage: map[string]*testpb.TestAllTypes_NestedMessage{
				"a": {
					A: proto.Int32(1),
				}, "b": {
					A: proto.Int32(2),
				},
			}},
			y: &testpb.TestAllTypes{MapStringNestedMessage: map[string]*testpb.TestAllTypes_NestedMessage{
				"a": {
					A: proto.Int32(1),
				}, "b": {
					A: proto.Int32(3),
				},
			}},
			want: -1,
		},
		{
			x: &testpb.TestAllTypes{MapStringNestedMessage: map[string]*testpb.TestAllTypes_NestedMessage{
				"a": {
					A: proto.Int32(1),
				}, "b": {
					A: proto.Int32(3),
				},
			}},
			y: &testpb.TestAllTypes{MapStringNestedMessage: map[string]*testpb.TestAllTypes_NestedMessage{
				"a": {
					A: proto.Int32(1),
				}, "b": {
					A: proto.Int32(2),
				},
			}},
			want: 1,
		},
		{
			x: &testpb.TestAllTypes{
				MapStringNestedEnum: map[string]testpb.TestAllTypes_NestedEnum{
					"a": testpb.TestAllTypes_FOO,
					"b": testpb.TestAllTypes_BAR,
				},
			},
			y: &testpb.TestAllTypes{
				MapStringNestedEnum: map[string]testpb.TestAllTypes_NestedEnum{
					"a": testpb.TestAllTypes_FOO,
					"b": testpb.TestAllTypes_BAZ,
				},
			},
			want: -1,
		},
		{
			x: &testpb.TestAllTypes{
				MapStringNestedEnum: map[string]testpb.TestAllTypes_NestedEnum{
					"a": testpb.TestAllTypes_FOO,
					"b": testpb.TestAllTypes_BAZ,
				},
			},
			y: &testpb.TestAllTypes{
				MapStringNestedEnum: map[string]testpb.TestAllTypes_NestedEnum{
					"a": testpb.TestAllTypes_FOO,
					"b": testpb.TestAllTypes_BAR,
				},
			},
			want: 1,
		},
		{
			x:    &testeditionspb.TestAllTypes{MapInt32Int32: map[int32]int32{1: 2, 3: 4}},
			y:    &testeditionspb.TestAllTypes{MapInt32Int32: map[int32]int32{1: 2, 3: 5}},
			want: -1,
		},
		{
			x:    &testeditionspb.TestAllTypes{MapInt32Int32: map[int32]int32{1: 2, 3: 5}},
			y:    &testeditionspb.TestAllTypes{MapInt32Int32: map[int32]int32{1: 2, 3: 4}},
			want: 1,
		},
		{
			x:    &testeditionspb.TestAllTypes{MapInt64Int64: map[int64]int64{1: 2, 3: 4}},
			y:    &testeditionspb.TestAllTypes{MapInt64Int64: map[int64]int64{1: 2, 3: 5}},
			want: -1,
		},
		{
			x:    &testeditionspb.TestAllTypes{MapInt64Int64: map[int64]int64{1: 2, 3: 5}},
			y:    &testeditionspb.TestAllTypes{MapInt64Int64: map[int64]int64{1: 2, 3: 4}},
			want: 1,
		},
		{
			x:    &testeditionspb.TestAllTypes{MapUint32Uint32: map[uint32]uint32{1: 2, 3: 4}},
			y:    &testeditionspb.TestAllTypes{MapUint32Uint32: map[uint32]uint32{1: 2, 3: 5}},
			want: -1,
		},
		{
			x:    &testeditionspb.TestAllTypes{MapUint32Uint32: map[uint32]uint32{1: 2, 3: 5}},
			y:    &testeditionspb.TestAllTypes{MapUint32Uint32: map[uint32]uint32{1: 2, 3: 4}},
			want: 1,
		},
		{
			x:    &testeditionspb.TestAllTypes{MapUint64Uint64: map[uint64]uint64{1: 2, 3: 4}},
			y:    &testeditionspb.TestAllTypes{MapUint64Uint64: map[uint64]uint64{1: 2, 3: 5}},
			want: -1,
		},
		{
			x:    &testeditionspb.TestAllTypes{MapUint64Uint64: map[uint64]uint64{1: 2, 3: 5}},
			y:    &testeditionspb.TestAllTypes{MapUint64Uint64: map[uint64]uint64{1: 2, 3: 4}},
			want: 1,
		},
		{
			x:    &testeditionspb.TestAllTypes{MapSint32Sint32: map[int32]int32{1: 2, 3: 4}},
			y:    &testeditionspb.TestAllTypes{MapSint32Sint32: map[int32]int32{1: 2, 3: 5}},
			want: -1,
		},
		{
			x:    &testeditionspb.TestAllTypes{MapSint32Sint32: map[int32]int32{1: 2, 3: 5}},
			y:    &testeditionspb.TestAllTypes{MapSint32Sint32: map[int32]int32{1: 2, 3: 4}},
			want: 1,
		},
		{
			x:    &testeditionspb.TestAllTypes{MapSint64Sint64: map[int64]int64{1: 2, 3: 4}},
			y:    &testeditionspb.TestAllTypes{MapSint64Sint64: map[int64]int64{1: 2, 3: 5}},
			want: -1,
		},
		{
			x:    &testeditionspb.TestAllTypes{MapSint64Sint64: map[int64]int64{1: 2, 3: 5}},
			y:    &testeditionspb.TestAllTypes{MapSint64Sint64: map[int64]int64{1: 2, 3: 4}},
			want: 1,
		},
		{
			x:    &testeditionspb.TestAllTypes{MapFixed32Fixed32: map[uint32]uint32{1: 2, 3: 4}},
			y:    &testeditionspb.TestAllTypes{MapFixed32Fixed32: map[uint32]uint32{1: 2, 3: 5}},
			want: -1,
		},
		{
			x:    &testeditionspb.TestAllTypes{MapFixed32Fixed32: map[uint32]uint32{1: 2, 3: 5}},
			y:    &testeditionspb.TestAllTypes{MapFixed32Fixed32: map[uint32]uint32{1: 2, 3: 4}},
			want: 1,
		},
		{
			x:    &testeditionspb.TestAllTypes{MapFixed64Fixed64: map[uint64]uint64{1: 2, 3: 4}},
			y:    &testeditionspb.TestAllTypes{MapFixed64Fixed64: map[uint64]uint64{1: 2, 3: 5}},
			want: -1,
		},
		{
			x:    &testeditionspb.TestAllTypes{MapFixed64Fixed64: map[uint64]uint64{1: 2, 3: 5}},
			y:    &testeditionspb.TestAllTypes{MapFixed64Fixed64: map[uint64]uint64{1: 2, 3: 4}},
			want: 1,
		},
		{
			x:    &testeditionspb.TestAllTypes{MapSfixed32Sfixed32: map[int32]int32{1: 2, 3: 4}},
			y:    &testeditionspb.TestAllTypes{MapSfixed32Sfixed32: map[int32]int32{1: 2, 3: 5}},
			want: -1,
		},
		{
			x:    &testeditionspb.TestAllTypes{MapSfixed32Sfixed32: map[int32]int32{1: 2, 3: 5}},
			y:    &testeditionspb.TestAllTypes{MapSfixed32Sfixed32: map[int32]int32{1: 2, 3: 4}},
			want: 1,
		},
		{
			x:    &testeditionspb.TestAllTypes{MapSfixed64Sfixed64: map[int64]int64{1: 2, 3: 4}},
			y:    &testeditionspb.TestAllTypes{MapSfixed64Sfixed64: map[int64]int64{1: 2, 3: 5}},
			want: -1,
		},
		{
			x:    &testeditionspb.TestAllTypes{MapSfixed64Sfixed64: map[int64]int64{1: 2, 3: 5}},
			y:    &testeditionspb.TestAllTypes{MapSfixed64Sfixed64: map[int64]int64{1: 2, 3: 4}},
			want: 1,
		},
		{
			x:    &testeditionspb.TestAllTypes{MapInt32Float: map[int32]float32{1: 2, 3: 4}},
			y:    &testeditionspb.TestAllTypes{MapInt32Float: map[int32]float32{1: 2, 3: 5}},
			want: -1,
		},
		{
			x:    &testeditionspb.TestAllTypes{MapInt32Float: map[int32]float32{1: 2, 3: 5}},
			y:    &testeditionspb.TestAllTypes{MapInt32Float: map[int32]float32{1: 2, 3: 4}},
			want: 1,
		},
		{
			x:    &testeditionspb.TestAllTypes{MapInt32Float: map[int32]float32{0: -1, 1: 1}},
			y:    &testeditionspb.TestAllTypes{MapInt32Float: map[int32]float32{0: 1, 1: -1}},
			want: -1,
		},
		{
			x:    &testeditionspb.TestAllTypes{MapInt32Float: map[int32]float32{0: 1, 1: -1}},
			y:    &testeditionspb.TestAllTypes{MapInt32Float: map[int32]float32{0: -1, 1: 1}},
			want: 1,
		},
		{
			x:    &testeditionspb.TestAllTypes{MapInt32Double: map[int32]float64{1: 2, 3: 4}},
			y:    &testeditionspb.TestAllTypes{MapInt32Double: map[int32]float64{1: 2, 3: 5}},
			want: -1,
		},
		{
			x:    &testeditionspb.TestAllTypes{MapInt32Double: map[int32]float64{1: 2, 3: 5}},
			y:    &testeditionspb.TestAllTypes{MapInt32Double: map[int32]float64{1: 2, 3: 4}},
			want: 1,
		},
		{
			x:    &testeditionspb.TestAllTypes{MapBoolBool: map[bool]bool{true: false, false: false}},
			y:    &testeditionspb.TestAllTypes{MapBoolBool: map[bool]bool{true: false, false: true}},
			want: -1,
		},
		{
			x:    &testeditionspb.TestAllTypes{MapBoolBool: map[bool]bool{true: false, false: true}},
			y:    &testeditionspb.TestAllTypes{MapBoolBool: map[bool]bool{true: false, false: false}},
			want: 1,
		},
		{
			x:    &testeditionspb.TestAllTypes{MapStringString: map[string]string{"a": "b", "c": "d"}},
			y:    &testeditionspb.TestAllTypes{MapStringString: map[string]string{"a": "b", "c": "e"}},
			want: -1,
		},
		{
			x:    &testeditionspb.TestAllTypes{MapStringString: map[string]string{"a": "b", "c": "e"}},
			y:    &testeditionspb.TestAllTypes{MapStringString: map[string]string{"a": "b", "c": "d"}},
			want: 1,
		},
		{
			x:    &testeditionspb.TestAllTypes{MapStringBytes: map[string][]byte{"a": []byte("b"), "c": []byte("d")}},
			y:    &testeditionspb.TestAllTypes{MapStringBytes: map[string][]byte{"a": []byte("b"), "c": []byte("e")}},
			want: -1,
		},
		{
			x:    &testeditionspb.TestAllTypes{MapStringBytes: map[string][]byte{"a": []byte("b"), "c": []byte("e")}},
			y:    &testeditionspb.TestAllTypes{MapStringBytes: map[string][]byte{"a": []byte("b"), "c": []byte("d")}},
			want: 1,
		},
		{
			x: &testeditionspb.TestAllTypes{MapStringNestedMessage: map[string]*testeditionspb.TestAllTypes_NestedMessage{
				"a": {
					A: proto.Int32(1),
				}, "b": {
					A: proto.Int32(2),
				},
			}},
			y: &testeditionspb.TestAllTypes{MapStringNestedMessage: map[string]*testeditionspb.TestAllTypes_NestedMessage{
				"a": {
					A: proto.Int32(1),
				}, "b": {
					A: proto.Int32(3),
				},
			}},
			want: -1,
		},
		{
			x: &testeditionspb.TestAllTypes{MapStringNestedMessage: map[string]*testeditionspb.TestAllTypes_NestedMessage{
				"a": {
					A: proto.Int32(1),
				}, "b": {
					A: proto.Int32(3),
				},
			}},
			y: &testeditionspb.TestAllTypes{MapStringNestedMessage: map[string]*testeditionspb.TestAllTypes_NestedMessage{
				"a": {
					A: proto.Int32(1),
				}, "b": {
					A: proto.Int32(2),
				},
			}},
			want: 1,
		},
		{
			x: &testeditionspb.TestAllTypes{
				MapStringNestedEnum: map[string]testeditionspb.TestAllTypes_NestedEnum{
					"a": testeditionspb.TestAllTypes_FOO,
					"b": testeditionspb.TestAllTypes_BAR,
				},
			},
			y: &testeditionspb.TestAllTypes{
				MapStringNestedEnum: map[string]testeditionspb.TestAllTypes_NestedEnum{
					"a": testeditionspb.TestAllTypes_FOO,
					"b": testeditionspb.TestAllTypes_BAZ,
				},
			},
			want: -1,
		},
		{
			x: &testeditionspb.TestAllTypes{
				MapStringNestedEnum: map[string]testeditionspb.TestAllTypes_NestedEnum{
					"a": testeditionspb.TestAllTypes_FOO,
					"b": testeditionspb.TestAllTypes_BAZ,
				},
			},
			y: &testeditionspb.TestAllTypes{
				MapStringNestedEnum: map[string]testeditionspb.TestAllTypes_NestedEnum{
					"a": testeditionspb.TestAllTypes_FOO,
					"b": testeditionspb.TestAllTypes_BAR,
				},
			},
			want: 1,
		},

		// Extensions.
		{
			x: build(&testpb.TestAllExtensions{},
				extend(testpb.E_OptionalInt32, int32(1)),
			),
			y: build(&testpb.TestAllExtensions{},
				extend(testpb.E_OptionalInt32, int32(2)),
			),
			want: -1,
		},
		{
			x: build(&testpb.TestAllExtensions{},
				extend(testpb.E_OptionalInt32, int32(2)),
			),
			y: build(&testpb.TestAllExtensions{},
				extend(testpb.E_OptionalInt32, int32(1)),
			),
			want: 1,
		},
		{
			x: &testpb.TestAllExtensions{},
			y: build(&testpb.TestAllExtensions{},
				extend(testpb.E_OptionalInt32, int32(2)),
			),
			want: -1,
		},
		{
			x: build(&testpb.TestAllExtensions{},
				extend(testpb.E_OptionalInt32, int32(2)),
			),
			y:    &testpb.TestAllExtensions{},
			want: 1,
		},
		{
			x: build(&testpb.TestAllTypes{}, unknown(protopack.Message{
				protopack.Tag{Number: 100000, Type: protopack.VarintType}, protopack.Varint(1),
			}.Marshal())),
			y: build(&testpb.TestAllTypes{}, unknown(protopack.Message{
				protopack.Tag{Number: 100000, Type: protopack.VarintType}, protopack.Varint(2),
			}.Marshal())),
			want: -1,
		},
		{
			x: build(&testpb.TestAllTypes{}, unknown(protopack.Message{
				protopack.Tag{Number: 100000, Type: protopack.VarintType}, protopack.Varint(2),
			}.Marshal())),
			y: build(&testpb.TestAllTypes{}, unknown(protopack.Message{
				protopack.Tag{Number: 100000, Type: protopack.VarintType}, protopack.Varint(1),
			}.Marshal())),
			want: 1,
		},
		{
			x: build(&testpb.TestAllTypes{}, unknown(protopack.Message{
				protopack.Tag{Number: 100000, Type: protopack.VarintType}, protopack.Varint(1),
			}.Marshal())),
			y:    &testpb.TestAllTypes{},
			want: 1,
		},
		{
			x: &testpb.TestAllTypes{},
			y: build(&testpb.TestAllTypes{}, unknown(protopack.Message{
				protopack.Tag{Number: 100000, Type: protopack.VarintType}, protopack.Varint(1),
			}.Marshal())),
			want: -1,
		},
		{
			x: build(&testeditionspb.TestAllTypes{}, unknown(protopack.Message{
				protopack.Tag{Number: 100000, Type: protopack.VarintType}, protopack.Varint(2),
			}.Marshal())),
			y: build(&testeditionspb.TestAllTypes{}, unknown(protopack.Message{
				protopack.Tag{Number: 100000, Type: protopack.VarintType}, protopack.Varint(1),
			}.Marshal())),
			want: 1,
		},
		{
			x: build(&testeditionspb.TestAllTypes{}, unknown(protopack.Message{
				protopack.Tag{Number: 100000, Type: protopack.VarintType}, protopack.Varint(1),
			}.Marshal())),
			y: build(&testeditionspb.TestAllTypes{}, unknown(protopack.Message{
				protopack.Tag{Number: 100000, Type: protopack.VarintType}, protopack.Varint(2),
			}.Marshal())),
			want: -1,
		},
		{
			x: &testeditionspb.TestAllTypes{},
			y: build(&testeditionspb.TestAllTypes{}, unknown(protopack.Message{
				protopack.Tag{Number: 100000, Type: protopack.VarintType}, protopack.Varint(1),
			}.Marshal())),
			want: -1,
		},
		{
			x: build(&testeditionspb.TestAllTypes{}, unknown(protopack.Message{
				protopack.Tag{Number: 100000, Type: protopack.VarintType}, protopack.Varint(1),
			}.Marshal())),
			y:    &testeditionspb.TestAllTypes{},
			want: 1,
		},
		{
			x: build(&testeditionspb.TestAllTypes{}, unknown(protopack.Message{
				protopack.Tag{Number: 1000, Type: protopack.VarintType}, protopack.Varint(1),
			}.Marshal())),
			y: build(&testeditionspb.TestAllTypes{}, unknown(protopack.Message{
				protopack.Tag{Number: 100000, Type: protopack.VarintType}, protopack.Varint(1),
			}.Marshal())),
			want: -1,
		},
		{
			x: build(&testeditionspb.TestAllTypes{}, unknown(protopack.Message{
				protopack.Tag{Number: 100000, Type: protopack.VarintType}, protopack.Varint(1),
			}.Marshal())),
			y: build(&testeditionspb.TestAllTypes{}, unknown(protopack.Message{
				protopack.Tag{Number: 1000, Type: protopack.VarintType}, protopack.Varint(1),
			}.Marshal())),
			want: 1,
		},
		{
			x: build(&testeditionspb.TestAllTypes{}, unknown(protopack.Message{
				protopack.Tag{Number: 100000, Type: protopack.Fixed32Type}, protopack.Int32(1),
			}.Marshal())),
			y: build(&testeditionspb.TestAllTypes{}, unknown(protopack.Message{
				protopack.Tag{Number: 100000, Type: protopack.VarintType}, protopack.Varint(1),
			}.Marshal())),
			want: 1,
		},
		{
			x: build(&testeditionspb.TestAllTypes{}, unknown(protopack.Message{
				protopack.Tag{Number: 100000, Type: protopack.VarintType}, protopack.Varint(1),
			}.Marshal())),
			y: build(&testeditionspb.TestAllTypes{}, unknown(protopack.Message{
				protopack.Tag{Number: 100000, Type: protopack.Fixed32Type}, protopack.Int32(1),
			}.Marshal())),
			want: -1,
		},
		{
			x: build(&testeditionspb.TestAllTypes{}, unknown(protopack.Message{
				protopack.Tag{Number: 100000, Type: protopack.VarintType}, protopack.Varint(1),
			}.Marshal())),
			y: build(&testeditionspb.TestAllTypes{}, unknown(protopack.Message{
				protopack.Tag{Number: 100000, Type: protopack.VarintType}, protopack.Varint(1),
			}.Marshal())),
			want: 0,
		},
	}

	for _, tt := range tests {
		if got := proto.Compare(tt.x, tt.y); got != tt.want {
			t.Errorf(
				"Compare(x, y) = %v, want %v\n==== x ====\n%v\n==== y ====\n%v",
				got,
				tt.want,
				prototext.Format(tt.x),
				prototext.Format(tt.y),
			)
		}
	}
}

func TestLessThan(t *testing.T) {
	identicalPtrPb := &testpb.TestAllTypes{MapStringString: map[string]string{"a": "b", "c": "d"}}

	type incomparableMessage struct {
		*testpb.TestAllTypes
		pragma.DoNotCompare
	}

	tests := []struct {
		x, y   proto.Message
		isLess bool
	}{
		{
			x:      nil,
			y:      nil,
			isLess: false,
		},
		{
			x:      nil,
			y:      (*testpb.TestAllTypes)(nil),
			isLess: true,
		},
		{
			x:      (*testpb.TestAllTypes)(nil),
			y:      nil,
			isLess: false,
		},
		{
			x:      (*testpb.TestAllTypes)(nil),
			y:      (*testpb.TestAllTypes)(nil),
			isLess: false,
		},
		{
			x:      new(testpb.TestAllTypes),
			y:      (*testpb.TestAllTypes)(nil),
			isLess: false,
		},
		{
			x:      (*testpb.TestAllTypes)(nil),
			y:      new(testpb.TestAllTypes),
			isLess: true,
		},
		{
			x:      new(testpb.TestAllTypes),
			y:      new(testpb.TestAllTypes),
			isLess: false,
		},
		{
			x:      (*testpb.TestAllTypes)(nil),
			y:      (*testpb.TestAllExtensions)(nil),
			isLess: false,
		},
		{
			x:      (*testpb.TestAllExtensions)(nil),
			y:      (*testpb.TestAllTypes)(nil),
			isLess: true,
		},
		{
			x:      (*testpb.TestAllTypes)(nil),
			y:      new(testpb.TestAllExtensions),
			isLess: true,
		},
		{
			x:      new(testpb.TestAllExtensions),
			y:      (*testpb.TestAllTypes)(nil),
			isLess: false,
		},
		{
			x:      (*testpb.TestAllExtensions)(nil),
			y:      new(testpb.TestAllTypes),
			isLess: true,
		},
		{
			x:      new(testpb.TestAllTypes),
			y:      (*testpb.TestAllExtensions)(nil),
			isLess: false,
		},
		{
			x:      new(testpb.TestAllExtensions),
			y:      new(testpb.TestAllTypes),
			isLess: true,
		},
		{
			x:      new(testpb.TestAllTypes),
			y:      new(testpb.TestAllExtensions),
			isLess: false,
		},

		// Identical input pointers
		{
			x:      identicalPtrPb,
			y:      identicalPtrPb,
			isLess: false,
		},

		// Incomparable types. The top-level types are not actually directly
		// compared (which would panic), but rather the comparison happens on the
		// objects returned by ProtoReflect(). These tests are here just to ensure
		// that any short-circuit checks do not accidentally try to compare
		// incomparable top-level types.
		{
			x:      incomparableMessage{TestAllTypes: identicalPtrPb},
			y:      incomparableMessage{TestAllTypes: identicalPtrPb},
			isLess: false,
		},
		{
			x:      identicalPtrPb,
			y:      incomparableMessage{TestAllTypes: identicalPtrPb},
			isLess: false,
		},
		{
			x:      identicalPtrPb,
			y:      &incomparableMessage{TestAllTypes: identicalPtrPb},
			isLess: false,
		},

		// Proto2 scalars.
		{
			x: &testpb.TestAllTypes{
				OptionalInt32: proto.Int32(-1),
				OptionalInt64: proto.Int64(1),
			},
			y: &testpb.TestAllTypes{
				OptionalInt32: proto.Int32(1),
				OptionalInt64: proto.Int64(-1),
			},
			isLess: true,
		},
		{
			x: &testpb.TestAllTypes{
				OptionalInt32: proto.Int32(1),
				OptionalInt64: proto.Int64(-1),
			},
			y: &testpb.TestAllTypes{
				OptionalInt32: proto.Int32(-1),
				OptionalInt64: proto.Int64(1),
			},
			isLess: false,
		},
		{
			x:      &testpb.TestAllTypes{OptionalInt32: proto.Int32(1)},
			y:      &testpb.TestAllTypes{OptionalInt32: proto.Int32(2)},
			isLess: true,
		},
		{
			x:      &testpb.TestAllTypes{OptionalInt32: proto.Int32(2)},
			y:      &testpb.TestAllTypes{OptionalInt32: proto.Int32(1)},
			isLess: false,
		},
		{
			x:      &testpb.TestAllTypes{OptionalInt64: proto.Int64(1)},
			y:      &testpb.TestAllTypes{OptionalInt64: proto.Int64(2)},
			isLess: true,
		},
		{
			x:      &testpb.TestAllTypes{OptionalInt64: proto.Int64(2)},
			y:      &testpb.TestAllTypes{OptionalInt64: proto.Int64(1)},
			isLess: false,
		},
		{
			x:      &testpb.TestAllTypes{OptionalUint32: proto.Uint32(1)},
			y:      &testpb.TestAllTypes{OptionalUint32: proto.Uint32(2)},
			isLess: true,
		},
		{
			x:      &testpb.TestAllTypes{OptionalUint32: proto.Uint32(2)},
			y:      &testpb.TestAllTypes{OptionalUint32: proto.Uint32(1)},
			isLess: false,
		},
		{
			x:      &testpb.TestAllTypes{OptionalUint64: proto.Uint64(1)},
			y:      &testpb.TestAllTypes{OptionalUint64: proto.Uint64(2)},
			isLess: true,
		},
		{
			x:      &testpb.TestAllTypes{OptionalUint64: proto.Uint64(2)},
			y:      &testpb.TestAllTypes{OptionalUint64: proto.Uint64(1)},
			isLess: false,
		},
		{
			x:      &testpb.TestAllTypes{OptionalSint32: proto.Int32(1)},
			y:      &testpb.TestAllTypes{OptionalSint32: proto.Int32(2)},
			isLess: true,
		},
		{
			x:      &testpb.TestAllTypes{OptionalSint32: proto.Int32(2)},
			y:      &testpb.TestAllTypes{OptionalSint32: proto.Int32(1)},
			isLess: false,
		},
		{
			x:      &testpb.TestAllTypes{OptionalSint64: proto.Int64(1)},
			y:      &testpb.TestAllTypes{OptionalSint64: proto.Int64(2)},
			isLess: true,
		},
		{
			x:      &testpb.TestAllTypes{OptionalSint64: proto.Int64(2)},
			y:      &testpb.TestAllTypes{OptionalSint64: proto.Int64(1)},
			isLess: false,
		},
		{
			x:      &testpb.TestAllTypes{OptionalFixed32: proto.Uint32(1)},
			y:      &testpb.TestAllTypes{OptionalFixed32: proto.Uint32(2)},
			isLess: true,
		},
		{
			x:      &testpb.TestAllTypes{OptionalFixed32: proto.Uint32(2)},
			y:      &testpb.TestAllTypes{OptionalFixed32: proto.Uint32(1)},
			isLess: false,
		},
		{
			x:      &testpb.TestAllTypes{OptionalFixed64: proto.Uint64(1)},
			y:      &testpb.TestAllTypes{OptionalFixed64: proto.Uint64(2)},
			isLess: true,
		},
		{
			x:      &testpb.TestAllTypes{OptionalFixed64: proto.Uint64(2)},
			y:      &testpb.TestAllTypes{OptionalFixed64: proto.Uint64(1)},
			isLess: false,
		},
		{
			x:      &testpb.TestAllTypes{OptionalSfixed32: proto.Int32(1)},
			y:      &testpb.TestAllTypes{OptionalSfixed32: proto.Int32(2)},
			isLess: true,
		},
		{
			x:      &testpb.TestAllTypes{OptionalSfixed32: proto.Int32(2)},
			y:      &testpb.TestAllTypes{OptionalSfixed32: proto.Int32(1)},
			isLess: false,
		},
		{
			x:      &testpb.TestAllTypes{OptionalSfixed64: proto.Int64(1)},
			y:      &testpb.TestAllTypes{OptionalSfixed64: proto.Int64(2)},
			isLess: true,
		},
		{
			x:      &testpb.TestAllTypes{OptionalSfixed64: proto.Int64(2)},
			y:      &testpb.TestAllTypes{OptionalSfixed64: proto.Int64(1)},
			isLess: false,
		},
		{
			x:      &testpb.TestAllTypes{OptionalFloat: proto.Float32(1)},
			y:      &testpb.TestAllTypes{OptionalFloat: proto.Float32(2)},
			isLess: true,
		},
		{
			x:      &testpb.TestAllTypes{OptionalFloat: proto.Float32(2)},
			y:      &testpb.TestAllTypes{OptionalFloat: proto.Float32(1)},
			isLess: false,
		},
		{
			x:      &testpb.TestAllTypes{OptionalDouble: proto.Float64(1)},
			y:      &testpb.TestAllTypes{OptionalDouble: proto.Float64(2)},
			isLess: true,
		},
		{
			x:      &testpb.TestAllTypes{OptionalDouble: proto.Float64(2)},
			y:      &testpb.TestAllTypes{OptionalDouble: proto.Float64(1)},
			isLess: false,
		},
		{
			x:      &testpb.TestAllTypes{OptionalFloat: proto.Float32(float32(math.NaN()))},
			y:      &testpb.TestAllTypes{OptionalFloat: proto.Float32(0)},
			isLess: true,
		},
		{
			x:      &testpb.TestAllTypes{OptionalFloat: proto.Float32(0)},
			y:      &testpb.TestAllTypes{OptionalFloat: proto.Float32(float32(math.NaN()))},
			isLess: false,
		},
		{
			x:      &testpb.TestAllTypes{OptionalDouble: proto.Float64(float64(math.NaN()))},
			y:      &testpb.TestAllTypes{OptionalDouble: proto.Float64(0)},
			isLess: true,
		},
		{
			x:      &testpb.TestAllTypes{OptionalDouble: proto.Float64(0)},
			y:      &testpb.TestAllTypes{OptionalDouble: proto.Float64(float64(math.NaN()))},
			isLess: false,
		},
		{
			x:      &testpb.TestAllTypes{OptionalBool: proto.Bool(false)},
			y:      &testpb.TestAllTypes{OptionalBool: proto.Bool(true)},
			isLess: true,
		},
		{
			x:      &testpb.TestAllTypes{OptionalBool: proto.Bool(true)},
			y:      &testpb.TestAllTypes{OptionalBool: proto.Bool(false)},
			isLess: false,
		},
		{
			x:      &testpb.TestAllTypes{OptionalString: proto.String("a")},
			y:      &testpb.TestAllTypes{OptionalString: proto.String("b")},
			isLess: true,
		},
		{
			x:      &testpb.TestAllTypes{OptionalString: proto.String("b")},
			y:      &testpb.TestAllTypes{OptionalString: proto.String("a")},
			isLess: false,
		},
		{
			x:      &testpb.TestAllTypes{OptionalBytes: []byte("a")},
			y:      &testpb.TestAllTypes{OptionalBytes: []byte("b")},
			isLess: true,
		},
		{
			x:      &testpb.TestAllTypes{OptionalBytes: []byte("b")},
			y:      &testpb.TestAllTypes{OptionalBytes: []byte("a")},
			isLess: false,
		},
		{
			x:      &testpb.TestAllTypes{OptionalNestedEnum: testpb.TestAllTypes_FOO.Enum()},
			y:      &testpb.TestAllTypes{OptionalNestedEnum: testpb.TestAllTypes_BAR.Enum()},
			isLess: true,
		},
		{
			x:      &testpb.TestAllTypes{OptionalNestedEnum: testpb.TestAllTypes_BAR.Enum()},
			y:      &testpb.TestAllTypes{OptionalNestedEnum: testpb.TestAllTypes_FOO.Enum()},
			isLess: false,
		},
		{
			x:      &testpb.TestAllTypes{OptionalInt32: proto.Int32(2)},
			y:      &testpb.TestAllTypes{OptionalInt32: proto.Int32(2)},
			isLess: false,
		},
		{
			x:      &testpb.TestAllTypes{OptionalInt64: proto.Int64(2)},
			y:      &testpb.TestAllTypes{OptionalInt64: proto.Int64(2)},
			isLess: false,
		},
		{
			x:      &testpb.TestAllTypes{OptionalUint32: proto.Uint32(2)},
			y:      &testpb.TestAllTypes{OptionalUint32: proto.Uint32(2)},
			isLess: false,
		},
		{
			x:      &testpb.TestAllTypes{OptionalUint64: proto.Uint64(2)},
			y:      &testpb.TestAllTypes{OptionalUint64: proto.Uint64(2)},
			isLess: false,
		},
		{
			x:      &testpb.TestAllTypes{OptionalSint32: proto.Int32(2)},
			y:      &testpb.TestAllTypes{OptionalSint32: proto.Int32(2)},
			isLess: false,
		},
		{
			x:      &testpb.TestAllTypes{OptionalSint64: proto.Int64(2)},
			y:      &testpb.TestAllTypes{OptionalSint64: proto.Int64(2)},
			isLess: false,
		},
		{
			x:      &testpb.TestAllTypes{OptionalFixed32: proto.Uint32(2)},
			y:      &testpb.TestAllTypes{OptionalFixed32: proto.Uint32(2)},
			isLess: false,
		},
		{
			x:      &testpb.TestAllTypes{OptionalFixed64: proto.Uint64(2)},
			y:      &testpb.TestAllTypes{OptionalFixed64: proto.Uint64(2)},
			isLess: false,
		},
		{
			x:      &testpb.TestAllTypes{OptionalSfixed32: proto.Int32(2)},
			y:      &testpb.TestAllTypes{OptionalSfixed32: proto.Int32(2)},
			isLess: false,
		},
		{
			x:      &testpb.TestAllTypes{OptionalSfixed64: proto.Int64(2)},
			y:      &testpb.TestAllTypes{OptionalSfixed64: proto.Int64(2)},
			isLess: false,
		},
		{
			x:      &testpb.TestAllTypes{OptionalFloat: proto.Float32(2)},
			y:      &testpb.TestAllTypes{OptionalFloat: proto.Float32(2)},
			isLess: false,
		},
		{
			x:      &testpb.TestAllTypes{OptionalDouble: proto.Float64(2)},
			y:      &testpb.TestAllTypes{OptionalDouble: proto.Float64(2)},
			isLess: false,
		},
		{
			x:      &testpb.TestAllTypes{OptionalFloat: proto.Float32(float32(math.NaN()))},
			y:      &testpb.TestAllTypes{OptionalFloat: proto.Float32(float32(math.NaN()))},
			isLess: false,
		},
		{
			x:      &testpb.TestAllTypes{OptionalDouble: proto.Float64(float64(math.NaN()))},
			y:      &testpb.TestAllTypes{OptionalDouble: proto.Float64(float64(math.NaN()))},
			isLess: false,
		},
		{
			x:      &testpb.TestAllTypes{OptionalBool: proto.Bool(true)},
			y:      &testpb.TestAllTypes{OptionalBool: proto.Bool(true)},
			isLess: false,
		},
		{
			x:      &testpb.TestAllTypes{OptionalString: proto.String("abc")},
			y:      &testpb.TestAllTypes{OptionalString: proto.String("abc")},
			isLess: false,
		},
		{
			x:      &testpb.TestAllTypes{OptionalBytes: []byte("abc")},
			y:      &testpb.TestAllTypes{OptionalBytes: []byte("abc")},
			isLess: false,
		},
		{
			x:      &testpb.TestAllTypes{OptionalNestedEnum: testpb.TestAllTypes_FOO.Enum()},
			y:      &testpb.TestAllTypes{OptionalNestedEnum: testpb.TestAllTypes_FOO.Enum()},
			isLess: false,
		},

		// Editions scalars.
		{
			x:      &testeditionspb.TestAllTypes{OptionalInt32: proto.Int32(1)},
			y:      &testeditionspb.TestAllTypes{OptionalInt64: proto.Int64(2)},
			isLess: true,
		},
		{
			x:      &testeditionspb.TestAllTypes{OptionalInt64: proto.Int64(2)},
			y:      &testeditionspb.TestAllTypes{OptionalInt32: proto.Int32(1)},
			isLess: false,
		},
		{
			x:      &testeditionspb.TestAllTypes{OptionalInt32: proto.Int32(1)},
			y:      &testeditionspb.TestAllTypes{OptionalInt32: proto.Int32(2)},
			isLess: true,
		},
		{
			x:      &testeditionspb.TestAllTypes{OptionalInt32: proto.Int32(2)},
			y:      &testeditionspb.TestAllTypes{OptionalInt32: proto.Int32(1)},
			isLess: false,
		},
		{
			x:      &testeditionspb.TestAllTypes{OptionalInt64: proto.Int64(1)},
			y:      &testeditionspb.TestAllTypes{OptionalInt64: proto.Int64(2)},
			isLess: true,
		},
		{
			x:      &testeditionspb.TestAllTypes{OptionalInt64: proto.Int64(2)},
			y:      &testeditionspb.TestAllTypes{OptionalInt64: proto.Int64(1)},
			isLess: false,
		},
		{
			x:      &testeditionspb.TestAllTypes{OptionalUint32: proto.Uint32(1)},
			y:      &testeditionspb.TestAllTypes{OptionalUint32: proto.Uint32(2)},
			isLess: true,
		},
		{
			x:      &testeditionspb.TestAllTypes{OptionalUint32: proto.Uint32(2)},
			y:      &testeditionspb.TestAllTypes{OptionalUint32: proto.Uint32(1)},
			isLess: false,
		},
		{
			x:      &testeditionspb.TestAllTypes{OptionalUint64: proto.Uint64(1)},
			y:      &testeditionspb.TestAllTypes{OptionalUint64: proto.Uint64(2)},
			isLess: true,
		},
		{
			x:      &testeditionspb.TestAllTypes{OptionalUint64: proto.Uint64(2)},
			y:      &testeditionspb.TestAllTypes{OptionalUint64: proto.Uint64(1)},
			isLess: false,
		},
		{
			x:      &testeditionspb.TestAllTypes{OptionalSint32: proto.Int32(1)},
			y:      &testeditionspb.TestAllTypes{OptionalSint32: proto.Int32(2)},
			isLess: true,
		},
		{
			x:      &testeditionspb.TestAllTypes{OptionalSint32: proto.Int32(2)},
			y:      &testeditionspb.TestAllTypes{OptionalSint32: proto.Int32(1)},
			isLess: false,
		},
		{
			x:      &testeditionspb.TestAllTypes{OptionalSint64: proto.Int64(1)},
			y:      &testeditionspb.TestAllTypes{OptionalSint64: proto.Int64(2)},
			isLess: true,
		},
		{
			x:      &testeditionspb.TestAllTypes{OptionalSint64: proto.Int64(2)},
			y:      &testeditionspb.TestAllTypes{OptionalSint64: proto.Int64(1)},
			isLess: false,
		},
		{
			x:      &testeditionspb.TestAllTypes{OptionalFixed32: proto.Uint32(1)},
			y:      &testeditionspb.TestAllTypes{OptionalFixed32: proto.Uint32(2)},
			isLess: true,
		},
		{
			x:      &testeditionspb.TestAllTypes{OptionalFixed32: proto.Uint32(2)},
			y:      &testeditionspb.TestAllTypes{OptionalFixed32: proto.Uint32(1)},
			isLess: false,
		},
		{
			x:      &testeditionspb.TestAllTypes{OptionalFixed64: proto.Uint64(1)},
			y:      &testeditionspb.TestAllTypes{OptionalFixed64: proto.Uint64(2)},
			isLess: true,
		},
		{
			x:      &testeditionspb.TestAllTypes{OptionalFixed64: proto.Uint64(2)},
			y:      &testeditionspb.TestAllTypes{OptionalFixed64: proto.Uint64(1)},
			isLess: false,
		},
		{
			x:      &testeditionspb.TestAllTypes{OptionalSfixed32: proto.Int32(1)},
			y:      &testeditionspb.TestAllTypes{OptionalSfixed32: proto.Int32(2)},
			isLess: true,
		},
		{
			x:      &testeditionspb.TestAllTypes{OptionalSfixed32: proto.Int32(2)},
			y:      &testeditionspb.TestAllTypes{OptionalSfixed32: proto.Int32(1)},
			isLess: false,
		},
		{
			x:      &testeditionspb.TestAllTypes{OptionalSfixed64: proto.Int64(1)},
			y:      &testeditionspb.TestAllTypes{OptionalSfixed64: proto.Int64(2)},
			isLess: true,
		},
		{
			x:      &testeditionspb.TestAllTypes{OptionalSfixed64: proto.Int64(2)},
			y:      &testeditionspb.TestAllTypes{OptionalSfixed64: proto.Int64(1)},
			isLess: false,
		},
		{
			x:      &testeditionspb.TestAllTypes{OptionalFloat: proto.Float32(1)},
			y:      &testeditionspb.TestAllTypes{OptionalFloat: proto.Float32(2)},
			isLess: true,
		},
		{
			x:      &testeditionspb.TestAllTypes{OptionalFloat: proto.Float32(2)},
			y:      &testeditionspb.TestAllTypes{OptionalFloat: proto.Float32(1)},
			isLess: false,
		},
		{
			x:      &testeditionspb.TestAllTypes{OptionalDouble: proto.Float64(1)},
			y:      &testeditionspb.TestAllTypes{OptionalDouble: proto.Float64(2)},
			isLess: true,
		},
		{
			x:      &testeditionspb.TestAllTypes{OptionalDouble: proto.Float64(2)},
			y:      &testeditionspb.TestAllTypes{OptionalDouble: proto.Float64(1)},
			isLess: false,
		},
		{
			x:      &testeditionspb.TestAllTypes{OptionalFloat: proto.Float32(float32(math.NaN()))},
			y:      &testeditionspb.TestAllTypes{OptionalFloat: proto.Float32(0)},
			isLess: true,
		},
		{
			x:      &testeditionspb.TestAllTypes{OptionalFloat: proto.Float32(0)},
			y:      &testeditionspb.TestAllTypes{OptionalFloat: proto.Float32(float32(math.NaN()))},
			isLess: false,
		},
		{
			x:      &testeditionspb.TestAllTypes{OptionalDouble: proto.Float64(float64(math.NaN()))},
			y:      &testeditionspb.TestAllTypes{OptionalDouble: proto.Float64(0)},
			isLess: true,
		},
		{
			x:      &testeditionspb.TestAllTypes{OptionalDouble: proto.Float64(0)},
			y:      &testeditionspb.TestAllTypes{OptionalDouble: proto.Float64(float64(math.NaN()))},
			isLess: false,
		},
		{
			x:      &testeditionspb.TestAllTypes{OptionalBool: proto.Bool(false)},
			y:      &testeditionspb.TestAllTypes{OptionalBool: proto.Bool(true)},
			isLess: true,
		},
		{
			x:      &testeditionspb.TestAllTypes{OptionalBool: proto.Bool(true)},
			y:      &testeditionspb.TestAllTypes{OptionalBool: proto.Bool(false)},
			isLess: false,
		},
		{
			x:      &testeditionspb.TestAllTypes{OptionalString: proto.String("a")},
			y:      &testeditionspb.TestAllTypes{OptionalString: proto.String("b")},
			isLess: true,
		},
		{
			x:      &testeditionspb.TestAllTypes{OptionalString: proto.String("b")},
			y:      &testeditionspb.TestAllTypes{OptionalString: proto.String("a")},
			isLess: false,
		},
		{
			x:      &testeditionspb.TestAllTypes{OptionalBytes: []byte("a")},
			y:      &testeditionspb.TestAllTypes{OptionalBytes: []byte("b")},
			isLess: true,
		},
		{
			x:      &testeditionspb.TestAllTypes{OptionalBytes: []byte("b")},
			y:      &testeditionspb.TestAllTypes{OptionalBytes: []byte("a")},
			isLess: false,
		},
		{
			x:      &testeditionspb.TestAllTypes{OptionalNestedEnum: testeditionspb.TestAllTypes_FOO.Enum()},
			y:      &testeditionspb.TestAllTypes{OptionalNestedEnum: testeditionspb.TestAllTypes_BAR.Enum()},
			isLess: true,
		},
		{
			x:      &testeditionspb.TestAllTypes{OptionalNestedEnum: testeditionspb.TestAllTypes_BAR.Enum()},
			y:      &testeditionspb.TestAllTypes{OptionalNestedEnum: testeditionspb.TestAllTypes_FOO.Enum()},
			isLess: false,
		},
		{
			x:      &testeditionspb.TestAllTypes{OptionalInt32: proto.Int32(2)},
			y:      &testeditionspb.TestAllTypes{OptionalInt32: proto.Int32(2)},
			isLess: false,
		},
		{
			x:      &testeditionspb.TestAllTypes{OptionalInt64: proto.Int64(2)},
			y:      &testeditionspb.TestAllTypes{OptionalInt64: proto.Int64(2)},
			isLess: false,
		},
		{
			x:      &testeditionspb.TestAllTypes{OptionalUint32: proto.Uint32(2)},
			y:      &testeditionspb.TestAllTypes{OptionalUint32: proto.Uint32(2)},
			isLess: false,
		},
		{
			x:      &testeditionspb.TestAllTypes{OptionalUint64: proto.Uint64(2)},
			y:      &testeditionspb.TestAllTypes{OptionalUint64: proto.Uint64(2)},
			isLess: false,
		},
		{
			x:      &testeditionspb.TestAllTypes{OptionalSint32: proto.Int32(2)},
			y:      &testeditionspb.TestAllTypes{OptionalSint32: proto.Int32(2)},
			isLess: false,
		},
		{
			x:      &testeditionspb.TestAllTypes{OptionalSint64: proto.Int64(2)},
			y:      &testeditionspb.TestAllTypes{OptionalSint64: proto.Int64(2)},
			isLess: false,
		},
		{
			x:      &testeditionspb.TestAllTypes{OptionalFixed32: proto.Uint32(2)},
			y:      &testeditionspb.TestAllTypes{OptionalFixed32: proto.Uint32(2)},
			isLess: false,
		},
		{
			x:      &testeditionspb.TestAllTypes{OptionalFixed64: proto.Uint64(2)},
			y:      &testeditionspb.TestAllTypes{OptionalFixed64: proto.Uint64(2)},
			isLess: false,
		},
		{
			x:      &testeditionspb.TestAllTypes{OptionalSfixed32: proto.Int32(2)},
			y:      &testeditionspb.TestAllTypes{OptionalSfixed32: proto.Int32(2)},
			isLess: false,
		},
		{
			x:      &testeditionspb.TestAllTypes{OptionalSfixed64: proto.Int64(2)},
			y:      &testeditionspb.TestAllTypes{OptionalSfixed64: proto.Int64(2)},
			isLess: false,
		},
		{
			x:      &testeditionspb.TestAllTypes{OptionalFloat: proto.Float32(2)},
			y:      &testeditionspb.TestAllTypes{OptionalFloat: proto.Float32(2)},
			isLess: false,
		},
		{
			x:      &testeditionspb.TestAllTypes{OptionalDouble: proto.Float64(2)},
			y:      &testeditionspb.TestAllTypes{OptionalDouble: proto.Float64(2)},
			isLess: false,
		},
		{
			x:      &testeditionspb.TestAllTypes{OptionalFloat: proto.Float32(float32(math.NaN()))},
			y:      &testeditionspb.TestAllTypes{OptionalFloat: proto.Float32(float32(math.NaN()))},
			isLess: false,
		},
		{
			x:      &testeditionspb.TestAllTypes{OptionalDouble: proto.Float64(float64(math.NaN()))},
			y:      &testeditionspb.TestAllTypes{OptionalDouble: proto.Float64(float64(math.NaN()))},
			isLess: false,
		},
		{
			x:      &testeditionspb.TestAllTypes{OptionalBool: proto.Bool(true)},
			y:      &testeditionspb.TestAllTypes{OptionalBool: proto.Bool(true)},
			isLess: false,
		},
		{
			x:      &testeditionspb.TestAllTypes{OptionalString: proto.String("abc")},
			y:      &testeditionspb.TestAllTypes{OptionalString: proto.String("abc")},
			isLess: false,
		},
		{
			x:      &testeditionspb.TestAllTypes{OptionalBytes: []byte("abc")},
			y:      &testeditionspb.TestAllTypes{OptionalBytes: []byte("abc")},
			isLess: false,
		},
		{
			x:      &testeditionspb.TestAllTypes{OptionalNestedEnum: testeditionspb.TestAllTypes_FOO.Enum()},
			y:      &testeditionspb.TestAllTypes{OptionalNestedEnum: testeditionspb.TestAllTypes_FOO.Enum()},
			isLess: false,
		},

		// Proto2 presence.
		{
			x:      &testpb.TestAllTypes{},
			y:      &testpb.TestAllTypes{OptionalInt32: proto.Int32(0)},
			isLess: true,
		},
		{
			x:      &testpb.TestAllTypes{OptionalInt32: proto.Int32(0)},
			y:      &testpb.TestAllTypes{},
			isLess: false,
		},
		{
			x:      &testpb.TestAllTypes{},
			y:      &testpb.TestAllTypes{OptionalInt64: proto.Int64(0)},
			isLess: true,
		},
		{
			x:      &testpb.TestAllTypes{OptionalInt64: proto.Int64(0)},
			y:      &testpb.TestAllTypes{},
			isLess: false,
		},
		{
			x:      &testpb.TestAllTypes{},
			y:      &testpb.TestAllTypes{OptionalUint32: proto.Uint32(0)},
			isLess: true,
		},
		{
			x:      &testpb.TestAllTypes{OptionalUint32: proto.Uint32(0)},
			y:      &testpb.TestAllTypes{},
			isLess: false,
		},
		{
			x:      &testpb.TestAllTypes{},
			y:      &testpb.TestAllTypes{OptionalUint64: proto.Uint64(0)},
			isLess: true,
		},
		{
			x:      &testpb.TestAllTypes{OptionalUint64: proto.Uint64(0)},
			y:      &testpb.TestAllTypes{},
			isLess: false,
		},
		{
			x:      &testpb.TestAllTypes{},
			y:      &testpb.TestAllTypes{OptionalSint32: proto.Int32(0)},
			isLess: true,
		},
		{
			x:      &testpb.TestAllTypes{OptionalSint32: proto.Int32(0)},
			y:      &testpb.TestAllTypes{},
			isLess: false,
		},
		{
			x:      &testpb.TestAllTypes{},
			y:      &testpb.TestAllTypes{OptionalSint64: proto.Int64(0)},
			isLess: true,
		},
		{
			x:      &testpb.TestAllTypes{OptionalSint64: proto.Int64(0)},
			y:      &testpb.TestAllTypes{},
			isLess: false,
		},
		{
			x:      &testpb.TestAllTypes{},
			y:      &testpb.TestAllTypes{OptionalFixed32: proto.Uint32(0)},
			isLess: true,
		},
		{
			x:      &testpb.TestAllTypes{OptionalFixed32: proto.Uint32(0)},
			y:      &testpb.TestAllTypes{},
			isLess: false,
		},
		{
			x:      &testpb.TestAllTypes{},
			y:      &testpb.TestAllTypes{OptionalFixed64: proto.Uint64(0)},
			isLess: true,
		},
		{
			x:      &testpb.TestAllTypes{OptionalFixed64: proto.Uint64(0)},
			y:      &testpb.TestAllTypes{},
			isLess: false,
		},
		{
			x:      &testpb.TestAllTypes{},
			y:      &testpb.TestAllTypes{OptionalSfixed32: proto.Int32(0)},
			isLess: true,
		},
		{
			x:      &testpb.TestAllTypes{OptionalSfixed32: proto.Int32(0)},
			y:      &testpb.TestAllTypes{},
			isLess: false,
		},
		{
			x:      &testpb.TestAllTypes{},
			y:      &testpb.TestAllTypes{OptionalSfixed64: proto.Int64(0)},
			isLess: true,
		},
		{
			x:      &testpb.TestAllTypes{OptionalSfixed64: proto.Int64(0)},
			y:      &testpb.TestAllTypes{},
			isLess: false,
		},
		{
			x:      &testpb.TestAllTypes{},
			y:      &testpb.TestAllTypes{OptionalFloat: proto.Float32(0)},
			isLess: true,
		},
		{
			x:      &testpb.TestAllTypes{OptionalFloat: proto.Float32(0)},
			y:      &testpb.TestAllTypes{},
			isLess: false,
		},
		{
			x:      &testpb.TestAllTypes{},
			y:      &testpb.TestAllTypes{OptionalDouble: proto.Float64(0)},
			isLess: true,
		},
		{
			x:      &testpb.TestAllTypes{OptionalDouble: proto.Float64(0)},
			y:      &testpb.TestAllTypes{},
			isLess: false,
		},
		{
			x:      &testpb.TestAllTypes{},
			y:      &testpb.TestAllTypes{OptionalBool: proto.Bool(false)},
			isLess: true,
		},
		{
			x:      &testpb.TestAllTypes{OptionalBool: proto.Bool(false)},
			y:      &testpb.TestAllTypes{},
			isLess: false,
		},
		{
			x:      &testpb.TestAllTypes{},
			y:      &testpb.TestAllTypes{OptionalString: proto.String("")},
			isLess: true,
		},
		{
			x:      &testpb.TestAllTypes{OptionalString: proto.String("")},
			y:      &testpb.TestAllTypes{},
			isLess: false,
		},
		{
			x:      &testpb.TestAllTypes{},
			y:      &testpb.TestAllTypes{OptionalBytes: []byte{}},
			isLess: true,
		},
		{
			x:      &testpb.TestAllTypes{OptionalBytes: []byte{}},
			y:      &testpb.TestAllTypes{},
			isLess: false,
		},
		{
			x:      &testpb.TestAllTypes{},
			y:      &testpb.TestAllTypes{OptionalNestedEnum: testpb.TestAllTypes_FOO.Enum()},
			isLess: true,
		},
		{
			x:      &testpb.TestAllTypes{OptionalNestedEnum: testpb.TestAllTypes_FOO.Enum()},
			y:      &testpb.TestAllTypes{},
			isLess: false,
		},

		// Editions presence.
		{
			x:      &testeditionspb.TestAllTypes{},
			y:      &testeditionspb.TestAllTypes{OptionalInt32: proto.Int32(0)},
			isLess: true,
		},
		{
			x:      &testeditionspb.TestAllTypes{OptionalInt32: proto.Int32(0)},
			y:      &testeditionspb.TestAllTypes{},
			isLess: false,
		},
		{
			x:      &testeditionspb.TestAllTypes{},
			y:      &testeditionspb.TestAllTypes{OptionalInt64: proto.Int64(0)},
			isLess: true,
		},
		{
			x:      &testeditionspb.TestAllTypes{OptionalInt64: proto.Int64(0)},
			y:      &testeditionspb.TestAllTypes{},
			isLess: false,
		},
		{
			x:      &testeditionspb.TestAllTypes{},
			y:      &testeditionspb.TestAllTypes{OptionalUint32: proto.Uint32(0)},
			isLess: true,
		},
		{
			x:      &testeditionspb.TestAllTypes{OptionalUint32: proto.Uint32(0)},
			y:      &testeditionspb.TestAllTypes{},
			isLess: false,
		},
		{
			x:      &testeditionspb.TestAllTypes{},
			y:      &testeditionspb.TestAllTypes{OptionalUint64: proto.Uint64(0)},
			isLess: true,
		},
		{
			x:      &testeditionspb.TestAllTypes{OptionalUint64: proto.Uint64(0)},
			y:      &testeditionspb.TestAllTypes{},
			isLess: false,
		},
		{
			x:      &testeditionspb.TestAllTypes{},
			y:      &testeditionspb.TestAllTypes{OptionalSint32: proto.Int32(0)},
			isLess: true,
		},
		{
			x:      &testeditionspb.TestAllTypes{OptionalSint32: proto.Int32(0)},
			y:      &testeditionspb.TestAllTypes{},
			isLess: false,
		},
		{
			x:      &testeditionspb.TestAllTypes{},
			y:      &testeditionspb.TestAllTypes{OptionalSint64: proto.Int64(0)},
			isLess: true,
		},
		{
			x:      &testeditionspb.TestAllTypes{OptionalSint64: proto.Int64(0)},
			y:      &testeditionspb.TestAllTypes{},
			isLess: false,
		},
		{
			x:      &testeditionspb.TestAllTypes{},
			y:      &testeditionspb.TestAllTypes{OptionalFixed32: proto.Uint32(0)},
			isLess: true,
		},
		{
			x:      &testeditionspb.TestAllTypes{OptionalFixed32: proto.Uint32(0)},
			y:      &testeditionspb.TestAllTypes{},
			isLess: false,
		},
		{
			x:      &testeditionspb.TestAllTypes{},
			y:      &testeditionspb.TestAllTypes{OptionalFixed64: proto.Uint64(0)},
			isLess: true,
		},
		{
			x:      &testeditionspb.TestAllTypes{OptionalFixed64: proto.Uint64(0)},
			y:      &testeditionspb.TestAllTypes{},
			isLess: false,
		},
		{
			x:      &testeditionspb.TestAllTypes{},
			y:      &testeditionspb.TestAllTypes{OptionalSfixed32: proto.Int32(0)},
			isLess: true,
		},
		{
			x:      &testeditionspb.TestAllTypes{OptionalSfixed32: proto.Int32(0)},
			y:      &testeditionspb.TestAllTypes{},
			isLess: false,
		},
		{
			x:      &testeditionspb.TestAllTypes{},
			y:      &testeditionspb.TestAllTypes{OptionalSfixed64: proto.Int64(0)},
			isLess: true,
		},
		{
			x:      &testeditionspb.TestAllTypes{OptionalSfixed64: proto.Int64(0)},
			y:      &testeditionspb.TestAllTypes{},
			isLess: false,
		},
		{
			x:      &testeditionspb.TestAllTypes{},
			y:      &testeditionspb.TestAllTypes{OptionalFloat: proto.Float32(0)},
			isLess: true,
		},
		{
			x:      &testeditionspb.TestAllTypes{OptionalFloat: proto.Float32(0)},
			y:      &testeditionspb.TestAllTypes{},
			isLess: false,
		},
		{
			x:      &testeditionspb.TestAllTypes{},
			y:      &testeditionspb.TestAllTypes{OptionalDouble: proto.Float64(0)},
			isLess: true,
		},
		{
			x:      &testeditionspb.TestAllTypes{OptionalDouble: proto.Float64(0)},
			y:      &testeditionspb.TestAllTypes{},
			isLess: false,
		},
		{
			x:      &testeditionspb.TestAllTypes{},
			y:      &testeditionspb.TestAllTypes{OptionalBool: proto.Bool(false)},
			isLess: true,
		},
		{
			x:      &testeditionspb.TestAllTypes{OptionalBool: proto.Bool(false)},
			y:      &testeditionspb.TestAllTypes{},
			isLess: false,
		},
		{
			x:      &testeditionspb.TestAllTypes{},
			y:      &testeditionspb.TestAllTypes{OptionalString: proto.String("")},
			isLess: true,
		},
		{
			x:      &testeditionspb.TestAllTypes{OptionalString: proto.String("")},
			y:      &testeditionspb.TestAllTypes{},
			isLess: false,
		},
		{
			x:      &testeditionspb.TestAllTypes{},
			y:      &testeditionspb.TestAllTypes{OptionalBytes: []byte{}},
			isLess: true,
		},
		{
			x:      &testeditionspb.TestAllTypes{OptionalBytes: []byte{}},
			y:      &testeditionspb.TestAllTypes{},
			isLess: false,
		},
		{
			x:      &testeditionspb.TestAllTypes{},
			y:      &testeditionspb.TestAllTypes{OptionalNestedEnum: testeditionspb.TestAllTypes_FOO.Enum()},
			isLess: true,
		},
		{
			x:      &testeditionspb.TestAllTypes{OptionalNestedEnum: testeditionspb.TestAllTypes_FOO.Enum()},
			y:      &testeditionspb.TestAllTypes{},
			isLess: false,
		},

		// Proto3 presence.
		{
			x:      &test3pb.TestAllTypes{},
			y:      &test3pb.TestAllTypes{OptionalInt32: proto.Int32(0)},
			isLess: true,
		},
		{
			x:      &test3pb.TestAllTypes{OptionalInt32: proto.Int32(0)},
			y:      &test3pb.TestAllTypes{},
			isLess: false,
		},
		{
			x:      &test3pb.TestAllTypes{},
			y:      &test3pb.TestAllTypes{OptionalInt64: proto.Int64(0)},
			isLess: true,
		},
		{
			x:      &test3pb.TestAllTypes{OptionalInt64: proto.Int64(0)},
			y:      &test3pb.TestAllTypes{},
			isLess: false,
		},
		{
			x:      &test3pb.TestAllTypes{},
			y:      &test3pb.TestAllTypes{OptionalUint32: proto.Uint32(0)},
			isLess: true,
		},
		{
			x:      &test3pb.TestAllTypes{OptionalUint32: proto.Uint32(0)},
			y:      &test3pb.TestAllTypes{},
			isLess: false,
		},
		{
			x:      &test3pb.TestAllTypes{},
			y:      &test3pb.TestAllTypes{OptionalUint64: proto.Uint64(0)},
			isLess: true,
		},
		{
			x:      &test3pb.TestAllTypes{OptionalUint64: proto.Uint64(0)},
			y:      &test3pb.TestAllTypes{},
			isLess: false,
		},
		{
			x:      &test3pb.TestAllTypes{},
			y:      &test3pb.TestAllTypes{OptionalSint32: proto.Int32(0)},
			isLess: true,
		},
		{
			x:      &test3pb.TestAllTypes{OptionalSint32: proto.Int32(0)},
			y:      &test3pb.TestAllTypes{},
			isLess: false,
		},
		{
			x:      &test3pb.TestAllTypes{},
			y:      &test3pb.TestAllTypes{OptionalSint64: proto.Int64(0)},
			isLess: true,
		},
		{
			x:      &test3pb.TestAllTypes{OptionalSint64: proto.Int64(0)},
			y:      &test3pb.TestAllTypes{},
			isLess: false,
		},
		{
			x:      &test3pb.TestAllTypes{},
			y:      &test3pb.TestAllTypes{OptionalFixed32: proto.Uint32(0)},
			isLess: true,
		},
		{
			x:      &test3pb.TestAllTypes{OptionalFixed32: proto.Uint32(0)},
			y:      &test3pb.TestAllTypes{},
			isLess: false,
		},
		{
			x:      &test3pb.TestAllTypes{},
			y:      &test3pb.TestAllTypes{OptionalFixed64: proto.Uint64(0)},
			isLess: true,
		},
		{
			x:      &test3pb.TestAllTypes{OptionalFixed64: proto.Uint64(0)},
			y:      &test3pb.TestAllTypes{},
			isLess: false,
		},
		{
			x:      &test3pb.TestAllTypes{},
			y:      &test3pb.TestAllTypes{OptionalSfixed32: proto.Int32(0)},
			isLess: true,
		},
		{
			x:      &test3pb.TestAllTypes{OptionalSfixed32: proto.Int32(0)},
			y:      &test3pb.TestAllTypes{},
			isLess: false,
		},
		{
			x:      &test3pb.TestAllTypes{},
			y:      &test3pb.TestAllTypes{OptionalSfixed64: proto.Int64(0)},
			isLess: true,
		},
		{
			x:      &test3pb.TestAllTypes{OptionalSfixed64: proto.Int64(0)},
			y:      &test3pb.TestAllTypes{},
			isLess: false,
		},
		{
			x:      &test3pb.TestAllTypes{},
			y:      &test3pb.TestAllTypes{OptionalFloat: proto.Float32(0)},
			isLess: true,
		},
		{
			x:      &test3pb.TestAllTypes{OptionalFloat: proto.Float32(0)},
			y:      &test3pb.TestAllTypes{},
			isLess: false,
		},
		{
			x:      &test3pb.TestAllTypes{},
			y:      &test3pb.TestAllTypes{OptionalDouble: proto.Float64(0)},
			isLess: true,
		},
		{
			x:      &test3pb.TestAllTypes{OptionalDouble: proto.Float64(0)},
			y:      &test3pb.TestAllTypes{},
			isLess: false,
		},
		{
			x:      &test3pb.TestAllTypes{},
			y:      &test3pb.TestAllTypes{OptionalBool: proto.Bool(false)},
			isLess: true,
		},
		{
			x:      &test3pb.TestAllTypes{OptionalBool: proto.Bool(false)},
			y:      &test3pb.TestAllTypes{},
			isLess: false,
		},
		{
			x:      &test3pb.TestAllTypes{},
			y:      &test3pb.TestAllTypes{OptionalString: proto.String("")},
			isLess: true,
		},
		{
			x:      &test3pb.TestAllTypes{OptionalString: proto.String("")},
			y:      &test3pb.TestAllTypes{},
			isLess: false,
		},
		{
			x:      &test3pb.TestAllTypes{},
			y:      &test3pb.TestAllTypes{OptionalBytes: []byte{}},
			isLess: true,
		},
		{
			x:      &test3pb.TestAllTypes{OptionalBytes: []byte{}},
			y:      &test3pb.TestAllTypes{},
			isLess: false,
		},
		{
			x:      &test3pb.TestAllTypes{},
			y:      &test3pb.TestAllTypes{OptionalNestedEnum: test3pb.TestAllTypes_FOO.Enum()},
			isLess: true,
		},
		{
			x:      &test3pb.TestAllTypes{OptionalNestedEnum: test3pb.TestAllTypes_FOO.Enum()},
			y:      &test3pb.TestAllTypes{},
			isLess: false,
		},

		// Proto2 default values are not considered by Equal, so the following are still unequal.
		{
			x:      &testpb.TestAllTypes{},
			y:      &testpb.TestAllTypes{DefaultInt32: proto.Int32(81)},
			isLess: true,
		},
		{
			x:      &testpb.TestAllTypes{DefaultInt32: proto.Int32(81)},
			y:      &testpb.TestAllTypes{},
			isLess: false,
		},
		{
			x:      &testpb.TestAllTypes{},
			y:      &testpb.TestAllTypes{DefaultInt64: proto.Int64(82)},
			isLess: true,
		},
		{
			x:      &testpb.TestAllTypes{DefaultInt64: proto.Int64(82)},
			y:      &testpb.TestAllTypes{},
			isLess: false,
		},
		{
			x:      &testpb.TestAllTypes{},
			y:      &testpb.TestAllTypes{DefaultUint32: proto.Uint32(83)},
			isLess: true,
		},
		{
			x:      &testpb.TestAllTypes{DefaultUint32: proto.Uint32(83)},
			y:      &testpb.TestAllTypes{},
			isLess: false,
		},
		{
			x:      &testpb.TestAllTypes{},
			y:      &testpb.TestAllTypes{DefaultUint64: proto.Uint64(84)},
			isLess: true,
		},
		{
			x:      &testpb.TestAllTypes{DefaultUint64: proto.Uint64(84)},
			y:      &testpb.TestAllTypes{},
			isLess: false,
		},
		{
			x:      &testpb.TestAllTypes{},
			y:      &testpb.TestAllTypes{DefaultSint32: proto.Int32(-85)},
			isLess: true,
		},
		{
			x:      &testpb.TestAllTypes{DefaultSint32: proto.Int32(-85)},
			y:      &testpb.TestAllTypes{},
			isLess: false,
		},
		{
			x:      &testpb.TestAllTypes{},
			y:      &testpb.TestAllTypes{DefaultSint64: proto.Int64(86)},
			isLess: true,
		},
		{
			x:      &testpb.TestAllTypes{DefaultSint64: proto.Int64(86)},
			y:      &testpb.TestAllTypes{},
			isLess: false,
		},
		{
			x:      &testpb.TestAllTypes{},
			y:      &testpb.TestAllTypes{DefaultFixed32: proto.Uint32(87)},
			isLess: true,
		},
		{
			x:      &testpb.TestAllTypes{DefaultFixed32: proto.Uint32(87)},
			y:      &testpb.TestAllTypes{},
			isLess: false,
		},
		{
			x:      &testpb.TestAllTypes{},
			y:      &testpb.TestAllTypes{DefaultFixed64: proto.Uint64(88)},
			isLess: true,
		},
		{
			x:      &testpb.TestAllTypes{DefaultFixed64: proto.Uint64(88)},
			y:      &testpb.TestAllTypes{},
			isLess: false,
		},
		{
			x:      &testpb.TestAllTypes{},
			y:      &testpb.TestAllTypes{DefaultSfixed32: proto.Int32(89)},
			isLess: true,
		},
		{
			x:      &testpb.TestAllTypes{DefaultSfixed32: proto.Int32(89)},
			y:      &testpb.TestAllTypes{},
			isLess: false,
		},
		{
			x:      &testpb.TestAllTypes{},
			y:      &testpb.TestAllTypes{DefaultSfixed64: proto.Int64(-90)},
			isLess: true,
		},
		{
			x:      &testpb.TestAllTypes{DefaultSfixed64: proto.Int64(-90)},
			y:      &testpb.TestAllTypes{},
			isLess: false,
		},
		{
			x:      &testpb.TestAllTypes{},
			y:      &testpb.TestAllTypes{DefaultFloat: proto.Float32(91.5)},
			isLess: true,
		},
		{
			x:      &testpb.TestAllTypes{DefaultFloat: proto.Float32(91.5)},
			y:      &testpb.TestAllTypes{},
			isLess: false,
		},
		{
			x:      &testpb.TestAllTypes{},
			y:      &testpb.TestAllTypes{DefaultDouble: proto.Float64(92e3)},
			isLess: true,
		},
		{
			x:      &testpb.TestAllTypes{DefaultDouble: proto.Float64(92e3)},
			y:      &testpb.TestAllTypes{},
			isLess: false,
		},
		{
			x:      &testpb.TestAllTypes{},
			y:      &testpb.TestAllTypes{DefaultBool: proto.Bool(true)},
			isLess: true,
		},
		{
			x:      &testpb.TestAllTypes{DefaultBool: proto.Bool(true)},
			y:      &testpb.TestAllTypes{},
			isLess: false,
		},
		{
			x:      &testpb.TestAllTypes{},
			y:      &testpb.TestAllTypes{DefaultString: proto.String("hello")},
			isLess: true,
		},
		{
			x:      &testpb.TestAllTypes{DefaultString: proto.String("hello")},
			y:      &testpb.TestAllTypes{},
			isLess: false,
		},
		{
			x:      &testpb.TestAllTypes{},
			y:      &testpb.TestAllTypes{DefaultBytes: []byte("world")},
			isLess: true,
		},
		{
			x:      &testpb.TestAllTypes{DefaultBytes: []byte("world")},
			y:      &testpb.TestAllTypes{},
			isLess: false,
		},
		{
			x:      &testpb.TestAllTypes{},
			y:      &testpb.TestAllTypes{DefaultNestedEnum: testpb.TestAllTypes_BAR.Enum()},
			isLess: true,
		},
		{
			x:      &testpb.TestAllTypes{DefaultNestedEnum: testpb.TestAllTypes_BAR.Enum()},
			y:      &testpb.TestAllTypes{},
			isLess: false,
		},

		// Edition default values are not considered by Equal, so the following are still unequal.
		{
			x:      &testeditionspb.TestAllTypes{},
			y:      &testeditionspb.TestAllTypes{DefaultInt32: proto.Int32(81)},
			isLess: true,
		},
		{
			x:      &testeditionspb.TestAllTypes{DefaultInt32: proto.Int32(81)},
			y:      &testeditionspb.TestAllTypes{},
			isLess: false,
		},
		{
			x:      &testeditionspb.TestAllTypes{},
			y:      &testeditionspb.TestAllTypes{DefaultInt64: proto.Int64(82)},
			isLess: true,
		},
		{
			x:      &testeditionspb.TestAllTypes{DefaultInt64: proto.Int64(82)},
			y:      &testeditionspb.TestAllTypes{},
			isLess: false,
		},
		{
			x:      &testeditionspb.TestAllTypes{},
			y:      &testeditionspb.TestAllTypes{DefaultUint32: proto.Uint32(83)},
			isLess: true,
		},
		{
			x:      &testeditionspb.TestAllTypes{DefaultUint32: proto.Uint32(83)},
			y:      &testeditionspb.TestAllTypes{},
			isLess: false,
		},
		{
			x:      &testeditionspb.TestAllTypes{},
			y:      &testeditionspb.TestAllTypes{DefaultUint64: proto.Uint64(84)},
			isLess: true,
		},
		{
			x:      &testeditionspb.TestAllTypes{DefaultUint64: proto.Uint64(84)},
			y:      &testeditionspb.TestAllTypes{},
			isLess: false,
		},
		{
			x:      &testeditionspb.TestAllTypes{},
			y:      &testeditionspb.TestAllTypes{DefaultSint32: proto.Int32(-85)},
			isLess: true,
		},
		{
			x:      &testeditionspb.TestAllTypes{DefaultSint32: proto.Int32(-85)},
			y:      &testeditionspb.TestAllTypes{},
			isLess: false,
		},
		{
			x:      &testeditionspb.TestAllTypes{},
			y:      &testeditionspb.TestAllTypes{DefaultSint64: proto.Int64(86)},
			isLess: true,
		},
		{
			x:      &testeditionspb.TestAllTypes{DefaultSint64: proto.Int64(86)},
			y:      &testeditionspb.TestAllTypes{},
			isLess: false,
		},
		{
			x:      &testeditionspb.TestAllTypes{},
			y:      &testeditionspb.TestAllTypes{DefaultFixed32: proto.Uint32(87)},
			isLess: true,
		},
		{
			x:      &testeditionspb.TestAllTypes{DefaultFixed32: proto.Uint32(87)},
			y:      &testeditionspb.TestAllTypes{},
			isLess: false,
		},
		{
			x:      &testeditionspb.TestAllTypes{},
			y:      &testeditionspb.TestAllTypes{DefaultFixed64: proto.Uint64(88)},
			isLess: true,
		},
		{
			x:      &testeditionspb.TestAllTypes{DefaultFixed64: proto.Uint64(88)},
			y:      &testeditionspb.TestAllTypes{},
			isLess: false,
		},
		{
			x:      &testeditionspb.TestAllTypes{},
			y:      &testeditionspb.TestAllTypes{DefaultSfixed32: proto.Int32(89)},
			isLess: true,
		},
		{
			x:      &testeditionspb.TestAllTypes{DefaultSfixed32: proto.Int32(89)},
			y:      &testeditionspb.TestAllTypes{},
			isLess: false,
		},
		{
			x:      &testeditionspb.TestAllTypes{},
			y:      &testeditionspb.TestAllTypes{DefaultSfixed64: proto.Int64(-90)},
			isLess: true,
		},
		{
			x:      &testeditionspb.TestAllTypes{DefaultSfixed64: proto.Int64(-90)},
			y:      &testeditionspb.TestAllTypes{},
			isLess: false,
		},
		{
			x:      &testeditionspb.TestAllTypes{},
			y:      &testeditionspb.TestAllTypes{DefaultFloat: proto.Float32(91.5)},
			isLess: true,
		},
		{
			x:      &testeditionspb.TestAllTypes{DefaultFloat: proto.Float32(91.5)},
			y:      &testeditionspb.TestAllTypes{},
			isLess: false,
		},
		{
			x:      &testeditionspb.TestAllTypes{},
			y:      &testeditionspb.TestAllTypes{DefaultDouble: proto.Float64(92e3)},
			isLess: true,
		},
		{
			x:      &testeditionspb.TestAllTypes{DefaultDouble: proto.Float64(92e3)},
			y:      &testeditionspb.TestAllTypes{},
			isLess: false,
		},
		{
			x:      &testeditionspb.TestAllTypes{},
			y:      &testeditionspb.TestAllTypes{DefaultBool: proto.Bool(true)},
			isLess: true,
		},
		{
			x:      &testeditionspb.TestAllTypes{DefaultBool: proto.Bool(true)},
			y:      &testeditionspb.TestAllTypes{},
			isLess: false,
		},
		{
			x:      &testeditionspb.TestAllTypes{},
			y:      &testeditionspb.TestAllTypes{DefaultString: proto.String("hello")},
			isLess: true,
		},
		{
			x:      &testeditionspb.TestAllTypes{DefaultString: proto.String("hello")},
			y:      &testeditionspb.TestAllTypes{},
			isLess: false,
		},
		{
			x:      &testeditionspb.TestAllTypes{},
			y:      &testeditionspb.TestAllTypes{DefaultBytes: []byte("world")},
			isLess: true,
		},
		{
			x:      &testeditionspb.TestAllTypes{DefaultBytes: []byte("world")},
			y:      &testeditionspb.TestAllTypes{},
			isLess: false,
		},
		{
			x:      &testeditionspb.TestAllTypes{},
			y:      &testeditionspb.TestAllTypes{DefaultNestedEnum: testeditionspb.TestAllTypes_BAR.Enum()},
			isLess: true,
		},
		{
			x:      &testeditionspb.TestAllTypes{DefaultNestedEnum: testeditionspb.TestAllTypes_BAR.Enum()},
			y:      &testeditionspb.TestAllTypes{},
			isLess: false,
		},

		// Groups.
		{
			x: &testpb.TestAllTypes{Optionalgroup: &testpb.TestAllTypes_OptionalGroup{
				A: proto.Int32(1),
			}},
			y: &testpb.TestAllTypes{Optionalgroup: &testpb.TestAllTypes_OptionalGroup{
				A: proto.Int32(2),
			}},
			isLess: true,
		},
		{
			x: &testpb.TestAllTypes{Optionalgroup: &testpb.TestAllTypes_OptionalGroup{
				A: proto.Int32(2),
			}},
			y: &testpb.TestAllTypes{Optionalgroup: &testpb.TestAllTypes_OptionalGroup{
				A: proto.Int32(1),
			}},
			isLess: false,
		},
		{
			x:      &testpb.TestAllTypes{},
			y:      &testpb.TestAllTypes{Optionalgroup: &testpb.TestAllTypes_OptionalGroup{}},
			isLess: true,
		},
		{
			x:      &testpb.TestAllTypes{Optionalgroup: &testpb.TestAllTypes_OptionalGroup{}},
			y:      &testpb.TestAllTypes{},
			isLess: false,
		},
		{
			x: &testeditionspb.TestAllTypes{Optionalgroup: &testeditionspb.TestAllTypes_OptionalGroup{
				A: proto.Int32(1),
			}},
			y: &testeditionspb.TestAllTypes{Optionalgroup: &testeditionspb.TestAllTypes_OptionalGroup{
				A: proto.Int32(2),
			}},
			isLess: true,
		},
		{
			x: &testeditionspb.TestAllTypes{Optionalgroup: &testeditionspb.TestAllTypes_OptionalGroup{
				A: proto.Int32(2),
			}},
			y: &testeditionspb.TestAllTypes{Optionalgroup: &testeditionspb.TestAllTypes_OptionalGroup{
				A: proto.Int32(1),
			}},
			isLess: false,
		},
		{
			x:      &testeditionspb.TestAllTypes{},
			y:      &testeditionspb.TestAllTypes{Optionalgroup: &testeditionspb.TestAllTypes_OptionalGroup{}},
			isLess: true,
		},
		{
			x:      &testeditionspb.TestAllTypes{Optionalgroup: &testeditionspb.TestAllTypes_OptionalGroup{}},
			y:      &testeditionspb.TestAllTypes{},
			isLess: false,
		},

		// Messages.
		{
			x: &testpb.TestAllTypes{OptionalNestedMessage: &testpb.TestAllTypes_NestedMessage{
				A: proto.Int32(1),
			}},
			y: &testpb.TestAllTypes{OptionalNestedMessage: &testpb.TestAllTypes_NestedMessage{
				A: proto.Int32(2),
			}},
			isLess: true,
		},
		{
			x: &testpb.TestAllTypes{OptionalNestedMessage: &testpb.TestAllTypes_NestedMessage{
				A: proto.Int32(2),
			}},
			y: &testpb.TestAllTypes{OptionalNestedMessage: &testpb.TestAllTypes_NestedMessage{
				A: proto.Int32(1),
			}},
			isLess: false,
		},
		{
			x:      &testpb.TestAllTypes{},
			y:      &testpb.TestAllTypes{OptionalNestedMessage: &testpb.TestAllTypes_NestedMessage{}},
			isLess: true,
		},
		{
			x:      &testpb.TestAllTypes{OptionalNestedMessage: &testpb.TestAllTypes_NestedMessage{}},
			y:      &testpb.TestAllTypes{},
			isLess: false,
		},
		{
			x:      &testeditionspb.TestAllTypes{},
			y:      &testeditionspb.TestAllTypes{OptionalNestedMessage: &testeditionspb.TestAllTypes_NestedMessage{}},
			isLess: true,
		},
		{
			x:      &testeditionspb.TestAllTypes{OptionalNestedMessage: &testeditionspb.TestAllTypes_NestedMessage{}},
			y:      &testeditionspb.TestAllTypes{},
			isLess: false,
		},
		{
			x:      &test3pb.TestAllTypes{},
			y:      &test3pb.TestAllTypes{OptionalNestedMessage: &test3pb.TestAllTypes_NestedMessage{}},
			isLess: true,
		},
		{
			x:      &test3pb.TestAllTypes{OptionalNestedMessage: &test3pb.TestAllTypes_NestedMessage{}},
			y:      &test3pb.TestAllTypes{},
			isLess: false,
		},

		// Lists.
		{
			x:      &testpb.TestAllTypes{RepeatedInt32: []int32{1}},
			y:      &testpb.TestAllTypes{RepeatedInt32: []int32{1}},
			isLess: false,
		},
		{
			x:      &testpb.TestAllTypes{RepeatedInt32: []int32{1}},
			y:      &testpb.TestAllTypes{RepeatedInt32: []int32{1, 2}},
			isLess: true,
		},
		{
			x:      &testpb.TestAllTypes{RepeatedInt32: []int32{1, 2}},
			y:      &testpb.TestAllTypes{RepeatedInt32: []int32{1}},
			isLess: false,
		},
		{
			x:      &testpb.TestAllTypes{RepeatedInt32: []int32{7}},
			y:      &testpb.TestAllTypes{RepeatedInt32: []int32{1, 2}},
			isLess: true,
		},
		{
			x:      &testpb.TestAllTypes{RepeatedInt32: []int32{1, 2}},
			y:      &testpb.TestAllTypes{RepeatedInt32: []int32{7}},
			isLess: false,
		},
		{
			x:      &testpb.TestAllTypes{RepeatedInt32: []int32{1, 2, 3}},
			y:      &testpb.TestAllTypes{RepeatedInt32: []int32{1, 3, 2}},
			isLess: true,
		},
		{
			x:      &testpb.TestAllTypes{RepeatedInt32: []int32{1, 3, 2}},
			y:      &testpb.TestAllTypes{RepeatedInt32: []int32{1, 2, 3}},
			isLess: false,
		},
		{
			x:      &testpb.TestAllTypes{RepeatedInt64: []int64{1, 2, 3}},
			y:      &testpb.TestAllTypes{RepeatedInt64: []int64{1, 3, 2}},
			isLess: true,
		},
		{
			x:      &testpb.TestAllTypes{RepeatedInt64: []int64{1, 3, 2}},
			y:      &testpb.TestAllTypes{RepeatedInt64: []int64{1, 2, 3}},
			isLess: false,
		},
		{
			x:      &testpb.TestAllTypes{RepeatedUint32: []uint32{1, 2, 3}},
			y:      &testpb.TestAllTypes{RepeatedUint32: []uint32{1, 3, 2}},
			isLess: true,
		},
		{
			x:      &testpb.TestAllTypes{RepeatedUint32: []uint32{1, 3, 2}},
			y:      &testpb.TestAllTypes{RepeatedUint32: []uint32{1, 2, 3}},
			isLess: false,
		},
		{
			x:      &testpb.TestAllTypes{RepeatedUint64: []uint64{1, 2, 3}},
			y:      &testpb.TestAllTypes{RepeatedUint64: []uint64{1, 3, 2}},
			isLess: true,
		},
		{
			x:      &testpb.TestAllTypes{RepeatedUint64: []uint64{1, 3, 2}},
			y:      &testpb.TestAllTypes{RepeatedUint64: []uint64{1, 2, 3}},
			isLess: false,
		},
		{
			x:      &testpb.TestAllTypes{RepeatedSint32: []int32{1, 2, 3}},
			y:      &testpb.TestAllTypes{RepeatedSint32: []int32{1, 3, 2}},
			isLess: true,
		},
		{
			x:      &testpb.TestAllTypes{RepeatedSint32: []int32{1, 3, 2}},
			y:      &testpb.TestAllTypes{RepeatedSint32: []int32{1, 2, 3}},
			isLess: false,
		},
		{
			x:      &testpb.TestAllTypes{RepeatedSint64: []int64{1, 2, 3}},
			y:      &testpb.TestAllTypes{RepeatedSint64: []int64{1, 3, 2}},
			isLess: true,
		},
		{
			x:      &testpb.TestAllTypes{RepeatedSint64: []int64{1, 3, 2}},
			y:      &testpb.TestAllTypes{RepeatedSint64: []int64{1, 2, 3}},
			isLess: false,
		},
		{
			x:      &testpb.TestAllTypes{RepeatedFixed32: []uint32{1, 2, 3}},
			y:      &testpb.TestAllTypes{RepeatedFixed32: []uint32{1, 3, 2}},
			isLess: true,
		},
		{
			x:      &testpb.TestAllTypes{RepeatedFixed32: []uint32{1, 3, 2}},
			y:      &testpb.TestAllTypes{RepeatedFixed32: []uint32{1, 2, 3}},
			isLess: false,
		},
		{
			x:      &testpb.TestAllTypes{RepeatedFixed64: []uint64{1, 2, 3}},
			y:      &testpb.TestAllTypes{RepeatedFixed64: []uint64{1, 3, 2}},
			isLess: true,
		},
		{
			x:      &testpb.TestAllTypes{RepeatedFixed64: []uint64{1, 3, 2}},
			y:      &testpb.TestAllTypes{RepeatedFixed64: []uint64{1, 2, 3}},
			isLess: false,
		},
		{
			x:      &testpb.TestAllTypes{RepeatedSfixed32: []int32{1, 2, 3}},
			y:      &testpb.TestAllTypes{RepeatedSfixed32: []int32{1, 3, 2}},
			isLess: true,
		},
		{
			x:      &testpb.TestAllTypes{RepeatedSfixed32: []int32{1, 3, 2}},
			y:      &testpb.TestAllTypes{RepeatedSfixed32: []int32{1, 2, 3}},
			isLess: false,
		},
		{
			x:      &testpb.TestAllTypes{RepeatedSfixed64: []int64{1, 2, 3}},
			y:      &testpb.TestAllTypes{RepeatedSfixed64: []int64{1, 3, 2}},
			isLess: true,
		},
		{
			x:      &testpb.TestAllTypes{RepeatedSfixed64: []int64{1, 3, 2}},
			y:      &testpb.TestAllTypes{RepeatedSfixed64: []int64{1, 2, 3}},
			isLess: false,
		},
		{
			x:      &testpb.TestAllTypes{RepeatedFloat: []float32{1, 2, 3}},
			y:      &testpb.TestAllTypes{RepeatedFloat: []float32{1, 3, 2}},
			isLess: true,
		},
		{
			x:      &testpb.TestAllTypes{RepeatedFloat: []float32{1, 3, 2}},
			y:      &testpb.TestAllTypes{RepeatedFloat: []float32{1, 2, 3}},
			isLess: false,
		},
		{
			x:      &testpb.TestAllTypes{RepeatedDouble: []float64{1, 2, 3}},
			y:      &testpb.TestAllTypes{RepeatedDouble: []float64{1, 3, 2}},
			isLess: true,
		},
		{
			x:      &testpb.TestAllTypes{RepeatedDouble: []float64{1, 3, 2}},
			y:      &testpb.TestAllTypes{RepeatedDouble: []float64{1, 2, 3}},
			isLess: false,
		},
		{
			x:      &testpb.TestAllTypes{RepeatedBool: []bool{true, false}},
			y:      &testpb.TestAllTypes{RepeatedBool: []bool{true, true}},
			isLess: true,
		},
		{
			x:      &testpb.TestAllTypes{RepeatedBool: []bool{true, true}},
			y:      &testpb.TestAllTypes{RepeatedBool: []bool{true, false}},
			isLess: false,
		},
		{
			x:      &testpb.TestAllTypes{RepeatedString: []string{"a", "b"}},
			y:      &testpb.TestAllTypes{RepeatedString: []string{"a", "c"}},
			isLess: true,
		},
		{
			x:      &testpb.TestAllTypes{RepeatedString: []string{"a", "c"}},
			y:      &testpb.TestAllTypes{RepeatedString: []string{"a", "b"}},
			isLess: false,
		},
		{
			x:      &testpb.TestAllTypes{RepeatedBytes: [][]byte{[]byte("a"), []byte("b")}},
			y:      &testpb.TestAllTypes{RepeatedBytes: [][]byte{[]byte("a"), []byte("c")}},
			isLess: true,
		},
		{
			x:      &testpb.TestAllTypes{RepeatedBytes: [][]byte{[]byte("a"), []byte("c")}},
			y:      &testpb.TestAllTypes{RepeatedBytes: [][]byte{[]byte("a"), []byte("b")}},
			isLess: false,
		},
		{
			x:      &testpb.TestAllTypes{RepeatedNestedEnum: []testpb.TestAllTypes_NestedEnum{testpb.TestAllTypes_FOO}},
			y:      &testpb.TestAllTypes{RepeatedNestedEnum: []testpb.TestAllTypes_NestedEnum{testpb.TestAllTypes_BAR}},
			isLess: true,
		},
		{
			x:      &testpb.TestAllTypes{RepeatedNestedEnum: []testpb.TestAllTypes_NestedEnum{testpb.TestAllTypes_BAR}},
			y:      &testpb.TestAllTypes{RepeatedNestedEnum: []testpb.TestAllTypes_NestedEnum{testpb.TestAllTypes_FOO}},
			isLess: false,
		},
		{
			x: &testpb.TestAllTypes{Repeatedgroup: []*testpb.TestAllTypes_RepeatedGroup{
				{
					A: proto.Int32(1),
				}, {
					A: proto.Int32(2),
				},
			}},
			y: &testpb.TestAllTypes{Repeatedgroup: []*testpb.TestAllTypes_RepeatedGroup{
				{
					A: proto.Int32(1),
				}, {
					A: proto.Int32(3),
				},
			}},
			isLess: true,
		},
		{
			x: &testpb.TestAllTypes{Repeatedgroup: []*testpb.TestAllTypes_RepeatedGroup{
				{
					A: proto.Int32(1),
				}, {
					A: proto.Int32(3),
				},
			}},
			y: &testpb.TestAllTypes{Repeatedgroup: []*testpb.TestAllTypes_RepeatedGroup{
				{
					A: proto.Int32(1),
				}, {
					A: proto.Int32(2),
				},
			}},
			isLess: false,
		},
		{
			x: &testpb.TestAllTypes{RepeatedNestedMessage: []*testpb.TestAllTypes_NestedMessage{
				{
					A: proto.Int32(1),
				}, {
					A: proto.Int32(2),
				},
			}},
			y: &testpb.TestAllTypes{RepeatedNestedMessage: []*testpb.TestAllTypes_NestedMessage{
				{
					A: proto.Int32(1),
				}, {
					A: proto.Int32(3),
				},
			}},
			isLess: true,
		},
		{
			x: &testpb.TestAllTypes{RepeatedNestedMessage: []*testpb.TestAllTypes_NestedMessage{
				{
					A: proto.Int32(1),
				}, {
					A: proto.Int32(3),
				},
			}},
			y: &testpb.TestAllTypes{RepeatedNestedMessage: []*testpb.TestAllTypes_NestedMessage{
				{
					A: proto.Int32(1),
				}, {
					A: proto.Int32(2),
				},
			}},
			isLess: false,
		},

		// Editions Lists.
		{
			x:      &testeditionspb.TestAllTypes{RepeatedInt32: []int32{1, 2, 3}},
			y:      &testeditionspb.TestAllTypes{RepeatedInt32: []int32{1, 2, 3}},
			isLess: false,
		},
		{
			x:      &testeditionspb.TestAllTypes{RepeatedInt32: []int32{1}},
			y:      &testeditionspb.TestAllTypes{RepeatedInt32: []int32{1, 2, 3}},
			isLess: true,
		},
		{
			x:      &testeditionspb.TestAllTypes{RepeatedInt32: []int32{1, 2, 3}},
			y:      &testeditionspb.TestAllTypes{RepeatedInt32: []int32{1}},
			isLess: false,
		},
		{
			x:      &testeditionspb.TestAllTypes{RepeatedInt32: []int32{1, 2, 3}},
			y:      &testeditionspb.TestAllTypes{RepeatedInt32: []int32{1, 3, 2}},
			isLess: true,
		},
		{
			x:      &testeditionspb.TestAllTypes{RepeatedInt32: []int32{1, 3, 2}},
			y:      &testeditionspb.TestAllTypes{RepeatedInt32: []int32{1, 2, 3}},
			isLess: false,
		},
		{
			x:      &testeditionspb.TestAllTypes{RepeatedInt64: []int64{1, 2, 3}},
			y:      &testeditionspb.TestAllTypes{RepeatedInt64: []int64{1, 3, 2}},
			isLess: true,
		},
		{
			x:      &testeditionspb.TestAllTypes{RepeatedInt64: []int64{1, 3, 2}},
			y:      &testeditionspb.TestAllTypes{RepeatedInt64: []int64{1, 2, 3}},
			isLess: false,
		},
		{
			x:      &testeditionspb.TestAllTypes{RepeatedUint32: []uint32{1, 2, 3}},
			y:      &testeditionspb.TestAllTypes{RepeatedUint32: []uint32{1, 3, 2}},
			isLess: true,
		},
		{
			x:      &testeditionspb.TestAllTypes{RepeatedUint32: []uint32{1, 3, 2}},
			y:      &testeditionspb.TestAllTypes{RepeatedUint32: []uint32{1, 2, 3}},
			isLess: false,
		},
		{
			x:      &testeditionspb.TestAllTypes{RepeatedUint64: []uint64{1, 2, 3}},
			y:      &testeditionspb.TestAllTypes{RepeatedUint64: []uint64{1, 3, 2}},
			isLess: true,
		},
		{
			x:      &testeditionspb.TestAllTypes{RepeatedUint64: []uint64{1, 3, 2}},
			y:      &testeditionspb.TestAllTypes{RepeatedUint64: []uint64{1, 2, 3}},
			isLess: false,
		},
		{
			x:      &testeditionspb.TestAllTypes{RepeatedSint32: []int32{1, 2, 3}},
			y:      &testeditionspb.TestAllTypes{RepeatedSint32: []int32{1, 3, 2}},
			isLess: true,
		},
		{
			x:      &testeditionspb.TestAllTypes{RepeatedSint32: []int32{1, 3, 2}},
			y:      &testeditionspb.TestAllTypes{RepeatedSint32: []int32{1, 2, 3}},
			isLess: false,
		},
		{
			x:      &testeditionspb.TestAllTypes{RepeatedSint64: []int64{1, 2, 3}},
			y:      &testeditionspb.TestAllTypes{RepeatedSint64: []int64{1, 3, 2}},
			isLess: true,
		},
		{
			x:      &testeditionspb.TestAllTypes{RepeatedSint64: []int64{1, 3, 2}},
			y:      &testeditionspb.TestAllTypes{RepeatedSint64: []int64{1, 2, 3}},
			isLess: false,
		},
		{
			x:      &testeditionspb.TestAllTypes{RepeatedFixed32: []uint32{1, 2, 3}},
			y:      &testeditionspb.TestAllTypes{RepeatedFixed32: []uint32{1, 3, 2}},
			isLess: true,
		},
		{
			x:      &testeditionspb.TestAllTypes{RepeatedFixed32: []uint32{1, 3, 2}},
			y:      &testeditionspb.TestAllTypes{RepeatedFixed32: []uint32{1, 2, 3}},
			isLess: false,
		},
		{
			x:      &testeditionspb.TestAllTypes{RepeatedFixed64: []uint64{1, 2, 3}},
			y:      &testeditionspb.TestAllTypes{RepeatedFixed64: []uint64{1, 3, 2}},
			isLess: true,
		},
		{
			x:      &testeditionspb.TestAllTypes{RepeatedFixed64: []uint64{1, 3, 2}},
			y:      &testeditionspb.TestAllTypes{RepeatedFixed64: []uint64{1, 2, 3}},
			isLess: false,
		},
		{
			x:      &testeditionspb.TestAllTypes{RepeatedSfixed32: []int32{1, 2, 3}},
			y:      &testeditionspb.TestAllTypes{RepeatedSfixed32: []int32{1, 3, 2}},
			isLess: true,
		},
		{
			x:      &testeditionspb.TestAllTypes{RepeatedSfixed32: []int32{1, 3, 2}},
			y:      &testeditionspb.TestAllTypes{RepeatedSfixed32: []int32{1, 2, 3}},
			isLess: false,
		},
		{
			x:      &testeditionspb.TestAllTypes{RepeatedSfixed64: []int64{1, 2, 3}},
			y:      &testeditionspb.TestAllTypes{RepeatedSfixed64: []int64{1, 3, 2}},
			isLess: true,
		},
		{
			x:      &testeditionspb.TestAllTypes{RepeatedSfixed64: []int64{1, 3, 2}},
			y:      &testeditionspb.TestAllTypes{RepeatedSfixed64: []int64{1, 2, 3}},
			isLess: false,
		},
		{
			x:      &testeditionspb.TestAllTypes{RepeatedFloat: []float32{1, 2, 3}},
			y:      &testeditionspb.TestAllTypes{RepeatedFloat: []float32{1, 3, 2}},
			isLess: true,
		},
		{
			x:      &testeditionspb.TestAllTypes{RepeatedFloat: []float32{1, 3, 2}},
			y:      &testeditionspb.TestAllTypes{RepeatedFloat: []float32{1, 2, 3}},
			isLess: false,
		},
		{
			x:      &testeditionspb.TestAllTypes{RepeatedDouble: []float64{1, 2, 3}},
			y:      &testeditionspb.TestAllTypes{RepeatedDouble: []float64{1, 3, 2}},
			isLess: true,
		},
		{
			x:      &testeditionspb.TestAllTypes{RepeatedDouble: []float64{1, 3, 2}},
			y:      &testeditionspb.TestAllTypes{RepeatedDouble: []float64{1, 2, 3}},
			isLess: false,
		},
		{
			x:      &testeditionspb.TestAllTypes{RepeatedBool: []bool{true, false}},
			y:      &testeditionspb.TestAllTypes{RepeatedBool: []bool{true, true}},
			isLess: true,
		},
		{
			x:      &testeditionspb.TestAllTypes{RepeatedBool: []bool{true, true}},
			y:      &testeditionspb.TestAllTypes{RepeatedBool: []bool{true, false}},
			isLess: false,
		},
		{
			x:      &testeditionspb.TestAllTypes{RepeatedString: []string{"a", "b"}},
			y:      &testeditionspb.TestAllTypes{RepeatedString: []string{"a", "c"}},
			isLess: true,
		},
		{
			x:      &testeditionspb.TestAllTypes{RepeatedString: []string{"a", "c"}},
			y:      &testeditionspb.TestAllTypes{RepeatedString: []string{"a", "b"}},
			isLess: false,
		},
		{
			x:      &testeditionspb.TestAllTypes{RepeatedBytes: [][]byte{[]byte("a"), []byte("b")}},
			y:      &testeditionspb.TestAllTypes{RepeatedBytes: [][]byte{[]byte("a"), []byte("c")}},
			isLess: true,
		},
		{
			x:      &testeditionspb.TestAllTypes{RepeatedBytes: [][]byte{[]byte("a"), []byte("c")}},
			y:      &testeditionspb.TestAllTypes{RepeatedBytes: [][]byte{[]byte("a"), []byte("b")}},
			isLess: false,
		},
		{
			x: &testeditionspb.TestAllTypes{
				RepeatedNestedEnum: []testeditionspb.TestAllTypes_NestedEnum{testeditionspb.TestAllTypes_FOO},
			},
			y: &testeditionspb.TestAllTypes{
				RepeatedNestedEnum: []testeditionspb.TestAllTypes_NestedEnum{testeditionspb.TestAllTypes_BAR},
			},
			isLess: true,
		},
		{
			x: &testeditionspb.TestAllTypes{
				RepeatedNestedEnum: []testeditionspb.TestAllTypes_NestedEnum{testeditionspb.TestAllTypes_BAR},
			},
			y: &testeditionspb.TestAllTypes{
				RepeatedNestedEnum: []testeditionspb.TestAllTypes_NestedEnum{testeditionspb.TestAllTypes_FOO},
			},
			isLess: false,
		},
		{
			x: &testeditionspb.TestAllTypes{Repeatedgroup: []*testeditionspb.TestAllTypes_RepeatedGroup{
				{
					A: proto.Int32(1),
				}, {
					A: proto.Int32(2),
				},
			}},
			y: &testeditionspb.TestAllTypes{Repeatedgroup: []*testeditionspb.TestAllTypes_RepeatedGroup{
				{
					A: proto.Int32(1),
				}, {
					A: proto.Int32(3),
				},
			}},
			isLess: true,
		},
		{
			x: &testeditionspb.TestAllTypes{Repeatedgroup: []*testeditionspb.TestAllTypes_RepeatedGroup{
				{
					A: proto.Int32(1),
				}, {
					A: proto.Int32(3),
				},
			}},
			y: &testeditionspb.TestAllTypes{Repeatedgroup: []*testeditionspb.TestAllTypes_RepeatedGroup{
				{
					A: proto.Int32(1),
				}, {
					A: proto.Int32(2),
				},
			}},
			isLess: false,
		},
		{
			x: &testeditionspb.TestAllTypes{RepeatedNestedMessage: []*testeditionspb.TestAllTypes_NestedMessage{
				{
					A: proto.Int32(1),
				}, {
					A: proto.Int32(2),
				},
			}},
			y: &testeditionspb.TestAllTypes{RepeatedNestedMessage: []*testeditionspb.TestAllTypes_NestedMessage{
				{
					A: proto.Int32(1),
				}, {
					A: proto.Int32(3),
				},
			}},
			isLess: true,
		},
		{
			x: &testeditionspb.TestAllTypes{RepeatedNestedMessage: []*testeditionspb.TestAllTypes_NestedMessage{
				{
					A: proto.Int32(1),
				}, {
					A: proto.Int32(3),
				},
			}},
			y: &testeditionspb.TestAllTypes{RepeatedNestedMessage: []*testeditionspb.TestAllTypes_NestedMessage{
				{
					A: proto.Int32(1),
				}, {
					A: proto.Int32(2),
				},
			}},
			isLess: false,
		},

		// Maps: various configurations.
		{
			x:      &testpb.TestAllTypes{MapInt32Int32: map[int32]int32{1: 2}},
			y:      &testpb.TestAllTypes{MapInt32Int32: map[int32]int32{1: 2}},
			isLess: false,
		},
		{
			x:      &testpb.TestAllTypes{MapInt32Int32: map[int32]int32{1: 2}},
			y:      &testpb.TestAllTypes{MapInt32Int32: map[int32]int32{3: 4}},
			isLess: true,
		},
		{
			x:      &testpb.TestAllTypes{MapInt32Int32: map[int32]int32{3: 4}},
			y:      &testpb.TestAllTypes{MapInt32Int32: map[int32]int32{1: 2}},
			isLess: false,
		},
		{
			x:      &testpb.TestAllTypes{MapInt32Int32: map[int32]int32{1: 2}},
			y:      &testpb.TestAllTypes{MapInt32Int32: map[int32]int32{1: 2, 3: 4}},
			isLess: true,
		},
		{
			x:      &testpb.TestAllTypes{MapInt32Int32: map[int32]int32{1: 2, 3: 4}},
			y:      &testpb.TestAllTypes{MapInt32Int32: map[int32]int32{1: 2}},
			isLess: false,
		},
		{
			x:      &testeditionspb.TestAllTypes{MapInt32Int32: map[int32]int32{1: 2}},
			y:      &testeditionspb.TestAllTypes{MapInt32Int32: map[int32]int32{3: 4}},
			isLess: true,
		},
		{
			x:      &testeditionspb.TestAllTypes{MapInt32Int32: map[int32]int32{3: 4}},
			y:      &testeditionspb.TestAllTypes{MapInt32Int32: map[int32]int32{1: 2}},
			isLess: false,
		},
		{
			x:      &testeditionspb.TestAllTypes{MapInt32Int32: map[int32]int32{1: 2}},
			y:      &testeditionspb.TestAllTypes{MapInt32Int32: map[int32]int32{1: 2, 3: 4}},
			isLess: true,
		},
		{
			x:      &testeditionspb.TestAllTypes{MapInt32Int32: map[int32]int32{1: 2, 3: 4}},
			y:      &testeditionspb.TestAllTypes{MapInt32Int32: map[int32]int32{1: 2}},
			isLess: false,
		},

		// Maps: various types.
		{
			x:      &testpb.TestAllTypes{MapInt32Int32: map[int32]int32{1: 2, 3: 4}},
			y:      &testpb.TestAllTypes{MapInt32Int32: map[int32]int32{1: 2, 3: 4}},
			isLess: false,
		},
		{
			x:      &testpb.TestAllTypes{MapInt32Int32: map[int32]int32{1: 2, 3: 4}},
			y:      &testpb.TestAllTypes{MapInt32Int32: map[int32]int32{1: 2, 3: 5}},
			isLess: true,
		},
		{
			x:      &testpb.TestAllTypes{MapInt32Int32: map[int32]int32{1: 2, 3: 5}},
			y:      &testpb.TestAllTypes{MapInt32Int32: map[int32]int32{1: 2, 3: 4}},
			isLess: false,
		},
		{
			x:      &testpb.TestAllTypes{MapInt64Int64: map[int64]int64{1: 2, 3: 4}},
			y:      &testpb.TestAllTypes{MapInt64Int64: map[int64]int64{1: 2, 3: 5}},
			isLess: true,
		},
		{
			x:      &testpb.TestAllTypes{MapInt64Int64: map[int64]int64{1: 2, 3: 5}},
			y:      &testpb.TestAllTypes{MapInt64Int64: map[int64]int64{1: 2, 3: 4}},
			isLess: false,
		},
		{
			x:      &testpb.TestAllTypes{MapUint32Uint32: map[uint32]uint32{1: 2, 3: 4}},
			y:      &testpb.TestAllTypes{MapUint32Uint32: map[uint32]uint32{1: 2, 3: 5}},
			isLess: true,
		},
		{
			x:      &testpb.TestAllTypes{MapUint32Uint32: map[uint32]uint32{1: 2, 3: 5}},
			y:      &testpb.TestAllTypes{MapUint32Uint32: map[uint32]uint32{1: 2, 3: 4}},
			isLess: false,
		},
		{
			x:      &testpb.TestAllTypes{MapUint64Uint64: map[uint64]uint64{1: 2, 3: 4}},
			y:      &testpb.TestAllTypes{MapUint64Uint64: map[uint64]uint64{1: 2, 3: 5}},
			isLess: true,
		},
		{
			x:      &testpb.TestAllTypes{MapUint64Uint64: map[uint64]uint64{1: 2, 3: 5}},
			y:      &testpb.TestAllTypes{MapUint64Uint64: map[uint64]uint64{1: 2, 3: 4}},
			isLess: false,
		},
		{
			x:      &testpb.TestAllTypes{MapSint32Sint32: map[int32]int32{1: 2, 3: 4}},
			y:      &testpb.TestAllTypes{MapSint32Sint32: map[int32]int32{1: 2, 3: 5}},
			isLess: true,
		},
		{
			x:      &testpb.TestAllTypes{MapSint32Sint32: map[int32]int32{1: 2, 3: 5}},
			y:      &testpb.TestAllTypes{MapSint32Sint32: map[int32]int32{1: 2, 3: 4}},
			isLess: false,
		},
		{
			x:      &testpb.TestAllTypes{MapSint64Sint64: map[int64]int64{1: 2, 3: 4}},
			y:      &testpb.TestAllTypes{MapSint64Sint64: map[int64]int64{1: 2, 3: 5}},
			isLess: true,
		},
		{
			x:      &testpb.TestAllTypes{MapSint64Sint64: map[int64]int64{1: 2, 3: 5}},
			y:      &testpb.TestAllTypes{MapSint64Sint64: map[int64]int64{1: 2, 3: 4}},
			isLess: false,
		},
		{
			x:      &testpb.TestAllTypes{MapFixed32Fixed32: map[uint32]uint32{1: 2, 3: 4}},
			y:      &testpb.TestAllTypes{MapFixed32Fixed32: map[uint32]uint32{1: 2, 3: 5}},
			isLess: true,
		},
		{
			x:      &testpb.TestAllTypes{MapFixed32Fixed32: map[uint32]uint32{1: 2, 3: 5}},
			y:      &testpb.TestAllTypes{MapFixed32Fixed32: map[uint32]uint32{1: 2, 3: 4}},
			isLess: false,
		},
		{
			x:      &testpb.TestAllTypes{MapFixed64Fixed64: map[uint64]uint64{1: 2, 3: 4}},
			y:      &testpb.TestAllTypes{MapFixed64Fixed64: map[uint64]uint64{1: 2, 3: 5}},
			isLess: true,
		},
		{
			x:      &testpb.TestAllTypes{MapFixed64Fixed64: map[uint64]uint64{1: 2, 3: 5}},
			y:      &testpb.TestAllTypes{MapFixed64Fixed64: map[uint64]uint64{1: 2, 3: 4}},
			isLess: false,
		},
		{
			x:      &testpb.TestAllTypes{MapSfixed32Sfixed32: map[int32]int32{1: 2, 3: 4}},
			y:      &testpb.TestAllTypes{MapSfixed32Sfixed32: map[int32]int32{1: 2, 3: 5}},
			isLess: true,
		},
		{
			x:      &testpb.TestAllTypes{MapSfixed32Sfixed32: map[int32]int32{1: 2, 3: 5}},
			y:      &testpb.TestAllTypes{MapSfixed32Sfixed32: map[int32]int32{1: 2, 3: 4}},
			isLess: false,
		},
		{
			x:      &testpb.TestAllTypes{MapSfixed64Sfixed64: map[int64]int64{1: 2, 3: 4}},
			y:      &testpb.TestAllTypes{MapSfixed64Sfixed64: map[int64]int64{1: 2, 3: 5}},
			isLess: true,
		},
		{
			x:      &testpb.TestAllTypes{MapSfixed64Sfixed64: map[int64]int64{1: 2, 3: 5}},
			y:      &testpb.TestAllTypes{MapSfixed64Sfixed64: map[int64]int64{1: 2, 3: 4}},
			isLess: false,
		},
		{
			x:      &testpb.TestAllTypes{MapInt32Float: map[int32]float32{1: 2, 3: 4}},
			y:      &testpb.TestAllTypes{MapInt32Float: map[int32]float32{1: 2, 3: 5}},
			isLess: true,
		},
		{
			x:      &testpb.TestAllTypes{MapInt32Float: map[int32]float32{1: 2, 3: 5}},
			y:      &testpb.TestAllTypes{MapInt32Float: map[int32]float32{1: 2, 3: 4}},
			isLess: false,
		},
		{
			x:      &testpb.TestAllTypes{MapInt32Float: map[int32]float32{0: -1, 1: 1}},
			y:      &testpb.TestAllTypes{MapInt32Float: map[int32]float32{0: 1, 1: -1}},
			isLess: true,
		},
		{
			x:      &testpb.TestAllTypes{MapInt32Float: map[int32]float32{0: 1, 1: -1}},
			y:      &testpb.TestAllTypes{MapInt32Float: map[int32]float32{0: -1, 1: 1}},
			isLess: false,
		},
		{
			x:      &testpb.TestAllTypes{MapInt32Double: map[int32]float64{1: 2, 3: 4}},
			y:      &testpb.TestAllTypes{MapInt32Double: map[int32]float64{1: 2, 3: 5}},
			isLess: true,
		},
		{
			x:      &testpb.TestAllTypes{MapInt32Double: map[int32]float64{1: 2, 3: 5}},
			y:      &testpb.TestAllTypes{MapInt32Double: map[int32]float64{1: 2, 3: 4}},
			isLess: false,
		},
		{
			x:      &testpb.TestAllTypes{MapBoolBool: map[bool]bool{true: false, false: false}},
			y:      &testpb.TestAllTypes{MapBoolBool: map[bool]bool{true: false, false: true}},
			isLess: true,
		},
		{
			x:      &testpb.TestAllTypes{MapBoolBool: map[bool]bool{true: false, false: true}},
			y:      &testpb.TestAllTypes{MapBoolBool: map[bool]bool{true: false, false: false}},
			isLess: false,
		},
		{
			x:      &testpb.TestAllTypes{MapStringString: map[string]string{"a": "b", "c": "d"}},
			y:      &testpb.TestAllTypes{MapStringString: map[string]string{"a": "b", "c": "e"}},
			isLess: true,
		},
		{
			x:      &testpb.TestAllTypes{MapStringString: map[string]string{"a": "b", "c": "e"}},
			y:      &testpb.TestAllTypes{MapStringString: map[string]string{"a": "b", "c": "d"}},
			isLess: false,
		},
		{
			x:      &testpb.TestAllTypes{MapStringBytes: map[string][]byte{"a": []byte("b"), "c": []byte("d")}},
			y:      &testpb.TestAllTypes{MapStringBytes: map[string][]byte{"a": []byte("b"), "c": []byte("e")}},
			isLess: true,
		},
		{
			x:      &testpb.TestAllTypes{MapStringBytes: map[string][]byte{"a": []byte("b"), "c": []byte("e")}},
			y:      &testpb.TestAllTypes{MapStringBytes: map[string][]byte{"a": []byte("b"), "c": []byte("d")}},
			isLess: false,
		},
		{
			x: &testpb.TestAllTypes{MapStringNestedMessage: map[string]*testpb.TestAllTypes_NestedMessage{
				"a": {
					A: proto.Int32(1),
				}, "b": {
					A: proto.Int32(2),
				},
			}},
			y: &testpb.TestAllTypes{MapStringNestedMessage: map[string]*testpb.TestAllTypes_NestedMessage{
				"a": {
					A: proto.Int32(1),
				}, "b": {
					A: proto.Int32(3),
				},
			}},
			isLess: true,
		},
		{
			x: &testpb.TestAllTypes{MapStringNestedMessage: map[string]*testpb.TestAllTypes_NestedMessage{
				"a": {
					A: proto.Int32(1),
				}, "b": {
					A: proto.Int32(3),
				},
			}},
			y: &testpb.TestAllTypes{MapStringNestedMessage: map[string]*testpb.TestAllTypes_NestedMessage{
				"a": {
					A: proto.Int32(1),
				}, "b": {
					A: proto.Int32(2),
				},
			}},
			isLess: false,
		},
		{
			x: &testpb.TestAllTypes{
				MapStringNestedEnum: map[string]testpb.TestAllTypes_NestedEnum{
					"a": testpb.TestAllTypes_FOO,
					"b": testpb.TestAllTypes_BAR,
				},
			},
			y: &testpb.TestAllTypes{
				MapStringNestedEnum: map[string]testpb.TestAllTypes_NestedEnum{
					"a": testpb.TestAllTypes_FOO,
					"b": testpb.TestAllTypes_BAZ,
				},
			},
			isLess: true,
		},
		{
			x: &testpb.TestAllTypes{
				MapStringNestedEnum: map[string]testpb.TestAllTypes_NestedEnum{
					"a": testpb.TestAllTypes_FOO,
					"b": testpb.TestAllTypes_BAZ,
				},
			},
			y: &testpb.TestAllTypes{
				MapStringNestedEnum: map[string]testpb.TestAllTypes_NestedEnum{
					"a": testpb.TestAllTypes_FOO,
					"b": testpb.TestAllTypes_BAR,
				},
			},
			isLess: false,
		},
		{
			x:      &testeditionspb.TestAllTypes{MapInt32Int32: map[int32]int32{1: 2, 3: 4}},
			y:      &testeditionspb.TestAllTypes{MapInt32Int32: map[int32]int32{1: 2, 3: 5}},
			isLess: true,
		},
		{
			x:      &testeditionspb.TestAllTypes{MapInt32Int32: map[int32]int32{1: 2, 3: 5}},
			y:      &testeditionspb.TestAllTypes{MapInt32Int32: map[int32]int32{1: 2, 3: 4}},
			isLess: false,
		},
		{
			x:      &testeditionspb.TestAllTypes{MapInt64Int64: map[int64]int64{1: 2, 3: 4}},
			y:      &testeditionspb.TestAllTypes{MapInt64Int64: map[int64]int64{1: 2, 3: 5}},
			isLess: true,
		},
		{
			x:      &testeditionspb.TestAllTypes{MapInt64Int64: map[int64]int64{1: 2, 3: 5}},
			y:      &testeditionspb.TestAllTypes{MapInt64Int64: map[int64]int64{1: 2, 3: 4}},
			isLess: false,
		},
		{
			x:      &testeditionspb.TestAllTypes{MapUint32Uint32: map[uint32]uint32{1: 2, 3: 4}},
			y:      &testeditionspb.TestAllTypes{MapUint32Uint32: map[uint32]uint32{1: 2, 3: 5}},
			isLess: true,
		},
		{
			x:      &testeditionspb.TestAllTypes{MapUint32Uint32: map[uint32]uint32{1: 2, 3: 5}},
			y:      &testeditionspb.TestAllTypes{MapUint32Uint32: map[uint32]uint32{1: 2, 3: 4}},
			isLess: false,
		},
		{
			x:      &testeditionspb.TestAllTypes{MapUint64Uint64: map[uint64]uint64{1: 2, 3: 4}},
			y:      &testeditionspb.TestAllTypes{MapUint64Uint64: map[uint64]uint64{1: 2, 3: 5}},
			isLess: true,
		},
		{
			x:      &testeditionspb.TestAllTypes{MapUint64Uint64: map[uint64]uint64{1: 2, 3: 5}},
			y:      &testeditionspb.TestAllTypes{MapUint64Uint64: map[uint64]uint64{1: 2, 3: 4}},
			isLess: false,
		},
		{
			x:      &testeditionspb.TestAllTypes{MapSint32Sint32: map[int32]int32{1: 2, 3: 4}},
			y:      &testeditionspb.TestAllTypes{MapSint32Sint32: map[int32]int32{1: 2, 3: 5}},
			isLess: true,
		},
		{
			x:      &testeditionspb.TestAllTypes{MapSint32Sint32: map[int32]int32{1: 2, 3: 5}},
			y:      &testeditionspb.TestAllTypes{MapSint32Sint32: map[int32]int32{1: 2, 3: 4}},
			isLess: false,
		},
		{
			x:      &testeditionspb.TestAllTypes{MapSint64Sint64: map[int64]int64{1: 2, 3: 4}},
			y:      &testeditionspb.TestAllTypes{MapSint64Sint64: map[int64]int64{1: 2, 3: 5}},
			isLess: true,
		},
		{
			x:      &testeditionspb.TestAllTypes{MapSint64Sint64: map[int64]int64{1: 2, 3: 5}},
			y:      &testeditionspb.TestAllTypes{MapSint64Sint64: map[int64]int64{1: 2, 3: 4}},
			isLess: false,
		},
		{
			x:      &testeditionspb.TestAllTypes{MapFixed32Fixed32: map[uint32]uint32{1: 2, 3: 4}},
			y:      &testeditionspb.TestAllTypes{MapFixed32Fixed32: map[uint32]uint32{1: 2, 3: 5}},
			isLess: true,
		},
		{
			x:      &testeditionspb.TestAllTypes{MapFixed32Fixed32: map[uint32]uint32{1: 2, 3: 5}},
			y:      &testeditionspb.TestAllTypes{MapFixed32Fixed32: map[uint32]uint32{1: 2, 3: 4}},
			isLess: false,
		},
		{
			x:      &testeditionspb.TestAllTypes{MapFixed64Fixed64: map[uint64]uint64{1: 2, 3: 4}},
			y:      &testeditionspb.TestAllTypes{MapFixed64Fixed64: map[uint64]uint64{1: 2, 3: 5}},
			isLess: true,
		},
		{
			x:      &testeditionspb.TestAllTypes{MapFixed64Fixed64: map[uint64]uint64{1: 2, 3: 5}},
			y:      &testeditionspb.TestAllTypes{MapFixed64Fixed64: map[uint64]uint64{1: 2, 3: 4}},
			isLess: false,
		},
		{
			x:      &testeditionspb.TestAllTypes{MapSfixed32Sfixed32: map[int32]int32{1: 2, 3: 4}},
			y:      &testeditionspb.TestAllTypes{MapSfixed32Sfixed32: map[int32]int32{1: 2, 3: 5}},
			isLess: true,
		},
		{
			x:      &testeditionspb.TestAllTypes{MapSfixed32Sfixed32: map[int32]int32{1: 2, 3: 5}},
			y:      &testeditionspb.TestAllTypes{MapSfixed32Sfixed32: map[int32]int32{1: 2, 3: 4}},
			isLess: false,
		},
		{
			x:      &testeditionspb.TestAllTypes{MapSfixed64Sfixed64: map[int64]int64{1: 2, 3: 4}},
			y:      &testeditionspb.TestAllTypes{MapSfixed64Sfixed64: map[int64]int64{1: 2, 3: 5}},
			isLess: true,
		},
		{
			x:      &testeditionspb.TestAllTypes{MapSfixed64Sfixed64: map[int64]int64{1: 2, 3: 5}},
			y:      &testeditionspb.TestAllTypes{MapSfixed64Sfixed64: map[int64]int64{1: 2, 3: 4}},
			isLess: false,
		},
		{
			x:      &testeditionspb.TestAllTypes{MapInt32Float: map[int32]float32{1: 2, 3: 4}},
			y:      &testeditionspb.TestAllTypes{MapInt32Float: map[int32]float32{1: 2, 3: 5}},
			isLess: true,
		},
		{
			x:      &testeditionspb.TestAllTypes{MapInt32Float: map[int32]float32{1: 2, 3: 5}},
			y:      &testeditionspb.TestAllTypes{MapInt32Float: map[int32]float32{1: 2, 3: 4}},
			isLess: false,
		},
		{
			x:      &testeditionspb.TestAllTypes{MapInt32Float: map[int32]float32{0: -1, 1: 1}},
			y:      &testeditionspb.TestAllTypes{MapInt32Float: map[int32]float32{0: 1, 1: -1}},
			isLess: true,
		},
		{
			x:      &testeditionspb.TestAllTypes{MapInt32Float: map[int32]float32{0: 1, 1: -1}},
			y:      &testeditionspb.TestAllTypes{MapInt32Float: map[int32]float32{0: -1, 1: 1}},
			isLess: false,
		},
		{
			x:      &testeditionspb.TestAllTypes{MapInt32Double: map[int32]float64{1: 2, 3: 4}},
			y:      &testeditionspb.TestAllTypes{MapInt32Double: map[int32]float64{1: 2, 3: 5}},
			isLess: true,
		},
		{
			x:      &testeditionspb.TestAllTypes{MapInt32Double: map[int32]float64{1: 2, 3: 5}},
			y:      &testeditionspb.TestAllTypes{MapInt32Double: map[int32]float64{1: 2, 3: 4}},
			isLess: false,
		},
		{
			x:      &testeditionspb.TestAllTypes{MapBoolBool: map[bool]bool{true: false, false: false}},
			y:      &testeditionspb.TestAllTypes{MapBoolBool: map[bool]bool{true: false, false: true}},
			isLess: true,
		},
		{
			x:      &testeditionspb.TestAllTypes{MapBoolBool: map[bool]bool{true: false, false: true}},
			y:      &testeditionspb.TestAllTypes{MapBoolBool: map[bool]bool{true: false, false: false}},
			isLess: false,
		},
		{
			x:      &testeditionspb.TestAllTypes{MapStringString: map[string]string{"a": "b", "c": "d"}},
			y:      &testeditionspb.TestAllTypes{MapStringString: map[string]string{"a": "b", "c": "e"}},
			isLess: true,
		},
		{
			x:      &testeditionspb.TestAllTypes{MapStringString: map[string]string{"a": "b", "c": "e"}},
			y:      &testeditionspb.TestAllTypes{MapStringString: map[string]string{"a": "b", "c": "d"}},
			isLess: false,
		},
		{
			x:      &testeditionspb.TestAllTypes{MapStringBytes: map[string][]byte{"a": []byte("b"), "c": []byte("d")}},
			y:      &testeditionspb.TestAllTypes{MapStringBytes: map[string][]byte{"a": []byte("b"), "c": []byte("e")}},
			isLess: true,
		},
		{
			x:      &testeditionspb.TestAllTypes{MapStringBytes: map[string][]byte{"a": []byte("b"), "c": []byte("e")}},
			y:      &testeditionspb.TestAllTypes{MapStringBytes: map[string][]byte{"a": []byte("b"), "c": []byte("d")}},
			isLess: false,
		},
		{
			x: &testeditionspb.TestAllTypes{MapStringNestedMessage: map[string]*testeditionspb.TestAllTypes_NestedMessage{
				"a": {
					A: proto.Int32(1),
				}, "b": {
					A: proto.Int32(2),
				},
			}},
			y: &testeditionspb.TestAllTypes{MapStringNestedMessage: map[string]*testeditionspb.TestAllTypes_NestedMessage{
				"a": {
					A: proto.Int32(1),
				}, "b": {
					A: proto.Int32(3),
				},
			}},
			isLess: true,
		},
		{
			x: &testeditionspb.TestAllTypes{MapStringNestedMessage: map[string]*testeditionspb.TestAllTypes_NestedMessage{
				"a": {
					A: proto.Int32(1),
				}, "b": {
					A: proto.Int32(3),
				},
			}},
			y: &testeditionspb.TestAllTypes{MapStringNestedMessage: map[string]*testeditionspb.TestAllTypes_NestedMessage{
				"a": {
					A: proto.Int32(1),
				}, "b": {
					A: proto.Int32(2),
				},
			}},
			isLess: false,
		},
		{
			x: &testeditionspb.TestAllTypes{
				MapStringNestedEnum: map[string]testeditionspb.TestAllTypes_NestedEnum{
					"a": testeditionspb.TestAllTypes_FOO,
					"b": testeditionspb.TestAllTypes_BAR,
				},
			},
			y: &testeditionspb.TestAllTypes{
				MapStringNestedEnum: map[string]testeditionspb.TestAllTypes_NestedEnum{
					"a": testeditionspb.TestAllTypes_FOO,
					"b": testeditionspb.TestAllTypes_BAZ,
				},
			},
			isLess: true,
		},
		{
			x: &testeditionspb.TestAllTypes{
				MapStringNestedEnum: map[string]testeditionspb.TestAllTypes_NestedEnum{
					"a": testeditionspb.TestAllTypes_FOO,
					"b": testeditionspb.TestAllTypes_BAZ,
				},
			},
			y: &testeditionspb.TestAllTypes{
				MapStringNestedEnum: map[string]testeditionspb.TestAllTypes_NestedEnum{
					"a": testeditionspb.TestAllTypes_FOO,
					"b": testeditionspb.TestAllTypes_BAR,
				},
			},
			isLess: false,
		},

		// Extensions.
		{
			x: build(&testpb.TestAllExtensions{},
				extend(testpb.E_OptionalInt32, int32(1)),
			),
			y: build(&testpb.TestAllExtensions{},
				extend(testpb.E_OptionalInt32, int32(2)),
			),
			isLess: true,
		},
		{
			x: build(&testpb.TestAllExtensions{},
				extend(testpb.E_OptionalInt32, int32(2)),
			),
			y: build(&testpb.TestAllExtensions{},
				extend(testpb.E_OptionalInt32, int32(1)),
			),
			isLess: false,
		},
		{
			x: &testpb.TestAllExtensions{},
			y: build(&testpb.TestAllExtensions{},
				extend(testpb.E_OptionalInt32, int32(2)),
			),
			isLess: true,
		},
		{
			x: build(&testpb.TestAllExtensions{},
				extend(testpb.E_OptionalInt32, int32(2)),
			),
			y:      &testpb.TestAllExtensions{},
			isLess: false,
		},
		{
			x: build(&testpb.TestAllTypes{}, unknown(protopack.Message{
				protopack.Tag{Number: 100000, Type: protopack.VarintType}, protopack.Varint(1),
			}.Marshal())),
			y: build(&testpb.TestAllTypes{}, unknown(protopack.Message{
				protopack.Tag{Number: 100000, Type: protopack.VarintType}, protopack.Varint(2),
			}.Marshal())),
			isLess: true,
		},
		{
			x: build(&testpb.TestAllTypes{}, unknown(protopack.Message{
				protopack.Tag{Number: 100000, Type: protopack.VarintType}, protopack.Varint(2),
			}.Marshal())),
			y: build(&testpb.TestAllTypes{}, unknown(protopack.Message{
				protopack.Tag{Number: 100000, Type: protopack.VarintType}, protopack.Varint(1),
			}.Marshal())),
			isLess: false,
		},
		{
			x: build(&testpb.TestAllTypes{}, unknown(protopack.Message{
				protopack.Tag{Number: 100000, Type: protopack.VarintType}, protopack.Varint(1),
			}.Marshal())),
			y:      &testpb.TestAllTypes{},
			isLess: false,
		},
		{
			x: &testpb.TestAllTypes{},
			y: build(&testpb.TestAllTypes{}, unknown(protopack.Message{
				protopack.Tag{Number: 100000, Type: protopack.VarintType}, protopack.Varint(1),
			}.Marshal())),
			isLess: true,
		},
		{
			x: build(&testeditionspb.TestAllTypes{}, unknown(protopack.Message{
				protopack.Tag{Number: 100000, Type: protopack.VarintType}, protopack.Varint(2),
			}.Marshal())),
			y: build(&testeditionspb.TestAllTypes{}, unknown(protopack.Message{
				protopack.Tag{Number: 100000, Type: protopack.VarintType}, protopack.Varint(1),
			}.Marshal())),
			isLess: false,
		},
		{
			x: build(&testeditionspb.TestAllTypes{}, unknown(protopack.Message{
				protopack.Tag{Number: 100000, Type: protopack.VarintType}, protopack.Varint(1),
			}.Marshal())),
			y: build(&testeditionspb.TestAllTypes{}, unknown(protopack.Message{
				protopack.Tag{Number: 100000, Type: protopack.VarintType}, protopack.Varint(2),
			}.Marshal())),
			isLess: true,
		},
		{
			x: &testeditionspb.TestAllTypes{},
			y: build(&testeditionspb.TestAllTypes{}, unknown(protopack.Message{
				protopack.Tag{Number: 100000, Type: protopack.VarintType}, protopack.Varint(1),
			}.Marshal())),
			isLess: true,
		},
		{
			x: build(&testeditionspb.TestAllTypes{}, unknown(protopack.Message{
				protopack.Tag{Number: 100000, Type: protopack.VarintType}, protopack.Varint(1),
			}.Marshal())),
			y:      &testeditionspb.TestAllTypes{},
			isLess: false,
		},
		{
			x: build(&testeditionspb.TestAllTypes{}, unknown(protopack.Message{
				protopack.Tag{Number: 1000, Type: protopack.VarintType}, protopack.Varint(1),
			}.Marshal())),
			y: build(&testeditionspb.TestAllTypes{}, unknown(protopack.Message{
				protopack.Tag{Number: 100000, Type: protopack.VarintType}, protopack.Varint(1),
			}.Marshal())),
			isLess: true,
		},
		{
			x: build(&testeditionspb.TestAllTypes{}, unknown(protopack.Message{
				protopack.Tag{Number: 100000, Type: protopack.VarintType}, protopack.Varint(1),
			}.Marshal())),
			y: build(&testeditionspb.TestAllTypes{}, unknown(protopack.Message{
				protopack.Tag{Number: 1000, Type: protopack.VarintType}, protopack.Varint(1),
			}.Marshal())),
			isLess: false,
		},
		{
			x: build(&testeditionspb.TestAllTypes{}, unknown(protopack.Message{
				protopack.Tag{Number: 100000, Type: protopack.Fixed32Type}, protopack.Int32(1),
			}.Marshal())),
			y: build(&testeditionspb.TestAllTypes{}, unknown(protopack.Message{
				protopack.Tag{Number: 100000, Type: protopack.VarintType}, protopack.Varint(1),
			}.Marshal())),
			isLess: false,
		},
		{
			x: build(&testeditionspb.TestAllTypes{}, unknown(protopack.Message{
				protopack.Tag{Number: 100000, Type: protopack.VarintType}, protopack.Varint(1),
			}.Marshal())),
			y: build(&testeditionspb.TestAllTypes{}, unknown(protopack.Message{
				protopack.Tag{Number: 100000, Type: protopack.Fixed32Type}, protopack.Int32(1),
			}.Marshal())),
			isLess: true,
		},
		{
			x: build(&testeditionspb.TestAllTypes{}, unknown(protopack.Message{
				protopack.Tag{Number: 100000, Type: protopack.VarintType}, protopack.Varint(1),
			}.Marshal())),
			y: build(&testeditionspb.TestAllTypes{}, unknown(protopack.Message{
				protopack.Tag{Number: 100000, Type: protopack.VarintType}, protopack.Varint(1),
			}.Marshal())),
			isLess: false,
		},
	}

	for _, tt := range tests {
		if got := proto.LessThan(tt.x, tt.y); got != tt.isLess {
			t.Errorf(
				"LessThan(x, y) = %t, want %t\n==== x ====\n%v\n==== y ====\n%v",
				got,
				tt.isLess,
				prototext.Format(tt.x),
				prototext.Format(tt.y),
			)
		}
	}
}
