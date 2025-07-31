package dataextractor

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"le-tour-dashmore-server/model"
	"net/http"
	"slices"
	"strings"

	"golang.org/x/net/html"
)

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

func getStageTables(doc *html.Node, tableName string) ([]*html.Node, error) {
	const dataId = "data-id"
	var err error
	stageTitle := findElementWithHtmlContent(doc, "a", tableName)
	if stageTitle == nil {
		err = errors.New("Unable to find stage table " + tableName)
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

	tables := getElementsByType(stageTable, "table")
	return tables, err
}

func getElementsByType(input *html.Node, elementType string) []*html.Node {
	res := make([]*html.Node, 0)

	for n := range input.Descendants() {
		if n.Type == html.ElementNode && n.Data == elementType {
			res = append(res, n)
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

func getValueAtColumn(input *html.Node, column int, isRider bool) string {
	counter := 0
	var res string
out:
	for n := range input.ChildNodes() {
		if n.Type == html.ElementNode && n.Data == "td" {
			counter++
			if counter == column {
				for m := range n.Descendants() {
					if m.Type == html.TextNode && strings.TrimSpace(m.Data) != "" {
						if isRider {
							res = fmt.Sprintf("%s%s", strings.ToUpper(m.Data), m.Parent.NextSibling.Data)
						} else {
							res = m.Data
						}
						break out
					}
				}
			}
		}
	}
	return res
}

func extractData[V model.DataItem](table *html.Node, fields []string, ignoreList []string, constructor func() V) ([]V, error) {
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

	_, isRidersTable := any(items).([]*model.Rider)
outer:
	for n := range tbody.ChildNodes() {
		if n.Type == html.ElementNode {
			itemCount++
			item := constructor()
			for k, v := range columns {
				// TODO and V is not model.Rider
				val := getValueAtColumn(n, v, k == "Rider" && !isRidersTable)
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

func extractDataFromList(list *html.Node, baseURL string) []model.Team {
	teams := make([]model.Team, 0)
	for n := range list.Descendants() {
		if n.Type == html.ElementNode && n.Data == "a" {
			teamURL := getAttribute(n, "href")
			for m := range n.Descendants() {
				if m.Type == html.TextNode && m.Data != "" {
					teams = append(teams, model.Team{
						Name: m.Data,
						TeamURL: fmt.Sprintf("%s/%s", baseURL, teamURL),
					})
				}
			}
		}
	}
	return teams
}

func addTimedResultToResultsMap(resultsMap map[string]*model.LegacyResult, timedResult model.TimedResult, resultType string) map[string]*model.LegacyResult {
	if resultsMap[timedResult.Rider] == nil {
		resultsMap[timedResult.Rider] = &model.LegacyResult{}
	}
	switch resultType {
	case GC:
		resultsMap[timedResult.Rider].Time = timedResult.TimeStr
		resultsMap[timedResult.Rider].TimeRanking = timedResult.Rank
	case YOUTH:
		resultsMap[timedResult.Rider].YoungRiderRanking = timedResult.Rank
	}
	return resultsMap
}

func addPointsResultToResultsMap(resultsMap map[string]*model.LegacyResult, pointsResult model.PointsResult, resultType string) map[string]*model.LegacyResult {
	if resultsMap[pointsResult.Rider] == nil {
		resultsMap[pointsResult.Rider] = &model.LegacyResult{}
	}
	switch resultType {
	case POINTS:
		resultsMap[pointsResult.Rider].Sprint = pointsResult.Points
		resultsMap[pointsResult.Rider].SprintRanking = pointsResult.Rank
	case KOM:
		resultsMap[pointsResult.Rider].Climber = pointsResult.Points
		resultsMap[pointsResult.Rider].ClimberRanking = pointsResult.Rank
	}
	return resultsMap
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

func getTeamsList(url string, element string, htmlContent string) (*html.Node, error) {
	doc, err := getNodesFromURL(url)
	if err != nil {
		fmt.Printf("Error parsing html %v\n", err)
		return nil, err
	}

	teamsTitle := findElementWithHtmlContent(doc, element, htmlContent)
	if teamsTitle == nil {
		fmt.Printf("Unable to find teams title\n")
		return nil, err
	}

	table := findElementByTagName(teamsTitle.Parent, "ul")
	if table == nil {
		err = errors.New("Unable to find teams list")
	}
	return table, err
}
