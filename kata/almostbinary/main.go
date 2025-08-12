package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// Problem: https://hackattic.com/kata/almost_binary

const (
	zeroSymbol = "."
	oneSymbol  = "#"
	zero       = "0"
	one        = "1"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		line := scanner.Text()

		var binaryStrBuilder strings.Builder
		for i := 0; i < len(line); i++ {
			if string(line[i]) == zeroSymbol {
				binaryStrBuilder.WriteString(zero)
				continue
			}

			if string(line[i]) == oneSymbol {
				binaryStrBuilder.WriteString(one)
				continue
			}

			panic(fmt.Errorf("received an invalid character %s", string(line[i])))
		}

		decimalNum, err := strconv.ParseInt(binaryStrBuilder.String(), 2, 64)
		if err != nil {
			panic(err)
		}

		fmt.Println(decimalNum)
	}
}
