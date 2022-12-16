package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
)

func main() {
	lines, err := readLines("input.txt")
	if err != nil {
		log.Fatalf("Failed to read input file: %s", err)
	}

	// Part 1
	solutionPart1 := solvePart1(lines)
	fmt.Printf("Solution for part 1 is %s\n", solutionPart1)

	// Part 2
	solutionPart2 := solvePart2(lines)
	fmt.Printf("Solution for part 2 is %s\n", solutionPart2)
}

func solvePart1(lines []string) string {
	stacksOfCratesLines, instructions := parseInput(lines)
	stacksOfCrates := identifyStacks(stacksOfCratesLines)
	stacksOfCrates = executeInstructions(stacksOfCrates, instructions, false)

	// Get top crates from each stack
	result := ""
	for i := 1; i <= len(stacksOfCrates); i++ {
		stack := stacksOfCrates[i]
		crate, _ := stack.Pop()
		result += crate
	}

	return result
}

func solvePart2(lines []string) string {
	stacksOfCratesLines, instructions := parseInput(lines)
	stacksOfCrates := identifyStacks(stacksOfCratesLines)
	stacksOfCrates = executeInstructions(stacksOfCrates, instructions, true)

	// Get top crates from each stack
	result := ""
	for i := 1; i <= len(stacksOfCrates); i++ {
		stack := stacksOfCrates[i]
		crate, _ := stack.Pop()
		result += crate
	}

	return result
}

func parseInput(lines []string) ([]string, []string) {
	stacksOfCratesLines := make([]string, 0)
	instructions := make([]string, 0)

	for _, line := range lines {
		if line != "" && !strings.Contains(line, "move") {
			stacksOfCratesLines = append(stacksOfCratesLines, line)
		} else if line != "" {
			instructions = append(instructions, line)
		}
	}

	return stacksOfCratesLines, instructions
}

func identifyStacks(stacksOfCratesLines []string) map[int]Stack {
	stacksOfCrates := make(map[int]Stack)

	for x, stackIndex := range stacksOfCratesLines[len(stacksOfCratesLines)-1] {
		if stackIndex != ' ' {
			tmpStack := make(Stack, 0)
			// -1 because array starts at 0 and -1 because we don't want the last line
			for y := len(stacksOfCratesLines) - 2; y >= 0; y-- {
				var stackChar = string(stacksOfCratesLines[y][x])
				if stackChar != " " {
					tmpStack = append(tmpStack, stackChar)
				}
			}
			stacksOfCrates[convertRuneToDigit(stackIndex)] = tmpStack
		}
	}

	return stacksOfCrates
}

func executeInstructions(stackOfCrates map[int]Stack, instructions []string, retainOrder bool) map[int]Stack {
	for _, instruction := range instructions {

		re := regexp.MustCompile(`\d+`)
		numbers := re.FindAllString(instruction, 3)

		moveCount, _ := strconv.Atoi(numbers[0])
		fromIndex, _ := strconv.Atoi(numbers[1])
		toIndex, _ := strconv.Atoi(numbers[2])

		tmpStack := make(Stack, 0)
		for i := 0; i < moveCount; i++ {
			fromStack := stackOfCrates[fromIndex]
			toStack := stackOfCrates[toIndex]
			movingCrate, success := fromStack.Pop()
			if success {
				if !retainOrder {
					toStack.Push(movingCrate)
				} else {
					// If we want to retain order, use intermediary stack to "invert order twice" ! See Hanoi tower
					tmpStack.Push(movingCrate)
				}
			}

			// Reassigning the modified values
			stackOfCrates[fromIndex] = fromStack
			if !retainOrder {
				stackOfCrates[toIndex] = toStack
			}
		}

		if retainOrder {
			toStack := stackOfCrates[toIndex]
			// Save the length in a var otherwise it changes as the loop goes on...
			tmpLength := len(tmpStack)
			for i := 0; i < tmpLength; i++ {
				movingCrate, _ := tmpStack.Pop()
				toStack.Push(movingCrate)
			}
			stackOfCrates[toIndex] = toStack
		}
	}

	return stackOfCrates
}

func convertRuneToDigit(r rune) int {
	digit := int(r - '0')
	if digit >= 0 && digit <= 10 {
		return digit
	}
	return -1
}

// readLines reads a whole file into memory
// and returns a slice of its lines.
func readLines(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}

type Stack []string

func (s *Stack) IsEmpty() bool {
	return len(*s) == 0
}

func (s *Stack) Push(str string) {
	*s = append(*s, str)
}

func (s *Stack) Pop() (string, bool) {
	if s.IsEmpty() {
		return "", false
	} else {
		indexOfTopElement := len(*s) - 1
		topElement := (*s)[indexOfTopElement]
		*s = (*s)[:indexOfTopElement]
		return topElement, true
	}
}
