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
		{
			name:    "Case 1: No input",
			input:   "",
			wantVal: 0,
			wantErr: nil,
		},
		{
			name:    "Case 2: Invalid Prefix",
			input:   "14\n",
			wantVal: 0,
			wantErr: ErrParseDirection,
		},
		{
			name:    "Case 3: Valid Single Line",
			input:   "R14\n",
			wantVal: 0,
			wantErr: nil,
		},
		{
			name:    "Case 4: Multi Line",
			input:   "R50\nR133",
			wantVal: 2,
			wantErr: nil,
		},
	}

	for _, test := range tests {
		tc := test //  only for parellel testing
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			fakeFile := strings.NewReader(tc.input)
			got, err := ParseLine(fakeFile)

			assertError(t, err, tc.wantErr)

			if tc.wantErr == nil && got != tc.wantVal {
				t.Errorf("ParseLine() got %d, want %d", got, tc.wantVal)
			}
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
		{
			name:    "Case 1: Click must be non negative",
			current: StartingDialPosition,
			clicks:  -1,
			dir:     "R",
			want:    result{resets: 0, pos: StartingDialPosition},
		},
		{
			name:    "R - Case 2: Normal Turn",
			current: StartingDialPosition,
			clicks:  49,
			dir:     "R",
			want:    result{resets: 0, pos: 99},
		},
		{
			name:    "R - Case 3: One Reset",
			current: StartingDialPosition,
			clicks:  52,
			dir:     "R",
			want:    result{resets: 1, pos: 2},
		},
		{
			name:    "R - Case 4: Multiple Resets",
			current: StartingDialPosition,
			clicks:  585,
			dir:     "R",
			want:    result{resets: 6, pos: 35},
		},
		// TODO: Check for L turning dials
	}

	for _, test := range tests {
		tc := test
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			gotNum, gotPos := TurnDial(tc.dir, tc.clicks, tc.current)
			checkDialResults(t, gotNum, tc.want.resets, gotPos, tc.want.pos)
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
		{
			name:    "Case 1: No Prefix",
			input:   "80",
			wantStr: "",
			wantErr: ErrParseDirection,
		},
		{
			name:    "Case 2: Valid Input",
			input:   "L80",
			wantStr: "L",
			wantErr: nil,
		},
	}

	for _, test := range tests {
		tc := test
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			gotStr, gotErr := ParseDirection(tc.input)

			assertError(t, gotErr, tc.wantErr)

			if gotStr != tc.wantStr {
				t.Errorf("ParseDirection string mismatch. Got %q, want %q", gotStr, tc.wantStr)
			}
		})
	}
}

// Helpers
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

func checkDialResults(t *testing.T, gotNum, wantNum, gotPos, wantPos int) {
	t.Helper()
	if gotNum != wantNum {
		t.Errorf("Number of Resets mismatch: got %d, want %d", gotNum, wantNum)
	}
	if gotPos != wantPos {
		t.Errorf("Dial Position mismatch: got %d, want %d", gotPos, wantPos)
	}
}
