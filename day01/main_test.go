package main

import (
	"strings"
	"testing"
)

func TestParseLine(t *testing.T) {
	cases := []struct {
		Test  string
		Input string
		Want  int
	}{
		{Test: "Case 1: No input", Input: "", Want: 0},
		{Test: "Case 2: Invalid Prefix", Input: "14\n", Want: 0},
		{Test: "Case 3", Input: "R14\n", Want: 0},
		{Test: "Case 4", Input: "R50\nR133", Want: 2},
		// {Input: "L68\nL30\nR48\nL5\nR60\nL55\nL1\nL99\nR14\nL82\n", Want: 3},
		// {Input: "L68\nL30\nR48\nL5\nR60\nL55", Want: 2},
		// {Input: "L68\nL30\nR48\nL5\nR60", Want: 1},
	}
	for _, test := range cases {
		t.Run(test.Test, func(t *testing.T) {
			fakeFile := strings.NewReader(test.Input)

			got, err := ParseLine(fakeFile)
			if err != nil {
				t.Fatalf("ParseLine() returned an error: %v", err)
			}

			if got != test.Want {
				t.Errorf("ParseLine() got %d, want %d", int(got), test.Want)
			}
		})
	}
}

func TestTurnDial(t *testing.T) {
	type Expected struct {
		numOfResets         int
		updatedDialPosition int
	}

	cases := []struct {
		test                string
		currentDialPosition int
		clicks              int
		direction           string
		want                Expected
	}{
		{
			test:                "Case 1: Click must be non negative",
			currentDialPosition: StartingDialPosition,
			clicks:              -1,
			direction:           "R",
			want: Expected{
				numOfResets:         0,
				updatedDialPosition: StartingDialPosition,
			},
		},
		{
			test:                "R - Case 2",
			currentDialPosition: StartingDialPosition,
			clicks:              49,
			direction:           "R",
			want: Expected{
				numOfResets:         0,
				updatedDialPosition: 99,
			},
		}, {
			test:                "R - Case 3",
			currentDialPosition: StartingDialPosition,
			clicks:              52,
			direction:           "R",
			want: Expected{
				numOfResets:         1,
				updatedDialPosition: 2,
			},
		},
		{
			test:                "R - Case 4",
			currentDialPosition: StartingDialPosition,
			clicks:              585,
			direction:           "R",
			want: Expected{
				numOfResets:         6,
				updatedDialPosition: 35,
			},
		},
	}
	for _, test := range cases {
		t.Run(test.test, func(t *testing.T) {
			gotNum, gotPos := TurnDial(test.direction, test.clicks, test.currentDialPosition)
			checkResults(t, gotNum, test.want.numOfResets, gotPos, test.want.updatedDialPosition)
		})
	}
}
func checkResults(t *testing.T, gotNum int, wantNum int, gotPos int, wantPos int) {
	t.Helper()
	if gotNum != wantNum {
		t.Errorf("Number of Resets mismatch: got %d, want %d", gotNum, wantNum)
	}
	if gotPos != wantPos {
		t.Errorf("Dial Position mismatch: got %d, want %d", gotPos, wantPos)
	}
}
