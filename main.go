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

type TableTitleLocator struct {
	Element   string
	TitleName string
}

func getNodesFromURL(url string) (*html.Node, error) {
	resBody, err := fetchPage(url)
	if err != nil {
		fmt.Printf("Error fetching page %v\n", err)
		return nil, err
	}

	doc, err := html.Parse(bytes.NewReader(resBody))
	if err != nil {
		fmt.Printf("Error parsing html %v\n", err)
		return nil, err
	}
	return doc, err
}

func getRaceDetailsTable(url string, element string, htmlContent string) (*html.Node, error) {
	doc, err := getNodesFromURL(url)
	if err != nil {
		fmt.Printf("Error parsing html %v\n", err)
		return nil, err
	}

	stagesTitle := findElementWithHtmlContent(doc, element, htmlContent)
	if stagesTitle == nil {
		fmt.Printf("Unable to find stages title\n")
		return nil, err
	}
	table := findElementByTagName(stagesTitle.Parent, "table")
	if table == nil {
		err = errors.New("Unable to find table")
	}
	return table, err
}

func getStageTable(doc *html.Node, tableName string) (*html.Node, error) {
	const dataId = "data-id"
	var err error
	stageTitle := findElementWithHtmlContent(doc, "a", tableName)
	if stageTitle == nil {
		err = errors.New("Unable to find stage table")
		return nil, err
	}
	stageTableId := getAttribute(stageTitle, dataId)
	if stageTableId == "" {
		err = errors.New("Error finding stage table id")
		return nil, err
	}

	stageTable := findElementWithAttribute(doc, "div", dataId, stageTableId)
	if stageTable == nil {
		err = errors.New("Error finding stage tage")
	}
	return stageTable, err
}

func main() {
	getRaceDetails()

	getStageResults("http://localhost:8000/stage.html")
}

func getRaceDetails() {
	//homePage := "https://www.procyclingstats.com/race/tour-de-france/2025"
	stagesTable, err := getRaceDetailsTable("http://localhost:8000/tour-home-page.html", "h3", "Stages")
	if err != nil {
		fmt.Println(err)
		return
	}

	_, err = extractData(stagesTable, []string{"Date", "Stage"}, []string{"Restday"}, NewStage)
	if err != nil {
		fmt.Printf("Error extracting stages %v\n", err)
		return
	}

	//ridersPage := "https://www.procyclingstats.com/race/tour-de-france/2025/startlist/alphabetical"
	ridersTable, err := getRaceDetailsTable("http://localhost:8000/riders.html", "h2", "Alphabetical")
	if err != nil {
		fmt.Println(err)
		return
	}

	_, err = extractData(ridersTable, []string{"Ridername", "Team"}, []string{}, NewRider)
	if err != nil {
		fmt.Printf("Error extracting data for riders %v\n", err)
		return
	}
}

func getStageResults(url string) {
	// STAGE //
	doc, err := getNodesFromURL(url)
	if err != nil {
		fmt.Printf("Error parsing html %v\n", err)
		return
	}

	timedTables := []string{"Stage", "GC"}
	for _, v := range timedTables {
		stageTable, err := getStageTable(doc, v)
		if err != nil {
			fmt.Println(err)
			return
		}
		_, err = extractData(stageTable, []string{"Rnk", "Rider", "Time"}, []string{}, NewTimedResult)
		if err != nil {
			fmt.Printf("Error extracting data for riders %v\n", err)
			return
		}
	}

	pointsTables := []string{"Points", "KOM"}
	for _, v := range pointsTables {
		pointsTable, err := getStageTable(doc, v)
		if err != nil {
			fmt.Println(err)
			return
		}

		_, err = extractData(pointsTable, []string{"Rnk", "Rider", "Points"}, []string{}, NewPointsResult)
		if err != nil {
			fmt.Printf("Error extracting data for riders %v\n", err)
			return
		}
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

func NewRider() *Rider {
	return &Rider{}
}

func (r *Rider) SetField(fieldName string, value string) {
	switch fieldName {
	case "Ridername":
		r.Name = value
	case "Team":
		r.Team = value
	}
}

type TimedResult struct {
	Rank  string
	Rider string
	Time  string
}

func NewTimedResult() *TimedResult {
	return &TimedResult{}
}

func (r *TimedResult) SetField(fieldName string, value string) {
	switch fieldName {
	case "Rnk":
		r.Rank = value
	case "Rider":
		r.Rider = value
	case "Time":
		r.Time = value
	}
}

type PointsResult struct {
	Rank   string
	Rider  string
	Points string
}

func NewPointsResult() *PointsResult {
	return &PointsResult{}
}

func (r *PointsResult) SetField(fieldName string, value string) {
	switch fieldName {
	case "Rnk":
		r.Rank = value
	case "Rider":
		r.Rider = value
	case "Points":
		r.Points = value
	}
}

func getAttribute(node *html.Node, attributeKey string) string {
	var res string
	for _, v := range node.Attr {
		if v.Key == attributeKey {
			res = v.Val
			break
		}
	}
	return res
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

func findElementWithAttribute(input *html.Node, element string, attributeKey string, attributeValue string) *html.Node {
	var output *html.Node = nil
	for n := range input.Descendants() {
		if n.Type == html.ElementNode && n.Data == element {
			for _, m := range n.Attr {
				if m.Key == attributeKey && m.Val == attributeValue {
					output = n
					break
				}
			}
		}
	}
	return output
}

func findElementWithHtmlContent(input *html.Node, element string, htmlContent string) *html.Node {
	var output *html.Node = nil
	for n := range input.Descendants() {
		if n.Type == html.ElementNode && n.Data == element {
			for m := range n.Descendants() {
				if m.Data == htmlContent {
					output = n
					break
				}
			}
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
		fmt.Printf("Error reading response body %v\n", err)
		return nil, err
	}
	return resBody, nil
}
