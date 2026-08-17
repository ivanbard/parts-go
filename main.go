package main 

import (
	"fmt"
	"os"
	"strings"
	"net/http"
	"log"
	"io"
)

func main() {
	argsWithoutProg := os.Args[1:]
	fmt.Println(argsWithoutProg)

	partResult := strings.Join(argsWithoutProg, " ")
	fmt.Println(partResult)

	// get request on rockauto
	resp, err := http.Get("https://www.rockauto.com/")
	if err != nil {
		log.Fatalln(err)
	}

	// read GET request
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	
	sb := string(body)
	log.Printf(sb)
}