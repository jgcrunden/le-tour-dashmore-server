package dataextractor

import (
	"errors"
	"fmt"

	"le-tour-dashmore-server/dbclient"
	"le-tour-dashmore-server/model"

	"golang.org/x/net/html"
)

const (
	GC     = "GC"
	STAGE  = "STAGE"
	POINTS = "POINTS"
	KOM    = "KOM"
	YOUTH  = "YOUTH"
	RNK    = "Rnk"
	RIDER  = "Rider"
	TIME   = "Time"
	PNT    = "Pnt"
)

type Extractor struct {
	client dbclient.DBClient
	config model.Config
}

func NewExtractor(dbClient dbclient.DBClient, config model.Config) Extractor {
	return Extractor{
		client: dbClient,
		config: config,
	}
}

func (e Extractor) GetStageResults(stage model.Stage) (map[string]*model.LegacyResult, error) {
	fmt.Println("Getting results for", stage.Name)
	// STAGE //
	doc, err := getNodesFromURL(stage.URL)
	if err != nil {
		fmt.Printf("Error parsing html %v\n", err)
		return nil, err
	}

	res := make(map[string]*model.LegacyResult)

	res, err = e.processTimedTables(doc, stage, res)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return e.processPointsTables(doc, stage, res)
}

func (e Extractor) processTimedTables(doc *html.Node, stage model.Stage, res map[string]*model.LegacyResult) (map[string]*model.LegacyResult, error) {
	timedTables := []string{STAGE, GC, YOUTH}
	for _, v := range timedTables {
		fmt.Println("Getting", v)
		stageTables, err := getStageTables(doc, v)
		if err != nil {
			fmt.Println(err)
			return nil, err
		}

		jersey, err := e.client.GetJersey(v)
		if err != nil {
			fmt.Println("Error getting jersey", err)
		}

		for i, stageTable := range stageTables {
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
			leader := timedResults[0]
			for _, timedResult := range timedResults {
				previousTime = timedResult.FixDitto(previousTime, leader.TimeStr)
				timedResult.ParseTime(leader.Time)
				res = addTimedResultToResultsMap(res, *timedResult, v)
				if v == STAGE || (v == YOUTH && i == 1) {
					// add to stage results table
					err := e.client.AddTimedResult(timedResult, stage.ID, jersey.ID)
					if err != nil {
						return nil, err
					}
				} else {
					// add to jersey rank table
					err := e.client.AddJerseyRanking(timedResult, stage.ID, jersey.ID)
					if err != nil {
						return nil, err
					}
				}
			}
		}
	}
	return res, nil
}

func (e Extractor) processPointsTables(doc *html.Node, stage model.Stage, res map[string]*model.LegacyResult) (map[string]*model.LegacyResult, error) {
	pointsTables := []string{POINTS, KOM}
	for _, v := range pointsTables {
		fmt.Println("Getting", v)
		pointsTables, err := getStageTables(doc, v)
		if err != nil {
			fmt.Println(err)
			return nil, err
		}

		jersey, err := e.client.GetJersey(v)
		if err != nil {
			fmt.Println("Error getting jersey", e)
			return nil, err
		}

		for i, pointsTable := range pointsTables {
			splitID := 0
			if i > 0 {
				// add to stage results table
				splitName := pointsTable.PrevSibling.PrevSibling.FirstChild.Data
				splitID, err = e.client.AddSplit(splitName, stage.ID)
				if err != nil {
					fmt.Println("Error adding split", splitName)
					return nil, err
				}
			}
			pointsResults, err := extractData(pointsTable, []string{RNK, RIDER, PNT}, []string{}, model.NewPointsResult)
			if err != nil {
				fmt.Printf("Error extracting data for riders %v\n", err)
				return nil, err
			}

			for _, pointsResult := range pointsResults {
				res = addPointsResultToResultsMap(res, *pointsResult, v)
				if i > 0 {
					err = e.client.AddPointsResult(pointsResult, stage.ID, jersey.ID, splitID)
					if err != nil {
						return nil, err
					}
				} else {
					err = e.client.AddJerseyRanking(pointsResult, stage.ID, jersey.ID)
					if err != nil {
						return nil, err
					}
					// add to classification results table
				}
			}
		}
	}
	return res, nil
}

func (e Extractor) GetRaceDetails() {
	fmt.Println("Getting race details")

	err := e.getStages()
	if err != nil {
		fmt.Println(err)
		return
	}

	err = e.getTeams()
	if err != nil {
		fmt.Println(err)
		return
	}

	err = e.getRiders()
	if err != nil {
		fmt.Println(err)
		return
	}
}

func (e Extractor) CalculatePoints() error {

	riders, err := e.client.GetAllRiders()
	if err != nil {
		fmt.Println("Unable to get riders", err)
		return err
	}

	for _, rider := range riders {
		timedResultPoints, err := e.client.GetTimedResultPoints(rider.ID)
		if err != nil {
			fmt.Println("Error getting timed results points", err)
			return err
		}
		pointsResultPoints, err := e.client.GetPointsResultPoints(rider.ID)
		if err != nil {
			fmt.Println("Error getting jersey ranking points", err)
			return err
		}
		jerseyResultPoints, err := e.client.GetJerseyRankingPoints(rider.ID)
		if err != nil {
			fmt.Println("Error getting jersey ranking points", err)
			return err
		}

		total := timedResultPoints + pointsResultPoints + jerseyResultPoints
		err = e.client.UpdateRiderPoints(rider.ID, total)
		if err != nil {
			fmt.Println(err)
			return err
		}
	}
	return nil
}

func (e Extractor) getStages() error {
	fmt.Println("Getting stages")
	stagesTable, err := getRaceDetailsTable(e.config.TourExtension, "h4", "Stages")
	if err != nil {
		fmt.Println("Error getting race details table", err)
		return err
	}

	stages, err := extractData(stagesTable, []string{"Date", "", "Stage"}, []string{"Restday"}, model.NewStage)
	for i := range len(stages) {
		stages[i].URL = fmt.Sprintf("%s/stage-%d", e.config.TourExtension, i+1)
	}
	if err != nil {
		fmt.Printf("Error extracting stages %v\n", err)
		return err
	}

	for _, stage := range stages {
		err = e.client.AddStage(stage)
		if err != nil {
			fmt.Println("Error adding stage", err)
			return err
		}
	}
	return nil
}

func (e Extractor) getTeams() error {
	fmt.Println("Getting teams")
	teamsList, err := getTeamsList(e.config.TourExtension, "h4", "Teams")
	if err != nil {
		fmt.Println("Error getting teams list", err)
		return err
	}

	teams := extractDataFromList(teamsList, e.config.CyclingStatsBaseURL)
	if teams == nil {
		return errors.New("Error extract teams from list")
	}

	for _, v := range teams {
		img, err := e.getJerseyImage(v.TeamURL, "div", "Jersey: ")
		if err != nil {
			fmt.Println(err)
			return err
		}
		v.JerseyImage = img
		err = e.client.AddTeam(&v)
		if err != nil {
			fmt.Println("Error adding team", v.Name)
			return err
		}
	}
	return nil
}

func (e Extractor) getRiders() error {
	fmt.Println("Getting riders")
	ridersTable, err := getRaceDetailsTable(fmt.Sprintf("%s/startlist/top-competitors", e.config.TourExtension), "h2", "Top competitors")
	if err != nil {
		fmt.Println("Error getting riders table", err)
		return err
	}

	riders, err := extractData(ridersTable, []string{"Rider", "Team", "Points"}, []string{}, model.NewRider)
	if err != nil {
		fmt.Printf("Error extracting data for riders %v\n", err)
		return err
	}

	for _, v := range riders {
		err := e.client.AddRider(v)
		if err != nil {
			fmt.Println("Error adding rider", v.Name)
			return err
		}
	}
	return nil
}

func (e Extractor) getJerseyImage(url string, element string, htmlContent string) ([]byte, error) {
	doc, err := getNodesFromURL(url)
	if err != nil {
		fmt.Printf("Error parsing html %v\n", err)
		return nil, err
	}

	jerseyImageSection := findElementWithHtmlContent(doc, element, htmlContent)
	if jerseyImageSection == nil {
		fmt.Printf("Unable to find teams jersey\n")
		return nil, err
	}

	jerseyImageTag := findElementByTagName(jerseyImageSection.Parent, "img")
	if jerseyImageTag == nil {
		err = errors.New("Unable to find teams list")
		return nil, err
	}
	jerseyImageLink := getAttribute(jerseyImageTag, "src")
	if jerseyImageLink == "" {
		return nil, errors.New("Unable to find jersey image")
	}

	return fetchPage(fmt.Sprintf("%s/%s", e.config.CyclingStatsBaseURL, jerseyImageLink))
}
