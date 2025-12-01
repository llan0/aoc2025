package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
)

const (
	MaxDialCycleClicks   = 100
	StartingDialPosition = 50
)

func main() {
	file, err := os.Open("day01/input.txt")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error opening file %v: %v\n", file, err)
		os.Exit(1)
	}
	defer file.Close()

	count, err := ParseLine(file)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error counting lines in file %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Password: %d\n", count)
}

func ParseLine(r io.Reader) (int, error) {
	scanner := bufio.NewScanner(r)
	var count int
	currentDialPosition := StartingDialPosition

	for scanner.Scan() {
		// parse rotation direction
		direction, err := ParseDirection(scanner.Text())
		if err != nil {
			return 0, nil // this doesnt make sense?? catch error
		}

		// parse the clicks
		clicks, err := ParseClicks(scanner.Text())
		if err != nil {
			return 0, err
		}
		// turn dial
		numOfResets, updatedDialPosition := TurnDial(direction, clicks, currentDialPosition)
		count += numOfResets
		currentDialPosition = updatedDialPosition
	}
	if err := scanner.Err(); err != nil {
		return 0, err
	}
	return count, nil
}
func ParseDirection(line string) (string, error) {
	if line[0] != 'L' && line[0] != 'R' {
		return "", errors.New("Invalid direction prefix")
	}
	return string(line[0]), nil
}
func ParseClicks(line string) (int, error) {
	clicks, err := strconv.ParseInt(line[1:], 10, 64)
	if err != nil {
		return -1, err
	}
	return int(clicks), nil
}

func TurnDial(direction string, clicks, currentDialPosition int) (int, int) {
	if clicks < 0 {
		return 0, currentDialPosition
	}
	switch direction {
	case "R":
		totalClicksFromZero := currentDialPosition + clicks
		numOfResets := totalClicksFromZero / MaxDialCycleClicks
		updatedDialPosition := totalClicksFromZero % MaxDialCycleClicks
		return numOfResets, updatedDialPosition
	default:
		return 0, 0
	}
}
