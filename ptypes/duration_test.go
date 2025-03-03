// Copyright 2016 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ptypes

import (
	"math"
	"testing"
	"time"

	"github.com/ryanleesmith/protobuf/proto"

	durpb "github.com/ryanleesmith/protobuf/ptypes/duration"
)

const (
	minGoSeconds = math.MinInt64 / int64(1e9)
	maxGoSeconds = math.MaxInt64 / int64(1e9)
)

var durationTests = []struct {
	proto   *durpb.Duration
	isValid bool
	inRange bool
	dur     time.Duration
}{
	// The zero duration.
	{&durpb.Duration{Seconds: 0, Nanos: 0}, true, true, 0},
	// Some ordinary non-zero durations.
	{&durpb.Duration{Seconds: 100, Nanos: 0}, true, true, 100 * time.Second},
	{&durpb.Duration{Seconds: -100, Nanos: 0}, true, true, -100 * time.Second},
	{&durpb.Duration{Seconds: 100, Nanos: 987}, true, true, 100*time.Second + 987},
	{&durpb.Duration{Seconds: -100, Nanos: -987}, true, true, -(100*time.Second + 987)},
	// The largest duration representable in Go.
	{&durpb.Duration{Seconds: maxGoSeconds, Nanos: int32(math.MaxInt64 - 1e9*maxGoSeconds)}, true, true, math.MaxInt64},
	// The smallest duration representable in Go.
	{&durpb.Duration{Seconds: minGoSeconds, Nanos: int32(math.MinInt64 - 1e9*minGoSeconds)}, true, true, math.MinInt64},
	{nil, false, false, 0},
	{&durpb.Duration{Seconds: -100, Nanos: 987}, false, false, 0},
	{&durpb.Duration{Seconds: 100, Nanos: -987}, false, false, 0},
	{&durpb.Duration{Seconds: math.MinInt64, Nanos: 0}, false, false, 0},
	{&durpb.Duration{Seconds: math.MaxInt64, Nanos: 0}, false, false, 0},
	// The largest valid duration.
	{&durpb.Duration{Seconds: maxSeconds, Nanos: 1e9 - 1}, true, false, 0},
	// The smallest valid duration.
	{&durpb.Duration{Seconds: minSeconds, Nanos: -(1e9 - 1)}, true, false, 0},
	// The smallest invalid duration above the valid range.
	{&durpb.Duration{Seconds: maxSeconds + 1, Nanos: 0}, false, false, 0},
	// The largest invalid duration below the valid range.
	{&durpb.Duration{Seconds: minSeconds - 1, Nanos: -(1e9 - 1)}, false, false, 0},
	// One nanosecond past the largest duration representable in Go.
	{&durpb.Duration{Seconds: maxGoSeconds, Nanos: int32(math.MaxInt64-1e9*maxGoSeconds) + 1}, true, false, 0},
	// One nanosecond past the smallest duration representable in Go.
	{&durpb.Duration{Seconds: minGoSeconds, Nanos: int32(math.MinInt64-1e9*minGoSeconds) - 1}, true, false, 0},
	// One second past the largest duration representable in Go.
	{&durpb.Duration{Seconds: maxGoSeconds + 1, Nanos: int32(math.MaxInt64 - 1e9*maxGoSeconds)}, true, false, 0},
	// One second past the smallest duration representable in Go.
	{&durpb.Duration{Seconds: minGoSeconds - 1, Nanos: int32(math.MinInt64 - 1e9*minGoSeconds)}, true, false, 0},
}

func TestValidateDuration(t *testing.T) {
	for _, test := range durationTests {
		err := validateDuration(test.proto)
		gotValid := (err == nil)
		if gotValid != test.isValid {
			t.Errorf("validateDuration(%v) = %t, want %t", test.proto, gotValid, test.isValid)
		}
	}
}

func TestDuration(t *testing.T) {
	for _, test := range durationTests {
		got, err := Duration(test.proto)
		gotOK := (err == nil)
		wantOK := test.isValid && test.inRange
		if gotOK != wantOK {
			t.Errorf("Duration(%v) ok = %t, want %t", test.proto, gotOK, wantOK)
		}
		if err == nil && got != test.dur {
			t.Errorf("Duration(%v) = %v, want %v", test.proto, got, test.dur)
		}
	}
}

func TestDurationProto(t *testing.T) {
	for _, test := range durationTests {
		if test.isValid && test.inRange {
			got := DurationProto(test.dur)
			if !proto.Equal(got, test.proto) {
				t.Errorf("DurationProto(%v) = %v, want %v", test.dur, got, test.proto)
			}
		}
	}
}
