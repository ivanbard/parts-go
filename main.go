package main 

import (
	"fmt"
	"os"
	"strings"
)

func main() {
	argsWithoutProg := os.Args[1:]
	fmt.Println(argsWithoutProg)

	partResult := strings.Join(argsWithoutProg, " ")
	fmt.Println(partResult)
}