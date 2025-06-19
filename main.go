package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"slices"
	"strings"

	"golang.org/x/net/html"
)

func main() {

	//homePage := "https://www.procyclingstats.com/race/tour-de-france/2025"
	homePage := "http://localhost:8000/tour-home-page.html"

	resBody, err := fetchPage(homePage)
	if err != nil {
		fmt.Printf("Error fetching page %v\n", err)
		return
	}

	doc, err := html.Parse(bytes.NewReader(resBody))
	if err != nil {
		fmt.Printf("Erro parsing html %v\n", err)
		return
	}

	stagesTitle := findElementWithHtmlContent(doc, "h3", "Stages")
	if stagesTitle == nil {
		fmt.Printf("Unable to find stages title\n")
		return
	}
	stagesTable := findElementByTagName(stagesTitle.Parent, "table")

	tableHeaderNames := []string{"Date", "Stage"}
	ignoreList := []string{"Restday"}
	stages, err := extractData(stagesTable, tableHeaderNames, ignoreList, NewStage)
	if err != nil {
		fmt.Printf("Error extracting stages %v\n", err)
		return
	}

	for _, v := range stages {
		fmt.Printf("Name: %v, Date: %v\n", v.Name, v.Date)
	}

	//ridersPage := "https://www.procyclingstats.com/race/tour-de-france/2025/startlist/alphabetical"
	ridersPage := "http://localhost:8000/riders.html"
	resBody, err = fetchPage(ridersPage)
	if err != nil {
		fmt.Printf("Error fetching page %v\n", err)
		return
	}

	doc, err = html.Parse(bytes.NewReader(resBody))
	if err != nil {
		fmt.Printf("Error parsing html %v\n", err)
		return
	}

	ridersTitle := findElementWithHtmlContent(doc, "h2", "Alphabetical")
	ridersTable := findElementByTagName(ridersTitle.Parent, "table")

	riderTableHeaderName := []string{"Ridername", "Team"}

	riders, err := extractData(ridersTable, riderTableHeaderName, []string{}, NewRider)
	if err != nil {
		fmt.Printf("Error extracting data for riders %v\n", err)
		return
	}

	for _, v := range riders {
		fmt.Printf("Name: %v, Team: %v\n", v.Name, v.Team)
	}
}

type DataItem interface {
	SetField(fieldName string, value string)
}

type Stage struct {
	Name string
	Date string
}

func NewStage() *Stage {
	return &Stage{}
}

func (s *Stage) SetField(fieldName string, value string) {
	switch fieldName {
	case "Date":
		s.Date = value
	case "Stage":
		s.Name = value
	}
}

type Rider struct {
	Name string
	Team string
}

func (r *Rider) SetField(fieldName string, value string) {
	switch fieldName {
	case "Ridername":
		r.Name = value
	case "Team":
		r.Team = value
	}
}

func NewRider() *Rider {
	return &Rider{}
}

func getColumnNumbersForHeaders(input *html.Node, headers []string) map[string]int {
	output := make(map[string]int)
	for _, v := range headers {
		position := 0
		for n := range input.Descendants() {
			if n.Type == html.ElementNode && n.Data == "th" {
				position++
				if n.FirstChild != nil && n.FirstChild.Data == v {
					output[v] = position
					break
				}
			}
		}
	}
	return output
}

func extractData[V DataItem](table *html.Node, fields []string, ignoreList []string, constructor func() V) ([]V, error) {
	thead := findElementByTagName(table, "thead")
	if thead == nil {
		return nil, errors.New("Unable to find table head")
	}
	columns := getColumnNumbersForHeaders(thead, fields)

	tbody := findElementByTagName(table, "tbody")
	if tbody == nil {
		return nil, errors.New("Unable to find table head")
	}

	itemCount := 0
	var items []V = make([]V, 0)
outer:
	for n := range tbody.ChildNodes() {
		if n.Type == html.ElementNode {
			itemCount++
			item := constructor()
			for k, v := range columns {
				val := getValueAtColumn(n, v)
				if val == "" || slices.Contains(ignoreList, val) {
					continue outer
				}
				item.SetField(k, val)
			}
			items = append(items, item)
		}
	}
	return items, nil
}

func getValueAtColumn(input *html.Node, column int) string {
	counter := 0
	var res string
out:
	for n := range input.ChildNodes() {
		if n.Type == html.ElementNode && n.Data == "td" {
			counter++
			if counter == column {
				for m := range n.Descendants() {
					if m.Type == html.TextNode && strings.TrimSpace(m.Data) != "" {
						res = m.Data
						break out
					}
				}
			}
		}
	}
	return res
}

func findElementByTagName(input *html.Node, tagName string) *html.Node {
	var output *html.Node = nil
	for n := range input.Descendants() {
		if n.Type == html.ElementNode && n.Data == tagName {
			output = n
			break
		}
	}
	return output
}

func findElementWithHtmlContent(input *html.Node, element string, htmlContent string) *html.Node {
	var output *html.Node = nil
	for n := range input.Descendants() {
		if n.Type == html.ElementNode && n.Data == element && n.FirstChild.Data == htmlContent {
			output = n
			break
		}
	}
	return output
}

func fetchPage(url string) ([]byte, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Printf("Error generating request %v\n", err)
		return nil, err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Printf("Error making request %v\n", err)
		return nil, err
	}
	defer res.Body.Close()

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Printf("Erro reading response body %v\n", err)
		return nil, err
	}
	return resBody, nil
}
