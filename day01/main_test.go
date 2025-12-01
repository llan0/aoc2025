package main

import (
	"errors"
	"strings"
	"testing"
)

func TestParseLine(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantVal int
		wantErr error
	}{
		{name: "Case 1: No input", input: "", wantVal: 0, wantErr: nil},
		{name: "Case 2: Invalid Prefix", input: "14\n", wantVal: 0, wantErr: ErrParseDirection},
		{name: "Case 3: Valid Single Line", input: "R14\n", wantVal: 0, wantErr: nil},
	}

	for _, test := range tests {
		tc := test //  only for parellel testing
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			fakeFile := strings.NewReader(tc.input)
			got, err := ParseLine(fakeFile)

			assertError(t, err, tc.wantErr)
			assertEqual(t, got, tc.wantVal)
		})
	}
}

func TestTurnDial(t *testing.T) {
	type result struct {
		resets int
		pos    int
	}

	tests := []struct {
		name    string
		current int
		clicks  int
		dir     string
		want    result
	}{
		{name: "R: Normal Turn (No Reset)", current: 50, clicks: 49, dir: "R", want: result{resets: 0, pos: 99}},
		{name: "R: Land exactly on 0 (Reset)", current: 50, clicks: 50, dir: "R", want: result{resets: 1, pos: 0}},
		{name: "R: Pass 0 but don't stop (No Reset)", current: 50, clicks: 52, dir: "R", want: result{resets: 0, pos: 2}},
		{name: "R: Multiple rotations landing on non-zero", current: 50, clicks: 585, dir: "R", want: result{resets: 0, pos: 35}},

		{name: "L: Normal Turn (No Reset)", current: 50, clicks: 10, dir: "L", want: result{resets: 0, pos: 40}},
		{name: "L: Land exactly on 0 (Reset)", current: 50, clicks: 50, dir: "L", want: result{resets: 1, pos: 0}},
		{name: "L: Wrap around (No Reset)", current: 50, clicks: 51, dir: "L", want: result{resets: 0, pos: 99}},
		{name: "L: Multiple rotations landing on 0", current: 50, clicks: 150, dir: "L", want: result{resets: 1, pos: 0}},
	}

	for _, test := range tests {
		tc := test
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			gotNum, gotPos := TurnDial(tc.dir, tc.clicks, tc.current)
			assertDialResults(t, gotNum, tc.want.resets, gotPos, tc.want.pos)
		})
	}
}

func TestParseDirection(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantStr string
		wantErr error
	}{
		{name: "Case 1: No Prefix", input: "80", wantStr: "", wantErr: ErrParseDirection},
		{name: "Case 2: Valid Input", input: "L80", wantStr: "L", wantErr: nil},
	}

	for _, test := range tests {
		tc := test
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			gotStr, gotErr := ParseDirection(tc.input)

			assertError(t, gotErr, tc.wantErr)
			assertEqual(t, gotStr, tc.wantStr)
		})
	}
}

// Helpers ----------------------------------------
func assertEqual[T comparable](t testing.TB, got, want T) {
	t.Helper()

	if got != want {
		t.Errorf("got %v, want %v", got, want)
	}
}

func assertError(t testing.TB, got, want error) {
	t.Helper()

	if want == nil {
		if got != nil {
			t.Fatalf("expected no error, but got: %v", got)
		}
		return
	}
	if got == nil {
		t.Fatalf("expected error %q, but got nil", want)
	}
	if !errors.Is(got, want) {
		t.Errorf("error mismatch: got %v, want %v", got, want)
	}
}

func assertDialResults(t *testing.T, gotNum, wantNum, gotPos, wantPos int) {
	t.Helper()
	if gotNum != wantNum {
		t.Errorf("number of resets mismatch: got %d, want %d", gotNum, wantNum)
	}
	if gotPos != wantPos {
		t.Errorf("dial position mismatch: got %d, want %d", gotPos, wantPos)
	}
}
