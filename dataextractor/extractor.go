package dataextractor

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"slices"
	"strings"

	"le-tour-dashmore-server/model"

	"golang.org/x/net/html"
)

const (
	year   int = 2025
	GC         = "GC"
	STAGE      = "STAGE"
	POINTS     = "POINTS"
	KOM        = "KOM"
	YOUTH      = "YOUTH"
	RNK        = "Rnk"
	RIDER      = "Rider"
	TIME       = "Time"
	PNT        = "Pnt"
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

func getStageTable(doc *html.Node, tableName string) (*html.Node, error) {
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
	return stageTable, err
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
outer:
	for n := range tbody.ChildNodes() {
		if n.Type == html.ElementNode {
			itemCount++
			item := constructor()
			for k, v := range columns {
				val := getValueAtColumn(n, v, k == "Rider")
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

func addTimedResultToResultsMap(resultsMap map[string]*model.Result, timedResult model.TimedResult, resultType string) map[string]*model.Result {
	if resultsMap[timedResult.Rider] == nil {
		resultsMap[timedResult.Rider] = &model.Result{}
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

func addPointsResultToResultsMap(resultsMap map[string]*model.Result, pointsResult model.PointsResult, resultType string) map[string]*model.Result {
	if resultsMap[pointsResult.Rider] == nil {
		resultsMap[pointsResult.Rider] = &model.Result{}
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

func GetStageResults(url string) (map[string]*model.Result, error) {
	// STAGE //
	doc, err := getNodesFromURL(url)
	if err != nil {
		fmt.Printf("Error parsing html %v\n", err)
		return nil, err
	}

	res := make(map[string]*model.Result)
	timedTables := []string{STAGE, GC, YOUTH}
	for _, v := range timedTables {
		stageTable, err := getStageTable(doc, v)
		if err != nil {
			fmt.Println(err)
			return nil, err
		}
		timedResults, err := extractData(stageTable, []string{RNK, RIDER, TIME}, []string{}, model.NewTimedResult)
		if err != nil {
			fmt.Printf("Error extracting data for riders %v\n", err)
			return nil, err
		}

		previousTime := ""
		if err := timedResults[0].ParseTime(0); err != nil {
			fmt.Println(err)
			return nil, err
		}
		leadersTime := timedResults[0].Time
		for _, timedResult := range timedResults {
			previousTime = timedResult.FixDitto(previousTime)
			timedResult.ParseTime(leadersTime)
			res = addTimedResultToResultsMap(res, *timedResult, v)
		}

	}

	pointsTables := []string{POINTS, KOM}
	for _, v := range pointsTables {
		pointsTable, err := getStageTable(doc, v)
		if err != nil {
			fmt.Println(err)
			return nil, err
		}

		pointsResults, err := extractData(pointsTable, []string{RNK, RIDER, PNT}, []string{}, model.NewPointsResult)
		if err != nil {
			fmt.Printf("Error extracting data for riders %v\n", err)
			return nil, err
		}

		for _, pointsResult := range pointsResults {
			res = addPointsResultToResultsMap(res, *pointsResult, v)
		}
	}
	return res, nil
}

func GetRaceDetails() {
	//homePage := "https://www.procyclingstats.com/race/tour-de-france/2025"
	stagesTable, err := getRaceDetailsTable("http://localhost:8000/tour-home-page.html", "h3", "Stages")
	if err != nil {
		fmt.Println(err)
		return
	}

	_, err = extractData(stagesTable, []string{"Date", "", "Stage"}, []string{"Restday"}, model.NewStage)
	if err != nil {
		fmt.Printf("Error extracting stages %v\n", err)
		return
	}

	//ridersPage := "https://www.procyclingstats.com/race/tour-de-france/2025/startlist/alphabetical"
	/*
		ridersTable, err := getRaceDetailsTable("http://localhost:8000/top-competitors.html", "h2", "Top competitors")
		if err != nil {
			fmt.Println(err)
			return
		}

		_, err = extractData(ridersTable, []string{"Rider", "Team", "Points"}, []string{}, NewRider)
		if err != nil {
			fmt.Printf("Error extracting data for riders %v\n", err)
			return
		}

	*/

	//getAllStagesFromDatabase()
	/*
		for _, v := range riders {
			err := addRiderToDatabase(v)
			if err != nil {
				fmt.Println(err)
				return
			}
		}
	*/

}
