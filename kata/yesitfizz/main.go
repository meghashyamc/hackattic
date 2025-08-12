package main

import (
	"bufio"
	"fmt"
	"log/slog"
	"os"
	"strconv"
	"strings"
)

// Problem: https://hackattic.com/kata/yes_it_fizz

const (
	fizz = "Fizz"
	buzz = "Buzz"
)

func main() {

	scanner := bufio.NewScanner(os.Stdin)

	scanner.Scan()
	line := scanner.Text()
	parts := strings.Split(line, " ")
	num1, err := strconv.Atoi(parts[0])
	if err != nil {
		slog.Error("got an unexpected error when converting string to int", "err", err, "str", num1)
		os.Exit(1)
	}
	num2, err := strconv.Atoi(parts[1])
	if err != nil {
		slog.Error("got an unexpected error when converting string to int", "err", err, "str", num2)
		os.Exit(1)
	}

	if num1 > num2 {
		slog.Error("first number must be less than second number", "num1", num1, "num2", num2)
		os.Exit(1)
	}

	printFizzBuzz(num1, num2)

}

func printFizzBuzz(start int, end int) {
	for i := start; i <= end; i++ {
		if i%3 == 0 && i%5 == 0 {
			fmt.Println(fizz + buzz)
			continue
		}
		if i%3 == 0 {
			fmt.Println(fizz)
			continue
		}
		if i%5 == 0 {
			fmt.Println(buzz)
			continue
		}

		fmt.Println(i)
	}
}
