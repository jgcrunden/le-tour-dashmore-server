package dataextractor

import (
	"le-tour-dashmore-server/model"
	"testing"
	"time"
)

type MockDB struct {
}

// GetGoalPoints implements dbclient.DBClient.
func (m MockDB) GetGoalPoints(teamID int) (int, error) {
	panic("unimplemented")
}

// GetJerseyRankingForRiderAtStage implements dbclient.DBClient.
func (m MockDB) GetJerseyRankingForRiderAtStage(riderID int, stageID int, goal model.GoalType) (int, error) {
	panic("unimplemented")
}

// UpdateGoalPoints implements dbclient.DBClient.
func (m MockDB) UpdateGoalPoints(teamID int, totalGoalPoints int) error {
	panic("unimplemented")
}

// GetTimedResultsForRider implements dbclient.DBClient.
func (m MockDB) GetTimedResultsForRider(riderID int, jersey model.GoalType) ([]int, error) {
	panic("unimplemented")
}

// GetGoalPointAllocations implements dbclient.DBClient.
func (m MockDB) GetGoalPointAllocations(goal model.GoalType) (map[int]int, error) {
	panic("unimplemented")
}

// CreateFantasyTeamGoal implements dbclient.DBClient.
func (m MockDB) CreateFantasyTeamGoal(id int, goal string) error {
	panic("unimplemented")
}

// UpdateFantasyTeamPoints implements dbclient.DBClient.
func (m MockDB) UpdateFantasyTeamPoints(teamID int, points int) error {
	panic("unimplemented")
}

// GetAllFantasyTeams implements dbclient.DBClient.
func (m MockDB) GetAllFantasyTeams() ([]model.FantasyTeam, error) {
	panic("unimplemented")
}

// GetRiderByName implements dbclient.DBClient.
func (m MockDB) GetRiderByName(riderName string) (model.Rider, error) {
	panic("unimplemented")
}

// AddFantasyRider implements dbclient.DBClient.
func (m MockDB) AddFantasyRider(id int, riderName string) error {
	panic("unimplemented")
}

// CreateFantasyTeam implements dbclient.DBClient.
func (m MockDB) CreateFantasyTeam(teamName string) (int, error) {
	panic("unimplemented")
}

// UpdateRiderPoints implements dbclient.DBClient.
func (m MockDB) UpdateRiderPoints(riderID int, points int) error {
	panic("unimplemented")
}

// GetJerseyRankingPoints implements dbclient.DBClient.
func (m MockDB) GetJerseyRankingPoints(riderID int) (int, error) {
	panic("unimplemented")
}

// GetPointsResultPoints implements dbclient.DBClient.
func (m MockDB) GetPointsResultPoints(riderID int) (int, error) {
	panic("unimplemented")
}

// GetTimedResultPoints implements dbclient.DBClient.
func (m MockDB) GetTimedResultPoints(riderID int) (int, error) {
	panic("unimplemented")
}

// AddSplit implements dbclient.DBClient.
func (m MockDB) AddSplit(splitName string, stageID int) (int, error) {
	panic("unimplemented")
}

// AddTeam implements dbclient.DBClient.
func (m MockDB) AddTeam(team *model.Team) error {
	panic("unimplemented")
}

// AddJerseyRanking implements dbclient.DBClient.
func (m MockDB) AddJerseyRanking(timedResult model.JerseyResult, stageID int, jerseyID int) error {
	panic("unimplemented")
}

// AddPointsResult implements dbclient.DBClient.
func (m MockDB) AddPointsResult(pointsResult *model.PointsResult, stageID int, jerseyID int, splitID int) error {
	panic("unimplemented")
}

// AddTimedResult implements dbclient.DBClient.
func (m MockDB) AddTimedResult(timedResult *model.TimedResult, stageID int, jerseyID int) error {
	panic("unimplemented")
}

// GetJersey implements dbclient.DBClient.
func (m MockDB) GetJersey(name string) (model.Jersey, error) {
	panic("unimplemented")
}

// AddRider implements dbclient.DBClient.
func (m MockDB) AddRider(rider *model.Rider) error {
	panic("unimplemented")
}

// AddStage implements dbclient.DBClient.
func (m MockDB) AddStage(stage *model.Stage) error {
	panic("unimplemented")
}

// GetAllRiders implements dbclient.DBClient.
func (m MockDB) GetAllRiders() ([]model.Rider, error) {
	panic("unimplemented")
}

// GetAllStages implements dbclient.DBClient.
func (m MockDB) GetAllStages() ([]model.Stage, error) {
	panic("unimplemented")
}

func TestGetStageResults(t *testing.T) {
	mockDB := MockDB{}
	config := model.Config{}
	extractor := NewExtractor(mockDB, config)
	stage := model.Stage{
		ID:   2,
		Name: "Stage 2",
		Date: time.Now(),
		URL:  "http://localhost:8000/stage-2.html",
	}
	extractor.GetStageResults(stage)

}
