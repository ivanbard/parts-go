package main

import (
	"fmt"
	"strings"
)

func buildURL(year int, make, model string) string {
	base_url := "https://www.rockauto.com/en/catalog/"
	// make sure inputs are uniform
	make = strings.ToLower(strings.TrimSpace(make))
	model = strings.ToLower(strings.TrimSpace(model))

	// merge together with inputs + commas as per rockauto formatting
	new_url := fmt.Sprintf("%s%s,%d,%s", base_url, make, year, model)
	return new_url
}

func main() {
	var make, model string
	var year int

	// very very basic car specifics prompting
	fmt.Printf("yo you launched parts-go\n")
	fmt.Printf("put in your cars year yea?\n")
	fmt.Scan(&year)
	fmt.Printf("now the make?\n")
	fmt.Scan(&make)
	fmt.Printf("and the model?\n")
	fmt.Scan(&model)

	u := buildURL(year, make, model)
	fmt.Printf("looking for parts for a %d %s %s\n", year, make, model)
	fmt.Printf("%s", u)
}
