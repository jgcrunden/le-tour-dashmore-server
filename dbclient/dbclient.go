package dbclient

import "le-tour-dashmore-server/model"

type DBClient interface {
	AddRider(rider *model.Rider) error
	GetAllRiders() ([]model.Rider, error)
	GetRiderByName(riderName string) (model.Rider, error)
	UpdateRiderPoints(riderID int, points int) error

	AddTeam(team *model.Team) error

	AddStage(stage *model.Stage) error
	GetAllStages() ([]model.Stage, error)

	AddSplit(splitName string, stageID int) (int, error)

	GetJersey(name string) (model.Jersey, error)

	AddTimedResult(timedResult *model.TimedResult, stageID int, jerseyID int) error
	GetTimedResultPoints(riderID int) (int, error)
	GetTimedResultsForRider(riderID int, jersey model.GoalType) ([]int, error)

	AddJerseyRanking(timedResult model.JerseyResult, stageID int, jerseyID int) error
	GetJerseyRankingPoints(riderID int) (int, error)
	GetJerseyRankingForRiderAtStage(riderID int, stageID int, goal model.GoalType) (int, error)

	AddPointsResult(pointsResult *model.PointsResult, stageID int, jerseyID int, splitID int) error
	GetPointsResultPoints(riderID int) (int, error)

	GetGoalPoints(teamID int) (int, error)
	UpdateGoalPoints(teamID int, totalGoalPoints int) error

	GetGoalPointAllocations(goal model.GoalType) (map[int]int, error)

	CreateFantasyTeam(teamName string) (int, error)
	GetAllFantasyTeams() ([]model.FantasyTeam, error)
	UpdateFantasyTeamPoints(teamID int, points int) error

	AddFantasyRider(id int, riderName string) error

	CreateFantasyTeamGoal(id int, goal string) error
}

