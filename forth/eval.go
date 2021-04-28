// +build !solution

package main

import (
	"fmt"
	"strconv"
	"strings"
)

type Evaluator struct {
	wordDefinitions map[string][]string
	stack           []int
}

// NewEvaluator creates evaluator.
func NewEvaluator() *Evaluator {
	evaluator := Evaluator{}
	evaluator.wordDefinitions = make(map[string][]string)
	return &evaluator
}

// Process evaluates sequence of words or definition.
//
// Returns resulting stack state and an error.
func (e *Evaluator) Process(row string) ([]int, error) {
	rowArray := strings.Split(row, " ")
	n := len(rowArray)

	rowArrayFormatted := make([]string, n)
	for i := 0; i < n; i++ {
		rowArrayFormatted[i] = strings.ToLower(rowArray[i])
	}

	if rowArrayFormatted[0] == ":" {
		return e.addDefinition(rowArrayFormatted[1], rowArrayFormatted[2:n-1])
	} else {
		return e.processArray(rowArrayFormatted)
	}
}

func (e *Evaluator) addDefinition(commandName string, commands []string) ([]int, error) {
	resultCommands := make([]string, 0)

	if _, err := strconv.Atoi(commandName); err == nil {
		return nil, fmt.Errorf("can't redefine number")
	}

	for _, command := range commands {
		commandsMap, contain := e.wordDefinitions[command]
		if contain {
			resultCommands = append(resultCommands, commandsMap...)
		} else {
			resultCommands = append(resultCommands, command)
		}
	}
	e.wordDefinitions[commandName] = resultCommands

	return e.stack, nil
}

func (e *Evaluator) processArray(commands []string) ([]int, error) {
	var rez []int
	var err error

	for _, v := range commands {
		rez, err = e.subprocess(v)
		if err != nil {
			return nil, err
		}
	}
	return rez, nil
}

func (e *Evaluator) subprocess(command string) ([]int, error) {
	if number, err := strconv.Atoi(command); err == nil {
		e.stack = append(e.stack, number)
		return e.stack, nil
	}

	cmd, contains := e.wordDefinitions[command]
	if contains {
		return e.processArray(cmd)
	}

	if len(e.stack) < 1 {
		return nil, fmt.Errorf("small stack")
	}
	e1 := e.stack[len(e.stack)-1]
	e.stack = e.stack[:len(e.stack)-1]

	switch command {
	case "drop":
	case "dup":
		e.stack = append(e.stack, e1, e1)
	default:
		if len(e.stack) < 1 {
			return nil, fmt.Errorf("small stack")
		}
		e2 := e.stack[len(e.stack)-1]
		e.stack = e.stack[:len(e.stack)-1]

		switch command {
		case "over":
			e.stack = append(e.stack, e2, e1, e2)
		case "swap":
			e.stack = append(e.stack, e1, e2)
		case "+":
			e.stack = append(e.stack, e2+e1)
		case "-":
			e.stack = append(e.stack, e2-e1)
		case "*":
			e.stack = append(e.stack, e2*e1)
		case "/":
			if e1 == 0 {
				return nil, fmt.Errorf("div zero")
			}
			e.stack = append(e.stack, e2/e1)
		default:
			return nil, fmt.Errorf("unsupported operation")
		}
	}

	return e.stack, nil
}
