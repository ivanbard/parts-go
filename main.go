package main

import (
	"fmt"
	"html"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strings"
)

// i will put vehicle info in a Vehicle struct later down the line
type Vehicle struct {
	Make  string
	Year  int
	Model string
}

type EngineOption struct {
	Name string // engine option
	Path string // href tail with engine option + carcode
}

func buildURL(year int, make, model string) string {
	base_url := "https://www.rockauto.com/en/catalog/"
	// make sure inputs are uniform
	make = strings.ToLower(strings.TrimSpace(make))
	model = strings.ToLower(strings.TrimSpace(model))

	// merge together with inputs + commas as per rockauto formatting
	new_url := fmt.Sprintf("%s%s,%d,%s", base_url, make, year, model)
	return new_url
}

func extractEngines(pageContent string) []string {
	decoded := html.UnescapeString(pageContent)

	// find all engine tags
	re := regexp.MustCompile(`"engine"\s*:\s*"([^"]+)"`)
	matches := re.FindAllStringSubmatch(decoded, -1)

	// use map to not have duplicates
	seen := make(map[string]bool)
	var engines []string

	for _, match := range matches {
		if len(match) < 2 {
			continue
		}

		engine := match[1]

		if !seen[engine] {
			seen[engine] = true
			engines = append(engines, engine)
		}
	}

	return engines
}

// TO-DO ONCE STRATEGY IS FIGURED OUT
// should i parse the html content for the href right in this function?
func buildEngineOptions(enginesList []string, currURL string, bodyString string) []EngineOption {
	var options []EngineOption

	// Parse full URL so we can get only:
	// /en/catalog/toyota,2003,camry
	parsedURL, err := url.Parse(currURL)
	if err != nil {
		fmt.Println("Error parsing current URL:", err)
		return options
	}

	basePath := parsedURL.Path

	for _, engine := range enginesList {
		// keep original engine name untouched
		engineSlug := strings.ToLower(engine)
		engineSlug = strings.ReplaceAll(engineSlug, " ", "+")

		// what we expect the href to START with
		searchPrefix := fmt.Sprintf(
			`href="%s,%s,`,
			basePath,
			engineSlug,
		)

		// find searchPrefix inside bodyString
		hrefStart := strings.Index(bodyString, searchPrefix)

		pathStart := hrefStart + len(`href="`)
		// extract everything until the next "
		extractedString := bodyString[pathStart : pathStart+strings.Index(bodyString[pathStart:], `"`)]

		// store it in struct
		options = append(options, EngineOption{
			Name: engine,
			Path: extractedString,
		})

		fmt.Println(searchPrefix) // temporary debugging
	}

	return options
}

func buildFullURL(tag string) string {
	base_url := "https://www.rockauto.com"
	new_url := fmt.Sprintf("%s%s", base_url, tag)
	return new_url
}

func main() {
	var make, model string
	var year, engineChoice int

	// very very basic car specifics prompting
	fmt.Printf("yo you launched parts-go\n")
	fmt.Printf("put in your cars year yea?\n")
	fmt.Scan(&year)
	fmt.Printf("now the make?\n")
	fmt.Scan(&make)
	fmt.Printf("and the model?\n")
	fmt.Scan(&model)

	// build URL with input info
	u := buildURL(year, make, model)
	fmt.Printf("looking for parts for a %d %s %s\n", year, make, model)
	fmt.Printf("%s\n", u)

	// send request to RA for vehicle
	resp, err := http.Get(u)
	if err != nil {
		fmt.Printf("error fetching page: %v\n", err)
		return
	}

	defer resp.Body.Close() // close body when GET function is done

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("error with status code: %d\n", resp.StatusCode)
		return
	}

	//
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("error reading body: %v\n", err)
		return
	}
	bodyString := string(bodyBytes)

	// raw scrape text dump
	fmt.Printf("%s\n", bodyString)

	// print dump of found engines in content scrape from car model
	engines := extractEngines(bodyString)
	for i, engine := range engines {
		fmt.Printf("%d. %s\n", i+1, engine)
	}

	// find hrefs for available engine options

	fmt.Printf("punch in the number of engine choice\n")
	fmt.Scan(&engineChoice)

	choices := buildEngineOptions(engines, u, bodyString)
	fmt.Printf("%s %s\n", choices[engineChoice-1].Name, choices[engineChoice-1].Path)

	fullURL := buildFullURL(choices[engineChoice-1].Path)
	fmt.Printf("%s\n", fullURL)

	//now just need to re-request http scrape to new URL
	resp, err = http.Get(fullURL)
	if err != nil {
		fmt.Printf("error fetching page: %v\n", err)
		return
	}

	defer resp.Body.Close() // close body when GET function is done

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("error with status code: %d\n", resp.StatusCode)
		return
	}

	//
	bodyBytes, err = io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("error reading body: %v\n", err)
		return
	}
	bodyString = string(bodyBytes)

	// raw scrape text dump
	fmt.Printf("%s\n", bodyString)
}
