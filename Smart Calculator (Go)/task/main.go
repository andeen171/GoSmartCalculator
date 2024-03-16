package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"unicode"
)

var variables = map[string]int{}
var precedence = map[string]int{
	"+": 1, "-": 1,
	"*": 2, "/": 2,
	"^": 3,
}

func isOperator(token string) bool {
	_, ok := precedence[token]
	return ok
}

func main() {
	reader := bufio.NewReader(os.Stdin)

	for {
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading input:", err)
			continue
		}

		input = strings.TrimRight(input, "\r\n")

		if input == "" {
			continue
		}

		if strings.HasPrefix(input, "/") {
			if input == "/exit" {
				fmt.Println("Bye!")
				break
			} else if input == "/help" {
				printHelp()
			} else {
				fmt.Println("Unknown command")
			}
			continue
		}

		processInput(input)
	}
}

func printHelp() {
	fmt.Println(`This program supports basic arithmetic operations: addition and subtraction, and variable assignment.
Enter an expression to calculate its value or use one of the commands:
- /help to display this message.
- /exit to quit the program.`)
}

func isValidVariable(input string) bool {
	for _, r := range input {
		if !unicode.IsLetter(r) && !unicode.IsSpace(r) {
			return false
		}
	}
	return true
}

func normalizeOperators(input string) string {
	input = strings.ReplaceAll(input, " ", "")
	input = strings.ReplaceAll(input, "---", "-")
	input = strings.ReplaceAll(input, "--", "+")
	input = strings.ReplaceAll(input, "+", " + ")
	input = strings.ReplaceAll(input, "-", " - ")
	input = strings.ReplaceAll(input, "=", " = ")
	input = strings.ReplaceAll(input, "^", " ^ ")
	input = strings.ReplaceAll(input, "*", " * ")
	input = strings.ReplaceAll(input, "/", " / ")
	input = strings.ReplaceAll(input, "(", " ( ")
	input = strings.ReplaceAll(input, ")", " ) ")
	return input
}

func processInput(input string) {
	if strings.Count(input, "(") != strings.Count(input, ")") {
		fmt.Println("invalid expression")
		return
	}
	if strings.Contains(input, "=") {
		handleAssignment(input)
	} else {
		input = normalizeOperators(input)
		postfix, err := infixToPostfix(strings.Fields(input))
		result, err := calculatePostfix(postfix)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(result)
	}
}

func handleAssignment(input string) {
	parts := strings.Split(input, "=")
	if len(parts) != 2 || !isValidVariable(strings.TrimSpace(parts[0])) {
		fmt.Println("Invalid assignment")
		return
	}

	variable := strings.TrimSpace(parts[0])
	valueStr := strings.TrimSpace(parts[1])
	value, err := getNumFromString(valueStr)
	if err != nil {
		fmt.Println("Invalid assignment")
		return
	}
	variables[variable] = value
}

func getNumFromString(input string) (int, error) {
	if isValidVariable(input) {
		if num, ok := variables[input]; ok {
			return num, nil
		} else {
			return 0, fmt.Errorf("unknown variable")
		}
	}
	value, err := strconv.Atoi(strings.TrimSpace(input))
	if err != nil {
		return 0, fmt.Errorf("invalid identifier")
	}
	return value, nil
}

func infixToPostfix(tokens []string) ([]string, error) {
	var postfix []string
	var stack []string

	for _, token := range tokens {
		switch token {
		case "(", ")":
			if token == "(" {
				stack = append(stack, token)
			} else {
				for len(stack) > 0 && stack[len(stack)-1] != "(" {
					top := stack[len(stack)-1]
					stack = stack[:len(stack)-1]
					postfix = append(postfix, top)
				}
				if len(stack) > 0 {
					stack = stack[:len(stack)-1]
				}
			}
		case "+", "-", "*", "/", "^":
			for len(stack) > 0 && precedence[token] <= precedence[stack[len(stack)-1]] && stack[len(stack)-1] != "(" {
				top := stack[len(stack)-1]
				stack = stack[:len(stack)-1]
				postfix = append(postfix, top)
			}
			stack = append(stack, token)
		default: // Operand
			postfix = append(postfix, token)
		}
	}

	for len(stack) > 0 {
		top := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		if top == "(" {
			return nil, fmt.Errorf("invalid expression")
		}
		postfix = append(postfix, top)
	}

	return postfix, nil
}

func calculatePostfix(postfix []string) (int, error) {
	var stack []int

	for _, token := range postfix {
		if isOperator(token) {
			if len(stack) < 2 {
				return 0, fmt.Errorf("invalid expression")
			}
			b, a := stack[len(stack)-1], stack[len(stack)-2]
			stack = stack[:len(stack)-2]

			var result int
			switch token {
			case "+":
				result = a + b
			case "-":
				result = a - b
			case "*":
				result = a * b
			case "/":
				if b == 0 {
					return 0, fmt.Errorf("division by zero")
				}
				result = a / b
			default:
				return 0, fmt.Errorf("invalid operator %s", token)
			}

			stack = append(stack, result)
			continue
		}
		val, err := getNumFromString(token)
		if err != nil {
			return 0, err
		}
		stack = append(stack, val)
	}

	if len(stack) != 1 {
		return 0, fmt.Errorf("invalid expression")
	}

	return stack[0], nil
}
